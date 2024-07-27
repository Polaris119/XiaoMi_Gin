package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
	"os"
	"time"
)

var contx = context.Background()
var rdbClient *redis.Client
var redisEnable bool

func init() {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("文件连接失败:%v", iniErr)
		os.Exit(1)
	}

	ip := config.Section("redis").Key("ip").String()
	port := config.Section("redis").Key("port").String()
	redisEnable, _ = config.Section("redis").Key("redisEnable").Bool()

	if redisEnable {
		// 连接Redis数据库
		rdbClient = redis.NewClient(&redis.Options{
			Addr:     ip + ":" + port,
			Password: "",
			DB:       0,
		})
		_, err := rdbClient.Ping(contx).Result()
		if err != nil {
			fmt.Println("redis数据库连接失败")
		} else {
			fmt.Println("redis数据库连接成功")
		}
	}
}

type cacheDb struct{}

func (c cacheDb) Set(key string, value interface{}, expiration int) {
	if redisEnable {
		v, err := json.Marshal(value)
		if err == nil {
			rdbClient.Set(contx, key, string(v), time.Second*time.Duration(expiration))
		}
	}
}

func (c cacheDb) Get(key string, obj interface{}) bool {
	if redisEnable {
		valueStr, err1 := rdbClient.Get(contx, key).Result()
		if err1 == nil && valueStr != "" {
			err2 := json.Unmarshal([]byte(valueStr), obj)
			return err2 == nil
		}
		return false // 成功在redis中获取了数据
	}
	return false // 不启用redis
}

// 清除缓存
func (c cacheDb) FlushAll() {
	if redisEnable {
		rdbClient.FlushAll(contx)
	}
}

var CacheDb = &cacheDb{}
