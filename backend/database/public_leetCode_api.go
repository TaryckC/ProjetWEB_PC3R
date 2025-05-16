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
* GENERIC OPERATIONS
**/

// GetChallengeContent essaye de récupérer le contenu du challenge.
func GetChallengeContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	titleSlug := vars["titleSlug"]
	content, err := findChallengeContentBySlug(titleSlug)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des contenus", http.StatusInternalServerError)
		return
	}
	if content == nil {
		http.Error(w, "Contenu non trouvé pour ce titleSlug", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

// FetchAndStoreChallengeContent : Se comporte comme un GET sauf si la donnée recherchée n'existe pas, alors on l'ajoute dans la BDD.
func FetchAndStoreChallengeContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	titleSlug := vars["titleSlug"]

	content, err := findChallengeContentBySlug(titleSlug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if content != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(content)
		return
	}

	// Si le contenu n'existe pas encore, on le récupère via l'API LeetCode et on l'écrit dans Firestore
	err = GlobalFirebaseService.writeChallengeContent(titleSlug)
	if err != nil {
		log.Printf("Erreur lors de l'écriture du challenge %s : %v", titleSlug, err)
		http.Error(w, "Erreur lors de l'enregistrement du challenge", http.StatusInternalServerError)
		return
	}

	// On récupère à nouveau le contenu maintenant qu'il est censé être écrit
	content, err = findChallengeContentBySlug(titleSlug)
	if err != nil || content == nil {
		http.Error(w, "Erreur lors de la récupération après écriture", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

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