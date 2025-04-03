package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"projetweb/backend/api/utils"
	"time"
)

// Structure de soumission pour Judge0
type Submission struct {
	SourceCode string `json:"source_code"`
	LanguageID int    `json:"language_id"`
	Stdin      string `json:"stdin"`
}

// Structure de réponse de Judge0
type SubmissionResponse struct {
	Token string `json:"token"`
}

type ExecutionResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Status struct {
		Description string `json:"description"`
	} `json:"status"`
	Time   string `json:"time"`
	Memory int    `json:"memory"`
}

func ExecuteCode(code string, languageID int, stdin string) (string, error) {
	apiURL := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=false&wait=false"
	utils.LoadEnv()
	apiKey, _ := utils.GetApiKey()
	if apiKey == "" {
		return "", fmt.Errorf("Clé API Judge0 manquante")
	}

	submission := Submission{
		SourceCode: code,
		LanguageID: languageID,
		Stdin:      stdin,
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
		fmt.Println("❌ Problème : Token non reçu !")
		return "", fmt.Errorf("Token non reçu de Judge0")
	}

	return submissionResp.Token, nil
}

func GetExecutionResult(token string) (ExecutionResult, error) {
	apiURL := fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/%s?base64_encoded=true&fields=*", token)

	utils.LoadEnv()
	apiKey, _ := utils.GetApiKey()

	req, _ := http.NewRequest("GET", apiURL, nil)
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

	return result, nil
}

func HandleCompiler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req Submission
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	token, err := ExecuteCode(req.SourceCode, req.LanguageID, req.Stdin)
	if err != nil {
		http.Error(w, "Erreur lors de la soumission : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Attendre que l’exécution soit prête (car wait=false)
	time.Sleep(10 * time.Second)

	// Récupérer le résultat
	result, err := GetExecutionResult(token)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre au frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
