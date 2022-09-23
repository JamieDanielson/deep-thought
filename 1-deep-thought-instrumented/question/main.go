package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
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

func determineQuestion() string {
	return "what is the answer to the ultimate question of life, the universe, and everything?"
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	question := func(ctx context.Context) string {
		_, span := tracer.Start(ctx, "✨ pondering the question ✨")
		defer span.End()
		return determineQuestion()
	}(ctx)

	_, _ = fmt.Fprintf(w, "%v", question)

}

func main() {

	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp)

	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	otel.SetTracerProvider(tp)

	tracer = tp.Tracer("deep-thought/question")

	mux := http.NewServeMux()
	mux.HandleFunc("/question", questionHandler)

	wrappedHandler := otelhttp.NewHandler(mux, "question")

	log.Println("Listening on http://localhost:1234/question")
	log.Fatal(http.ListenAndServe(":1234", wrappedHandler))
}
