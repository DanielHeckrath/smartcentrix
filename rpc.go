package main

import (
	"time"

	"github.com/DanielHeckrath/smartcentrix-notifications/notification"

	"github.com/DanielHeckrath/singleton"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/juju/errors"
	"google.golang.org/grpc"
)

const rpcTimeout = time.Second * 20

var notificationClientSingleton singleton.Singleton

func getNotificationServiceClient() (notification.NotificationServiceClient, error) {
	entity, err := notificationClientSingleton.Get(func() (interface{}, error) {
		conn, err := dial("notifications:8081")
		if err != nil {
			return nil, errors.Annotate(err, "Unable to create grpc connection")
		}
		return notification.NewNotificationServiceClient(conn), nil
	})

	if err != nil {
		return nil, err
	}

	client, ok := entity.(notification.NotificationServiceClient)

	if !ok {
		return nil, errors.New("Cannot cast singleton to notification service client")
	}

	return client, nil
}

func dial(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithTimeout(rpcTimeout),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	)
}
