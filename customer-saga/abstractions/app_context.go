package abstractions

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"
)

type IAppContext interface {
	CreateSpan(ref string) (*AppContext, trace.Span)
	Logger() *log.Logger
	GetTraceId() string
	GetSpanId() string
	Wait()
	Terminate()
}

type AppContext struct {
	context.Context
	tracer           trace.Tracer
	span             trace.Span
	terminationToken chan bool
}

func CreateAppContext(ctx context.Context, tracer trace.Tracer, span trace.Span) *AppContext {
	return &AppContext{
		Context:          ctx,
		tracer:           tracer,
		span:             span,
		terminationToken: make(chan bool),
	}
}

func (w *AppContext) CreateSpan(ref string) (*AppContext, trace.Span) {
	ctx, span := w.tracer.Start(w.Context, ref)
	appCtx := CreateAppContext(ctx, w.tracer, span)

	return appCtx, span
}

func (w *AppContext) Logger() *log.Logger {
	return w.Context.Value("logger").(*log.Logger)
}

func (w *AppContext) GetTraceId() string {
	return w.span.SpanContext().TraceID().String()
}

func (w *AppContext) GetSpanId() string {
	return w.span.SpanContext().SpanID().String()
}

func (w *AppContext) Wait() {
	<-w.terminationToken
}

func (w *AppContext) Terminate() {
	w.span.End()
	w.terminationToken <- true
}
