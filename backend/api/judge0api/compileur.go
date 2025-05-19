package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"projetweb/backend/backend/api/utils"
	"strings"
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

func BatchExecuteCodes(subs []Submission) ([]string, error) {
	apiURL := "https://judge0-ce.p.rapidapi.com/submissions/batch?base64_encoded=true"
	utils.LoadEnv()
	apiKey, err := utils.GetApiKey()
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("clé API Judge0 manquante")
	}

	// Encoder tous les inputs en base64
	for i := range subs {
		subs[i].SourceCode = utils.ToBase64(subs[i].SourceCode)
		subs[i].Stdin = utils.ToBase64(subs[i].Stdin)
	}

	bodyData, err := json.Marshal(BatchSubmission{Submissions: subs})
	if err != nil {
		return nil, fmt.Errorf("erreur encodage JSON : %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Batch POST response:", string(body))

	var batchResp []BatchResponseItem
	err = json.Unmarshal(body, &batchResp)
	if err != nil {
		return nil, err
	}

	if len(batchResp) == 0 {
		return nil, fmt.Errorf("aucun token reçu")
	}

	tokens := make([]string, len(batchResp))
	for i, item := range batchResp {
		tokens[i] = item.Token
	}

	return tokens, nil

}

func BatchPollResults(tokens []string) ([]ExecutionResult, error) {
	apiURL := fmt.Sprintf(
		"https://judge0-ce.p.rapidapi.com/submissions/batch?tokens=%s&base64_encoded=true&fields=*",
		strings.Join(tokens, ","))

	utils.LoadEnv()
	apiKey, err := utils.GetApiKey()
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("clé API manquante")
	}

	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-RapidAPI-Key", apiKey)
		req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var batchResp BatchResultResponse
		err = json.Unmarshal(body, &batchResp)
		if err != nil {
			return nil, err
		}

		results := batchResp.Submissions
		done := true
		for _, res := range results {
			if res.Status.Description == "In Queue" || res.Status.Description == "Processing" {
				done = false
				break
			}
		}
		if done {
			// décoder stdout/stderr
			for i := range results {
				results[i].Stdout, _ = utils.FromBase64(results[i].Stdout)
				results[i].Stderr, _ = utils.FromBase64(results[i].Stderr)
			}
			return results, nil
		}
	}
	return nil, fmt.Errorf("timeout : résultats toujours en attente")
}
