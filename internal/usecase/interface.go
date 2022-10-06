package usecase

import (
	"context"
	groupRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/group/request"
	entity "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user"
	userRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user/request"
	"gitlab.boquar.tech/galileosky/pkg/http_errors"

	customerRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/customer/request"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/customer/response"
)

type (

	//customers

	ICustomer interface {
		CustomerList(ctx context.Context, req customerRequest.CustomerListReq) (res response.CustomerList, err error)
		GetCustomer(ctx context.Context, customerId string) (res response.CustomerRes, err error)
		CreateCustomer(ctx context.Context, req customerRequest.CreateCustomerReq) (err error)
		UpdateCustomer(ctx context.Context, req customerRequest.UpdateCustomerReq) (err error)
		SetDisabledStatusCustomer(ctx context.Context, req customerRequest.SetEnabledStatusCustomer) error
	}
	ICustomerRepo interface {
		CustomerList(ctx context.Context, req customerRequest.CustomerListReq) (res response.CustomerList, err error)
		CustomerListTotal(ctx context.Context, req customerRequest.CustomerListReq) (int, error)
		GetCustomer(ctx context.Context, customerId string) (res response.CustomerRes, err error)
		CreateCustomer(ctx context.Context, req customerRequest.CreateCustomerReq) error
		UpdateCustomer(ctx context.Context, req customerRequest.UpdateCustomerReq) (err error)
		ReadApiKey(ctx context.Context, customerId string) (res string, err error)
		SetApiKey(ctx context.Context, req customerRequest.CustomerSetApiKeyReq) (err error)
		SetDisabledStatusCustomer(ctx context.Context, req customerRequest.SetEnabledStatusCustomer) error
	}

	//groups

	IGroup interface {
		GroupList(ctx context.Context, req groupRequest.GroupListReq) (res response.GroupList, err error)
		GetGroup(ctx context.Context, groupId string) (res response.GroupRes, err error)
		CreateGroup(ctx context.Context, req groupRequest.CreateGroupReq) (err error)
		UpdateGroup(ctx context.Context, req groupRequest.UpdateGroupReq) (err error)
		//SetDisabledStatusGroup(ctx context.Context, req request.DisabledStatusReq) error
	}

	IGroupRepo interface {
		GroupList(ctx context.Context, req groupRequest.GroupListReq) (res response.GroupList, err error)
		GroupListTotal(ctx context.Context) (int, error)
		GetGroup(ctx context.Context, groupId string) (res response.GroupRes, err error)
		CreateGroup(ctx context.Context, req groupRequest.CreateGroupReq) error
		UpdateGroup(ctx context.Context, req groupRequest.UpdateGroupReq) (err error)
		//SetDisabledStatusCustomer(ctx context.Context, req request.DisabledStatusReq) error
	}

	//users

	IUser interface {
		UserList(ctx context.Context, req userRequest.UserListReq) (res entity.UserList, err error)
		GetUserInfo(ctx context.Context, req userRequest.GetUserInfoReq) (res entity.UserInfo, err error)
		//GetUserInfoPublic(ctx context.Context, userId string) (res entity.UserInfo, err error)
		CreateUser(ctx context.Context, req userRequest.CreateUserReq) (err error)
		UpdateUser(ctx context.Context, req userRequest.UpdateUserReq) (err error)
		DeleteUser(ctx context.Context, req userRequest.DeleteUserReq) error
		SetEnabledStatusUser(ctx context.Context, req userRequest.SetEnabledStatusUserReq) error
	}

	IUserRepo interface {
		UserList(ctx context.Context, req userRequest.UserListReq) (res entity.UserList, err error)
		UserListTotal(ctx context.Context) (int, error)
		GetUserInfo(ctx context.Context, req userRequest.GetUserInfoReq, requestPermission string) (res entity.UserInfo, err error)
		//GetUserInfoPublic(ctx context.Context, userId string, requestPermission string) (res entity.UserInfo, err error)
		CreateUser(ctx context.Context, req userRequest.CreateUserReq) error
		UpdateUser(ctx context.Context, req userRequest.UpdateUserReq) (err error)
		SetEnabledStatusUser(ctx context.Context, req userRequest.SetEnabledStatusUserReq, requestPermission string) (*string, *http_errors.RestError)
		GetCustomerUUIDByUserName(ctx context.Context, userName string) (string, error)
	}

	IKeycloak interface {
		CreateCustomer(ctx context.Context, req customerRequest.CreateCustomerReq) (string, error)
		CreateGroup(ctx context.Context, req groupRequest.CreateGroupReq) (string, error)
		CreateUser(ctx context.Context, req userRequest.CreateUserReq) (string, error)
		GetUserInfo(ctx context.Context, req userRequest.GetUserInfoReq) (res entity.UserInfo, err error)
		UpdateUser(ctx context.Context, req userRequest.UpdateUserReq) (err error)
		DeleteUser(ctx context.Context, req userRequest.DeleteUserReq) error
		SetEnabledStatusUser(ctx context.Context, req userRequest.SetEnabledStatusUserReq) error
	}
)
