package main

import (
	"github.com/DanielHeckrath/smartcentrix-notifications/notification"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	errorNotImplemented = grpc.Errorf(codes.Unimplemented, "This method is not yet implemented")
)

// notificationAPI implements the notification.NotificationServiceServer interface
type notificationAPI struct{}

func (s *notificationAPI) SendNotification(context.Context, *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	return nil, errorNotImplemented
}
