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
	return &SessionRepository{db, userRep, labRep}
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
