package main

import (
	"net/http"
	"time"
	"os"
	"strings"
)

type ApplicationGetResponse struct {
	ApplicationId		string	`json:"id"`
	TemplateId		string	`json:"templateId"`
	BaseImage		string	`json:"image"`
	Replication		string	`json:"replication"`
	Type			string	`json:"type"`
	State			string	`json:"state"`
}

type ApplicationStatePutRequest struct {
	ApplicationId		string	`json:"application_id"`
	State			string	`json:"state"`
}

type BuildImagePostRequest struct {
	ApplicationId    	string	`json:"application_id"`
}

func CurrentEnv() map[string]string {
	return Env(os.Environ())
}

func Env(env []string) map[string]string {
	vars := mapEnv(env, splitEnv())
	return vars
}

func splitEnv() func(item string) (key, val string) {
	return func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = strings.Join(splits[1:], "=")
		return
	}
}

func mapEnv(data []string, keyFunc func(item string) (key, val string)) map[string]string {
	items := make(map[string]string)
	for _, item := range data {
		key, val := keyFunc(item)
		items[key] = val
	}
	return items
}

func GetCatalogAddress() string {
	return CurrentEnv()["CATALOG_ADDRESS"]
}

type CatalogConnector struct {
	Server     string
	Client     *http.Client
}

func NewCatalogConnector() *CatalogConnector {
	transport := &http.Transport{}
	clientCreator := &http.Client{Transport: transport, Timeout: time.Duration(30 * time.Minute)}

	return &CatalogConnector{
		Server: GetCatalogAddress(),
		Client: clientCreator,

	}
}