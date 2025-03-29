package htest

import (
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

// GomegaSuite is a specialized [suite.Suite] that add [gomega.Gomega] assertion
// semantics to the test suite.
//
// This can provide more expressive assertions when combined with custom
// mathers.
type GomegaSuite struct {
	suite.Suite
	gomega gomega.Gomega
}

func (s *GomegaSuite) Expect(actual any, extra ...any) gomega.Assertion {
	if s.gomega == nil {
		s.gomega = gomega.NewWithT(s.T())
	}
	return s.gomega.Expect(actual, extra...)
}
