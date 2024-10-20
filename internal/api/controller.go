package api

import (
	"auther/internal/auth"
	"auther/internal/auth/oauth2"
	"fmt"
	"net/http"
	"os"
)

func TokenController(w http.ResponseWriter, r *http.Request) {
	tokenRequest, err := oauth2.ToTokenRequest(r)
	if err != nil {
		JsonErrorResponse(w, err)
		return
	}

	var aud string
	switch tokenRequest.(type) {
	case *oauth2.AuthorizationCodeRequest:
		JsonErrorResponse(w, fmt.Errorf("Not Implemented"))
		return
	case *oauth2.PasswordRequest:
		JsonErrorResponse(w, fmt.Errorf("Not Implemented"))
		return
	case *oauth2.ClientCredentialsRequest:
		basicAuth, err := auth.ToBasicAuth(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}
		if basicAuth.Username != os.Getenv("CLIENT_ID") || basicAuth.Password != os.Getenv("CLIENT_SECRET") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		aud = basicAuth.Username
	case *oauth2.RefreshTokenRequest:
		JsonErrorResponse(w, fmt.Errorf("Not Implemented"))
		return
	}

	requestedScope := r.FormValue("scope")
	token, err := auth.GenerateToken(aud, requestedScope)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err.Error())
		return
	}

	response := oauth2.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       requestedScope,
	}

	JsonResponseNoCache(w, response)
	return
}
