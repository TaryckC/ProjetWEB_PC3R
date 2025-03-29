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

type DailyChallenge struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion ActiveDailyCodingChallenge `firestore:"activeDailyCodingChallengeQuestion"`
	} `firestore:"data"`
}
