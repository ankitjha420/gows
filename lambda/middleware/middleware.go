package middleware

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// extract headers ->
		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body:       "Missing JWT Auth",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		// extract claims ->
		claims, err := parseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "User is unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, err
		}

		expires := int64(claims["expires"].(float64))
		if expires < (time.Now()).Unix() {
			return events.APIGatewayProxyResponse{
				Body:       "JWT expired",
				StatusCode: http.StatusUnauthorized,
			}, err
		}

		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := "secret"
		return []byte(secret), nil
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized: token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims of unauthorized type")
	}

	return claims, nil
}
