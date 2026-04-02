package controllers_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/query/getbook"
	"github.com/YK4651/ibrary-clean-architectur/internal/http/controllers"
	queryimpl "github.com/YK4651/ibrary-clean-architectur/internal/infrastructure/query"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetBookAPI_Returns200ForExistingBook(t *testing.T) {
	// Arrange - テストデータベースをセットアップ
	db := setupTestDatabase(t)
	defer db.Close()

	// テスト用の書籍データをデータベースに挿入
	_, err := db.Exec(`
		INSERT INTO books (id, isbn, title, author)
		VALUES ('b-12345', '978-0-13-468599-1', 'Clean Architecture', 'Robert C. Martin')
	`)
	if err != nil {
		t.Fatal(err)
	}

	// 依存関係とコントローラをセットアップ
	deps := setupDependencies(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/books/{bookID}", deps.BookController.GetBook)

	// Act - HTTPリクエストを送信
	req := httptest.NewRequest("GET", "/api/books/b-12345", nil)
	// リクエストのcontextにDBを注入
	req = req.WithContext(context.WithValue(req.Context(), "db", db))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Assert（検証）
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["success"] != true {
		t.Error("Expected success to be true")
	}

	data := response["data"].(map[string]interface{})
	if data["id"] != "b-12345" {
		t.Errorf("Expected book ID 'b-12345', got %v", data["id"])
	}
	if data["title"] != "Clean Architecture" {
		t.Errorf("Expected title 'Clean Architecture', got %v", data["title"])
	}
}

func TestGetBookAPI_Returns404ForNonExistentBook(t *testing.T) {
	// Arrange（準備）
	db := setupTestDatabase(t)
	defer db.Close()

	deps := setupDependencies(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/books/{bookID}", deps.BookController.GetBook)

	// Act（実行）
	req := httptest.NewRequest("GET", "/api/books/b-99999", nil)
	// リクエストのcontextにDBを注入
	req = req.WithContext(context.WithValue(req.Context(), "db", db))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Assert（検証）
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	if response["success"] != false {
		t.Error("Expected success to be false")
	}
	if response["error"] != "Book not found" {
		t.Errorf("Expected error 'Book not found', got %v", response["error"])
	}
}

// Dependencies はテストの依存関係を保持します
type Dependencies struct {
	BookController *controllers.BookController
}

// setupDependencies はテスト用に依存関係を構築します
func setupDependencies(db *sql.DB) *Dependencies {
	bookQueryService := queryimpl.NewBookQueryService()
	getBookUseCase := getbook.NewGetBookUseCase(bookQueryService)
	bookController := controllers.NewBookController(getBookUseCase)

	return &Dependencies{
		BookController: bookController,
	}
}

func setupTestDatabase(t *testing.T) *sql.DB {
	// MySQLテストデータベースに接続
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/library")
	if err != nil {
		t.Fatal(err)
	}

	// 各テスト前にテーブルをクリーンアップ
	cleanupTables(t, db)

	return db
}

func cleanupTables(t *testing.T, db *sql.DB) {
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	tables := []string{"loans", "books", "users"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
