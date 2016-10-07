/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"testing"

	"archive/tar"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"github.com/trustedanalytics/tap-catalog/models"
)

var (
	testBaseImage            = ImagesMap["JAVA"]
	testImageArtifactContent = "Test artifact file content"
	testDockerfileContent    = "FROM " + testBaseImage + "\n" + "RUN mkdir /test_dir\n" + "CMD [~/run.sh]"
)

func TestCreateDockerfileForTarGzBlob(t *testing.T) {
	Convey("Test CreateDockerfile", t, func() {
		Convey("Dockerfile should contains proper lines", func() {
			dockerfile, _ := createDockerfile("JAVA", models.BlobTypeTarGz)
			dockerfileStr, _ := StreamToString(dockerfile)
			dockerfileStringArray := strings.Split(dockerfileStr, "\n")

			So(dockerfileStringArray[0], ShouldEqual, "FROM " + GetImageWithHubAddressWithoutProtocol(testBaseImage))
			So(dockerfileStringArray[1], ShouldEqual, "ADD " + blobTypeFileNameMap[models.BlobTypeTarGz]+" /root")
			So(dockerfileStringArray[2], ShouldEqual, "ENV PORT 80")
			So(dockerfileStringArray[3], ShouldEqual, "EXPOSE $PORT")
			So(dockerfileStringArray[4], ShouldEqual, "WORKDIR /root")
			So(dockerfileStringArray[5], ShouldEqual, "CMD [\"/root/run.sh\"]")
		})
	})
}

func TestCreateDockerfileForJarBlob(t *testing.T) {
	Convey("Test CreateDockerfile", t, func() {
		Convey("Dockerfile should contains proper lines", func() {
			dockerfile, _ := createDockerfile("JAVA", models.BlobTypeJar)
			dockerfileStr, _ := StreamToString(dockerfile)
			dockerfileStringArray := strings.Split(dockerfileStr, "\n")

			So(dockerfileStringArray[0], ShouldEqual, "FROM " + GetImageWithHubAddressWithoutProtocol(testBaseImage))
			So(dockerfileStringArray[1], ShouldEqual, "ADD " + blobTypeFileNameMap[models.BlobTypeJar]+" /root")
			So(dockerfileStringArray[2], ShouldEqual, "ADD run.sh /root")
			So(dockerfileStringArray[3], ShouldEqual, "ENV PORT 80")
			So(dockerfileStringArray[4], ShouldEqual, "EXPOSE $PORT")
			So(dockerfileStringArray[5], ShouldEqual, "WORKDIR /root")
			So(dockerfileStringArray[6], ShouldEqual, "CMD [\"/root/run.sh\"]")
		})
	})
}

func TestCreateBuildContextForTarGzBlob(t *testing.T) {
	testImageArtifact := bytes.NewBufferString(testImageArtifactContent)
	testDockerfile    := bytes.NewBufferString(testDockerfileContent)
	Convey("Test CreateBuildContext", t, func() {
		Convey("Context should contains proper files", func() {
			context, err := createBuildContext(testImageArtifact, testDockerfile, models.BlobTypeTarGz)
			So(err, ShouldBeNil)
			tr := tar.NewReader(context)
			hdr, err := tr.Next()

			So(hdr.Name, ShouldEqual, blobTypeFileNameMap[models.BlobTypeTarGz])
			trStr, _ := StreamToString(tr)
			So(trStr, ShouldEqual, "Test artifact file content")

			hdr, err = tr.Next()

			So(hdr.Name, ShouldEqual, "Dockerfile")
			trStr, _ = StreamToString(tr)
			So(trStr, ShouldEqual, testDockerfileContent)
		})
	})
}

func TestCreateBuildContextForJarBlob(t *testing.T) {
	testImageArtifact := bytes.NewBufferString(testImageArtifactContent)
	testDockerfile    := bytes.NewBufferString(testDockerfileContent)
	Convey("Test CreateBuildContext", t, func() {
		Convey("Context should contains proper files", func() {
			context, err := createBuildContext(testImageArtifact, testDockerfile, models.BlobTypeJar)
			So(err, ShouldBeNil)
			tr := tar.NewReader(context)
			hdr, err := tr.Next()

			So(hdr.Name, ShouldEqual, blobTypeFileNameMap[models.BlobTypeJar])
			trStr, _ := StreamToString(tr)
			So(trStr, ShouldEqual, "Test artifact file content")

			hdr, err = tr.Next()

			So(hdr.Name, ShouldEqual, "Dockerfile")
			trStr, _ = StreamToString(tr)
			So(trStr, ShouldEqual, testDockerfileContent)

			hdr, err = tr.Next()

			So(hdr.Name, ShouldEqual, "run.sh")
			trStr, _ = StreamToString(tr)
			So(trStr, ShouldEqual, blobTypeRunCommandMap[models.BlobTypeJar])
		})
	})
}
