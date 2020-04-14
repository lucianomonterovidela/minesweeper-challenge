package repositories

import (
	"context"
	"fmt"
	"github.com/pedidosya/minesweeper-API/app/models"
	"github.com/pedidosya/minesweeper-API/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
)

type IUserRepository interface {
	InsertUser(user *models.User) (interface{}, error)
	Login(userName string, password string) (*models.User, error)
	GetUser(userName string) (*models.User, error)
}

const userCollection string = "users"

type UserRepository struct {
	dataBaseProvider infrastructure.IDataBaseProvider
}

func (userRepository *UserRepository) InsertUser(user *models.User) (interface{}, error) {
	return userRepository.dataBaseProvider.Insert(userCollection, user)
}

func (userRepository *UserRepository) Login(userName string, password string) (*models.User, error) {
	query := bson.M{}

	query["_id"] = userName
	query["password"] = password

	filter := bson.M{"$and": []bson.M{query}}
	cur, err := userRepository.dataBaseProvider.Find(userCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	var user *models.User
	for cur.Next(context.TODO()) {
		if err := cur.Decode(&user); err != nil {
			return nil, fmt.Errorf("error marshal from database: %v", err)
		}
	}

	return user, nil
}

func (userRepository *UserRepository) GetUser(userName string) (*models.User, error) {
	var user *models.User

	sr, err := userRepository.dataBaseProvider.GetById(userCollection, userName)
	if err != nil {
		return nil, err
	}

	if err := sr.Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository() IUserRepository {
	dataBaseProvider := infrastructure.NewDataBaseClient()
	return &UserRepository{
		dataBaseProvider: dataBaseProvider,
	}
}
