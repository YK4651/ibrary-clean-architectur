package createuser

import (
	"errors"
	"fmt"
	"time"

	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

// 入力DTO
type CreateUserInput struct {
	Name  string
	Email string
}

// 出力DTO
type CreateUserOutput struct {
	ID          string
	Name        string
	Email       string
	Status      string
	OverdueFees float64
	CreatedAt   time.Time
}

type CreateUserUseCase struct {
	// コンストラクター: リポジトリインターフェースに依存
	userRepository userdm.UserRepository
}

func NewCreateUserUseCase(repo userdm.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	// ビジネスルール: 重複メールをチェック
	existingUser, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// ファクトリーメソッドでユーザーを作成
	u := userdm.NewUser(input.Name, input.Email)

	// リポジトリに永続化
	if err := uc.userRepository.Save(u); err != nil {
		return nil, err
	}

	// DTOを返す
	return &CreateUserOutput{
		ID:          u.ID().Value(),
		Name:        u.Name(),
		Email:       u.Email(),
		Status:      string(u.Status()),
		OverdueFees: u.OverdueFees(),
		CreatedAt:   u.CreatedAt(),
	}, nil
}
