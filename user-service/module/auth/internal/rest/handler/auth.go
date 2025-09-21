package handler

import (
	"encoding/json"
	"net/http"
	"user-service/internal/util/liberr"
	"user-service/internal/util/librest"
	"user-service/module/auth/entity"
)

type AuthHandler struct {
	authUsecase AuthUsecase
}

func NewAuthHandler(authUsecase AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (a *AuthHandler) Authentication(w http.ResponseWriter, r *http.Request) error {
	params := new(entity.AuthenticationUserRequest)
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		return liberr.NewBaseError(entity.ErrorInvalidBodyJSON)
	}

	authUser, err := a.authUsecase.Authentication(r.Context(), params)
	if err != nil {
		return err
	}

	code := http.StatusOK
	librest.WriteHTTPResponse(w, entity.AuthenticationUserResponse{
		Auth: authUser,
		Meta: &entity.Meta{
			HttpStatusCode: code,
		},
	}, code)
	return nil
}
