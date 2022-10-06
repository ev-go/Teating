package adapter

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/ev-go/Testing/internal/entity/customer/response"
	groupRequest "github.com/ev-go/Testing/internal/entity/group/request"
	"github.com/ev-go/Testing/pkg/postgres"
	"github.com/ev-go/Testing/pkg/tracer"
	"github.com/jackc/pgx/v4"
)

type GroupRepoImpl struct {
	*postgres.Postgres
}

func NewGroupRepo(pg *postgres.Postgres) *GroupRepoImpl {
	return &GroupRepoImpl{pg}
}

func (r *GroupRepoImpl) GroupList(ctx context.Context, req groupRequest.GroupListReq) (res response.GroupList, err error) {
	newCtx, spn := tracer.Start(ctx, "CustomerList run")
	defer tracer.End(spn)

	res = response.GroupList{}

	q := r.Builder.
		Select("customer_id",
			"name",
			"inn",
			"fullname",
			"country",
			"region",
			"contacts",
			"dopinfo",
			"disabled",
		).
		From("customer.customer")

	rows, err := r.PaginationQuery(newCtx, q, req.Pagination)
	if err != nil {
		return res, fmt.Errorf("adapter - CustomerList - r.PaginationQuery: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := response.GroupListItem{}
		err = rows.Scan(
			&item.GroupId,
			&item.Name,
			&item.Inn,
			&item.FullName,
			&item.Country,
			&item.Region,
			&item.Contacts,
			&item.DopInfo,
			&item.Disabled,
		)
		if err != nil {
			return res, fmt.Errorf("adapter - CustomerList - rows.Scan: %w", err)
		}
		res.Items = append(res.Items, item)
	}

	//TODO pagination, total count
	return res, nil
}

func (r *GroupRepoImpl) GroupListTotal(ctx context.Context) (int, error) {
	newCtx, span := tracer.Start(ctx, "repo - CustomerListTotal")
	defer tracer.End(span)

	sql, args, err := r.Builder.
		Select("count(*)").
		From("customer_administration.customer").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("CustomerRepoImpl - CustomerListTotal - r.Builder: %w", err)
	}
	var count int
	err = r.Pool.QueryRow(newCtx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CustomerRepoImpl - CustomerListTotal - r.Pool.QueryRow: %w", err)
	}
	return count, nil
}

func (r *GroupRepoImpl) PaginationQuery(ctx context.Context, builder squirrel.SelectBuilder, page groupRequest.Pagination) (pgx.Rows, error) {

	newCtx, span := tracer.Start(ctx, "repo - PaginationQuery")
	defer tracer.End(span)

	sql, args, err := builder.
		OrderBy(page.Order).
		Limit(page.Limit).
		Offset(page.Offset).
		ToSql()
	if err != nil {
		fmt.Println(fmt.Errorf("CustomerepoImpl - PaginationQuery - builder: %w", err))
	}

	return r.Pool.Query(newCtx, sql, args...)
}

func (r *GroupRepoImpl) GetGroup(ctx context.Context, groupId string) (res response.GroupRes, err error) {
	_, spn := tracer.Start(ctx, "GetCustomer run")
	defer tracer.End(spn)

	res = response.GroupRes{}

	q := r.Builder.
		Select("customer_id",
			"name",
			"inn",
			"fullname",
			"country",
			"region",
			"contacts",
			"dopinfo",
			"disabled",
		).
		From("customer.customer").
		Where("customer_id = ?", groupId)

	sql, args, err := q.ToSql()

	if err != nil {
		return res, fmt.Errorf("adapter - GetCustomer - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)

	if err != nil {
		return res, fmt.Errorf("adapter - GetCustomer - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&res.GroupId,
			&res.Name,
			&res.Inn,
			&res.FullName,
			&res.Country,
			&res.Region,
			&res.Contacts,
			&res.DopInfo,
			&res.Disabled,
		)
		if err != nil {
			return res, fmt.Errorf("adapter - GetCustomer - rows.Scan: %w", err)
		}
	}

	return res, nil
}

func (r *GroupRepoImpl) UpdateGroup(ctx context.Context, req groupRequest.UpdateGroupReq) error {
	_, spn := tracer.Start(ctx, "UpdateCustomer run")
	defer tracer.End(spn)

	q := r.Builder.
		Update("customer.customer")

	if req.Name != nil {
		q = q.Set("name", req.Name)
	}
	if req.Inn != nil {
		q = q.Set("inn", req.Inn)
	}
	if req.FullName != nil {
		q = q.Set("fullname", req.FullName)
	}
	if req.Country != nil {
		q = q.Set("country", req.Country)
	}
	if req.Region != nil {
		q = q.Set("region", req.Region)
	}
	if req.Contacts != nil {
		q = q.Set("contacts", req.Contacts)
	}
	if req.DopInfo != nil {
		q = q.Set("dopinfo", req.DopInfo)
	}
	if req.Disabled != nil {
		q = q.Set("disabled", req.Disabled)
	}
	q = q.Where("customer_id = ?", req.GroupId)

	sql, args, err := q.ToSql()

	if err != nil {
		return fmt.Errorf("adapter - UpdateCustomer - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("adapter - UpdateCustomer - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *GroupRepoImpl) CreateGroup(ctx context.Context, req groupRequest.CreateGroupReq) error {
	_, spn := tracer.Start(ctx, "ReadApiKey run")
	defer tracer.End(spn)

	q := r.Builder.
		Insert("customer.customer").
		Columns("customer_id",
			"name",
			"inn",
			"fullname",
			"country",
			"region",
			"contacts",
			"dopinfo",
			"disabled",
			"apikey",
		).
		Values(req.GroupId,
			req.Name,
			req.Inn,
			req.FullName,
			req.Country,
			req.Region,
			req.Contacts,
			req.DopInfo,
			req.Disabled,
			req.ApiKey,
		)

	sql, args, err := q.ToSql()

	if err != nil {
		return fmt.Errorf("adapter - ReadApiKey - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("adapter - ReadApiKey - r.Pool.Exec: %w", err)
	}

	return nil
}

/*
func (r *CustomerRepoImpl) SetDisabledStatusCustomer(ctx context.Context, req request.DisabledStatusReq) error {
	_, spn := tracer.Start(ctx, "SetDisabledStatusCustomer run")
	defer tracer.End(spn)

	q := r.Builder.
		Update("customer.customer").
		Set("disabled", req.Disabled)

	sql, args, err := q.ToSql()

	if err != nil {
		return fmt.Errorf("adapter - SetDisabledStatusCustomer - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("adapter - SetDisabledStatusCustomer - r.Pool.Exec: %w", err)
	}

	return nil
}
*/
