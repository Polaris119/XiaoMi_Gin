package routers

import (
	"github.com/gin-gonic/gin"
	"mian.go/controllers/polaris"
	"mian.go/middlewares"
)

func DefaultRouters(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", polaris.DefaultController{}.Index)
		defaultRouters.GET("/category:id", polaris.ProductController{}.Category)
		defaultRouters.GET("/detail", polaris.ProductController{}.Detail)
		defaultRouters.GET("/product/getImgList", polaris.ProductController{}.GetImgList)

		defaultRouters.GET("/cart", polaris.CartController{}.Get)
		defaultRouters.GET("/cart/addCart", polaris.CartController{}.AddCart)

		defaultRouters.GET("/cart/successTip", polaris.CartController{}.AddCartSuccess)

		defaultRouters.GET("/cart/decCart", polaris.CartController{}.DecCart)
		defaultRouters.GET("/cart/incCart", polaris.CartController{}.IncCart)

		defaultRouters.GET("/cart/changeOneCart", polaris.CartController{}.ChangeOneCart)
		defaultRouters.GET("/cart/changeAllCart", polaris.CartController{}.ChangeAllCart)
		defaultRouters.GET("/cart/delCart", polaris.CartController{}.DelCart)

		defaultRouters.GET("/pass/login", polaris.PassController{}.Login)
		defaultRouters.GET("/pass/captcha", polaris.PassController{}.Captcha)

		defaultRouters.GET("/pass/registerStep1", polaris.PassController{}.RegisterStep1)
		defaultRouters.GET("/pass/registerStep2", polaris.PassController{}.RegisterStep2)
		defaultRouters.GET("/pass/registerStep3", polaris.PassController{}.RegisterStep3)
		defaultRouters.GET("/pass/sendCode", polaris.PassController{}.SendCode)
		defaultRouters.GET("/pass/validateSmsCode", polaris.PassController{}.ValidateSmsCode)
		defaultRouters.POST("/pass/doRegister", polaris.PassController{}.DoRegister)
		defaultRouters.POST("/pass/doLogin", polaris.PassController{}.DoLogin)
		defaultRouters.GET("/pass/loginOut", polaris.PassController{}.LoginOut)
		//判断用户权限
		defaultRouters.GET("/buy/checkout", middlewares.InitUserAuthMiddleware, polaris.BuyController{}.Checkout)
		defaultRouters.POST("/buy/doCheckout", middlewares.InitUserAuthMiddleware, polaris.BuyController{}.DoCheckout)
		defaultRouters.GET("/buy/pay", middlewares.InitUserAuthMiddleware, polaris.BuyController{}.Pay)
		defaultRouters.GET("/buy/orderPayStatus", middlewares.InitUserAuthMiddleware, polaris.BuyController{}.OrderPayStatus)

		defaultRouters.POST("/address/addAddress", middlewares.InitUserAuthMiddleware, polaris.AddressController{}.AddAddress)
		defaultRouters.POST("/address/editAddress", middlewares.InitUserAuthMiddleware, polaris.AddressController{}.EditAddress)
		defaultRouters.GET("/address/changeDefaultAddress", middlewares.InitUserAuthMiddleware, polaris.AddressController{}.ChangeDefaultAddress)
		defaultRouters.GET("/address/getOneAddressList", middlewares.InitUserAuthMiddleware, polaris.AddressController{}.GetOneAddressList)

		defaultRouters.GET("/alipay", middlewares.InitUserAuthMiddleware, polaris.AlipayController{}.Alipay)
		defaultRouters.POST("/alipayNotify", polaris.AlipayController{}.AlipayNotify)
		defaultRouters.GET("/alipayReturn", middlewares.InitUserAuthMiddleware, polaris.AlipayController{}.AlipayReturn)

		defaultRouters.GET("/wxpay", middlewares.InitUserAuthMiddleware, polaris.WxpayController{}.Wxpay)
		defaultRouters.POST("/wxpay/notify", polaris.WxpayController{}.WxpayNotify)

		defaultRouters.GET("/user", middlewares.InitUserAuthMiddleware, polaris.UserController{}.Index)
		defaultRouters.GET("/user/order", middlewares.InitUserAuthMiddleware, polaris.UserController{}.OrderList)
		defaultRouters.GET("/user/orderinfo", middlewares.InitUserAuthMiddleware, polaris.UserController{}.OrderInfo)

		//defaultRouters.GET("/", polaris.DefaultController{}.Index)
	}
}
