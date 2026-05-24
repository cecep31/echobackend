package service

import (
	"context"
	"errors"
	"testing"

	apperrors "echobackend/internal/errors"
	"echobackend/internal/model"
)

func ptr[T any](v T) *T { return &v }

func TestUserService_GetByID_Success(t *testing.T) {
	first := "Alice"
	last := "Smith"
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return &model.User{ID: id, Email: "a@b.com", FirstName: &first, LastName: &last}, nil
		},
	}
	svc := NewUserService(repo)
	resp, err := svc.GetByID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "user-1" || resp.Email != "a@b.com" {
		t.Errorf("unexpected response %+v", resp)
	}
	if resp.Name != "Alice Smith" {
		t.Errorf("Name = %q, want %q", resp.Name, "Alice Smith")
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, apperrors.ErrUserNotFound
		},
	}
	svc := NewUserService(repo)
	_, err := svc.GetByID(context.Background(), "missing")
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetMe_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return &model.User{ID: id, Email: "me@example.com", IsSuperAdmin: ptr(true)}, nil
		},
	}
	svc := NewUserService(repo)
	resp, err := svc.GetMe(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.IsSuperAdmin == nil || !*resp.IsSuperAdmin {
		t.Fatalf("expected IsSuperAdmin true, got %v", resp.IsSuperAdmin)
	}
}

func TestUserService_GetByUsername_Success(t *testing.T) {
	uname := "alice"
	repo := &mockUserRepo{
		getByUsernameFn: func(ctx context.Context, username string) (*model.User, error) {
			return &model.User{ID: "1", Email: "alice@x.com", Username: &uname}, nil
		},
	}
	svc := NewUserService(repo)
	resp, err := svc.GetByUsername(context.Background(), uname)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Username == nil || *resp.Username != uname {
		t.Fatalf("got %v", resp.Username)
	}
}

func TestUserService_GetUsers_Success(t *testing.T) {
	repo := &mockUserRepo{
		getUsersFn: func(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
			return []*model.User{
				{ID: "u1", Email: "a@x.com"},
				nil, // nil entries should be skipped per service contract
				{ID: "u2", Email: "b@x.com"},
			}, 2, nil
		},
	}
	svc := NewUserService(repo)
	resp, total, err := svc.GetUsers(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(resp) != 2 {
		t.Fatalf("len = %d, want 2", len(resp))
	}
}

func TestUserService_GetUsers_RepoError(t *testing.T) {
	wantErr := errors.New("db down")
	repo := &mockUserRepo{
		getUsersFn: func(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
			return nil, 0, wantErr
		},
	}
	svc := NewUserService(repo)
	_, _, err := svc.GetUsers(context.Background(), 0, 10)
	if err == nil || !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestUserService_Delete(t *testing.T) {
	deleted := false
	repo := &mockUserRepo{
		softDeleteFn: func(ctx context.Context, id string) error {
			if id != "u1" {
				t.Errorf("unexpected id %q", id)
			}
			deleted = true
			return nil
		},
	}
	svc := NewUserService(repo)
	if err := svc.Delete(context.Background(), "u1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Fatal("expected SoftDeleteByID to be called")
	}
}
