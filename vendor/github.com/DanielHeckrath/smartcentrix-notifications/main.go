package main

import (
	"log"
	"net"
	"os"

	"github.com/DanielHeckrath/smartcentrix-notifications/notification"
	"github.com/DanielHeckrath/smartcentrix-notifications/signals"

	"github.com/juju/errors"
	"google.golang.org/grpc"
)

var errSystemInterupt = errors.New("Received system interupt")

func main() {
	errc := make(chan error)

	// Transport: grpc
	go func() {
		ln, err := net.Listen("tcp", ":8081")
		if err != nil {
			errc <- err
		}
		s := grpc.NewServer()
		notification.RegisterNotificationServiceServer(s, &notificationAPI{})
		errc <- s.Serve(ln)
	}()

	go signals.Handle(quit(errc))

	err := <-errc

	if err != errSystemInterupt && os.Getenv("ENVIRONMENT") != "staging" {
		log.Println("Api service is terminating because of an unexpected error")
		log.Printf("%s\n", err)
		return
	}
}

func quit(out chan error) func() {
	return func() {
		out <- errSystemInterupt
	}
}
