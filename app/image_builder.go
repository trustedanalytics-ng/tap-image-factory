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
	cli *client.Client
}

type ImageBuilder interface {
	CreateImage(artifact io.Reader, baseImage string)
	BuildImage(buildContext io.Reader)
	TagImage(imageId, tag string)
	PushImage(tag string)
}

func CreateClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(GetDockerHostAddress(), GetDocerApiVersion(), nil, defaultHeaders)
	if err != nil {
		logger.Error("Couldn't create docker client", err)
		return nil, err
	}
	return cli, nil
}

func CreateDockerfile(baseImage string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString("FROM " + baseImage + "\n")
	buf.WriteString("ADD " + artifactFileName + " /root\n")
	buf.WriteString("CMD [/root/run.sh]")
	return buf
}

func CreateBuildContext(applicationArtifact io.Reader, dockerfile io.Reader) (io.Reader, error) {
	applicationArtifactBytes := StreamToByte(applicationArtifact)
	dockerfileBytes := StreamToByte(dockerfile)

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	hdr := createHeader(artifactFileName, int64(len(applicationArtifactBytes)))
	err := tw.WriteHeader(hdr)
	if err != nil {
		logger.Error("Couldn't write tar header for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	_, err = tw.Write(applicationArtifactBytes)
	if err != nil {
		logger.Error("Couldn't write tar body for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	hdr = createHeader("Dockerfile", int64(len(dockerfileBytes)))
	err = tw.WriteHeader(hdr)
	if err != nil {
		logger.Error("Couldn't write tar header for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	_, err = tw.Write(dockerfileBytes)
	if err != nil {
		logger.Error("Couldn't write tar body for file: ", hdr.Name, "Error:", err)
		return nil, err
	}
	tw.Close()

	return buf, nil
}

func createHeader(fileName string, fileSize int64) *tar.Header {
	return &tar.Header{
		Name: fileName,
		Mode: 0700,
		Size: fileSize,
	}
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
	response, err := d.cli.ImageBuild(context.Background(), buildContext, buildOptions)
	logger.Info(StreamToString(response.Body))
	if err != nil {
		logger.Error("Couldn't build docker image", err)
		return err
	}
	return nil
}
