package admin

import (
	"github.com/gin-gonic/gin"
	"mian.go/models"
	"net/http"
	"strings"
)

type GoodsTypeAttributeController struct {
	BaseController
}

func (con GoodsTypeAttributeController) Index(c *gin.Context) {
	cateId, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数1不正确", "/admin/goodsType")
		return
	}
	// 获取商品类型属性
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	models.DB.Where("cate_id=?", cateId).Find(&goodsTypeAttributeList)
	// 获取 商品类型属性 对应的类型
	goodsType := models.GoodsType{}
	models.DB.Where("id=?", cateId).Find(&goodsType)

	c.HTML(http.StatusOK, "admin/goodsTypeAttribute/index.html", gin.H{
		"cateId":                 cateId,
		"goodsTypeAttributeList": goodsTypeAttributeList,
		"goodsType":              goodsType,
	})
}

func (con GoodsTypeAttributeController) Add(c *gin.Context) {
	// 获取  当前商品类型属性  对应的  类型id
	cateId, err := models.Int(c.Query("cate_id"))
	if err != nil {
		con.Error(c, "传入参数2不正确", "/admin/goodsType")
		return
	}

	// 获取所有商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	c.HTML(http.StatusOK, "admin/goodsTypeAttribute/add.html", gin.H{
		"cateId":        cateId,
		"goodsTypeList": goodsTypeList,
	})
}

func (con GoodsTypeAttributeController) DoAdd(c *gin.Context) {
	cateId, err1 := models.Int(c.PostForm("cate_id"))
	title := strings.Trim(c.PostForm("title"), " ")
	attrType, err2 := models.Int(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, err3 := models.Int(c.PostForm("sort"))

	if err1 != nil || err2 != nil {
		con.Error(c, "add传入参数有问题", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}

	if title == "" {
		con.Error(c, "商品类型属性名称不能为空", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}

	if err3 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
		return
	}

	goodsTypeAttribute := models.GoodsTypeAttribute{
		Title:     title,
		CateId:    cateId,
		AttrType:  attrType,
		AttrValue: attrValue,
		Sort:      sort,
	}

	err := models.DB.Create(&goodsTypeAttribute).Error
	if err != nil {
		con.Error(c, "增加数据失败", "/admin/goodsTypeAttribute/add?cate_id="+models.String(cateId))
	} else {
		con.Success(c, "增加数据成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
	}
}

func (con GoodsTypeAttributeController) Edit(c *gin.Context) {
	// 获取当前需要修改的数据
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数3不正确", "/admin/goodsType")
		return
	}
	// 获取当前id对应的商品类型属性
	goodsTypeAttribute := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttribute)

	// 获取所有商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	// 模板路径应当是相对于你在 LoadHTMLGlob 中指定的模板目录的相对路径，而不是绝对路径。
	// 所以  admin/goodsTypeAttribute/edit.html  中 admin 的前面无需加 /
	c.HTML(http.StatusOK, "admin/goodsTypeAttribute/edit.html", gin.H{
		"goodsTypeAttribute": goodsTypeAttribute,
		"goodsTypeList":      goodsTypeList,
	})
}

func (con GoodsTypeAttributeController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	title := strings.Trim(c.PostForm("title"), " ")
	cateId, err2 := models.Int(c.PostForm("cate_id"))
	attrType, err3 := models.Int(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, err4 := models.Int(c.PostForm("sort"))

	if err1 != nil || err2 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/goodsType")
		return
	}
	if title == "" {
		con.Error(c, "商品类型属性名称不能为空", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}
	if err4 != nil {
		con.Error(c, "排序值不对", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}

	goodsTypeAttr := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&goodsTypeAttr)
	goodsTypeAttr.Title = title
	goodsTypeAttr.CateId = cateId
	goodsTypeAttr.AttrType = attrType
	goodsTypeAttr.AttrValue = attrValue
	goodsTypeAttr.Sort = sort
	err := models.DB.Save(&goodsTypeAttr).Error
	if err != nil {
		// 失败，则跳转到 商品类型属性 编辑 页面  /admin/goodsTypeAttribute?id=商品类型属性的id
		con.Error(c, "修改数据失败", "/admin/goodsTypeAttribute/edit?id="+models.String(id))
		return
	}
	// 成功，则跳转到 商品类型展示页面  /admin/goodsTypeAttribute?id=商品类型的id
	con.Success(c, "需改数据成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
}

func (con GoodsTypeAttributeController) Delete(c *gin.Context) {
	// 用于删除数据（id用于标记  商品类型属性）
	id, err1 := models.Int(c.Query("id"))
	// 用于页面返回（cate_id用于标记  商品类型）
	cateId, err2 := models.Int(c.Query("cate_id"))
	if err1 != nil || err2 != nil {
		con.Error(c, "传入参数错误", "/admin/goodsType")
	} else {
		goodsTypeAttr := models.GoodsTypeAttribute{Id: id}
		models.DB.Delete(&goodsTypeAttr)
		con.Success(c, "删除数据成功", "/admin/goodsTypeAttribute?id="+models.String(cateId))
	}
}
