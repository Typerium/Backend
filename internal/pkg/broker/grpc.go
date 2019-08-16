package broker

import (
	"context"
	"net"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"typerium/internal/pkg/logging"
)

func NewGRPCServer(errorsMapper map[error]error) *GRPCServer {
	if errorsMapper == nil {
		errorsMapper = make(map[error]error)
	}

	server := &GRPCServer{
		log:          logging.New("grpc_server"),
		errorsMapper: errorsMapper,
	}

	server.Server = grpc.NewServer(
		grpc.UnaryInterceptor(server.unaryInterceptor),
		grpc.StreamInterceptor(server.streamInterceptor),
	)
	grpc_health_v1.RegisterHealthServer(server.Server, health.NewServer())

	return server
}

type GRPCServer struct {
	*grpc.Server
	wg           sync.WaitGroup
	log          *zap.Logger
	errorsMapper map[error]error
}

// grpc errors
var (
	InternalGRPCError          = status.Error(codes.Unknown, "internal error")
	NotFoundGRPCError          = status.Error(codes.NotFound, "not found")
	BadInputArgumentsGRPCError = status.Error(codes.InvalidArgument, "bad input arguments")
)

func (s *GRPCServer) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	if req == nil {
		return nil, BadInputArgumentsGRPCError
	}

	resp, err = handler(ctx, req)
	if err != nil {
		err = s.handlerError(err)
		return
	}
	return
}

func (s *GRPCServer) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	err = handler(srv, ss)
	if err != nil {
		err = s.handlerError(err)
		return
	}
	return
}

func (s *GRPCServer) handlerError(input error) error {
	input = errors.Cause(input)
	_, ok := status.FromError(input)
	if ok {
		return input
	}

	output, ok := s.errorsMapper[input]
	if ok {
		return output
	}

	s.log.Warn("unprocessed error", zap.Error(input))
	return InternalGRPCError
}

func (s *GRPCServer) ServeOnAddress(addr string) {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		s.log.Fatal("can't start server", zap.Error(err))
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err = s.Serve(ln)
		if err != nil {
			s.log.Fatal("staring is failed", zap.Error(err))
		}
	}()
	s.log.Info("server is started", zap.String("address", addr))
}

func (s *GRPCServer) GracefulStop() {
	s.Server.GracefulStop()
	s.wg.Wait()
	s.log.Info("server is stopped")
}

type GRPCClientFactory interface {
	Acquire(ctx context.Context) (*grpc.ClientConn, error)
	Release(client *grpc.ClientConn)
}

func NewGRPCClientFactory(log *zap.Logger, uri string, opts ...grpc.DialOption) GRPCClientFactory {
	log = log.Named("grpc_client").With(zap.String("uri", uri))
	defer log.Info("client is created")
	if len(opts) == 0 {
		opts = append(opts, grpc.WithInsecure())
	}
	return &grpcClientFactory{
		log:  log,
		uri:  uri,
		opts: opts,
	}
}

type grpcClientFactory struct {
	log  *zap.Logger
	uri  string
	opts []grpc.DialOption
}

func (f *grpcClientFactory) Acquire(ctx context.Context) (*grpc.ClientConn, error) {
	client, err := grpc.DialContext(ctx, f.uri, f.opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client, nil
}

func (f *grpcClientFactory) Release(client *grpc.ClientConn) {
	err := client.Close()
	if err != nil {
		f.log.Error("can't close connection", zap.Error(err))
	}
}
