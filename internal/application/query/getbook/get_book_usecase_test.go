package getbook_test

import (
	"context"
	"testing"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/query"
	"github.com/YK4651/ibrary-clean-architectur/internal/application/query/getbook"
	"github.com/YK4651/ibrary-clean-architectur/internal/application/query/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetBookUseCase_ReturnsBookWhenAvailable(t *testing.T) {
	// Arrange（準備）
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockBookQueryService(ctrl)
	useCase := getbook.NewGetBookUseCase(mockQueryService)

	ctx := context.Background()
	bookID := "b-12345"

	// モックの期待値を設定
	mockQueryService.EXPECT().
		GetBookByID(ctx, bookID).
		Return(&query.BookReadModel{
			ID:          "b-12345",
			ISBN:        "978-0-123456-78-9",
			Title:       "Clean Architecture",
			Author:      "Robert C. Martin",
			IsAvailable: true,
			CurrentLoan: nil,
		}, nil)

	request, err := getbook.NewGetBookRequest(bookID)
	if err != nil {
		t.Fatal(err)
	}

	// Act（実行）
	response := useCase.Execute(ctx, request)

	// Assert（検証）
	if !response.Success {
		t.Errorf("Expected success, got failure: %v", *response.ErrorMessage)
	}
	if response.Data == nil {
		t.Error("Expected book data, got nil")
	}
	if response.Data.ID != "b-12345" {
		t.Errorf("Expected book ID 'b-12345', got %s", response.Data.ID)
	}
	if response.Data.CurrentLoan != nil {
		t.Error("Expected no current loan")
	}
}

func TestGetBookUseCase_ReturnsBookWithLoanWhenBorrowed(t *testing.T) {
	// Arrange（準備）
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockBookQueryService(ctrl)
	useCase := getbook.NewGetBookUseCase(mockQueryService)

	ctx := context.Background()
	bookID := "b-12345"

	// モックの期待値を設定（貸出中の書籍）
	mockQueryService.EXPECT().
		GetBookByID(ctx, bookID).
		Return(&query.BookReadModel{
			ID:          "b-12345",
			ISBN:        "978-0-123456-78-9",
			Title:       "Domain-Driven Design",
			Author:      "Eric Evans",
			IsAvailable: false,
			CurrentLoan: &query.LoanReadModel{
				LoanID:       "l-111",
				UserID:       "u-222",
				BorrowedDate: "2024-01-01",
				DueDate:      "2024-01-15",
			},
		}, nil)

	request, err := getbook.NewGetBookRequest(bookID)
	if err != nil {
		t.Fatal(err)
	}

	// Act（実行）
	response := useCase.Execute(ctx, request)

	// Assert（検証）
	if !response.Success {
		t.Errorf("Expected success, got failure")
	}
	if response.Data.CurrentLoan == nil {
		t.Error("Expected current loan, got nil")
	}
}

func TestGetBookUseCase_ReturnsErrorWhenBookNotFound(t *testing.T) {
	// Arrange（準備）
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockBookQueryService(ctrl)
	useCase := getbook.NewGetBookUseCase(mockQueryService)

	ctx := context.Background()
	bookID := "b-99999"

	// モックの期待値を設定（書籍が見つからない）
	mockQueryService.EXPECT().
		GetBookByID(ctx, bookID).
		Return(nil, nil) // 書籍が見つかりません

	request, err := getbook.NewGetBookRequest(bookID)
	if err != nil {
		t.Fatal(err)
	}

	// Act（実行）
	response := useCase.Execute(ctx, request)

	// Assert（検証）
	if response.Success {
		t.Error("Expected failure, got success")
	}
	if *response.ErrorMessage != "Book not found" {
		t.Errorf("Expected 'Book not found', got %s", *response.ErrorMessage)
	}
}

func TestGetBookRequest_RejectsEmptyBookID(t *testing.T) {
	// Act（実行）
	_, err := getbook.NewGetBookRequest("")

	// Assert（検証）
	if err == nil {
		t.Error("Expected error for empty book ID")
	}
}
