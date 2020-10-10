package main

import (
	"flag"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/mtreload/zh/internal/api"
	"github.com/mtreload/zh/pkg/config"
)

func main() {
	configPath := flag.String("c", "./config.yml", "path to config")
	flag.Parse()

	log := logrus.New()

	cfg := &config.Config{}

	err := config.Read(*configPath, cfg)
	if err != nil {
		log.WithError(err).Fatal("can't read config")
	}

	db, err := sqlx.Open("postgres", cfg.Database)
	if err != nil {
		log.WithError(err).Fatal("can't create connection")
	}
	err = db.Ping()
	if err != nil {
		log.WithError(err).Fatal("can't ping")
	}

	s := api.Server{DB: db}

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodPost, http.MethodGet, http.MethodHead, http.MethodPatch, http.MethodPut},
	}).Handler)

	r.Route("/api", func(r chi.Router) {
		r.Post("/game", s.CreateGame)
		r.Get("/game/{gameSN}", s.GetGame)
		r.Get("/game/{gameSN}/info", s.GetGameInfo)
		// also creates new player
		r.Get("/game/{gameSN}/info/{playerName}", s.GetPersonalGameInfo)

		r.Post("/game/{gameSN}/player/{playerName}/task/complete", s.CompleteTask)

		r.Post("/game/{gameSN}/tasks", s.CreateTasks)

		// r.Post("/player/{gameSN}")
	})

	log.Info("running server...")
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, r))
	log.Info("..stop")
}
