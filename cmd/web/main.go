package main

import (
	"auther/api/handler"
	"log"
	"net/http"
)

var router *http.ServeMux

func init() {
	router = handler.Router()
}

func hostStaticWeb() {
	router.Handle("GET /", http.FileServer(http.Dir("./web/static")))
}

func main() {
	log.Println("Listening for requests")
	hostStaticWeb()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
