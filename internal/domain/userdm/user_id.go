package userdm

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type UserID struct {
	value string
}

// NewUserID は検証付きでUserIDを作成します（8桁）
func NewUserID(value string) (*UserID, error) {
	matched, err := regexp.MatchString(`^\d{8}$`, value)
	if err != nil {
		return nil, fmt.Errorf("failed to validate UserID format: %w", err)
	}
	if !matched {
		return nil, errors.New("UserIDは正確に8桁である必要があります。入力値: " + value)
	}
	return &UserID{value: value}, nil
}

// GenerateUserID はランダムな8桁のUserIDを作成します
func GenerateUserID() *UserID {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomId := r.Intn(90000000) + 10000000
	return &UserID{value: strconv.Itoa(randomId)}
}

func (u *UserID) Value() string {
	return u.value
}

func (u *UserID) Equals(other *UserID) bool {
	return u.value == other.value
}

func (u *UserID) String() string {
	return u.value
}
