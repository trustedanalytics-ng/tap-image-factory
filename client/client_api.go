package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	brokerHttp "github.com/trustedanalytics/tapng-go-common/http"
)

type TapApiImageFactoryApi interface {
	BuildImage(imageId string) error
}

func NewTapImageFactoryApiWithBasicAuth(address, username, password string) (*TapImageFactoryApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithBasicAuth()
	if err != nil {
		return nil, err
	}
	return &TapImageFactoryApiConnector{address, username, password, client}, nil
}

func NewTapImageFactoryApiWithSSLAndBasicAuth(address, username, password, certPemFile, keyPemFile, caPemFile string) (*TapImageFactoryApiConnector, error) {
	client, _, err := brokerHttp.GetHttpClientWithCertAndCaFromFile(certPemFile, keyPemFile, caPemFile)
	if err != nil {
		return nil, err
	}
	return &TapImageFactoryApiConnector{address, username, password, client}, nil
}

type TapImageFactoryApiConnector struct {
	Address  string
	Username string
	Password string
	Client   *http.Client
}

func (c *TapImageFactoryApiConnector) getApiConnector(url string) brokerHttp.ApiConnector {
	return brokerHttp.ApiConnector{
		BasicAuth: &brokerHttp.BasicAuth{c.Username, c.Password},
		Client:    c.Client,
		Url:       url,
	}
}

type RequestJson struct {
	ImageId string
}

func (c *TapImageFactoryApiConnector) BuildImage(imageId string) error {
	req_json := RequestJson{
		ImageId: imageId,
	}
	requestBodyByte, err := json.Marshal(req_json)
	if err != nil {
		return err
	}
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/image", c.Address))
	_, _, err = brokerHttp.RestPOST(connector.Url, string(requestBodyByte), connector.BasicAuth, connector.Client)
	return err
}

func (c *TapImageFactoryApiConnector) GetImageFactoryHealth() error {
	connector := c.getApiConnector(fmt.Sprintf("%s/api/v1/healthz", c.Address))
	status, _, err := brokerHttp.RestGET(connector.Url, connector.BasicAuth, connector.Client)
	if status != http.StatusOK {
		err = errors.New("Invalid health status: " + string(status))
	}
	return err
}
