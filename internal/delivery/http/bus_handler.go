package http

import (
	"github.com/gin-gonic/gin"
)

type BusHandler struct {
	// TODO: Inject BusUsecase
}

func NewBusHandler(r *gin.RouterGroup) {
	handler := &BusHandler{}
	r.GET("", handler.ListBuses)
	r.GET("/:id", handler.GetBus)
}

func (h *BusHandler) ListBuses(c *gin.Context) {
	// TODO: Implement list buses logic
}

func (h *BusHandler) GetBus(c *gin.Context) {
	// TODO: Implement get bus detail logic
}
