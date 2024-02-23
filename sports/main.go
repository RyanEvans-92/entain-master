package main

import (
	"database/sql"
	"flag"
	"log"
	"net"

	"sports/db"
	"sports/proto/sports"
	"sports/service"

	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9999", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s\n", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9999")
	if err != nil {
		return err
	}

	sportingDB, err := sql.Open("sqlite3", "./db/events.db")
	if err != nil {
		return err
	}

	eventsRepo := db.NewEventsRepo(sportingDB)
	if err := eventsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sports.RegisterEventsServer(
		grpcServer,
		service.NewEventsService(
			eventsRepo,
		),
	)

	log.Printf("gRPC server listening on: %s\n", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
