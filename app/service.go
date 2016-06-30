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
	"encoding/json"
	"errors"
	. "github.com/trustedanalytics/tap-go-common/http"
	"net/http"
	"strconv"
	"time"
)

type CatalogApi interface {
	UpdateApplicationState(applicationId, state string) error
	GetApplicationDetails(applicationId string) (*ApplicationGetResponse, error)
}

type BlobStoreApi interface {
	GetBlob(blobId string) ([]byte, error)
	GetApplicationBlob(applicationId string) ([]byte, error)
	DeleteBlob(blobId string) error
	DeleteApplicationBlob(applicationId string) error
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

func NewBlobStoreConnector() *Connector {
	transport := &http.Transport{}
	clientCreator := &http.Client{Transport: transport, Timeout: time.Duration(30 * time.Minute)}

	return &Connector{
		Server: GetBlobStoreAddress(),
		Client: clientCreator,
	}
}

func (c *Connector) GetApplicationDetails(applicationId string) (*ApplicationGetResponse, error) {
	response := ApplicationGetResponse{}

	status, body, err := RestGET(c.Server+"/applications/"+applicationId, nil, c.Client)

	if status != 200 || err != nil {
		if err == nil {
			err = errors.New("Invalid status: " + strconv.Itoa(status))
		}
		logger.Error(err)
		return &ApplicationGetResponse{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error(err)
		return &ApplicationGetResponse{}, err
	}
	return &response, nil
}

func (c *Connector) UpdateApplicationState(applicationId, state string) error {

	req, err := json.Marshal(ApplicationStatePutRequest{state})
	if err != nil {
		return err
	}

	status, _, err := RestPATCH(c.Server+"/applications/"+applicationId, string(req), nil, c.Client)

	if status != 200 || err != nil {
		if err == nil {
			err = errors.New("Invalid status: " + strconv.Itoa(status))
		}
		logger.Error(err)
		return err
	}
	return nil
}

func (c *Connector) GetApplicationBlob(applicationId string) ([]byte, error) {
	blobId := "app_" + applicationId
	return c.GetBlob(blobId)
}

func (c *Connector) GetBlob(blobId string) ([]byte, error) {
	status, res, err := RestGET(c.Server+"/blobs/"+blobId, nil, c.Client)
	if status != 200 || err != nil {
		if err == nil {
			err = errors.New("Invalid status: " + strconv.Itoa(status) + " Response: " + string(res))
		}
		logger.Error(err)
		return nil, err
	}
	return res, err
}

func (c *Connector) DeleteApplicationBlob(applicationId string) error {
	blobId := "app_" + applicationId
	return c.DeleteBlob(blobId)
}

func (c *Connector) DeleteBlob(blobId string) error {
	status, res, err := RestDELETE(c.Server+"/blobs/"+blobId, "", nil, c.Client)
	if status != 204 || err != nil {
		if err == nil {
			err = errors.New("Invalid status: " + strconv.Itoa(status) + " Response: " + string(res))
		}
		logger.Error(err)
		return err
	}
	return nil
}
