package bookdm

import "errors"

// Bookエンティティ - 図書館の本を表す
type Book struct {
  id          BookID
  title       string
  author      string
  isbn        *ISBN
  totalCopies uint32
  // availableCopies is NOT stored - derive it from the Loan table
  // availableCopies = totalCopies - activeLoansForThisBook
}

func NewBook(
  id BookID,
  title string,
  author string,
  isbn *ISBN,
  totalCopies uint32,
) (*Book, error) {
  // コンストラクタでの検証（フェイルファスト）
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

// ゲッター（Goの慣習: "Get"プレフィックスなし）
func (b *Book) Id() BookID {
  return b.id
}

func (b *Book) Title() string {
  return b.title
}

func (b *Book) Author() string {
  return b.author
}

func (b *Book) ISBN() *ISBN {
  return b.isbn
}

func (b *Book) TotalCopies() uint32 {
  return b.totalCopies
}

// ビジネスロジックメソッド
// Book entity focuses on RULES, not STATE
// Available copies are derived from the Loan table

// IsAvailable validates if book has available copies based on provided active loan count
// The count is provided by the use case (derived from Loan table)
func (b *Book) IsAvailable(currentActiveLoans uint32) bool {
  return currentActiveLoans < b.totalCopies
}

// State changes are tracked in Loan table (no BorrowCopy/ReturnCopy methods needed)

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