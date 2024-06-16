package repository

import (
	"context"
)

// CreateEstate this function is to store new estate
func (r *Repository) CreateEstate(ctx context.Context, input Estate) (output Estate, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO estates (length, width) VALUES ($1, $2) RETURNING id, width, length, created_at, updated_at",
		input.Length, input.Width,
	).Scan(&output.Id, &output.Width, &output.Length, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return
}

// GetEstateById this function is for get estate by id
func (r *Repository) GetEstateById(ctx context.Context, input GetEstateByIdInput) (output Estate, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, width, length, created_at, updated_at FROM estates WHERE id = $1",
		input.Id,
	).Scan(&output.Id, &output.Width, &output.Length, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return
}

// CreateTree this function is for store tree
func (r *Repository) CreateTree(ctx context.Context, input Tree) (output Tree, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO trees (estate_id, x, y, height) VALUES ($1, $2, $3, $4) RETURNING id, estate_id, x, y, height, created_at, updated_at",
		input.EstateId, input.X, input.Y, input.Height,
	).Scan(&output.Id, &output.EstateId, &output.X, &output.Y, &output.Height, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return
}

// GetTreeByPlot this function is for get tree by plot x and y
func (r *Repository) GetTreeByPlot(ctx context.Context, input GetTreeByPlot) (output Tree, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, estate_id, x, y, height, created_at, updated_at FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3",
		input.EstateId, input.X, input.Y,
	).Scan(&output.Id, &output.EstateId, &output.X, &output.Y, &output.Height, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		return
	}
	return
}

// ListTreesByEstateId this function is for get list trees by estate id
func (r *Repository) ListTreesByEstateId(ctx context.Context, input ListTreesByEstateIdInput) (output []Tree, err error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, estate_id, x, y, height, created_at, updated_at FROM trees WHERE estate_id = $1", input.EstateId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trees []Tree

	// Iterate over the rows
	for rows.Next() {
		var tree Tree
		if err := rows.Scan(&tree.Id, &tree.EstateId, &tree.X, &tree.Y, &tree.Height, &tree.CreatedAt, &tree.UpdatedAt); err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trees, nil
}
