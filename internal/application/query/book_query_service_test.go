package query_test

import (
	"context"
	"database/sql"
	"testing"

	queryimpl "github.com/YK4651/ibrary-clean-architectur/internal/infrastructure/query"
	_ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/library")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// テスト前にテーブルをクリーンアップ
	cleanupTables(t, db)

	return db
}

func cleanupTables(t *testing.T, db *sql.DB) {
	// 外部キー制約を一時的に無効化
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// テーブルをトランケート
	tables := []string{"loans", "books", "users"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}

	// 外部キー制約を再度有効化
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func TestBookQueryService_GetBookByID_FindsAvailableBook(t *testing.T) {
	// Arrange（準備）
	db := setupTestDB(t)
	defer db.Close()

	// テストデータを挿入
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-12345', '978-0-123456-78-9', 'Clean Architecture',
				'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	queryService := queryimpl.NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	// Act（実行）
	book, err := queryService.GetBookByID(ctx, "b-12345")

	// Assert（検証）
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if book == nil {
		t.Fatal("Expected book to be found, got nil")
	}
	if book.ID != "b-12345" {
		t.Errorf("Expected ID 'b-12345', got: %s", book.ID)
	}
	if book.Title != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got: %s", book.Title)
	}
	if !book.IsAvailable {
		t.Error("Expected book to be available")
	}
	if book.CurrentLoan != nil {
		t.Error("Expected no current loan for available book")
	}
}

func TestBookQueryService_GetBookByID_FindsBorrowedBook(t *testing.T) {
	// Arrange（準備）
	db := setupTestDB(t)
	defer db.Close()

	// テストデータと貸出を挿入
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author, created_at)
		VALUES ('b-12345', '978-0-123456-78-9', 'Clean Architecture',
				'Robert C. Martin', '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}

	// Insert test user (required for foreign key)
	_, err = db.Exec(`
		INSERT INTO users (id, name, email, status, created_at)
		VALUES ('u-001', 'Test User', 'test@example.com', 1, '2024-01-01 00:00:00')
	`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO loans (id, book_id, user_id, borrowed_at, due_date, returned_at)
		VALUES ('l-001', 'b-12345', 'u-001', '2024-01-15 10:00:00',
				'2024-02-15 10:00:00', NULL)
	`)
	if err != nil {
		t.Fatalf("Failed to insert loan: %v", err)
	}

	queryService := queryimpl.NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	// Act（実行）
	book, err := queryService.GetBookByID(ctx, "b-12345")

	// Assert（検証）
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if book == nil {
		t.Fatal("Expected book to be found, got nil")
	}
	if book.IsAvailable {
		t.Error("Expected book to be borrowed (not available)")
	}
	if book.CurrentLoan == nil {
		t.Fatal("Expected current loan to exist")
	}
	if book.CurrentLoan.LoanID != "l-001" {
		t.Errorf("Expected loan ID 'l-001', got: %s", book.CurrentLoan.LoanID)
	}
}

func TestBookQueryService_GetBookByID_ReturnsNilWhenNotFound(t *testing.T) {
	// Arrange（準備）
	db := setupTestDB(t)
	defer db.Close()

	queryService := queryimpl.NewBookQueryService()
	ctx := context.WithValue(context.Background(), "db", db)

	// Act（実行）
	book, err := queryService.GetBookByID(ctx, "b-99999")

	// Assert（検証）
	if err != nil {
		t.Errorf("Expected no error for not found case, got: %v", err)
	}
	if book != nil {
		t.Error("Expected nil for non-existent book")
	}
}
