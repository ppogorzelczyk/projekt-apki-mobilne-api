package listitems

import (
	"buymeagiftapi/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateListItemRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Link        *string `json:"link"`
	Photo       *string `json:"imageUrl"`
	Price       float64 `json:"price" binding:"required"`
}

func (li *ListItemsController) Update(ctx *gin.Context) {
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

	itemIdRaw := ctx.Param("itemId")
	itemId, err := strconv.Atoi(itemIdRaw)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var listItemBody UpdateListItemRequest
	err = ctx.BindJSON(&listItemBody)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	listItem := domain.ListItem{
		Title:       listItemBody.Title,
		Description: listItemBody.Description,
		Link:        listItemBody.Link,
		Photo:       listItemBody.Photo,
		Price:       listItemBody.Price,
	}

	service := NewService(db)
	err = service.Update(listId, itemId, listItem, *userId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)

}
