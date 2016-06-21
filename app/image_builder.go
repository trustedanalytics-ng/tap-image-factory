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

func CreateImage(artifact io.Reader, baseImage string) error {
	dockerfile := createDockerfile(baseImage)
	buildContext, err := createBuildContext(artifact, dockerfile)
	if err != nil {
		return err
	}
	buildImage(buildContext)
	return nil
}

func createClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.23", nil, defaultHeaders)
	if err != nil {
		logger.Error("Couldn't create docker client", err)
		return nil, err
	}
	return cli, nil
}

func createDockerfile(baseImage string) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString("FROM ")
	buf.WriteString(baseImage)
	buf.WriteString("\n")
	buf.WriteString("ADD " + artifactFileName + " ~\n")
	buf.WriteString("CMD [~/run.sh]")
	return buf
}

func createBuildContext(applicationArtifact io.Reader, dockerfile io.Reader) (io.Reader, error) {
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

func buildImage(buildContext io.Reader) error {
	dockerClient, err := createClient()
	if err != nil {
		return err
	}
	buildOptions := types.ImageBuildOptions{}
	_, err = dockerClient.ImageBuild(context.Background(), buildContext, buildOptions)
	if err != nil {
		logger.Error("Couldn't build docker image", err)
		return err
	}
	return nil
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
