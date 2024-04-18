package main

import (
	"html/template"
	"log"
	"net/http"
)

// Configuration des routes et des gestionnaires associés.
func setupRoutes() {
	fssrcWithMIME := setMIMEType(http.FileServer(http.Dir("src")))
	http.Handle("/src/", http.StripPrefix("/src/", fssrcWithMIME))

	fsImg := http.FileServer(http.Dir("img"))
	http.Handle("/img/", http.StripPrefix("/img/", fsImg))

	// Routes pour les différentes parties du quiz.
	http.HandleFunc("/", difficultySelectionHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/finish", finishHandler)
	http.HandleFunc("/reset-quiz", resetQuizHandler)
}

func difficultySelectionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/difficulty_selection.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Gestionnaire pour les requêtes de la page du quiz.
func quizHandler(w http.ResponseWriter, r *http.Request) {
	setUTF8Header(w)
	if r.Method == "POST" {
		r.ParseForm()
		difficulty := r.FormValue("difficulty")

		if difficulty != "" {
			filteredQuiz := filterQuestionsByDifficulty(difficulty)
			selectedQuestions := selectRandomQuestions(filteredQuiz, 5)

			currentDifficulty = difficulty
			quiz = selectedQuestions
			currentQuestionIndex = 0

			http.Redirect(w, r, "/quiz", http.StatusSeeOther)
			return
		} else {
			processUserAnswer(w, r)

		}
	} else if r.Method == "GET" {
		if currentQuestionIndex < len(quiz) {
			question := quiz[currentQuestionIndex]

			data := struct {
				Text       string
				Answers    []string
				Index      int
				Difficulty string
			}{
				Text:       question.Text,
				Answers:    question.Answers,
				Index:      currentQuestionIndex + 1,
				Difficulty: currentDifficulty,
			}

			tmpl, err := template.ParseFiles("src/templates/question_form.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Redirect(w, r, "/finish", http.StatusSeeOther)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Affiche la page de retour après la soumission d'une réponse.
func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	if lastFeedback != nil {
		renderFeedback(w, lastFeedback)
		lastFeedback = nil
	} else {
		http.Error(w, "No feedback available", http.StatusInternalServerError)
	}
}

// Affiche la page de résultats à la fin du quiz.
func finishHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/results.html")
	if err != nil {
		log.Fatalf("Error loading template: %s", err)
	}

	score := scores["defaultUser"]
	totalQuestions := len(quiz)
	data := map[string]interface{}{
		"Score":          score,
		"TotalQuestions": totalQuestions,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

// Réinitialise le quiz et redirige vers la page de sélection de difficulté.
func resetQuizHandler(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	scores[userID] = 0
	currentQuestionIndex = 0
	quiz, _ = generateQuestions("questions.csv", true)
	lastFeedback = nil

	http.SetCookie(w, &http.Cookie{
		Name:   "quiz-session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "difficulty",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Gère l'affichage des résultats.
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	displayResultsPage(w, r)
}
