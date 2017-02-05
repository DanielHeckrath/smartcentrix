package main

import (
	_ "expvar"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/DanielHeckrath/smartcentrix/proto"
	"github.com/DanielHeckrath/smartcentrix/signals"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var errSystemInterupt = errors.New("Received system interupt")

func main() {
	// setup database connection
	db, err := newDatabase()
	if err != nil {
		log.Println("Unable to create database connection")
		log.Printf("%s\n", err)
		return
	}
	defer db.Close()
	db.LogMode(true)

	errc := make(chan error)

	// Transport: debug/metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Serving debug metrics (expvar, pprof, metrics) on HTTP :8082")
		errc <- http.ListenAndServe(":8082", nil)
	}()

	// Transport: http
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := smartcentrix.RegisterSensorApiServiceHandlerFromEndpoint(ctx, mux, "localhost:8081", opts)
		if err != nil {
			errc <- err
		}

		log.Println("Serving transport HTTP on :8080")
		errc <- http.ListenAndServe(":8080", mux)
	}()

	// Transport: grpc
	go func() {
		ln, err := net.Listen("tcp", ":8081")
		if err != nil {
			errc <- err
		}
		s := grpc.NewServer()
		smartcentrix.RegisterSensorApiServiceServer(s, &sensorAPI{
			userRepo:   &sqlUserRepository{db},
			sensorRepo: &sqlSensorRepository{db},
			roomRepo:   &sqlRoomRepository{db},
		})

		log.Println("Serving transport gRPC on :8081")
		errc <- s.Serve(ln)
	}()

	go signals.Handle(quit(errc))

	err = <-errc

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
