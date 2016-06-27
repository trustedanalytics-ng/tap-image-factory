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
			dockerfileStr, _ := StreamToString(dockerfile)
			dockerfileStringArray := strings.Split(dockerfileStr, "\n")

			So(dockerfileStringArray[0], ShouldEqual, "FROM "+testBaseImage)
			So(dockerfileStringArray[1], ShouldEqual, "ADD "+artifactFileName+" /root")
			So(dockerfileStringArray[2], ShouldEqual, "CMD [/root/run.sh]")
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
			trStr, _ := StreamToString(tr)
			So(trStr, ShouldEqual, "Test artifact file content")

			hdr, err = tr.Next()

			So(hdr.Name, ShouldEqual, "Dockerfile")
			trStr, _ = StreamToString(tr)
			So(trStr, ShouldEqual, testDockerfileContent)
		})
	})
}
