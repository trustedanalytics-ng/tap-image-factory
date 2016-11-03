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
	blobStoreApi "github.com/trustedanalytics/tap-blob-store/client"
	catalogApi "github.com/trustedanalytics/tap-catalog/client"
	"github.com/trustedanalytics/tap-go-common/util"
)

type BlobStoreApi interface {
	GetBlob(blobId string) ([]byte, error)
	GetImageBlob(imageId string) ([]byte, error)
	DeleteBlob(blobId string) error
	DeleteImageBlob(imageId string) error
}

func GetCatalogConnector() (*catalogApi.TapCatalogApiConnector, error) {
	componentName := "CATALOG"

	address, err := util.GetConnectionAddressFromEnvs(componentName)
	if err != nil {
		panic(err)
	}

	username, password, err := util.GetConnectionCredentialsFromEnvs(componentName)
	if err != nil {
		panic(err)
	}

	return catalogApi.NewTapCatalogApiWithBasicAuth("https://"+address, username, password)
}

func GetBlobStoreConnector() (*blobStoreApi.TapBlobStoreApiConnector, error) {
	componentName := "BLOB_STORE"

	address, err := util.GetConnectionAddressFromEnvs(componentName)
	if err != nil {
		panic(err)
	}

	username, password, err := util.GetConnectionCredentialsFromEnvs(componentName)
	if err != nil {
		panic(err)
	}
	return blobStoreApi.NewTapBlobStoreApiWithBasicAuth("https://"+address, username, password)
}
