package v1

import (
	customerRequest "github.com/ev-go/Testing/internal/entity/customer/request"
	"github.com/ev-go/Testing/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type customerRoutes struct {
	uc usecase.ICustomer
}

func newCustomerRoutes(handler *gin.RouterGroup, uc usecase.ICustomer) {
	r := &customerRoutes{uc}

	handler.GET("/list", r.GetCustomerList)
	handler.GET("", r.GetCustomer)
	handler.POST("", r.CreateCustomer)
	handler.PUT("", r.UpdateCustomer)
	handler.PATCH("", r.SetDisabledStatusCustomer)

}

// @Summary CustomerList
// @Param offset query int false "0"
// @Param limit query int false "100"
// @Param order query string false "id desc"
// @Schemes
// @Description GET CustomerList
// @Tags customer
// @ID list-customer
// @Accept json
// @Produce json
// @Success      200              {object}	response.CustomerList    "CustomerList"
// @failure      400	          {string}  string    "error"
// @Router /customer/list [GET]
func (r *customerRoutes) GetCustomerList(c *gin.Context) {
	req := customerRequest.CustomerListReq{
		Pagination: customerRequest.Pagination{
			Order: "name asc",
		},
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)

		return
	}

	res, err := r.uc.CustomerList(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary GetCustomer
// @Param customerId query int false "0da3b22f-ec3f-4383-bc25-480b6dcb82a1"
// @Schemes
// @Description GET GetCustomer
// @Tags customer
// @ID read-customer
// @Accept json
// @Produce json
// @Success      200              {object}	response.CustomerRes    "CustomerRes"
// @failure      400	          {string}  string    "error"
// @Router /customer [GET]
func (r *customerRoutes) GetCustomer(c *gin.Context) {
	var req customerRequest.CustomerReq

	if err := c.ShouldBindQuery(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)

		return
	}

	res, err := r.uc.GetCustomer(c.Request.Context(), req.CustomerId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary CreateCustomer
// @Param body body request.CreateCustomerReq true "JSON"
// @Schemes
// @Description GET CreateCustomer
// @Tags customer
// @ID create-customer
// @Accept json
// @Produce json
// @Success      200              {string}  string    "ok"
// @failure      400	          {string}  string    "error"
// @Router /customer [POST]
func (r *customerRoutes) CreateCustomer(c *gin.Context) {
	var req customerRequest.CreateCustomerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.CreateCustomer(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, err)
}

// @Summary UpdateCustomer
// @Param body body request.UpdateCustomerReq true "JSON"
// @Schemes
// @Description GET UpdateCustomer
// @Tags customer
// @ID update-customer
// @Accept json
// @Produce json
// @Success      200              {string}  string    "ok"
// @failure      400	          {string}  string    "error"
// @Router /customer [PUT]
func (r *customerRoutes) UpdateCustomer(c *gin.Context) {
	var req customerRequest.UpdateCustomerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.UpdateCustomer(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")
}

// @Summary SetDisabledStatusCustomer
// @Param body body request.UpdateCustomerReq true "JSON"
// @Schemes
// @Description GET SetDisabledStatusCustomer
// @Tags customer
// @ID toggle-customer
// @Accept json
// @Produce json
// @Success      200              {string}  string    "ok"
// @failure      400	          {string}  string    "error"
// @Router /customer [PATCH]
func (r *customerRoutes) SetDisabledStatusCustomer(c *gin.Context) {
	var req customerRequest.UpdateCustomerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	err := r.uc.UpdateCustomer(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "")
}
