package http

import (
	"net/http"
	"strconv"
	"log/slog"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

type OrderHandler struct {
	usecase   domain.OrderUsecase
	aliClient *alipay.Client
}

func NewOrderHandler(r *gin.RouterGroup, usecase domain.OrderUsecase, aliClient *alipay.Client, authMiddleware gin.HandlerFunc) {
	handler := &OrderHandler{usecase: usecase, aliClient: aliClient}
	r.Use(authMiddleware)
	r.POST("", handler.CreateOrder)
	r.GET("", handler.ListUserOrders)
	r.POST("/:id/cancel", handler.CancelOrder)
	r.POST("/:id/pay", handler.PayOrder)
	r.POST("/:id/verify", handler.VerifyOrder)
}

func RegisterPublicOrderHandler(r *gin.RouterGroup, usecase domain.OrderUsecase, aliClient *alipay.Client) {
	handler := &OrderHandler{usecase: usecase, aliClient: aliClient}
	r.POST("/alipay/notify", handler.AlipayNotify)
}

func (h *OrderHandler) AlipayNotify(c *gin.Context) {
	if h.aliClient == nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	req := c.Request
	req.ParseForm()

	noti, err := h.aliClient.DecodeNotification(c.Request.Context(), req.Form)
	if err != nil {
		slog.Error("failed to decode alipay notification", "error", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if noti.TradeStatus == "TRADE_SUCCESS" {
		orderID, err := strconv.ParseUint(noti.OutTradeNo, 10, 64)
		if err == nil {
			err = h.usecase.ConfirmPayment(c.Request.Context(), orderID)
			if err != nil {
				slog.Error("failed to confirm payment", "order_id", orderID, "error", err)
			} else {
				slog.Info("alipay payment confirmed", "order_id", orderID)
			}
		}
	}

	h.aliClient.ACKNotification(c.Writer)
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
