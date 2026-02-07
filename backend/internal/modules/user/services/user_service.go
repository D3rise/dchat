package services

import (
	"context"
	"errors"

	"github.com/D3rise/dchat/internal/modules/user/entities"
	"github.com/D3rise/dchat/internal/modules/user/repositories"
	"golang.org/x/crypto/bcrypt"
)

var UsernameTakenErr = errors.New("username already taken")

type SignUpOptions struct {
	Username string
	Password string
}

type UserService interface {
	SignUp(ctx context.Context, options SignUpOptions) (entities.UserEntity, error)
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}

func (u *userServiceImpl) SignUp(ctx context.Context, options SignUpOptions) (entities.UserEntity, error) {
	existingUser, err := u.userRepository.GetUserByUsername(ctx, options.Username)
	if err != nil {
		return entities.UserEntity{}, err
	}

	if existingUser != nil {
		return entities.UserEntity{}, UsernameTakenErr
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(options.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.UserEntity{}, err
	}

	user, err := u.userRepository.CreateUser(ctx, repositories.CreateUserOptions{
		Username:     options.Username,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return entities.UserEntity{}, err
	}

	return user, nil
}
