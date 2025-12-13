// Package sqlsetpgx provides helper functions for using sqlset with pgx.
package sqlsetpgx

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theprogrammer67/sqlset"
)

// DBHelper is a helper struct that holds a pgxpool.Pool and a sqlset.SQLSet.
type DBHelper struct {
	pool    *pgxpool.Pool
	sqlset  *sqlset.SQLSet
	scanAPI *pgxscan.API
}

// NewDBHelper creates a new DBHelper.
func NewDBHelper(pool *pgxpool.Pool, sqlset *sqlset.SQLSet) *DBHelper {
	return &DBHelper{
		pool:    pool,
		sqlset:  sqlset,
		scanAPI: mustNewAPI(mustNewDBScanAPI(dbscan.WithAllowUnknownColumns(true))),
	}
}

// Get retrieves a query from sqlset and scans exactly one row into dest.
func (h *DBHelper) Get(ctx context.Context, dest any, setID, queryID string, args ...any) error {
	sql, err := h.sqlset.Get(setID, queryID)
	if err != nil {
		return err
	}

	return h.scanAPI.Get(ctx, h.pool, dest, sql, args...)
}

// Select retrieves a query from sqlset and scans all rows into a slice.
func (h *DBHelper) Select(ctx context.Context, dest any, setID, queryID string, args ...any) error {
	sql, err := h.sqlset.Get(setID, queryID)
	if err != nil {
		return err
	}

	return h.scanAPI.Select(ctx, h.pool, dest, sql, args...)
}

// Exec retrieves a query from sqlset and executes it, returning the number of affected rows.
// It should be used for queries that do not return rows (e.g. INSERT, UPDATE, DELETE).
func (h *DBHelper) Exec(ctx context.Context, setID, queryID string, args ...any) (int64, error) {
	sql, err := h.sqlset.Get(setID, queryID)
	if err != nil {
		return 0, err
	}

	tag, err := h.pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("exec '%s.%s' failed: %w", setID, queryID, err)
	}

	return tag.RowsAffected(), nil
}

// pivate

func mustNewDBScanAPI(opts ...dbscan.APIOption) *dbscan.API {
	api, err := pgxscan.NewDBScanAPI(opts...)
	if err != nil {
		panic(err)
	}

	return api
}

func mustNewAPI(dbscanAPI *dbscan.API) *pgxscan.API {
	api, err := pgxscan.NewAPI(dbscanAPI)
	if err != nil {
		panic(err)
	}

	return api
}
