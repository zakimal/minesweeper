package apiserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/zakimal/minesweeper/types"
	"net/http"
)

func (s *Services) createGame(w http.ResponseWriter, r *http.Request) {
	var game types.Game

	logger := s.logger.WithFields(logrus.Fields{
		"service": "game",
		"method": "create",
	})

	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		logger.Error(err)
		ErrInvalidJSON.Send(w)
		return
	}

	if err := s.GameService.Create(&game); err != nil {
		logger.WithField("err", err).Error("can not create game")
		ErrInternalServer.Send(w)
		return
	}

	Success(game, http.StatusCreated).Send(w)
}

func (s *Services) startGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	logger := s.logger.WithFields(logrus.Fields{
		"service": "game",
		"method": "start",
	})

	game, err := s.GameService.Start(name)
	if err != nil {
		logger.WithField("err", err).Error("can not start game")
		ErrInternalServer.Send(w)
		return
	}

	game_ := *game
	game_.Grid = nil

	Success(game_, http.StatusOK).Send(w)
}

func (s *Services) clickCell(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	logger := s.logger.WithFields(logrus.Fields{
		"service": "game",
		"method": "click",
	})

	var cellPos struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}

	if err := json.NewDecoder(r.Body).Decode(&cellPos); err != nil {
		logger.Error(err)
		ErrInvalidJSON.Send(w)
		return
	}

	game, err := s.GameService.Click(name, cellPos.Row, cellPos.Col)
	if err != nil {
		logger.WithField("err", err).Error("can not click cell")
		ErrInternalServer.Send(w)
		return
	}

	cell := game.Grid[cellPos.Row][cellPos.Col]

	if game.Status != "over" && game.Status != "won" {
		game.Grid = nil
	}

	var result struct {
		Cell types.Cell
		Game types.Game
	}

	result.Cell = cell
	result.Game = *game

	Success(&result, http.StatusOK).Send(w)
}