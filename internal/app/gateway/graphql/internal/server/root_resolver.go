package server

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SignUp(ctx context.Context, email string, userName *string, phone *string, password string) (*User, error) {
	panic("not implemented")
}
func (r *mutationResolver) SignIn(ctx context.Context, login string, password string) (*Tokens, error) {
	panic("not implemented")
}
func (r *mutationResolver) SignInByAccount(ctx context.Context, accountType LoginAccount) (*Tokens, error) {
	panic("not implemented")
}
func (r *mutationResolver) ResetPassword(ctx context.Context, email string) (bool, error) {
	panic("not implemented")
}
func (r *mutationResolver) SaveImage(ctx context.Context, image graphql.Upload, album *string) (*Image, error) {
	panic("not implemented")
}
func (r *mutationResolver) PublishToMarketPlace(ctx context.Context, imageID string) (*MarketPlaceImage, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) CurrentUser(ctx context.Context) (*User, error) {
	panic("not implemented")
}
func (r *queryResolver) Albums(ctx context.Context) ([]*Album, error) {
	panic("not implemented")
}
func (r *queryResolver) Images(ctx context.Context, limit int, offset int) ([]*Image, error) {
	panic("not implemented")
}
func (r *queryResolver) MarketPlaceLots(ctx context.Context, limit int, offset int) ([]*MarketPlaceImage, error) {
	panic("not implemented")
}
