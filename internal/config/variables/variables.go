package variables

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func get(key string) string {
	return os.Getenv(key)
}

func GetJwtVariables() Jwt {
	expirationMinutes, err := strconv.ParseInt(os.Getenv("JWT_EXPIRATION_MINUTES"), 10, 64)

	if err != nil {
		slog.Error(fmt.Sprintf("Error parsing JWT_EXPIRATION_MINUTES: %v", err))
		expirationMinutes = 60
	}

	return Jwt{
		SecretKey:         os.Getenv("JWT_SECRET_KEY"),
		Issuer:            os.Getenv("JWT_ISSUER"),
		ExpirationMinutes: expirationMinutes,
	}
}
