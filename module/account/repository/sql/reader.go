package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/QuickAmethyst/monosvc/module/account/domain"
	"github.com/QuickAmethyst/monosvc/stdlibgo/errors"
	qb "github.com/QuickAmethyst/monosvc/stdlibgo/querybuilder/sql"
	sdkSql "github.com/QuickAmethyst/monosvc/stdlibgo/sql"
	"github.com/google/uuid"
)

type Reader interface {
	GetUser(ctx context.Context, stmt UserStatement) (user domain.User, err error)
	GetUserByEmail(ctx context.Context, email string) (user domain.User, err error)
	GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (user domain.UserProfile, err error)
}

type reader struct {
	db sdkSql.DB
}

func (r *reader) GetUserByEmail(ctx context.Context, email string) (user domain.User, err error) {
	return r.GetUser(ctx, UserStatement{Email: email})
}

func (r *reader) GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (profile domain.UserProfile, err error) {
	whereClause, whereArgs, err := qb.NewWhereClause(UserProfileStatement{UserID: userID})
	query := fmt.Sprintf("SELECT id, full_name, created_at, updated_at FROM user_profiles %s", whereClause)

	if err = r.db.QueryRowContext(ctx, r.db.Rebind(query), whereArgs...).Scan(&profile); err != nil {
		err = errors.PropagateWithCode(err, EcodeGetUserProfileFailed, "Failed to get user profile")
		return
	}

	return
}

func (r *reader) GetUser(ctx context.Context, stmt UserStatement) (user domain.User, err error) {
	whereClause, whereArgs, err := qb.NewWhereClause(stmt)
	query := fmt.Sprintf("SELECT id, email, password, created_at, updated_at FROM users %s", whereClause)

	err = r.db.QueryRowContext(ctx, r.db.Rebind(query), whereArgs...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.PropagateWithCode(err, EcodeGetUserNotFound, "User not found")
			return
		}

		err = errors.PropagateWithCode(err, EcodeGetUserFailed, "Failed to get user")
		return
	}

	return
}

func NewReader(opt *Options) Reader {
	return &reader{db: opt.SlaveDB}
}
