package api

import (
	"net/http"
)

func registerControllers(mux *http.ServeMux) {
	mux.HandleFunc("POST /as/{id}/token", TokenController)
}

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	registerControllers(mux)
	return mux
}
