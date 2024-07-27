package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"mian.go/models"
	"net/http"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type GoodsController struct {
	BaseController
}

func (con GoodsController) Index(c *gin.Context) {
	// 当前页数
	page, _ := models.Int(c.Query("page")) // 将 page为空 的情况，转换成 page=0
	if page == 0 {
		page = 1
	}
	//条件
	where := "is_delete = 0"
	// 获取keyword
	keyword := c.Query("keyword")
	if len(keyword) > 0 {
		where += " AND title like \"%" + keyword + "%\""
	}

	// 每页查询的数量
	pageSize := 5
	goodsList := []models.Goods{}
	models.DB.Where(where).Offset((page - 1) * pageSize).Limit(pageSize).Find(&goodsList)

	// 获取总数量
	var count int64
	models.DB.Where(where).Table("goods").Count(&count)

	//判断最后一页有没有数据 如果没有跳转到第一页
	if len(goodsList) > 0 { // 页面中 存在商品 则直接渲染
		c.HTML(http.StatusOK, "admin/goods/index.html", gin.H{
			"goodsList": goodsList,
			//注意float64类型
			"totalPages": math.Ceil(float64(count) / float64(pageSize)),
			"page":       page,
			"keyword":    keyword,
		})
	} else { // 页面中 不存在商品
		if page != 1 { // 若是页面不是第一页  则跳转到第一页
			c.Redirect(302, "/admin/goods")
		} else { // 若页面是第一页，则正常渲染
			c.HTML(http.StatusOK, "admin/goods/index.html", gin.H{
				"goodsList": goodsList,
				//注意float64类型
				"totalPages": math.Ceil(float64(count) / float64(pageSize)),
				"page":       page,
				"keyword":    keyword,
			})
		}

	}
}

func (con GoodsController) Add(c *gin.Context) {
	// 获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	// 获取所有颜色信息
	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)

	// 获取商品规格包装
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "admin/goods/add.html", gin.H{
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
	})
}

func (con GoodsController) DoAdd(c *gin.Context) {
	//1、获取表单提交过来的数据 进行判断

	title := c.PostForm("title")                             // 商品通用信息  商品标题
	subTitle := c.PostForm("sub_title")                      // 商品通用信息  附属标题
	goodsSn := c.PostForm("goods_sn")                        // 商品通用信息
	cateId, _ := models.Int(c.PostForm("cate_id"))           // 商品通用信息  所属分类
	goodsNumber, _ := models.Int(c.PostForm("goods_number")) //  库存
	//价格  注意小数点  （数据库中保存的是float64形式）
	marketPrice, _ := models.Float(c.PostForm("market_price")) // 商品通用信息  商品原价
	price, _ := models.Float(c.PostForm("price"))              // 商品通用信息  商品价格

	relationGoods := c.PostForm("relation_goods") // 商品属性  关联商品
	goodsAttr := c.PostForm("goods_attr")         // 商品属性  更多属性
	goodsVersion := c.PostForm("goods_version")   // 商品通用信息  商品版本
	goodsGift := c.PostForm("goods_gift")         // 商品属性  关联赠品
	goodsFitting := c.PostForm("goods_fitting")   // 商品属性  关联配件
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color") // 商品属性  商品颜色

	goodsKeywords := c.PostForm("goods_keywords") // 商品属性  SEO关键词
	goodsDesc := c.PostForm("goods_desc")         // 商品属性  SEO描述
	goodsContent := c.PostForm("goods_content")   // 详细描述
	isDelete, _ := models.Int(c.PostForm("is_delete"))
	isHot, _ := models.Int(c.PostForm("is_hot"))              // 商品通用信息  热销
	isBest, _ := models.Int(c.PostForm("is_best"))            // 商品通用信息  精品
	isNew, _ := models.Int(c.PostForm("is_new"))              // 商品通用信息  新品
	goodsTypeId, _ := models.Int(c.PostForm("goods_type_id")) // 规格与包装  商品类型
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status")) // 商品通用信息  商品状态
	addTime := int(models.GetUnix())              // 添加时间

	// 2、获取颜色信息  将颜色转化成字符串
	goodsColorStr := strings.Join(goodsColorArr, ",")

	// 3、上传图片  生成缩略图
	goodsImg, _ := models.UploadImg(c, "goods_img")
	if len(goodsImg) > 0 {
		//判断 本地图片才需要处理
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(goodsImg)
				wg.Done()
			}()
		}

	}

	// 4、增加商品数据  放到 goods表
	goods := models.Goods{
		Title:         title,         // 商品标题
		SubTitle:      subTitle,      // 附属标题
		GoodsSn:       goodsSn,       // 商品SN号
		CateId:        cateId,        // 商品所属分类（除了顶级分类，还有别的分类）
		ClickCount:    100,           // 点击量
		GoodsNumber:   goodsNumber,   // 商品库存
		MarketPrice:   marketPrice,   // 商品市场价
		Price:         price,         // 商品当前价格
		RelationGoods: relationGoods, // 关联商品
		GoodsAttr:     goodsAttr,     // 商品更多属性
		GoodsVersion:  goodsVersion,  // 商品版本
		GoodsGift:     goodsGift,     // 商品赠品
		GoodsFitting:  goodsFitting,  // 商品配件
		GoodsKeywords: goodsKeywords, // 商品关键词
		GoodsDesc:     goodsDesc,     // 商品描述
		GoodsContent:  goodsContent,  // 商品详情
		IsDelete:      isDelete,      // 商品是否删除
		IsHot:         isHot,         // 商品是否热销
		IsBest:        isBest,        // 商品是否精品
		IsNew:         isNew,         // 商品是否新品
		GoodsTypeId:   goodsTypeId,   // 商品分类（顶级分类）
		Sort:          sort,          // 商品排序
		Status:        status,        // 商品状态
		AddTime:       addTime,       // 商品增加时间
		GoodsColor:    goodsColorStr, // 商品颜色
		GoodsImg:      goodsImg,      // 商品图片
	}
	err := models.DB.Create(&goods).Error
	if err != nil {
		con.Error(c, "增加失败", "/admin/goods/add")
	}
	// 5、增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id // 先创建goods表，等goods表创建完毕，就能得到goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10 // sort和status设为默认值
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()

	// 6、增加规格包装
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrValueList); i++ {
			goodsTypeAttributeId, attributeIdErr := models.Int(attrIdList[i]) // 获取商品类型属性Id
			if attributeIdErr == nil {
				// 根据id 获取  商品类型属性  对应的数据
				goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
				models.DB.Find(&goodsTypeAttributeObj)
				// 给商品属性里面增加  规格包装  数据
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id                             // 用于 与goods表关联
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title   // 商品类型属性  的标题
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType // 商品类型属性  的属性分类（AttrType=1 ：单行文本框；AttrType=2：多行文本框）
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id         // 商品类型属性  的id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId // 商品类型属性  的cate_id（用于判断 类型属性 属于哪个顶级分类）
				goodsAttrObj.AttributeValue = attrValueList[i]              // 商品类型属性  的下拉框的值（AttrType=3）
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}
		}
		wg.Done()
	}()
	con.Success(c, "增加数据成功", "/admin/goods")
}

func (con GoodsController) Edit(c *gin.Context) {
	// 1、获取要修改的商品数据
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "edit传入参数错误", "/admin/goods/edit")
		return
	}
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)

	// 2、获取商品分类
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	// 3、获取所有颜色 以及 需要选中的颜色
	goodsColorSlice := strings.Split(goods.GoodsColor, ",") // 需要选中的颜色，先转换成切片
	goodsColorMap := make(map[string]string)                // 将需要选中的颜色放到 映射中，方便与所有颜色作对比，从而设置checked字段
	for _, v := range goodsColorSlice {
		goodsColorMap[v] = v
	}

	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList) // 获取所有颜色
	for i := 0; i < len(goodsColorList); i++ {
		if _, ok := goodsColorMap[models.String(goodsColorList[i].Id)]; ok {
			goodsColorList[i].Checked = true
		}
	}

	// 4、获取商品图库信息
	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id=?", id).Find(&goodsImageList)

	// 5、获取商品类型
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	// 6、获取商品规格信息
	goodsAttr := []models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsAttr) // 获取商品属性（商品属性类型）
	goodsAttrStr := ""

	// 根据商品属性类型，判断是 单行文本框、多行文本框、下拉框
	for _, v := range goodsAttr {
		if v.AttributeType == 1 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: </span> <input type="hidden" name="attr_id_list" value="%v" />   <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else if v.AttributeType == 2 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: 　</span><input type="hidden" name="attr_id_list" value="%v" />  <textarea cols="50" rows="3" name="attr_value_list">%v</textarea></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else {
			// 获取 下拉框 当前对应的值
			goodsTypeAttribute := models.GoodsTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&goodsTypeAttribute) // 根据属性id  获取商品 类型属性 表
			attrValueSlice := strings.Split(goodsTypeAttribute.AttrValue, "\n")

			goodsAttrStr += fmt.Sprintf(`<li><span>%v: 　</span>  <input type="hidden" name="attr_id_list" value="%v" /> `, v.AttributeTitle, v.AttributeId)
			goodsAttrStr += fmt.Sprintf(`<select name="attr_value_list">`)
			for i := 0; i < len(attrValueSlice); i++ {
				if attrValueSlice[i] == v.AttributeValue {
					goodsAttrStr += fmt.Sprintf(`<option value="%v" selected >%v</option>`, attrValueSlice[i], attrValueSlice[i])
				} else {
					goodsAttrStr += fmt.Sprintf(`<option value="%v">%v</option>`, attrValueSlice[i], attrValueSlice[i])
				}
			}
			goodsAttrStr += fmt.Sprintf(`</select>`)
			goodsAttrStr += fmt.Sprintf(`</li>`)
		}
	}

	c.HTML(http.StatusOK, "admin/goods/edit.html", gin.H{
		"goods":          goods,
		"goodsCateList":  goodsCateList,
		"goodsImageList": goodsImageList,
		"goodsTypeList":  goodsTypeList,
		"goodsColorList": goodsColorList,
		"goodsAttrStr":   goodsAttrStr,
		"prevPage":       c.Request.Referer(), // 获取上一页地址
	})
}

func (con GoodsController) DoEdit(c *gin.Context) {
	//1、获取表单提交过来的数据
	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/admin/goods")
	}
	// 获取上一页地址
	prevPage := c.PostForm("prevPage")

	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	goodsSn := c.PostForm("goods_sn")
	cateId, _ := models.Int(c.PostForm("cate_id"))
	goodsNumber, _ := models.Int(c.PostForm("goods_number"))
	//注意小数点
	marketPrice, _ := models.Float(c.PostForm("market_price"))
	price, _ := models.Float(c.PostForm("price"))
	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr")
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	//获取的是切片
	goodsColorArr := c.PostFormArray("goods_color")
	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	isDelete, _ := models.Int(c.PostForm("is_delete"))
	isHot, _ := models.Int(c.PostForm("is_hot"))
	isBest, _ := models.Int(c.PostForm("is_best"))
	isNew, _ := models.Int(c.PostForm("is_new"))
	goodsTypeId, _ := models.Int(c.PostForm("goods_type_id"))
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	//2、获取颜色信息 把颜色转化成字符串
	goodsColorStr := strings.Join(goodsColorArr, ",")
	//3、修改数据
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	goods.Title = title
	goods.SubTitle = subTitle
	goods.GoodsSn = goodsSn
	goods.CateId = cateId
	goods.GoodsNumber = goodsNumber
	goods.MarketPrice = marketPrice
	goods.Price = price
	goods.RelationGoods = relationGoods
	goods.GoodsAttr = goodsAttr
	goods.GoodsVersion = goodsVersion
	goods.GoodsGift = goodsGift
	goods.GoodsFitting = goodsFitting
	goods.GoodsKeywords = goodsKeywords
	goods.GoodsDesc = goodsDesc
	goods.GoodsContent = goodsContent
	goods.IsDelete = isDelete
	goods.IsHot = isHot
	goods.IsBest = isBest
	goods.IsNew = isNew
	goods.GoodsTypeId = goodsTypeId
	goods.Sort = sort
	goods.Status = status
	goods.GoodsColor = goodsColorStr

	//4、上传图片   生成缩略图
	goodsImg, err2 := models.UploadImg(c, "goods_img")
	if err2 == nil && len(goodsImg) > 0 {
		goods.GoodsImg = goodsImg
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(goodsImg)
				wg.Done()
			}()
		}
	}

	err3 := models.DB.Save(&goods).Error
	if err3 != nil {
		con.Error(c, "修改失败", "/admin/goods/edit?id="+models.String(id))
		return
	}

	//5、修改图库 增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()
	// 6、修改规格包装  1、删除当前商品下面的规格包装   2、重新执行增加

	// 6.1删除当前商品下面的规格包装
	goodsAttrObj := models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Delete(&goodsAttrObj)
	// 6.2、重新执行增加
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, attributeIdErr := models.Int(attrIdList[i])
			if attributeIdErr == nil {
				//获取商品类型属性的数据
				goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
				models.DB.Find(&goodsTypeAttributeObj)
				//给商品属性里面增加数据  规格包装
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
				goodsAttrObj.AttributeValue = attrValueList[i]
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}

		}
		wg.Done()
	}()
	wg.Wait()
	if len(prevPage) > 0 {
		con.Success(c, "修改数据成功", prevPage)
	} else {
		con.Success(c, "修改数据成功", "/admin/goods")
	}

}

// 配置 商品类型属性
func (con GoodsController) GoodsTypeAttribute(c *gin.Context) {
	cateId, err1 := models.Int(c.Query("cateId"))
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	err2 := models.DB.Where("cate_id=?", cateId).Find(&goodsTypeAttributeList).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"result":  "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  goodsTypeAttributeList,
		})
	}
}

// 富文本编辑器上传图片
func (con GoodsController) EditorImageUpload(c *gin.Context) {
	//上传图片
	imgDir, err := models.UploadImg(c, "file") //注意：可以在网络里面看到传递的参数
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
	} else {
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(imgDir)
				wg.Done()
			}()
			c.JSON(http.StatusOK, gin.H{
				"link": "/" + imgDir,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"link": models.GetSettingFromColumn("OssDomain") + imgDir,
			})
		}

	}
}

// 图库上传图片
func (con GoodsController) GoodsImageUpload(c *gin.Context) {
	//上传图片
	imgDir, err := models.UploadImg(c, "file") //注意：可以在网络里面看到传递的参数
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
	} else {
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(imgDir)
				wg.Done()
			}()

		}
		c.JSON(http.StatusOK, gin.H{
			"link": imgDir,
		})

	}
}

// 修改商品图库关联的颜色
func (con GoodsController) ChangeGoodsImageColor(c *gin.Context) {
	//获取图片id 获取颜色id
	goodsImageId, err1 := models.Int(c.Query("goods_image_id"))
	colorId, err2 := models.Int(c.Query("color_id"))
	goodsImage := models.GoodsImage{Id: goodsImageId}
	models.DB.Find(&goodsImage)
	goodsImage.ColorId = colorId
	err3 := models.DB.Save(&goodsImage).Error
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  "更新失败",
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result":  "更新成功",
			"success": true,
		})
	}
}

// 删除图库
func (con GoodsController) RemoveGoodsImage(c *gin.Context) {
	//获取图片id
	goodsImageId, err1 := models.Int(c.Query("goods_image_id"))
	goodsImage := models.GoodsImage{Id: goodsImageId}
	err2 := models.DB.Delete(&goodsImage).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  "删除失败",
			"success": false,
		})
	} else {
		//删除服务器图片
		// os.Remove()
		c.JSON(http.StatusOK, gin.H{
			"result":  "删除成功",
			"success": true,
		})
	}

}

// 删除数据
func (con GoodsController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/goods")
	} else {
		goods := models.Goods{Id: id}
		models.DB.Find(&goods)
		goods.IsDelete = 1
		goods.Status = 0
		models.DB.Save(&goods)
		//获取上一页
		prevPage := c.Request.Referer()
		if len(prevPage) > 0 {
			con.Success(c, "删除数据成功", prevPage)
		} else {
			con.Success(c, "删除数据成功", "/admin/goods")
		}
	}
}
