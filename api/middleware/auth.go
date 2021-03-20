package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func Auth(w http.ResponseWriter, r *http.Request, params httprouter.Params) (int, error) {
	var (
		authHeader string
		ctx        context.Context
	)

	if authHeader = r.Header.Get("Authorization"); len(authHeader) < 8 {
		return http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest))
	}

	tokenString := authHeader[7:len(authHeader)]

	claims, err := verifyToken(tokenString)
	if err != nil {
		return http.StatusBadRequest, err
	}

	ctx = r.Context()
	ctx = context.WithValue(ctx, "userID", claims.(jwt.MapClaims)["id"])

	r = r.WithContext(ctx)

	return http.StatusOK, nil
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	// TODO, move secret to ENV
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims, err
}
