package oauth2

import (
	"net/http"
)

type GrantType string

const (
	AuthorizationCode GrantType = "authorization_code"
	Password          GrantType = "password"
	ClientCredentials GrantType = "client_credentials"
	RefreshToken      GrantType = "refresh_token"
)

type TokenRequest interface {
}

type AuthorizationCodeRequest struct {
	GrantType    GrantType `json:"grant_type"`
	Code         string    `json:"code"`
	RedirectURI  string    `json:"redirect_uri"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Scope        string    `json:"scope"`
}

type PasswordRequest struct {
	GrantType    GrantType `json:"grant_type"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Scope        string    `json:"scope"`
}

type ClientCredentialsRequest struct {
	GrantType    GrantType `json:"grant_type"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Scope        string    `json:"scope"`
}

type RefreshTokenRequest struct {
	GrantType    GrantType `json:"grant_type"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
}

// TokenResponse /token
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func ToTokenRequest(r *http.Request) (TokenRequest, error) {
	grantType := r.FormValue("grant_type")
	switch grantType {
	case "authorization_code":
		return &AuthorizationCodeRequest{
			GrantType:    GrantType(grantType),
			Code:         r.FormValue("code"),
			RedirectURI:  r.FormValue("redirect_uri"),
			ClientID:     r.FormValue("client_id"),
			ClientSecret: r.FormValue("client_secret"),
			Scope:        r.FormValue("scope"),
		}, nil
	case "password":
		return &PasswordRequest{
			GrantType:    GrantType(grantType),
			Username:     r.FormValue("username"),
			Password:     r.FormValue("password"),
			ClientID:     r.FormValue("client_id"),
			ClientSecret: r.FormValue("client_secret"),
			Scope:        r.FormValue("scope"),
		}, nil
	case "client_credentials":
		return &ClientCredentialsRequest{
			GrantType:    GrantType(grantType),
			ClientID:     r.FormValue("client_id"),
			ClientSecret: r.FormValue("client_secret"),
			Scope:        r.FormValue("scope"),
		}, nil
	case "refresh_token":
		return &RefreshTokenRequest{
			GrantType:    GrantType(grantType),
			RefreshToken: r.FormValue("refresh_token"),
			Scope:        r.FormValue("scope"),
		}, nil
	default:
		return nil, InvalidGrantError{}
	}
}
