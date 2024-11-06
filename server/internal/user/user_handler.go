package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var userReq CreateUserRequest

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.CreateUser(c.Request.Context(), &userReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": user})
}

func (h *Handler) Login(c *gin.Context) {
	var userReq LoginUserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := h.Service.Login(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error----jjkjkjk": err.Error()})
	}

	c.SetCookie("jwt", user.AccessToken, 3600, "/", "localhost", false, false)

	res := &LoginUserResponse{
		AccessToken: user.AccessToken,
		ID:          user.ID,
		Username:    user.Username,
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": res})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful", "data": "ok"})
}
