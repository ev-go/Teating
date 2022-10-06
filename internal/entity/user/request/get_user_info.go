package request

type GetUserInfoReq struct {
	UserUuid      string `form:"userUuid" json:"userUuid" example:"not used in request"`
	CustomersUuid string `form:"customersUuid" json:"customersUuid" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	UserName      string `form:"userName" json:"userName" example:"ivanovAA"`
}
