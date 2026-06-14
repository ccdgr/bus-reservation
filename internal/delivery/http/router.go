package http

import (
	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, userUsecase domain.UserUsecase, busUsecase domain.BusUsecase, orderUsecase domain.OrderUsecase, jwtSecret, frontendCancelURL string) {
	authMiddleware := AuthMiddleware(jwtSecret)

	// Public routes
	publicV1 := engine.Group("/api/v1")
	{
		RegisterPublicOrderHandler(publicV1.Group("/payments"), orderUsecase, frontendCancelURL)
	}

	// Protected routes
	v1 := engine.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		NewUserHandler(userGroup, userUsecase, authMiddleware)

		busGroup := v1.Group("/buses")
		NewBusHandler(busGroup, busUsecase)

		orderGroup := v1.Group("/orders")
		NewOrderHandler(orderGroup, orderUsecase, authMiddleware)
	}
}
