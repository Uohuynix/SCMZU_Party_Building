package api

import (
	"SCMZU_Party_Building/pkg/jwt"
	"SCMZU_Party_Building/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"name":       user.Name,
			"role":       user.Role,
			"branch":     user.Branch,
			"group_name": user.GroupName,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Role      string `json:"role"`
		Branch    string `json:"branch"`     // 学生/教师的所属党支部
		GroupName string `json:"group_name"` // 学生的所属党小组
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if registerRequest.Role == "" {
		registerRequest.Role = "student"
	}

	// 学生必须绑定具体党支部和党小组；教师只需绑定党支部即可，便于管理本支部下所有小组；
	// 管理员不属于任何具体党支部/小组，忽略前端传入的组织信息。
	if registerRequest.Role == "student" {
		if registerRequest.Branch == "" || registerRequest.GroupName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "branch and group_name are required for student role"})
			return
		}
	} else if registerRequest.Role == "teacher" {
		if registerRequest.Branch == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "branch is required for teacher role"})
			return
		}
	} else if registerRequest.Role == "admin" {
		registerRequest.Branch = ""
		registerRequest.GroupName = ""
	}

	user, err := h.authService.Register(
		registerRequest.Username,
		registerRequest.Password,
		registerRequest.Name,
		registerRequest.Role,
		registerRequest.Branch,
		registerRequest.GroupName,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"name":       user.Name,
			"role":       user.Role,
			"branch":     user.Branch,
			"group_name": user.GroupName,
		},
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"name":       user.Name,
			"role":       user.Role,
			"branch":     user.Branch,
			"group_name": user.GroupName,
		},
	})
}
