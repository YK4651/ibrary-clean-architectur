package userdm

// ドメイン層で定義されたリポジトリインターフェース
type UserRepository interface {
	Save(user *User) error
	FindByID(id *UserID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	Delete(id *UserID) error
	FindUsersWithOverdueFees() ([]*User, error)
}
