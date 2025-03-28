package authdomain

import "golang.org/x/crypto/bcrypt"

type AccountID string

type Account struct {
	Id       AccountID
	Email    string
	Password PasswordHash
}

func (a Account) ID() AccountID { return a.Id }

type Password struct{ password string }

func NewPassword(pw string) Password { return Password{pw} }

type PasswordHash struct{ hash []byte }

func NewHash(pw Password) (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw.password), 0)
	return PasswordHash{hash}, err
}

func (h PasswordHash) Validate(pw Password) bool {
	return bcrypt.CompareHashAndPassword(h.hash, []byte(pw.password)) == nil
}

func (p Password) Validate(comp string) bool {
	return comp == p.password
}

type AccountRegistered struct {
	AccountID
}
