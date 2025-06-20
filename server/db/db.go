// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

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
	if q.historyAddStmt, err = db.PrepareContext(ctx, historyAdd); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryAdd: %w", err)
	}
	if q.historyGetStmt, err = db.PrepareContext(ctx, historyGet); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryGet: %w", err)
	}
	if q.historyGetExclusiveStmt, err = db.PrepareContext(ctx, historyGetExclusive); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryGetExclusive: %w", err)
	}
	if q.historyGetInclusiveStmt, err = db.PrepareContext(ctx, historyGetInclusive); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryGetInclusive: %w", err)
	}
	if q.historyRemoveStmt, err = db.PrepareContext(ctx, historyRemove); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryRemove: %w", err)
	}
	if q.historyUpdateCategoriesStmt, err = db.PrepareContext(ctx, historyUpdateCategories); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryUpdateCategories: %w", err)
	}
	if q.historyUpdateOwnerStmt, err = db.PrepareContext(ctx, historyUpdateOwner); err != nil {
		return nil, fmt.Errorf("error preparing query HistoryUpdateOwner: %w", err)
	}
	if q.updateHistoryRemoveFileStmt, err = db.PrepareContext(ctx, updateHistoryRemoveFile); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateHistoryRemoveFile: %w", err)
	}
	if q.userCategoryAddStmt, err = db.PrepareContext(ctx, userCategoryAdd); err != nil {
		return nil, fmt.Errorf("error preparing query UserCategoryAdd: %w", err)
	}
	if q.userCategoryRemoveStmt, err = db.PrepareContext(ctx, userCategoryRemove); err != nil {
		return nil, fmt.Errorf("error preparing query UserCategoryRemove: %w", err)
	}
	if q.userGetStmt, err = db.PrepareContext(ctx, userGet); err != nil {
		return nil, fmt.Errorf("error preparing query UserGet: %w", err)
	}
	if q.userUpdateHashStmt, err = db.PrepareContext(ctx, userUpdateHash); err != nil {
		return nil, fmt.Errorf("error preparing query UserUpdateHash: %w", err)
	}
	if q.userUpdateInstagramSessionStmt, err = db.PrepareContext(ctx, userUpdateInstagramSession); err != nil {
		return nil, fmt.Errorf("error preparing query UserUpdateInstagramSession: %w", err)
	}
	if q.userUserStmt, err = db.PrepareContext(ctx, userUser); err != nil {
		return nil, fmt.Errorf("error preparing query UserUser: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.historyAddStmt != nil {
		if cerr := q.historyAddStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyAddStmt: %w", cerr)
		}
	}
	if q.historyGetStmt != nil {
		if cerr := q.historyGetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyGetStmt: %w", cerr)
		}
	}
	if q.historyGetExclusiveStmt != nil {
		if cerr := q.historyGetExclusiveStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyGetExclusiveStmt: %w", cerr)
		}
	}
	if q.historyGetInclusiveStmt != nil {
		if cerr := q.historyGetInclusiveStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyGetInclusiveStmt: %w", cerr)
		}
	}
	if q.historyRemoveStmt != nil {
		if cerr := q.historyRemoveStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyRemoveStmt: %w", cerr)
		}
	}
	if q.historyUpdateCategoriesStmt != nil {
		if cerr := q.historyUpdateCategoriesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyUpdateCategoriesStmt: %w", cerr)
		}
	}
	if q.historyUpdateOwnerStmt != nil {
		if cerr := q.historyUpdateOwnerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing historyUpdateOwnerStmt: %w", cerr)
		}
	}
	if q.updateHistoryRemoveFileStmt != nil {
		if cerr := q.updateHistoryRemoveFileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateHistoryRemoveFileStmt: %w", cerr)
		}
	}
	if q.userCategoryAddStmt != nil {
		if cerr := q.userCategoryAddStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userCategoryAddStmt: %w", cerr)
		}
	}
	if q.userCategoryRemoveStmt != nil {
		if cerr := q.userCategoryRemoveStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userCategoryRemoveStmt: %w", cerr)
		}
	}
	if q.userGetStmt != nil {
		if cerr := q.userGetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userGetStmt: %w", cerr)
		}
	}
	if q.userUpdateHashStmt != nil {
		if cerr := q.userUpdateHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userUpdateHashStmt: %w", cerr)
		}
	}
	if q.userUpdateInstagramSessionStmt != nil {
		if cerr := q.userUpdateInstagramSessionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userUpdateInstagramSessionStmt: %w", cerr)
		}
	}
	if q.userUserStmt != nil {
		if cerr := q.userUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userUserStmt: %w", cerr)
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
	db                             DBTX
	tx                             *sql.Tx
	historyAddStmt                 *sql.Stmt
	historyGetStmt                 *sql.Stmt
	historyGetExclusiveStmt        *sql.Stmt
	historyGetInclusiveStmt        *sql.Stmt
	historyRemoveStmt              *sql.Stmt
	historyUpdateCategoriesStmt    *sql.Stmt
	historyUpdateOwnerStmt         *sql.Stmt
	updateHistoryRemoveFileStmt    *sql.Stmt
	userCategoryAddStmt            *sql.Stmt
	userCategoryRemoveStmt         *sql.Stmt
	userGetStmt                    *sql.Stmt
	userUpdateHashStmt             *sql.Stmt
	userUpdateInstagramSessionStmt *sql.Stmt
	userUserStmt                   *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                             tx,
		tx:                             tx,
		historyAddStmt:                 q.historyAddStmt,
		historyGetStmt:                 q.historyGetStmt,
		historyGetExclusiveStmt:        q.historyGetExclusiveStmt,
		historyGetInclusiveStmt:        q.historyGetInclusiveStmt,
		historyRemoveStmt:              q.historyRemoveStmt,
		historyUpdateCategoriesStmt:    q.historyUpdateCategoriesStmt,
		historyUpdateOwnerStmt:         q.historyUpdateOwnerStmt,
		updateHistoryRemoveFileStmt:    q.updateHistoryRemoveFileStmt,
		userCategoryAddStmt:            q.userCategoryAddStmt,
		userCategoryRemoveStmt:         q.userCategoryRemoveStmt,
		userGetStmt:                    q.userGetStmt,
		userUpdateHashStmt:             q.userUpdateHashStmt,
		userUpdateInstagramSessionStmt: q.userUpdateInstagramSessionStmt,
		userUserStmt:                   q.userUserStmt,
	}
}
