package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"github.com/zakimal/minesweeper"
	"github.com/zakimal/minesweeper/storage/memorydb"
	"github.com/zakimal/minesweeper/types"
	"net/http"
)

type Services struct {
	logger *logrus.Logger
	GameService types.GameService
}

func Start(logger *logrus.Logger) error {
	db := memorydb.New()
	services := Services{
		logger:      logger,
		GameService: &minesweeper.GameService{Store: memorydb.NewGameStore(db)},
	}

	// setup router
	router := Router(&services)

	// middleware
	n := negroni.Classic()
	n.UseHandler(router)

	logger.Infoln("Minesweeper API Server running on port :3000")
	http.ListenAndServe(":3000", n)
	return nil
}

func Router(services *Services) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", services.healthcheck).Methods("GET")
	router.HandleFunc("/game", services.createGame).Methods("POST")
	router.HandleFunc("/game/{name}/start", services.startGame).Methods("POST")
	router.HandleFunc("/game/{name}/click", services.clickCell).Methods("POST")
	return router
}