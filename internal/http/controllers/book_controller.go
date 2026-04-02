package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/query/getbook"
)

// BookController - 書籍APIコントローラー（プレゼンテーション層）
type BookController struct {
	getBookUseCase *getbook.GetBookUseCase // ユースケース（依存性注入）
}

// NewBookController - 新しいBookControllerを作成
func NewBookController(getBookUseCase *getbook.GetBookUseCase) *BookController {
	return &BookController{
		getBookUseCase: getBookUseCase,
	}
}

// GetBook - GET /api/books/{bookID} を処理
func (c *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	// 1. URLパスからbookIDを抽出（Go 1.22+ 標準ライブラリ）
	bookID := r.PathValue("bookID")

	// 2. リクエストDTOを作成
	request, err := getbook.NewGetBookRequest(bookID)
	if err != nil {
		// リクエストDTOのバリデーション失敗
		c.sendJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 3. ユースケースを実行
	response := c.getBookUseCase.Execute(r.Context(), request)

	// 4. 失敗レスポンスを処理
	if !response.Success {
		c.handleError(w, *response.ErrorMessage)
		return
	}

	// 5. 成功レスポンスを返す
	c.sendJSON(w, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":          response.Data.ID,
			"isbn":        response.Data.ISBN,
			"title":       response.Data.Title,
			"author":      response.Data.Author,
			"isAvailable": response.Data.IsAvailable,
			"currentLoan": response.Data.CurrentLoan,
		},
	}, http.StatusOK) // 200 OK
}

// handleError - エラーメッセージを適切なHTTPステータスコードにマッピング
func (c *BookController) handleError(w http.ResponseWriter, errorMessage string) {
	statusCode := http.StatusInternalServerError // デフォルト: 500

	if strings.Contains(strings.ToLower(errorMessage), "not found") {
		statusCode = http.StatusNotFound // 404 Not Found
	} else if strings.Contains(errorMessage, "Invalid") {
		statusCode = http.StatusBadRequest // 400 Bad Request
	}

	c.sendJSON(w, map[string]interface{}{
		"success": false,
		"error":   errorMessage,
	}, statusCode)
}

// sendJSON - JSONレスポンスを送信するヘルパー関数
func (c *BookController) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
