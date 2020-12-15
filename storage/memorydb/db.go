package memorydb

import "github.com/zakimal/minesweeper/types"

type memoryDB struct {
	games map[string]*types.Game
}

func New() *memoryDB {
	return &memoryDB{
		games: make(map[string]*types.Game),
	}
}
