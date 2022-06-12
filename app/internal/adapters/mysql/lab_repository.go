package mysql

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type LabType string

const (
	Kali    LabType = "kali"
	Windows LabType = "windows"
)

type Lab struct {
	ID         uint64  `db:"id" json:"id"`
	UUID       string  `db:"uuid" json:"uuid"`
	Name       string  `db:"name" json:"name"`
	Type       LabType `db:"type" json:"type"`
	Hostname   string  `db:"hostname" json:"hostname"`
	InstanceID string  `db:"instance_id" json:"instance_id"`
	Available  bool    `db:"available" json:"available"`
}

type LabRepository struct {
	db *sqlx.DB
}

func NewLabRepository(db *sqlx.DB) *LabRepository {
	return &LabRepository{db: db}
}

func (rep LabRepository) ListLabs(ctx context.Context) (result []Lab, err error) {
	qry := `SELECT * FROM labs ORDER BY id ASC`
	err = rep.db.SelectContext(ctx, &result, qry)

	return
}

func (rep LabRepository) GetLabById(ctx context.Context, id uint64) (lab Lab, err error) {
	qry := `SELECT * FROM labs WHERE id = ?`
	row := rep.db.QueryRowxContext(ctx, qry, id)

	err = row.StructScan(&lab)

	return
}
