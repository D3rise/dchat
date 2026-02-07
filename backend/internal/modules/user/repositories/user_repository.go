package repositories

import (
	"context"

	"github.com/D3rise/dchat/internal/infrastructure/database"
	"github.com/D3rise/dchat/internal/modules/user/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*entities.UserEntity, error)
	CreateUser(ctx context.Context, options CreateUserOptions) (entities.UserEntity, error)
}

type userServiceImpl struct {
	db *gorm.DB
}

func NewUserRepository(db database.Database) UserRepository {
	return &userServiceImpl{
		db: db.Gorm,
	}
}

type CreateUserOptions struct {
	Username     string
	PasswordHash string
}

func (u *userServiceImpl) GetUserByUsername(ctx context.Context, username string) (*entities.UserEntity, error) {
	foundUser, err := gorm.G[entities.UserEntity](u.db).Where("username = ?", username).First(ctx)

	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (u *userServiceImpl) CreateUser(ctx context.Context, options CreateUserOptions) (entities.UserEntity, error) {
	newUser := entities.UserEntity{
		Username:     options.Username,
		PasswordHash: options.PasswordHash,
	}

	result := gorm.WithResult()
	err := gorm.G[entities.UserEntity](u.db, result).Create(ctx, &newUser)
	if err != nil {
		return entities.UserEntity{}, err
	}

	return newUser, nil
}
