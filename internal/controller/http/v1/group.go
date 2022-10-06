package v1

import (
	"github.com/gin-gonic/gin"

	groupRequest "github.com/ev-go/Testing/internal/entity/group/request"

	"github.com/ev-go/Testing/internal/usecase"
	"net/http"
)

type groupRoutes struct {
	uc usecase.IGroup
}

func newGroupRoutes(handler *gin.RouterGroup, uc usecase.IGroup) {
	r := &groupRoutes{uc}

	handler.GET("/list", r.GetGroupList)
	handler.GET("", r.GetGroup)
	handler.POST("", r.CreateGroup)
	handler.PUT("", r.UpdateGroup)
	//handler.PATCH("", r.SetDisabledStatusCustomer)

}

// @Summary GroupList
// @Param parentId query int false "1"
// @Schemes
// @Description GET GroupList
// @Tags group
// @ID usergroup-list
// @Accept json
// @Produce json
// @Success      200              {object}	response.CustomerRes    "CustomerRes"
// @failure      400	          {string}  string    "error"
// @Router /group/list [GET]
func (r *groupRoutes) GetGroupList(c *gin.Context) {
	req := groupRequest.GroupListReq{
		Pagination: groupRequest.Pagination{
			Order: "name asc",
		},
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)

		return
	}

	res, err := r.uc.GroupList(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary GetGroup
// @Param parentId query int false "1"
// @Schemes
// @Description GET GetGroup
// @Tags group
// @ID group-read
// @Accept json
// @Produce json
// @Success      200              {object}	response.CustomerRes    "CustomerRes"
// @failure      400	          {string}  string    "error"
// @Router /group [GET]
func (r *groupRoutes) GetGroup(c *gin.Context) {
	var req groupRequest.GetGroupReq

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)

		return
	}

	res, err := r.uc.GetGroup(c.Request.Context(), req.CustomerId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func (r *groupRoutes) CreateGroup(c *gin.Context) {
	var req groupRequest.CreateGroupReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.CreateGroup(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, err)
}

func (r *groupRoutes) UpdateGroup(c *gin.Context) {
	var req groupRequest.UpdateGroupReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.UpdateGroup(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")
}

/*
func (r *customerRoutes) SetDisabledStatusCustomer(c *gin.Context) {
	var req request.UpdateCustomerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.UpdateCustomer(c.request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")
}
*/
