// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package post

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, arg CreateParams) (*Post, error)
	Edit(ctx context.Context, arg EditParams) (*Post, error)
	GetByID(ctx context.Context, arg GetByIDParams) (*Post, error)
	List(ctx context.Context, arg ListParams) ([]*ListRow, error)
}

var _ Querier = (*Queries)(nil)