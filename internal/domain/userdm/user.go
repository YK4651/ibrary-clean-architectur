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

const MaxLoans = 5

type User struct {
	id          *UserID
	name        string
	email       string
	status      UserStatus
	overdueFees float64
	createdAt   time.Time
}

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

// Userエンティティはルールに焦点を当て、状態は持たない
// 貸出数はLoanテーブルから導出

// CanBorrow - 提供された貸出数に基づいて貸出可能性を検証
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

// 状態の変更はLoanテーブルで追跡される（BorrowBook/ReturnBookメソッドは不要）

func (u *User) AddOverdueFee(amount float64) (*User, error) {
	if amount <= 0 {
		return nil, errors.New("overdue fee must be greater than 0")
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

func (u *User) PayOverdueFee(amount float64) (*User, error) {
	if amount < 0 {
		return nil, errors.New("payment amount cannot be negative")
	}
	if amount > u.overdueFees {
		return nil, errors.New("payment exceeds current overdue fees")
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
func (u *User) Id() *UserID          { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) Status() UserStatus   { return u.status }
func (u *User) OverdueFees() float64 { return u.overdueFees }
func (u *User) CreatedAt() time.Time { return u.createdAt }
