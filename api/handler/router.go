package handler

import (
	"auther/api/handler/application"
	"auther/api/handler/oauth2"
	"net/http"
)

func registerControllers(mux *http.ServeMux) {
	application.RegisterControllers(mux)
	oauth2.RegisterControllers(mux)
}

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	registerControllers(mux)
	return mux
}
