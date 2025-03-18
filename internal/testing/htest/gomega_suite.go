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
	gomega.Gomega
}

func (s *GomegaSuite) SetupTest() {
	s.Gomega = gomega.NewWithT(s.T())
}
