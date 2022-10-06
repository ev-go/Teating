package v1

import (
	"github.com/VictoriaMetrics/metrics"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/controller/http/middleware"
	"gitlab.boquar.tech/galileosky/device/customer-administration/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title           GS Customer Administration API
// @version         1.0.0
// @description     customer-administration
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost

// @securityDefinitions.basic  BasicAuth
func NewRouter(handler *gin.Engine, c usecase.ICustomer, g usecase.IGroup, u usecase.IUser) {
	// Options

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.Any("/metrics", func(c *gin.Context) { metrics.WritePrometheus(c.Writer, true) })

	handler.Use(
		//middleware.CORSMiddleware(),
		middleware.ACLMiddleware,
		middleware.TracerMiddleware,
		middleware.MetricsMiddleware,
		middleware.LoggerMiddleware,
		gin.CustomRecovery(middleware.RecoveryMiddleware),
	)

	// Routers
	hCustomer := handler.Group("/customer")
	hGroup := handler.Group("/group")
	hUser := handler.Group("/user")
	hAdminUser := handler.Group("/admin/user")
	{
		newCustomerRoutes(hCustomer, c)
		newGroupRoutes(hGroup, g)
		newUserRoutes(hUser, u)
		newAdminUserRoutes(hAdminUser, u)
	}
}
