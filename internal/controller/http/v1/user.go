package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	userRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user/request"
	uResp "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user/user/response"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/usecase"
	"gitlab.boquar.tech/galileosky/pkg/acl"
	"net/http"
)

type userRoutes struct {
	uc usecase.IUser
}

func newUserRoutes(handler *gin.RouterGroup, uc usecase.IUser) {
	r := &userRoutes{uc}

	handler.GET("/list", r.GetUserList)
	//handler.GET("/info", r.GetUserInfoPublic)
	handler.GET("/info", r.GetUserInfo)
	handler.POST("", r.CreateUser)
	handler.PATCH("/info", r.UpdateUser) // TODO change for patch
	handler.PATCH("/enabled", r.SetEnabledStatusUser)

}

// @Summary UserList
// @Param parentId query int false "1"
// @Schemes
// @Description GET UserList
// @Tags User
// @ID user-list
// @Accept json
// @Produce json
// @Success      200              {object}	response.UserListRes    "UserListRes"
// @failure      400	          {string}  string    "error"
// @Router /user/list [GET]
func (r *userRoutes) GetUserList(c *gin.Context) {
	req := userRequest.UserListReq{
		Pagination: userRequest.Pagination{
			Order: "name asc",
		},
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)

		return
	}

	res, err := r.uc.UserList(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary GetUserInfo
// @Param   userId     query     string     false  "string default"     default(28b255fc-755b-4141-810e-283e08ebe836)
// @Schemes
// @Description GET GetUserInfo
// @Tags User
// @ID read-user-info
// @Accept json
// @Produce json
// @Success      200              {object}	response.UserInfoRes    "UserInfoRes"
// @failure      400	          {string}  string    "error"
// @Router /user/info [GET]
func (r *userRoutes) GetUserInfo(c *gin.Context) {
	var req userRequest.GetUserInfoReq

	ctx := c.Request.Context()
	tokenParams := ctx.Value(acl.CtxNameTokenParams).(acl.TokenParams)
	req.CustomersUuid = tokenParams.CustomerUUID

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	userInfoAllRes, err := r.uc.GetUserInfo(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userInfoAllJSON, _ := json.Marshal(userInfoAllRes) //TODO verify
	userInfoPubRes := uResp.UserInfoRes{}
	err = json.Unmarshal(userInfoAllJSON, &userInfoPubRes)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userInfoPubRes)
}

// @Summary CreateUser
// @Param request body request.CreateUserReq true "query params"
// @Schemes
// @Description GET CreateUser
// @Tags User
// @ID create-user
// @Accept json
// @Produce json
// @Success      200
// @failure      400	          {string}  string    "error"
// @Router /user [POST]
func (r *userRoutes) CreateUser(c *gin.Context) {
	var req userRequest.CreateUserReq

	ctx := c.Request.Context()
	tokenParams := ctx.Value(acl.CtxNameTokenParams).(acl.TokenParams)
	req.CustomersUuid = tokenParams.CustomerUUID
	req.CreateUser = tokenParams.Username
	req.UpdateUser = tokenParams.Username
	req.Enabled = true

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.CreateUser(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, err)
}

// @Summary UpdateUser
// @Param request body request.UpdateUserReq true "query params"
// @Schemes
// @Description GET UpdateUser
// @Tags User
// @ID update-user-info
// @Accept json
// @Produce json
// @Success      200
// @failure      400	          {string}  string    "error"
// @Router /user/info [PATCH]
func (r *userRoutes) UpdateUser(c *gin.Context) {
	var req userRequest.UpdateUserReq

	ctx := c.Request.Context()
	tokenParams := ctx.Value(acl.CtxNameTokenParams).(acl.TokenParams)
	req.CustomersUuid = tokenParams.CustomerUUID
	req.UpdateName = tokenParams.Username

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.UpdateUser(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")
}

// @Summary SetEnabledStatusUser
// @Param userId query string false "1"
// @Schemes
// @Description GET SetEnabledStatusUser
// @Tags User
// @ID enabled-user
// @Accept json
// @Produce json
// @Success      200
// @failure      400	          {string}  string    "error"
// @Router /user/enabled [PATCH]
func (r *userRoutes) SetEnabledStatusUser(c *gin.Context) {
	var req userRequest.SetEnabledStatusUserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.SetEnabledStatusUser(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")
}
