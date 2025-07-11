package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func ToBase64(code string) string {
	return base64.StdEncoding.EncodeToString([]byte(code))
}

func FromBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("❌ Erreur lors du décodage Base64 : %v", err)
	}
	return string(decodedBytes), nil
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ ERREUR: .env introuvable")
	} else {
		log.Println("✅ .env chargé")
		log.Println("clé =", os.Getenv("JUDGE0_API_KEY"))
	}
}

func GetApiKey() (string, error) {
	apiKey := os.Getenv("RAPID_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Clé API Judge0 manquante")
	}
	return apiKey, nil
}

// Genérer par chatgpt
func PrintPrettyJSON(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("❌ Erreur lors du formatage JSON :", err)
		return
	}
	fmt.Println(string(prettyJSON))
}
