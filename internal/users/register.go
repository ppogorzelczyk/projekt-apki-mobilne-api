package users

import (
	"buymeagiftapi/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required"`
	Username  string  `json:"username" binding:"required"`
	Phone     *string `json:"phone"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

func (c *UsersController) Register(ctx *gin.Context) {
	db := c.GetDatabase(ctx)

	if db == nil {
		return
	}

	var req NewUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user := domain.User{
		Email:     req.Email,
		Username:  req.Username,
		Phone:     req.Phone,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	user.HashPassword(req.Password)

	service := NewService(db)
	createdUser, err := service.CreateUser(user)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if createdUser == nil {
		ctx.Status(http.StatusConflict)
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}
