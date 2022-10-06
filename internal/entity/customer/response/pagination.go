package response

type Pagination struct {
	Offset uint64 `json:"offset" form:"offset" default:"0"`
	Limit  uint64 `json:"limit" form:"limit"`
	Order  string `json:"order" form:"order" default:"id desc"`
}
