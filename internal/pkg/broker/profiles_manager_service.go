package broker

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"typerium/internal/pkg/broker/proto"
	"typerium/internal/pkg/logging"
)

func NewProfilesManagerClient(uri string, opts ...grpc.DialOption) proto.ProfilesManagerServiceClient {
	log := logging.New("").With(zap.String("service", "profiles_manager"))
	return &profilesManagerClient{NewGRPCClientFactory(log, uri, opts...)}
}

type profilesManagerClient struct {
	GRPCClientFactory
}

func (c *profilesManagerClient) CreateUser(ctx context.Context, in *proto.NewProfilesUser, opts ...grpc.CallOption,
) (*proto.ProfilesUser, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewProfilesManagerServiceClient(client).CreateUser(ctx, in, opts...)
}

func (c *profilesManagerClient) DeleteUser(ctx context.Context, in *proto.UserIdentifier, opts ...grpc.CallOption,
) (*empty.Empty, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewProfilesManagerServiceClient(client).DeleteUser(ctx, in, opts...)
}

func (c *profilesManagerClient) GetUserByID(ctx context.Context, in *proto.UserIdentifier, opts ...grpc.CallOption,
) (*proto.ProfilesUser, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewProfilesManagerServiceClient(client).GetUserByID(ctx, in, opts...)
}
