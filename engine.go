package minesweeper

import (
	"errors"
	"fmt"
	"github.com/zakimal/minesweeper/types"
	"math/rand"
	"time"
)

const (
	defaultRows = 9
	defaultCols = 9
	defaultMines = 18
	maxRows = 30
	maxCols = 30
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ types.GameService = (*GameService)(nil)

type GameService struct {
	Store types.GameStore
}

func (s *GameService) Create(game *types.Game) error {
	if game.Name == "" {
		return errors.New("can not create game without name")
	}

	if game.Rows == 0 {
		game.Rows = defaultRows
	}
	if game.Cols == 0 {
		game.Cols = defaultCols
	}
	if game.Mines == 0 {
		game.Mines = defaultMines
	}

	if maxRows < game.Rows {
		game.Rows = maxRows
	}
	if maxCols < game.Cols {
		game.Cols = maxCols
	}

	if game.Mines < game.Rows * game.Cols { // too many mines!
		game.Mines = game.Rows * game.Cols
	}

	game.Status = "new"

	err := s.Store.Insert(game)
	return err
}

func (s *GameService) Start(name string) (*types.Game, error) {
	game, err := s.Store.GetByName(name)
	if err != nil {
		return nil, err
	}

	createBoard(game)

	game.Status = "started"
	err = s.Store.Update(game)
	fmt.Printf("%#v\n", game.Grid)

	return game, err
}

func (s *GameService) Click(name string, i, j int) (*types.Game, error) {
	game, err := s.Store.GetByName(name)
	if err != nil {
		return nil, err
	}

	if err := clickCell(game, i, j); err != nil {
		return nil, err
	}

	if err := s.Store.Update(game); err != nil {
		return nil, err
	}

	return game, nil
}

func createBoard(game *types.Game) {
	numCells := game.Rows * game.Cols
	cells := make(types.CellGrid, numCells)

	// FIXME: more efficient
	i := 0
	for i < game.Mines {
		idx := rand.Intn(numCells)
		if !cells[idx].Mine {
			cells[idx].Mine = true;
			i++
		}
	}

	game.Grid = make([]types.CellGrid, game.Rows)
	for row := range game.Grid {
		game.Grid[row] = cells[(game.Cols * row):game.Cols*(row + 1)]
	}

	for i, row := range game.Grid {
		for j, cell := range row {
			if cell.Mine {
				setAdjValues(game, i, j)
			}
		}
	}
}

func setAdjValues(game *types.Game, i, j int) {
	for z := i - 1; z < i + 2; z++ {
		if z < 0 || game.Rows - 1 < z { // out of board
			continue
		}
		for w := j - 1; w < j + 2; w++ {
			if w < 0 || game.Cols-1 < w { // out of board
				continue
			}

			game.Grid[z][w].Value++;
		}
	}
}

// TODO: add feature to open connected "zero" cells at once
func clickCell(game *types.Game, i, j int) error {
	if game.Grid[i][j].Clicked {
		return errors.New("cell already clicked")
	}

	game.Grid[i][j].Clicked = true
	game.Clicks++
	if game.Grid[i][j].Mine {
		game.Status = "over"
		return nil
	}

	if checkWon(game) {
		game.Status = "won"
	}

	return nil
}

func checkWon(game *types.Game) bool {
	return game.Clicks == ((game.Rows * game.Cols) - game.Mines)
}