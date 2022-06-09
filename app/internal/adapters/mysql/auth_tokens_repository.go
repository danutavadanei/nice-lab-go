package mysql

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type AuthTokenRepository struct {
	db      *sqlx.DB
	userRep *UserRepository
}

func NewAuthTokenRepository(db *sqlx.DB, userRep *UserRepository) *AuthTokenRepository {
	return &AuthTokenRepository{db: db, userRep: userRep}
}

func (rep AuthTokenRepository) GetUserByAuthToken(ctx context.Context, token string) (user User, err error) {
	var userId uint64
	qry := `SELECT user_id FROM auth_tokens WHERE token = ? AND expire_at > NOW()`
	row := rep.db.QueryRowContext(ctx, qry, token)

	if err = row.Scan(&userId); err != nil {
		return
	}

	return rep.userRep.GetUserById(ctx, userId)
}
