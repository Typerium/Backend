package handlers

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"

	"typerium/internal/app/profiles_manager/store"
	"typerium/internal/pkg/broker/proto"
)

func NewGRPCServer(db store.Store) proto.ProfilesManagerServiceServer {
	return &grpcServer{db: db}
}

type grpcServer struct {
	db store.Store
}

func userToProto(in *store.User) *proto.ProfilesUser {
	id := in.ID.String()
	out := &proto.ProfilesUser{
		ID:       &id,
		Username: in.Username,
		Email:    in.Email,
		Phone:    in.Phone,
		Tm: proto.TimeMark{
			CreatedAt: in.CreatedAt,
			UpdatedAt: in.UpdatedAt,
		},
	}
	return out
}

func (s *grpcServer) CreateUser(ctx context.Context, in *proto.ProfilesUser) (out *proto.ProfilesUser, err error) {
	user := &store.User{
		ID:       uuid.Nil,
		Username: in.Username,
		Email:    in.Email,
		Phone:    in.Phone,
	}
	if in.ID != nil {
		user.ID, err = uuid.FromString(*in.ID)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	user, err = s.db.CreateUser(ctx, user)
	if err != nil {
		return
	}

	return userToProto(user), nil
}

func (s *grpcServer) DeleteUser(ctx context.Context, in *proto.ProfilesUserIdentifier) (out *empty.Empty, err error) {
	userID, err := uuid.FromString(in.ID)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = s.db.DeleteUser(ctx, userID)
	if err != nil {
		return
	}

	out = new(empty.Empty)
	return
}

func (s *grpcServer) GetUserByID(ctx context.Context, in *proto.ProfilesUserIdentifier) (out *proto.ProfilesUser, err error) {
	userID, err := uuid.FromString(in.ID)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return
	}

	return userToProto(user), nil
}
