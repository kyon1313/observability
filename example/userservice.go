package main

import (
	"context"

	apw_tracing "github.com/kyon1313/observability/tracing"
)

type UserService interface {
	GetUser(ctx context.Context, name string) (*User, error)
	Tracer() apw_tracing.OtelTracing
}

type userService struct {
	repo   UserRepository
	tracer apw_tracing.OtelTracing
}

func NewUserService(repo UserRepository, tracer apw_tracing.OtelTracing) UserService {
	return &userService{repo: repo, tracer: tracer}
}

func (s *userService) GetUser(ctx context.Context, name string) (*User, error) {
	ctx, span := s.tracer.StartSpan(ctx, "service.GetUser")
	defer s.tracer.EndSpan(span)

	user, err := s.repo.GetUserByName(ctx, name)
	if err != nil {
		s.tracer.RecordError(span, err, "service")
		return nil, err
	}

	s.tracer.SetOKStatus(span, "User found")
	return user, nil
}

func (s *userService) Tracer() apw_tracing.OtelTracing {
	return s.tracer
}
