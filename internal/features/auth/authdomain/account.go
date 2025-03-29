package authdomain

import "golang.org/x/crypto/bcrypt"

type AccountID string

type Account struct {
	ID          AccountID
	Email       string
	Name        string
	DisplayName string
}

type PasswordAuthentication struct {
	AccountID
	PasswordHash
}

type Password struct{ password []byte }

func (p Password) String() string {
	return "······"
}

func NewPassword(pw string) Password { return Password{[]byte(pw)} }

type PasswordHash struct{ hash []byte }

func NewHash(pw Password) (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword(pw.password, 0)
	return PasswordHash{hash}, err
}

func (h PasswordHash) Validate(pw Password) bool {
	return bcrypt.CompareHashAndPassword(h.hash, pw.password) == nil
}

type AccountRegistered struct {
	AccountID
}
