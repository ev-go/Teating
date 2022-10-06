package usecase

import (
	"context"
	"fmt"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/customer/response"
	request2 "gitlab.boquar.tech/galileosky/device/customer-administration/internal/entity/group/request"
	"gitlab.boquar.tech/galileosky/device/customer-administration/pkg/tracer"
	"gitlab.boquar.tech/galileosky/pkg/acl"
)

type GroupUCImpl struct {
	repo IGroupRepo
	a    *acl.ACL
	k    IKeycloak
}

func NewGroup(r IGroupRepo, k IKeycloak) *GroupUCImpl {
	return &GroupUCImpl{
		repo: r,
		k:    k,
	}
}

func (uc *GroupUCImpl) GroupList(ctx context.Context, req request2.GroupListReq) (res response.GroupList, err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - GroupList run")
	defer tracer.End(spn)

	res, err = uc.repo.GroupList(ctxNew, req)
	if err != nil {
		return res, fmt.Errorf("uc - GroupList - repo.GroupList: %w", err)
	}

	count, err := uc.repo.GroupListTotal(ctxNew)
	if err != nil {
		return res, fmt.Errorf("uc - GroupList - repo.GroupListTotal: %w", err)
	}
	res.Total = count
	return res, nil
}

func (uc *GroupUCImpl) GetGroup(ctx context.Context, groupId string) (res response.GroupRes, err error) {
	ctxNew, spn := tracer.Start(ctx, "uc - GetGroup run")
	defer tracer.End(spn)

	res, err = uc.repo.GetGroup(ctxNew, groupId)
	if err != nil {
		return res, fmt.Errorf("uc - GetGroup - repo.GetGroup: %w", err)
	}

	return res, nil
}

func (uc *GroupUCImpl) CreateGroup(ctx context.Context, req request2.CreateGroupReq) (err error) {
	ctx, spn := tracer.Start(ctx, "uc - CreateGroup run")
	defer tracer.End(spn)

	//req.CustomerId = uuid.NewString()

	groupId, err := uc.k.CreateGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("uc - CreateGroup - k.CreateGroup: %w", err)
	}

	req.GroupId = groupId
	//err = uc.repo.CreateGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("uc - CreateGroup - repo.CreateGroup: %w", err)
	}

	return nil
}

func (uc *GroupUCImpl) UpdateGroup(ctx context.Context, req request2.UpdateGroupReq) error {
	ctxNew, spn := tracer.Start(ctx, "uc - UpdateGroup run")
	defer tracer.End(spn)

	err := uc.repo.UpdateGroup(ctxNew, req)
	if err != nil {
		return fmt.Errorf("uc - UpdateGroup - repo.UpdateGroup: %w", err)
	}

	return nil
}

/*
func (uc *CustomerUCImpl) SetDisabledStatusCustomer(ctx context.Context, req request.DisabledStatusReq) error {
	ctxNew, spn := tracer.Start(ctx, "uc - SetDisabledStatusCustomer run")
	defer tracer.End(spn)

	err := uc.repo.SetDisabledStatusCustomer(ctxNew, req)
	if err != nil {
		return fmt.Errorf("uc - SetDisabledStatusCustomer - repo.SetDisabledStatusCustomer: %w", err)
	}

	return nil
}
*/
