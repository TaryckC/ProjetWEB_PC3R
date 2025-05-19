package news

import (
	"fmt"
	"io"
	"net/http"
	"projetweb/backend/backend/api/utils"
	"time"
)



var cache = make(map[string]cacheEntry)
const cacheDuration = 1 * time.Hour



func GetNews(topic string) ([]byte, error) {
	now := time.Now()

	// Vérifie si on a une entrée valide en cache
	if entry, exists := cache[topic]; exists {
		if now.Sub(entry.timestamp) < cacheDuration {
			return entry.data, nil
		}
	}

	// Sinon, requête à l'API
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

	// Mise à jour du cache
	cache[topic] = cacheEntry{
		data:      body,
		timestamp: now,
	}

	return body, nil
}


func RefreshTopicsPeriodically(topics []string) {
	go func() {
		for {
			for _, topic := range topics {
				_, err := GetNews(topic) // force mise à jour si plus vieux que 1h
				if err != nil {
					fmt.Println("Erreur mise à jour news pour topic :", topic, err)
				} else {
					fmt.Println("News mises à jour pour topic :", topic)
				}
			}
			time.Sleep(1 * time.Hour)
		}
	}()
}
