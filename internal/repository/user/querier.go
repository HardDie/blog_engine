// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package user

import (
	"context"
)

type Querier interface {
	Create(ctx context.Context, arg CreateParams) (*User, error)
	GetByIDPrivate(ctx context.Context, id int64) (*User, error)
	GetByIDPublic(ctx context.Context, id int64) (*GetByIDPublicRow, error)
	GetByName(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, arg UpdateParams) (*User, error)
}

var _ Querier = (*Queries)(nil)
