package main

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() {
	runtime.HTTPError = httpError
}

func httpError(ctx context.Context, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	cause := errors.Cause(err)
	code := grpc.Code(cause)
	msg := err.Error()
	runtime.DefaultHTTPError(ctx, marshaler, w, req, grpc.Errorf(code, msg))
}
