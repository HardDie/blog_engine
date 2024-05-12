// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package invite

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.activateStmt, err = db.PrepareContext(ctx, activate); err != nil {
		return nil, fmt.Errorf("error preparing query Activate: %w", err)
	}
	if q.createOrUpdateStmt, err = db.PrepareContext(ctx, createOrUpdate); err != nil {
		return nil, fmt.Errorf("error preparing query CreateOrUpdate: %w", err)
	}
	if q.deleteStmt, err = db.PrepareContext(ctx, delete); err != nil {
		return nil, fmt.Errorf("error preparing query Delete: %w", err)
	}
	if q.getActiveByUserIDStmt, err = db.PrepareContext(ctx, getActiveByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetActiveByUserID: %w", err)
	}
	if q.getAllByUserIDStmt, err = db.PrepareContext(ctx, getAllByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllByUserID: %w", err)
	}
	if q.getByIDStmt, err = db.PrepareContext(ctx, getByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetByID: %w", err)
	}
	if q.getByInviteHashStmt, err = db.PrepareContext(ctx, getByInviteHash); err != nil {
		return nil, fmt.Errorf("error preparing query GetByInviteHash: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.activateStmt != nil {
		if cerr := q.activateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing activateStmt: %w", cerr)
		}
	}
	if q.createOrUpdateStmt != nil {
		if cerr := q.createOrUpdateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createOrUpdateStmt: %w", cerr)
		}
	}
	if q.deleteStmt != nil {
		if cerr := q.deleteStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteStmt: %w", cerr)
		}
	}
	if q.getActiveByUserIDStmt != nil {
		if cerr := q.getActiveByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getActiveByUserIDStmt: %w", cerr)
		}
	}
	if q.getAllByUserIDStmt != nil {
		if cerr := q.getAllByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllByUserIDStmt: %w", cerr)
		}
	}
	if q.getByIDStmt != nil {
		if cerr := q.getByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getByIDStmt: %w", cerr)
		}
	}
	if q.getByInviteHashStmt != nil {
		if cerr := q.getByInviteHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getByInviteHashStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                    DBTX
	tx                    *sql.Tx
	activateStmt          *sql.Stmt
	createOrUpdateStmt    *sql.Stmt
	deleteStmt            *sql.Stmt
	getActiveByUserIDStmt *sql.Stmt
	getAllByUserIDStmt    *sql.Stmt
	getByIDStmt           *sql.Stmt
	getByInviteHashStmt   *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                    tx,
		tx:                    tx,
		activateStmt:          q.activateStmt,
		createOrUpdateStmt:    q.createOrUpdateStmt,
		deleteStmt:            q.deleteStmt,
		getActiveByUserIDStmt: q.getActiveByUserIDStmt,
		getAllByUserIDStmt:    q.getAllByUserIDStmt,
		getByIDStmt:           q.getByIDStmt,
		getByInviteHashStmt:   q.getByInviteHashStmt,
	}
}
