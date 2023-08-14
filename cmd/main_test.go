package main

import (
	"context"
	"log"
	"net"
	"testing"

	"amitshekar-clean-arch/bootstrap"
	repository "amitshekar-clean-arch/repository/mysql"
	pb "amitshekar-clean-arch/todogrpc"
	"amitshekar-clean-arch/usecase"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestTodoUsecase_GetTodo(t *testing.T) {
	// ---- Server Initialization ----
	app = bootstrap.App()

	db := app.MySql

	tr = repository.NewMysqlTodoRepo(db)
	tu = usecase.NewTodoUsecase(tr)
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
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))

	t.Cleanup(func() {
		conn.Close()
	})

	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	client := pb.NewTodoCRUDClient(conn)
	taskName := "Task 1"
	todoStatus := "Done"
	todo := pb.NewTodo{
		Name:   taskName,
		Status: todoStatus,
	}
	todoUpdate := pb.NewTodo{
		Name:   taskName,
		Status: "Not Done",
	}

	pbErrorCreate, _ := client.CreateTodo(context.Background(), &todo)
	assert.Equal(t, "Success!", pbErrorCreate.Message)

	pbErrorUpdate, _ := client.UpdateTodo(context.Background(), &todoUpdate)
	assert.Equal(t, "Success!", pbErrorUpdate.Message)

	pbTodo, err := client.GetTodo(context.Background(), &pb.TodoName{Name: taskName})
	assert.ErrorIs(t, nil, err)
	assert.Equal(t, "Task 1", pbTodo.Name)
	// Test status update
	assert.Equal(t, todoUpdate.Status, pbTodo.Status)

	pbErrorDelete, _ := client.DeleteTodo(context.Background(), &pb.TodoName{Name: taskName})
	assert.Equal(t, "Success!", pbErrorDelete.Message)
}
