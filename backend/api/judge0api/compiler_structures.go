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

type ExecutionResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Status struct {
		Description string `json:"description"`
	} `json:"status"`
	Time   string `json:"time"`
	Memory int    `json:"memory"`
}

type BatchSubmission struct {
	Submissions []Submission `json:"submissions"`
}

type BatchResponseItem struct {
	Token string `json:"token"`
}
type BatchResultResponse struct {
	Submissions []ExecutionResult `json:"submissions"`
}
