package loandm

import (
	"fmt"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

// LoanEligibilityService - Domain Service for loan eligibility checks
type LoanEligibilityService struct{}

func NewLoanEligibilityService() *LoanEligibilityService {
	return &LoanEligibilityService{}
}

func (s *LoanEligibilityService) CanBorrow(u *userdm.User, currentLoanCount uint32, b *bookdm.Book, currentActiveLoans uint32) bool {
	// Rule 1 & 2: User must be eligible (loan limit + overdue check)
	if !u.CanBorrow(currentLoanCount) {
		return false
	}

	// Rule 3: Book must have available copies
	if !b.IsAvailable(currentActiveLoans) {
		return false
	}

	return true
}

func (s *LoanEligibilityService) IneligibilityReason(u *userdm.User, currentLoanCount uint32, b *bookdm.Book, currentActiveLoans uint32) *string {
	if u.Status() == userdm.UserStatusSuspended {
		reason := "ユーザーアカウントが停止されています"
		return &reason
	}

	if currentLoanCount >= userdm.MaxLoans {
		reason := fmt.Sprintf("最大貸出数に達しています（%d/%d冊）", currentLoanCount, userdm.MaxLoans)
		return &reason
	}

	if u.HasOverdueBooks() {
		reason := "延滞中の本があるため貸出できません"
		return &reason
	}

	if !b.IsAvailable(currentActiveLoans) {
		reason := fmt.Sprintf("書籍「%s」に貸出可能な在庫がありません", b.Title())
		return &reason
	}

	return nil
}
