package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/danglebary/beatstore-backend-go/graph/generated"
	"github.com/danglebary/beatstore-backend-go/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input model.DeleteUserInput) (*model.DeleteUserResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateBeat(ctx context.Context, input model.CreateBeatInput) (*model.Beat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateBeat(ctx context.Context, input model.UpdateBeatInput) (*model.Beat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteBeat(ctx context.Context, input model.DeleteBeatInput) (*model.DeleteBeatResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Beats(ctx context.Context, obj *model.User) ([]*model.Beat, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
