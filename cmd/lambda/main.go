package main

import (
	"auther/internal/api"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"net/http"
)

var router *http.ServeMux

func init() {
	router = api.Router()
}

func main() {
	fmt.Println("Listening for requests")
	lambda.Start(httpadapter.NewV2(router).ProxyWithContext)
}
