package main

import (
	"amitshekar-clean-arch/bootstrap"
	"amitshekar-clean-arch/domain"
	repository "amitshekar-clean-arch/repository/mysql"
	"amitshekar-clean-arch/usecase"
	"context"
	"log"
	"net"

	pb "amitshekar-clean-arch/todogrpc"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	app bootstrap.Application
	tr  domain.TodoRepo
	tu  domain.TodoUsecase
)

type TodoGrpcServer struct {
	pb.UnimplementedTodoCRUDServer
}

func init() {
	app = bootstrap.App()
}

func main() {
	env := app.Env
	db := app.MySql

	tr = repository.NewMysqlTodoRepo(db)
	tu = usecase.NewTodoUsecase(tr)

	lis, err := net.Listen("tcp", env.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodoCRUDServer(s, &TodoGrpcServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *TodoGrpcServer) CreateTodo(ctx context.Context, in *pb.NewTodo) (*pb.Error, error) {
	log.Printf("Task to be created: %v", in.GetName())
	todo := domain.Todo{
		Name:   in.GetName(),
		Status: in.GetStatus(),
	}

	err := tu.CreateTodo(&todo)
	if err == nil {
		return &pb.Error{Message: "Success!"}, nil
	} else {
		log.Println(err)
		return &pb.Error{Message: err.Error()}, nil
	}
}

func (s *TodoGrpcServer) GetAll(ctx context.Context, empty *emptypb.Empty) (*pb.RepeatedTodo, error) {
	log.Println("Getting All data from database")
	todos, err := tu.GetAll()

	var todosPb []*pb.Todo
	if err == nil {
		for i := 0; i < len(todos); i++ {
			todosPb = append(todosPb, &pb.Todo{ID: int32(todos[i].ID), Name: todos[i].Name, Status: todos[i].Status})
		}

		return &pb.RepeatedTodo{Todo: todosPb}, nil
	} else {
		log.Println(err)
		return &pb.RepeatedTodo{Todo: todosPb}, err
	}
}

func (s *TodoGrpcServer) GetTodo(ctx context.Context, in *pb.TodoName) (*pb.Todo, error) {
	log.Printf("Task to be searched: %v", in.GetName())

	todo, err := tu.GetTodo(&in.Name)

	if err == nil {
		return &pb.Todo{ID: int32(todo.ID), Name: todo.Name, Status: todo.Status}, nil
	} else {
		log.Println(err)
		return &pb.Todo{}, err
	}
}

func (s *TodoGrpcServer) DeleteTodo(ctx context.Context, in *pb.TodoName) (*pb.Error, error) {
	log.Printf("Task to be deleted: %v", in.GetName())

	err := tu.DeleteTodo(&in.Name)
	if err == nil {
		return &pb.Error{Message: "Success!"}, nil
	} else {
		log.Println(err)
		return &pb.Error{Message: err.Error()}, nil
	}
}

func (s *TodoGrpcServer) UpdateTodo(ctx context.Context, in *pb.NewTodo) (*pb.Error, error) {
	log.Printf("Task to be updated: %v", in.GetName())

	todo := domain.Todo{
		Name:   in.GetName(),
		Status: in.GetStatus(),
	}

	err := tu.UpdateTodo(&todo)

	if err == nil {
		return &pb.Error{Message: "Success!"}, nil
	} else {
		log.Println(err)
		return &pb.Error{Message: err.Error()}, nil
	}
}
