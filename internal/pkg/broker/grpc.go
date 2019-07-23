package broker

import (
	"context"
	"net"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func NewGRPCServer(log *zap.Logger, errorsMapper map[error]error) *GRPCServer {
	if errorsMapper == nil {
		errorsMapper = make(map[error]error)
	}

	server := &GRPCServer{
		log:          log.Named("grpc"),
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
	InternalGRPCErr          = status.Error(codes.Unknown, "internal error")
	NotFoundGRPCErr          = status.Error(codes.NotFound, "not found")
	BadInputArgumentsGRPCErr = status.Error(codes.InvalidArgument, "bad input arguments")
	UnavailableGRPCErr       = status.Error(codes.Unavailable, "service isn't available")
)

func (s *GRPCServer) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	if req == nil {
		return nil, BadInputArgumentsGRPCErr
	}

	validator, ok := req.(validation.Validatable)
	if ok {
		err = validator.Validate()
		if err != nil {
			return nil, s.handlerError(err)
		}
	}

	resp, err = handler(ctx, req)
	if err != nil {
		return nil, s.handlerError(err)
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
	return InternalGRPCErr
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
	s.log.Info("server is started")
}

func (s *GRPCServer) GracefulStop() {
	s.Server.GracefulStop()
	s.wg.Wait()
	s.log.Info("server is stopped")
}

type grpcClient struct {
	uri            string
	defaultTimeout time.Duration
	log            *zap.Logger
}

func (client *grpcClient) Acquire(ctx context.Context) *grpc.ClientConn {
	if ctx == nil {
		ctx = context.Background()
	}
	_, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, client.defaultTimeout)
		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, client.uri)
	if err != nil {
		client.log.Error("can't connect to grpc server",
			zap.Error(err),
			zap.String("uri", client.uri),
		)
		return nil
	}

	return conn
}

func (client *grpcClient) Release(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}

	err := conn.Close()
	if err != nil {
		client.log.Error("failed closing grpc connection",
			zap.Error(err),
			zap.String("uri", client.uri),
		)
	}
}
