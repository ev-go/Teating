package request

type UserListReq struct {
	UserId     string     `form:"customedId" json:"customedId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	Pagination Pagination `json:"pagination"`
}
