package nacos

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
	"log"
	"xorm.io/xorm"
)

type Database interface {
	Init() error
}

// 初始化数据库连接的方法
func NewDatabasesConn(data ...Database) {
	for _, v := range data {
		err := v.Init()
		if err != nil {
			panic(err)
		}
	}
}

var (
	XDB *xorm.Engine
	RDB *redis.Client
)

// mysql初始化
func (m Mysql) Init() error {
	var err error
	XDB, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4", m.Username, m.Password, m.DatabasesName))
	if err != nil {
		panic(err)
	}
	log.Println("mysql连接成功")
	return nil
}

// redis初始化
func (r Redis) Init() error {
	RDB = redis.NewClient(&redis.Options{
		Addr: r.Host + ":" + r.Port,
	})
	if _, err := RDB.Ping().Result(); err != nil {
		panic(err)
	}
	log.Println("redis连接失败")
	return nil
}

// 初始化服务
func Initialisation(ipAddr string, port uint64, dataId, group string) Config {
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
			IpAddr: ipAddr,
			Port:   port,
		},
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	config, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		panic(err)
	}
	var Configs Config
	err = yaml.Unmarshal([]byte(config), &Configs)
	if err != nil {
		panic(err)
	}
	log.Println(Configs)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生更改")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		panic(err)
	}
	return Configs
}

type Config struct {
	Mysql Mysql `yaml:"Mysql"`
	Redis Redis `yaml:"Redis"`
}

type Mysql struct {
	Username      string `yaml:"Username"`
	Password      string `yaml:"Password"`
	DatabasesName string `yaml:"DatabasesName"`
}
type Redis struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}
