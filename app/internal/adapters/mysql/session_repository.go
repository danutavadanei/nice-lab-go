package mysql

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type dbSession struct {
	ID     uint64 `db:"id"`
	UserID uint64 `db:"user_id"`
	LabID  uint64 `db:"lab_id"`
}

type Session struct {
	ID   uint64 `json:"id"`
	User User   `json:"user"`
	Lab  Lab    `json:"lab"`
}

type SessionRepository struct {
	db      *sqlx.DB
	userRep *UserRepository
	labRep  *LabRepository
}

func NewSessionRepository(db *sqlx.DB, userRep *UserRepository, labRep *LabRepository) *SessionRepository {
	return &SessionRepository{
		db:      db,
		userRep: userRep,
		labRep:  labRep,
	}
}

func (rep SessionRepository) ListSessions(ctx context.Context) (sessions []Session, err error) {
	rows := make([]dbSession, 0)
	qry := `SELECT * FROM sessions ORDER BY id DESC`

	if err = rep.db.SelectContext(ctx, &rows, qry); err != nil {
		return
	}

	for _, row := range rows {
		var user User
		var lab Lab

		if user, err = rep.userRep.GetUserById(ctx, row.UserID); err != nil {
			return nil, err
		}
		if lab, err = rep.labRep.GetLabById(ctx, row.LabID); err != nil {
			return nil, err
		}

		sessions = append(sessions, Session{
			ID:   row.ID,
			User: user,
			Lab:  lab,
		})
	}

	return
}

func (rep SessionRepository) GetSessionById(ctx context.Context, id uint64) (session Session, err error) {
	var dbSes dbSession
	var user User
	var lab Lab

	qry := `SELECT * FROM sessions WHERE id = ?`
	row := rep.db.QueryRowxContext(ctx, qry, id)

	err = row.StructScan(&dbSes)

	if user, err = rep.userRep.GetUserById(ctx, dbSes.UserID); err != nil {
		return Session{}, err
	}
	if lab, err = rep.labRep.GetLabById(ctx, dbSes.LabID); err != nil {
		return Session{}, err
	}

	return Session{
		ID:   dbSes.ID,
		User: user,
		Lab:  lab,
	}, nil
}

func (rep SessionRepository) CreateSession(ctx context.Context, user User, lab Lab) (Session, error) {
	tx, err := rep.db.Beginx()

	qry := `INSERT INTO sessions (user_id, lab_id) VALUES (?, ?)`
	qry = tx.Rebind(qry)

	res, err := tx.ExecContext(ctx, qry, user.ID, lab.ID)

	if err != nil {
		return Session{}, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return Session{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Session{}, err
	}

	return rep.GetSessionById(ctx, uint64(id))
}
