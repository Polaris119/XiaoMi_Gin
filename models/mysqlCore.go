package models

import (
	"fmt"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB
var err error

func init() {
	// 读取 .ini里的数据库配置

	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Println("%v 加载失败", iniErr)
		os.Exit(1)
	}

	ip := config.Section("mysql").Key("ip").String()
	port := config.Section("mysql").Key("port").String()
	user := config.Section("mysql").Key("user").String()
	password := config.Section("mysql").Key("password").String()
	database := config.Section("mysql").Key("database").String()

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, ip, port, database)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		fmt.Println(err)
	}
}
