package oauth2

import (
	"auther/internal/service/oauth2"
	"net/http"
)

func ToTokenRequest(r *http.Request) (oauth2.TokenRequest, error) {
	grantType := r.FormValue("grant_type")
	switch grantType {
	case "authorization_code":
		return &oauth2.AuthorizationCodeRequest{
			GrantType:    oauth2.GrantType(grantType),
			Code:         r.FormValue("code"),
			RedirectURI:  r.FormValue("redirect_uri"),
			ClientID:     r.FormValue("client_id"),
			ClientSecret: r.FormValue("client_secret"),
			Scope:        r.FormValue("scope"),
		}, nil
	case "password":
		return &oauth2.PasswordRequest{
			GrantType:    oauth2.GrantType(grantType),
			Username:     r.FormValue("username"),
			Password:     r.FormValue("password"),
			ClientID:     r.FormValue("client_id"),
			ClientSecret: r.FormValue("client_secret"),
			Scope:        r.FormValue("scope"),
		}, nil
	case "client_credentials":
		return &oauth2.ClientCredentialsRequest{
			GrantType:    oauth2.GrantType(grantType),
			ClientID:     r.FormValue("client_id"),
			ClientSecret: r.FormValue("client_secret"),
			Scope:        r.FormValue("scope"),
		}, nil
	case "refresh_token":
		return &oauth2.RefreshTokenRequest{
			GrantType:    oauth2.GrantType(grantType),
			RefreshToken: r.FormValue("refresh_token"),
			Scope:        r.FormValue("scope"),
		}, nil
	default:
		return nil, oauth2.InvalidGrantError{}
	}
}
