package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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

	render.DefaultResponder(w, r, *game)
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

	game, err := model.GetGameInfoForPlayer(r.Context(), s.DB, gameSN, playerName)
	if err != nil {
		log.WithError(err).Error("can't get personal game info")
		render.DefaultResponder(w, r, err)
		return
	}

	render.DefaultResponder(w, r, game)
}
