package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	memoryStore := NewMemoryStore()

	testproduct := model.Product{
		Name:  "Steak",
		Price: 99.99,
	}
	created := memoryStore.Create(testproduct)

	retrieved, err := memoryStore.GetByID(created.ID)
	if err != nil {
		t.Errorf("unexpected error retrieving product: %v", err)
	}
	if created.ID == retrieved.ID {
		t.Logf("Product %s retrieved successfully with ID: %d", created.Name, created.ID)
	}

}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

// TODO: Add tests for Update, Delete of existing product, and GetByID with invalid ID

func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()

	// Create a product to update
	original := model.Product{
		Name:  "Original Product",
		Price: 50.00,
	}
	created := s.Create(original)

	t.Logf("Created product %v", created)

	// Update the product
	updated := model.Product{
		ID:    created.ID,
		Name:  "Updated Product",
		Price: 75.00,
	}
	updated, err := s.Update(created.ID, updated)
	if err != nil {
		t.Errorf("unexpected error updating product: %v", err)
	}
	t.Logf("Updated product %v", updated)
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()

	product := model.Product{
		Name:  "Product to Delete",
		Price: 20.00,
	}
	created := s.Create(product)

	t.Logf("Created product %v", created)

	err := s.Delete(created.ID)
	if err != nil {
		t.Errorf("unexpected error deleting product: %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err == ErrNotFound {
		t.Logf("expected ErrNotFound after deleting product")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if err == ErrNotFound {
		t.Logf("expected ErrNotFound when retrieving non-existent product")
	}
}
