package userdm

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	t.Run("CanBorrow", func(t *testing.T) {
		u := NewUser("田中太郎", "tanaka@example.com")

		if !u.CanBorrow(0) {
			t.Error("新規ユーザーは本を借りられるべきです")
		}

		if u.Name() != "田中太郎" {
			t.Errorf("名前は '田中太郎' であるべきですが、'%s' でした", u.Name())
		}
		if u.Email() != "tanaka@example.com" {
			t.Errorf("メールは 'tanaka@example.com' であるべきですが、'%s' でした", u.Email())
		}
		if u.Status() != UserStatusActive {
			t.Errorf("ステータスは 'active' であるべきですが、'%s' でした", u.Status())
		}
	})

	t.Run("Suspend", func(t *testing.T) {
		u := NewUser("田中太郎", "tanaka@example.com")

		u2 := u.Suspend()

		if u2.Status() != UserStatusSuspended {
			t.Errorf("ステータスは 'suspended' であるべきですが、'%s' でした", u2.Status())
		}

		// 元のユーザーは変更されないべき（不変性）
		if u.Status() != UserStatusActive {
			t.Error("元のユーザーは変更されないべきです")
		}
	})

	t.Run("CannotBorrowWhenSuspended", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}

		// 停止中のユーザーを再構築
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusSuspended,
			0,
			time.Now(),
		)

		if u.CanBorrow(0) {
			t.Error("停止中のユーザーは本を借りられないべきです")
		}
	})

	t.Run("CannotBorrowAtMaxLimit", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}

		// 最大貸出数のユーザーを再構築
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusActive,
			0,
			time.Now(),
		)

		if u.CanBorrow(MaxLoans) {
			t.Error("最大貸出数に達したユーザーは本を借りられないべきです")
		}
	})

	t.Run("CannotBorrowWithOverdueFees", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}

		// 延滞料金があるユーザーを再構築
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusActive,
			10.50, // 延滞料金あり
			time.Now(),
		)

		if u.CanBorrow(0) {
			t.Error("延滞料金があるユーザーは本を借りられないべきです")
		}
	})

	t.Run("FeeManagement", func(t *testing.T) {
		u := NewUser("田中太郎", "tanaka@example.com")

		u2, err := u.AddOverdueFee(5.00)
		if err != nil {
			t.Fatalf("延滞料金の追加に失敗しました: %v", err)
		}

		if u2.OverdueFees() != 5.00 {
			t.Errorf("延滞料金は 5.00 であるべきですが、%.2f でした", u2.OverdueFees())
		}

		// 負の料金をテスト
		_, err = u.AddOverdueFee(-1.00)
		if err == nil {
			t.Error("負の料金に対してエラーが期待されました")
		}
	})

	t.Run("Immutability", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusActive,
			0,
			time.Now(),
		)

		// AddOverdueFeeの不変性をテスト
		u2, err := u.AddOverdueFee(5.00)
		if err != nil {
			t.Fatalf("延滞料金の追加に失敗しました: %v", err)
		}

		if u2.OverdueFees() != 5.00 {
			t.Errorf("延滞料金は 5.00 であるべきですが、%.2f でした", u2.OverdueFees())
		}

		// 元のユーザーは変更されないべき
		if u.OverdueFees() != 0 {
			t.Error("元のユーザーは変更されないべきです")
		}

		// Suspendの不変性をテスト
		u3 := u.Suspend()
		if u3.Status() != UserStatusSuspended {
			t.Errorf("ステータスは 'suspended' であるべきですが、'%s' でした", u3.Status())
		}

		// 元のユーザーは変更されないべき
		if u.Status() != UserStatusActive {
			t.Error("元のユーザーは変更されないべきです")
		}
	})
}
