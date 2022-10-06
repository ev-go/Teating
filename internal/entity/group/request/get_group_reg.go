package request

type GetGroupReq struct {
	GroupId    string `json:"groupId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	CustomerId string `json:"customerId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	ParentId   string `json:"parentId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	Name       string `json:"name" example:"OOO Galileosky"`
	Inn        string `json:"inn" example:"5904254657"`
	FullName   string `json:"fullName" example:"Limited liability company Galileosky"`
	Country    string `json:"country" example:"Russia"`
	Region     string `json:"region" example:"Perm region"`
	Contacts   string `json:"contacts" example:"8 495 001 3930"`
	DopInfo    string `json:"dopInfo" example:"vendor"`
	Disabled   bool   `json:"disabled" example:"false"`
	ApiKey     string `json:"apiKey" example:"cSPZV2BtniuCyGynsww7PY.LsMsj3TEnYMbezinoA6NsL"`
}
