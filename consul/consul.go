package consul

import (
	"errors"
	"fmt"
	nacos "github.com/dinzhen12306/framework/config"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"log"
)

func ServiceRegister(serviceName, serviceHost string, servicePort int) error {
	config, err := nacos.GetConfig()
	if err != nil {
		return err
	}
	data := struct {
		Consul struct {
			Host string `yaml:"Host"`
			Port int    `yaml:"Port"`
		} `yaml:"Consul"`
	}{}
	err = yaml.Unmarshal([]byte(config), &data)
	if err != nil {
		return err
	}
	a := api.DefaultConfig()
	a.Address = fmt.Sprintf("%s:%d", data.Consul.Host, data.Consul.Port)
	client, err := api.NewClient(a)
	if err != nil {
		return err
	}
	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    serviceName,
		Tags:    []string{"GRPC"},
		Port:    servicePort,
		Address: serviceHost,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serviceHost, servicePort),
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func AgentHealthServiceByName(name string) (string, error) {
	config, err := nacos.GetConfig()
	if err != nil {
		return "", err
	}
	data := struct {
		Consul struct {
			Host string `yaml:"Host"`
			Port int    `yaml:"Port"`
		} `yaml:"Consul"`
	}{}
	log.Println(config)
	err = yaml.Unmarshal([]byte(config), &data)
	if err != nil {
		return "", err
	}
	client, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%s:%d", data.Consul.Host, data.Consul.Port)})
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
	return fmt.Sprintf("%s:%d", info[0].Service.Address, info[0].Service.Port), nil
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
