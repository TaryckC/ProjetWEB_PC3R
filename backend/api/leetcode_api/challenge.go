package leetcodeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TODO - Il faudra rendre le client global pour pas avoir à en recréer un pour chaque requête. il sera à fermer à la fin du main.

const leetcodeGraphQLEndpoint = "https://leetcode.com/graphql"

type questionResponse struct {
	Data struct {
		Question struct {
			Content string `json:"content"`
		} `json:"question"`
	} `json:"data"`
}

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables"`
}

func DoGraphQLRequest(res GraphQLRequest) ([]byte, error) {
	body, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("failed to encode GraphQL request: %v", err)
	}

	req, err := http.NewRequest("POST", leetcodeGraphQLEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LeetCodeBot/1.0)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("leetcode API error: %s", string(respBody))
	}

	return respBody, nil
}

func RequestChallengeDescription(slug string) (string, error) {
	query := `
		query getDescription($titleSlug: String!) {
			question(titleSlug: $titleSlug) {
				content
			}
		}
	`

	request := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"titleSlug": slug,
		},
	}

	respBody, err := DoGraphQLRequest(request)
	if err != nil {
		return "", err
	}

	var parsed questionResponse
	err = json.Unmarshal(respBody, &parsed)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if parsed.Data.Question.Content == "" {
		log.Println("LEETCODEAPI : description is empty (question not found or not accessible)")
		return "", fmt.Errorf("no description found for slug %q", slug)
	}

	return parsed.Data.Question.Content, nil
}

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

	request := GraphQLRequest{
		Query:         query,
		OperationName: "questionOfToday",
		Variables:     map[string]interface{}{},
	}

	respBody, err := DoGraphQLRequest(request)
	if err != nil {
		log.Printf("LEETCODEAPI : error when requesting daily question : %v\n", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Printf("LEETCODEAPI : raw body: %s\n", respBody)
		log.Printf("LEETCODEAPI : error decoding JSON response : %v\n", err)
		return nil, err
	}
	return result, nil
}

func RequestChallengeList(year int, month int) (map[string]interface{}, error) {
	query := `
	query dailyCodingQuestionRecords($year: Int!, $month: Int!) {
		dailyCodingChallengeV2(year: $year, month: $month) {
			challenges {
				date
				userStatus
				link
				question {
					questionFrontendId
					title
					titleSlug
				}
			}
			weeklyChallenges {
				date
				userStatus
				link
				question {
					questionFrontendId
					title
					titleSlug
					isPaidOnly
				}
			}
		}
	}`

	request := GraphQLRequest{
		Query:         query,
		OperationName: "dailyCodingQuestionRecords",
		Variables: map[string]interface{}{
			"year":  year,
			"month": month,
		},
	}

	respBody, err := DoGraphQLRequest(request)
	if err != nil {
		log.Printf("LEETCODEAPI : error when requesting challenge list : %v\n", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Printf("LEETCODEAPI : raw body: %s\n", respBody)
		log.Printf("LEETCODEAPI : error decoding JSON response : %v\n", err)
		return nil, err
	}

	return result, nil
}

func RequestQuestionsData(titleSlug string) (map[string]interface{}, error) {
	query := `
	query questionData($titleSlug: String!) {
		question(titleSlug: $titleSlug) {
			questionId
			title
			titleSlug
			content
			difficulty
			codeSnippets {
				lang
				langSlug
				code
			}
		}
	}`

	request := GraphQLRequest{
		Query:         query,
		OperationName: "questionData",
		Variables:     map[string]interface{}{"titleSlug": titleSlug},
	}

	respBody, err := DoGraphQLRequest(request)
	if err != nil {
		log.Printf("LEETCODE : error when requesting for challenge's questions data : %v\n", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Printf("LEETCODEAPI : raw body: %s\n", respBody)
		log.Printf("LEETCODEAPI : error decoding JSON response : %v\n", err)
		return nil, err
	}

	return result, nil
}