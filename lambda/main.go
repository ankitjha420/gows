package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}
	return fmt.Sprintf("Called by - %s", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: "This is a secret endpoint",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUser(request)

		case "/login":
			return myApp.ApiHandler.LoginUser(request)

		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler) (request)

		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
