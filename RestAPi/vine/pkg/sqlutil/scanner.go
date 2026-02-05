package sqlutil

import (
	"database/sql"
	"fmt"
)

// ScanRow is a helper to scan a single row
func ScanRow(rows *sql.Rows, dest ...interface{}) error {
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	if err := rows.Scan(dest...); err != nil {
		return err
	}

	return rows.Close()
}

// ScanRows scans multiple rows
func ScanRows(rows *sql.Rows, scanFunc func() error) error {
	defer rows.Close()

	for rows.Next() {
		if err := scanFunc(); err != nil {
			return err
		}
	}

	return rows.Err()
}

// MustScan scans or panics (useful in tests)
func MustScan(rows *sql.Rows, dest ...interface{}) {
	if err := ScanRow(rows, dest...); err != nil {
		panic(fmt.Sprintf("scan failed: %v", err))
	}
}
