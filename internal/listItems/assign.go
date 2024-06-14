package listitems

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (li *ListItemsController) Assign(ctx *gin.Context) {
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

	amountRaw := ctx.Query("amount")
	amount, err := strconv.ParseFloat(amountRaw, 64)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	service := NewService(db)
	err = service.Assign(listId, itemId, *userId, amount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
