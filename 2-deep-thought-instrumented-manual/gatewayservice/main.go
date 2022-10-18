package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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
	questionserviceUrl = getEnv("QUESTION_ENDPOINT", "http://localhost:1234") + "/questionservice"
	answerserviceUrl   = getEnv("ANSWER_ENDPOINT", "http://localhost:5678") + "/answerservice"
	tracer             trace.Tracer
)

// set up an OTLP Trace Exporter
func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	client := otlptracegrpc.NewClient()
	return otlptrace.New(ctx, client)
}

// set up a Tracer Provider
func newTracerProvider(exp *otlptrace.Exporter) *sdktrace.TracerProvider {
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
	tp := newTracerProvider(exp)

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	// Set the Tracer Provider and the W3C Trace Context propagator as globals.
	// Important, otherwise this won't let you see distributed traces be connected!
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	tracer = tp.Tracer("deep-thought/gatewayservice")

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

	// let's add a manual span!
	var getQuestionSpan trace.Span
	ctx, getQuestionSpan = tracer.Start(ctx, "✨ call /questionservice ✨")
	defer getQuestionSpan.End()

	return makeRequest(ctx, questionserviceUrl)
}

func getAnswer(ctx context.Context) string {

	// let's add a manual span!
	var getAnswerSpan trace.Span
	ctx, getAnswerSpan = tracer.Start(ctx, "✨ call /answerservice ✨")
	time.Sleep(1 * time.Second)
	// add interesting detail to this span
	getAnswerSpan.SetAttributes(attribute.String("important_note", "don't panic"))
	defer getAnswerSpan.End()

	return makeRequest(ctx, answerserviceUrl)
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
