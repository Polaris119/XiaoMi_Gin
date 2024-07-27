package test

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/go_image"
	qrcode "github.com/skip2/go-qrcode"
	"net/http"
	"os"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func (con DefaultController) Thumbnail1(c *gin.Context) {
	// 按宽度进行比例缩放，输入输出都是文件
	filename := "static/1.jpg"
	savepath := "static/1_1000.jpg"
	err := go_image.ScaleF2F(filename, savepath, 1000)
	if err != nil {
		c.String(http.StatusOK, "Thumbnail1失败")
		return
	}
	c.String(http.StatusOK, "Thumbnail1成功")
}

func (con DefaultController) Thumbnail2(c *gin.Context) {
	filename := "static/1.jpg"
	savepath := "static/1_500_400.jpg"
	// 按宽度及高度进行比例缩放，输入输出都是文件
	err := go_image.ThumbnailF2F(filename, savepath, 500, 400)
	if err != nil {
		c.String(http.StatusOK, "Thumbnail2失败")
		return
	}
	c.String(http.StatusOK, "Thumbnail2成功")
}

func (con DefaultController) Qrcode1(c *gin.Context) {
	var png []byte
	png, err := qrcode.Encode("https://www.qq.com", qrcode.Medium, 256)
	if err != nil {
		c.String(http.StatusOK, "二维码生成失败！！！")
		return
	}
	// string()  将切片转换成图片
	c.String(http.StatusOK, string(png))
}

// 显示二维码的同时，将二维码保存到本地
func (con DefaultController) Qrcode2(c *gin.Context) {
	savepath := "static/2.jpg"
	err := qrcode.WriteFile("https://www.qq.com", qrcode.Medium, 256, savepath)
	if err != nil {
		c.String(http.StatusOK, "二维码保存失败！！！")
		return
	}
	file, _ := os.ReadFile(savepath)
	c.String(http.StatusOK, string(file))
}
