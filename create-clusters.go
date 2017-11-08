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

func create(cli *cli.Context) {

	cnf := clientcmd.GetConfigFromFileOrDie(kubeConfigPath)
	c := *cnf
	override := clientcmd.ConfigOverrides{
		ClusterInfo: api.Cluster{
			Server: "",
		},
	}
	for contextName := range c.Contexts {
		logger := log.WithFields(log.Fields{
			"context_name": contextName,
		})
		logger.Info("Found context")
		config := clientcmd.NewNonInteractiveClientConfig(c, contextName, &override, nil)
		clientCnf, e := config.ClientConfig()

		if e != nil {
			msg := fmt.Sprintf("Failed to create config with error:\n%s\n\n", e)
			logger.Warn(msg)
			continue
		}
		logger.Info("Created config for context")

		clientset, e := kubernetes.NewForConfig(clientCnf)
		if e != nil {
			msg := fmt.Sprintf("Failed to create kubernetes client with error:\n%s\n\n", e)
			logger.Warn(msg)
			continue
		}
		logger.Info("Created client set for context")

		sa, e := clientset.CoreV1().ServiceAccounts("default").Get("default", metav1.GetOptions{})
		if e != nil {
			msg := fmt.Sprintf("Failed to get service account token with error:\n%s\n\n", e)
			logger.Warn(msg)
			continue
		}
		secretName := string(sa.Secrets[0].Name)
		namespace := sa.Secrets[0].Namespace
		logger.Info(fmt.Sprintf("Found service account accisiated with secret: %s in namespace %s\n", secretName, namespace))

		secret, e := clientset.CoreV1().Secrets("default").Get(secretName, metav1.GetOptions{})
		if e != nil {
			msg := fmt.Sprintf("Failed to get secrets with error:\n%s\n\n", e)
			logger.Warn(msg)
			continue
		}
		logger.Info(fmt.Sprintf("Found secret"))

		logger.Info(fmt.Sprintf("Creating cluster in Codefresh"))
		body, e := addCluser(clientCnf.Host, contextName, secret.Data["token"], secret.Data["ca.crt"])
		if e != nil {
			msg := fmt.Sprintf("Failed to add cluster with error:\n%s\n\n", e)
			logger.Warn(msg)
			continue
		}
		logger.Warn(fmt.Sprintf("Added!\n%s\n\n", string(body)))
	}
	log.Info("Operation is done, check your account setting")
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
