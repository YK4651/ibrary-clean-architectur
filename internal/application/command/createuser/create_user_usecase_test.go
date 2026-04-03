package createuser_test

import (
	"testing"

	"github.com/YK4651/ibrary-clean-architectur/internal/application/command/createuser"
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

// テスト用のモックリポジトリ
type mockUserRepository struct {
	users map[string]*userdm.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*userdm.User),
	}
}

func (m *mockUserRepository) Save(u *userdm.User) error {
	m.users[u.Id().Value()] = u
	return nil
}

func (m *mockUserRepository) FindByID(id *userdm.UserID) (*userdm.User, error) {
	u, exists := m.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) FindByEmail(email string) (*userdm.User, error) {
	for _, u := range m.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) FindAll() ([]*userdm.User, error) {
	users := make([]*userdm.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

func (m *mockUserRepository) Delete(id *userdm.UserID) error {
	delete(m.users, id.Value())
	return nil
}

func (m *mockUserRepository) FindUsersWithOverdueFees() ([]*userdm.User, error) {
	var users []*userdm.User
	for _, u := range m.users {
		if u.OverdueFees() > 0 {
			users = append(users, u)
		}
	}
	return users, nil
}

func TestCreateUserUseCase(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := newMockUserRepository()
		uc := createuser.NewCreateUserUseCase(repo)

		input := createuser.CreateUserInput{
			Name:  "田中太郎",
			Email: "tanaka@example.com",
		}

		output, err := uc.Execute(input)
		if err != nil {
			t.Fatalf("エラーが発生しないべきですが、%v が発生しました", err)
		}

		if output.Name != "田中太郎" {
			t.Errorf("名前は '田中太郎' であるべきですが、'%s' でした", output.Name)
		}
		if output.Email != "tanaka@example.com" {
			t.Errorf("メールは 'tanaka@example.com' であるべきですが、'%s' でした", output.Email)
		}
		if output.Status != "active" {
			t.Errorf("ステータスは 'active' であるべきですが、'%s' でした", output.Status)
		}
		if output.OverdueFees != 0 {
			t.Errorf("延滞料金は 0 であるべきですが、%.2f でした", output.OverdueFees)
		}

		// ユーザーがリポジトリに保存されたか確認
		users, err := repo.FindAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 1 {
			t.Errorf("リポジトリには 1 人のユーザーが存在するべきですが、%d 人でした", len(users))
		}
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		repo := newMockUserRepository()
		uc := createuser.NewCreateUserUseCase(repo)

		// 最初のユーザーを作成
		input1 := createuser.CreateUserInput{
			Name:  "田中太郎",
			Email: "tanaka@example.com",
		}
		_, err := uc.Execute(input1)
		if err != nil {
			t.Fatalf("最初のユーザーの作成に失敗しました: %v", err)
		}

		// 同じメールで2番目のユーザーを作成しようとする
		input2 := createuser.CreateUserInput{
			Name:  "佐藤花子",
			Email: "tanaka@example.com", // 重複メール
		}
		_, err = uc.Execute(input2)
		if err == nil {
			t.Error("重複メールに対してエラーが期待されました")
		}

		// ユーザーが1人だけ存在することを確認
		users, err := repo.FindAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 1 {
			t.Errorf("リポジトリには 1 人のユーザーが存在するべきですが、%d 人でした", len(users))
		}
	})
}
