package oauth2

import (
	"auther/api/handler/util"
	"auther/internal/database/application"
	"auther/internal/service"
	"auther/internal/service/jwt"
	"auther/internal/service/oauth2"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func RegisterControllers(mux *http.ServeMux) {
	mux.HandleFunc("POST /oauth2/token", TokenController)
}

func TokenController(w http.ResponseWriter, r *http.Request) {
	log.Println("TokenController Called")
	tokenRequest, err := ToTokenRequest(r)
	if err != nil {
		util.JsonResponse(w, http.StatusBadRequest, err)
		return
	}
	appId := strings.Split(r.Host, ".")[0]
	app, err := application.Repository.GetByHashKey(r.Context(), appId) // Add cache

	var aud []string
	var sub string
	switch tokenRequest.(type) {
	case *oauth2.AuthorizationCodeRequest:
		util.JsonResponse(w, http.StatusBadRequest, fmt.Errorf("Not Implemented"))
		return
	case *oauth2.PasswordRequest:
		util.JsonResponse(w, http.StatusBadRequest, fmt.Errorf("Not Implemented"))
		return
	case *oauth2.ClientCredentialsRequest:
		basicAuth, err := service.ToBasicAuth(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		for _, client := range app.Clients {
			if basicAuth.Username == client.ClientId && basicAuth.Password == client.ClientSecret {
				aud = client.Audiences
				sub = client.ClientId
				break
			}
		}
		if sub == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	case *oauth2.RefreshTokenRequest:
		util.JsonResponse(w, http.StatusBadRequest, fmt.Errorf("Not Implemented"))
		return
	}

	issuer := fmt.Sprintf("https://%s/", r.Host)
	requestedScope := r.FormValue("scope")
	token, err := jwt.GenerateToken(app, sub, aud, issuer, requestedScope)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err.Error())
		return
	}

	response := oauth2.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       requestedScope,
	}

	util.JsonResponseNoCache(w, http.StatusOK, response)
	return
}
