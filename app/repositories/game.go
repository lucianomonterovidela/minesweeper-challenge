package repositories

import (
	"context"
	"fmt"
	"github.com/pedidosya/minesweeper-API/app/models"
	"github.com/pedidosya/minesweeper-API/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IGameRepository interface {
	NewGame(board *models.Board, userName string) (interface{}, error)
	PauseGame(gameId string) (bool, error)
	ResumeGame(gameId string, userName string) (*models.Game, error)
	UpdateGame(gameId string, game *models.Game) error
	GetGame(gameId string) (*models.Game, error)
	FindGames(user string) (*models.GameDto, error)
}

const gameCollection string = "games"

type GameRepository struct {
	dataBaseProvider infrastructure.IDataBaseProvider
}

func (gameRepository *GameRepository) NewGame(board *models.Board, userName string) (interface{}, error) {
	newGame := &models.Game{
		Board:      board,
		UserName:   userName,
		State:      models.Playing,
		CreationAt: time.Now(),
	}

	id, err := gameRepository.dataBaseProvider.Insert(gameCollection, newGame)
	if err != nil {
		return nil, err
	}

	newGame.Id = id.(primitive.ObjectID)
	return newGame, nil
}

func (gameRepository *GameRepository) UpdateGame(gameId string, game *models.Game) error {
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		return fmt.Errorf("not can create object_id: %v", err)
	}
	return gameRepository.dataBaseProvider.ReplaceById(gameCollection, objID, game)
}

func (gameRepository *GameRepository) PauseGame(gameId string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		return false, fmt.Errorf("not can create object_id: %v", err)
	}

	query := bson.M{}
	newState := bson.M{}
	newState["state"] = models.Paused
	query["$set"] = newState

	return gameRepository.dataBaseProvider.Update(gameCollection, objID, query)
}

func (gameRepository *GameRepository) ResumeGame(gameId string, userName string) (*models.Game, error) {
	game, err := gameRepository.GetGame(gameId)
	if err != nil {
		return nil, 	err
	}
	if game.UserName != userName {
		return nil, fmt.Errorf("the game belongs to another user")
	}
	gameRepository.resumeGame(gameId)
	return game, nil
}

func (gameRepository *GameRepository) FindGames(user string) (*models.GameDto, error) {
	// create empty map for query
	query := bson.M{}

	if user != "" {
		query["user_name"] = user
	}

	filter := bson.M{"$and": []bson.M{query}}
	cur, err := gameRepository.dataBaseProvider.Find(gameCollection, filter, nil)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	results := &models.GameDto{Data: []*models.Game{}}
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var game *models.Game
		if err := cur.Decode(&game); err != nil {
			return nil, fmt.Errorf("error marshal from database: %v", err)
		}
		results.Data = append(results.Data, game)
	}

	return results, nil
}

func (gameRepository *GameRepository) GetGame(gameId string) (*models.Game, error) {
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		return nil, fmt.Errorf("not can create object_id: %v", err)
	}

	var result *models.Game

	sr, err := gameRepository.dataBaseProvider.GetById(gameCollection, objID)
	if err != nil {
		return nil, err
	}

	if err := sr.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}


func (gameRepository *GameRepository) resumeGame(gameId string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		return false, fmt.Errorf("not can create object_id: %v", err)
	}

	query := bson.M{}
	newState := bson.M{}
	newState["state"] = models.Playing
	query["$set"] = newState

	return gameRepository.dataBaseProvider.Update(gameCollection, objID, query)
}

func NewGameRepository() IGameRepository {
	dataBaseProvider := infrastructure.NewDataBaseClient()
	return &GameRepository{
		dataBaseProvider: dataBaseProvider,
	}
}
