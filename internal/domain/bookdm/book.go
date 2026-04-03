package bookdm

import "errors"

type Book struct {
	id          BookID
	title       string
	author      string
	isbn        *ISBN
	totalCopies uint32
	// availableCopies is NOT stored - derive it from the Loan table
	// availableCopies = totalCopies - activeLoansForThisBook
}

// NewBook - 新しいBookを作成
func NewBook(
	id BookID,
	title string,
	author string,
	isbn *ISBN,
	totalCopies uint32,
) (*Book, error) {
	if len(title) == 0 {
		return nil, errors.New("book title cannot be empty")
	}
	if len(author) == 0 {
		return nil, errors.New("book author cannot be empty")
	}
	if totalCopies < 1 {
		return nil, errors.New("total copies must be at least 1")
	}

	return &Book{
		id:          id,
		title:       title,
		author:      author,
		isbn:        isbn,
		totalCopies: totalCopies,
	}, nil
}

// Book エンティティはルールに焦点を当て、状態は持たない
// 在庫数はLoanテーブルから導出

// IsAvailable - 提供されたアクティブ貸出数に基づいて在庫の有無を検証
// 貸出数はユースケースによって提供される（Loanテーブルから導出）
func (b *Book) IsAvailable(currentActiveLoans uint32) bool {
	return currentActiveLoans < b.totalCopies
}

// 状態の変更はLoanテーブルで追跡される（BorrowCopy/ReturnCopyメソッドは不要）

// 状態変更メソッド
func (b *Book) UpdateTitle(newTitle string) error {
	if len(newTitle) == 0 {
		return errors.New("book title cannot be empty")
	}
	b.title = newTitle
	return nil
}

func (b *Book) UpdateAuthor(newAuthor string) error {
	if len(newAuthor) == 0 {
		return errors.New("book author cannot be empty")
	}
	b.author = newAuthor
	return nil
}

// Getters
func (b *Book) Id() BookID          { return b.id }
func (b *Book) Title() string       { return b.title }
func (b *Book) Author() string      { return b.author }
func (b *Book) ISBN() *ISBN         { return b.isbn }
func (b *Book) TotalCopies() uint32 { return b.totalCopies }
