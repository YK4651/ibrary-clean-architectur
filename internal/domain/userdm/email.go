package userdm

import (
	"fmt"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	// メールアドレスのフォーマットを検証
	matched, err := regexp.MatchString(`^[^\s@]+@[^\s@]+\.[^\s@]+$`, value)
	if err != nil {
		return nil, fmt.Errorf("メールアドレスフォーマットの検証に失敗: %w", err)
	}
	if !matched {
		return nil, fmt.Errorf("無効なメールアドレスフォーマット: %s", value)
	}

	return &Email{value: value}, nil
}

func (e *Email) Value() string {
	return e.value
}

func (e *Email) Domain() string {
	parts := strings.Split(e.value, "@")
	return parts[1] // @以降の部分
}

func (e *Email) Equals(other *Email) bool {
	return e.value == other.value
}

func (e *Email) String() string {
	return e.value
}
