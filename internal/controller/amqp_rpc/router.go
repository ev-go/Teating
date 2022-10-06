package amqprpc

import (
	"github.com/ev-go/Testing/internal/usecase"
	"github.com/ev-go/Testing/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}
