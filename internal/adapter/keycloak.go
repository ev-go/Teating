package adapter

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v11"
	customerRequest "github.com/ev-go/Testing/internal/entity/customer/request"
	groupRequest "github.com/ev-go/Testing/internal/entity/group/request"
	entity "github.com/ev-go/Testing/internal/entity/user"
	"net/http"
	"time"

	"github.com/ev-go/Testing/config"
	userRequest "github.com/ev-go/Testing/internal/entity/user/request"
	"github.com/ev-go/Testing/pkg/tracer"
)

const defaultTimeout = 1

type KeyCloak struct {
	url              string
	clientid         string
	username         string
	clientsecret     string
	accessToken      string
	refreshToken     string
	expiresIn        time.Time
	refreshExpiresIn time.Time
	client           http.Client
}

type GoCloakImpl struct {
	Token        *gocloak.JWT
	Client       gocloak.GoCloak
	ClientId     string
	UserName     string
	ClientSecret string
	Realm        string
}

func NewKeyCloak(ctx context.Context, config config.Config) (k *GoCloakImpl, err error) {
	client := gocloak.NewClient(config.KeyCloak.URL)

	token, err := client.LoginAdmin(ctx, config.KeyCloak.UserName, config.ClientSecret, config.KeyCloak.Realm)
	if err != nil {
		return k, fmt.Errorf("adapter - New - client.LoginAdmin: %w", err)
	}
	return &GoCloakImpl{
		Token:        token,
		Client:       client,
		ClientId:     config.KeyCloak.ClientId,
		UserName:     config.KeyCloak.UserName,
		ClientSecret: config.KeyCloak.ClientSecret,
		Realm:        config.KeyCloak.Realm,
	}, nil
}

func (k *GoCloakImpl) CreateCustomer(ctx context.Context, req customerRequest.CreateCustomerReq) (string, error) {
	ctxNew, spn := tracer.Start(ctx, "adapter - CreateGroup run")
	defer tracer.End(spn)

	group := gocloak.Group{
		Name: &req.Name,
	}

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)

	if err != nil {
		return "", fmt.Errorf("adapter - CreateGroup - client.LoginAdmin: %w", err)
	}
	k.Token = token

	resp, err := k.Client.CreateGroup(ctxNew, k.Token.AccessToken, k.Realm, group)

	if err != nil {
		return "", fmt.Errorf("adapter - CreateGroup - CreateGroup: %w", err)
	}

	return resp, nil
}
func (k *GoCloakImpl) CreateGroup(ctx context.Context, req groupRequest.CreateGroupReq) (string, error) {
	ctxNew, spn := tracer.Start(ctx, "adapter - CreateGroup run")
	defer tracer.End(spn)

	group := gocloak.Group{
		Name: &req.Name,
	}

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)
	if err != nil {
		return "", fmt.Errorf("adapter - CreateGroup - client.LoginAdmin: %w", err)
	}
	k.Token = token

	resp, err := k.Client.CreateChildGroup(ctxNew, k.Token.AccessToken, k.Realm, req.ParentId, group)

	if err != nil {
		return "", fmt.Errorf("adapter - CreateGroup - CreateGroup: %w", err)
	}

	return resp, nil
}
func (k *GoCloakImpl) CreateUser(ctx context.Context, req userRequest.CreateUserReq) (string, error) {
	ctxNew, spn := tracer.Start(ctx, "adapter - CreateUser run")
	defer tracer.End(spn)
	//secretData := "123"
	//credential := gocloak.CredentialRepresentation{SecretData: &secretData}
	user := gocloak.User{
		Username: &req.UserName,
		ID:       &req.UserUuid,
		Enabled:  &req.Enabled,
		//EmailVerified:
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Email:     &req.Email,
		//Attributes: &req.CustomerId
		//DisableableCredentialTypes
		//RequiredActions
		//Access
		//ClientRoles
		//RealmRoles
		//Groups
		//ServiceAccountClientID
		//Credentials: &[]gocloak.CredentialRepresentation{credential},
	}

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)
	if err != nil {
		return "", fmt.Errorf("adapter - CreateUser - client.LoginAdmin: %w", err)
	}
	k.Token = token

	resp, err := k.Client.CreateUser(ctxNew, k.Token.AccessToken, k.Realm, user)

	if err != nil {
		return "", fmt.Errorf("adapter - CreateUser - CreateUser: %w", err)
	}
	password := req.UserName
	err = k.Client.SetPassword(ctxNew, k.Token.AccessToken, resp, k.Realm, password, false)
	if err != nil {
		return "", fmt.Errorf("adapter - CreateUser - CreateUser: %w", err)
	}
	return resp, nil
}

func (k *GoCloakImpl) UpdateUser(ctx context.Context, req userRequest.UpdateUserReq) error {
	ctxNew, spn := tracer.Start(ctx, "adapter - UpdateUser run")
	defer tracer.End(spn)
	user := gocloak.User{
		Username:  &req.UserName,
		ID:        &req.UserUuid,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Email:     &req.Email,
	}

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)
	if err != nil {
		return fmt.Errorf("adapter - CreateGroup - client.LoginAdmin: %w", err)
	}

	k.Token = token
	err = k.Client.UpdateUser(ctxNew, k.Token.AccessToken, k.Realm, user)
	return err
}

func (k *GoCloakImpl) GetUserInfo(ctx context.Context, req userRequest.GetUserInfoReq) (res entity.UserInfo, err error) {
	ctxNew, spn := tracer.Start(ctx, "adapter - UpdateUser run")
	defer tracer.End(spn)

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)
	if err != nil {
		return res, fmt.Errorf("adapter - CreateGroup - client.LoginAdmin: %w", err)
	}

	k.Token = token
	var userInfo *gocloak.User
	userInfo, err = k.Client.GetUserByID(ctxNew, k.Token.AccessToken, k.Realm, req.UserUuid)
	if err != nil {
		return res, fmt.Errorf("adapter - DeleteUser - Client.DeleteUser: %w", err)
	}

	res.UserUuid = *userInfo.ID
	//res.CustomerId  = *userInfo.Attributes
	res.UserName = *userInfo.Username
	res.FirstName = *userInfo.FirstName
	//res.MiddleName = *userInfo.ID
	res.LastName = *userInfo.LastName
	res.Email = *userInfo.Email
	return res, err
}

func (k *GoCloakImpl) DeleteUser(ctx context.Context, req userRequest.DeleteUserReq) error {
	ctxNew, spn := tracer.Start(ctx, "adapter - UpdateUser run")
	defer tracer.End(spn)

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)
	if err != nil {
		return fmt.Errorf("adapter - CreateGroup - client.LoginAdmin: %w", err)
	}

	k.Token = token
	err = k.Client.DeleteUser(ctxNew, k.Token.AccessToken, k.Realm, req.UserUuid)
	if err != nil {
		return fmt.Errorf("adapter - DeleteUser - Client.DeleteUser: %w", err)
	}
	return err
}

func (k *GoCloakImpl) SetEnabledStatusUser(ctx context.Context, req userRequest.SetEnabledStatusUserReq) error {
	ctxNew, spn := tracer.Start(ctx, "adapter - SetEnabledStatusUser run")
	defer tracer.End(spn)
	user := gocloak.User{
		Username: &req.UserName,
		ID:       &req.UserUuid,
		Enabled:  &req.Enabled,
	}

	token, err := k.Client.LoginAdmin(ctx, k.UserName, k.ClientSecret, k.Realm)
	if err != nil {
		return fmt.Errorf("adapter - CreateGroup - client.LoginAdmin: %w", err)
	}

	k.Token = token
	erro := k.Client.UpdateUser(ctxNew, k.Token.AccessToken, k.Realm, user)
	return erro
}
