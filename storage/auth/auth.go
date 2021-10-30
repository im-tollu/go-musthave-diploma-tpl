package auth

type Storage interface {
	CreateUser(u UserToCreate) (User, error)
	SetUserSession(u UserSession) error
	GetUserByLogin(login string) (*User, error)
}
