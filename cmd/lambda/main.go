package main

import (
	"auther/internal/api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"log"
	"net/http"
	"os"
)

var router *http.ServeMux

func init() {
	router = api.Router()
}

type LoggerHandler struct {
	handler http.Handler
}

func (l LoggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.Path)
	l.handler.ServeHTTP(w, r)
	log.Printf("Execution Complete")

}

func main() {
	log.Printf("Listening for requests")
	l := LoggerHandler{
		handler: router,
	}
	adapter := httpadapter.NewV2(l)
	adapter.StripBasePath(os.Getenv("API_GATEWAY_BASE_PATH"))
	lambda.Start(adapter.ProxyWithContext)
}
