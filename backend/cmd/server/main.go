package main

import (
	"fmt"
	"net/http"
	compiler "projetweb/backend/api/judge0api"
	news "projetweb/backend/api/news_api"
)

func main() {
	// firebaseService, err := database.InitFireBase()
	// if err != nil {
	// 	log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", err)
	// }
	// firebaseService.WriteDailyChallenge(2025, 3)
	// firebaseService.UpdateDailyQuestionDescription()
	// http.HandleFunc("/", handlers.HandleRoot)

	// fmt.Println("Server running on http://localhost:8080")
	// firebaseService.Client.Close()
	// log.Fatal(http.ListenAndServe(":8080", nil))

	http.HandleFunc("/api/news", news.HandleNews)
	http.HandleFunc("/api/compile", compiler.HandleCompiler)
	fmt.Println("Serveur en écoute sur :8080")
	http.ListenAndServe(":8080", nil)

}
