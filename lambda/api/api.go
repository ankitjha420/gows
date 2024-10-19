package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api *ApiHandler) RegisterUser(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("username or password is empty")
	}

	doesUserExist, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("error registering the user %w", err)
	}

	if doesUserExist {
		return fmt.Errorf("a user with that username already exists")
	}

	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("error inserting the user %w", err)
	}

	return nil
}