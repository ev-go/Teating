package request

type DeleteUserReq struct {
	UserUuid string `form:"userUuid" json:"userUuid" example:"not used in request"`
}
