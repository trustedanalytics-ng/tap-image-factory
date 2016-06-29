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
	"github.com/gocraft/web"
	"github.com/trustedanalytics/tap-go-common/util"
)

type Context struct {
	BlobStoreConnector *BlobStoreConnector
	CatalogConnector   *CatalogConnector
	DockerConnector    *DockerHandler
}

func (c *Context) SetupContext(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.BlobStoreConnector = NewBlobStoreConnector()
	c.CatalogConnector = NewCatalogConnector()
	dockerClient, err := NewDockerClient()
	if err != nil {
		logger.Critical("Could not instantiate Docker Client!")
	}
	c.DockerConnector = &DockerHandler{dockerClient, &DockerClient{}}
	next(rw, req)
}

func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	req_json := BuildImagePostRequest{}
	err := util.ReadJson(req, &req_json)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(400)
		return
	}
	appDetails, err := c.CatalogConnector.Api.GetApplicationDetails(req_json.ApplicationId)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	blobBytes, err := c.BlobStoreConnector.Api.GetApplicationBlob(appDetails.ApplicationId)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.CatalogConnector.Api.UpdateApplicationState(appDetails.ApplicationId, "BUILDING")
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	tag := GetDockerHostAddress() + "/" + appDetails.ApplicationId
	err = c.DockerConnector.Api.CreateImage(bytes.NewReader(blobBytes), appDetails.BaseImage, tag)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.DockerConnector.Api.PushImage(tag)
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.CatalogConnector.Api.UpdateApplicationState(appDetails.ApplicationId, "READY")
	if err != nil {
		logger.Error(err.Error())
		rw.WriteHeader(500)
		return
	}
	err = c.BlobStoreConnector.Api.DeleteApplicationBlob(appDetails.ApplicationId)
	if err != nil {
		logger.Error(err.Error())
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(201)
}
