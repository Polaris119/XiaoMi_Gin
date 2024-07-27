package admin

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mian.go/models"
	"net/http"
)

type MainController struct {
	BaseController
}

func (con MainController) Index(c *gin.Context) {
	// 获取userinfo对应的session
	session := sessions.Default(c)
	userinfo := session.Get("userinfo")
	// 类型断言  判断 userinfo是否为string
	userinfoStr, ok := userinfo.(string)

	if ok {
		// 1、获取用户信息
		var userinfoStruct []models.Manager
		json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		// 2、获取所有权限
		accessList := []models.Access{}
		// 关联查询排序
		models.DB.Where("module_id = ?", 0).Preload("AccessItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("access.sort DESC")
		}).Order("sort DESC").Find(&accessList)
		// 3、获取当前角色拥有的权限，把权限id放到一个map对象里
		roleAccess := []models.RoleAccess{}
		models.DB.Where("role_id = ?", userinfoStruct[0].RoleId).Find(&roleAccess)
		roleAccessMap := make(map[int]int)
		for _, v := range roleAccess {
			roleAccessMap[v.AccessId] = v.AccessId
		}
		// 4、循环遍历所有的权限数据，判断当前权限id是否在角色权限的map对象中，如果是，给当前数据加入checked属性
		for i := 0; i < len(accessList); i++ {
			if _, ok := roleAccessMap[accessList[i].Id]; ok {
				accessList[i].Checked = true
			}
			for j := 0; j < len(accessList[i].AccessItem); j++ {
				if _, ok := roleAccessMap[accessList[i].AccessItem[j].Id]; ok {
					accessList[i].AccessItem[j].Checked = true
				}
			}
		}

		c.HTML(http.StatusOK, "admin/main/index.html", gin.H{
			"username":   userinfoStruct[0].Username,
			"accessList": accessList,
			"isSuper":    userinfoStruct[0].IsSuper,
		})
	} else {
		c.Redirect(302, "/admin/login")
	}

}

func (con MainController) Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/main/welcome.html", gin.H{})
}

// 公共修改状态的方法
func (con MainController) ChangeStatus(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入的参数错误",
		})
		return
	}
	table := c.Query("table")
	field := c.Query("field")

	//  Exec 用于执行 mysql 语句
	// ABS(status-1)  ABS的作用是  取绝对值
	// 假设status之前是 1 则 ABS(1-1)=0 ; 假设status之前是0  则 ABS(0-1)=1
	//models.DB.Exec("update ? set ? = ABS(?-1) where id = ?", table, field, field, id)这样写 空格会被去掉，选择使用拼接字符串的方法
	err2 := models.DB.Exec("update "+table+" set "+field+" = ABS("+field+"-1) where id = ?", id).Error
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改失败，请重新尝试~",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功！！！",
	})
}

// 公共修改排序值的方法
func (con MainController) ChangeNum(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "传入的参数错误",
		})
		return
	}
	table := c.Query("table")
	field := c.Query("field")
	num := c.Query("num")

	err2 := models.DB.Exec("update "+table+" set "+field+"="+num+" where id = ?", id).Error
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "修改失败，请重新尝试~",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "修改成功！！！",
		})
	}
}

func (con MainController) FlushAll(c *gin.Context) {
	models.CacheDb.FlushAll()
	con.Success(c, "成功清除Redis缓存数据", "/admin")
}
