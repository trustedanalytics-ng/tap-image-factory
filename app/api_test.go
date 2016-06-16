package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"

	"github.com/gocraft/web"
	"strings"
	. "github.com/smartystreets/goconvey/convey"
)

const URLbuildImage = "/rest/v1/app"

func prepareMocksAndRouter(t *testing.T) (router *web.Router) {
	router = web.New(Context{})
	return router
}

func TestBuildImage(t *testing.T) {
	request := BuildImagePostRequest{ApplicationId: "test-app-id"}
	router := prepareMocksAndRouter(t)
	router.Post(URLbuildImage, (*Context).BuildImage)

	Convey("Test BuildImage", t, func() {
		Convey("Should returns proper response", func() {
			response := sendRequest("POST", URLbuildImage, marshallToJson(t, request), router)
			assertResponse(response, "", 201)
		})
	})
}

func sendRequest(rType, path string, body []byte, r *web.Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(rType, path, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func assertResponse(rr *httptest.ResponseRecorder, body string, code int) {
	if body != "" {
		So(strings.TrimSpace(string(rr.Body.Bytes())), ShouldContainSubstring, body)
	}
	So(rr.Code, ShouldEqual, code)
}

func marshallToJson(t *testing.T, serviceInstance interface{}) []byte {
	if body, err := json.Marshal(serviceInstance); err != nil {
		t.Errorf(err.Error())
		t.FailNow()
		return nil
	} else {
		return body
	}
}
