SOURCE := $(shell find . -name '*.proto')

PROJECT=api-service
ORGANIZATION=smartcentrix
DOCKER_REPO = $(ORGANIZATION)/$(PROJECT)
VERSION_TAG = 0.0.1

GO_SOURCE := $(shell find . -name '*.go')

build: $(GO_SOURCE)
	docker build -t $(DOCKER_REPO):$(VERSION_TAG) .

push: build
	docker push $(DOCKER_REPO):$(VERSION_TAG)
	docker tag $(DOCKER_REPO):$(VERSION_TAG) $(DOCKER_REPO):latest
	docker push $(DOCKER_REPO):latest

proto: $(SOURCE)
	protoc -I . \
	  -I $(GOPATH)/src \
	  -I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  -I /usr/local/include \
	  --go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	  ./proto/*.proto

	protoc -I . \
	  -I $(GOPATH)/src \
	  -I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  -I /usr/local/include \
	  --grpc-gateway_out=logtostderr=true:. \
	  ./proto/*.proto

	protoc -I . \
	  -I $(GOPATH)/src \
	  -I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  -I /usr/local/include \
	  --swagger_out=logtostderr=true:. \
	  ./proto/*.proto