APDIR=$(shell go list ./... | grep -v /vendor/)
GOBIN=$(GOPATH)/bin

build:
	CGO_ENABLED=0 go install -tags netgo ${APDIR}
	go fmt $(APDIR)

run: build
	${GOPATH}/bin/tapng-image-factory

run-local: build
	PORT=8086 HUB_ADDRESS="http://localhost:5000" BLOB_STORE_HOST="http://localhost" BLOB_STORE_PORT=8084 \
	 $(GOPATH)/bin/tapng-image-factory

docker_build: build
	rm -Rf application && mkdir application
	cp -Rf $(GOBIN)/tapng-image-factory application/
	docker build -t tapng-image-factory .

push_docker: docker_build
	docker tag tapng-image-factory $(REPOSITORY_URL)/tapng-image-factory:latest
	docker push $(REPOSITORY_URL)/tapng-image-factory:latest

kubernetes_deploy:
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

deps_update: verify_gopath
	$(GOBIN)/govendor remove +all
	$(GOBIN)/govendor add +external
	@echo "Done"

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetchinf specific deps in newest versions"

	$(GOBIN)/govendor fetch -v $(DEP_URL)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

tests: verify_gopath
	go test --cover $(APP_DIR_LIST)
	
prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tapng-image-factory
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tapng-image-factory

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tapng-image-factory/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp $(GOPATH)/bin/tapng-image-factory ./application/tapng-image-factory	
