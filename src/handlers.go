package main

import (
	"fmt"
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
	http.HandleFunc("/debug-questions", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Liste des questions</h1>"))
		for _, q := range quiz {
			fmt.Fprintf(w, "<p>Question : %s | Difficulté : %s</p>", q.Text, q.Difficulty)
		}
	})

}

func difficultySelectionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/difficulty_selection.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérification des paramètres pour un éventuel message d'erreur
	errorMessage := ""
	if r.URL.Query().Get("error") != "" {
		errorMessage = r.URL.Query().Get("error")
	}

	// Structure des données envoyées au template
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}

	// Exécution du template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Gestionnaire pour les requêtes de la page du quiz.
func quizHandler(w http.ResponseWriter, r *http.Request) {
	setUTF8Header(w)
	if r.Method == "POST" {
		r.ParseForm()
		difficulty := r.FormValue("difficulty")

		if difficulty != "" {
			// Filtrer les questions à partir de la variable globale allQuestions
			filteredQuiz := filterQuestionsByDifficulty(difficulty)
			if len(filteredQuiz) == 0 {
				log.Printf("Aucune question disponible pour la difficulté : %s\n", difficulty)
				http.Redirect(w, r, "/?error=Aucune+question+disponible+pour+la+difficulté+sélectionnée.+Veuillez+choisir+une+autre+difficulté.", http.StatusSeeOther)
				return
			}

			// Sélectionner des questions aléatoires parmi les questions filtrées
			selectedQuestions := selectRandomQuestions(filteredQuiz, 5)

			log.Printf("Difficulté sélectionnée : %s, Questions disponibles : %d\n", difficulty, len(filteredQuiz))

			if len(selectedQuestions) == 0 {
				http.Redirect(w, r, "/?error=Aucune+question+disponible+pour+la+difficulté+sélectionnée.+Veuillez+choisir+une+autre+difficulté.", http.StatusSeeOther)
				return
			}

			// Mettre à jour les variables globales pour le quiz actuel
			currentDifficulty = difficulty
			quiz = selectedQuestions
			currentQuestionIndex = 0

			// Rediriger vers la première question du quiz
			http.Redirect(w, r, "/quiz", http.StatusSeeOther)
			return
		} else {
			// Traiter la réponse de l'utilisateur pour la question actuelle
			processUserAnswer(w, r)
		}
	} else if r.Method == "GET" {
		if currentQuestionIndex < len(quiz) {
			question := quiz[currentQuestionIndex]

			// Préparer les données pour le template de question
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

			// Charger et exécuter le template de question
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
			// Rediriger vers la page de fin si toutes les questions ont été posées
			http.Redirect(w, r, "/finish", http.StatusSeeOther)
		}
	} else {
		// Gérer les méthodes HTTP non autorisées
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
