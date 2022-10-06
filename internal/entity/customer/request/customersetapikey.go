package request

type CustomerSetApiKeyReq struct {
	CustomerUUID string `json:"customerId" binding:"required"`
	ApiKey       string `json:"apiKey" binding:"required"`
}
