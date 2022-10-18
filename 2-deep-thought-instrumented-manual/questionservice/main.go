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

// set up an OTLP Trace Exporter
func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	client := otlptracegrpc.NewClient()
	return otlptrace.New(ctx, client)
}

// set up a Tracer Provider
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

		// let's add a manual span!
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

  // Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle this error in a sensible manner where possible
	defer func() { _ = tp.Shutdown(ctx) }()

	// Set the Tracer Provider and the W3C Trace Context propagator as globals.
	// Important, otherwise this won't let you see distributed traces be connected!
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	tracer = tp.Tracer("deep-thought/questionservice")

	mux := http.NewServeMux()
	mux.HandleFunc("/questionservice", questionHandler)

	wrappedHandler := otelhttp.NewHandler(mux, "questionservice")

	log.Println("Listening on http://localhost:1234/questionservice")
	log.Fatal(http.ListenAndServe(":1234", wrappedHandler))
}
