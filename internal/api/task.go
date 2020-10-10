package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

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
}

func (s *Server) CreateTasks(w http.ResponseWriter, r *http.Request) {

}
