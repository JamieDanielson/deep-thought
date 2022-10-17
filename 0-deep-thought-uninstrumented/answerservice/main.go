package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func provideAnswer(ctx context.Context) string {
	for {
		min := 1
		max := 10000
		answer := strconv.Itoa((rand.Intn(max-min) + min))
		if answer == "42" {
			return answer
		}
	}
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	answer := func(ctx context.Context) string {
		return provideAnswer(ctx)
	}(ctx)

	_, _ = fmt.Fprintf(w, "%s", answer)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/answerservice", answerHandler)

	wrappedHandler := http.Handler(mux)

	log.Println("Listening on http://localhost:5678/answerservice")
	log.Fatal(http.ListenAndServe(":5678", wrappedHandler))
}
