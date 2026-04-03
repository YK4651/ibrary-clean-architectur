package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/YK4651/ibrary-clean-architectur/internal/infrastructure/http/middleware"
)

func main() {
	// データベース接続の初期化
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/library_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ginルーターのセットアップ
	r := gin.Default()

	// トランザクションミドルウェアをグローバルに適用
	r.Use(middleware.TransactionMiddleware(db))

	// ルーティング（ハンドラーはStep4で実装）
	// r.POST("/api/books/borrow", handler.BorrowBook)
	// ハンドラー内で ctx.Request.Context().Value("db") からDB/txを取得

	r.Run(":8080")
}
