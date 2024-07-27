package models

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

// 定义结构体  缓存结构体 私有
type ginCookie struct{}

// 写入数据的方法
func (cookie ginCookie) Set(c *gin.Context, key string, value interface{}) {

	bytes, _ := json.Marshal(value)
	//des加密
	desKey := []byte("polarism") //注意：key必须是8位
	encData, _ := DesEncrypt(bytes, desKey)

	host := c.Request.Host
	// 解析主机名部分
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	c.SetCookie(key, string(encData), 3600*24*30, "/", host, false, true)

	//c.SetCookie(key, string(encData), 3600*24*30, "/", c.Request.Host, false, true)
	//fmt.Println(c.Request.Host)
}

// 获取数据的方法
func (cookie ginCookie) Get(c *gin.Context, key string, obj interface{}) bool {

	valueStr, err1 := c.Cookie(key)
	if err1 == nil && valueStr != "" && valueStr != "[]" {
		//des解密
		desKey := []byte("polarism") //注意：key必须是8位
		decData, e := DesDecrypt([]byte(valueStr), desKey)
		if e != nil {
			return false
		} else {
			err2 := json.Unmarshal([]byte(decData), obj)
			return err2 == nil
		}

	}
	return false
}
func (cookie ginCookie) Remove(c *gin.Context, key string) bool {
	host := c.Request.Host
	// 解析主机名部分
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	c.SetCookie(key, "", -1, "/", host, false, true)

	//c.SetCookie(key, "", -1, "/", c.Request.Host, false, true)
	return true
}

// 实例化结构体
var Cookie = &ginCookie{}
