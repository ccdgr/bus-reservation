package http

import (
	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine, userUsecase domain.UserUsecase, busUsecase domain.BusUsecase, orderUsecase domain.OrderUsecase, jwtSecret string) {
	authMiddleware := AuthMiddleware(jwtSecret)

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
