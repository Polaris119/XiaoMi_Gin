package models

type Nav struct {
	Id         int
	Title      string
	Link       string
	Position   int // 1表示顶级导航  2表示中间导航
	IsOpennew  int
	Relation   string
	Sort       int
	Status     int
	AddTime    int
	GoodsItems []Goods `gorm:"-"`
}

func (Nav) TableName() string {
	return "nav"
}
