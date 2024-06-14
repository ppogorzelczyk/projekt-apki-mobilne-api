package auth

import (
	"buymeagiftapi/internal/config/variables"
	"buymeagiftapi/internal/constants"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenRepository := NewRepository(db)
		jwtWrapper := NewJwtWrapper(variables.GetJwtVariables(), tokenRepository)
		clientToken, err := jwtWrapper.ExtractToken(clientToken)
		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		claims, validationError, err := jwtWrapper.ValidateToken(db, clientToken)

		if validationError != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": validationError.Error()})
			return
		}

		if err != nil {
			slog.Error(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(constants.GIN_CLAIMS_EMAIL_KEY, claims.Email)
		c.Set(constants.GIN_CLAIMS_USER_ID_KEY, claims.UserId)

		c.Next()
	}
}
