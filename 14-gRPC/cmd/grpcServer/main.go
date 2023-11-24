package main

import (
	"database/sql"
	"net"

	"github.com/christianferraz/goexpert/14-gRPC/internal/database"
	"github.com/christianferraz/goexpert/14-gRPC/internal/pb"
	"github.com/christianferraz/goexpert/14-gRPC/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDB)
	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}
}
