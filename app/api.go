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
	"net/http"
	"os"
	"time"

	"github.com/gocraft/web"
	"github.com/op/go-logging"

	blobStoreApi "github.com/trustedanalytics-ng/tap-blob-store/client"
	catalogApi "github.com/trustedanalytics-ng/tap-catalog/client"
	commonHttp "github.com/trustedanalytics-ng/tap-go-common/http"
	commonLogger "github.com/trustedanalytics-ng/tap-go-common/logger"
	"github.com/trustedanalytics-ng/tap-image-factory/models"
	"github.com/trustedanalytics-ng/tap-go-common/util"
)

var logger = initLogger()

func initLogger() *logging.Logger {
	logger, err := commonLogger.InitLogger("app")
	if err != nil {
		panic(err)
	}
	return logger
}

type Context struct {
	BlobStoreConnector     blobStoreApi.TapBlobStoreApi
	TapCatalogApiConnector catalogApi.TapCatalogApi
	DockerConnector        ImageBuilder
	Factory                FactoryAPI
	Reader                 ArchiveReader
}

var ctx *Context

var (
	dockerRegistryURL = os.Getenv("HUB_ADDRESS")
	caCertFilepath    = os.Getenv("KUBE_CA_CERT_PATH")
	certFilepath      = os.Getenv("KUBE_CA_CLIENT_CERT_PATH")
	keyFilepath       = os.Getenv("KUBE_CA_CLIENT_KEY_PATH")
)

const (
	imageReadinessProbeRetriesEnvVarName      = "IMAGE_READY_MAX_RETRIES"
	defaultImageReadinessProbeRetries         = uint32(20)
	imageReadinessProbeDelaySecondsEnvVarName = "IMAGE_READY_BACK_OFF_DELAY_SECS"
	defaultImageReadinessProbeDelaySeconds    = 5
)

func parseUint32EnvVarOrWarning(envVarName string, defaultValue uint32) uint32 {
	result, err := util.GetUint32EnvValueOrDefault(envVarName, defaultValue)
	if err != nil {
		logger.Warningf("parsing %v failed - using default value: %v . err: %v", envVarName, defaultValue, err)
	}
	return result
}

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

		dockerRegistryProber, err := NewDockerRegistryWithTLSClientCerts("https://"+dockerRegistryURL, caCertFilepath, certFilepath, keyFilepath)
		if err != nil {
			logger.Panic(err)
		}

		imageReadinessProbeRetries := parseUint32EnvVarOrWarning(imageReadinessProbeRetriesEnvVarName, defaultImageReadinessProbeRetries)
		imageReadinessProbeDelaySeconds := parseUint32EnvVarOrWarning(imageReadinessProbeDelaySecondsEnvVarName, defaultImageReadinessProbeDelaySeconds)

		ctx.Factory = NewFactoryWithCustomProbeTimes(dockerRegistryProber,
			imageReadinessProbeRetries, time.Duration(imageReadinessProbeDelaySeconds)*time.Second)

		ctx.Reader = &DefaultArchiveReader{}
	}
	return ctx
}

func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	buildRequest := models.BuildImagePostRequest{}
	if err := commonHttp.ReadJson(req, &buildRequest); err != nil {
		commonHttp.Respond400(rw, err)
		return
	}

	go func() {
		if err := c.Factory.BuildAndPushImage(buildRequest.ImageId); err != nil {
			logger.Error("Building image error:", err)
		}
	}()
	commonHttp.WriteJson(rw, "", http.StatusAccepted)
}
