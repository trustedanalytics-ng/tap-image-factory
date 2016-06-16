package main

import (
	"github.com/gocraft/web"
)

func (c *Context) BuildImage(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(201)
}

