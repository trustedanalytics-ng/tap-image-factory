package main

import (
	"github.com/gocraft/web"
)

type BuildImagePostRequest struct {
	ApplicationId    string          `json:"organization_guid"`
}

func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(201)
}
