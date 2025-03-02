package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

var cities [100]City

type City struct {
	City     string    `json:"city"`
	Country  string    `json:"country"`
	Clues    [2]string `json:"clues"`
	FunFacts [2]string `json:"fun_facts"`
	Trivia   [2]string `json:"trivia"`
}

type QuestionOption struct {
	Label string `json:"label"`
	Id    string `json:"id"`
}

type Question struct {
	Index   int               `json:"index,omitempty"`
	Clues   [2]string         `json:"clues,omitempty"`
	Options [4]QuestionOption `json:"options,omitempty"`
}

type UserSelectedOption struct {
	Index          int            `json:"index,omitempty"`
	SelectedOption QuestionOption `json:"selected_option,omitempty"`
}

type UserSelectionAnswer struct {
	ValidAnswer bool `json:"valid_answer"`
	City        City `json:"city"`
}

func loadData() {
	file, err := os.Open("cities_dataset.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(bytes, &cities)
	if jsonErr != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
}

func httpRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httprate.Limit(getHttpRateLimit(), time.Second)
		next.ServeHTTP(w, r)
	})
}

func getRandomQuestion() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var question Question
		randQuestionIndex := rand.Intn(100)
		randQuestion := cities[randQuestionIndex]
		qOptions := [4]QuestionOption{}
		for i := randQuestionIndex; i <= randQuestionIndex+3; i++ {
			qOption := QuestionOption{}
			qOption.Label = cities[i].City
			qOption.Id = getMD5Hex(cities[i].City)
			qOptions[i-randQuestionIndex] = qOption
		}
		shuffleOptions(qOptions)
		question.Index = randQuestionIndex
		question.Clues = randQuestion.Clues
		question.Options = qOptions
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(toJson(question))
	})
}

func checkUserAnswer() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestData UserSelectedOption
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, badRequest, http.StatusBadRequest)
			return
		}
		city := cities[requestData.Index]
		var userSelectionAnswer UserSelectionAnswer
		userSelectionAnswer.ValidAnswer = isSelectedOptionCorrect(city, requestData.SelectedOption)
		userSelectionAnswer.City = city
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(toJson(userSelectionAnswer))
	})
}

func main() {
	loadData()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(httpRateLimit)
	router.Use(middleware.AllowContentType("application/json"))
	router.Get("/api/question", getRandomQuestion())
	router.Post("/api/question", checkUserAnswer())
	serverAddr := serverListenerAddress()
	log.Println("Sever running at", serverAddr)
	log.Fatalln(http.ListenAndServe(serverAddr, router))
}
