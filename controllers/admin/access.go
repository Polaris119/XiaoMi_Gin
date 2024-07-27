package admin

import (
	"github.com/gin-gonic/gin"
	"mian.go/models"
	"net/http"
	"strings"
)

type AccessController struct {
	BaseController
}

func (con AccessController) Index(c *gin.Context) {
	accessList := []models.Access{}
	models.DB.Where("module_id = ?", 0).Preload("AccessItem").Find(&accessList)

	c.HTML(http.StatusOK, "admin/access/index.html", gin.H{
		"accessList": accessList,
	})
}

func (con AccessController) Add(c *gin.Context) {
	// 获取顶级模块
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Find(&accessList)
	c.HTML(http.StatusOK, "admin/access/add.html", gin.H{
		"accessList": accessList,
	})
}

func (con AccessController) DoAdd(c *gin.Context) {
	// 获取表单数据
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	accseeType, err1 := models.Int(c.PostForm("type"))
	actionName := strings.Trim(c.PostForm("action_name"), " ")
	url := c.PostForm("url")
	moduleId, err2 := models.Int(c.PostForm("module_id"))
	sort, err3 := models.Int(c.PostForm("sort"))
	description := c.PostForm("description")
	status, err4 := models.Int(c.PostForm("status"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		con.Error(c, "传入数据有误", "/admin/access/add")
		return
	}

	if moduleName == "" {
		con.Error(c, "模块名不可为空", "/admin/access/add")
		return
	}

	access := models.Access{
		ModuleName:  moduleName,
		Type:        accseeType,
		ActionName:  actionName,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Description: description,
		Status:      status,
	}

	err5 := models.DB.Create(&access).Error
	if err5 != nil {
		con.Error(c, "数据增加失败", "/admin/access/add")
		return
	}

	con.Success(c, "数据增加成功", "/admin/access")

}

func (con AccessController) Edit(c *gin.Context) {
	// 获取要修改的数据
	id, err1 := models.Int(c.Query("id"))

	if err1 != nil {
		con.Error(c, "传入参数有误", "/admin/access")
	}

	access := models.Access{Id: id}
	models.DB.Find(&access)

	//获取顶级模块
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Find(&accessList)

	c.HTML(http.StatusOK, "admin/access/edit.html", gin.H{
		"accessList": accessList,
		"access":     access,
	})
}

func (con AccessController) DoEdit(c *gin.Context) {
	// 获取表单数据
	id, err1 := models.Int(c.PostForm("id"))
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	accseeType, err2 := models.Int(c.PostForm("type"))
	actionName := strings.Trim(c.PostForm("action_name"), " ")
	url := c.PostForm("url")
	moduleId, err3 := models.Int(c.PostForm("module_id"))
	sort, err4 := models.Int(c.PostForm("sort"))
	description := c.PostForm("description")
	status, err5 := models.Int(c.PostForm("status"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		con.Error(c, "传入数据有误", "/admin/access/edit?id="+models.String(id))
		return
	}

	if moduleName == "" {
		con.Error(c, "模块名不可为空", "/admin/access/edit?id="+models.String(id))
		return
	}

	access := models.Access{Id: id}
	models.DB.Find(&access)
	access.ModuleName = moduleName
	access.Type = accseeType
	access.ActionName = actionName
	access.Url = url
	access.ModuleId = moduleId
	access.Sort = sort
	access.Description = description
	access.Status = status

	err := models.DB.Save(&access).Error
	if err != nil {
		con.Error(c, "数据修改失败", "/admin/access/edit?id="+models.String(id))
		return
	} else {
		con.Success(c, "数据修改成功", "/admin/access")
	}
}

func (con AccessController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "数据有误", "/admin/access")
	} else {
		// 获取需要删除的数据
		access := models.Access{Id: id}
		models.DB.Find(&access)
		if access.ModuleId == 0 { // 顶级模块
			accessList := []models.Access{}
			// 查看 属于 顶级模块 的数据 是否存在
			models.DB.Where("module_id=?", access.Id).Find(&accessList)
			if len(accessList) > 0 {
				con.Error(c, "当前模块下存在菜单或操作，谨慎删除", "/admin/access")
			} else {
				models.DB.Delete(&access)
				con.Success(c, "删除数据成功", "/admin/access")
			}
		} else { //菜单/操作 可以直接删除
			models.DB.Delete(&access)
			con.Success(c, "删除数据成功", "/admin/access")
		}
	}
}
