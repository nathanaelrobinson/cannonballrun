package main

import (
	"net/http"
)

// Index Page as Follows are all URL Pathways

func index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./templates/index.html")
}

func newIndexx(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./templates/index.html")
}
