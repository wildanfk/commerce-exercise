package usecase

import (
	"context"
	"user-service/module/auth/entity"
)

//go:generate mockgen -destination=mock/repository.go -package=mock -source=repository.go

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
