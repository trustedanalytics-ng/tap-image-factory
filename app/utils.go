/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
	ApplicationId		string	`json:"id"`
	State			string	`json:"state"`
}

type BuildImagePostRequest struct {
	ApplicationId    	string	`json:"id"`
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