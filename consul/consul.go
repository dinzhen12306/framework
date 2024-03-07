package consul

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceRegister(serverName, address string, port int) error {
	client, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%s:%d", address, 8500)})
	if err != nil {
		return err
	}
	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    serverName,
		Tags:    []string{"GRPC"},
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", address, port),
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func AgentHealthServiceByName(name string) (string, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", err
	}
	sr, info, err := client.Agent().AgentHealthServiceByName(name)
	if err != nil {
		return "", err
	}
	if sr != "passing" {
		return "", errors.New("services without health")
	}
	return fmt.Sprintf("%s:%d", info[0].Service.Service, info[0].Service.Port), nil
}

func Dial(ServerHost, ServerPort, ServerName string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		// consul服务
		fmt.Sprintf("%s:%s", ServerHost, ServerPort),
		// 指定round_robin策略
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "round_robin", "service": {"name": "%s"}}`, ServerName)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

