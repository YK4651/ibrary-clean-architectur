package main

import (
	"database/sql"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/query/getbook"
	"github.com/YK4651/ibrary-clean-architectur/internal/http/controllers"
	queryimpl "github.com/YK4651/ibrary-clean-architectur/internal/infrastructure/query"
)

// Dependencies はすべてのアプリケーション依存関係を保持します
type Dependencies struct {
	BookController *controllers.BookController
	// 他のコントローラー/サービスをここに追加
}

// SetupDependencies はすべての依存関係を配線します
// db は context に注入するミドルウェアで使用する（将来用）
func SetupDependencies(db *sql.DB) *Dependencies {
	_ = db // ミドルウェアで context に渡すまで未使用
	// QueryServiceを登録（CQRSクエリサイド）
	bookQueryService := queryimpl.NewBookQueryService()

	// ユースケースを登録
	getBookUseCase := getbook.NewGetBookUseCase(bookQueryService)

	// コントローラーを登録
	bookController := controllers.NewBookController(getBookUseCase)

	return &Dependencies{
		BookController: bookController,
	}
}

// main.goでの使用例:
// db := connectToDatabase()
// deps := SetupDependencies(db)
// setupRoutes(router, deps.BookController)
