package main

import (
	"fmt"
	"log"
	"net/http"
	"projetweb/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HandleRoot)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
