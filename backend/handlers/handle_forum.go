package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projetweb/backend/backend/api/utils"
	"projetweb/backend/backend/database"
	"time"

	"github.com/gorilla/mux"
)

func PostForumMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["titleSlug"]

	log.Println("🔵 POST forum reçu pour titleSlug =", problemID)

	var post database.ForumPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Println("🔴 Erreur décodage JSON :", err)
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("🟡 Message reçu : author=%s | content=%s\n", post.Author, post.Content)

	post.CreatedAt = time.Now()

	err := database.GlobalFirebaseService.PostForumMessage(problemID, post)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Error saving post")
		log.Println("🔴 Erreur Firebase :", err)
		return
	}

	log.Println("🟢 Message enregistré avec succès")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Post added successfully")
}

func GetForumMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["titleSlug"]

	posts, err := database.GlobalFirebaseService.GetForumMessage(problemID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Error fetching forum posts")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
