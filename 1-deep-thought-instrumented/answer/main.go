package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	// "go.opentelemetry.io/otel/attribute"
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

func provideAnswer(ctx context.Context) string {

	// this is a slow computation!
	for {
		min := 1
		max := 100000
		answer := strconv.Itoa((rand.Intn(max-min) + min))
		trace.SpanFromContext(ctx).AddEvent(answer)
		if answer == "42" {
			return answer
		}
	}
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	answer := func(ctx context.Context) string {
		_, span := tracer.Start(ctx, "✨ thinking about the answer ✨")
		time.Sleep(1 * time.Second)
		defer span.End()
		return provideAnswer(ctx)
	}(ctx)

	_, _ = fmt.Fprintf(w, "%s", answer)

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

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	tracer = tp.Tracer("deep-thought/answer")

	mux := http.NewServeMux()
	mux.HandleFunc("/answer", answerHandler)

	wrappedHandler := otelhttp.NewHandler(mux, "answer")

	log.Println("Listening on http://localhost:5678/answer")
	log.Fatal(http.ListenAndServe(":5678", wrappedHandler))
}
