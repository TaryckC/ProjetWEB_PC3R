package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projetweb/backend/backend/database"
	"projetweb/backend/backend/handlers"

	"github.com/gorilla/mux"
)

var FirebaseService *database.FirebaseService

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func setUpLeetCodeAPIRoute(r *mux.Router) {
	r.HandleFunc("/", handlers.HandleRoot).Methods("GET")
	r.HandleFunc("/classic-challenges", database.GetAllClassicChallenges).Methods("GET")
	r.HandleFunc("/classic-challenges/{id}", database.GetClassicChallenge).Methods("GET")
	r.HandleFunc("/daily-challenge", database.GetTodayChallenge).Methods("GET")
	r.HandleFunc("/challengeContent/{titleSlug}", database.GetChallengeContent).Methods("GET")
	r.HandleFunc("/challengeContent/{titleSlug}", database.FetchAndStoreChallengeContent).Methods("POST", "OPTIONS")
}

func setUpCompilerRoutes(r *mux.Router) {
	r.HandleFunc("/compile", handlers.HandleCompiler).Methods("POST", "OPTIONS")
}

func setUpNewsRoutes(r *mux.Router) {
	r.HandleFunc("/news", handlers.HandleNews).Methods("GET")
}
func setUpForum(r *mux.Router) {
	r.HandleFunc("/forum/challengeContent/{titleSlug}", handlers.GetForumMessages).Methods("GET")
	r.HandleFunc("/forum/challengeContent/{titleSlug}", handlers.PostForumMessage).Methods("POST", "OPTIONS")
}

func main() {
	var err error
	FirebaseService, err = database.InitFireBase()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", err)
	}
	database.GlobalFirebaseService = FirebaseService
	// FirebaseService.WriteDailyChallenge(2025, 5)
	// FirebaseService.UpdateDailyQuestionDescription()
	// FirebaseService.WriteDailyAndWeeklyChallenges(2025, 4)
	// FirebaseService.WriteChallengeComplementaryData()
	r := mux.NewRouter()
	r.Use(middlewareCors)

	setUpLeetCodeAPIRoute(r)
	setUpCompilerRoutes(r)
	setUpNewsRoutes(r)
	setUpForum(r)
	fmt.Println("Server running on http://localhost:8080")
	defer FirebaseService.Client.Close()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
