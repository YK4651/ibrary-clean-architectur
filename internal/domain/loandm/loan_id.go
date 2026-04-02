package loandm

import "github.com/YK4651/ibrary-clean-architectur/internal/domain/shared"

// LoanID は貸出エンティティのID
type LoanID shared.ID

// NewLoanID は新しいLoanIDを生成
func NewLoanID() LoanID {
	return LoanID(shared.NewID())
}

// Value は文字列としての値を返す
func (id LoanID) Value() string {
	return shared.ID(id).Value()
}

// Equals は値オブジェクトの等価性を比較
func (id LoanID) Equals(other LoanID) bool {
	return shared.ID(id).Equals(shared.ID(other))
}
