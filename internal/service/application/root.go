package application

import (
	"auther/internal/database/application"
	"context"
	"github.com/lestrrat-go/jwx/jwa"
	"log"
)

func init() {
	entity, err := application.Repository.GetByHashKey(context.TODO(), "root")
	if err != nil {
		panic(err)
	}

	if entity == nil {
		log.Println("root application does not exist, creating...")
		entity, err := Create("root")
		if err != nil {
			panic(err)
		}
		err = AddGeneratedKey(entity, jwa.ES256)
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("root application already exists")
	}
}
