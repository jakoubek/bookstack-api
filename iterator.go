package bookstack

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
)

const defaultPageSize = 100

// listAllResponse is used internally to get both data and total from paginated endpoints.
type listAllResponse struct {
	Data  json.RawMessage `json:"data"`
	Total int             `json:"total"`
}

// listAll returns an iterator that paginates through all results for the given path.
func listAll[T any](ctx context.Context, c *Client, path string) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		offset := 0
		for {
			var resp listAllResponse
			url := fmt.Sprintf("%s?count=%d&offset=%d", path, defaultPageSize, offset)
			if err := c.do(ctx, "GET", url, nil, &resp); err != nil {
				var zero T
				yield(zero, err)
				return
			}

			var items []T
			if err := json.Unmarshal(resp.Data, &items); err != nil {
				var zero T
				yield(zero, fmt.Errorf("unmarshaling page data: %w", err))
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			offset += len(items)
			if offset >= resp.Total || len(items) == 0 {
				return
			}
		}
	}
}
