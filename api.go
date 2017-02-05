package main

import (
	"strings"

	"github.com/DanielHeckrath/smartcentrix/proto"

	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var (
	errorNotImplemented      = grpc.Errorf(codes.Unimplemented, "This method is not yet implemented")
	errorMissingMetadata     = grpc.Errorf(codes.InvalidArgument, "Unable to load request metadata")
	errorMissingCredentials  = grpc.Errorf(codes.InvalidArgument, "User credentials are missing in request metadata")
	errorUnauthorized        = grpc.Errorf(codes.Unauthenticated, "Cannot authenticate user with given credentials")
	errorPermissionDenied    = grpc.Errorf(codes.PermissionDenied, "You are not allowed to access this resource")
	errorInternalServerError = grpc.Errorf(codes.Internal, "Some internal operation had an error")
)

// sensorAPI implements the smartcentrix.SensorApiServiceServer interface
type sensorAPI struct {
	userRepo   UserRepository
	sensorRepo SensorRepository
}

func (s *sensorAPI) ListSensorMeasurement(context.Context, *smartcentrix.ListSensorMeasurementRequest) (*smartcentrix.ListSensorMeasurementResponse, error) {
	return nil, errorNotImplemented
}

func (s *sensorAPI) UpdateSensorMeasurement(context.Context, *smartcentrix.UpdateSensorMeasurementRequest) (*smartcentrix.UpdateSensorMeasurementResponse, error) {
	return nil, errorNotImplemented
}

// authorizeUser validates a users existence with the jwt token thats supplied in the request context
func (s *sensorAPI) authorizeUser(ctx context.Context) (*User, error) {
	// try to get metadata from request context
	md, ok := metadata.FromContext(ctx)

	if !ok {
		return nil, errorMissingMetadata
	}

	// get authorization header and check if not empty
	authorization := md["authorization"]

	if len(authorization) == 0 || authorization[0] == "" {
		return nil, errors.Annotate(errorMissingCredentials, "Authorization Header is missing")
	}

	// header should have format "Authorization: Bearer Token"
	header := strings.SplitN(authorization[0], " ", 2)
	if len(header) != 2 {
		return nil, errors.Annotate(errorMissingCredentials, "Authorization Header has wrong number of parts")
	}

	// validate access token
	claims, err := validateToken(header[1])

	if err != nil {
		return nil, errors.Wrapf(err, errorUnauthorized, "Access Token could not be validated")
	}

	// load user with user id from claims
	user, err := s.userRepo.GetUser(ctx, claims.UserID)

	if err != nil {
		return nil, errors.Wrapf(err, errorUnauthorized, "Could not load user with id %s", claims.UserID)
	}

	return user, nil
}
