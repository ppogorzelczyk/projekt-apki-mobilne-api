package lists

import (
	"buymeagiftapi/internal/domain"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllResponse struct {
	Lists []domain.List `json:"lists" binding:"required"`
}

func (l *ListsController) GetMyLists(ctx *gin.Context) {
	db := l.GetDatabase(ctx)
	userId := l.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	service := NewService(db)
	lists, err := service.GetUserLists(*userId)

	if err != nil {
		slog.Error("Error getting user lists: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := GetAllResponse{Lists: lists}
	ctx.JSON(http.StatusOK, response)
}
