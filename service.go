package storage

import (
	"sync"
)

type StorageService interface {
	Add(productId int, qty int) int
	Take(productId int, qty int) (int, error)
	Get(productId int) int
}

type MStorageService struct {
	m    sync.Map
	repo StorageRepository
}

func NewMStorageService(repo StorageRepository) *MStorageService {
	return &MStorageService{
		repo: repo,
	}
}

func (ss *MStorageService) Add(productId int, qty int) int {
	// блокировка записи по продукту
	// можно реализовать на уровне бд через select for update
	m, _ := ss.m.LoadOrStore(productId, &sync.Mutex{})
	m.(*sync.Mutex).Lock()
	defer m.(*sync.Mutex).Unlock()

	// проверка на наличие продукта по id опущена
	ps := ss.repo.GetByProductId(productId)
	ps.Add(qty)
	ss.repo.Store(ps)

	ss.m.Delete(productId)

	return ps.Quantity
}

func (ss *MStorageService) Take(productId int, qty int) (int, error) {
	// блокировка записи по продукту
	// можно реализовать на уровне бд через select for update
	m, _ := ss.m.LoadOrStore(productId, &sync.Mutex{})
	m.(*sync.Mutex).Lock()
	defer m.(*sync.Mutex).Unlock()

	ps := ss.repo.GetByProductId(productId)

	err := ps.Take(qty)
	if err != nil {
		return 0, err
	}

	ss.repo.Store(ps)

	ss.m.Delete(productId)
	return ps.Quantity, nil
}

func (ss *MStorageService) Get(productId int) int {
	ps := ss.repo.GetByProductId(productId)
	return ps.Quantity
}
