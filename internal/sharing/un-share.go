package sharing

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *SharingController) Unshare(ctx *gin.Context) {
	db := s.GetDatabase(ctx)
	currentUserId := s.GetUserId(ctx)

	if db == nil || currentUserId == nil {
		return
	}

	listIdRaw := ctx.Param("listId")
	listId, err := strconv.Atoi(listIdRaw)

	userIdRaw := ctx.Param("userId")
	userId, err := strconv.Atoi(userIdRaw)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	service := NewService(db)
	result, err := service.Unshare(listId, *currentUserId, userId)

	if err != nil {
		slog.Error("Error un-sharing list: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch result {
	case UnshareUserNotOwner:
		slog.Error("User is not the owner of the list")
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	case UnshareNotShared:
		slog.Error("User is not shared with the list")
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Status(http.StatusNoContent)
}
