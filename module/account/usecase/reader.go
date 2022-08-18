package usecase

import (
	"context"
	"github.com/QuickAmethyst/monosvc/module/account/domain"
	"github.com/QuickAmethyst/monosvc/module/account/repository/sql"
	"github.com/QuickAmethyst/monosvc/stdlibgo/auth"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
)

type Reader interface {
	GetUser(ctx context.Context, stmt sql.UserStatement) (user domain.User, err error)
	SignInWithEmail(ctx context.Context, params SignInWithEmailParams) (accessTokenDetail auth.TokenDetail, refreshTokenDetail auth.TokenDetail, err error)
	RefreshAccessToken(rToken string) (accessTokenDetail auth.TokenDetail, refreshTokenDetail auth.TokenDetail, err error)
}

type reader struct {
	AccountSQL sql.SQL
	Auth       auth.Auth
}

func (r *reader) SignInWithEmail(ctx context.Context, params SignInWithEmailParams) (accessTokenDetail auth.TokenDetail, refreshTokenDetail auth.TokenDetail, err error) {
	user, err := r.AccountSQL.GetUserByEmail(ctx, params.Email)
	errCode := errors.GetCode(err)

	if errCode == sql.EcodeGetUserNotFound {
		err = errors.PropagateWithCode(err, EcodeInvalidCredential, "User with email %+s not found", params.Email)
		return
	} else if err != nil {
		err = errors.PropagateWithCode(err, EcodeSignInWithEmailFailed, "Get user by email failed")
		return
	}

	if err = user.PasswordValid(params.Password); err != nil {
		err = errors.PropagateWithCode(err, EcodeInvalidCredential, "Invalid password")
		return
	}

	accessTokenDetail, refreshTokenDetail, err = r.Auth.CreateTokenPair(user.ID.String())
	if err != nil {
		err = errors.PropagateWithCode(err, EcodeSignInWithEmailFailed, "Create token pair failed")
		return
	}

	return
}

func (r *reader) RefreshAccessToken(rToken string) (accessTokenDetail auth.TokenDetail, refreshTokenDetail auth.TokenDetail, err error) {
	return r.Auth.RefreshAccessToken(rToken)
}

func (r *reader) GetUser(ctx context.Context, stmt sql.UserStatement) (user domain.User, err error) {
	return r.AccountSQL.GetUser(ctx, stmt)
}

func NewReader(opt *Options) Reader {
	return &reader{AccountSQL: opt.AccountSQL, Auth: opt.Auth}
}
