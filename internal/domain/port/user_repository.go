package port

import (
	"context"
	"ip_detector/internal/domain/model"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}
