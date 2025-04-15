package htest

import (
	"context"

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
}

func (s *GomegaSuite) Expect(actual any, extra ...any) gomega.Assertion {
	return gomega.NewWithT(s.T()).Expect(actual, extra...)
}

func (s *GomegaSuite) Context() context.Context { return s.T().Context() }
