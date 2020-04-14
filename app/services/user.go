package services

import (
	"fmt"
	"github.com/pedidosya/minesweeper-API/app/models"
	"github.com/pedidosya/minesweeper-API/app/repositories"
	"github.com/pedidosya/minesweeper-API/utils"
)

type IUserService interface {
	InsertUser(user *models.User) (interface{}, error)
	Login(userName string, password string) (string, error)
	UserLogin(token string) string
}

type UserService struct {
	userRepository repositories.IUserRepository
}

const encryptKey string = "123456789012345678901234"

var tokens = make(map[string]string)

func (service *UserService) InsertUser(user *models.User) (interface{}, error) {
	if err := encryptPasswordUser(user); err != nil {
		return nil, err
	}
	return service.userRepository.InsertUser(user)
}

func (service *UserService) Login(userName string, password string) (string, error) {
	encryptPassword, err := encryptPasswordWithUserName(password)
	if err != nil {
		return "", err
	}

	user, err := service.userRepository.Login(userName, encryptPassword)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", nil
	}

	uuid, err := utils.NewUUID()

	if err != nil {
		return "", err
	}

	tokens[uuid] = userName

	return uuid, nil
}

func (service *UserService) UserLogin(token string) string {
	return tokens[token]
}

func encryptPasswordUser(user *models.User) error {
	encryptPassword, err := encryptPasswordWithUserName(user.Password)
	if err != nil {
		return fmt.Errorf("it's not possible to encrypt the password, error :%v", err)
	}
	user.Password = encryptPassword
	return nil
}

func encryptPasswordWithUserName(password string) (string, error) {
	encryptPassword, err := utils.Encrypt(encryptKey, password)
	if err != nil {
		return "", fmt.Errorf("it's not possible to encrypt the password, error :%v", err)
	}
	return encryptPassword, nil
}

func decryptPasswordUser(user *models.User) error {
	decryptPassword, err := utils.Decrypt(encryptKey, user.Password)
	if err != nil {
		return fmt.Errorf("it's not possible to decrypt the password, error :%v", err)
	}
	user.Password = decryptPassword
	return nil
}

func NewUserService() IUserService {
	userRepository := repositories.NewUserRepository()
	return &UserService{
		userRepository: userRepository,
	}
}
