package main

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Service interface {
	GenerateRandomNumber(ctx context.Context) (int, error)
	GenerateError(ctx context.Context) error
	Tracer() trace.Tracer
}

type service struct {
	repo   Repository
	tracer trace.Tracer
}

func NewService(r Repository, t trace.Tracer) Service {
	return &service{repo: r, tracer: t}
}

func (s *service) GenerateRandomNumber(ctx context.Context) (int, error) {
	ctx, span := s.tracer.Start(ctx, "service.GenerateRandomNumber")
	defer span.End()

	return s.repo.GetRandomNumber(ctx)
}

func (s *service) GenerateError(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "service.GenerateError")
	defer span.End()

	return s.repo.GetError(ctx)
}

func (s *service) Tracer() trace.Tracer {
	return s.tracer
}
