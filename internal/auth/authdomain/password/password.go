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

// HashFromBytes constructs a PasswordHash from a stored byte slice. This is
// only intended to be used from database layers, recreating an instance with a
// value retrieved from [PasswordHash.UnsecureRead].
func HashFromBytes(b []byte) PasswordHash { return PasswordHash{b} }

// UnsecureRead retrieves the underlying data. This is _only_ intended for use
// by database layers, allowing them to persist the data.
func (h PasswordHash) UnsecureRead() []byte { return h.hash }

func (h PasswordHash) Validate(pw Password) bool {
	return bcrypt.CompareHashAndPassword(h.hash, pw.password) == nil
}
