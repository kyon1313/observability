package main

import (
	"context"
	"errors"
	"math/rand"

	"go.opentelemetry.io/otel/trace"
)

type Repository interface {
	GetRandomNumber(ctx context.Context) (int, error)
	GetError(ctx context.Context) error
}

type repository struct {
	tracer trace.Tracer
}

func NewRepository(t trace.Tracer) Repository {
	return &repository{tracer: t}
}

func (r *repository) GetRandomNumber(ctx context.Context) (int, error) {
	ctx, span := r.tracer.Start(ctx, "repository.GetRandomNumber")
	defer span.End()

	return rand.Intn(5) + 1, nil
}

func (r *repository) GetError(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "repository.GetError")
	defer span.End()

	return errors.New("error request")
}
