// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/egor-denisov/wallet-infotecs/docs"
	"github.com/egor-denisov/wallet-infotecs/internal/usecase"
	"github.com/egor-denisov/wallet-infotecs/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       EWallet
// @version     1.0
// @host        localhost:8000
// @BasePath    /api/v1
func NewRouter(handler *gin.Engine, l logger.Interface, w usecase.Wallet) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// For generation new Swagger documentation: 
	// swag init -dir internal/controller/http/v1/ -generalInfo router.go --parseDependency internal/entity/ 
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers
	h := handler.Group("/api/v1")
	{
		newWalletRoutes(h, w, l)
	}
}
