package main

import (
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"GoWebApp/user"
)

func main() {
	router := user.NewRouter()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    router.PathPrefix("/static/").Handler(s)

	// launch server
	log.Fatal(http.ListenAndServe(":9000",
		handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
