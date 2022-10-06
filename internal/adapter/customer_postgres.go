package adapter

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	customerRequest "github.com/ev-go/Testing/internal/entity/customer/request"
	"github.com/ev-go/Testing/internal/entity/customer/response"
	"github.com/jackc/pgx/v4"

	"github.com/ev-go/Testing/pkg/postgres"
	"github.com/ev-go/Testing/pkg/tracer"
)

type CustomerRepoImpl struct {
	*postgres.Postgres
}

func NewCustomerRepo(pg *postgres.Postgres) *CustomerRepoImpl {
	return &CustomerRepoImpl{pg}
}

func (r *CustomerRepoImpl) CustomerList(ctx context.Context, req customerRequest.CustomerListReq) (res response.CustomerList, err error) {
	newCtx, spn := tracer.Start(ctx, "CustomerList run")
	defer tracer.End(spn)

	res = response.CustomerList{}

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
		item := response.CustomerListItem{}
		err = rows.Scan(
			&item.CustomerId,
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

func (r *CustomerRepoImpl) CustomerListTotal(ctx context.Context, req customerRequest.CustomerListReq) (int, error) {
	newCtx, span := tracer.Start(ctx, "repo - CustomerListTotal")
	defer tracer.End(span)

	sql, args, err := r.Builder.
		Select("count(*)").
		From("customer.customer").
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

func (r *CustomerRepoImpl) PaginationQuery(ctx context.Context, builder squirrel.SelectBuilder, page customerRequest.Pagination) (pgx.Rows, error) {

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

func (r *CustomerRepoImpl) GetCustomer(ctx context.Context, customerId string) (res response.CustomerRes, err error) {
	_, spn := tracer.Start(ctx, "GetCustomer run")
	defer tracer.End(spn)

	res = response.CustomerRes{}

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
		Where("customer_id = ?", customerId)

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
			&res.CustomerId,
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

func (r *CustomerRepoImpl) UpdateCustomer(ctx context.Context, req customerRequest.UpdateCustomerReq) error {
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
	q = q.Where("customer_id = ?", req.CustomerId)

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

func (r *CustomerRepoImpl) ReadApiKey(ctx context.Context, customerUUID string) (res string, err error) {
	_, spn := tracer.Start(ctx, "ReadApiKey run")
	defer tracer.End(spn)

	res = ""

	q := r.Builder.
		Select("apikey").
		From("customers.customers").
		Where("customers_uuid = ?", customerUUID)

	sql, args, err := q.ToSql()

	if err != nil {
		return res, fmt.Errorf("adapter - ReadApiKey - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)

	if err != nil {
		return res, fmt.Errorf("adapter - ReadApiKey - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&res,
		)
		if err != nil {
			return res, fmt.Errorf("adapter - ReadApiKey - rows.Scan: %w", err)
		}
	}

	return res, nil
}

func (r *CustomerRepoImpl) SetApiKey(ctx context.Context, req customerRequest.CustomerSetApiKeyReq) (err error) {
	_, spn := tracer.Start(ctx, "SetApiKey run")
	defer tracer.End(spn)

	q := r.Builder.
		Update("customers.customers").
		Set("apikey", req.ApiKey).
		Where("customers_uuid = ?", req.CustomerUUID)

	sql, args, err := q.ToSql()

	if err != nil {
		return fmt.Errorf("adapter - SetApiKey - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("adapter - SetApiKey - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *CustomerRepoImpl) CreateCustomer(ctx context.Context, req customerRequest.CreateCustomerReq) error {
	_, spn := tracer.Start(ctx, "CreateCustomer run")
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
		Values(req.CustomerId,
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

func (r *CustomerRepoImpl) SetDisabledStatusCustomer(ctx context.Context, req customerRequest.SetEnabledStatusCustomer) error {
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
