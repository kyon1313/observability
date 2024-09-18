package main

import (
	"context"
	"fmt"

	apw_tracing "otel-test/tracing"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, name string) (*User, error)
}

type userRepository struct {
	tracer apw_tracing.OtelTracing
}

func NewUserRepository(tracer apw_tracing.OtelTracing) UserRepository {
	return &userRepository{tracer: tracer}
}

func (r *userRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	_, span := r.tracer.StartSpan(ctx, "repo.GetUserByName")
	defer r.tracer.EndSpan(span)

	for _, user := range users {
		if user.Name == name {
			r.tracer.SetOKStatus(span, "User found")
			return &user, nil
		}
	}

	err := fmt.Errorf("name:%s user not found", name)
	r.tracer.RecordError(span, err, "repository")

	return nil, err
}
