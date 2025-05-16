package news

type NewsResponse struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Source      string `json:"source"`
		Url         string `json:"url"`
		Image       string `json:"image"`
	} `json:"articles"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
