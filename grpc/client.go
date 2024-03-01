package grpc

import "google.golang.org/grpc"

func Client(port string, with grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(port, with)
}
