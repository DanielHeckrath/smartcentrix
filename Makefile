SOURCE := $(shell find . -name '*.proto')

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