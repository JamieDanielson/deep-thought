package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/attribute"

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
	tracer             trace.Tracer
)

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	client := otlptracegrpc.NewClient()
	return otlptrace.New(ctx, client)
}

func newTraceProvider(exp *otlptrace.Exporter) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)
}

func main() {
	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Set the Tracer Provider and the W3C Trace Context propagator as globals.
	// Important, otherwise this won't let you see distributed traces be connected!
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	tracer = tp.Tracer("deep-thought/frontend")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		question := getQuestion(r.Context())
		answer := getAnswer(r.Context())

		_, _ = fmt.Fprintf(w, "%s\n%s\n", question, answer)
	})

	wrappedHandler := otelhttp.NewHandler(mux, "main")

	log.Println("Listening on http://localhost:4242/")
	log.Fatal(http.ListenAndServe(":4242", wrappedHandler))
}

func getQuestion(ctx context.Context) string {
	var getQuestionSpan trace.Span
	ctx, getQuestionSpan = tracer.Start(ctx, "✨ call /question ✨")
	defer getQuestionSpan.End()
	return makeRequest(ctx, questionServiceUrl)
}

func getAnswer(ctx context.Context) string {
	var getAnswerSpan trace.Span
	ctx, getAnswerSpan = tracer.Start(ctx, "✨ call /answer ✨")
	getAnswerSpan.SetAttributes(attribute.String("important_note", "don't panic"))
	defer getAnswerSpan.End()
	return makeRequest(ctx, answerServiceUrl)
}

func makeRequest(ctx context.Context, url string) string {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
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
