package users

import (
	"buymeagiftapi/internal/auth"
	"buymeagiftapi/internal/config/variables"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *UsersController) Delete(ctx *gin.Context) {
	db := c.GetDatabase(ctx)
	userId := c.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	service := NewService(db)
	err := service.DeleteUser(*userId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	repository := auth.NewRepository(db)
	jwtWrapper := auth.NewJwtWrapper(variables.GetJwtVariables(), repository)

	token := ctx.GetHeader("Authorization")
	extractedToken, _ := jwtWrapper.ExtractToken(token)

	err = repository.BlacklistToken(extractedToken)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
