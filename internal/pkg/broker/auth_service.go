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

func NewAuthClient(uri string, opts ...grpc.DialOption) proto.AuthServiceClient {
	log := logging.New("").With(zap.String("service", "auth_service"))
	return &authClient{NewGRPCClientFactory(log, uri, opts...)}
}

type authClient struct {
	GRPCClientFactory
}

func (c *authClient) CreateUser(ctx context.Context, in *proto.NewAuthUser, opts ...grpc.CallOption,
) (*proto.AuthUser, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewAuthServiceClient(client).CreateUser(ctx, in, opts...)
}

func (c *authClient) DeleteUser(ctx context.Context, in *proto.UserIdentifier, opts ...grpc.CallOption,
) (*empty.Empty, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewAuthServiceClient(client).DeleteUser(ctx, in, opts...)
}

func (c *authClient) SignIn(ctx context.Context, in *proto.AuthCredentials, opts ...grpc.CallOption,
) (*proto.Session, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewAuthServiceClient(client).SignIn(ctx, in, opts...)
}

func (c *authClient) SignOut(ctx context.Context, in *proto.AccessCredentials, opts ...grpc.CallOption,
) (*proto.Session, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewAuthServiceClient(client).SignOut(ctx, in, opts...)
}

func (c *authClient) RefreshSession(ctx context.Context, in *proto.AccessCredentials, opts ...grpc.CallOption,
) (*proto.Session, error) {
	client, err := c.Acquire(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer c.Release(client)

	return proto.NewAuthServiceClient(client).RefreshSession(ctx, in, opts...)
}
