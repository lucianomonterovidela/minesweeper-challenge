package repositories

import (
	"fmt"
	"github.com/pedidosya/minesweeper-API/app/models"
	"github.com/pedidosya/minesweeper-API/mocks/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUserRepository_InsertUser(t *testing.T) {

	dataBaseProviderMock := &mocks.DataBaseProviderMock{}

	type args struct {
		user *models.User
	}
	tests := []struct {
		name        string
		initMocks   func()
		args        args
		assertMocks func(*testing.T)
		assertError func(*testing.T, error)
	}{
		{
			name: "Error - Insert User",
			initMocks: func() {
				dataBaseProviderMock.On("Insert", userCollection, mock.Anything).
					Return(nil, fmt.Errorf("error when invoke database")).Once()
			},
			args: args{
				user: &models.User{
					UserName:"luciano",
					Password:"password",
				},
			},
			assertMocks: func(t *testing.T) {
				dataBaseProviderMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.NotNil(t, e)
			},
		},
		{
			name: "Success - Insert User",
			initMocks: func() {
				dataBaseProviderMock.On("Insert", userCollection, mock.Anything).
					Return(&models.User{}, nil).Once()
			},
			args: args{
				user: &models.User{
					UserName:"luciano",
					Password:"password",
				},
			},
			assertMocks: func(t *testing.T) {
				dataBaseProviderMock.AssertExpectations(t)
			},
			assertError: func(t *testing.T, e error) {
				assert.Nil(t, e)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &UserRepository{
				dataBaseProvider: dataBaseProviderMock,
			}
			tt.initMocks()
			_, err := userRepository.InsertUser(tt.args.user)
			tt.assertMocks(t)
			tt.assertError(t, err)
		})
	}
}