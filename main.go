package main

import (
	"fmt"
	"log"
	"net/http"
	"projetweb/backend/backend/database"
	"projetweb/backend/backend/handlers"
	"time"

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
	FirebaseService.WriteDailyChallenge(2025, 5)
	// FirebaseService.UpdateDailyQuestionDescription()
	// FirebaseService.WriteChallengeComplementaryData()
	r := mux.NewRouter()
	r.Use(middlewareCors)

	setUpLeetCodeAPIRoute(r)
	setUpCompilerRoutes(r)
	setUpNewsRoutes(r)
	setUpForum(r)
	go scheduleDailyChallenge()

	fmt.Println("Server running on http://localhost:8080")
	defer FirebaseService.Client.Close()
	log.Fatal(http.ListenAndServe("[::]:8100", r))
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

func scheduleDailyChallenge() {

	leetcodeTZ := "America/Los_Angeles" // PST/PDT (LeetCode's timezone)
	loc, err := time.LoadLocation(leetcodeTZ)
	if err != nil {
		loc = time.UTC
		log.Printf("⚠️ Failed to load LeetCode timezone, defaulting to UTC")
	}

	const (
		targetHour   = 0                // Minuit (PST/PDT)
		targetMinute = 5                // 5 min après minuit pour laisser un buffer
		retryDelay   = 30 * time.Minute // En cas d'échec
	)

	var lastExecDate string

	// Fonction pour exécuter et mettre à jour la date
	fetchAndUpdate := func() error {
		now := time.Now().In(loc)
		today := now.Format("2006-01-02")
		if today == lastExecDate {
			return nil // Déjà fait aujourd'hui
		}
		err := FirebaseService.WriteDailyChallenge(now.Year(), int(now.Month()))
		if err != nil {
			return err
		}
		lastExecDate = today
		log.Println("✅ Successfully updated Daily Challenge")
		return nil
	}

	// Premier fetch immédiat si nécessaire
	if err := fetchAndUpdate(); err != nil {
		log.Printf("⚠️ Initial fetch failed: %v. Retrying in %v", err, retryDelay)
		time.Sleep(retryDelay)
		// Nouvelle tentative (peut être étendu en backoff exponentiel)
	}

	// Boucle principale pour les prochains jours
	for {
		now := time.Now().In(loc)
		next := time.Date(
			now.Year(), now.Month(), now.Day(),
			targetHour, targetMinute, 0, 0, loc,
		)

		// Si l'heure est passée, programmer pour demain
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		waitDuration := next.Sub(now)
		log.Printf("⏳ Next fetch at %s (in %v)", next.Format("2006-01-02 15:04 MST"), waitDuration)

		time.Sleep(waitDuration)

		// Exécution normale
		if err := fetchAndUpdate(); err != nil {
			log.Printf("⚠️ Fetch failed: %v. Retrying in %v", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}
	}

}
