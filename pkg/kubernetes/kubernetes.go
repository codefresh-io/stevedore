package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/codefresh-io/stevedore/pkg/codefresh"
	"github.com/codefresh-io/stevedore/pkg/reporter"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubeConfig "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type (
	API interface {
		GoOverAllContexts(context.Context)
		GoOverContextByName(context.Context, string, string, string, bool, string)
		GoOverCurrentContext(context.Context)
	}

	kubernetes struct {
		config    *api.Config
		codefresh codefresh.API
		reporter  reporter.Reporter
	}
)

func getDefaultOverride() clientcmd.ConfigOverrides {
	return clientcmd.ConfigOverrides{
		ClusterInfo: api.Cluster{
			Server: "",
		},
	}
}

type getOverContextOptions struct {
	contextName    string
	namespace      string
	serviceaccount string
	config         clientcmd.ClientConfig
	logger         *log.Entry
	codefresh      codefresh.API
	reporter       reporter.Reporter
	behindFirewall bool
	name           string
}

func goOverContext(ctx context.Context, options *getOverContextOptions) error {
	var host string
	var ca []byte
	var token []byte
	clientCnf, e := options.config.ClientConfig()
	if e != nil {
		message := fmt.Sprintf("Failed to create config with error:\n%s", e)
		options.logger.Warn(message)
		clientCnf, e = rest.InClusterConfig()
		if e != nil {
			message = fmt.Sprintf("Failed to create in cluster config with error:\n%s", e)
			options.logger.Warn(message)
			return e
		}
	}
	options.logger.Info("Created config for context")
	host = clientCnf.Host

	options.logger.Info("Creating rest client")
	clientset, e := kubeConfig.NewForConfig(clientCnf)
	if e != nil {
		message := fmt.Sprintf("Failed to create kubernetes client with error:\n%s", e)
		options.logger.Warn(message)

		return e
	}
	options.logger.Info("Created client set for context")

	options.logger.Info("Generating service account secret")

	secret, err := getServiceAccountTokenSecret(ctx, clientset, options)
	if err != nil {
		message := fmt.Sprintf("Failed to generate service account secret: error:\n%s", err.Error())
		options.logger.Error(message)
		return err
	}

	token = secret.Data["token"]
	ca = secret.Data["ca.crt"]
	options.logger.Info(fmt.Sprint("Found secret"))

	options.logger.Info(fmt.Sprint("Creating cluster in Codefresh"))
	result, e := options.codefresh.Create(host, options.name, token, ca, options.behindFirewall)
	if e != nil {
		message := fmt.Sprintf("Failed to add cluster with error:\n%s", e)
		options.logger.Error(message)
		return e
	}
	options.reporter.AddToReport(options.contextName, reporter.SUCCESS, string(result))
	options.logger.Info(fmt.Sprint("Cluster added!"))
	return nil
}

func getServiceAccountTokenSecret(ctx context.Context, client kubeConfig.Interface, options *getOverContextOptions) (*v1.Secret, error) {
	saName := options.serviceaccount
	namespace := options.namespace
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-token-", saName),
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": saName,
			},
		},
		Type: v1.SecretTypeServiceAccountToken,
	}

	options.logger.Debug("Creating secret for service-account token", "service-account", saName)

	secret, err := client.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create service-account token secret: %w", err)
	}
	secretName := secret.Name

	options.logger.Debug("Created secret for service-account token", "service-account", saName, "secret", secret.Name)

	patch := []byte(fmt.Sprintf("{\"secrets\": [{\"name\": \"%s\"}]}", secretName))
	_, err = client.CoreV1().ServiceAccounts(namespace).Patch(ctx, saName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to patch service-account with new secret: %w", err)
	}

	options.logger.Debug("Added secret to service-account secrets", "service-account", saName, "secret", secret.Name)

	// try to read the token from the secret
	ticker := time.NewTicker(time.Second)
	retries := 15
	defer ticker.Stop()

	for try := 0; try < retries; try++ {
		select {
		case <-ticker.C:
			secret, err = client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		options.logger.Debug("Checking secret for service-account token", "service-account", saName, "secret", secret.Name)

		if err != nil {
			return nil, fmt.Errorf("failed to get service-account secret: %w", err)
		}

		if secret.Data == nil || len(secret.Data["token"]) == 0 {
			options.logger.Debug("Secret is missing service-account token", "service-account", saName, "secret", secret.Name)
			continue
		}

		options.logger.Debug("Got service-account token from secret", "service-account", saName, "secret", secret.Name)

		return secret, nil
	}

	return nil, fmt.Errorf("timed out waiting for secret to contain token")
}

func (kube *kubernetes) GoOverAllContexts(ctx context.Context) {
	contexts := kube.config.Contexts
	for contextName := range contexts {
		logger := log.WithFields(log.Fields{
			"context_name": contextName,
		})
		logger.Info("Working on context")
		logger.Info("Creating config")
		override := getDefaultOverride()
		config := clientcmd.NewNonInteractiveClientConfig(*kube.config, contextName, &override, nil)
		options := &getOverContextOptions{
			contextName:    contextName,
			config:         config,
			logger:         logger,
			codefresh:      kube.codefresh,
			reporter:       kube.reporter,
			behindFirewall: false,
			name:           contextName,
		}
		err := goOverContext(ctx, options)
		if err != nil {
			kube.reporter.AddToReport(contextName, reporter.FAILED, err.Error())
			continue
		}
	}
}

func (kube *kubernetes) GoOverContextByName(ctx context.Context, contextName string, namespace string, serviceaccount string, bf bool, name string) {
	var override clientcmd.ConfigOverrides
	var config clientcmd.ClientConfig
	override = getDefaultOverride()
	config = clientcmd.NewNonInteractiveClientConfig(*kube.config, contextName, &override, nil)
	logger := log.WithFields(log.Fields{
		"context_name":    contextName,
		"namespace":       namespace,
		"serviceaccount":  serviceaccount,
		"behind_firewall": bf,
		"name":            name,
	})
	options := &getOverContextOptions{
		contextName:    contextName,
		config:         config,
		logger:         logger,
		codefresh:      kube.codefresh,
		reporter:       kube.reporter,
		namespace:      namespace,
		serviceaccount: serviceaccount,
		behindFirewall: bf,
		name:           name,
	}
	err := goOverContext(ctx, options)
	if err != nil {
		kube.reporter.AddToReport(contextName, reporter.FAILED, err.Error())
	}
}

func (kube *kubernetes) GoOverCurrentContext(ctx context.Context) {
	override := getDefaultOverride()
	config := clientcmd.NewDefaultClientConfig(*kube.config, &override)
	rawConfig, err := config.RawConfig()
	if err != nil {
		kube.reporter.AddToReport("current-context", reporter.FAILED, err.Error())
	}
	contextName := rawConfig.CurrentContext
	logger := log.WithFields(log.Fields{
		"context_name": contextName,
	})
	options := &getOverContextOptions{
		contextName:    contextName,
		config:         config,
		logger:         logger,
		codefresh:      kube.codefresh,
		reporter:       kube.reporter,
		behindFirewall: false,
		name:           contextName,
	}
	err = goOverContext(ctx, options)
	if err != nil {
		kube.reporter.AddToReport(contextName, reporter.FAILED, err.Error())
	}
}

func NewKubernetesAPI(kubeConfigPath string, codefresh codefresh.API, reporter reporter.Reporter) API {
	return &kubernetes{
		config:    clientcmd.GetConfigFromFileOrDie(kubeConfigPath),
		codefresh: codefresh,
		reporter:  reporter,
	}
}
