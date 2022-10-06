package v1

import (
	"fmt"
	userRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user/request"
	"gitlab.boquar.tech/galileosky/pkg/acl"
	"net/http"

	"github.com/gin-gonic/gin"
	uResp "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user/admin/response"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/usecase"
)

type adminUserRoutes struct {
	uc usecase.IUser
}

func newAdminUserRoutes(handler *gin.RouterGroup, uc usecase.IUser) {
	r := &adminUserRoutes{uc}

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
// @Tags AdminUser
// @ID admin-user-list
// @Accept json
// @Produce json
// @Success      200              {object}	response.UserListRes    "UserListRes"
// @failure      400	          {string}  string    "error"
// @Router /admin/user/list [GET]
func (r *adminUserRoutes) GetUserList(c *gin.Context) {
	isAdmin := c.Request.Context().Value(acl.CtxNameTokenParams).(acl.TokenParams).IsSupport
	if !isAdmin {
		respondWithError(c, http.StatusForbidden, fmt.Errorf("ТЫ НЕ ПРОЙДЕШЬ!!!")) // todo accept to all endpoints
		return
	}
	req := userRequest.UserListReq{
		Pagination: userRequest.Pagination{
			Order: "name asc",
		},
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)

		return
	}

	//token := c.request.Header.Get("Authorization")
	//tokenParams, err := acl.ParseToken(token)

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
// @Tags AdminUser
// @ID admin-read-user-info
// @Accept json
// @Produce json
// @Success      200              {object}	response.AdminUserInfoRes    "AdminUserInfoRes"
// @failure      400	          {string}  string    "error"
// @Router /admin/user/info [GET]
func (r *adminUserRoutes) GetUserInfo(c *gin.Context) {
	var req userRequest.GetUserInfoReq

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	res, err := r.uc.GetUserInfo(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := uResp.AdminUserInfoRes(res)
	c.JSON(http.StatusOK, resp)
}

// @Summary CreateUser
// @Param request body request.CreateUserReq true "query params"
// @Schemes
// @Description GET CreateUser
// @Tags AdminUser
// @ID admin-create-user
// @Accept json
// @Produce json
// @Success      200
// @failure      400	          {string}  string    "error"
// @Router /admin/user [POST]
func (r *adminUserRoutes) CreateUser(c *gin.Context) {
	var req userRequest.CreateUserReq

	ctx := c.Request.Context()
	tokenParams := ctx.Value(acl.CtxNameTokenParams).(acl.TokenParams)
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
// @Tags AdminUser
// @ID admin-update-user-info
// @Accept json
// @Produce json
// @Success      200
// @failure      400	          {string}  string    "error"
// @Router /admin/user/info [PATCH]
func (r *adminUserRoutes) UpdateUser(c *gin.Context) {
	var req userRequest.UpdateUserReq

	ctx := c.Request.Context()
	tokenParams := ctx.Value(acl.CtxNameTokenParams).(acl.TokenParams)

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
// @Param request body request.SetEnabledStatusUserReq true "query params"
// @Schemes
// @Description GET SetEnabledStatusUser
// @Tags AdminUser
// @ID admin-enabled-user
// @Accept json
// @Produce json
// @Success      200
// @failure      400	          {string}  string    "error"
// @Router /admin/user/enabled [PATCH]
func (r *adminUserRoutes) SetEnabledStatusUser(c *gin.Context) {
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
