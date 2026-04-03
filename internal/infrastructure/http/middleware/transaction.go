package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TransactionMiddleware - HTTPリクエストごとにトランザクション境界を設定
// POST/PUT/PATCH/DELETEは自動的にトランザクション内で実行
// GETは通常の接続を使用（トランザクションなし）
func TransactionMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		ctx := c.Request.Context()

		// 読み取り専用操作（GET）はトランザクション不要
		if method == http.MethodGet {
			// contextに通常のDB接続を設定
			ctx = context.WithValue(ctx, "db", db)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		// 書き込み操作（POST/PUT/PATCH/DELETE）はトランザクション内で実行
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			c.Abort()
			return
		}

		// contextにトランザクションを設定
		ctx = context.WithValue(ctx, "db", tx)
		c.Request = c.Request.WithContext(ctx)

		// ハンドラーを実行
		c.Next()

		// エラーがあればROLLBACK
		if len(c.Errors) > 0 {
			tx.Rollback()
			return
		}

		// 成功時は自動COMMIT
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		}
	}
}
