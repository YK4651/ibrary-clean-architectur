package getbook

import "errors"

type GetBookRequest struct {
	BookID string
}

func NewGetBookRequest(bookID string) (*GetBookRequest, error) {
	if bookID == "" {
		return nil, errors.New("bookID cannot be empty")
	}
	return &GetBookRequest{BookID: bookID}, nil
}
