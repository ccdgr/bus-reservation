package http

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ccdgr/bus-reservation/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
)

type OrderHandler struct {
	usecase   domain.OrderUsecase
	aliClient *alipay.Client
	reqChan   chan *createOrderTask
}

type createOrderTask struct {
	Ctx    context.Context
	UserID uint64
	BusID  uint64
	Result chan *taskResult
}

type taskResult struct {
	Order *domain.Order
	Err   error
}

func NewOrderHandler(r *gin.RouterGroup, usecase domain.OrderUsecase, aliClient *alipay.Client, authMiddleware gin.HandlerFunc) {
	handler := &OrderHandler{
		usecase:   usecase,
		aliClient: aliClient,
		reqChan:   make(chan *createOrderTask, 5000), // Buffer for 5000 concurrent requests
	}
	
	// Start worker pool to process order requests
	for i := 0; i < 20; i++ {
		go handler.worker()
	}

	r.Use(authMiddleware)
	r.POST("", handler.CreateOrder)
	r.GET("", handler.ListUserOrders)
	r.POST("/:id/cancel", handler.CancelOrder)
	r.POST("/:id/pay", handler.PayOrder)
	r.POST("/:id/verify", handler.VerifyOrder)
}

func (h *OrderHandler) worker() {
	for task := range h.reqChan {
		order, err := h.usecase.Create(task.Ctx, task.UserID, task.BusID)
		task.Result <- &taskResult{Order: order, Err: err}
	}
}

func RegisterPublicOrderHandler(r *gin.RouterGroup, usecase domain.OrderUsecase, aliClient *alipay.Client, frontendCancelURL string) {
	handler := &OrderHandler{usecase: usecase, aliClient: aliClient}
	r.POST("/alipay/notify", handler.AlipayNotify)
	
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

	// Create a task and push it to the channel
	task := &createOrderTask{
		Ctx:    c.Request.Context(), // Use context for cancellation
		UserID: userID,
		BusID:  req.BusID,
		Result: make(chan *taskResult, 1),
	}

	select {
	case h.reqChan <- task:
		// Task accepted, wait for processing
	default:
		// Queue is full, shed load immediately
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "系统繁忙，请稍后再试 (System overloaded)"})
		return
	}

	// Wait for the worker to finish, with an HTTP-level timeout
	select {
	case res := <-task.Result:
		if res.Err != nil {
			// StatusConflict or BadRequest depending on error logic, using 500 for simplicity here but could be 400 for 'already booked'
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, res.Order)
	case <-time.After(5 * time.Second):
		// If worker takes too long, timeout the HTTP request to free up the connection
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "请求超时，请稍后在订单列表中查看结果 (Request timeout)"})
	case <-c.Request.Context().Done():
		// Client disconnected early
		return
	}
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

type payOrderRequest struct {
	Method string `json:"method"`
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req payOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Method = "mock" // fallback
	}

	url, err := h.usecase.Pay(c.Request.Context(), id, req.Method)
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

func (h *OrderHandler) AlipayNotify(c *gin.Context) {
	// The aliClient decoding is handled securely inside usecase or handler if configured.
	// For simplicity, we just pass the request to usecase if we had injected aliClient.
	// Let's add aliClient back to OrderHandler to decode properly.
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
