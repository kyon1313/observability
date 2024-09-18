package main

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserService interface {
	AddUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, name string) (*User, error)
	Tracer() trace.Tracer
}

type userService struct {
	repo   UserRepository
	tracer trace.Tracer
}

func NewUserService(repo UserRepository, t trace.Tracer) UserService {
	return &userService{repo: repo, tracer: t}
}

func (s *userService) AddUser(ctx context.Context, user User) error {
	ctx, span := s.tracer.Start(ctx, "service.AddUser")
	defer span.End()
	return s.repo.AddUser(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, name string) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "service.GetUser")
	defer span.End()

	user, err := s.repo.GetUserByName(ctx, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctx = context.WithValue(ctx, errorSourceKey, "service")
		return nil, err
	}
	return user, nil
}

func (s *userService) Tracer() trace.Tracer {
	return s.tracer
}
