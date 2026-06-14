package http

import (
	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

func NewRouter(engine *gin.Engine, userUsecase domain.UserUsecase, busUsecase domain.BusUsecase, orderUsecase domain.OrderUsecase, aliClient *alipay.Client, jwtSecret, frontendCancelURL string) {
	authMiddleware := AuthMiddleware(jwtSecret)

	// Public routes
	publicV1 := engine.Group("/api/v1")
	{
		RegisterPublicOrderHandler(publicV1.Group("/payments"), orderUsecase, aliClient, frontendCancelURL)
	}

	// Protected routes
	v1 := engine.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		NewUserHandler(userGroup, userUsecase, authMiddleware)

		busGroup := v1.Group("/buses")
		NewBusHandler(busGroup, busUsecase)

		orderGroup := v1.Group("/orders")
		NewOrderHandler(orderGroup, orderUsecase, aliClient, authMiddleware)
	}
}
