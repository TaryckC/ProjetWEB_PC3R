package compiler

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

type Example struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type EvaluationRequest struct {
	SourceCode string    `json:"source_code"`
	LanguageID int       `json:"language_id"`
	Stdin      string    `json:"stdin"`
	Examples   []Example `json:"examples"`
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
