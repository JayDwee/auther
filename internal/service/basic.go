package service

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

type BasicAuth struct {
	Username string
	Password string
}

func ToBasicAuth(r *http.Request) (auth BasicAuth, err error) {
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(r.Header.Get("Authorization"), "Basic "))
	if err != nil {
		log.Println(err.Error())
		return BasicAuth{}, err
	}
	parts := strings.Split(string(decoded), ":")
	return BasicAuth{Username: parts[0], Password: parts[1]}, nil
}
