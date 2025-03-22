package database

import (
	"context"
	"log"
	leetcodeapi "projetweb/api/leetcode_api"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App
var FirestoreClient *firestore.Client

func InitFireBase() {

	opt := option.WithCredentialsFile("keys/serviceAccountKey.json")
	config := &firebase.Config{ProjectID: "projetwebpc3r"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	FirebaseApp = app
	log.Println("Connextion à Firebase établie")

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}
	FirestoreClient = client
	log.Println("Client Firestore crée.")
}

func WriteDailyChallenge(year int, month int) {
	challenge, err := leetcodeapi.RequestDailyChallenge(year, month)
	if err != nil {
		log.Printf("LEETCODEAPI : Error fetching daily challenge : %v\n", err)
	}
	result, err := FirestoreClient.Collection("Challenges").Doc("daily_challenge").Set(context.Background(), challenge)
	if err != nil {
		log.Printf("FIREBASE : Error writing daily challenge : %v\n", err)
	}
	log.Println(result)
}
