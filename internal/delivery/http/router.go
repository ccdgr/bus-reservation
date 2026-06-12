package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		NewUserHandler(userGroup)

		busGroup := v1.Group("/buses")
		NewBusHandler(busGroup)

		orderGroup := v1.Group("/orders")
		NewOrderHandler(orderGroup)
	}
}
