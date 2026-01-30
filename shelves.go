package bookstack

import "context"

// ShelvesService handles operations on shelves.
type ShelvesService struct {
	client *Client
}

// List returns a list of shelves with optional filtering.
// TODO: Implement API call to GET /api/shelves
func (s *ShelvesService) List(ctx context.Context, opts *ListOptions) ([]Shelf, error) {
	// Placeholder for future implementation
	return nil, nil
}

// Get retrieves a single shelf by ID.
// TODO: Implement API call to GET /api/shelves/{id}
func (s *ShelvesService) Get(ctx context.Context, id int) (*Shelf, error) {
	// Placeholder for future implementation
	return nil, nil
}

// Create creates a new shelf.
// TODO: Implement API call to POST /api/shelves
func (s *ShelvesService) Create(ctx context.Context, shelf *Shelf) (*Shelf, error) {
	// Placeholder for future implementation
	return nil, nil
}

// Update updates an existing shelf.
// TODO: Implement API call to PUT /api/shelves/{id}
func (s *ShelvesService) Update(ctx context.Context, id int, shelf *Shelf) (*Shelf, error) {
	// Placeholder for future implementation
	return nil, nil
}

// Delete deletes a shelf by ID.
// TODO: Implement API call to DELETE /api/shelves/{id}
func (s *ShelvesService) Delete(ctx context.Context, id int) error {
	// Placeholder for future implementation
	return nil
}
