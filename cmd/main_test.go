package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "amitshekar-clean-arch/todogrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestTodoUsecase_GetTodo(t *testing.T) {
	// ---- Server Initialization ----

	// setup a listener
	lis := bufconn.Listen(1024 * 1024)
	// setup defer to clean
	t.Cleanup(func() {
		lis.Close()
	})
	// setup grpc server

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	tgs := TodoGrpcServer{}
	pb.RegisterTodoCRUDServer(srv, &tgs)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	// ---- Test ----
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())

	t.Cleanup(func() {
		conn.Close()
	})

	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	client := pb.NewTodoCRUDClient(conn)
	taskName := "Task 1"
	res, err := client.GetTodo(context.Background(), &pb.TodoName{Name: taskName})

	if err != nil {
		t.Fatalf("client.gettodo %v", err)
	}

	if res.Name != "Task 1" {
		t.Fatal("Values that is returned is error!")
	}
}
