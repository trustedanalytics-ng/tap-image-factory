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
	"bytes"
	"io"
	"net/http"
	"os"
	"time"
)

type ApplicationGetResponse struct {
	ApplicationId string `json:"id"`
	TemplateId    string `json:"templateId"`
	BaseImage     string `json:"image"`
	Replication   string `json:"replication"`
	Type          string `json:"type"`
	State         string `json:"state"`
}

type ApplicationStatePutRequest struct {
	State string `json:"state"`
}

type BuildImagePostRequest struct {
	ApplicationId string `json:"id"`
}

func GetCatalogAddress() string {
	return os.Getenv("CATALOG_ADDRESS") + "/api/v1"
}

func GetBlobStoreAddress() string {
	return os.Getenv("BLOB_STORE_ADDRESS") + "/api/v1"
}

type Connector struct {
	Server string
	Client *http.Client
}

func NewCatalogConnector() *Connector {
	transport := &http.Transport{}
	clientCreator := &http.Client{Transport: transport, Timeout: time.Duration(30 * time.Minute)}

	return &Connector{
		Server: GetCatalogAddress(),
		Client: clientCreator,
	}
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func NewBlobStoreConnector() *Connector {
	transport := &http.Transport{}
	clientCreator := &http.Client{Transport: transport, Timeout: time.Duration(30 * time.Minute)}

	return &Connector{
		Server: GetBlobStoreAddress(),
		Client: clientCreator,
	}
}
