package grpc

import (
	"github.com/dinzhen12306/framework/consul"
	"google.golang.org/grpc"
	"log"
)

func Client(name string, with grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := consul.AgentHealthServiceByName(name)
	log.Println(conn, err)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(conn, with)
}
