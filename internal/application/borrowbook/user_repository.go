package borrowbook

import (
	"context"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

//go:generate mockgen -destination=mocks/mock_user_repository.go -package=mocks . UserRepository

type UserRepository interface {
	FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error)
}
