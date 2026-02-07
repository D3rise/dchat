package services

import (
	"testing"

	"github.com/D3rise/dchat/internal/modules/user/entities"
	"github.com/D3rise/dchat/internal/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_userServiceImpl_SignUp(t *testing.T) {
	t.Run("user does not exist", func(st *testing.T) {
		ctrl := gomock.NewController(st)
		userRepoMock := NewMockUserRepository(ctrl)
		service := NewUserService(userRepoMock)

		userRepoMock.
			EXPECT().
			GetUserByUsername(gomock.Any(), "testUsername").
			Times(1).
			Return(nil, nil)

		userRepoMock.EXPECT().CreateUser(
			gomock.Any(),
			testutils.StructMatcher().Field("Username", "testUsername"),
		).Times(1).Return(entities.UserEntity{Username: "testUsername", ID: "id"}, nil)

		user, err := service.SignUp(t.Context(), SignUpOptions{
			Username: "testUsername",
			Password: "testPassword",
		})

		assert.Nil(st, err)
		assert.Equal(st, user.Username, "testUsername", "usernames should be equal")
	})

	t.Run("user already exists", func(st *testing.T) {
		ctrl := gomock.NewController(st)
		userRepoMock := NewMockUserRepository(ctrl)
		service := NewUserService(userRepoMock)

		userRepoMock.
			EXPECT().
			GetUserByUsername(gomock.Any(), "testUsername").
			Times(1).
			Return(&entities.UserEntity{Username: "testUsername"}, nil)

		user, err := service.SignUp(t.Context(), SignUpOptions{
			Username: "testUsername",
			Password: "testPassword",
		})

		assert.Equal(st, user, entities.UserEntity{}, "user should be empty")
		assert.Error(st, err, UsernameTakenErr)
	})
}
