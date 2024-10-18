package main

import (
	"auther/internal/api"
	"fmt"
	"net/http"
)

var router *http.ServeMux

func init() {
	router = api.Router()
}

func main() {
	fmt.Println("Listening for requests")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
