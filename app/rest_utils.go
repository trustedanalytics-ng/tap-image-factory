package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func RestGET(url string, client *http.Client) (int, []byte, error) {
	return makeRequest("GET", url, "", "application/json", client)
}

func RestPUT(url, body string, client *http.Client) (int, []byte, error) {
	return makeRequest("PUT", url, body, "application/json", client)
}

func RestPOST(url, body string, client *http.Client) (int, []byte, error) {
	return makeRequest("POST", url, body, "application/json", client)
}

func RestDELETE(url, body string, client *http.Client) (int, []byte, error) {
	return makeRequest("DELETE", url, "", "application/json", client)
}


func makeRequest(reqType, url, body, contentType string, client *http.Client) (int, []byte, error) {
	logger.Info("Doing:  ", reqType, url)

	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(reqType, url, bytes.NewBuffer([]byte(body)))
	} else {
		req, _ = http.NewRequest(reqType, url, nil)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ERROR: Make http request "+reqType, err)
		return -1, nil, err
	}
	ret_code := resp.StatusCode
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("ERROR: Make http request "+reqType, err)
		return -1, nil, err
	}

	logger.Info("CODE:", ret_code, "BODY:", string(data))
	return ret_code, data, nil
}

