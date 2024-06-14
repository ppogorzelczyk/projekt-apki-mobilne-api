package lists

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (l *ListsController) Delete(ctx *gin.Context) {
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
	err = service.DeleteList(listId, *userId)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusNoContent)
}
