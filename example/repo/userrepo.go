package repo

import (
	"context"
	"fmt"

	"github.com/kyon1313/observability/example/model"
	apw_tracing "github.com/kyon1313/observability/tracing"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}

type userRepository struct {
	tracer apw_tracing.OtelTracing
}

func NewUserRepository(tracer apw_tracing.OtelTracing) UserRepository {
	return &userRepository{tracer: tracer}
}

func (r *userRepository) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var (
		err      error
		response *model.User
	)

	_, span := r.tracer.StartSpan(ctx, "repo.GetUser")

	for _, user := range model.Users {
		if user.Name == name {

			response = &user
			r.tracer.LogTrace(span, &err, "repo.GetUser.Success", response)()
			return &user, nil
		}
	}

	err = fmt.Errorf("name:%s user not found", name)
	r.tracer.LogTrace(span, &err, "repo.GetUser.Error", response)()

	return nil, err
}
