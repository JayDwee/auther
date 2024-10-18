package main

import (
	"auther/internal/api"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
	"os"
)

var router *http.ServeMux

func init() {
	router = api.Router()
}

func main() {
	fmt.Println("Listening for requests")
	adapter := httpadapter.NewV2(router)
	adapter.StripBasePath(os.Getenv("API_GATEWAY_BASE_PATH"))
	lambda.Start(adapter.ProxyWithContext)
}
