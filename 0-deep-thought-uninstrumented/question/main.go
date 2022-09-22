package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func determineQuestion() string {
	questions := []string{
		"what is the answer to the ultimate question of life, the universe, and everything?", "what is 6 * 7?",
	}
	return questions[rand.Intn(len(questions))]
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	question := func(ctx context.Context) string {
		return determineQuestion()
	}(ctx)

	_, _ = fmt.Fprintf(w, "%v", question)

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/question", questionHandler)

	wrappedHandler := http.Handler(mux)

	log.Println("Listening on http://localhost:1234/question")
	log.Fatal(http.ListenAndServe(":1234", wrappedHandler))
}
