package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projetweb/backend/api/utils"
	"projetweb/backend/database"
	"time"

	"github.com/gorilla/mux"
)

func PostForumMessage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    problemID := vars["id"]

    var post database.ForumPost
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    post.CreatedAt = time.Now()

    err := database.GlobalFirebaseService.PostForumMessage(problemID,post)

    if err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError,  "Error saving post")
        log.Println("Error adding forum post:", err)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "Post added successfully")
}



func GetForumMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["id"]

	posts, err := database.GlobalFirebaseService.GetForumMessage(problemID)
	if err != nil {
		utils.WriteJSONError(w,  http.StatusInternalServerError, "Error fetching forum posts")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
