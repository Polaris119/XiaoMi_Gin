package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mian.go/models"
	"net/http"
)

type FocusController struct {
	BaseController
}

func (con FocusController) Index(c *gin.Context) {
	focusList := []models.Focus{}
	models.DB.Find(&focusList)
	c.HTML(http.StatusOK, "admin/focus/index.html", gin.H{
		"focusList": focusList,
	})
}

func (con FocusController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/focus/add.html", gin.H{})
}

func (con FocusController) DoAdd(c *gin.Context) {
	title := c.PostForm("title")
	focusType, err1 := models.Int(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, err2 := models.Int(c.PostForm("sort"))
	status, err3 := models.Int(c.PostForm("status"))
	// focusType和sort的数据是写死的，如果出错，只能是非法请求
	if err1 != nil || err3 != nil {
		con.Error(c, "非法请求", "/admin/focus/add")
	}
	if err2 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/add")
	}

	// 上传文件
	focusImgSrc, err4 := models.UploadImg(c, "focus_img")
	//fmt.Println(focusImgSrc)  // 文件格式：  static/upload/20240703/1720017170.jpg
	// 所以在前端  <img src="/{{$value.FocusImg}}" width="200" />  需要在  {{$value.FocusImg}}  前面加  /

	if err4 != nil {
		fmt.Println(err4)
	}
	fmt.Println(focusImgSrc)

	focus := models.Focus{
		Title:     title,
		FocusType: focusType,
		FocusImg:  focusImgSrc,
		Link:      link,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}

	err5 := models.DB.Create(&focus).Error
	if err5 != nil {
		con.Error(c, "增加轮播图失败", "/admin/focus/add")
	} else {
		con.Success(c, "增加轮播图成功", "/admin/focus")
	}
}

func (con FocusController) Edit(c *gin.Context) {
	id, err1 := models.Int(c.Query("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/focus")
		return
	}
	focus := models.Focus{Id: id}
	models.DB.Find(&focus)
	c.HTML(http.StatusOK, "admin/focus/edit.html", gin.H{
		"focus": focus,
	})
}

func (con FocusController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	title := c.PostForm("title")
	focusType, err2 := models.Int(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, err3 := models.Int(c.PostForm("sort"))
	status, err4 := models.Int(c.PostForm("status"))

	if err1 != nil || err2 != nil || err4 != nil {
		con.Error(c, "非法请求", "/admin/focus/edit？id="+models.String(id))
	}
	if err3 != nil {
		con.Error(c, "请输入正确的排序值", "/admin/focus/edit/"+models.String(id))
	}

	focusImg, _ := models.UploadImg(c, "focus_img")

	focus := models.Focus{Id: id}
	models.DB.Find(&focus)

	focus.Title = title
	focus.FocusType = focusType
	focus.Link = link
	focus.Sort = sort
	focus.Status = status

	// 图片地址为空，表示图片不需要修改
	if focusImg != "" {
		focus.FocusImg = focusImg
	}

	err := models.DB.Save(&focus).Error
	if err != nil {
		con.Error(c, "轮播图修改失败", "/admin/focus/edit/"+models.String(id))
	} else {
		con.Success(c, "轮播图修改成功", "/admin/focus")
	}

}

func (con FocusController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/focus")
	} else {
		focus := models.Focus{Id: id}
		models.DB.Delete(&focus)
		// 可以根据需求，考虑删除图片
		// 如果删除图片，可以使用  os.Remove("图片地址")
		con.Success(c, "删除数据成功！", "/admin/focus")
	}
}
