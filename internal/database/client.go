package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itsamirhn/dongetobede/internal/database/entities"
)

var ErrNotFound = errors.New("not found")

type Client interface {
	AddUser(ctx context.Context, user *entities.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	GetUsersByIDs(ctx context.Context, ids []int64) ([]*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	AddDong(ctx context.Context, dong *entities.Dong) (*primitive.ObjectID, error)
	GetDongByID(ctx context.Context, id primitive.ObjectID) (*entities.Dong, error)
	UpdateDong(ctx context.Context, dong *entities.Dong) error
}
