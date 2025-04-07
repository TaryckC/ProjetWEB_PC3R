package leetcodeapi

type TopicTag struct {
	ID   string `firestore:"id"`
	Name string `firestore:"name"`
	Slug string `firestore:"slug"`
}

type Question struct {
	ACRate           float64    `firestore:"acRate"`
	Difficulty       string     `firestore:"difficulty"`
	FreqBar          *float64   `firestore:"freqBar"`
	FrontendID       string     `firestore:"frontendQuestionId"`
	HasSolution      bool       `firestore:"hasSolution"`
	HasVideoSolution bool       `firestore:"hasVideoSolution"`
	IsFavor          bool       `firestore:"isFavor"`
	PaidOnly         bool       `firestore:"paidOnly"`
	Status           *string    `firestore:"status"`
	Title            string     `firestore:"title"`
	Description      *string    `firestore:"description"` // Initialement nul, ajouté après une première demande d'accès à la question.
	TitleSlug        string     `firestore:"titleSlug"`
	TranslatedTitle  *string    `firestore:"translatedTitle"`
	UserStatus       string     `firestore:"userStatus"`
	TopicTags        []TopicTag `firestore:"topicTags"`
}

type ActiveDailyCodingChallenge struct {
	Date     string   `firestore:"date"`
	Link     string   `firestore:"link"`
	Question Question `firestore:"question"`
}

// Challenge quelconque
type MinimalQuestion struct {
	FrontendID string `firestore:"questionFrontendId" json:"questionFrontendId"`
	Title      string `firestore:"title" json:"title"`
	TitleSlug  string `firestore:"titleSlug" json:"titleSlug"`
	UserStatus string `firestore:"userStatus" json:"userStatus"`
}

type ChallengeItem struct {
	Date     string          `firestore:"date" json:"date"`
	Link     string          `firestore:"link" json:"link"`
	Question MinimalQuestion `firestore:"question" json:"question"`
}
