package postgres

import (
	"context"
	"crisplite/internal/domain"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := "postgres://postgres:postgres@localhost:5432/crisplite"
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	t.Cleanup(func() { pool.Close() })
	return pool
}

func cleanupUsers(t *testing.T, pool *pgxpool.Pool, userIDs ...string) {
	t.Helper()
	ctx := context.Background()
	for _, id := range userIDs {
		pool.Exec(ctx, "DELETE FROM contacts WHERE user_id = $1 OR contact_id = $1", id)
		pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	}
}

func TestSave_Success(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{Username: "testuser_save", Password: "hashedpass123"}
	userID, err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID == "" {
		t.Fatal("expected a user ID, got empty string")
	}
	t.Cleanup(func() { cleanupUsers(t, pool, userID) })

	var username string
	err = pool.QueryRow(ctx, "SELECT username FROM users WHERE id = $1", userID).Scan(&username)
	if err != nil {
		t.Fatalf("user not found in db: %v", err)
	}
	if username != "testuser_save" {
		t.Errorf("expected username testuser_save, got %s", username)
	}
}

func TestSave_DuplicateUsername(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{Username: "testuser_dup", Password: "hashedpass123"}
	userID, err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("first save failed: %v", err)
	}
	t.Cleanup(func() { cleanupUsers(t, pool, userID) })

	_, err = repo.Save(ctx, &domain.User{Username: "testuser_dup", Password: "otherpass"})
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAddContact_Success(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	u1, _ := repo.Save(ctx, &domain.User{Username: "contact_user1", Password: "pass"})
	u2, _ := repo.Save(ctx, &domain.User{Username: "contact_user2", Password: "pass"})
	t.Cleanup(func() { cleanupUsers(t, pool, u1, u2) })

	err := repo.AddContact(ctx, u1, u2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var count int
	pool.QueryRow(ctx, "SELECT COUNT(*) FROM contacts WHERE user_id = $1 AND contact_id = $2", u1, u2).Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 contact row, got %d", count)
	}
}

func TestAddContact_DuplicateContact(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	u1, _ := repo.Save(ctx, &domain.User{Username: "dup_contact_u1", Password: "pass"})
	u2, _ := repo.Save(ctx, &domain.User{Username: "dup_contact_u2", Password: "pass"})
	t.Cleanup(func() { cleanupUsers(t, pool, u1, u2) })

	repo.AddContact(ctx, u1, u2)
	err := repo.AddContact(ctx, u1, u2)
	if !errors.Is(err, domain.ErrContactAlreadyExists) {
		t.Errorf("expected ErrContactAlreadyExists, got %v", err)
	}
}

func TestAddContact_UserNotFound(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	fakeUUID1 := "00000000-0000-0000-0000-000000000001"
	fakeUUID2 := "00000000-0000-0000-0000-000000000002"
	err := repo.AddContact(ctx, fakeUUID1, fakeUUID2)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestRemoveContact_Success(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	u1, _ := repo.Save(ctx, &domain.User{Username: "rm_contact_u1", Password: "pass"})
	u2, _ := repo.Save(ctx, &domain.User{Username: "rm_contact_u2", Password: "pass"})
	t.Cleanup(func() { cleanupUsers(t, pool, u1, u2) })

	repo.AddContact(ctx, u1, u2)
	err := repo.RemoveContact(ctx, u1, u2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var count int
	pool.QueryRow(ctx, "SELECT COUNT(*) FROM contacts WHERE user_id = $1 AND contact_id = $2", u1, u2).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 contacts after removal, got %d", count)
	}
}
