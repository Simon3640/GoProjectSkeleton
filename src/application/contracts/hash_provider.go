package contracts

type IHashProvider interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) (bool, error)
}
