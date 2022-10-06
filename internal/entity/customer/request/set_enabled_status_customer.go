package request

type SetEnabledStatusCustomer struct {
	CustomerId string `json:"customerId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	Disabled   bool   `json:"disabled" example:"false"`
}
