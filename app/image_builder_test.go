package main

import (
	"testing"

	"archive/tar"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
)

const (
	testBaseImage = "test-base-image:0.1"
)

var (
	testApplicationArtifactContent = "Test artifact file content"
	testDockerfileContent          = "FROM " + testBaseImage + "\n" + "RUN mkdir /test_dir\n" + "CMD [~/run.sh]"

	testApplicationArtifact = bytes.NewBufferString(testApplicationArtifactContent)
	testDockerfile          = bytes.NewBufferString(testDockerfileContent)
)

func TestCreateDockerfile(t *testing.T) {
	Convey("Test CreateDockerfile", t, func() {
		Convey("Dockerfile should contains proper lines", func() {
			dockerfile := CreateDockerfile(testBaseImage)
			dockerfileStringArray := strings.Split(StreamToString(dockerfile), "\n")

			So(dockerfileStringArray[0], ShouldEqual, "FROM "+testBaseImage)
			So(dockerfileStringArray[len(dockerfileStringArray)-1], ShouldEqual, "CMD [~/run.sh]")
		})
	})
}

func TestCreateBuildContext(t *testing.T) {
	Convey("Test CreateBuildContext", t, func() {
		Convey("Context should contains proper files", func() {
			context, err := CreateBuildContext(testApplicationArtifact, testDockerfile)
			So(err, ShouldBeNil)
			tr := tar.NewReader(context)
			hdr, err := tr.Next()

			So(hdr.Name, ShouldEqual, artifactFileName)
			So(StreamToString(tr), ShouldEqual, "Test artifact file content")

			hdr, err = tr.Next()

			So(hdr.Name, ShouldEqual, "Dockerfile")
			So(StreamToString(tr), ShouldEqual, testDockerfileContent)
		})
	})
}
