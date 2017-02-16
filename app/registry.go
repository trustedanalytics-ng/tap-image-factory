/**
 * Copyright (c) 2017 Intel Corporation
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
	"fmt"
	"net/http"
	"strings"
	"time"

	manifestV2 "github.com/docker/distribution/manifest/schema2"
	commonHTTP "github.com/trustedanalytics-ng/tap-go-common/http"
)

type Registry struct {
	URL    string
	Client *http.Client
}

func NewDockerRegistry(url string) (*Registry, error) {
	client, _, err := commonHTTP.GetHttpClient()
	if err != nil {
		return nil, err
	}
	return newFromCustomClient(url, client)
}

func NewDockerRegistryWithCustomConnectionTimeout(url string, timeout time.Duration) (*Registry, error) {
	client, _, err := commonHTTP.GetHttpClient()
	if err != nil {
		return nil, err
	}
	client.Timeout = timeout
	return newFromCustomClient(url, client)
}

func NewDockerRegistryWithTLSClientCerts(url, caCertFilepath, clientCertFilepath, clientKeyFilepath string) (*Registry, error) {
	tlsHTTPClient, _, err := commonHTTP.GetHttpClientWithCertAndCaFromFile(clientCertFilepath, clientKeyFilepath, caCertFilepath)
	if err != nil {
		return nil, err
	}
	return newFromCustomClient(url, tlsHTTPClient)
}

func newFromCustomClient(url string, httpClient *http.Client) (*Registry, error) {
	registryURL := strings.TrimSuffix(url, "/")

	registryConfig := Registry{
		URL:    registryURL,
		Client: httpClient,
	}

	err := registryConfig.Ping()
	if err != nil {
		return nil, err
	}

	return &registryConfig, nil
}

func (r *Registry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *Registry) Ping() error {
	url := r.url("/v2/")
	logger.Debugf("registry.ping url=%s", url)
	resp, err := r.Client.Get(url)
	if err != nil {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping status code was %v instead of %v", resp.StatusCode, http.StatusOK)
	}
	return nil
}

func (r *Registry) IsImageReady(imageID, imageTag string) (bool, error) {
	url := r.url("/v2/%s/manifests/%s", imageID, imageTag)
	logger.Debugf("registry.manifest.get url=%v repository=%v reference=%v", url, imageID, imageTag)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Accept", manifestV2.MediaTypeManifest)
	resp, err := r.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("bad status code returned from docker-registry: %v", resp.StatusCode)
}
