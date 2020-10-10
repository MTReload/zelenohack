package main

import (
	"log"
	"net/http"
)

func main() {

	fs := http.FileServer(http.Dir("spa"))
	log.Fatal(http.ListenAndServe(":80", fs))
}