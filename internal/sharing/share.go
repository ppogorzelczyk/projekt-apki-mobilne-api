package sharing

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShareRequest struct {
	Email string `json:"email" binding:"required"`
}

func (s *SharingController) Share(ctx *gin.Context) {
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

	var req ShareRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	service := NewService(db)
	shareResult, err := service.Share(listId, *userId, req.Email)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch shareResult {
	case ShareSuccess:
		ctx.Status(http.StatusOK)
	case ShareUserNotFound:
		ctx.Status(http.StatusNotFound)
	case ShareAlreadyShared:
		ctx.Status(http.StatusConflict)
	default:
		ctx.Status(http.StatusInternalServerError)
	}
}
