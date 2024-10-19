package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda-func/app"
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

func main() {
	myApp := app.NewApp()
	lambda.Start(myApp.ApiHandler.RegisterUser)
}
