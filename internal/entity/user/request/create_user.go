package request

type CreateUserReq struct {
	UserUuid      string `form:"userUuid" json:"userUuid" example:"not used in request"`
	CustomersUuid string `form:"customersUuid" json:"customersUuid" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	Email         string `form:"email" json:"email" example:"user1@gmail.com"`
	UserName      string `form:"userName" json:"userName" example:"ivanovAA"`
	FirstName     string `form:"firstName" json:"firstName" example:"Anton"`
	LastName      string `form:"lastName" json:"lastName" example:"Ivanov"`
	MiddleName    string `form:"middleName" json:"middleName" example:"Antonovich"`
	Phone         string `form:"phone" json:"phone" example:"+79091234567"`
	Enabled       bool   `form:"enabled" json:"enabled" example:"true"`
	Password      string `form:"password" json:"password" example:"qwerty123"`
	CustomerRoles int    `form:"customerRoles" json:"customerRoles" example:"3"`
	Dopinfo       string `form:"dopinfo" json:"dopinfo" example:"vendor"`
	CreateUser    string `form:"createName" json:"createName" example:"PetrovAA"`
	//CreateTime    time.Time `form:"createTime" json:"createTime" example:"2022-09-23 13:17:56.079"`
	UpdateUser string `form:"updateName" json:"updateName" example:"PetrovAA"`
	//UpdateTime    time.Time `form:"updateTime" json:"updateTime" example:"2022-09-23 13:17:56.079"`

	//Groups        []string `form:"groups" example:"["/Yandex LLC"]"`
}
