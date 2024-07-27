package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"mian.go/models"
	"os"
	"strings"
)

func InitAdminAuthMiddleware(c *gin.Context) {
	//fmt.Println("InitAdminAuthMiddleware")

	//进行权限判断 ==》  没有登录的用户  不能进入后台管理中心

	// 1、获取url访问的地址
	// 使用  strings.Split  分割字符串   （ /admin/captcha?t=0.8706946438889653  需要被  分割成  /admin/captcha）
	pathname := strings.Split(c.Request.URL.String(), "?")[0]

	// 2、获取session保存的用户信息
	session := sessions.Default(c)
	userinfo := session.Get("userinfo")

	// 3、判断session中的用户信息是否存在  如果不存在就跳转到登录页面（注意需要判断）  如果存在继续便向下执行
	// 类型断言 判断 userinfo 是不是 一个string
	userinfoStr, ok := userinfo.(string)

	if ok {
		var userinfoStruct []models.Manager
		// 将 json字符串 转换成 结构体数据
		err := json.Unmarshal([]byte(userinfoStr), &userinfoStruct)

		// json数据成功转换  或者  用户信息不存在（ userinfoStr存在数据  并且 userinfoStr的username的数据  不是空 ）
		if err != nil || !(len(userinfoStr) > 0 && userinfoStruct[0].Username != "") {
			// 判断当前访问的url是否是【login  doLogin captcha】中的一个 ， 如果不是，跳转到登录页面， 如果是，不执行任何操作
			if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
				c.Redirect(302, "/admin/login")
			}
		} else { // 用户登录成功，需要进行权限判断
			// 因为 pathname 的格式为 "/admin/manager"  而数据库保存的路径为  "manager"  所以需要将"/admin/"替换为空
			// urlPath为当前访问的url
			urlPath := strings.Replace(pathname, "/admin/", "", 1)
			// 不是管理员  并且  不属于“不需要权限判断”的路径
			// 对于 是否属于“不需要权限判断”路径 的判断：  格式为"/welcome"  所以需要给urlPath的前面加上"/"
			if userinfoStruct[0].IsSuper == 0 && !excludeAuthPath("/"+urlPath) {
				// 1、 根据角色获取当前角色的权限列表， 然后把权限id放在一个map对象中
				roleAccess := []models.RoleAccess{}
				models.DB.Where("role_id = ?", userinfoStruct[0].RoleId).Find(&roleAccess)
				roleAccessMap := make(map[int]int)
				for _, v := range roleAccess {
					roleAccessMap[v.AccessId] = v.AccessId
				}
				// 2、获取 当前访问的url对应的权限id ； 判断权限id是否在角色id的map中
				access := models.Access{}
				models.DB.Where("url = ?", urlPath).Find(&access)
				// 3、 判断 当前访问的url对应的权限id  是否在角色id的map中
				// 不在map中，则代表  无法访问
				if _, ok := roleAccessMap[access.Id]; !ok {
					c.String(200, "你没有权限访问该网页")
					c.Abort() // 终止当前请求
				}
			}
		}
	} else { // 没有登录的情况（session不存在）
		if pathname != "/admin/login" && pathname != "/admin/doLogin" && pathname != "/admin/captcha" {
			c.Redirect(302, "/admin/login")
		}
	}

}

// 排除 “不需要权限判断” 的路径
// 不需要权限判断  则返回 true
func excludeAuthPath(url string) bool {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Println("app.ini文件读取失败", iniErr)
		os.Exit(1)
	}
	// 读取配置文件中的路径
	excludeAuthPath := config.Section("").Key("excludeAuthPath").String()
	// 将路径放到切片中，方便判断
	excludeAuthPathSlice := strings.Split(excludeAuthPath, ",")
	for _, v := range excludeAuthPathSlice {
		if v == url {
			return true
		}
	}
	return false
}
