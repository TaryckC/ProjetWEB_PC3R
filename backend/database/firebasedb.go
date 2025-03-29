package database

import (
	"context"
	"fmt"
	"log"
	leetcodeapi "projetweb/api/leetcode_api"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	DailyChallengeDoc   = "daily_challenge"
	FirestoreCollection = "Challenges"
	FirebaseProjectID   = "pc3rprojet-ce4a7"
	FirebaseKeyPath     = "keys/serviceAccountKey.json"
)

type FirebaseService struct {
	App    *firebase.App
	Client *firestore.Client
}

func InitFireBase() (*FirebaseService, error) {
	opt := option.WithCredentialsFile(FirebaseKeyPath)
	config := &firebase.Config{ProjectID: FirebaseProjectID}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp failed: %w", err)
	}
	log.Println("Connexion à Firebase établie")

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("firebase.Firestore failed: %w", err)
	}
	log.Println("Client Firestore créé.")

	return &FirebaseService{
		App:    app,
		Client: client,
	}, nil
}

func (fs *FirebaseService) WriteDailyChallenge(year int, month int) error {
	challenge, err := leetcodeapi.RequestDailyChallenge(year, month)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Error fetching daily challenge : %v", err)
	}
	_, err = fs.Client.Collection(FirestoreCollection).Doc(DailyChallengeDoc).Set(context.Background(), challenge)
	if err != nil {
		return fmt.Errorf("FIREBASE : Error writing daily challenge : %v", err)
	}
	return nil
}

// TODO : Peut-être plus tard stocker la référence dans une variable ?
// TODO : Afficher l'erreur et laisser le traitement aux autres ?
func (fs *FirebaseService) getDailyChallengeFromDataBase() (*leetcodeapi.DailyChallenge, *firestore.DocumentSnapshot, error) {
	doc, err := fs.Client.Collection(FirestoreCollection).Doc(DailyChallengeDoc).Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("FIRESTORE : failed to read daily challenge: %v", err)
	}

	challenge := new(leetcodeapi.DailyChallenge)
	if err := doc.DataTo(challenge); err != nil {
		return nil, doc, fmt.Errorf("FIRESTORE : failed to decode daily challenge: %v", err)
	}

	return challenge, doc, nil
}

func (fs *FirebaseService) WriteDailyChallengeComplementaryData() error {
	err := fs.UpdateDailyQuestionDescription()
	if err != nil {
		return fmt.Errorf("FIRESTORE : failed to add daily challenge description to challenge : %v", err)
	}
	return nil
}

func (fs *FirebaseService) UpdateDailyQuestionDescription() error {
	challenge, doc, err := fs.getDailyChallengeFromDataBase()
	if err != nil {
		return fmt.Errorf("FIRESTORE : failed to fetch daily challenge: %v", err)
	}

	titleSlug := challenge.Data.ActiveDailyCodingChallengeQuestion.Question.TitleSlug

	description, err := leetcodeapi.RequestChallengeDescription(titleSlug)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Failed to get Question description: %v", err)
	}

	_, err = doc.Ref.Update(context.Background(), []firestore.Update{
		{
			Path:  "activeDailyCodingChallengeQuestion.question.description",
			Value: description,
		},
	})
	if err != nil {
		return fmt.Errorf("FIRESTORE : Failed to update question description: %v", err)
	}

	return nil
}
