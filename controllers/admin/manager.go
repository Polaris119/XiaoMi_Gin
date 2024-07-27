package admin

import (
	"github.com/gin-gonic/gin"
	"mian.go/models"
	"net/http"
	"strings"
)

type ManagerController struct {
	BaseController
}

func (con ManagerController) Index(c *gin.Context) {

	managerList := []models.Manager{}
	models.DB.Preload("Roles").Find(&managerList)

	c.HTML(http.StatusOK, "admin/manager/index.html", gin.H{
		"managerList": managerList,
	})
}

func (con ManagerController) Add(c *gin.Context) {
	// 获取所有的角色
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	c.HTML(http.StatusOK, "admin/manager/add.html", gin.H{
		"roleList": roleList,
	})
}

func (con ManagerController) DoAdd(c *gin.Context) {
	roleId, err1 := models.Int(c.PostForm("role_id"))
	if err1 != nil {
		con.Error(c, "数据有误", "/admin/manager/add")
		return
	}

	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")
	email := strings.Trim(c.PostForm("email"), " ")

	// 判断用户名和密码的长度是否符合要求
	if len(username) < 2 || len(password) < 6 {
		con.Error(c, "用户名或密码长度不合法", "/admin/manager/add")
		return
	}

	// 判断管理是否存在
	managerList := []models.Manager{}
	models.DB.Where("username = ?", username).Find(&managerList)
	if len(managerList) > 0 {
		con.Error(c, "该管理员已存在", "/admin/manager/add")
		return
	}

	// 执行 增加管理员 操作
	manager := models.Manager{
		Username: username,
		Password: models.Md5(password), //密码加密
		Mobile:   mobile,
		Email:    email,
		RoleId:   roleId,
		Status:   1,
		AddTime:  int(models.GetUnix()),
	}
	err2 := models.DB.Create(&manager).Error
	if err2 != nil {
		con.Error(c, "增加管理员失败", "/admin/manager/add")
		return
	}

	con.Success(c, "增加管理员成功", "/admin/manager")
}

func (con ManagerController) Edit(c *gin.Context) {
	// 获取管理员
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/manager")
	}
	manager := models.Manager{Id: id}
	models.DB.Find(&manager)

	// 获取所有角色
	roleList := []models.Role{}
	models.DB.Find(&roleList)

	c.HTML(http.StatusOK, "admin/manager/edit.html", gin.H{
		"manager":  manager,
		"roleList": roleList,
	})
}

func (con ManagerController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误1", "/admin/manager")
		return
	}

	roleId, err2 := models.Int(c.PostForm("role_id"))
	if err2 != nil {
		con.Error(c, "数据有误2", "/admin/manager")
		return
	}

	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	if len(mobile) > 11 {
		con.Error(c, "手机号长度不合法", "/admin/manager/edit?id="+models.String(id))
		return
	}

	// 执行修改
	manager := models.Manager{Id: id}
	models.DB.Find(&manager)
	manager.Username = username
	manager.Email = email
	manager.Mobile = mobile
	manager.RoleId = roleId

	// 判断密码是否为空  为空表示不修改密码  不为空表示修改密码

	if password != "" {
		// 判断密码长度是否合法
		if len(password) < 6 {
			con.Error(c, "密码长度不能小于6位，请重新设置密码", "/admin/manager/edit?id="+models.String(id))
			return
		}
		manager.Password = models.Md5(password)
	}

	err3 := models.DB.Save(&manager).Error
	if err3 != nil {
		con.Error(c, "数据修改失败", "/admin/manager/edit?id="+models.String(id))
		return
	}
	con.Success(c, "数据修改成功", "/admin/manager")
}

func (con ManagerController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "数据有误", "/admin/manager")
	} else { //执行删除
		manager := models.Manager{Id: id}
		models.DB.Delete(&manager)
		con.Success(c, "删除数据成功", "/admin/manager")
	}
}
