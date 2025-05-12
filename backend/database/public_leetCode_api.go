package database

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	leetcodeapi "projetweb/backend/api/leetcode_api"

	"github.com/gorilla/mux"
)

// API REST PERMETTANT LES OPÉRATION DE BASES SUR LES STRUCTURES DE DONNÉES ASSOCIÉES À L'API LEETCODE

/**
* CLASSIC CHALLENGES
**/

// GET /classic-challenges
func GetAllClassicChallenges(w http.ResponseWriter, r *http.Request) {
	docs, err := GlobalFirebaseService.Client.Collection(ClassicChallengesDoc).Documents(context.Background()).GetAll()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des challenges", http.StatusInternalServerError)
		return
	}

	var challenges []leetcodeapi.ChallengeItem
	for _, doc := range docs {
		var c leetcodeapi.ChallengeItem
		if err := doc.DataTo(&c); err == nil {
			challenges = append(challenges, c)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

// GET /classic-challenges/{id}
func GetClassicChallenge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("DEBUG: Tentative de récupération du challenge avec ID = %s", id)

	challenge, _, err := GlobalFirebaseService.GetChallengeFromDataBase(ClassicChallengesDoc, id)
	if err != nil {
		log.Printf("ERREUR: Échec de la récupération du challenge ID = %s : %v", id, err)
		http.Error(w, "Erreur d'accès à la base de données", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}

//############################################################################################################################################
//############################################################################################################################################

/**
* DAILY CHALLENGE
**/

// GET /daily-challenge
func GetTodayChallenge(w http.ResponseWriter, r *http.Request) {
	challenge, _, err := GlobalFirebaseService.getDailyChallengeFromDataBase()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du challenge du jour", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}
