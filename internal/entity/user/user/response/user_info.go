package response

import "time"

type UserInfoRes struct {
	CustomerId string `json:"customerId" example:"0da3b22f-ec3f-4383-bc25-480b6dcb82a1"`
	//Id         int    `json:"id" example:"1"`
	Uuid       string    `json:"uuid" example:"c5074793-9d82-478c-9853-125c04bdb626"`
	UserName   string    `json:"username" example:"user1"`
	Email      string    `json:"email" example:"user1@gmail.com"`
	FirstName  string    `json:"firstName" example:"Sidr"`
	MiddleName string    `json:"middleName" example:"Sidorovich"`
	LastName   string    `json:"lastName" example:"Sidorov"`
	Phone      string    `json:"phone" example:"+79091234567"`
	RoleId     string    `json:"roleId" example:"???"` //fixme roleId example? // not BD column
	CreateTime time.Time `json:"createTime" example:"7/14/22 9:04:36 AM"`
	CreateUser string    `json:"createUser" example:"user1"`
	UpdateTime time.Time `json:"updateTime" example:"7/14/22 9:04:36 AM"`
	UpdateUser string    `json:"updateUser" example:"user1"`
}
