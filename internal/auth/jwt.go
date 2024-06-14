package auth

import (
	"buymeagiftapi/internal/config/variables"
	"buymeagiftapi/internal/domain"
	"database/sql"
	"errors"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	ERROR_CLAIMS_PARSE             = "Couldn't parse claims"
	ERROR_TOKEN_EXPIRED            = "Token is expired"
	ERROR_TOKEN_INVALID            = "Invalid token, 2 parts expected"
	ERROR_TOKEN_PARSE              = "Invalid token format"
	ERROR_TOKEN_BLACKLIST_CHECK    = "Error checking token blacklist"
	ERROR_TOKEN_BLACKLISTED        = "Token is blacklisted"
	ERROR_TOKEN_GENERATION         = "Error generating token"
	ERROR_TOKEN_REFRESH_GENERATION = "Error generating refresh token"
	refreshTokenBufferMinutes      = 10
)

type JwtWrapper struct {
	secretKey         string
	issuer            string
	expirationMinutes int64
	repository        tokenRepository
}

type tokenRepository interface {
	IsTokenBlacklisted(token string) (bool, error)
}

type JwtClaim struct {
	Email  string
	UserId int
	jwt.RegisteredClaims
}

func NewJwtWrapper(
	variables variables.Jwt,
	repository tokenRepository,
) *JwtWrapper {
	return &JwtWrapper{
		secretKey:         variables.SecretKey,
		issuer:            variables.Issuer,
		expirationMinutes: variables.ExpirationMinutes,
		repository:        repository,
	}
}

func (j *JwtWrapper) GenerateToken(user domain.User) (signedToken string, err error) {
	claims := j.generateClaims(user, nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JwtWrapper) GenerateRefreshToken(user domain.User) (signedToken string, err error) {
	expirationMinutes := j.expirationMinutes + refreshTokenBufferMinutes
	claims := j.generateClaims(user, &expirationMinutes)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(db *sql.DB, signedToken string) (claims *JwtClaim, validationError error, err error) {
	token, _ := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
	)

	if token == nil {
		validationError = errors.New(ERROR_TOKEN_PARSE)
		return
	}

	claims, ok := token.Claims.(*JwtClaim)

	if !ok {
		err = errors.New(ERROR_CLAIMS_PARSE)
		return
	}

	if claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		validationError = errors.New(ERROR_TOKEN_EXPIRED)
		return
	}

	isBlacklisted, err := j.repository.IsTokenBlacklisted(signedToken)

	if err != nil {
		return
	}

	if isBlacklisted {
		validationError = errors.New(ERROR_TOKEN_BLACKLISTED)
	}

	return
}

func (j *JwtWrapper) ExtractToken(token string) (string, error) {
	tokenParts := strings.Split(token, "Bearer ")
	if len(tokenParts) != 2 {
		err := errors.New(ERROR_TOKEN_INVALID)
		return "", err
	}

	extractedToken := strings.TrimSpace(tokenParts[1])

	return extractedToken, nil
}

func (j *JwtWrapper) generateClaims(
	user domain.User,
	expirationMinutes *int64,
) (claims *JwtClaim) {
	if expirationMinutes == nil {
		expirationMinutes = &j.expirationMinutes
	}

	claims = &JwtClaim{
		Email:  user.Email,
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(*expirationMinutes))),
			Issuer:    j.issuer,
		},
	}
	return
}
