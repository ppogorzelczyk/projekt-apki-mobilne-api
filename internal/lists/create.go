package lists

import (
	"buymeagiftapi/internal/domain"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateList struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	EventDate   *string `json:"eventDate"`
}

func (l *ListsController) Create(ctx *gin.Context) {
	db := l.GetDatabase(ctx)
	userId := l.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	requestBody := CreateList{}
	err := ctx.ShouldBindJSON(&requestBody)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var eventDate *time.Time
	if requestBody.EventDate != nil {
		parsedEventDate, err := time.Parse(time.DateOnly, *requestBody.EventDate)
		if err != nil {
			slog.Error("Error parsing event date: %v", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		eventDate = &parsedEventDate
	}

	list := domain.List{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		EventDate:   eventDate,
		OwnerId:     *userId,
	}

	service := NewService(db)
	list, err = service.CreateList(list)

	if err != nil {
		slog.Error("Error creating list: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)
}
