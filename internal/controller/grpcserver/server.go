package grpcserver

import (
	"context"
	"fmt"
	"github.com/ev-go/Testing/internal/entity/customer/request"
	"net"

	"github.com/ev-go/Testing/internal/usecase"
	"github.com/ev-go/Testing/pkg/tracer"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	customerRepo usecase.ICustomerRepo
	userRepo     usecase.IUserRepo
}

func (grpc *GRPCServer) GetApiKey(ctx context.Context, req *GetApiKeyRequest) (res *GetApiKeyResponse, err error) {
	ctxNew, span := tracer.Start(ctx, "GetAPiKey run")
	defer tracer.End(span)

	resRepo, err := grpc.customerRepo.ReadApiKey(ctxNew, req.CustomerId)

	if err != nil {
		return nil, fmt.Errorf("grpcserver - GetApiKey - customerRepo.ReadApiKey: %w", err)
	}

	res = &GetApiKeyResponse{
		Apikey: resRepo,
	}

	return res, nil
}

func (grpc *GRPCServer) SetApiKey(ctx context.Context, req *SetApiKeyRequest) (res *SetApiKeyResponse, err error) {
	ctxNew, span := tracer.Start(ctx, "SetApiKey run")
	defer tracer.End(span)
	reqRepo := request.CustomerSetApiKeyReq{
		CustomerUUID: req.CustomerId,
		ApiKey:       req.Apikey,
	}

	err = grpc.customerRepo.SetApiKey(ctxNew, reqRepo)

	if err != nil {
		return res, fmt.Errorf("grpcserver - SetApiKey - customerRepo.SetApiKey: %w", err)
	}

	res = &SetApiKeyResponse{
		Apikey: reqRepo.ApiKey,
	}

	return res, nil
}

func (grpc *GRPCServer) GetCustomerUUIDByUserName(ctx context.Context, req *GetCustomerUUIDByUserNameRequest) (res *GetCustomerUUIDByUserNameResponse, err error) {
	ctxNew, span := tracer.Start(ctx, "GetCustomerUUIDByUserName run")
	defer tracer.End(span)

	customerUUID, err := grpc.userRepo.GetCustomerUUIDByUserName(ctxNew, req.UserName)
	if err != nil {
		return res, fmt.Errorf("grpcserver - GetCustomerUUIDByUserName - userRepo.GetCustomerUUIDByUserName: %w", err)
	}

	res = &GetCustomerUUIDByUserNameResponse{
		CustomerUuid: customerUUID,
	}
	return res, nil
}

func (grpc *GRPCServer) mustEmbedUnimplementedCustomerServer() {

}

func New(port string, customerRepo usecase.ICustomerRepo, userRepo usecase.IUserRepo) error {

	s := grpc.NewServer()
	srv := &GRPCServer{
		customerRepo: customerRepo,
		userRepo:     userRepo,
	}
	RegisterCustomerServer(s, srv)
	l, err := net.Listen("tcp", port)

	if err != nil {
		return fmt.Errorf("grpcserver - New - net.Listen: %w", err)
	}

	if err := s.Serve(l); err != nil {
		return fmt.Errorf("grpcserver - New - s.Serve: %w", err)
	}
	return nil
}
