package broker

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"typerium/internal/pkg/broker/proto"
)

type ProfilesManagerServiceClient interface {
	Acquire(ctx context.Context) (proto.ProfilesManagerServiceClient, error)
	Release(client proto.ProfilesManagerServiceClient)
}

func NewProfilesManagerServiceClient(log *zap.Logger, uri string, ttl time.Duration) ProfilesManagerServiceClient {
	return &profilesManagerImpl{
		&grpcClient{
			log:            log.Named("grpc_client").With(zap.String("service", "profiles_manager")),
			uri:            uri,
			defaultTimeout: ttl,
		},
	}
}

type profilesManagerImpl struct {
	*grpcClient
}

type profilesManagerClient struct {
	conn *grpc.ClientConn
	proto.ProfilesManagerServiceClient
}

func (s *profilesManagerImpl) Acquire(ctx context.Context) (proto.ProfilesManagerServiceClient, error) {
	conn := s.grpcClient.Acquire(ctx)
	if conn == nil {
		return nil, UnavailableGRPCErr
	}

	return &profilesManagerClient{
		conn:                         conn,
		ProfilesManagerServiceClient: proto.NewProfilesManagerServiceClient(conn),
	}, nil
}

func (s *profilesManagerImpl) Release(client proto.ProfilesManagerServiceClient) {
	cc, ok := client.(*profilesManagerClient)
	if !ok {
		return
	}

	s.grpcClient.Release(cc.conn)
}
