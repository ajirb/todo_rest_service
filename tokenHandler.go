package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwtv4 "github.com/dgrijalva/jwt-go/v4"
)

var jwtKey = []byte("secret_key")

var user = "admin"

type Token struct {
	Token string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	jwtv4.StandardClaims
}

func getToken(w http.ResponseWriter, r *http.Request) {
	var expirationTime *jwtv4.Time = new(jwtv4.Time)
	expirationTime.Time = time.Now().Add(time.Minute * 30)
	claims := &Claims{
		Username: user,
		StandardClaims: jwtv4.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tok Token
	tok.Token = tokenString

	finalToken, err := json.Marshal(tok)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(finalToken))
}

func checkToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("X-Security-Token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request must contain X-Security-Token header"))
		return false
	}

	claims := &Claims{}

	tkn, err := jwtv4.ParseWithClaims(token, claims,
		func(t *jwtv4.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwtv4.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return false
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Token not valid"))
		return false
	}
	fmt.Println("correct token")
	return true
}
