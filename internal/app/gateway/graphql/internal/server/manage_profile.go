package server

import (
	"context"

	"github.com/pkg/errors"

	"typerium/internal/app/gateway/graphql/internal/response"
	"typerium/internal/pkg/broker/proto"
	"typerium/internal/pkg/web"
)

func (r *Resolvers) SignUp(ctx context.Context, email string, userName *string, phone *string,
	password string) (*User, error) {
	profilesUser, err := r.ProfilesManagerService.CreateUser(ctx, &proto.NewProfilesUser{
		Email:    email,
		Username: userName,
		Phone:    phone,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	authUser, err := r.AuthService.CreateUser(ctx, &proto.NewAuthUser{
		ID:       &profilesUser.ID,
		Logins:   []string{email, profilesUser.Username},
		Password: password,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out := &User{
		ID:       profilesUser.ID,
		Email:    profilesUser.Email,
		UserName: profilesUser.Username,
		Phone:    profilesUser.Phone,
		TimeMark: &TimeMark{
			CreatedAt: profilesUser.TimeMark.CreatedAt,
			UpdatedAt: authUser.TimeMark.UpdatedAt,
		},
	}
	return out, nil
}

func (r *Resolvers) SignIn(ctx context.Context, login string, password string) (*Tokens, error) {
	session, err := r.AuthService.SignIn(ctx, &proto.AuthCredentials{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out := &Tokens{
		Type:         web.BearerSchema,
		AccessToken:  session.Tokens.AccessToken,
		RefreshToken: session.Tokens.RefreshToken,
	}
	return out, nil
}

func (r *Resolvers) SignOut(ctx context.Context, refreshToken *string) (string, error) {
	credentials := new(proto.AccessCredentials)
	if refreshToken != nil {
		credentials.Token = *refreshToken
	} else {
		var ok bool
		credentials.Token, ok = r.TokenAttr.Get(ctx)
		if !ok {
			return "", response.BadRequestError
		}
	}

	session, err := r.AuthService.SignOut(ctx, credentials)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return session.ID, nil
}

func (r *Resolvers) SignInByAccount(ctx context.Context, accountType LoginAccount) (*Tokens, error) {
	panic("not implemented")
}

func (r *Resolvers) ResetPassword(ctx context.Context, email string) (bool, error) {
	panic("not implemented")
}

func (r *Resolvers) CurrentUser(ctx context.Context) (*User, error) {
	panic("not implemented")
}
