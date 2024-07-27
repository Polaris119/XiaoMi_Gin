package models

type Manager struct {
	Id       int
	Username string
	Password string
	Mobile   string
	Email    string
	Status   int
	RoleId   int
	AddTime  int
	IsSuper  int
	Roles    Role `gorm:"foreignkey:RoleId;references:Id"`
}

func (Manager) TableName() string {
	return "manager"
}
