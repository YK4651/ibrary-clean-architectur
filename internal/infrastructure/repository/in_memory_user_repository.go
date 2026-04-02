package repository

import (
	"github.com/YK4651/ibrary-clean-architectur/internal/domain/userdm"
)

type InMemoryUserRepository struct {
	// メモリ内ストレージ (テスト/開発用)
	users map[string]*userdm.User
}

func NewInMemoryUserRepository() userdm.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*userdm.User),
	}
}

func (r *InMemoryUserRepository) Save(u *userdm.User) error {
	r.users[u.ID().Value()] = u
	return nil
}

func (r *InMemoryUserRepository) FindByID(id *userdm.UserID) (*userdm.User, error) {
	u, exists := r.users[id.Value()]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*userdm.User, error) {
	for _, u := range r.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, nil
}

func (r *InMemoryUserRepository) FindAll() ([]*userdm.User, error) {
	users := make([]*userdm.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}

func (r *InMemoryUserRepository) Delete(id *userdm.UserID) error {
	delete(r.users, id.Value())
	return nil
}

func (r *InMemoryUserRepository) FindUsersWithOverdueFees() ([]*userdm.User, error) {
	var users []*userdm.User
	for _, u := range r.users {
		if u.OverdueFees() > 0 {
			users = append(users, u)
		}
	}
	return users, nil
}
