package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"webber/graph/generated"
	"webber/graph/model"
	"webber/models"
)

func (r *mutationResolver) CreateNote(ctx context.Context, input model.NewNote) (string, error) {
	//gc, err := GinContextFromContext(ctx)

	n := &models.Note{
		Title:     input.Title,
		Text:      &input.Text,
		Citation:  input.Citation,
		RelatedTo: nil,
	}

	uuid, err := r.Dao.CreateNote(n)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
