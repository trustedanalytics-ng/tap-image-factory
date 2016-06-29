package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
	"io"
)

const (
	artifactFileName = "artifact.tar"
)

type DockerHandler struct {
	Client *DockerClient
	Api    ImageBuilder
}

type DockerClient struct {
	cli *client.Client
}

type ImageBuilder interface {
	CreateImage(artifact io.Reader, baseImage, applicationId string) error
	buildImage(buildContext io.Reader, applicationId string) error
	TagImage(imageId, tag string) error
	PushImage(tag string) error
}

func NewDockerClient() (*DockerClient, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(GetDockerHostAddress(), GetDockerApiVersion(), nil, defaultHeaders)
	if err != nil {
		return nil, errors.New("Couldn't create docker client: " + err.Error())
	}
	dockerClient := DockerClient{cli}
	return &dockerClient, nil
}

func (d *DockerClient) CreateImage(artifact io.Reader, baseImage, tag string) error {
	dockerfile := createDockerfile(baseImage)
	buildContext, err := createBuildContext(artifact, dockerfile)
	if err != nil {
		return err
	}
	d.buildImage(buildContext, tag)
	return nil
}

func (d *DockerClient) buildImage(buildContext io.Reader, tag string) error {
	buildOptions := types.ImageBuildOptions{Tags: []string{tag}}
	response, err := d.cli.ImageBuild(context.Background(), buildContext, buildOptions)
	responseStr, _ := StreamToString(response.Body)
	logger.Info(responseStr)
	if err != nil {
		return errors.New("Couldn't build docker image: " + err.Error())
	}
	return nil
}

func (d *DockerClient) TagImage(imageId, tag string) error {
	return d.cli.ImageTag(context.Background(), imageId, tag)
}

func (d *DockerClient) PushImage(tag string) error {
	_, err := d.cli.ImagePush(context.Background(), tag, types.ImagePushOptions{})
	return err
}

func createDockerfile(baseImage string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString("FROM " + baseImage + "\n")
	buf.WriteString("ADD " + artifactFileName + " /root\n")
	buf.WriteString("CMD [/root/run.sh]")
	return buf
}

func createBuildContext(applicationArtifact io.Reader, dockerfile io.Reader) (io.Reader, error) {
	applicationArtifactBytes, err := StreamToByte(applicationArtifact)
	if err != nil {
		return nil, err
	}
	dockerfileBytes, err := StreamToByte(dockerfile)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err = writeToTar(artifactFileName, applicationArtifactBytes, tw)
	if err != nil {
		return nil, err
	}
	err = writeToTar("Dockerfile", dockerfileBytes, tw)
	if err != nil {
		return nil, err
	}
	tw.Close()

	return buf, nil
}

func writeToTar(fileName string, fileContents []byte, tw *tar.Writer) error {
	hdr := createHeader(fileName, int64(len(fileContents)))
	err := tw.WriteHeader(hdr)
	if err != nil {
		return errors.New("Couldn't write tar header for file: " + hdr.Name + ". Error: " + err.Error())
	}
	_, err = tw.Write(fileContents)
	if err != nil {
		return errors.New("Couldn't write tar body for file: " + hdr.Name + ". Error: " + err.Error())
	}
	return nil
}

func createHeader(fileName string, fileSize int64) *tar.Header {
	return &tar.Header{
		Name: fileName,
		Mode: 0700,
		Size: fileSize,
	}
}
