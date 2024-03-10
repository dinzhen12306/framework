package grpc

import (
	"fmt"
	"github.com/dinzhen12306/framework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func NewGrpcRegister(serverName, host string, port int, register func(grpcServer *grpc.Server)) error {
	err := consul.ServiceRegister(serverName, host, port)
	if err != nil {
		return err
	}
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	//反射
	reflection.Register(grpcServer)
	//支持健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	register(grpcServer)
	err = grpcServer.Serve(listen)
	if err != nil {
		return err
	}
	log.Println("The server port is located at ", listen.Addr())
	return nil
}

func GetIp() (ip []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, addr := range addrs {
		ipNet, isVailIpNet := addr.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = append(ip, ipNet.IP.String())
			}
		}

	}
	return ip
}
