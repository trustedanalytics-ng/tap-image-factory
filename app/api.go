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
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gocraft/web"

	blobStoreApi "github.com/trustedanalytics/tap-blob-store/client"
	catalogApi "github.com/trustedanalytics/tap-catalog/client"
	catalogModels "github.com/trustedanalytics/tap-catalog/models"
	"github.com/trustedanalytics/tap-go-common/util"
	"github.com/trustedanalytics/tap-image-factory/logger"
	"github.com/trustedanalytics/tap-image-factory/models"
)

var logger = logger_wrapper.InitLogger("app")

type Context struct {
	BlobStoreConnector     *blobStoreApi.TapBlobStoreApiConnector
	TapCatalogApiConnector *catalogApi.TapCatalogApiConnector
	DockerConnector        *DockerClient
}

var ctx *Context

func SetupContext() *Context {
	if ctx == nil {
		ctx = &Context{}
		tapBlobStoreConnector, err := GetBlobStoreConnector()
		if err != nil {
			logger.Panic(err)
		}
		ctx.BlobStoreConnector = tapBlobStoreConnector

		tapCatalogApiConnector, err := GetCatalogConnector()
		if err != nil {
			logger.Panic(err)
		}
		ctx.TapCatalogApiConnector = tapCatalogApiConnector

		dockerClient, err := NewDockerClient()
		if err != nil {
			logger.Panic(err)
		}
		ctx.DockerConnector = dockerClient
	}
	return ctx
}
func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	buildRequest := models.BuildImagePostRequest{}
	if err := util.ReadJson(req, &buildRequest); err != nil {
		util.Respond400(rw, err)
		return
	}

	go func() {
		if err := BuildAndPushImage(buildRequest); err != nil {
			logger.Error("Building image error:", err)
		}
	}()
	util.WriteJson(rw, "", http.StatusAccepted)
}

func BuildAndPushImage(buildRequest models.BuildImagePostRequest) error {
	if err := updateImageWithState(buildRequest.ImageId, "BUILDING", "PENDING"); err != nil {
		return err
	} else {
		imgDetails, _, err := ctx.TapCatalogApiConnector.GetImage(buildRequest.ImageId)
		if err != nil {
			updateImageWithState(buildRequest.ImageId, "ERROR", "")
			return err
		}

		buffer := bytes.Buffer{}
		if err = ctx.BlobStoreConnector.GetBlob(imgDetails.Id, &buffer); err != nil {
			updateImageWithState(buildRequest.ImageId, "ERROR", "")
			return err
		}

		tag := GetHubAddressWithoutProtocol() + "/" + imgDetails.Id
		if err = ctx.DockerConnector.CreateImage(bytes.NewReader(buffer.Bytes()), imgDetails.Type, tag); err != nil {
			updateImageWithState(buildRequest.ImageId, "ERROR", "")
			return err
		}

		if err = ctx.DockerConnector.PushImage(tag); err != nil {
			updateImageWithState(buildRequest.ImageId, "ERROR", "")
			return err
		}

		updateImageWithState(buildRequest.ImageId, "READY", "")

		status, err := ctx.BlobStoreConnector.DeleteBlob(imgDetails.Id)
		if err != nil {
			return err

		}
		if status != http.StatusNoContent {
			logger.Warning("Blob removal failed. Actual status %v", status)
		}
		logger.Info("Image build SUCCESS! Id:", buildRequest.ImageId)
		return nil
	}
}

func updateImageWithState(imageId, state, previousState string) error {
	marshalledValue, err := json.Marshal(state)
	if err != nil {
		return err
	}

	patch := catalogModels.Patch{Operation: catalogModels.OperationUpdate, Field: "State", Value: marshalledValue}
	if previousState != "" {
		previousStateByte, err := json.Marshal(previousState)
		if err != nil {
			return err
		}
		patch.PrevValue = previousStateByte
	}

	_, _, err = ctx.TapCatalogApiConnector.UpdateImage(imageId, []catalogModels.Patch{patch})
	return err
}
