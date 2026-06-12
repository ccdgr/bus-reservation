package http

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	// TODO: Inject UserUsecase
}

func NewUserHandler(r *gin.RouterGroup) {
	handler := &UserHandler{}
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.GET("/profile", handler.GetProfile)
}

func (h *UserHandler) Register(c *gin.Context) {
	// TODO: Implement registration logic
}

func (h *UserHandler) Login(c *gin.Context) {
	// TODO: Implement login logic
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: Implement profile retrieval logic
}
