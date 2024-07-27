package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"mian.go/models"
	"mian.go/routers"
)

func main() {
	//创建默认路由
	r := gin.Default()

	// 自定义模板函数  注意要把这个函数放在加载模板前
	r.SetFuncMap(template.FuncMap{
		"UnixToTime": models.UnixToTime,
		"Str2Html":   models.Str2Html,
		"FormatImg":  models.FormatImg,
		"Sub":        models.Sub,
		"Mul":        models.Mul,
		"Substr":     models.Substr,
		"FormatAttr": models.FormatAttr,
	})

	//加载模板  放在配置路由前面
	r.LoadHTMLGlob("templates/**/**/*")

	//配置静态web目录
	r.Static("/static", "./static")

	//创建基于cookie的存储引擎， secret111 参数是用于加密的密钥
	store := cookie.NewStore([]byte("secret111"))
	// 配置session的中间件 store是前面的存储引擎， 可以替换成其它的存储引擎
	r.Use(sessions.Sessions("mysession", store))

	//配置路由
	routers.AdminRoutersInit(r)
	routers.DefaultRouters(r)

	r.Run()
}
