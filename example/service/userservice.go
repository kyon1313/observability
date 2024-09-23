package service

import (
	"context"

	"github.com/kyon1313/observability/example/model"
	"github.com/kyon1313/observability/example/repo"
	apw_tracing "github.com/kyon1313/observability/tracing"
)

type UserService interface {
	GetUser(ctx context.Context, name string) (*model.User, error)
	Tracer() apw_tracing.OtelTracing
}

type userService struct {
	repo   repo.UserRepository
	tracer apw_tracing.OtelTracing
}

func NewUserService(repo repo.UserRepository, tracer apw_tracing.OtelTracing) UserService {
	return &userService{repo: repo, tracer: tracer}
}

func (s *userService) GetUser(ctx context.Context, name string) (*model.User, error) {
	var (
		err      error
		response *model.User
	)

	ctx, span := s.tracer.StartSpan(ctx, "service.GetUser")

	response, err = s.repo.GetUserByName(ctx, name)
	if err != nil {
		s.tracer.LogTrace(span, &err, "service.GetUser.Error", response)()
		return nil, err
	}

	s.tracer.LogTrace(span, &err, "service.GetUser.Success", response)()
	return response, nil
}

func (s *userService) Tracer() apw_tracing.OtelTracing {
	return s.tracer
}
