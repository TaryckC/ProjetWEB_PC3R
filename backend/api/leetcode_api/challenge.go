package leetcodeapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Il faudra rendre le client global pour pas avoir à en recréer un pour chaque requête. il sera à fermer à la fin du main.

func RequestDailyChallenge(year int, month int) (map[string]interface{}, error) {
	query := `
	query questionOfToday {
		activeDailyCodingChallengeQuestion {
			date
			userStatus
			link
			question {
				titleSlug
				title
				translatedTitle
				acRate
				difficulty
				freqBar
				frontendQuestionId: questionFrontendId
				isFavor
				paidOnly: isPaidOnly
				status
				hasVideoSolution
				hasSolution
				topicTags {
					name
					id
					slug
				}
			}
		}
	}`

	payload := map[string]interface{}{
		"operationName": "questionOfToday",
		"query":         query,
		"variables":     map[string]interface{}{},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("LEETCODEAPI : error when encoding daily question payload to JSON : %v\n", err)
		return nil, err
	}
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("LEETCODEAPI : error when requesting daily question : %v\n", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", "https://leetcode.com")
	request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LeetCodeBot/1.0)")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("LEETCODEAPI : error sending HTTP request to API  : %v\n", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("LEETCODEAPI : error reading response body : %v\n", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("LEETCODEAPI : raw body: %s\n", body)
		log.Printf("LEETCODEAPI : error decoding JSON response : %v\n", err)
		return nil, err
	}
	return result, nil
}
