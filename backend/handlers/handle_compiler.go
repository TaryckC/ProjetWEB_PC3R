package handlers

import (
	"encoding/json"
	"net/http"
	compiler "projetweb/backend/backend/api/judge0api"
	"projetweb/backend/backend/api/utils"
	"sync"
	"time"
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

	results := make([]TestResult, len(reqs))
	var wg sync.WaitGroup

	// 💡 Limite à 50 goroutines simultanées
	maxWorkers := 50
	semaphore := make(chan struct{}, maxWorkers)

	for i, req := range reqs {
		wg.Add(1)
		semaphore <- struct{}{} // ✋ bloque si trop de goroutines en cours

		go func(i int, req compiler.Submission) {
			defer wg.Done()
			defer func() { <-semaphore }() // ✅ libère un slot

			token, err := compiler.ExecuteCode(req.SourceCode, req.LanguageID, req.Stdin)
			if err != nil {
				results[i] = TestResult{Index: i, Error: "Erreur de soumission : " + err.Error()}
				return
			}

			res, err := compiler.PollForResult(token, 5, 2*time.Second)
			if err != nil {
				results[i] = TestResult{Index: i, Error: "Erreur d'exécution : " + err.Error()}
				return
			}

			results[i] = TestResult{Index: i, Output: res}
		}(i, req)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
