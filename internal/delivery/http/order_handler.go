package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase domain.OrderUsecase
}

func NewOrderHandler(r *gin.RouterGroup, usecase domain.OrderUsecase, authMiddleware gin.HandlerFunc) {
	handler := &OrderHandler{usecase: usecase}
	r.Use(authMiddleware)
	r.POST("", handler.CreateOrder)
	r.GET("", handler.ListUserOrders)
	r.POST("/:id/cancel", handler.CancelOrder)
	r.POST("/:id/pay", handler.PayOrder)
	r.POST("/:id/verify", handler.VerifyOrder)
}

func RegisterPublicOrderHandler(r *gin.RouterGroup, usecase domain.OrderUsecase, frontendCancelURL string) {
	handler := &OrderHandler{usecase: usecase}
	r.GET("/paypal/capture", func(c *gin.Context) {
		orderIDStr := c.Query("order_id")
		token := c.Query("token")

		if orderIDStr == "" || token == "" {
			c.Redirect(http.StatusTemporaryRedirect, frontendCancelURL)
			return
		}

		orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
		if err != nil {
			slog.Error("invalid order_id in paypal return", "error", err)
			c.Redirect(http.StatusTemporaryRedirect, frontendCancelURL)
			return
		}

		err = handler.usecase.CapturePayPalPayment(c.Request.Context(), orderID, token)
		if err != nil {
			slog.Error("failed to capture paypal payment", "order_id", orderID, "error", err)
		} else {
			slog.Info("paypal payment captured", "order_id", orderID)
		}

		// Always redirect to frontend orders page
		c.Redirect(http.StatusTemporaryRedirect, frontendCancelURL)
	})
}

type createOrderRequest struct {
	BusID uint64 `json:"bus_id" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint64)
	order, err := h.usecase.Create(c.Request.Context(), userID, req.BusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, order)
}

func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(uint64)
	orders, err := h.usecase.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	err = h.usecase.Cancel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	url, err := h.usecase.Pay(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if url != "" {
		c.JSON(http.StatusOK, gin.H{"payment_url": url})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment successful"})
}

func (h *OrderHandler) VerifyOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	err = h.usecase.Verify(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification successful"})
}
