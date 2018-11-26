package codefresh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	API interface {
		Test(*requestPayload) error
		Create(string, string, []byte, []byte) ([]byte, error)
	}

	codefreshAPI struct {
		baseURL string
		token   string
	}

	requestPayload struct {
		Type                string `json:"type"`
		ClientCa            []byte `json:"clientCa"`
		ProviderAgent       string `json:"providerAgent"`
		Selector            string `json:"selector"`
		ServiceAccountToken []byte `json:"serviceAccountToken"`
		Host                string `json:"host"`
	}
)

func (api *codefreshAPI) Test(payload *requestPayload) error {
	url := api.baseURL + "api/kubernetes/test"
	mar, _ := json.Marshal(payload)
	p := strings.NewReader(string(mar))
	req, err := http.NewRequest("POST", url, p)
	if err != nil {
		return err
	}
	req.Header.Add("authorization", api.token)
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

func (api *codefreshAPI) Create(host string, name string, saToken []byte, crt []byte) ([]byte, error) {
	url := api.baseURL + "api/clusters/local/cluster"
	payload := &requestPayload{
		Type:                "sat",
		ProviderAgent:       "custom",
		Host:                host,
		Selector:            name,
		ServiceAccountToken: saToken,
		ClientCa:            crt,
	}
	mar, _ := json.Marshal(payload)
	p := strings.NewReader(string(mar))
	err := api.Test(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, p)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", api.token)
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

func NewCodefreshAPI(baseUrl string, token string) API {
	return &codefreshAPI{
		baseURL: baseUrl,
		token:   token,
	}
}
