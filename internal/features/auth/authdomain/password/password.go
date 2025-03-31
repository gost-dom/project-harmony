package password

import "golang.org/x/crypto/bcrypt"

type Password struct{ password []byte }

func Parse(pw string) Password { return Password{[]byte(pw)} }

func (p Password) String() string { return "······" }

func (p Password) GoString() string { return p.String() }

func (p Password) Hash() (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword(p.password, bcrypt.MinCost)
	return PasswordHash{hash}, err
}

func (p Password) Equals(other Password) bool {
	return string(p.password) == string(other.password)
}

type PasswordHash struct{ hash []byte }

func (h PasswordHash) Validate(pw Password) bool {
	return bcrypt.CompareHashAndPassword(h.hash, pw.password) == nil
}
