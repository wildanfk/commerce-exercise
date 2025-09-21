package handler

import (
	"context"
	"user-service/module/auth/entity"
)

//go:generate mockgen -destination=mock/usecase.go -package=mock -source=usecase.go

type AuthUsecase interface {
	Authentication(ctx context.Context, params *entity.AuthenticationUserRequest) (*entity.AuthenticationUser, error)
}
