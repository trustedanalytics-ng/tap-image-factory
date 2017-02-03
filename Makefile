# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
GOBIN=$(GOPATH)/bin
APP_NAME=tap-image-factory

build: verify_gopath
	go fmt $(APP_DIR_LIST)
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	mkdir -p application && cp -f $(GOBIN)/$(APP_NAME) ./application/$(APP_NAME)

run: build
	${GOPATH}/bin/tap-image-factory

run-local: build
	PORT=8086 HUB_ADDRESS="http://localhost:5000" BLOB_STORE_HOST="http://localhost" BLOB_STORE_PORT=8084 \
	QUEUE_HOST=127.0.0.1 QUEUE_PORT=5672 QUEUE_USER=guest QUEUE_PASS=guest QUEUE_NAME=image-factory \
	$(GOPATH)/bin/tap-image-factory

docker_build: build_anywhere
	docker build -t tap-image-factory .

push_docker: docker_build
	docker tag tap-image-factory $(REPOSITORY_URL)/tap-image-factory:latest
	docker push $(REPOSITORY_URL)/tap-image-factory:latest

kubernetes_deploy:
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetching specific dependency in newest versions"
	$(GOBIN)/govendor fetch -v $(DEP_URL)

deps_update_tap: verify_gopath
	$(GOBIN)/govendor update github.com/trustedanalytics/...
	$(GOBIN)/govendor remove github.com/trustedanalytics/$(APP_NAME)/...
	@echo "Done"

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

test: verify_gopath
	go test --cover $(APP_DIR_LIST)
	
prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tap-image-factory
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tap-image-factory

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-image-factory/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp ./tap-image-factory ./application/tap-image-factory
	rm -Rf ./temp

mock_update:
	$(GOBIN)/mockgen -source=app/factory.go -package=app -destination=app/factory_mock_test.go
	$(GOBIN)/mockgen -source=app/archive_reader.go -package=app -destination=app/archive_reader_mock_test.go
	$(GOBIN)/mockgen -source=app/image_builder.go -package=app -destination=app/image_builder_mock.go
	$(GOBIN)/mockgen -source=vendor/github.com/trustedanalytics/tap-catalog/client/client.go -package=app -destination=app/catalog_mock_test.go
	$(GOBIN)/mockgen -source=vendor/github.com/trustedanalytics/tap-blob-store/client/client.go -package=app -destination=app/blob_store_mock_test.go
	./add_license.sh
