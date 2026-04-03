package loandm

import (
	"fmt"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/bookdm"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

// LoanEligibilityService - 貸出可否判定のドメインサービス
type LoanEligibilityService struct{}

func NewLoanEligibilityService() *LoanEligibilityService {
	return &LoanEligibilityService{}
}

// CanBorrow - currentLoanCount: ユーザーの現在の貸出数, activeLoansForBook: 当該書籍のアクティブ貸出数（Loanテーブルから取得）
func (s *LoanEligibilityService) CanBorrow(u *userdm.User, currentLoanCount uint32, b *bookdm.Book, activeLoansForBook uint32) bool {
	// Rule 1 & 2: ユーザーが貸出可能か（停止中でないか、上限以下か、延滞料なしか）
	if !u.CanBorrow(currentLoanCount) {
		return false
	}

	// Rule 3: 書籍に利用可能なコピーがあるか（Loanテーブルから導出）
	if !b.IsAvailable(activeLoansForBook) {
		return false
	}

	return true
}

// IneligibilityReason - 貸出不可の理由を返す
func (s *LoanEligibilityService) IneligibilityReason(u *userdm.User, currentLoanCount uint32, b *bookdm.Book, activeLoansForBook uint32) *string {
	if !u.CanBorrow(currentLoanCount) {
		if currentLoanCount >= userdm.MaxLoans {
			reason := fmt.Sprintf("ユーザーは最大貸出制限に達しています（%d冊）", userdm.MaxLoans)
			return &reason
		}
		if u.OverdueFees() > 0 {
			reason := "ユーザーは延滞中の本があります"
			return &reason
		}
		reason := "ユーザーは貸出できません（停止中または条件を満たしていません）"
		return &reason
	}

	if !b.IsAvailable(activeLoansForBook) {
		reason := fmt.Sprintf("書籍「%s」の利用可能なコピーがありません", b.Title())
		return &reason
	}

	return nil
}
