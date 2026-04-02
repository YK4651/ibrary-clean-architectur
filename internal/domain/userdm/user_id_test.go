package userdm

import (
	"testing"
)

func TestUserID(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		// 有効な8桁ID
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatalf("有効なUserIDの作成に失敗しました: %v", err)
		}
		if userID.Value() != "12345678" {
			t.Errorf("IDは '12345678' であるべきですが、'%s' でした", userID.Value())
		}
	})

	t.Run("ValidateFormat", func(t *testing.T) {
		// 無効: 8桁ではない
		_, err := NewUserID("123")
		if err == nil {
			t.Error("8桁ではないUserIDに対してエラーが期待されました")
		}

		// 無効: 文字が含まれている
		_, err = NewUserID("1234abcd")
		if err == nil {
			t.Error("数字以外を含むUserIDに対してエラーが期待されました")
		}
	})

	t.Run("Generate", func(t *testing.T) {
		userID := GenerateUserID()
		if userID == nil {
			t.Fatal("GenerateUserIDがnilを返しました")
		}
		if len(userID.Value()) != 8 {
			t.Errorf("生成されたUserIDは8桁であるべきですが、%d桁でした", len(userID.Value()))
		}
	})
}
