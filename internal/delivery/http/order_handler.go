package http

import (
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	// TODO: Inject OrderUsecase
}

func NewOrderHandler(r *gin.RouterGroup) {
	handler := &OrderHandler{}
	r.POST("", handler.CreateOrder)
	r.GET("", handler.ListUserOrders)
	r.GET("/:id", handler.GetOrder)
	r.POST("/:id/cancel", handler.CancelOrder)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// TODO: Implement order creation logic
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	// TODO: Implement list user orders logic
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	// TODO: Implement get order detail logic
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	// TODO: Implement order cancellation logic
}
