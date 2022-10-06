package usecase

import (
	"context"
	"fmt"
	entity "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user"
	userRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/user/request"
	"gitlab.boquar.tech/galileosky/device/customer-administration/pkg/tracer"
	"gitlab.boquar.tech/galileosky/pkg/acl"
)

type UserUCImpl struct {
	repo IUserRepo
	a    *acl.ACL
	k    IKeycloak
}

func NewUser(r IUserRepo, k IKeycloak) *UserUCImpl {
	return &UserUCImpl{
		repo: r,
		k:    k,
	}
}

func (uc *UserUCImpl) UserList(ctx context.Context, req userRequest.UserListReq) (res entity.UserList, err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - CustomerList run")
	defer tracer.End(spn)

	//tokenParams := ctx.Value("TokenParams").(acl.TokenBody)
	//requestParams := ctx.Value("RequestParams").(acl.Operation)

	res, err = uc.repo.UserList(ctxNew, req)
	if err != nil {
		return res, fmt.Errorf("uc - CustomerList - repo.CustomerList: %w", err)
	}

	count, err := uc.repo.UserListTotal(ctxNew)
	if err != nil {
		return res, fmt.Errorf("uc - CustomerList - repo.CustomerListTotal: %w", err)
	}
	res.Total = count
	return res, nil
}

func (uc *UserUCImpl) GetUserInfo(ctx context.Context, req userRequest.GetUserInfoReq) (res entity.UserInfo, err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - GetUserInfoAll run")
	defer tracer.End(spn)

	tokenParams := ctx.Value("TokenParams").(acl.TokenParams)
	requestParams := ctx.Value("RequestParams").(acl.RequestParams) // объект acl.RequestParams
	//GetQueryWithAclOnUsers(ctx context.Context, subjectUUID string, operationId string, customerUUID string, isAdmin bool) (string, error) {
	sqlPermission, err := acl.ACLRepo.GetQueryWithAclOnUsers(
		ctxNew, tokenParams.UserUUID, requestParams.OperationId, tokenParams.CustomerUUID, tokenParams.IsSupport,
	)
	if err != nil {
		return res, fmt.Errorf("uc - GetUserInfoAll - acl.ACLRepo.HasPermission: %w", err)
	}

	res, err = uc.repo.GetUserInfo(ctxNew, req, sqlPermission)
	if err != nil {
		return res, fmt.Errorf("uc - GetUserInfoAll - repo.GetUserInfoAll: %w", err)
	}

	return res, nil
}

//func (uc *UserUCImpl) GetUserInfoPublic(ctx context.Context, userId string) (res entity.UserInfo, err error) {
//	ctxNew, spn := tracer.Start(ctx, "uc - GetUserInfoPublic run")
//	defer tracer.End(spn)
//
//	sqlPermission, err := acl.ACLRepo.HasPermission(ctxNew, "")
//	if err != nil {
//		return res, fmt.Errorf("uc - GetUserInfoAll - acl.ACLRepo.HasPermission: %w", err)
//	}
//
//	res, err = uc.repo.GetUserInfoPublic(ctxNew, userId, sqlPermission)
//	if err != nil {
//		return res, fmt.Errorf("uc - GetUserInfoPublic - repo.GetUserInfoPublic: %w", err)
//	}
//
//	return res, nil
//}

// Deprecated: 123
func (uc *UserUCImpl) CreateUser(ctx context.Context, req userRequest.CreateUserReq) (err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - CreateCustomer run")
	defer tracer.End(spn)

	userUuid, err := uc.k.CreateUser(ctx, req)
	if err != nil {
		return fmt.Errorf("uc - CreateUSer - k.CreateUser: %w", err)
	}

	req.UserUuid = userUuid
	err = uc.repo.CreateUser(ctxNew, req)
	if err != nil {
		DeleteUserReq := userRequest.DeleteUserReq{UserUuid: req.UserUuid}
		err := uc.k.DeleteUser(ctxNew, DeleteUserReq)
		if err != nil {
			return fmt.Errorf("uc - CreateUser - k.CreateUser: %w", err)
		}
		return fmt.Errorf("uc - CreateUser - repo.CreateUser: %w", err)
	}

	return nil
}

func (uc *UserUCImpl) UpdateUser(ctx context.Context, req userRequest.UpdateUserReq) error {
	ctxNew, spn := tracer.Start(ctx, "uc - UpdateUser run")
	defer tracer.End(spn)

	var (
		kCloakUserInfoBeforeTx entity.UserInfo
		err                    error
	)
	getUserInfoReq := userRequest.GetUserInfoReq{
		UserUuid:      req.UserUuid,
		UserName:      req.UserName,
		CustomersUuid: req.CustomersUuid,
	}
	kCloakUserInfoBeforeTx, err = uc.k.GetUserInfo(ctx, getUserInfoReq)
	RollbackUpdateUserReq := userRequest.UpdateUserReq{
		UserUuid:  kCloakUserInfoBeforeTx.UserUuid,
		UserName:  kCloakUserInfoBeforeTx.UserName,
		FirstName: kCloakUserInfoBeforeTx.FirstName,
		LastName:  kCloakUserInfoBeforeTx.LastName,
		Email:     kCloakUserInfoBeforeTx.Email,
	}

	if err != nil {

		return fmt.Errorf("uc - UpdateUser - k.UpdateUser: %w", err)
	}

	err = uc.k.UpdateUser(ctx, req)
	if err != nil {

		return fmt.Errorf("uc - UpdateUser - k.UpdateUser: %w", err)
	}

	err = uc.repo.UpdateUser(ctxNew, req)
	if err != nil {
		err = uc.k.UpdateUser(ctx, RollbackUpdateUserReq)
		if err != nil {
			return fmt.Errorf("uc - UpdateUser - k.UpdateUser: %w", err)
		}
		return fmt.Errorf("uc - UpdateUser - repo.UpdateUser: %w", err)
	}

	return nil
}

func (uc *UserUCImpl) DeleteUser(ctx context.Context, req userRequest.DeleteUserReq) error {
	ctxNew, spn := tracer.Start(ctx, "uc - UpdateUser run")
	defer tracer.End(spn)

	err := uc.k.DeleteUser(ctxNew, req)
	if err != nil {
		return fmt.Errorf("uc - DeleteUser - k.DeleteUser: %w", err)
	}

	return nil
}

func (uc *UserUCImpl) SetEnabledStatusUser(ctx context.Context, req userRequest.SetEnabledStatusUserReq) error {
	ctxNew, spn := tracer.Start(ctx, "uc - SetEnabledStatusUser run")
	defer tracer.End(spn)

	tokenParams := ctx.Value("TokenParams").(acl.TokenParams)
	requestParams := ctx.Value("RequestParams").(acl.RequestParams) // объект acl.RequestParams
	//GetQueryWithAclOnUsers(ctx context.Context, subjectUUID string, operationId string, customerUUID string, isAdmin bool) (string, error) {
	sqlPermission, err := acl.ACLRepo.GetQueryWithAclOnUsers(
		ctxNew, tokenParams.UserUUID, requestParams.OperationId, tokenParams.CustomerUUID, tokenParams.IsSupport,
	)

	if err != nil {
		return fmt.Errorf("uc - GetUserInfoAll - acl.ACLRepo.HasPermission: %w", err)
	}

	username, errHTTP := uc.repo.SetEnabledStatusUser(ctxNew, req, sqlPermission)
	//fixme разобраться почему здесь не воспринимает ошибку как nil если она по факту nil
	if errHTTP != nil {
		return fmt.Errorf("uc - SetEnabledStatusUser - repo.SetEnabledStatusUser: %w", errHTTP)
	}

	req.UserName = *username
	erro := uc.k.SetEnabledStatusUser(ctx, req)
	if err != nil {
		return fmt.Errorf("uc - CreateCustomer - k.CreateGroup: %w", err)
	}

	return erro
}
