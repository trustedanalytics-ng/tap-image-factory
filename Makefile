export GOBIN=$(GOPATH)/bin
export APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)
COMMIT_COUNT=`git rev-list --count origin/master`
COMMIT_SHA=`git rev-parse HEAD`
VERSION=0.1.0
all: build

build: bin/app
	@echo "build complete."

bin/app: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

local_bin/app: verify_gopath
	CGO_ENABLED=0 go install -tags local $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)

run: local_bin/app
	$(GOPATH)/bin/app


pack: build
	cp -Rf $(GOBIN)/app application
	echo "commit_sha=$(COMMIT_SHA)" > build_info.ini
	zip -r -q image-factory-${VERSION}.zip application build_info.ini


