package main

import (
	"fmt"
	"log"
	"net/http"
	"projetweb/database"
	"projetweb/internal/handlers"
)

func main() {
	database.InitFireBase()
	database.WriteDailyChallenge(2025, 3)
	http.HandleFunc("/", handlers.HandleRoot)

	fmt.Println("Server running on http://localhost:8080")
	database.FirestoreClient.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
