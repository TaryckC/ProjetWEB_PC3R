package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"projetweb/backend/backend/api/utils"
	"time"
)

func ExecuteCode(code string, languageID int, stdin string) (string, error) {
	apiURL := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=true&wait=false"
	utils.LoadEnv()
	apiKey, err := utils.GetApiKey()
	if apiKey == "" || err != nil {
		return "", fmt.Errorf("clé API Judge0 manquante")
	}

	submission := Submission{
		SourceCode: utils.ToBase64(code),
		LanguageID: languageID,
		Stdin:      utils.ToBase64(stdin),
	}
	jsonData, err := json.Marshal(submission)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Réponse brute de Judge0 :", string(body))

	var submissionResp SubmissionResponse
	json.Unmarshal(body, &submissionResp)

	if submissionResp.Token == "" {
		fmt.Println("Problème : Token non reçu !")
		return "", fmt.Errorf("token non reçu de Judge0")
	}

	return submissionResp.Token, nil
}

func GetExecutionResult(token string) (ExecutionResult, error) {
	apiURL := fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/%s?base64_encoded=true&fields=*", token)

	utils.LoadEnv()
	apiKey, err := utils.GetApiKey()
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("erreur récuperation de clé : %w", err)
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("erreur création requête : %w", err)
	}
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ExecutionResult{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result ExecutionResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ExecutionResult{}, err
	}

	result.Stdout, _ = utils.FromBase64(result.Stdout)
	result.Stderr, _ = utils.FromBase64(result.Stderr)

	return result, nil
}

// func HandleCompiler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req Submission
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		http.Error(w, "JSON invalide", http.StatusBadRequest)
// 		return
// 	}

// 	sourceCodeEncoded := utils.ToBase64(req.SourceCode)
// 	stdin := utils.ToBase64(req.Stdin)

// 	token, err := ExecuteCode(sourceCodeEncoded, req.LanguageID, stdin)
// 	if err != nil {
// 		http.Error(w, "Erreur lors de la soumission : "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Attendre que l’exécution soit prête (car wait=false)
// 	time.Sleep(10 * time.Second)

// 	// Récupérer le résultat
// 	result, err := GetExecutionResult(token)
// 	if err != nil {
// 		http.Error(w, "Erreur lors de la récupération : "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Répondre au frontend
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(result)
// }

// pollForResult interroge l'API Judge0 toutes les secondes jusqu'à obtenir un résultat final.
// Elle retourne une erreur si le résultat est toujours en attente après le timeout.
func PollForResult(token string, maxAttempts int, delay time.Duration) (ExecutionResult, error) {
	var result ExecutionResult
	var err error

	for i := 0; i < maxAttempts; i++ {
		time.Sleep(delay)

		result, err = GetExecutionResult(token)
		if err != nil {
			continue
		}

		status := result.Status.Description
		if status != "In Queue" && status != "Processing" {
			return result, nil
		}
	}

	return result, fmt.Errorf("temps dépassé : l'exécution est toujours en attente après %d tentatives", maxAttempts)
}
