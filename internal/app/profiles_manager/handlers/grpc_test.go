package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"typerium/internal/app/profiles_manager/store"
	"typerium/internal/pkg/broker/proto"
)

//go:generate go run github.com/golang/mock/mockgen -package=handlers -self_package=typerium/internal/app/profiles_manager/handlers -destination=store_mock_test.go typerium/internal/app/profiles_manager/store Store

func TestNewGRPCServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockStore(ctrl)

	server := NewGRPCServer(db)

	assert.Implements(t, (*proto.ProfilesManagerServiceServer)(nil), server)
	assert.IsType(t, &grpcServer{}, server)
	assert.Equal(t, db, server.(*grpcServer).db)
}

func Test_grpcServer_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockStore(ctrl)
	server := NewGRPCServer(db)

	t.Run("user id is set", func(t *testing.T) {
		t.Run("invalid id", func(t *testing.T) {
			id := "test"
			in := &proto.NewProfilesUser{
				ID: &id,
			}
			out, err := server.CreateUser(context.TODO(), in)
			assert.Nil(t, out)
			assert.Error(t, err)
		})

		t.Run("valid id", func(t *testing.T) {
			id := uuid.NewV4().String()
			in := &proto.NewProfilesUser{
				ID: &id,
			}
			grpcServer_CreateUser(t, server, db, in)
		})
	})

	t.Run("user id isn't set", func(t *testing.T) {
		grpcServer_CreateUser(t, server, db, new(proto.NewProfilesUser))
	})
}

func grpcServer_CreateUser(t *testing.T, server proto.ProfilesManagerServiceServer, db *MockStore,
	in *proto.NewProfilesUser) {
	t.Run("creating user in store is failed", func(t *testing.T) {
		ctx := context.TODO()
		db.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, errors.New("test"))
		out, err := server.CreateUser(ctx, in)
		assert.Nil(t, out)
		assert.Error(t, err)
	})

	t.Run("creating user in store is success", func(t *testing.T) {
		ctx := context.TODO()
		db.EXPECT().CreateUser(ctx, gomock.Any()).Return(new(store.User), nil)
		out, err := server.CreateUser(ctx, in)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func Test_grpcServer_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockStore(ctrl)
	server := NewGRPCServer(db)

	t.Run("invalid user id", func(t *testing.T) {
		out, err := server.DeleteUser(context.TODO(), &proto.UserIdentifier{
			ID: "test",
		})
		assert.Nil(t, out)
		assert.Error(t, err)
	})

	id := uuid.NewV4()
	in := &proto.UserIdentifier{
		ID: id.String(),
	}

	t.Run("failed deleting user", func(t *testing.T) {
		ctx := context.TODO()
		db.EXPECT().DeleteUser(ctx, id).Return(errors.New("test"))
		out, err := server.DeleteUser(ctx, in)
		assert.Nil(t, out)
		assert.Error(t, err)
	})

	t.Run("user is deleted", func(t *testing.T) {
		ctx := context.TODO()
		db.EXPECT().DeleteUser(ctx, id).Return(nil)
		out, err := server.DeleteUser(ctx, in)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func Test_grpcServer_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockStore(ctrl)
	server := NewGRPCServer(db)

	t.Run("invalid user id", func(t *testing.T) {
		out, err := server.GetUserByID(context.TODO(), &proto.UserIdentifier{
			ID: "test",
		})
		assert.Nil(t, out)
		assert.Error(t, err)
	})

	id := uuid.NewV4()
	in := &proto.UserIdentifier{
		ID: id.String(),
	}

	t.Run("failed getting user", func(t *testing.T) {
		ctx := context.TODO()
		db.EXPECT().GetUserByID(ctx, id).Return(nil, errors.New("test"))
		out, err := server.GetUserByID(ctx, in)
		assert.Nil(t, out)
		assert.Error(t, err)
	})

	t.Run("return the existing user", func(t *testing.T) {
		ctx := context.TODO()
		db.EXPECT().GetUserByID(ctx, id).Return(new(store.User), nil)
		out, err := server.GetUserByID(ctx, in)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}
