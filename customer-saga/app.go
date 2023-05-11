package main

import (
	"context"
	"log"
	"os"

	"github.com/gdata/customer-saga/abstractions"
	"github.com/gdata/customer-saga/workers"

	"go.opentelemetry.io/otel"
)

type IApp interface {
	RegisterWorker(worker workers.IWorker) IApp
	Run()
}

type App struct {
	Context *abstractions.AppContext
	workers []workers.IWorker
}

func Create() IApp {
	logger := log.New(os.Stdout, "", 0)

	tracer := otel.Tracer("app-tracer")
	ctx, span := tracer.Start(context.Background(), "app-spanner")
	ctx = context.WithValue(ctx, "logger", logger)
	ctx = context.WithValue(ctx, "span", span)

	return &App{
		Context: abstractions.CreateAppContext(ctx, tracer, span),
	}
}

func (a *App) Run() {
	a.Context.Logger().Println("Starting application")

	for _, worker := range a.workers {
		go worker.ExecuteWithContext(a.Context)
	}

	a.Context.Wait()
	a.Context.Logger().Println("Shuting down completed")
}

func (a *App) RegisterWorker(worker workers.IWorker) IApp {
	a.workers = append(a.workers, worker)
	return a
}
