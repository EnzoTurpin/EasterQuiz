package main

// Importation des packages nécessaires.
import (
	"bufio"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Structure représentant une question de quiz.
type Question struct {
	Text       string   // Le texte de la question.
	Answers    []string // Les options de réponses.
	CorrectAns string   // La réponse correcte.
	Difficulty string   // Le niveau de difficulté.
}

// Structure pour les données de retour après réponse d'une question.
type FeedbackData struct {
	Correct        bool   // Si la réponse de l'utilisateur était correcte.
	UserAnswer     string // Réponse donnée par l'utilisateur.
	CorrectAns     string // Bonne réponse.
	QuestionText   string // Texte de la question.
	IsLastQuestion bool   // Si c'est la dernière question du quiz.
}

// Fonction pour générer les questions à partir d'un fichier, avec option de mélange.
func generateQuestions(filename string, shuffle bool) ([]Question, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	var tempQuiz []Question
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 7 {
			continue // Ignore les lignes mal formatées.
		}
		// Construction et ajout d'une nouvelle question à tempQuiz.
		tempQuiz = append(tempQuiz, Question{
			Text:       parts[0],
			Answers:    parts[1:5],
			CorrectAns: parts[5],
			Difficulty: parts[6],
		})
	}

	if shuffle {
		// Mélange aléatoire des questions si requis.
		rand.Shuffle(len(tempQuiz), func(i, j int) { tempQuiz[i], tempQuiz[j] = tempQuiz[j], tempQuiz[i] })
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return tempQuiz, nil
}

// selectRandomQuestions sélectionne des questions aléatoires dans le quiz filtré.
func selectRandomQuestions(questions []Question, n int) []Question {
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	if len(questions) > n {
		return questions[:n]
	}
	return questions
}

// Charge les questions à partir d'un fichier CSV et les mélange.
func loadQuestions(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open the CSV file: %s", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		line, error := csvReader.Read()
		if error != nil {
			break
		}
		correctAns, _ := strconv.Atoi(line[2])
		quiz = append(quiz, Question{
			Text:       line[0],
			Answers:    []string{line[1], line[3], line[4], line[5]},
			CorrectAns: strconv.Itoa(correctAns),
		})

	}
}

// Wrapper qui ajuste le type MIME pour les fichiers .js
func setMIMEType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		next.ServeHTTP(w, r)
	})
}

// filterQuestionsByDifficulty filtre les questions par difficulté.
func filterQuestionsByDifficulty(difficulty string) []Question {
	filtered := []Question{}
	for _, q := range quiz {
		if q.Difficulty == difficulty {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

// Affiche un formulaire de sélection de difficulté.
func renderDifficultySelection(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/difficulty_selection.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Affiche le formulaire de question avec les options de réponse.
func renderQuestionForm(w http.ResponseWriter, q Question, index int) {
	tmpl, err := template.ParseFiles("src/templates/question_form.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du formulaire de question", http.StatusInternalServerError)
		return
	}

	// Mélanger les réponses ici
	rand.Shuffle(len(q.Answers), func(i, j int) { q.Answers[i], q.Answers[j] = q.Answers[j], q.Answers[i] })

	data := struct {
		Difficulty string
		Text       string
		Answers    []string
		Index      int
	}{
		Difficulty: q.Difficulty,
		Text:       q.Text,
		Answers:    q.Answers,
		Index:      index,
	}
	tmpl.Execute(w, data)

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors du rendu du formulaire de question", http.StatusInternalServerError)
	}
}

// Affiche un retour à l'utilisateur après soumission d'une réponse.
func renderFeedback(w http.ResponseWriter, f *FeedbackData) {
	tmpl, err := template.ParseFiles("src/templates/feedback.html")
	if err != nil {
		http.Error(w, "Error loading feedback template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, f)
	if err != nil {
		http.Error(w, "Error rendering feedback", http.StatusInternalServerError)
	}
}

// Affiche la page de résultats avec le score de l'utilisateur.
func displayResultsPage(w http.ResponseWriter, r *http.Request) {
	userID := getSessionUserID(r)
	score := scores[userID]
	totalQuestions := len(quiz)

	scores[userID] = 0

	quiz = nil

	http.SetCookie(w, &http.Cookie{
		Name:   "quiz-session",
		Value:  "0",
		Path:   "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "difficulty",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	fmt.Fprintf(w, `<html><body>
		<p>Votre score : %d sur %d</p>
		<form action="/quiz" method="get">
			<button type="submit">Recommencer le quiz</button>
		</form>
		</body></html>`, score, totalQuestions)
}

// Obtient l'indice de la question actuelle à partir d'un cookie.
func getCurrentQuestionIndex(r *http.Request) (int, error) {
	cookie, err := r.Cookie("quiz-session")
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(cookie.Value)
}

// Vérifie si la réponse donnée par l'utilisateur est correcte.
func isValidAnswer(userAnswer string, q Question) bool {
	return userAnswer == q.CorrectAns
}

// Traite la réponse de l'utilisateur et redirige vers le retour ou la question suivante.
func processUserAnswer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	userAnswer := r.FormValue("answer")
	correct := isValidAnswer(userAnswer, quiz[currentQuestionIndex])

	if correct {
		scores["defaultUser"]++
	}

	correctAnswer := quiz[currentQuestionIndex].CorrectAns

	isLastQuestion := currentQuestionIndex == len(quiz)-1

	currentQuestionIndex++

	lastFeedback = &FeedbackData{
		Correct:        correct,
		UserAnswer:     userAnswer,
		CorrectAns:     correctAnswer,
		QuestionText:   quiz[currentQuestionIndex-1].Text,
		IsLastQuestion: isLastQuestion,
	}

	http.Redirect(w, r, "/feedback", http.StatusSeeOther)
}

// Met à jour le cookie de session pour la prochaine question.
func updateSessionCookie(w http.ResponseWriter, r *http.Request, nextQuestionIndex int) {
	http.SetCookie(w, &http.Cookie{
		Name:  "quiz-session",
		Value: strconv.Itoa(nextQuestionIndex),
		Path:  "/",
	})
}

// Retourne un identifiant utilisateur par défaut; pourrait être amélioré pour gérer plusieurs utilisateurs.
func getSessionUserID(r *http.Request) string {
	return "defaultUser"
}

// Définit l'en-tête HTTP pour le contenu en UTF-8.
func setUTF8Header(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
