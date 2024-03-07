package redis

import (
	"context"
	"fmt"
	nacos "github.com/dinzhen12306/framework/config"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"time"
)

func withClient(ctx context.Context, fun func(cli *redis.Client)) error {
	config, err := nacos.GetConfig()
	if err != nil {
		return err
	}
	data := struct {
		Redis struct {
			Host string `yaml:"Host"`
			Port string `yaml:"Port"`
		} `yaml:"Redis"`
	}{}
	err = yaml.Unmarshal([]byte(config), &data)
	if err != nil {
		return err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", data.Redis.Host, data.Redis.Port),
	})
	if _, err = rdb.Ping(ctx).Result(); err != nil {
		return err
	}
	defer rdb.Close()
	fun(rdb)
	return nil
}

func Get(ctx context.Context, key string) (str string, err error) {
	err1 := withClient(ctx, func(cli *redis.Client) {
		str, err = cli.Get(ctx, key).Result()
	})
	if err1 != nil {
		return "", err1
	}
	return
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (str string, err error) {
	err1 := withClient(ctx, func(cli *redis.Client) {
		str, err = cli.Set(ctx, key, value, expiration).Result()
	})
	if err1 != nil {
		return "", err1
	}
	return
}

func Exists(ctx context.Context, key string) (res int64, err error) {
	err1 := withClient(ctx, func(cli *redis.Client) {
		res, err = cli.Exists(ctx, key).Result()
	})
	if err1 != nil {
		return 0, err1
	}
	return
}
