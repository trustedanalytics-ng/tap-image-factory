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
	"github.com/trustedanalytics/tap-go-common/http"
)

type CatalogApi interface {
	UpdateApplicationState(applicationId string)
	GetApplicationDetails(applicationId string)
}

func (c *CatalogConnector) GetApplicationDetails(applicationId string) (*ApplicationGetResponse, error) {
	response := ApplicationGetResponse{}

	status, body, err := http.RestGET(c.Server+"/applications/"+applicationId, nil, c.Client)

	if status != 200 || err != nil {
		logger.Error("[GetApplicationDetails] Status: ", status)
		logger.Error("[GetApplicationDetails] Error: ", err)
		return &ApplicationGetResponse{}, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("[GetApplicationDetails] Error: ", err)
		return &ApplicationGetResponse{}, err
	}

	return &response, nil
}

func (c *CatalogConnector) UpdateApplicationState(applicationId, state string) error {

	req, err := json.Marshal(ApplicationStatePutRequest{applicationId, state})
	if err != nil {
		return err
	}

	status, _, err := http.RestPUT(c.Server+"/applications/"+applicationId, nil, string(req), c.Client)

	if status != 200 || err != nil {
		logger.Error("[UpdateApplicationState] Status: ", status)
		logger.Error("[UpdateApplicationState] Error: ", err)
		return err
	}

	if err != nil {
		logger.Error("[UpdateApplicationState] Error: ", err)
		return err
	}

	return nil
}
