package controllers

import (
	"buymeagiftapi/internal/constants"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

// Tries to get the database connection from the gin context.
//
// Returns nil if it doesn't exist or if it's not a *sql.DB.
func (u *BaseController) GetDatabase(ctx *gin.Context) *sql.DB {
	dbRaw, exists := ctx.Get(constants.GIN_CTX_DB_KEY)

	if !exists {
		slog.Error("Database connection not found in context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	db, ok := dbRaw.(*sql.DB)

	if !ok {
		slog.Error("Database connection is not a *sql.DB")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	return db
}

// Tries to get the user id from the gin context.
//
// Returns nil if it doesn't exist or if it's not a uuid.UUID.
func (u *BaseController) GetUserId(ctx *gin.Context) *int {
	userIdRaw, exists := ctx.Get(constants.GIN_CLAIMS_USER_ID_KEY)

	if !exists {
		slog.Error("User id not found in context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	userId, ok := userIdRaw.(int)

	if !ok {
		slog.Error("User id is not an int")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return nil
	}

	return &userId
}
