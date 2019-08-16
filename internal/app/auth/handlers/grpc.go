package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"typerium/internal/app/auth/password"
	"typerium/internal/app/auth/signature"
	"typerium/internal/app/auth/store"
	"typerium/internal/pkg/broker/proto"
)

func NewGRPCServer(db store.Store,
	passProcessor password.Processor, jwtSigningMethod jwt.SigningMethod, keyCreator signature.KeyCreator,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) proto.AuthServiceServer {
	return &grpcServer{db: db, passProcessor: passProcessor, jwtSigningMethod: jwtSigningMethod, keyCreator: keyCreator, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL}
}

type grpcServer struct {
	db            store.Store
	passProcessor password.Processor

	jwtSigningMethod jwt.SigningMethod
	keyCreator       signature.KeyCreator
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func userToProto(in *store.User) *proto.AuthUser {
	id := in.ID.String()
	out := &proto.AuthUser{
		ID:    id,
		Login: in.Login,
		TimeMark: proto.TimeMark{
			CreatedAt: in.CreatedAt,
			UpdatedAt: in.UpdatedAt,
		},
	}
	return out
}

func (s *grpcServer) CreateUser(ctx context.Context, in *proto.NewAuthUser) (out *proto.AuthUser, err error) {
	id := uuid.Nil
	if in.ID != nil {
		id, err = uuid.FromString(*in.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	hashPassword, err := s.passProcessor.Encode(in.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := s.db.CreateUser(ctx, id, hashPassword, in.Logins...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return userToProto(user), nil
}

func (s *grpcServer) DeleteUser(ctx context.Context, in *proto.UserIdentifier) (out *empty.Empty, err error) {
	userID, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = s.db.DeleteUser(ctx, userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out = new(empty.Empty)
	return
}

func (s *grpcServer) SignIn(ctx context.Context, in *proto.AuthCredentials) (out *proto.Session, err error) {
	user, err := s.db.GetUserByLogin(ctx, in.Login)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	passValid, err := s.passProcessor.Equal(user.Password, in.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !passValid {
		return nil, status.New(codes.Unauthenticated, "login or password is incorrect").Err()
	}

	secret, err := s.keyCreator.Acquire()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer s.keyCreator.Release(secret)

	session, err := s.db.CreateSession(ctx, user.ID, secret.Bytes(), s.refreshTokenTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if err != nil {
			sessionErr := s.db.DeleteSession(ctx, session.ID)
			if sessionErr != nil {
				// 	todo wrap error
			}
		}
	}()

	return s.createTokens(session, secret)
}

func (s *grpcServer) SignOut(ctx context.Context, in *proto.AccessCredentials) (out *proto.Session, err error) {
	session, err := s.verifyToken(ctx, in.Token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = s.db.DeleteSession(ctx, session.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out = &proto.Session{
		ID: session.ID.String(),
		TimeMark: proto.TimeMark{
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		},
	}
	return
}

func (s *grpcServer) RefreshSession(ctx context.Context, in *proto.AccessCredentials) (out *proto.Session, err error) {
	session, err := s.verifyToken(ctx, in.Token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	secret, err := s.keyCreator.Acquire()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer s.keyCreator.Release(secret)

	session, err = s.db.UpdateSession(ctx, session.ID, secret.Bytes(), s.refreshTokenTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.createTokens(session, secret)
}

const (
	sessionKeyJWT  = "kid"
	contentTypeJWT = "cty"
)

func (s *grpcServer) generateToken(sessionID uuid.UUID, key interface{}, ttl time.Duration) (token string, err error) {
	currentTime := jwt.TimeFunc()
	claims := &jwt.StandardClaims{
		Id:        uuid.NewV4().String(),
		IssuedAt:  currentTime.Unix(),
		ExpiresAt: currentTime.Add(ttl).Unix(),
		NotBefore: currentTime.Add(time.Second).Unix(),
	}
	jwtToken := jwt.NewWithClaims(s.jwtSigningMethod, claims)
	jwtToken.Header[sessionKeyJWT] = sessionID.String()
	jwtToken.Header[contentTypeJWT] = "JWT"

	token, err = jwtToken.SignedString(key)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (s *grpcServer) createTokens(session *store.Session, secret signature.Key) (out *proto.Session, err error) {
	out = &proto.Session{
		ID: session.ID.String(),
		TimeMark: proto.TimeMark{
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		},
	}

	key := secret.Key()

	out.Tokens.AccessToken, err = s.generateToken(session.ID, key, s.accessTokenTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out.Tokens.RefreshToken, err = s.generateToken(session.ID, key, s.refreshTokenTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return
}

func (s *grpcServer) verifyToken(ctx context.Context, token string) (session *store.Session, err error) {
	_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		sessionIDStr, ok := token.Header[sessionKeyJWT].(string)
		if !ok {
			return nil, nil
		}
		sessionID, err := uuid.FromString(sessionIDStr)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		session, err = s.db.GetSessionByID(ctx, sessionID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		secret, err := s.keyCreator.Create(session.KeySignature)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return secret.Key(), nil
	})

	return
}
