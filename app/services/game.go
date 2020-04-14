package services

import (
	"fmt"
	"github.com/pedidosya/minesweeper-API/app/models"
	"github.com/pedidosya/minesweeper-API/app/repositories"
)

type IGameService interface {
	NewGame(rows int, columns int, mines int, userName string) (interface{}, error)
	PauseGame(id string) (bool, error)
	ResumeGame(id string, userName string) (*models.Game, error)
	MarkRed(id string, row int, column int) (bool, error)
	MarkQuestion(id string, row int, column int) (bool, error)
	Uncover(id string, row int, column int) (*models.Game, error)
	FindGames(user string) (*models.GameDto, error)
}

type GameService struct {
	gameRepository repositories.IGameRepository
}

func (service *GameService) NewGame(boardSizeX int, boardSizeY int, mines int, userName string) (interface{}, error) {
	board := generateBoard(boardSizeX, boardSizeY, mines)
	return service.gameRepository.NewGame(board, userName)
}

func (service *GameService) PauseGame(id string) (bool, error) {
	return service.gameRepository.PauseGame(id)
}

func (service *GameService) ResumeGame(id string, userName string) (*models.Game, error) {
	return service.gameRepository.ResumeGame(id, userName)
}

func (service *GameService) MarkRed(id string, row int, column int) (bool, error) {
	game, err := service.gameRepository.GetGame(id)
	if err != nil {
		return false, err
	}
	if err := validateSizeGameToAction(game.Board, row, column); err != nil {
		return false, err
	}
	game.Board.MarkRed(row, column)
	go service.gameRepository.UpdateGame(id, game)
	return true, nil
}

func (service *GameService) MarkQuestion(id string, row int, column int) (bool, error) {
	game, err := service.gameRepository.GetGame(id)
	if err != nil {
		return false, err
	}
	if err := validateSizeGameToAction(game.Board, row, column); err != nil {
		return false, err
	}
	game.Board.MarkQuestion(row, column)
	go service.gameRepository.UpdateGame(id, game)
	return true, nil
}

func (service *GameService) Uncover(id string, row int, column int) (*models.Game, error) {
	game, err := service.gameRepository.GetGame(id)
	if err != nil {
		return nil, err
	}
	if err := validateSizeGameToAction(game.Board, row, column); err != nil {
		return nil, err
	}

	game.UncoverCell(row, column)
	go service.gameRepository.UpdateGame(id, game)
	return game, nil
}

func validateSizeGameToAction(board *models.Board, row int, column int) error {
	if board.Rows < row {
		return fmt.Errorf("the row number must be less than: %d", board.Rows)
	}

	if board.Columns < column {
		return fmt.Errorf("the column number must be less than: %d", board.Columns)
	}
	return nil
}

func (service *GameService) FindGames(user string) (*models.GameDto, error) {
	return service.gameRepository.FindGames(user)
}

func generateBoard(rows int, columns int, mines int) *models.Board {
	var board = &models.Board{
		Rows:    rows,
		Columns: columns,
		OpenCells: 0,
		Mines:   mines,
		Cells:   make([]*models.Cell, rows*columns),
	}
	board.InitBoard()

	return board
}

func NewGameService() IGameService {
	gameRepository := repositories.NewGameRepository()
	return &GameService{
		gameRepository: gameRepository,
	}
}
