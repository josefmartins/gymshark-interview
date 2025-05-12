package storage

import (
	"context"
	"database/sql"
	"errors"
	"gymshark-interview/internal/model"
	"log"

	sqlite "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	ErrFailedToCreateProduct = errors.New("failed to create product")
	ErrFailedToDeleteProduct = errors.New("failed to delete product")
	ErrFailedToListProducts  = errors.New("failed to list products")
	ErrFailedToGetProduct    = errors.New("failed to get product")
	ErrProductNotFound       = errors.New("product not found")
	ErrConstraintViolation   = errors.New("database constraint violation")
)

func (s *Storage) GetProductWithPackageSizes(ctx context.Context, productID string) (*model.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rows, err := s.db.QueryxContext(ctx, `
		SELECT p.id AS product_id, p.name, pkg.size FROM products p 
		LEFT JOIN package_sizes pkg ON pkg.product_id = p.id WHERE p.id = ?
	`, productID)
	if err != nil {
		log.Printf("failed to get product with package sizes in DB: %v", err)
		return nil, ErrFailedToGetProduct
	}
	defer rows.Close()

	var prod *model.Product

	for rows.Next() {
		var (
			pID, pName string
			pkgSize    sql.NullInt64
		)

		if err := rows.Scan(&pID, &pName, &pkgSize); err != nil {
			log.Printf("failed to scan row: %v", err)
			return nil, ErrFailedToGetProduct
		}

		// only register once
		if prod == nil {
			prod = &model.Product{
				ID:   pID,
				Name: pName,
			}

		}

		if pkgSize.Valid {
			prod.PackageSizes = append(prod.PackageSizes, int(pkgSize.Int64))
		}
	}

	if prod == nil {
		return nil, ErrProductNotFound
	}

	return prod, nil
}

func handleCreateProductError(tx *sql.Tx, err error) error {
	txErr := tx.Rollback()
	if txErr != nil {
		err = errors.Join(err, txErr)
	}
	log.Printf("failed to create product in DB: %v", err)
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) {
		if sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return ErrConstraintViolation
		}
	}
	return ErrFailedToCreateProduct
}

func (s *Storage) CreateProduct(ctx context.Context, product model.Product) (*model.Product, error) {
	id, _ := uuid.NewV7()
	var res model.Product
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, handleCreateProductError(tx, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO products (id,name) VALUES (?,?)",
		id.String(), product.Name)
	if err != nil {
		return nil, handleCreateProductError(tx, err)
	}
	if len(product.PackageSizes) != 0 {
		sizes, err := s.createPackageSizes(ctx, tx, id.String(), product.PackageSizes)
		if err != nil {
			return nil, handleCreateProductError(tx, err)
		}
		res.PackageSizes = sizes
	}

	err = tx.Commit()
	if err != nil {
		return nil, handleCreateProductError(tx, err)
	}

	return &res, nil
}

func (s *Storage) ListProducts(ctx context.Context) ([]model.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rows, err := s.db.QueryxContext(ctx, `SELECT p.id AS product_id, p.name, pkg.size FROM products p 
		LEFT JOIN package_sizes pkg ON pkg.product_id = p.id`)
	if err != nil {
		log.Printf("failed to list products in DB: %v", err)
		return nil, ErrFailedToListProducts
	}

	defer rows.Close()

	products := make(map[string]*model.Product, 0)
	for rows.Next() {
		var (
			pID, pName string
			pkgSize    sql.NullInt64
		)

		if err := rows.Scan(&pID, &pName, &pkgSize); err != nil {
			log.Printf("failed to scan row: %v", err)
			return nil, ErrFailedToGetProduct
		}

		// only register once
		if _, exists := products[pID]; !exists {
			products[pID] = &model.Product{
				ID:   pID,
				Name: pName,
			}
		}

		if pkgSize.Valid {
			products[pID].PackageSizes = append(products[pID].PackageSizes, int(pkgSize.Int64))
		}
	}

	res := make([]model.Product, 0)
	for _, p := range products {
		res = append(res, *p)
	}
	return res, nil
}

func (s *Storage) DeleteProduct(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)
	if err != nil {
		log.Printf("failed to delete product from DB: %v", err)
		return ErrFailedToDeleteProduct
	}
	return nil
}
