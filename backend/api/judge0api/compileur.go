package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func ExecuteCode(code string, languageID int, stdin string) (string, error) {
	apiURL := "https://judge0-ce.p.rapidapi.com/submissions?base64_encoded=false&wait=true"

	apiKey := os.Getenv("JUDGE0_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Clé API Judge0 manquante")
	}

	// Construire la requête JSON
	submission := Submission{
		SourceCode: code,
		LanguageID: languageID,
		Stdin:      stdin,
	}
	jsonData, err := json.Marshal(submission)
	if err != nil {
		return "", err
	}

	// Afficher le JSON pour vérifier
	fmt.Println("JSON envoyé :", string(jsonData))

	// Faire la requête POST vers Judge0
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

	// Lire la réponse
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Réponse brute de Judge0 :", string(body)) // Debugging

	// Décoder la réponse JSON
	var submissionResp SubmissionResponse
	json.Unmarshal(body, &submissionResp)

	if submissionResp.Token == "" {
		fmt.Println("❌ Problème : Token non reçu !")
		return "", fmt.Errorf("Token non reçu de Judge0")
	}

	return submissionResp.Token, nil
}

func GetExecutionResult(token string) {
	apiURL := fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/%s?base64_encoded=true&fields=*", token)

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("X-RapidAPI-Key", "503c87796dmshb4daef3cd7bc808p10c392jsnb4c28ac49ba1")
	req.Header.Set("X-RapidAPI-Host", "judge029.p.rapidapi.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("❌ Erreur lors de la requête GET :", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("📥 Résultat de l'exécution :", string(body))
}

func man() {
	url := "https://judge029.p.rapidapi.com/submissions?base64_encoded=true&wait=false&fields=*"

	payload := strings.NewReader(`{
		"source_code": "",
		"language_id": 71,
		"stdin": "SnVkZ2Uw",
		
	}`)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("x-rapidapi-key", "503c87796dmshb4daef3cd7bc808p10c392jsnb4c28ac49ba1")
	req.Header.Add("x-rapidapi-host", "judge029.p.rapidapi.com")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//fmt.Println(res)
	fmt.Println(string(body))

	var submissionResp SubmissionResponse
	err := json.Unmarshal(body, &submissionResp)
	if err != nil {
		return
	}

}
