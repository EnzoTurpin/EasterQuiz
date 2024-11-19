package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	allQuestions         []Question
	quiz                 []Question
	scores               map[string]int = make(map[string]int)
	currentQuestionIndex int
	lastFeedback         *FeedbackData
	currentDifficulty    string
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialisation du générateur de nombres aléatoires

	// Génération des questions
	var err error
	quiz, err = generateQuestions("questions.csv", true)
	if err != nil {
		log.Fatalf("Failed to load questions: %s", err)
	}

	allQuestions, err = generateQuestions("questions.csv", true)
	if err != nil {
		log.Fatalf("Failed to load questions: %s", err)
	}

	setupRoutes() // Configuration des routes HTTP

	// Démarrage du serveur
	fmt.Println("Server started on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
