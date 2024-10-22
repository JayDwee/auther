package api

import (
	"auther/internal/auth"
	"auther/internal/auth/oauth2"
	"auther/internal/auth_server"
	"fmt"
	"net/http"
)

func TokenController(w http.ResponseWriter, r *http.Request) {
	tokenRequest, err := oauth2.ToTokenRequest(r)
	if err != nil {
		JsonErrorResponse(w, err)
		return
	}

	as, err := auth_server.AuthorizationServerRepository.GetByHashKey(r.Context(), r.PathValue("id")) // Add cache

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
		for _, client := range as.Clients {
			if basicAuth.Username == client.ClientId && basicAuth.Password == client.ClientSecret {
				aud = client.ClientId
				break
			}
		}
		if aud == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	case *oauth2.RefreshTokenRequest:
		JsonErrorResponse(w, fmt.Errorf("Not Implemented"))
		return
	}

	requestedScope := r.FormValue("scope")
	token, err := auth.GenerateToken(as, aud, requestedScope)
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
