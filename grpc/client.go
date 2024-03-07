package grpc

import (
	"github.com/dinzhen12306/framework/consul"
	"google.golang.org/grpc"
)

func Client(name string, with grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := consul.AgentHealthServiceByName(name)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(conn, with)
}
