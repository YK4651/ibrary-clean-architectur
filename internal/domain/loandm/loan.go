package loandm

import (
  "errors"
  "time"

  "github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
  "github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

// Loanエンティティ - 本の貸出を表す
type Loan struct {
  id         LoanID
  userID     userdm.UserID
  bookID     bookdm.BookID
  borrowedAt time.Time
  dueDate    time.Time
  returnedAt *time.Time
}

const LoanPeriodDays = 14

func NewLoan(
  id LoanID,
  userID userdm.UserID,
  bookID bookdm.BookID,
  borrowedAt time.Time,
  returnedAt *time.Time,
) *Loan {
  if borrowedAt.IsZero() {
    borrowedAt = time.Now()
  }

  dueDate := calculateDueDate(borrowedAt)

  return &Loan{
    id:         id,
    userID:     userID,
    bookID:     bookID,
    borrowedAt: borrowedAt,
    dueDate:    dueDate,
    returnedAt: returnedAt,
  }
}

func calculateDueDate(borrowedAt time.Time) time.Time {
  return borrowedAt.AddDate(0, 0, LoanPeriodDays)
}

func (l *Loan) IsOverdue(currentDate *time.Time) bool {
  var now time.Time
  if currentDate != nil {
    now = *currentDate
  } else {
    now = time.Now()
  }
  return now.After(l.dueDate) && l.returnedAt == nil
}

func (l *Loan) IsReturned() bool {
  return l.returnedAt != nil
}

func (l *Loan) MarkAsReturned(returnedAt *time.Time) (*Loan, error) {
  if l.IsReturned() {
    return nil, errors.New("loan has already been returned")
  }

  var returnTime *time.Time
  if returnedAt != nil {
    returnTime = returnedAt
  } else {
    now := time.Now()
    returnTime = &now
  }

  return &Loan{
    id:         l.id,
    userID:     l.userID,
    bookID:     l.bookID,
    borrowedAt: l.borrowedAt,
    dueDate:    l.dueDate,
    returnedAt: returnTime,
  }, nil
}

func (l *Loan) DaysUntilDue(currentDate *time.Time) int32 {
  var now time.Time
  if currentDate != nil {
    now = *currentDate
  } else {
    now = time.Now()
  }

  duration := l.dueDate.Sub(now)
  days := int32(duration.Hours() / 24)

  return days
}

// ゲッター（Goの慣習: "Get"プレフィックスなし）
func (l *Loan) Id() LoanID {
  return l.id
}

func (l *Loan) UserID() userdm.UserID {
  return l.userID
}

func (l *Loan) BookID() bookdm.BookID {
  return l.bookID
}

func (l *Loan) BorrowedAt() time.Time {
  return l.borrowedAt
}

func (l *Loan) DueDate() time.Time {
  return l.dueDate
}

func (l *Loan) ReturnedAt() *time.Time {
  return l.returnedAt
}

func (l *Loan) PeriodDays() uint32 {
  return LoanPeriodDays
}