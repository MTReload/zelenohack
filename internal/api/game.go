package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/mtreload/zh/internal/model"
)

func (s *Server) CreateGame(w http.ResponseWriter, r *http.Request) {
	var err error

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("can't read data")
		render.DefaultResponder(w, r, err)
		return
	}

	var game = &model.Game{}

	err = json.Unmarshal(data, game)
	if err != nil {
		log.WithError(err).Error("can't unmarshal game")
		render.DefaultResponder(w, r, err)
		return
	}

	if game.Name == "" {
		render.DefaultResponder(w, r, "empty name")
	}
	game.ShortName = strings.Replace(uuid.New().String(), "-", "", -1)[:6]

	game, err = model.NewGame(r.Context(), s.DB, *game)
	if err != nil {
		log.WithError(err).Error("can't create game")
		render.DefaultResponder(w, r, err)
		return
	}

	render.DefaultResponder(w, r, *game)
}

func (s *Server) GetGame(w http.ResponseWriter, r *http.Request) {
	var err error

	gameSN := chi.URLParam(r, "gameSN")

	game, err := model.GameByShortName(r.Context(), s.DB, gameSN)
	if err != nil {
		log.WithError(err).Error("can't get game")
		render.DefaultResponder(w, r, err)
		return
	}

	if game != nil {
		render.DefaultResponder(w, r, *game)
	}

}

func (s *Server) GetGameInfo(w http.ResponseWriter, r *http.Request) {
	var err error

	gameSN := chi.URLParam(r, "gameSN")

	game, err := model.GameInfo(r.Context(), s.DB, gameSN)
	if err != nil {
		log.WithError(err).Error("can't get game info")
		render.DefaultResponder(w, r, err)
		return
	}

	render.DefaultResponder(w, r, game)
}
func (s *Server) GetPersonalGameInfo(w http.ResponseWriter, r *http.Request) {
	var err error

	gameSN := chi.URLParam(r, "gameSN")
	playerName := chi.URLParam(r, "playerName")

	gameInfo, err := model.GetGameInfoForPlayer(r.Context(), s.DB, gameSN, playerName)
	if err != nil {
		log.WithError(err).Error("can't get personal game info")
		render.DefaultResponder(w, r, err)
		return
	}

	// crating new player
	if gameInfo == nil {
		gameInfo, err = model.NewPlayer(r.Context(), s.DB, playerName, gameSN)
	}

	render.DefaultResponder(w, r, gameInfo)
}
