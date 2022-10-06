package adapter

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	entity "github.com/ev-go/Testing/internal/entity/user"
	userRequest "github.com/ev-go/Testing/internal/entity/user/request"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/ev-go/Testing/pkg/postgres"
	"github.com/ev-go/Testing/pkg/tracer"
	//"gitlab.boquar.tech/galileosky/pkg/http_errors"
	//"net/http"
)

const _resultNotUpdate = "UPDATE 0"

type UserRepoImpl struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepoImpl {
	return &UserRepoImpl{pg}
}

func (r *UserRepoImpl) UserList(ctx context.Context, req userRequest.UserListReq) (res entity.UserList, err error) {
	newCtx, spn := tracer.Start(ctx, "UserList run")
	defer tracer.End(spn)

	res = entity.UserList{}

	q := r.Builder.
		Select("customers_uuid",
			"username",
			"inn",
			"firstname",
			"middlename",
			"lastname",
			"country",
			"region",
			"contacts",
			"dop_info",
			"enabled",
		).
		From("customers.customer")

	rows, err := r.PaginationQuery(newCtx, q, req.Pagination)
	if err != nil {
		return res, fmt.Errorf("adapter - CustomerList - r.PaginationQuery: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := entity.UserListItem{}
		err = rows.Scan(
			&item.UserId,
			&item.Name,
			&item.Inn,
			&item.FullName,
			&item.Country,
			&item.Region,
			&item.Contacts,
			&item.DopInfo,
			&item.Enabled,
		)
		if err != nil {
			return res, fmt.Errorf("adapter - CustomerList - rows.Scan: %w", err)
		}
		res.Items = append(res.Items, item)
	}

	//TODO pagination, total count
	return res, nil
}

func (r *UserRepoImpl) UserListTotal(ctx context.Context) (int, error) {
	newCtx, span := tracer.Start(ctx, "repo - CustomerListTotal")
	defer tracer.End(span)

	sql, args, err := r.Builder.
		Select("count(*)").
		From("customers.customer").
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

func (r *UserRepoImpl) PaginationQuery(ctx context.Context, builder squirrel.SelectBuilder, page userRequest.Pagination) (pgx.Rows, error) {

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

func (r *UserRepoImpl) GetUserInfo(ctx context.Context, req userRequest.GetUserInfoReq, requestPermission string) (res entity.UserInfo, err error) {
	_, spn := tracer.Start(ctx, "GetUserInfo run")
	defer tracer.End(spn)

	res = entity.UserInfo{}
	//	res.LastNetworkAddress = "192.168.0.1"

	q := r.Builder.
		Select(
			"u.customers_uuid",
			//	"id",
			"u.user_uuid",
			"u.username",
			"u.email",
			"u.enabled",
			"u.firstname",
			"u.middlename",
			"u.lastname",
			"u.phone",
			"u.create_time",
			"u.create_user",
			"u.update_time",
			"u.update_user",
			"u.dop_info",
			"u.password",
			"u.last_connection_time",
			//"u.last_network_address",
		).
		From("(" + requestPermission + ") as u ").
		Where(squirrel.Eq{"u.username": req.UserName})

	sql, args, err := q.ToSql()

	if err != nil {
		return res, fmt.Errorf("adapter - GetUserInfoAll - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)

	if err != nil {
		return res, fmt.Errorf("adapter - GetUserInfoAll - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			//	&res.Id,
			&res.CustomerId,
			&res.UserUuid,
			&res.UserName,
			&res.Email,
			&res.Enabled,
			&res.FirstName,
			&res.MiddleName,
			&res.LastName,
			&res.Phone,
			&res.CreateTime,
			&res.CreateUser,
			&res.UpdateTime,
			&res.UpdateUser,
			&res.DopInfo,
			&res.PasswordHash,
			&res.LastСonnectionTime,
			//&res.LastNetworkAddress,
		)
		if err != nil {
			return res, fmt.Errorf("adapter - GetUserInfoAll - rows.Scan: %w", err)
		}
	}

	return res, nil
}

// Deprecated: 123
func (r *UserRepoImpl) UpdateUser(ctx context.Context, req userRequest.UpdateUserReq) error {
	_, spn := tracer.Start(ctx, "UpdateUser run")
	defer tracer.End(spn)

	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		//err = tx.Rollback(ctx)
		//if err != nil {
		//	return fmt.Errorf("adapter - UpdateUser 1 - r.Pool.Exec: %w", err)
		//}
		return fmt.Errorf("adapter - ReadApiKey - tx.Commit: %w", err)
	}

	q := r.Builder.
		Update("customers.users")

	if req.UserName != "" {
		q = q.Set("username", req.UserName)
	}
	if req.FirstName != "" {
		q = q.Set("firstname", req.FirstName)
	}
	if req.MiddleName != "" {
		q = q.Set("middlename", req.MiddleName)
	}
	if req.LastName != "" {
		q = q.Set("lastname", req.LastName)
	}
	if req.Phone != "" {
		q = q.Set("phone", req.Phone)
	}
	if req.Email != "" {
		q = q.Set("email", req.Email)
	}
	if req.Dopinfo != "" {
		q = q.Set("dop_info", req.Dopinfo)
	}

	q = q.Set("update_user", "user")

	q = q.Set("update_time", "now()")

	q = q.Where(squirrel.Eq{"username": req.UserName})
	//q = q.Where("customers_uuid = ?", req.UserId)

	sql, args, err := q.ToSql()
	sql = sql + " RETURNING users_pk, username, firstname, middlename, lastname, phone, email, update_time, update_user, dop_info" // fixme

	if err != nil {
		return fmt.Errorf("adapter - UpdateUser 1 - r.Builder: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)
	type updateUserRes struct {
		UsersPk    int //fixme
		UserName   string
		FirstName  string
		MiddleName string
		LastName   string
		Phone      string
		Email      string
		UpdateTime time.Time
		UpdateUser string
		Dopinfo    string
	}

	var res updateUserRes
	err = row.Scan(
		&res.UsersPk,
		&res.UserName,
		&res.FirstName,
		&res.MiddleName,
		&res.LastName,
		&res.Phone,
		&res.Email,
		&res.UpdateTime,
		&res.UpdateUser,
		&res.Dopinfo,
	)
	//fixme err

	//if err != nil {
	//	return res, fmt.Errorf("adapter - GetUserInfoAll - rows.Scan: %w", err)
	//}

	if err != nil {
		return fmt.Errorf("adapter - UpdateUser 1 - r.Pool.Exec: %w", err)
	}

	q = r.Builder.
		Update("acl.users_customer_roles")

	if req.CustomerRoles != 0 {
		q = q.Set("customer_roles_pk", req.CustomerRoles)
	}
	//q = q.Where(squirrel.Eq{"user_uuid": req.UserId})
	q = q.Where("users_pk = ?", res.UsersPk)

	sql, args, err = q.ToSql()

	if err != nil {
		return fmt.Errorf("adapter - UpdateUser 2 - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("adapter - UpdateUser 1 - r.Pool.Exec: %w", err)
		}
		return fmt.Errorf("adapter - ReadApiKey - tx.Commit: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("adapter - UpdateUser 1 - r.Pool.Exec: %w", err)
		}
		return fmt.Errorf("adapter - ReadApiKey - tx.Commit: %w", err)
	}

	return nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, req userRequest.CreateUserReq) error {
	_, spn := tracer.Start(ctx, "CreateUser run")
	defer tracer.End(spn)
	//tx, err := r.Pool.Begin(ctx)

	//txOptions := pgx.TxOptions{"ReadCommitted", "ReadWrite", "Deferrable"}

	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("adapter - CreateUser - r.Pool.BeginTx: %w", err)
	}

	sql, args, err := r.Builder.
		Insert("customers.users").
		Columns(
			"customers_uuid",
			"user_uuid",
			"username",
			"email",
			"enabled",
			"firstname",
			"middlename",
			"lastname",
			"phone",
			"create_user",
			"update_user",
			"dop_info",
			"password",
			//	"last_network_address",
		).
		Values(
			req.CustomersUuid,
			req.UserUuid,
			req.UserName,
			req.Email,
			req.Enabled,
			req.FirstName,
			req.MiddleName,
			req.LastName,
			req.Phone,
			req.CreateUser,
			req.UpdateUser,
			req.Dopinfo,
			req.Password,
			//	"192.168.0.1",
		).ToSql()

	if err != nil {
		return fmt.Errorf("adapter - CreateUser - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("adapter - CreateUser - tx.Exec: %w", err)
	}

	//....................
	var usersPk int
	sql, args, err = r.Builder.
		Select(
			"users_pk",
		).
		From("customers.users").
		Where(squirrel.Eq{"username": req.UserName}).ToSql()

	rows, err := tx.Query(ctx, sql, args...)

	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("adapter - CreateUser - tx.Rollback: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&usersPk,
		)
	}
	//....................

	sql, args, err = r.Builder.
		Insert("acl.users_customer_roles").
		Columns(
			"customer_roles_pk",
			"users_pk",
		).
		Values(
			req.CustomerRoles,
			usersPk,
		).ToSql()

	if err != nil {
		return fmt.Errorf("adapter - CreateUser - r.Builder: %w", err)
	}

	//_, err = r.Pool.Exec(ctx, sql, args...)
	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("adapter - ReadApiKey - r.Pool.Exec: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("adapter - ReadApiKey - tx.Commit: %w", err)
	}
	return nil
}

// SetEnabledStatusUser выполняет запрос вида
/*
update customers.users
set enabled = {req.Enabled}
where user_uuid = {req.UserId}
and user_uuid in (
	select u.user_uuid from (
		{requestPermission}
	) as u
)
*/
func (r *UserRepoImpl) SetEnabledStatusUser(ctx context.Context, req userRequest.SetEnabledStatusUserReq, requestPermission string) (*string, error) {
	_, spn := tracer.Start(ctx, "SetEnabledStatusUser run")
	defer tracer.End(spn)
	/*	sqlStringPermission, err := acl.ACLRepo.GetPermissions(ctx, reqPermission)

		if err != nil {
			return fmt.Errorf("adapter - SetEnabledStatusUser - acl.ACLRepo.GetPermissions: %w", err)
		}

		q := r.Builder.
			Select(
				"g.id",
				"g.name",
				"g.parent_id",
				"g.created_at",
				"g.updated_at",
			).Suffix("FROM (" + sqlStringPermission + ") AS g")*/

	q := r.Builder.
		Update("customers.users").
		Set("enabled", req.Enabled).
		Where(squirrel.Eq{"user_uuid": req.UserUuid}).
		Where("user_uuid in (select u.user_uuid from (" + requestPermission + ") as u)")

	sql, args, err := q.ToSql()
	sql = sql + " RETURNING username"
	if err != nil {
		return nil, err
	}

	var username string

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&username)
	if err != nil {
		return nil, err
	}

	if username == "" { //fixme изменить на err?
		return nil, err
	}

	return &username, nil
}

func (userRepo *UserRepoImpl) GetCustomerUUIDByUserName(ctx context.Context, userName string) (customerUUID string, err error) {
	_, spn := tracer.Start(ctx, "GetCustomerUUIDByUserName run")
	defer tracer.End(spn)

	sql, args, err := userRepo.Builder.
		Select("u.customers_uuid").
		From("customers.users u").
		Where("u.username = ?", userName).
		ToSql()
	if err != nil {
		return customerUUID, fmt.Errorf("adapter - GetCustomerUUIDByUserName - userRepo.Builder: %s", err)
	}

	if err = userRepo.Pool.QueryRow(ctx, sql, args...).Scan(&customerUUID); err != nil {
		return customerUUID, fmt.Errorf("adapter - GetCustomerUUIDByUserName - userRepo.Pool.QueryRow: %s", err)
	}

	return customerUUID, nil
}
