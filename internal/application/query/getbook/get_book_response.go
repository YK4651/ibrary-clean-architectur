package getbook

import "github.com/YK4651/ibrary-clean-architectur/internal/application/query"

type GetBookResponse struct {
	Success      bool
	Data         *query.BookReadModel
	ErrorMessage *string
}

func NewSuccessResponse(data *query.BookReadModel) *GetBookResponse {
	return &GetBookResponse{
		Success:      true,
		Data:         data,
		ErrorMessage: nil,
	}
}

func NewNotFoundResponse() *GetBookResponse {
	msg := "Book not found"
	return &GetBookResponse{
		Success:      false,
		Data:         nil,
		ErrorMessage: &msg,
	}
}

func NewFailureResponse(message string) *GetBookResponse {
	return &GetBookResponse{
		Success:      false,
		Data:         nil,
		ErrorMessage: &message,
	}
}
