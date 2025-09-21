package handler

import (
	"net/http"
	"order-service/internal/util/liberr"
	"order-service/module/order/entity"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func UserAuth(r *http.Request, authServiceJWTSecret string) (*entity.User, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, liberr.NewBaseError(entity.ErrorTokenNotFound)
	}

	parts := strings.SplitN(authHeader, " ", -1)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, liberr.NewBaseError(entity.ErrorTokenInvalidBearer)
	}
	tokenString := parts[1]

	jwtSecret := []byte(authServiceJWTSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, liberr.NewBaseError(entity.ErrorTokenInvalid)
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, liberr.NewBaseError(entity.ErrorTokenInvalid)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, liberr.NewBaseError(entity.ErrorTokenExpired)
	}

	userID := claims["sub"].(string)
	return &entity.User{
		ID: userID,
	}, nil
}
