package request

type UpdateUserReq struct {
	UserUuid      string `form:"userUuid" json:"userUuid" example:"c5074793-9d82-478c-9853-125c04bdb626"`
	CustomersUuid string `form:"customersUuid" json:"customersUuid" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	Email         string `form:"email" example:"user1@gmail.com"`
	UserName      string `form:"userName" example:"ivanovAA"`
	FirstName     string `form:"firstName" example:"Anton"`
	LastName      string `form:"lastName" example:"Ivanov"`
	MiddleName    string `form:"middleName" example:"Antonovich"`
	Phone         string `form:"phone" example:"+79091234567"`
	CustomerRoles int    `form:"customerRoles" json:"customerRoles" example:"3"`
	Dopinfo       string `form:"dopinfo" example:"vendor"`
	UpdateName    string `form:"updateName" json:"updateName" example:"PetrovAA"`
	//Groups        []string `form:"groups" example:"[/Yandex LLC]"`
}
