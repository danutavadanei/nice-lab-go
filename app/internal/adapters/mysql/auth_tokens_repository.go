package mysql

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (rep AuthTokenRepository) ListTokenUsers(ctx context.Context) (map[string]User, error) {
	tokenUsers := make(map[string]User)
	qry := `SELECT token, user_id FROM auth_tokens WHERE expire_at > NOW()`

	var tokens []struct {
		Token  string `db:"token"`
		UserID uint64 `db:"user_id"`
	}

	if err := rep.db.SelectContext(ctx, &tokens, qry); err != nil {
		return nil, err
	}

	for _, token := range tokens {
		user, _ := rep.userRep.GetUserById(ctx, token.UserID)
		tokenUsers[token.Token] = user
	}

	return tokenUsers, nil
}

func (rep AuthTokenRepository) NewTokenForUserId(ctx context.Context, id uint64) (token string, err error) {
	tx, err := rep.db.Beginx()

	qry := `DELETE FROM auth_tokens WHERE user_id = ?`
	if _, err = rep.db.ExecContext(ctx, qry, id); err != nil {
		return
	}

	qry = `INSERT INTO auth_tokens (user_id, token, expire_at) VALUES (?, ?, ?)`
	qry = tx.Rebind(qry)

	var u uuid.UUID
	if u, err = uuid.NewUUID(); err != nil {
		return
	}

	token = u.String()

	if _, err = tx.ExecContext(ctx, qry, id, token, time.Now().Add(24*time.Hour)); err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return
}
