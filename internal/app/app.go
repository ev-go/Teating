// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"github.com/ev-go/Testing/internal/controller/grpcserver"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/ev-go/Testing/config"
	v1 "github.com/ev-go/Testing/internal/controller/http/v1"

	"github.com/ev-go/Testing/internal/adapter"
	"github.com/ev-go/Testing/internal/usecase"
	"github.com/ev-go/Testing/pkg/httpserver"
	"github.com/ev-go/Testing/pkg/logger"
	"github.com/ev-go/Testing/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	pg, err := postgres.New(cfg.PG.DatabaseURL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		logger.I.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	//err = acl.NewACL(cfg.Acl)
	//if err != nil {
	//	logger.I.Fatal(fmt.Errorf("app - Run - newAcl: %w", err))
	//}

	keycloak, err := adapter.NewKeyCloak(context.Background(), *cfg)
	if err != nil {
		logger.I.Fatal(fmt.Errorf("app - Run - newKeyCloak: %w", err))
	}

	gin.SetMode(gin.DebugMode)
	handler := gin.New()

	customerRepo := adapter.NewCustomerRepo(pg)
	groupRepo := adapter.NewGroupRepo(pg)
	userRepo := adapter.NewUserRepo(pg)

	customerUseCase := usecase.NewCustomer(customerRepo, keycloak)
	groupUseCase := usecase.NewGroup(groupRepo, keycloak)
	userUseCase := usecase.NewUser(userRepo, keycloak)

	v1.NewRouter(handler, customerUseCase, groupUseCase, userUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	go func() {
		err = grpcserver.New(cfg.GRPC.Port, customerRepo, userRepo)
		if err != nil {
			logger.I.Fatal(fmt.Errorf("app - Run - grpcserver.New: %w", err))
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.I.Error(fmt.Errorf("app - Run - signal: " + s.String()))
	case err = <-httpServer.Notify():
		logger.I.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.I.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
