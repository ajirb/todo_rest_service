package main

import (
	"fmt"
	"net/http"
)

func main() {
	todoHandlers := newTodoHandlers()
	initialDbSetup()
	http.HandleFunc("/getToken", getToken)
	http.HandleFunc("/todo/", todoHandlers.mapIdByMethod)
	http.HandleFunc("/todo", todoHandlers.mapRequest)
	http.HandleFunc("/", greet)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home url hit")
	w.Write([]byte("Welcome to TODO service"))
}
