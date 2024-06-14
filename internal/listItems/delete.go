package listitems

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (li *ListItemsController) Delete(ctx *gin.Context) {
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

	service := NewService(db)
	err = service.Delete(listId, itemId, *userId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}
