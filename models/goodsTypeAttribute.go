package models

// 为了与 /admin/goods/add.html 中的 response.result 数组里的元素相对应，更改json数据的显示
type GoodsTypeAttribute struct {
	Id        int    `json:"id"`
	CateId    int    `json:"cate_id"`
	Title     string `json:"title"`
	AttrType  int    `json:"attr_type"`
	AttrValue string `json:"attr_value"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	AddTime   int    `json:"add_time"`
}

func (GoodsTypeAttribute) TableName() string {
	return "goods_type_attribute"
}
