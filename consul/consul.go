package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
)

func ServiceRegister(serverName, address string, port int) error {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      serverName,
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
