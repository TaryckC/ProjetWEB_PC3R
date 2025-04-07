package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projetweb/backend/database"
	"projetweb/backend/internal/handlers"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from Go!",
	})
}

func main() {
	firebaseService, err := database.InitFireBase()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", err)
	}
	firebaseService.WriteDailyChallenge(2025, 4)
	firebaseService.UpdateDailyQuestionDescription()
	//firebaseService.WriteDailyAndWeeklyChallenges(2025, 4)
	//firebaseService.WriteChallengeComplementaryData()
	http.HandleFunc("/", handlers.HandleRoot)

	fmt.Println("Server running on http://localhost:8080")
	firebaseService.Client.Close()
	log.Fatal(http.ListenAndServe(":8080", nil))

	// http.HandleFunc("/api/hello", helloHandler)
	// http.HandleFunc("/api/news", news.HandleNews)
	// http.HandleFunc("/api/compile", compiler.HandleCompiler)
	// fmt.Println("Serveur en écoute sur :8080")
	// http.ListenAndServe(":8080", nil)

}
