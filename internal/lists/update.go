package lists

import (
	"buymeagiftapi/internal/domain"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateListRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	EventDate   *string `json:"eventDate"`
}

func (l *ListsController) Update(ctx *gin.Context) {
	db := l.GetDatabase(ctx)
	userId := l.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	listIdRaw := ctx.Param("listId")
	listId, err := strconv.Atoi(listIdRaw)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var listBody UpdateListRequest
	err = ctx.BindJSON(&listBody)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	list := domain.List{
		Title:       listBody.Title,
		Description: listBody.Description,
	}

	if listBody.EventDate != nil {
		parsedEventDate, err := time.Parse(time.DateOnly, *listBody.EventDate)
		if err != nil {
			slog.Error("Error parsing event date: %v", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		list.EventDate = &parsedEventDate
	}

	service := NewService(db)
	err = service.UpdateList(listId, list, *userId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}
