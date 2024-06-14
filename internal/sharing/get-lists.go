package sharing

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAllResponse struct {
	Lists []SharedList `json:"lists" binding:"required"`
}

func (s *SharingController) GetSharedLists(ctx *gin.Context) {
	db := s.GetDatabase(ctx)
	userId := s.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	service := NewService(db)
	lists, err := service.GetSharedLists(*userId)

	if err != nil {
		slog.Error("Error getting user lists: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := GetAllResponse{Lists: lists}
	ctx.JSON(http.StatusOK, response)
}
