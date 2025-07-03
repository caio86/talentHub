package talenthub

import (
	"context"
	"time"
)

type RHUser struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *RHUser) Validate() error {
	if r.Name == "" {
		return Errorf(EINVALID, "name required")
	}
	if r.Email == "" {
		return Errorf(EINVALID, "email required")
	}
	if r.Password == "" {
		return Errorf(EINVALID, "password required")
	}

	return nil
}

type RHUserService interface {
	FindRHUserByID(ctx context.Context, id int) (*RHUser, error)
	FindRHUsers(ctx context.Context, filter RHUserFilter) ([]*RHUser, int, error)
	CreateRHUser(ctx context.Context, user *RHUser) (*RHUser, error)
	UpdateRHUser(ctx context.Context, id int, upd RHUserUpdate) (*RHUser, error)
	DeleteRHUser(ctx context.Context, id int) error
}

type RHUserFilter struct {
	Offset int32
	Limit  int32
}

type RHUserUpdate struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

