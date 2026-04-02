package main

import (
	"net/http"

	"github.com/YK4651/ibrary-clean-architectur/internal/http/controllers"
)

func setupRoutes(mux *http.ServeMux, controller *controllers.BookController) {
	// パスパラメータを持つルート（Go 1.22+）
	mux.HandleFunc("GET /api/books/{bookID}", controller.GetBook)
}
