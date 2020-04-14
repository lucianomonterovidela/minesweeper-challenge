package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

type StateGame int

const (
	Playing StateGame = iota + 1
	Paused
	Won
	Lose
)

type GameDto struct {
	Data []*Game `json:"data"`
}

type NewGameRequest struct {
	Rows    int `json:"rows"`
	Columns int `json:"columns"`
	Mines   int `json:"mines"`
}

type CellRequest struct {
	Row    int `json:"row"`
	Column int `json:"column"`}

type Game struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Board      *Board             `bson:"board" json:"board"`
	UserName   string             `bson:"user_name" json:"userName"`
	State      StateGame          `bson:"state" json:"state"`
	CreationAt time.Time          `bson:"creation_at" json:"createAt,omitempty"`
	EndedAt    *time.Time         `bson:"ended_at" json:"endedAt, omitempty"`
}

type Board struct {
	Rows    int     `bson:"rows" json:"rows"`
	Columns int     `bson:"columns" json:"columns"`
	Mines   int     `bson:"mines" json:"mines"`
	OpenCells int     `bson:"open_cells"`
	Cells   []*Cell `bson:"cells" json:"cells"`
}

type Cell struct {
	IsMined      bool `bson:"is_mined" json:"isMined"`
	MinesAround  int  `bson:"mines_around" json:"minesAround"`
	RedFlag      bool `bson:"red_flag" json:"redFlag"`
	QuestionFlag bool `bson:"question_flag" json:"questionFlag"`
	IsOpen       bool `bson:"is:_open" json:"isOpen"`
}

func (game *Game) UncoverCell(row int, column int) {
	minedCellIndex := game.Board.calculateCell(row, column)
	if game.Board.Cells[minedCellIndex].IsMined {
		game.State = Lose
		return
	}
	if !game.Board.Cells[minedCellIndex].IsOpen {
		game.Board.Cells[minedCellIndex].IsOpen = true
		game.Board.OpenCells = game.Board.OpenCells + 1

		game.recursivelyUncover(minedCellIndex, true)
		if game.Board.OpenCells + game.Board.Mines == game.Board.Rows*game.Board.Columns {
			game.State = Won
		}
	}
}

func (game *Game) recursivelyUncover(minedCellIndex int, firsTime bool) {
	if game.Board.Cells[minedCellIndex].IsOpen == false || firsTime {
		if game.Board.Cells[minedCellIndex].MinesAround == 0 {
			game.Board.Cells[minedCellIndex].IsOpen = true
			if !firsTime {
				game.Board.OpenCells = game.Board.OpenCells + 1
			}
			isNotOnLeftEdge := game.Board.isNotOnLeftEdge(minedCellIndex)
			isNotOnRightEdge := game.Board.isNotOnRightEdge(minedCellIndex)
			isNotOnTopEdge := game.Board.isNotOnTopEdge(minedCellIndex)
			isNotOnBottomEdge := game.Board.isNotOnBottomEdge(minedCellIndex)

			if isNotOnLeftEdge {
				left := minedCellIndex - 1
				game.recursivelyUncover(left, false)
			}

			if isNotOnRightEdge {
				right := minedCellIndex + 1
				game.recursivelyUncover(right, false)
			}

			if isNotOnTopEdge {
				up := minedCellIndex - game.Board.Columns
				game.recursivelyUncover(up, false)
			}

			if isNotOnBottomEdge {
				bottom := minedCellIndex + game.Board.Columns
				game.recursivelyUncover(bottom, false)
			}

			if isNotOnLeftEdge && isNotOnTopEdge {
				leftUp := minedCellIndex - 1 - game.Board.Columns
				game.recursivelyUncover(leftUp, false)
			}

			if isNotOnRightEdge && isNotOnTopEdge {
				rightUp := minedCellIndex + 1 - game.Board.Columns
				game.recursivelyUncover(rightUp, false)
			}

			if isNotOnLeftEdge && isNotOnBottomEdge {
				leftBottom := minedCellIndex - 1 + game.Board.Columns
				game.recursivelyUncover(leftBottom, false)
			}

			if isNotOnRightEdge && isNotOnBottomEdge {
				rightBottom := minedCellIndex + 1 + game.Board.Columns
				game.recursivelyUncover(rightBottom, false)
			}
		}
	}
}

func (board *Board) InitBoard() {
	board.fillEmptyCellsToBoard()
	board.fillMinesToBoard()
}

func (board *Board) fillEmptyCellsToBoard() {
	for i := 0; i < board.Rows*board.Columns; i++ {
		board.Cells[i] = &Cell{
			IsMined:      false,
			MinesAround:  0,
			RedFlag:      false,
			QuestionFlag: false,
			IsOpen:       false,
		}
	}
}

func (board *Board) fillMinesToBoard() {
	for i := 0; i < board.Mines; i++ {
		minedCellIndex := board.getRandomValueWithoutMine()
		board.fillMine(minedCellIndex)
	}
}

func (board *Board) fillMine(minedCellIndex int) {
	board.Cells[minedCellIndex].IsMined = true
	board.sumMinesAroundToAdjacent(minedCellIndex)
}

func (board *Board) sumMinesAroundToAdjacent(minedCellIndex int) {
	isNotOnLeftEdge := board.isNotOnLeftEdge(minedCellIndex)
	isNotOnRightEdge := board.isNotOnRightEdge(minedCellIndex)
	isNotOnTopEdge := board.isNotOnTopEdge(minedCellIndex)
	isNotOnBottomEdge := board.isNotOnBottomEdge(minedCellIndex)

	if isNotOnLeftEdge {
		left := minedCellIndex - 1
		board.Cells[left].MinesAround = board.Cells[left].MinesAround + 1
	}

	if isNotOnRightEdge {
		right := minedCellIndex + 1
		board.Cells[right].MinesAround = board.Cells[right].MinesAround + 1
	}

	if isNotOnTopEdge {
		up := minedCellIndex - board.Columns
		board.Cells[up].MinesAround = board.Cells[up].MinesAround + 1
	}

	if isNotOnBottomEdge {
		bottom := minedCellIndex + board.Columns
		board.Cells[bottom].MinesAround = board.Cells[bottom].MinesAround + 1
	}

	if isNotOnLeftEdge && isNotOnTopEdge {
		leftUp := minedCellIndex - 1 - board.Columns
		board.Cells[leftUp].MinesAround = board.Cells[leftUp].MinesAround + 1
	}

	if isNotOnRightEdge && isNotOnTopEdge {
		rightUp := minedCellIndex + 1 - board.Columns
		board.Cells[rightUp].MinesAround = board.Cells[rightUp].MinesAround + 1
	}

	if isNotOnLeftEdge && isNotOnBottomEdge {
		leftBottom := minedCellIndex - 1 + board.Columns
		board.Cells[leftBottom].MinesAround = board.Cells[leftBottom].MinesAround + 1
	}

	if isNotOnRightEdge && isNotOnBottomEdge {
		rightBottom := minedCellIndex + 1 + board.Columns
		board.Cells[rightBottom].MinesAround = board.Cells[rightBottom].MinesAround + 1
	}
}

func (board *Board) isNotOnLeftEdge(minedCellIndex int) bool {
	return minedCellIndex%board.Columns != 0
}

func (board *Board) isNotOnRightEdge(minedCellIndex int) bool {
	return (minedCellIndex+1)%board.Columns != 0
}

func (board *Board) isNotOnTopEdge(minedCellIndex int) bool {
	return minedCellIndex/board.Columns != 0
}

func (board *Board) isNotOnBottomEdge(minedCellIndex int) bool {
	return minedCellIndex/board.Columns < board.Rows-1
}

func (board *Board) getRandomValueWithoutMine() int {
	foundRandom := false
	random := 0
	for foundRandom != true {
		random = rand.Intn(board.Rows * board.Columns)
		fmt.Print(random)
		if !board.Cells[random].IsMined {
			foundRandom = true
		}
	}
	return random
}

func (board *Board) MarkRed(row int, column int) {
	board.Cells[board.calculateCell(row, column)].RedFlag = true
}

func (board *Board) MarkQuestion(row int, column int) {
	board.Cells[board.calculateCell(row, column)].QuestionFlag = true
}

func (board *Board) calculateCell(row int, column int) int {
	return ((row - 1)* board.Columns) + column - 1
}