package news

import (
	"fmt"
	"io"
	"net/http"

	"projetweb/backend/api/utils"
)

type NewsResponse struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Source      string `json:"source"`
		Url         string `json:"url"`
		Image       string `json:"image"`
	} `json:"articles"`
}

func GetNews(topic string) ([]byte, error) {
	url := fmt.Sprintf("https://news-api14.p.rapidapi.com/v2/search/articles?query=%s&language=en", topic)
	utils.LoadEnv()
	apiKey, _ := utils.GetApiKey()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "news-api14.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HandleNews(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "Missing topic parameter", http.StatusBadRequest)
		return
	}

	data, err := GetNews(topic)
	if err != nil {
		http.Error(w, "Erreur API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
