package models

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

//// 创建store
//var store = base64Captcha.DefaultMemStore

// 配置RedisStore
// RedisStore结构体实现了base64Captcha.Store接口中的所有方法
// 变量store是一个实现了base64Captcha.Store接口的RedisStore实例
var store base64Captcha.Store = RedisStore{}

// 获取验证码
func MakeCaptcha(height int, width int, length int) (string, string, error) {
	var driver base64Captcha.Driver
	driverString := base64Captcha.DriverString{
		Height:          height,
		Width:           width,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          length,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driver = driverString.ConvertFonts()

	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := c.Generate()
	return id, b64s, err

}

// 验证验证码
func VerifyCaptcha(id string, VerifyValue string) bool {
	fmt.Println(id, VerifyValue)
	if store.Verify(id, VerifyValue, true) {
		return true
	} else {
		return false
	}
}
