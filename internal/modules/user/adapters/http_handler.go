package adapters

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/modules/user/ports"
)

type UserHandler struct {
	Svc ports.UseCase
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/login", h.LoginByEmail)
}

func (h *UserHandler) LoginByEmail(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Svc.LoginByEmail(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}