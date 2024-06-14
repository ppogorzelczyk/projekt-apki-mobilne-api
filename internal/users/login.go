package users

import (
	"buymeagiftapi/internal/auth"
	"buymeagiftapi/internal/config/variables"
	"buymeagiftapi/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	domain.User
}

func (u *UsersController) Login(ctx *gin.Context) {
	db := u.GetDatabase(ctx)

	if db == nil {
		return
	}

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	service := NewService(db)
	user, err := service.GetUserByEmail(req.Email)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if user == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if !user.CheckPassword(req.Password) {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	token, err := auth.NewJwtWrapper(variables.GetJwtVariables(), nil).GenerateToken(*user)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}
