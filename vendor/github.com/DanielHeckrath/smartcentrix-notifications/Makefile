PROJECT=notification-service
ORGANIZATION=smartcentrix
DOCKER_REPO = $(ORGANIZATION)/$(PROJECT)
VERSION_TAG = 0.0.2

GO_SOURCE := $(shell find . -name '*.go')
GOOS := linux
GOARCH := amd64

compile: $(GO_SOURCE)
	echo Building for $(GOOS)/$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -v -o ./build/$(PROJECT)

build: compile
	docker build -t $(DOCKER_REPO):$(VERSION_TAG) .

push: build
	docker tag $(DOCKER_REPO):$(VERSION_TAG) gcr.io/fabric-157610/$(DOCKER_REPO):$(VERSION_TAG)
	gcloud docker -- push gcr.io/fabric-157610/$(DOCKER_REPO):$(VERSION_TAG)

proto:
	protoc -I . \
	  -I $(GOPATH)/src \
	  -I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  -I /usr/local/include \
	  --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	  ./notification/*.proto