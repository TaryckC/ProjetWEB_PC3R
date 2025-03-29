package main

import (
	"fmt"
	"log"
	"net/http"
	"projetweb/database"
	"projetweb/internal/handlers"
)

func main() {
	firebaseService, err := database.InitFireBase()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", err)
	}
	firebaseService.WriteDailyChallenge(2025, 3)
	firebaseService.UpdateDailyQuestionDescription()
	http.HandleFunc("/", handlers.HandleRoot)

	fmt.Println("Server running on http://localhost:8080")
	firebaseService.Client.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
