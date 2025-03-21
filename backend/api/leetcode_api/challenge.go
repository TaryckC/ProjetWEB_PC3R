package leetcodeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TopicTag struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

type Question struct {
	AcRate             float64    `json:"acRate"`
	Difficulty         string     `json:"difficulty"`
	FreqBar            *float64   `json:"freqBar"`
	FrontendQuestionID string     `json:"frontendQuestionId"`
	IsFavor            bool       `json:"isFavor"`
	PaidOnly           bool       `json:"paidOnly"`
	Status             *string    `json:"status"`
	Title              string     `json:"title"`
	TitleSlug          string     `json:"titleSlug"`
	HasVideoSolution   bool       `json:"hasVideoSolution"`
	HasSolution        bool       `json:"hasSolution"`
	TopicTags          []TopicTag `json:"topicTags"`
}

type ActiveDailyCodingChallengeQuestion struct {
	Date       string   `json:"date"`
	UserStatus string   `json:"userStatus"`
	Link       string   `json:"link"`
	Question   Question `json:"question"`
}

type APIResponse struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion ActiveDailyCodingChallengeQuestion `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

// https://medium.com/hprog99/working-with-json-in-golang-a-comprehensive-guide-5a94ca5961a1

func RequestDailyChallenge(apiKey string) (*ActiveDailyCodingChallengeQuestion, error) {
	url := "https://leetcode-api.p.rapidapi.com/leetcode/todays-question"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("server: could not create request: %v", err)
	}

	// Header pour l'API
	req.Header.Set("x-rapidapi-host", "leetcode-api.p.rapidapi.com") // Risque de bloquage si non réponse de rapid API ?
	req.Header.Set("x-rapidapi-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("server: request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server: API returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("server: could not read response: %v", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("server: could not parse JSON: %v", err)
	}

	return &apiResponse.Data.ActiveDailyCodingChallengeQuestion, nil
}

// https://www.digitalocean.com/community/tutorials/how-to-make-http-requests-in-go
// https://pkg.go.dev/net/http#hdr-Clients_and_Transports
