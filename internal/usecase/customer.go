package usecase

import (
	"context"
	"fmt"
	customerRequest "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/customer/request"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/customer/response"

	"gitlab.boquar.tech/galileosky/device/customer-administration/pkg/tracer"
	"gitlab.boquar.tech/galileosky/pkg/acl"
)

type CustomerUCImpl struct {
	repo ICustomerRepo
	a    *acl.ACL
	k    IKeycloak
}

func NewCustomer(r ICustomerRepo, k IKeycloak) *CustomerUCImpl {
	return &CustomerUCImpl{
		repo: r,
		k:    k,
	}
}

func (uc *CustomerUCImpl) CustomerList(ctx context.Context, req customerRequest.CustomerListReq) (res response.CustomerList, err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - CustomerList run")
	defer tracer.End(spn)

	res, err = uc.repo.CustomerList(ctxNew, req)
	if err != nil {
		return res, fmt.Errorf("uc - CustomerList - repo.CustomerList: %w", err)
	}

	count, err := uc.repo.CustomerListTotal(ctxNew, req)
	if err != nil {
		return res, fmt.Errorf("uc - CustomerList - repo.CustomerListTotal: %w", err)
	}
	res.Total = count
	return res, nil
}

func (uc *CustomerUCImpl) GetCustomer(ctx context.Context, customerId string) (res response.CustomerRes, err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - GetCustomer run")
	defer tracer.End(spn)

	res, err = uc.repo.GetCustomer(ctxNew, customerId)
	if err != nil {
		return res, fmt.Errorf("uc - GetCustomer - repo.GetCustomer: %w", err)
	}

	return res, nil
}

func (uc *CustomerUCImpl) CreateCustomer(ctx context.Context, req customerRequest.CreateCustomerReq) (err error) {
	ctx, spn := tracer.Start(ctx, "uc - CreateCustomer run")
	defer tracer.End(spn)

	//req.CustomerId = uuid.NewString()

	customerId, err := uc.k.CreateCustomer(ctx, req)

	if err != nil {
		return fmt.Errorf("uc - CreateCustomer - k.CreateGroup: %w", err)
	}

	req.CustomerId = customerId
	err = uc.repo.CreateCustomer(ctx, req)
	if err != nil {
		return fmt.Errorf("uc - CreateCustomer - repo.CreateCustomer: %w", err)
	}

	return nil
}

func (uc *CustomerUCImpl) UpdateCustomer(ctx context.Context, req customerRequest.UpdateCustomerReq) error {
	ctxNew, spn := tracer.Start(ctx, "uc - UpdateCustomer run")
	defer tracer.End(spn)

	err := uc.repo.UpdateCustomer(ctxNew, req)
	if err != nil {
		return fmt.Errorf("uc - UpdateCustomer - repo.UpdateCustomer: %w", err)
	}

	return nil
}

func (uc *CustomerUCImpl) SetDisabledStatusCustomer(ctx context.Context, req customerRequest.SetEnabledStatusCustomer) error {
	ctxNew, spn := tracer.Start(ctx, "uc - SetDisabledStatusCustomer run")
	defer tracer.End(spn)

	err := uc.repo.SetDisabledStatusCustomer(ctxNew, req)
	if err != nil {
		return fmt.Errorf("uc - SetDisabledStatusCustomer - repo.SetDisabledStatusCustomer: %w", err)
	}

	return nil
}
