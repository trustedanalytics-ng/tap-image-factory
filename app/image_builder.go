package main

import (
	"github.com/docker/engine-api/client"
	"golang.org/x/net/context"
	"io"
	"archive/tar"
	"bytes"
	"fmt"
	"os"
)

func createClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.23", nil, defaultHeaders)
	if err != nil {
		logger.Error("Couldn't create docker client", err)
		return nil, err
	}
	return cli, nil
}

func prapareBuildContext(applicationArtefact []byte, dockerfile []byte) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	hdr := &tar.Header{
		Name: "artifact.tar",
		Mode: 0777,
		Size: int64(len(applicationArtefact)),
	}
	tw.WriteHeader(hdr)
	tw.Write(applicationArtefact)

	hdr = &tar.Header{
		Name: "Dockerfile",
		Mode: 0777,
		Size: int64(len(dockerfile)),
	}

	tw.WriteHeader(hdr)
	tw.Write(dockerfile)

	r := bytes.NewReader(buf)

	return r
}

func buildImage(buildContext io.Reader) {
	dockerClient, err := createClient()
	if err != nil {
		return err
	}

	dockerClient.ImageBuild(context.Background(), buildContext, nil)
}
