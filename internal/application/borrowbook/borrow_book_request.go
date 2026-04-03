package borrowbook

import "errors"

type BorrowBookRequest struct {
	UserID string
	BookID string
}

func NewBorrowBookRequest(userID, bookID string) (*BorrowBookRequest, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if bookID == "" {
		return nil, errors.New("bookID cannot be empty")
	}
	return &BorrowBookRequest{
		UserID: userID,
		BookID: bookID,
	}, nil
}
