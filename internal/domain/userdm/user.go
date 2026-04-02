package userdm

import (
	"errors"
	"time"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
)

type User struct {
	id          *UserID
	name        string
	email       string
	status      UserStatus
	overdueFees float64
	createdAt   time.Time
}

const MaxLoans = 5

// NewUser は新しいアクティブなユーザーを作成します
func NewUser(name, email string) *User {
	return &User{
		id:          GenerateUserID(),
		name:        name,
		email:       email,
		status:      UserStatusActive,
		overdueFees: 0,
		createdAt:   time.Now(),
	}
}

// ReconstructUser は永続化からユーザーを再構築します
func ReconstructUser(
	id *UserID,
	name, email string,
	status UserStatus,
	overdueFees float64,
	createdAt time.Time,
) *User {
	return &User{id, name, email, status, overdueFees, createdAt}
}

// 提供された貸出数に基づいて貸出可能性を検証
// 貸出数はユースケースによって提供される（Loanテーブルから導出）
func (u *User) CanBorrow(currentLoanCount uint32) bool {
	if u.status == UserStatusSuspended {
		return false
	}
	if currentLoanCount >= MaxLoans {
		return false
	}
	if u.overdueFees > 0 {
		return false
	}
	return true
}

func (u *User) AddOverdueFee(amount float64) (*User, error) {
	if amount <= 0 {
		return nil, errors.New("延滞料金は0より大きい必要があります")
	}
	return &User{
		id:          u.id,
		name:        u.name,
		email:       u.email,
		status:      u.status,
		overdueFees: u.overdueFees + amount,
		createdAt:   u.createdAt,
	}, nil
}

// 不変の状態変更: 延滞料金を支払う
func (u *User) PayOverdueFee(amount float64) (*User, error) {
	if amount < 0 {
		return nil, errors.New("支払い金額は負の値にできません")
	}
	if amount > u.overdueFees {
		return nil, errors.New("支払い金額が現在の延滞料金を超えています")
	}
	return &User{
		id:          u.id,
		name:        u.name,
		email:       u.email,
		status:      u.status,
		overdueFees: u.overdueFees - amount,
		createdAt:   u.createdAt,
	}, nil
}

func (u *User) Suspend() *User {
	return &User{
		id:          u.id,
		name:        u.name,
		email:       u.email,
		status:      UserStatusSuspended,
		overdueFees: u.overdueFees,
		createdAt:   u.createdAt,
	}
}

// 不変の状態変更: アカウントを有効化
func (u *User) Activate() *User {
	return &User{
		id:          u.id,
		name:        u.name,
		email:       u.email,
		status:      UserStatusActive,
		overdueFees: u.overdueFees,
		createdAt:   u.createdAt,
	}
}

// ゲッター
func (u *User) ID() *UserID          { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) Status() UserStatus   { return u.status }
func (u *User) OverdueFees() float64 { return u.overdueFees }
func (u *User) CreatedAt() time.Time { return u.createdAt }
