package http

import (
	"net/http"
	"strconv"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
)

type BusHandler struct {
	usecase domain.BusUsecase
}

func NewBusHandler(r *gin.RouterGroup, usecase domain.BusUsecase) {
	handler := &BusHandler{usecase: usecase}
	r.GET("", handler.ListBuses)
	r.GET("/:id", handler.GetBus)
}

func (h *BusHandler) ListBuses(c *gin.Context) {
	origin := c.Query("origin")
	dest := c.Query("dest")
	date := c.Query("date")

	buses, err := h.usecase.List(c.Request.Context(), origin, dest, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buses)
}

func (h *BusHandler) GetBus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bus id"})
		return
	}

	bus, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bus not found"})
		return
	}
	c.JSON(http.StatusOK, bus)
}
