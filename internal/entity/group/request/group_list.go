package request

type GroupListReq struct {
	GroupId    string     `form:"groupId" json:"groupId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	Pagination Pagination `json:"pagination"`
}
