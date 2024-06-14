package sharing

import (
	"buymeagiftapi/internal/domain"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetSharedUsersResponse struct {
	Users []domain.User `json:"users" binding:"required"`
}

func (s *SharingController) GetSharedUsers(ctx *gin.Context) {
	db := s.GetDatabase(ctx)
	userId := s.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	listIdRaw := ctx.Param("listId")
	listId, err := strconv.Atoi(listIdRaw)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	service := NewService(db)
	users, err := service.GetUsersThatListIsShareWith(listId, *userId)

	if err != nil {
		slog.Error("Error getting shared users: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := GetSharedUsersResponse{Users: users}
	ctx.JSON(http.StatusOK, response)
}
