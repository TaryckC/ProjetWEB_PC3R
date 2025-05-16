package handlers

import (
	"net/http"
	news "projetweb/backend/api/news_api"
	"projetweb/backend/api/utils"
)

func HandleNews(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Topic manquant")
		return
	}

	data, err := news.GetNews(topic)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Erreur API: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
