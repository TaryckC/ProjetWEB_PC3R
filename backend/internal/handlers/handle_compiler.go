package handlers

import (
	"encoding/json"
	"net/http"
	compiler "projetweb/backend/api/judge0api"
	"projetweb/backend/api/utils"
	"time"
)

func HandleCompiler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var req compiler.Submission
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Requête JSON invalide")
		return
	}

	// Soumission à Judge0
	token, err := compiler.ExecuteCode(req.SourceCode, req.LanguageID, req.Stdin)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Erreur de soumission : "+err.Error())
		return
	}

	// Polling pour attendre la réponse
	result, err := compiler.PollForResult(token, 5, 2*time.Second)
	if err != nil {
		utils.WriteJSONError(w, http.StatusGatewayTimeout, err.Error())
		return
	}

	// Succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
