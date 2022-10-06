package request

type SetEnabledStatusUserReq struct {
	UserUuid      string `json:"userUuid" example:"28b255fc-755b-4141-810e-283e08ebe836"`
	CustomersUuid string `form:"customersUuid" json:"customersUuid" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	UserName      string `json:"userName" example:"user2"`
	Enabled       bool   `json:"enabled" example:"false"`
}
