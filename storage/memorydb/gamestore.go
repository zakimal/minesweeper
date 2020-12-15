package memorydb

import (
	"errors"
	"github.com/zakimal/minesweeper/types"
)

var _ types.GameStore = (*GameStore)(nil)

type GameStore struct {
	db *memoryDB
}

func NewGameStore(db *memoryDB) *GameStore {
	return &GameStore{db: db}
}

func (s *GameStore) Insert(game *types.Game) error {
	if _, ok := s.db.games[game.Name]; ok {
		return errors.New("game already exists")
	}

	s.db.games[game.Name] = game
	return nil
}

func (s *GameStore) Update(game *types.Game) error {
	if _, ok := s.db.games[game.Name]; !ok {
		return errors.New("game does not exit")
	}

	g := *game
	s.db.games[game.Name] = &g
	return nil
}

func (s *GameStore) GetByName(name string) (*types.Game, error) {
	if game, ok := s.db.games[name]; ok {
		g := *game
		return &g, nil
	}

	return nil, errors.New("game not found")
}