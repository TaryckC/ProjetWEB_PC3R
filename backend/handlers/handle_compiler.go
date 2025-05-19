package handlers

import (
	"encoding/json"
	"net/http"
	compiler "projetweb/backend/backend/api/judge0api"
	"projetweb/backend/backend/api/utils"
)

type TestResult struct {
	Index  int         `json:"index"`
	Output interface{} `json:"output,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func HandleCompiler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
		return
	}

	var reqs []compiler.Submission
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Requête JSON invalide")
		return
	}

	// Étape 1 : Envoi batch
	tokens, err := compiler.BatchExecuteCodes(reqs)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Erreur d'envoi batch : "+err.Error())
		return
	}

	if len(tokens) != len(reqs) {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Nombre de tokens incorrect")
		return
	}

	// Étape 2 : Attente des résultats (en batch)
	results, err := compiler.BatchPollResults(tokens)
	if err != nil {
		utils.WriteJSONError(w, http.StatusGatewayTimeout, "Erreur de récupération des résultats : "+err.Error())
		return
	}

	// Étape 3 : Réponse directe (liste d’ExecutionResult)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
