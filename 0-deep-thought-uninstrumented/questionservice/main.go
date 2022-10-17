package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func determineQuestion() string {
	return "what is the answer to the ultimate question of life, the universe, and everything?"
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
	mux.HandleFunc("/questionservice", questionHandler)

	wrappedHandler := http.Handler(mux)

	log.Println("Listening on http://localhost:1234/questionservice")
	log.Fatal(http.ListenAndServe(":1234", wrappedHandler))
}
