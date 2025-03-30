package authdomain_test

import (
	"fmt"
	"harmony/internal/features/auth/authdomain/password"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordIsNotContertibleToString(t *testing.T) {
	var pw any = password.Parse("s3cret")
	_, ok := pw.(string)
	assert.False(t, ok)
}

func TestPasswordStringerDoesntRevealPassword(t *testing.T) {
	assert := assert.New(t)
	// It is important that the _real_ password is not accidentally revealed in
	// e.g., log files, console output, etc.
	input := "s3cret"
	binaryForm1 := "115 51 99 114 101 116"
	binaryForm2 := "0x73, 0x33, 0x63, 0x72, 0x65, 0x74"
	assert.Contains(
		fmt.Sprintf("%+v", []byte(input)), binaryForm1,
		"Test error, binaryForm1 not valid input for NotContains later",
	)
	assert.Contains(
		fmt.Sprintf("%#v", []byte(input)), binaryForm2,
		"Test error, binaryForm2 not valid input for NotContains later",
	)

	pw := password.Parse(input)
	for _, searchVal := range []string{input, binaryForm1, binaryForm2} {
		assert.NotContains(
			fmt.Sprintf("%s", pw), searchVal,
			fmt.Sprintf("Found string %s when formatting password with '%%s'", searchVal),
		)
		assert.NotContains(
			fmt.Sprintf("%+s", pw), searchVal,
			fmt.Sprintf("Found string %s when formatting password with '%%+s'", searchVal),
		)
		assert.NotContains(
			fmt.Sprintf("%+v", pw), searchVal,
			fmt.Sprintf("Found string %s when formatting password with '%%+v'", searchVal),
		)
		assert.NotContains(
			fmt.Sprintf("%#v", pw), searchVal,
			fmt.Sprintf("Found string %s when formatting password with '%%#v'", searchVal),
		)
	}
}
