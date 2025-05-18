package news

import (
	"fmt"
	"io"
	"net/http"
	"projetweb/backend/backend/api/utils"
)

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
