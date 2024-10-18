package api

import (
	"net/http"
)

func registerControllers(mux *http.ServeMux) {
	mux.HandleFunc("POST /token", TokenController)
}

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	registerControllers(mux)
	return mux
}
