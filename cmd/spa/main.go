package main

import (
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"

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

	index, err := ioutil.ReadFile(cfg.Index)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(cfg.FrontAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Add("Content-type", "text/html")

		w.Write(index)
	})))
}
