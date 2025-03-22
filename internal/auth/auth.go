package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("not authenticated")
	}

	authHeaderValues := strings.Split(authHeader, " ")
	if len(authHeaderValues) != 2 {
		return "", errors.New("malformed auth header")
	}

	if authHeaderValues[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}

	return authHeaderValues[1], nil
}
