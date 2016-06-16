package main

import (
	"net/http"

	"github.com/gocraft/web"
	"strconv"

	"github.com/trustedanalytics/image-factory/logger"
)

type Context struct{}

var logger = logger_wrapper.InitLogger("main")

var port = 8080

func main() {
	r := web.New(Context{})
	r.Post("/rest/v1/app", (*Context).BuildImage)

	err := http.ListenAndServe("localhost"+":"+strconv.Itoa(port), r)
	if err != nil {
		logger.Critical("Couldn't serve app on port ", port, " Application will be closed now.")
	}
}
