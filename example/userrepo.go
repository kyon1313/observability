package main

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserRepository interface {
	AddUser(ctx context.Context, user User) error
	GetUserByName(ctx context.Context, name string) (*User, error)
}

type userRepository struct {
	tracer trace.Tracer
}

func NewUserRepository(tracer trace.Tracer) UserRepository {
	return &userRepository{tracer: tracer}
}

func (r *userRepository) AddUser(ctx context.Context, u User) error {
	_, span := r.tracer.Start(ctx, "repo.AddUser")
	defer span.End()

	for _, user := range users {
		if user.ID == u.ID {
			return errors.New("userid already exist")
		}
	}

	users = append(users, u)
	return nil
}

func (r *userRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	_, span := r.tracer.Start(ctx, "repo.GetUserByName")
	defer span.End()

	for _, user := range users {
		if user.Name == name {
			return &user, nil
		}
	}
	err := fmt.Errorf("name:%s user not found", name)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.String("error.source", "repository"))
	return nil, err
}
