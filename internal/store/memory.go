package store

import (
	"errors"
	"sync"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

var (
	ErrNotFound      = errors.New("product not found")
	ErrAlreadyExists = errors.New("product already exists")
)

// MemoryStore is a simple in-memory product store (used in exercises 1 & 2).
type MemoryStore struct {
	mu       sync.RWMutex
	products map[int]model.Product
	nextID   int
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		products: make(map[int]model.Product),
		nextID:   1,
	}
}

// GetAll returns all products.
func (s *MemoryStore) GetAll() []model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Product, 0, len(s.products))
	for _, p := range s.products {
		result = append(result, p)
	}
	return result
}

// GetByID returns a product by ID.
func (s *MemoryStore) GetByID(id int) (model.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.products[id]
	if !ok {
		return model.Product{}, ErrNotFound
	}
	return p, nil
}

// Create adds a new product and returns it with the assigned ID.
func (s *MemoryStore) Create(p model.Product) model.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextID
	s.nextID++
	s.products[p.ID] = p
	return p
}

// Update modifies an existing product.
func (s *MemoryStore) Update(id int, p model.Product) (model.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return model.Product{}, ErrNotFound
	}
	p.ID = id
	s.products[id] = p
	return p, nil
}

// Delete removes a product by ID.
func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.products[id]; !ok {
		return ErrNotFound
	}
	delete(s.products, id)
	return nil
}
