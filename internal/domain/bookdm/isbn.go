package bookdm

import (
	"fmt"
	"regexp"
	"strings"
)

type ISBN struct {
	value string
}

func NewISBN(value string) (*ISBN, error) {
	// ステップ1: 入力を正規化 - スペースとハイフンを削除
	// 受け入れ: "978-3-16-148410-0" または "9783161484100"
	cleanValue := strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", "")

	// ステップ2: フォーマットを検証 - 正確に13桁でなければならない
	matched, err := regexp.MatchString(`^\d{13}$`, cleanValue)
	if err != nil {
		return nil, fmt.Errorf("failed to validate ISBN format: %w", err)
	}
	if !matched {
		return nil, fmt.Errorf("ISBNは13桁でなければなりません。取得: %s", value)
	}

	// ステップ3: クリーンな値を保存（ハイフンなし）
	return &ISBN{value: cleanValue}, nil
}

func (i *ISBN) Value() string {
	return i.value
}

// ステップ4: 必要に応じてフォーマット済み出力を提供
// 戻り値: "978-3-16-148410-0"
func (i *ISBN) Formatted() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		i.value[0:3],
		i.value[3:4],
		i.value[4:6],
		i.value[6:12],
		i.value[12:13])
}

func (i *ISBN) Equals(other *ISBN) bool {
	return i.value == other.value
}

func (i *ISBN) String() string {
	return i.value
}
