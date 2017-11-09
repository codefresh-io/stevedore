package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var baseCodefreshURL = "https://g.codefresh.io/"

func getConfigOrDie() *api.Config {
	return clientcmd.GetConfigFromFileOrDie(kubeConfigPath)
}

func getDefaultOverride() clientcmd.ConfigOverrides {
	return clientcmd.ConfigOverrides{
		ClusterInfo: api.Cluster{
			Server: "",
		},
	}
}
func create(cli *cli.Context) {
	cnf := getConfigOrDie()
	c := *cnf
	if runOnAllContexts {
		goOverAllContexts(c)
	} else if runOnContext != "" {
		getOverContextByName(c, runOnContext)
	} else {
		goOverCurrentContext(c)
	}
	reportResult()
	log.Info("Operation is done, check your account setting")
}

func getLogger(name string) *log.Entry {
	return log.WithFields(log.Fields{
		"context_name": name,
	})
}

func goOverAllContexts(cnf api.Config) {
	contexts := cnf.Contexts
	for contextName := range contexts {
		logger := getLogger(contextName)
		logger.Info("Working on context")
		logger.Info("Creating config")
		override := getDefaultOverride()
		config := clientcmd.NewNonInteractiveClientConfig(cnf, contextName, &override, nil)
		err := goOverContext(contextName, config, logger)
		if err != nil {
			addClusterToFinalReport(contextName, failed, err.Error())
			continue
		}
	}

}

func getOverContextByName(cnf api.Config, contextName string) {
	override := getDefaultOverride()
	config := clientcmd.NewNonInteractiveClientConfig(cnf, contextName, &override, nil)
	logger := getLogger(contextName)
	err := goOverContext(contextName, config, logger)
	if err != nil {
		addClusterToFinalReport(contextName, failed, err.Error())
	}
}

func goOverCurrentContext(cnf api.Config) {
	override := getDefaultOverride()
	config := clientcmd.NewDefaultClientConfig(cnf, &override)
	rawConfig, err := config.RawConfig()
	if err != nil {
		addClusterToFinalReport("current-context", failed, err.Error())
	}
	contextName := rawConfig.CurrentContext
	logger := getLogger(contextName)

	err = goOverContext(contextName, config, logger)
	if err != nil {
		addClusterToFinalReport(contextName, failed, err.Error())
	}
}

func goOverContext(contextName string, config clientcmd.ClientConfig, logger *log.Entry) error {
	clientCnf, e := config.ClientConfig()
	if e != nil {
		message := fmt.Sprintf("Failed to create config with error:\n%s", e)
		logger.Warn(message)
		return e
	}
	logger.Info("Created config for context")

	logger.Info("Creating rest client")
	clientset, e := kubernetes.NewForConfig(clientCnf)
	if e != nil {
		message := fmt.Sprintf("Failed to create kubernetes client with error:\n%s", e)
		logger.Warn(message)
		return e
	}
	logger.Info("Created client set for context")

	logger.Info("Fetching service account from cluster")
	sa, e := clientset.CoreV1().ServiceAccounts("default").Get("default", metav1.GetOptions{})
	if e != nil {
		message := fmt.Sprintf("Failed to get service account token with error:\n%s", e)
		logger.Warn(message)
		return e
	}
	secretName := string(sa.Secrets[0].Name)
	namespace := sa.Namespace
	logger.WithFields(log.Fields{
		"secret_name": secretName,
		"namespace":   namespace,
	}).Info(fmt.Sprint("Found service account accisiated with secret"))

	logger.Info("Fetching secret from cluster")
	secret, e := clientset.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if e != nil {
		message := fmt.Sprintf("Failed to get secrets with error:\n%s", e)
		logger.Warn(message)
		return e
	}
	logger.Info(fmt.Sprint("Found secret"))

	logger.Info(fmt.Sprint("Creating cluster in Codefresh"))
	result, e := addCluser(clientCnf.Host, contextName, secret.Data["token"], secret.Data["ca.crt"])
	if e != nil {
		message := fmt.Sprintf("Failed to add cluster with error:\n%s", e)
		logger.Error(message)
		return e
	}
	addClusterToFinalReport(contextName, success, string(result))
	logger.Warn(fmt.Sprint("Cluster added!"))
	return nil
}

func addCluser(host string, contextName string, token []byte, crt []byte) ([]byte, error) {
	url := baseCodefreshURL + "api/clusters/local/cluster"
	payload := &requestPayload{
		Type:                "sat",
		ProviderAgent:       "custom",
		Host:                host,
		Selector:            contextName,
		ServiceAccountToken: token,
		ClientCa:            crt,
	}
	mar, _ := json.Marshal(payload)
	p := strings.NewReader(string(mar))
	err := testConnection(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, p)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Access-Token", codefreshJwt)
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 201 {
		err := errors.New(string(body))
		return nil, fmt.Errorf("Failed to create cluster %s", err)
	}
	return body, nil
}

func testConnection(payload *requestPayload) error {
	url := baseCodefreshURL + "api/kubernetes/test"
	mar, _ := json.Marshal(payload)
	p := strings.NewReader(string(mar))
	req, err := http.NewRequest("POST", url, p)
	if err != nil {
		return err
	}
	req.Header.Add("X-Access-Token", codefreshJwt)
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Failed to test cluster")
	}
	return nil
}

type requestPayload struct {
	Type                string `json:"type"`
	ClientCa            []byte `json:"clientCa"`
	ProviderAgent       string `json:"providerAgent"`
	Selector            string `json:"selector"`
	ServiceAccountToken []byte `json:"serviceAccountToken"`
	Host                string `json:"host"`
}
