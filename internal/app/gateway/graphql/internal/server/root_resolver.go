package server

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	"typerium/internal/app/gateway/attributes"
	"typerium/internal/pkg/broker/proto"
)

type Resolvers struct {
	TokenAttr  attributes.AttributeString
	UserIDAttr attributes.AttributeUUID

	AuthService            proto.AuthServiceClient
	ProfilesManagerService proto.ProfilesManagerServiceClient
}

func NewResolvers(authService proto.AuthServiceClient,
	profilesManagerService proto.ProfilesManagerServiceClient) *Resolvers {
	return &Resolvers{
		TokenAttr:  attributes.NewSchemaTokenAttribute(),
		UserIDAttr: attributes.NewUserIDAttribute(),

		AuthService:            authService,
		ProfilesManagerService: profilesManagerService,
	}
}

func (r *Resolvers) Mutation() MutationResolver {
	return r
}
func (r *Resolvers) Query() QueryResolver {
	return r
}

func (r *Resolvers) SaveImage(ctx context.Context, image graphql.Upload, album *string) (*Image, error) {
	panic("not implemented")
}
func (r *Resolvers) PublishToMarketPlace(ctx context.Context, imageID string) (*MarketPlaceImage, error) {
	panic("not implemented")
}

func (r *Resolvers) Albums(ctx context.Context) ([]*Album, error) {
	panic("not implemented")
}
func (r *Resolvers) Images(ctx context.Context, limit int, offset int) ([]*Image, error) {
	panic("not implemented")
}
func (r *Resolvers) MarketPlaceLots(ctx context.Context, limit int, offset int) ([]*MarketPlaceImage, error) {
	panic("not implemented")
}
