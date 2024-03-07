package nacos

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
)

type Nacos struct {
	IpAddr string //主机
	Port   uint64 //端口号
	DataId string //nacos数据库名称
	Group  string //nacos分组名称
}

var (
	client       config_client.IConfigClient
	namingClient naming_client.INamingClient
	err          error
)

// 初始化服务
func Initialisation(nacos *Nacos) error {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacos.IpAddr,
			Port:   nacos.Port,
		},
	}
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	// 创建动态配置客户端的另一种方式 (推荐)
	client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func ListenServer(nacos *Nacos, fun ...func()) error {
	err = client.ListenConfig(vo.ConfigParam{
		DataId: nacos.DataId,
		Group:  nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生更改")
			for _, v := range fun {
				v()
			}
			if err != nil {
				log.Println("参数修改失败")
			}
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	return nil
}

func GetConfig(dataId, group string) (string, error) {
	return client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

//	vo.RegisterInstanceParam{
//			Ip:          ip,
//			Port:        8848,
//			ServiceName: "main.go",
//			Weight:      10,
//			Enable:      true,
//			Healthy:     true,
//			Ephemeral:   true,
//			Metadata:    map[string]string{"idc": "shanghai"},
//			ClusterName: "cluster-a",     // 默认值DEFAULT
//			GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
//		}
func ServerRegisterInstance(registerInstanceParam vo.RegisterInstanceParam) (bool, error) {
	return namingClient.RegisterInstance(registerInstanceParam)
}

// Clusters:    nil, // 查询所有集群的实例
//
//	ServiceName: "example-service", // 要查询的服务名称
//	GroupName:   "DEFAULT_GROUP", // 服务所属的分组名称
//	HealthyOnly: false, // 查询所有实例，包括不健康的实例
func SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return namingClient.SelectInstances(param)
}
