package database

import (
	"context"
	"fmt"
	"log"
	leetcodeapi "projetweb/backend/backend/api/leetcode_api"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DailyChallengeDoc  = "daily_challenge"
	Daily_challenge_id = "0"

	WeeklyChallengesDoc  = "weekly_challenges"
	ClassicChallengesDoc = "classic_challenges"
	FirestoreCollection  = "Challenges"
	ChallengeContentDoc  = "Challenge_content"
	FirebaseProjectID    = "pc3rprojet-ce4a7"
	FirebaseKeyPath      = "backend/keys/serviceAccountKey.json"
)

var GlobalFirebaseService *FirebaseService

type FirebaseService struct {
	App    *firebase.App
	Client *firestore.Client
}

func InitFireBase() (*FirebaseService, error) {
	opt := option.WithCredentialsFile(FirebaseKeyPath)
	log.Printf("🔍 Chargement de la clé Firebase depuis : %s", FirebaseKeyPath)
	config := &firebase.Config{ProjectID: FirebaseProjectID}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp failed: %w", err)
	}
	log.Println("✅ Clé Firebase acceptée, application initialisée.")
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
/* GENERIC FUNCTIONS FOR CHALLENGES
/**/

// findChallengeContentBySlug récupère le contenu d'un challenge à partir du titleSlug
func findChallengeContentBySlug(titleSlug string) (map[string]interface{}, error) {
	log.Printf("🔍 Recherche du challenge avec titleSlug = %s", titleSlug)
	doc, err := GlobalFirebaseService.Client.Collection(ChallengeContentDoc).Doc(titleSlug).Get(context.Background())
	if err != nil {
		log.Printf("❌ Erreur dans findChallengeContentBySlug (Firestore): %v", err)
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}
	var content map[string]interface{}
	if err := doc.DataTo(&content); err != nil {
		log.Printf("❌ Erreur dans findChallengeContentBySlug (Firestore): %v", err)
		return nil, err
	}
	log.Printf("✅ Contenu du challenge trouvé : %+v", content)
	return content, nil
}

/**/
/*  Daily Challenge gestion
/**/

func (fs *FirebaseService) WriteDailyChallenge(year int, month int) error {
	log.Println("🟡 Début de WriteDailyChallenge")
	challenge, err := leetcodeapi.RequestDailyChallenge(year, month)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Error fetching daily challenge : %v", err)
	}

	log.Println("📥 Daily challenge récupéré depuis l'API")

	challengeData := challenge["data"].(map[string]interface{})
	activeChallenge := challengeData["activeDailyCodingChallengeQuestion"].(map[string]interface{})
	log.Printf("🔍 Challenge actif extrait : %+v", activeChallenge)
	if err := fs.writeChallenge(DailyChallengeDoc, Daily_challenge_id, activeChallenge); err != nil {
		log.Printf("❌ Échec d’écriture du daily challenge : %v", err)
		return fmt.Errorf("FIREBASE : Error writing daily challenge : %w", err)
	}

	log.Println("📝 Daily challenge écrit dans Firestore")

	fs.WriteDailyChallengeComplementaryData()

	log.Println("🧩 Données complémentaires du daily challenge ajoutées")

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

	log.Println("✅ SUCCÈS : description du challenge écrite dans Firestore")

	return nil
}

/**/
/*  Daily and Weekly challenges gestion
/**/

// TODO : Utiliser write Challenge dans dailyChallenge aussi pour la factoriser
func (fs *FirebaseService) writeChallenge(collection string, doc string, data map[string]interface{}) error {
	_, err := fs.Client.Collection(collection).Doc(doc).Set(context.Background(), data)
	if err != nil {
		return fmt.Errorf("FIREBASE : Error writing challenge : %v", err)
	}
	return nil
}

func (fs *FirebaseService) writeChallengeContent(titleSlug string) error {
	content, err := leetcodeapi.RequestQuestionsData(titleSlug)
	if err != nil {
		return fmt.Errorf("LEETCODEAPI : Error fetching challenge content : %v", err)
	}
	_, err = fs.Client.Collection(ChallengeContentDoc).Doc(titleSlug).Set(context.Background(), content)
	if err != nil {
		log.Println("FIREBASE : Error while writing challenge content in the database")
		return err
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

type ForumPost struct {
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (fs *FirebaseService) PostForumMessage(challengeId string, post ForumPost) error {
	if challengeId == "" {
		log.Println("Erreur : challengeId vide dans PostForumMessage")
		return nil
	}
	_, _, err := fs.Client.
		Collection("Challenge_content").
		Doc(challengeId).
		Collection("forum").
		Add(context.Background(), post)

	if err != nil {
		log.Println("Error adding forum post:", err)
	}

	return err
}

func (fs *FirebaseService) GetForumMessage(challengeId string) ([]ForumPost, error) {
	if challengeId == "" {
		log.Println("Erreur : challengeId vide dans PostForumMessage")
		return nil, nil
	}
	ctx := context.Background()
	docs, err := fs.Client.
		Collection("Challenge_content").
		Doc(challengeId).
		Collection("forum").
		OrderBy("CreatedAt", firestore.Asc).
		Documents(ctx).
		GetAll()

	if err != nil {
		log.Println("Error fetching forum posts:", err)
		return nil, err
	}

	var posts []ForumPost
	for _, doc := range docs {
		var p ForumPost
		if err := doc.DataTo(&p); err == nil {
			posts = append(posts, p)
		}
	}

	return posts, nil
}
