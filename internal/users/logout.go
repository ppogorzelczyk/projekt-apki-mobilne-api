package users

import (
	"buymeagiftapi/internal/auth"
	"buymeagiftapi/internal/config/variables"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *UsersController) Logout(ctx *gin.Context) {
	db := c.GetDatabase(ctx)
	userId := c.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	token := ctx.GetHeader("Authorization")

	if token == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	repository := auth.NewRepository(db)
	jwtWrapper := auth.NewJwtWrapper(variables.GetJwtVariables(), repository)

	extractedToken, _ := jwtWrapper.ExtractToken(token)

	err := repository.BlacklistToken(extractedToken)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
