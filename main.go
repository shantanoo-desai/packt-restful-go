package main

import (
	"fmt"
	"os"
	"net/http"

	"github.com/shantanoo-desai/packt-restful-go/handlers"
)

func main() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/users", handlers.UsersRouter)
	http.HandleFunc("/users/", handlers.UsersRouter)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
