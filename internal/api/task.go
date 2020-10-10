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

// /api/game/{gameSN}/player/{playerName}/task/complete
func (s *Server) CompleteTask(w http.ResponseWriter, r *http.Request) {
	gameSN := chi.URLParam(r, "gameSN")
	playerName := chi.URLParam(r, "playerName")

	err := model.CompleteTask(r.Context(), s.DB, playerName, gameSN)
	if err != nil {
		render.DefaultResponder(w, r, map[string]string{"error": err.Error()})
		return
	}
}

type TasksRequest struct {
	Tasks []model.Task `json:"tasks"`
}

// /game/{gameSN}/tasks
func (s *Server) CreateTasks(w http.ResponseWriter, r *http.Request) {
	var req TasksRequest

	gameSN := chi.URLParam(r, "gameSN")

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.WithError(err).Error()
		render.DefaultResponder(w, r, err.Error())
		return
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		log.WithError(err).Error()
		render.DefaultResponder(w, r, err.Error())
		return
	}

	err = model.NewTaskBatch(r.Context(), s.DB, gameSN, req.Tasks)
	if err != nil {
		log.WithError(err).Error()
		render.DefaultResponder(w, r, err.Error())
		return
	}
}
