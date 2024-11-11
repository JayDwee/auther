package main

import (
	"auther/internal/service/jwt"
	"github.com/lestrrat-go/jwx/jwa"
	"log"
)

func main() {
	jwkNew, err := jwt.GenerateJWK(jwa.NoSignature)
	if err != nil {
		return
	}
	log.Printf("%v", jwkNew)

}
