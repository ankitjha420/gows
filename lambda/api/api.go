package api

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusNotFound,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "missing fields",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	doesUserExist, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if doesUserExist {
		return events.APIGatewayProxyResponse{
			Body:       "username already exists",
			StatusCode: http.StatusConflict,
		}, err
	}

	user, err := types.NewUser(&registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("could not create new user - %w", err)
	}

	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user - %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "successfully inserted user",
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "invalid credentials",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	accessToken := types.CreateToken(user)
	msg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: http.StatusOK,
	}, nil

}
