package main

import (
	"archive/tar"
	"bytes"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
	"io"
)

const (
	artifactFileName = "artifact.tar"
)

type DockerClient struct {
	cli client.Client
}

type ImageBuilder interface {
	CreateImage(artifact io.Reader, baseImage string)
	BuildImage(buildContext io.Reader)
	TagImage(imageId, tag string)
	PushImage(tag string)
}

func CreateClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.23", nil, defaultHeaders)
	if err != nil {
		logger.Error("Couldn't create docker client", err)
		return nil, err
	}
	return cli, nil
}

func CreateDockerfile(baseImage string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString("FROM ")
	buf.WriteString(baseImage)
	buf.WriteString("\n")
	buf.WriteString("ADD " + artifactFileName + " ~\n")
	buf.WriteString("CMD [~/run.sh]")
	return buf
}

func CreateBuildContext(applicationArtifact io.Reader, dockerfile io.Reader) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	hdr := &tar.Header{
		Name: artifactFileName,
		Mode: 0777,
		Size: int64(len(StreamToByte(applicationArtifact))),
	}
	err := tw.WriteHeader(hdr)
	if err != nil {
		logger.Error("Couldn't write tar header for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	_, err = tw.Write(StreamToByte(applicationArtifact))
	if err != nil {
		logger.Error("Couldn't write tar body for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	hdr = &tar.Header{
		Name: "Dockerfile",
		Mode: 0777,
		Size: int64(len(StreamToByte(dockerfile))),
	}
	err = tw.WriteHeader(hdr)
	if err != nil {
		logger.Error("Couldn't write tar header for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	_, err = tw.Write(StreamToByte(dockerfile))
	if err != nil {
		logger.Error("Couldn't write tar body for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	tw.Close()

	return buf, nil
}

func (d *DockerClient) CreateImage(artifact io.Reader, baseImage string) error {
	dockerfile := CreateDockerfile(baseImage)
	buildContext, err := CreateBuildContext(artifact, dockerfile)
	if err != nil {
		return err
	}
	d.BuildImage(buildContext)
	return nil
}

func (d *DockerClient) BuildImage(buildContext io.Reader) error {
	buildOptions := types.ImageBuildOptions{}
	_, err := d.cli.ImageBuild(context.Background(), buildContext, buildOptions)
	if err != nil {
		logger.Error("Couldn't build docker image", err)
		return err
	}
	return nil
}
