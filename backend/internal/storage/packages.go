package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sqlite "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	ErrFailedToCreatePackageSize = errors.New("failed to create package size")
	ErrFailedToDeletePackageSize = errors.New("failed to delete package size")
)

func (s *Storage) AddPackageSize(ctx context.Context, productID string, size int) error {
	id, _ := uuid.NewV7()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.db.ExecContext(ctx, "INSERT INTO package_sizes (id,product_id,size) VALUES (?,?,?)",
		id, productID, size)
	if err != nil {
		log.Printf("failed to create package size in DB: %v", err)
		var sqliteError *sqlite.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return ErrConstraintViolation
			}
		}
		return ErrFailedToCreatePackageSize
	}
	return nil
}

func (s *Storage) RemovePackageSize(ctx context.Context, productID string, size int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.db.ExecContext(ctx, "DELETE FROM package_sizes WHERE product_id=? AND size=?", productID, size)
	if err != nil {
		log.Printf("failed to delete package size from DB: %v", err)
		var sqliteError *sqlite.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return ErrConstraintViolation
			}
		}
		return ErrFailedToDeletePackageSize
	}
	return nil
}

func (s *Storage) createPackageSizes(ctx context.Context, tx *sql.Tx, productID string, sizes []int) ([]int, error) {
	command := "INSERT INTO package_sizes (id,product_id,size) VALUES"
	args := []interface{}{}
	for _, size := range sizes {
		command += " (?,?,?),"
		id, _ := uuid.NewV7()
		args = append(args, id.String(), productID, size)
	}
	// remove last comma
	command = command[:len(command)-1]
	_, err := tx.ExecContext(ctx, command, args...)
	if err != nil {
		log.Printf("failed to create package size in DB: %v", err)
		return nil, err
	}
	return sizes, nil
}
