package mysql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	Student   UserType = "student"
	Professor UserType = "professor"
)

type User struct {
	ID    uint64   `db:"id" json:"id"`
	UUID  string   `db:"uuid" json:"uuid"`
	Name  string   `db:"name" json:"name"`
	Email string   `db:"email" json:"email"`
	Type  UserType `db:"type" json:"type"`
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (rep UserRepository) ListUsers(ctx context.Context) (result []User, err error) {
	qry := `SELECT id,uuid,name,email,type FROM users ORDER BY id ASC`
	err = rep.db.SelectContext(ctx, &result, qry)

	return
}

func (rep UserRepository) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	qry := `SELECT id,uuid,name,email,type FROM users WHERE email = ?`
	row := rep.db.QueryRowxContext(ctx, qry, email)

	err = row.StructScan(&user)
	return
}

func (rep UserRepository) GetUserById(ctx context.Context, id uint64) (user User, err error) {
	qry := `SELECT id,uuid,name,email,type FROM users WHERE id = ?`
	row := rep.db.QueryRowxContext(ctx, qry, id)

	err = row.StructScan(&user)
	return
}

func (rep UserRepository) CheckUserPassword(ctx context.Context, email string, password string) (err error) {
	var hashedPassword string
	qry := `SELECT password FROM users WHERE email = ?`
	row := rep.db.QueryRowContext(ctx, qry, email)

	if err = row.Scan(&hashedPassword); err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return
	}

	return nil
}
