package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func NewGrpcRegister(port int, register func(grpcServer *grpc.Server)) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	//反射
	reflection.Register(grpcServer)
	register(grpcServer)
	err = grpcServer.Serve(listen)
	if err != nil {
		return err
	}
	log.Println("The server port is located at ", listen.Addr())
	return nil
}
