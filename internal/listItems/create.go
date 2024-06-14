package listitems

import (
	"buymeagiftapi/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateListItemRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Link        *string `json:"link"`
	Photo       *string `json:"photo"`
}

func (li *ListItemsController) Create(ctx *gin.Context) {
	db := li.GetDatabase(ctx)
	userId := li.GetUserId(ctx)

	if db == nil || userId == nil {
		return
	}

	listIdRaw := ctx.Param("listId")
	listId, err := strconv.Atoi(listIdRaw)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req CreateListItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	item := domain.ListItem{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Link:        req.Link,
		Photo:       req.Photo,
		ListId:      listId,
	}

	service := NewService(db)
	listItem, err := service.Create(item, *userId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if listItem == nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.JSON(http.StatusCreated, item)
}
