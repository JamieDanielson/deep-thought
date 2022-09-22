package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var (
	questionServiceUrl = getEnv("QUESTION_ENDPOINT", "http://localhost:1234") + "/question"
	answerServiceUrl   = getEnv("ANSWER_ENDPOINT", "http://localhost:5678") + "/answer"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		question := getQuestion(r.Context())
		answer := getAnswer(r.Context())

		_, _ = fmt.Fprintf(w, "%s\n%s\n", question, answer)
	})

	wrappedHandler := http.Handler(mux)

	log.Println("Listening on http://localhost:4242/")
	log.Fatal(http.ListenAndServe(":4242", wrappedHandler))
}

func getQuestion(ctx context.Context) string {
	return makeRequest(ctx, questionServiceUrl)
}

func getAnswer(ctx context.Context) string {
	return makeRequest(ctx, answerServiceUrl)
}

func makeRequest(ctx context.Context, url string) string {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := http.Client{Transport: http.DefaultTransport}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	return string(body)
}
