package database

import (
	"buymeagiftapi/internal/constants"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func DatabaseMiddleware(db Database) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn := db.GetConnection()
		slog.Info("Setting database connection in context")
		ctx.Set(constants.GIN_CTX_DB_KEY, conn)
		ctx.Next()
	}
}
