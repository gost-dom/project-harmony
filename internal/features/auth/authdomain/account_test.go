package authdomain_test

import (
	"harmony/internal/features/auth/authdomain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordIsNotContertibleToString(t *testing.T) {
	var pw any = authdomain.NewPassword("s3cret")
	_, ok := pw.(string)
	assert.False(t, ok)
}

func TestPasswordStringifies(t *testing.T) {
	// This isn't a good test, as the return value from String() isn't
	// important. But it is important that the password isn't revealed
	// accidentally in log-files.
	pw := authdomain.NewPassword("s3cret")
	assert.Equal(t, "······", pw.String())
}
