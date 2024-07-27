package models

type Access struct {
	Id          int
	ModuleName  string
	ActionName  string
	Type        int
	Url         string
	ModuleId    int
	Sort        int
	Description string
	Status      int
	AddTime     int
	AccessItem  []Access `gorm:"foreignkey:ModuleId;references:Id"`
	Checked     bool     `gorm:"-"` // "-" 表示 忽略本字段
}

func (Access) TableName() string {
	return "access"
}
