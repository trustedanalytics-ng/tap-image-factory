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
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/trustedanalytics-ng/tap-catalog/builder"
	catalogModels "github.com/trustedanalytics-ng/tap-catalog/models"
)

const (
	runFileName     = "run.sh"
	defaultImageTag = "latest"
)

type ImageProcessingErr struct {
	imageID   string
	parentErr error
}

func (e ImageProcessingErr) Error() string {
	return e.parentErr.Error()
}

type FactoryAPI interface {
	BuildAndPushImage(imageID string) error
}

type ImageReadinessChecker interface {
	IsImageReady(imageID, imageTag string) (bool, error)
}

type Factory struct {
	imageProber                ImageReadinessChecker
	imageReadinessProbeRetries uint32
	imageReadinessProbeDelay   time.Duration
}

func NewFactoryWithCustomProbeTimes(prober ImageReadinessChecker, imageReadinessProbeRetries uint32, imageReadinessProbeDelaySeconds time.Duration) FactoryAPI {
	return &Factory{
		imageProber:                prober,
		imageReadinessProbeRetries: imageReadinessProbeRetries,
		imageReadinessProbeDelay:   imageReadinessProbeDelaySeconds,
	}
}

func (f *Factory) BuildAndPushImage(imageID string) error {
	logger.Debug("started")

	var err error
	defer func() {
		ipErr, ok := err.(ImageProcessingErr)
		if ok {
			err := markImageAs(catalogModels.ImageStateError, ipErr.imageID)
			if err != nil {
				logger.Errorf("cannot set image (id: %v) state to error: %v", ipErr.imageID, err)
			}
		}
	}()

	err = markImageAs(catalogModels.ImageStateBuilding, imageID)
	if err != nil {
		return err
	}

	imageInfo, err := getImageInfoFromCatalog(imageID)
	if err != nil {
		return err
	}
	defer cleanupBlob(imageID)

	imageAddress, err := buildImage(imageInfo)
	if err != nil {
		appID := catalogModels.GetApplicationId(imageID)
		updateLastStateChangeReason(appID, err.Error())
		return err
	}

	err = pushImage(imageID, imageAddress)
	if err != nil {
		return err
	}

	err = f.ensureImageReady(imageID, defaultImageTag)
	if err != nil {
		return err
	}

	err = markImageAs(catalogModels.ImageStateReady, imageID)
	if err != nil {
		return err
	}

	logger.Info("Image build SUCCESS! Id:", imageID)
	return nil
}

func (f *Factory) ensureImageReady(imageID, imageTag string) error {
	for i := uint32(0); i < f.imageReadinessProbeRetries; i++ {
		imageReady, err := f.imageProber.IsImageReady(imageID, defaultImageTag)
		if err != nil {
			logger.Errorf("Failed to check image readiness: %v", err)
			return ImageProcessingErr{imageID: imageID, parentErr: err}
		}

		if imageReady {
			return nil
		}

		logger.Warningf("Image: %v still not ready in docker-registry. retrying...", imageID)
		time.Sleep(f.imageReadinessProbeDelay)
	}

	err := fmt.Errorf("Image build ERROR! Id: %v - docker registry readiness timeout!", imageID)
	logger.Error(err)
	return ImageProcessingErr{imageID: imageID, parentErr: err}
}

func markImageAs(state catalogModels.ImageState, imageID string) error {
	switch state {
	case catalogModels.ImageStateBuilding:
		return updateImageWithState(imageID, state, catalogModels.ImageStatePending)
	case catalogModels.ImageStateReady, catalogModels.ImageStateError:
		return updateImageWithState(imageID, state, catalogModels.ImageStateBuilding)
	default:
		return fmt.Errorf("inappriopriate state (%v) to mark image as!", state)
	}
}

func updateLastStateChangeReason(id, message string) {
	lastStateChangeReasonMetadata := catalogModels.Metadata{
		Id: catalogModels.LAST_STATE_CHANGE_REASON, Value: message,
	}
	metadataPatch, err := builder.MakePatch("Metadata", lastStateChangeReasonMetadata, catalogModels.OperationAdd)
	if err != nil {
		logger.Errorf("cannot prepare last state reason patch %q for application %q: %v", message, id, err)
		return
	}

	if _, _, err = ctx.TapCatalogApiConnector.UpdateApplication(id, []catalogModels.Patch{metadataPatch}); err != nil {
		logger.Errorf("cannot patch last state change reason %q for application %q: %v", message, id, err)
	}

	return
}

func getImageInfoFromCatalog(imageID string) (catalogModels.Image, error) {
	imgInfo, _, err := ctx.TapCatalogApiConnector.GetImage(imageID)
	if err != nil {
		return imgInfo, ImageProcessingErr{imageID: imageID, parentErr: err}
	}
	return imgInfo, nil
}

func buildImage(imageInfo catalogModels.Image) (string, error) {
	tempBlobFile, err := ioutil.TempFile("", "blob-")
	if err != nil {
		return "", ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
	}
	defer cleanupTempfile(tempBlobFile)

	blobFileWriter := bufio.NewWriter(tempBlobFile)

	if err = ctx.BlobStoreConnector.GetBlob(imageInfo.Id, blobFileWriter); err != nil {
		return "", ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
	}

	err = blobFileWriter.Flush()
	if err != nil {
		return "", ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
	}

	_, err = tempBlobFile.Seek(0, 0)
	if err != nil {
		return "", ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
	}

	if imageInfo.BlobType == catalogModels.BlobTypeTarGz {
		if err = validateTarGz(tempBlobFile); err != nil {
			return "", ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
		}
	}

	_, err = tempBlobFile.Seek(0, 0)
	if err != nil {
		return "", ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
	}

	imageTag := GetImageWithHubAddressWithoutProtocol(imageInfo.Id)

	err = ctx.DockerConnector.CreateImage(tempBlobFile, imageInfo.Type, imageInfo.BlobType, imageTag)
	if err != nil {
		return imageTag, ImageProcessingErr{imageID: imageInfo.Id, parentErr: err}
	}

	return imageTag, nil
}

func validateTarGz(reader io.Reader) error {
	gzf, err := ctx.Reader.NewGzipReader(reader)
	if err != nil {
		err := fmt.Errorf("cannot parse file: %v", err)
		logger.Error(err.Error())
		return err
	}
	defer gzf.Close()

	tarReader := ctx.Reader.NewTarReader(gzf)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			err := fmt.Errorf("parsing error file: %v", err)
			logger.Error(err.Error())
			return err
		}

		if (header.Name == ("./"+runFileName) || header.Name == runFileName) &&
			(header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeSymlink) {
			return nil
		}
	}

	return fmt.Errorf("%s file is missing", runFileName)
}

func pushImage(imageID, imageTag string) error {
	logger.Debug("pushing image with tag: ", imageTag)
	err := ctx.DockerConnector.PushImage(imageTag)
	if err != nil {
		return ImageProcessingErr{imageID: imageID, parentErr: err}
	}
	logger.Debugf("pushing image with tag: %v succeeded", imageTag)
	return nil
}

func cleanupBlob(imageID string) {
	logger.Debugf("removing blob (id: %v) from blob-store", imageID)
	status, err := ctx.BlobStoreConnector.DeleteBlob(imageID)
	if err != nil || status != http.StatusNoContent {
		logger.Errorf("Blob removal failed. Actual status: %v , error: %v", status, err)
		return
	}
	logger.Debugf("blob removal (id: %v) succeeded", imageID)
}

func cleanupTempfile(tempFile *os.File) {
	if tempFile != nil {
		err := tempFile.Close()
		if err != nil {
			logger.Errorf("closing temporary blob file (named: %v) failed: %v", tempFile.Name(), err)
		}
		err = os.Remove(tempFile.Name())
		if err != nil {
			logger.Errorf("removing temporary blob file (named: %v) failed: %v", tempFile.Name(), err)
		}
	}
}

func updateImageWithState(imageId string, state catalogModels.ImageState, previousState catalogModels.ImageState) error {
	patch, err := builder.MakePatchWithPreviousValue("state", state, previousState, catalogModels.OperationUpdate)
	if err != nil {
		return err
	}

	_, _, err = ctx.TapCatalogApiConnector.UpdateImage(imageId, []catalogModels.Patch{patch})
	return err
}
