package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ProductStorage struct {
	ProductId int `db:"product_id"`
	Quantity  int `db:"quantity"`
}

func (ps *ProductStorage) Add(qty int) {
	ps.Quantity += qty
}

func (ps *ProductStorage) Take(qty int) error {
	if ps.Quantity < qty {
		return fmt.Errorf("insufficient quantity")
	}

	ps.Quantity -= qty

	return nil
}

type StorageRepository interface {
	GetByProductId(productId int) *ProductStorage
	Store(ps *ProductStorage)
}

type PgStorageRepository struct {
	db *sqlx.DB
}

func NewPgStorageRepository(db *sqlx.DB) *PgStorageRepository {
	return &PgStorageRepository{
		db: db,
	}
}

func (sr *PgStorageRepository) GetByProductId(productId int) *ProductStorage {
	ps := ProductStorage{}
	err := sr.db.Get(&ps, "select * from storage where product_id = $1", productId)
	if err != nil {
		return &ProductStorage{
			ProductId: productId,
			Quantity:  0,
		}
	}

	return &ps
}

func (sr *PgStorageRepository) Store(ps *ProductStorage) {
	sr.db.Exec("insert into storage (product_id, quantity) values ($1, $2) on conflict (product_id) do update set product_id = $1, quantity = $2", ps.ProductId, ps.Quantity)
}
