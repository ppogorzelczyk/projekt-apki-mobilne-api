package lists

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (l *ListsController) GetList(ctx *gin.Context) {
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

	service := NewService(db)
	list, err := service.GetListById(listId, *userId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if list == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, list)
}
