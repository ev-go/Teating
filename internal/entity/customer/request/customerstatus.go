package request

type CustomerStatusReq struct {
	CustomerId int  `json:"customerId" binding:"required"`
	Disabled   bool `json:"disabled" binding:"required"`
}
