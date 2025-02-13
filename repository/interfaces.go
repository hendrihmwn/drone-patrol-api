// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, input Estate) (output Estate, err error)
	GetEstateById(ctx context.Context, input GetEstateByIdInput) (output Estate, err error)
	CreateTree(ctx context.Context, input Tree) (output Tree, err error)
	GetTreeByPlot(ctx context.Context, input GetTreeByPlot) (output Tree, err error)
	ListTreesByEstateId(ctx context.Context, input ListTreesByEstateIdInput) (output []Tree, err error)
}
