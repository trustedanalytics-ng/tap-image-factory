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

package app

import (
	"encoding/json"
	"errors"
	catalogApi "github.com/trustedanalytics/tapng-catalog/client"
	. "github.com/trustedanalytics/tapng-go-common/http"
	"net/http"
	"os"
	"strconv"
	"time"
)

type BlobStoreApi interface {
	GetBlob(blobId string) ([]byte, error)
	GetImageBlob(imageId string) ([]byte, error)
	DeleteBlob(blobId string) error
	DeleteImageBlob(imageId string) error
}

type Connector struct {
	Server string
	Client *http.Client
}

func GetCatalogConnector() (*catalogApi.TapCatalogApiConnector, error) {
	address := GetCatalogAddressWithoutProtocol()
	if os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION") != "" {
		return catalogApi.NewTapCatalogApiWithSSLAndBasicAuth(
			"https://"+address,
			os.Getenv("CATALOG_USER"),
			os.Getenv("CATALOG_PASS"),
			os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION"),
			os.Getenv("CATALOG_SSL_KEY_FILE_LOCATION"),
			os.Getenv("CATALOG_SSL_CA_FILE_LOCATION"),
		)
	} else {
		return catalogApi.NewTapCatalogApiWithBasicAuth(
			"http://"+address,
			os.Getenv("CATALOG_USER"),
			os.Getenv("CATALOG_PASS"),
		)
	}
}

func NewCatalogConnector() *Connector {
	transport := &http.Transport{}
	clientCreator := &http.Client{Transport: transport, Timeout: time.Duration(30 * time.Minute)}

	return &Connector{
		Server: GetCatalogAddressWithoutProtocol(),
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

func (c *Connector) GetImageDetails(imageId string) (*ImageGetResponse, error) {
	response := ImageGetResponse{}

	status, body, err := RestGET(c.Server+"/images/"+imageId, nil, c.Client)

	if status != 200 || err != nil {
		if err == nil {
			err = errors.New("Invalid status: " + strconv.Itoa(status))
		}
		logger.Error(err)
		return &ImageGetResponse{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error(err)
		return &ImageGetResponse{}, err
	}
	return &response, nil
}

func (c *Connector) UpdateImageState(imageId, state string) error {

	req, err := json.Marshal(ImageStatePutRequest{state})
	if err != nil {
		return err
	}

	status, _, err := RestPATCH(c.Server+"/images/"+imageId, string(req), nil, c.Client)

	if status != 200 || err != nil {
		if err == nil {
			err = errors.New("Invalid status: " + strconv.Itoa(status))
		}
		logger.Error(err)
		return err
	}
	return nil
}

func (c *Connector) GetImageBlob(imageId string) ([]byte, error) {
	return c.GetBlob(imageId)
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

func (c *Connector) DeleteImageBlob(imageId string) error {
	return c.DeleteBlob(imageId)
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
