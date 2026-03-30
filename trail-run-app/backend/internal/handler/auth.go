package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"trail-run-app/internal/model"
	"trail-run-app/pkg/middleware"
	"trail-run-app/pkg/utils"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Token    string `json:"token"`
}

// Register creates a new user account
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	// Check if email already exists
	var existingUser model.User
	if err := model.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.RespondError(c, http.StatusBadRequest, "EMAIL_EXISTS", "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondInternalError(c, "Failed to hash password")
		return
	}

	user := model.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Nickname:     req.Nickname,
	}

	if err := model.DB.Create(&user).Error; err != nil {
		utils.RespondInternalError(c, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		utils.RespondInternalError(c, "Failed to generate token")
		return
	}

	utils.RespondCreated(c, AuthResponse{
		UserID:   user.ID.String(),
		Nickname: user.Nickname,
		Token:    token,
	})
}

// Login authenticates a user
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	var user model.User
	if err := model.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondUnauthorized(c, "Invalid email or password")
			return
		}
		utils.RespondInternalError(c, "Database error")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.RespondUnauthorized(c, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		utils.RespondInternalError(c, "Failed to generate token")
		return
	}

	utils.RespondSuccess(c, AuthResponse{
		UserID:    user.ID.String(),
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Token:     token,
	})
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(c *gin.Context) {
	userIDStr, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.RespondValidationError(c, "Invalid user ID")
		return
	}

	var user model.User
	if err := model.DB.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondNotFound(c, "User not found")
			return
		}
		utils.RespondInternalError(c, "Database error")
		return
	}

	utils.RespondSuccess(c, user.ToResponse())
}
