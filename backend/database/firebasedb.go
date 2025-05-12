package database

import (
	"context"
	"fmt"
	"log"
	leetcodeapi "projetweb/backend/api/leetcode_api"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	DailyChallengeDoc  = "daily_challenge"
	Daily_challenge_id = "0"

	WeeklyChallengesDoc  = "weekly_challenges"
	ClassicChallengesDoc = "classic_challenges"
	FirestoreCollection  = "Challenges"
	FirebaseProjectID    = "pc3rprojet-ce4a7"
	FirebaseKeyPath      = "keys/serviceAccountKey.json"
)

var GlobalFirebaseService *FirebaseService

// TODO : Va falloir réorganiser la BDD :
/*
*  Challenges -> TYPES (3+?) -> challenges -> data
 */

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

	res := &FirebaseService{App: app, Client: client}
	GlobalFirebaseService = res

	return res, nil
}

// TODO : Peut-être plus tard stocker la référence dans une variable ?
// TODO : Afficher l'erreur et laisser le traitement aux autres ?

/**/
/*  Daily Challenge gestion
/**/

func (fs *FirebaseService) WriteDailyChallenge(year int, month int) error {
	challenge, err := leetcodeapi.RequestDailyChallenge(year, month)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Error fetching daily challenge : %v", err)
	}

	challengeData := challenge["data"].(map[string]interface{})
	activeChallenge := challengeData["activeDailyCodingChallengeQuestion"].(map[string]interface{})
	if fs.writeChallenge(DailyChallengeDoc, Daily_challenge_id, activeChallenge) != nil {
		return fmt.Errorf("FIREBASE : Error writing daily challenge : %v", err)
	}
	return nil
}

func (fs *FirebaseService) getDailyChallengeFromDataBase() (*leetcodeapi.ActiveDailyCodingChallenge, *firestore.DocumentSnapshot, error) {
	doc, err := fs.Client.Collection(DailyChallengeDoc).Doc(Daily_challenge_id).Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("FIRESTORE : failed to read daily challenge: %v", err)
	}

	challenge := new(leetcodeapi.ActiveDailyCodingChallenge)
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

	titleSlug := challenge.Question.TitleSlug

	description, err := leetcodeapi.RequestChallengeDescription(titleSlug)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Failed to get Question description: %v", err)
	}

	_, err = doc.Ref.Update(context.Background(), []firestore.Update{
		{
			Path:  "question.description",
			Value: description,
		},
	})
	if err != nil {
		return fmt.Errorf("FIRESTORE : Failed to update question description: %v", err)
	}

	return nil
}

/**/
/*  Daily and Weekly challenges gestion
/**/

// TODO : Utiliser write Challenge dans dailyChallenge aussi pour la factoriser
// TODO : Il va falloir modifier la façon de stocker les données des dailychallenges dans la BDD (cf firebase pour voir le problème)
func (fs *FirebaseService) writeChallenge(collection string, doc string, data map[string]interface{}) error {
	_, err := fs.Client.Collection(collection).Doc(doc).Set(context.Background(), data)
	if err != nil {
		return fmt.Errorf("FIREBASE : Error writing challenge : %v", err)
	}
	return nil
}

// TODO : Ajouter une fonction de cleanup pour ne pas avoir d'ancienne valeur dans la BDD avant d'en ajouter de nouvelle ?
/*
* Récupère et écrit dans la BDD la liste des challenges disponibles
 */
func (fs *FirebaseService) WriteDailyAndWeeklyChallenges(year int, month int) error {
	challenges, err := leetcodeapi.RequestChallengeList(year, month)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Error fetching daily and weekly challenges : %v", err)
	}
	everyChallenges := challenges["data"].(map[string]interface{})["dailyCodingChallengeV2"].(map[string]interface{})

	// Writing Weekly Challenges
	weeklyChallenges := everyChallenges["weeklyChallenges"].([]interface{})
	for i, c := range weeklyChallenges {
		challenge := c.(map[string]interface{})
		if fs.writeChallenge(WeeklyChallengesDoc, strconv.Itoa(i), challenge) != nil {
			log.Println("FIREBASE : Error while writing a weekly challenge in the database")
		}
	}

	// Writing Classic Challenges
	classicChallenges := everyChallenges["challenges"].([]interface{})
	for i, c := range classicChallenges {
		challenge := c.(map[string]interface{})
		if fs.writeChallenge(ClassicChallengesDoc, strconv.Itoa(i), challenge) != nil {
			log.Println("FIREBASE : Error while writing a weekly challenge in the database")
		}
	}

	return nil
}

func (fs *FirebaseService) GetChallengeFromDataBase(collection string, doc string) (*leetcodeapi.ChallengeItem, *firestore.DocumentSnapshot, error) {
	_doc, err := fs.Client.Collection(collection).Doc(doc).Get(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("FIRESTORE : failed to read daily challenge: %v", err)
	}

	challenge := new(leetcodeapi.ChallengeItem)
	if err := _doc.DataTo(challenge); err != nil {
		return nil, _doc, fmt.Errorf("FIRESTORE : failed to decode daily challenge: %v", err)
	}

	return challenge, _doc, nil
}

// Définir une fonction UpdateAll.. Qui mettra à jours tous les challenges de la BDD
// Ecrire une foction update/ecrire description pour chaque différents type de challenges

func (fs *FirebaseService) UpdateChallengeDescription(collection string, doc string) error {
	challenge, _doc, err := fs.GetChallengeFromDataBase(collection, doc)
	if err != nil {
		return fmt.Errorf("FIRESTORE : failed to fetch daily challenge: %v", err)
	}

	titleSlug := challenge.Question.TitleSlug

	description, err := leetcodeapi.RequestChallengeDescription(titleSlug)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Failed to get Question description: %v", err)
	}

	_, err = _doc.Ref.Update(context.Background(), []firestore.Update{
		{
			Path:  "question.description",
			Value: description,
		},
	})
	if err != nil {
		return fmt.Errorf("FIRESTORE : Failed to update question description: %v", err)
	}

	return nil
}

/**/
/*  Utility functions for challenges manipulation
/**/

func (fs *FirebaseService) WriteChallengeComplementaryData() error {
	collections := []string{ClassicChallengesDoc, WeeklyChallengesDoc}
	for _, collection := range collections {
		docs, err := fs.Client.Collection(collection).Documents(context.Background()).GetAll()
		if err != nil {
			log.Printf("FIRESTORE : failed to get documents from %s: %v", collection, err)
			continue
		}
		for _, doc := range docs {
			err := fs.UpdateChallengeDescription(collection, doc.Ref.ID)
			if err != nil {
				log.Printf("FIRESTORE : failed to update challenge %s/%s: %v", collection, doc.Ref.ID, err) // TODO : REGARDER SI IL Y A ENCORE DES ERREURES AU MOMENT DE RÉCUPER LA DESCRIPTION (06/04 + 1 semaine)
			}
		}
	}
	return nil
}
