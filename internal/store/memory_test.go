package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Widget", Price: 9.99})

	retrieved, err := s.GetByID(p.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if retrieved.Name != p.Name || retrieved.Price != p.Price {
		t.Errorf("expected %v, got %v", p, retrieved)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestUpdateExisting(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Widget", Price: 9.99})

	updated, err := s.Update(p.ID, model.Product{Name: "Updated Widget", Price: 19.99})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if updated.Name != "Updated Widget" || updated.Price != 19.99 {
		t.Errorf("expected updated product, got %v", updated)
	}

	retrieved, _ := s.GetByID(p.ID)
	if retrieved.Name != "Updated Widget" {
		t.Errorf("expected updated product in store, got %v", retrieved)
	}
}

func TestUpdateNonExistent(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Update(999, model.Product{Name: "Widget", Price: 9.99})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExisting(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Widget", Price: 9.99})

	err := s.Delete(p.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = s.GetByID(p.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestGetAllWithMultipleProducts(t *testing.T) {
	s := NewMemoryStore()
	s.Create(model.Product{Name: "Widget", Price: 9.99})
	s.Create(model.Product{Name: "Gadget", Price: 19.99})
	s.Create(model.Product{Name: "Doohickey", Price: 29.99})

	products := s.GetAll()
	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}
}
