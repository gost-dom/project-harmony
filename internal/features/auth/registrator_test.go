package auth_test

import (
	"context"
	. "harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/testing/mocks/features/auth_mock"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	suite.Suite
	ctx context.Context
	Registrator
	repoMock *auth_mock.MockAccountRepository
}

func (s *RegisterTestSuite) SetupTest() {
	s.repoMock = auth_mock.NewMockAccountRepository(s.T())
	s.repoMock.EXPECT().Insert(mock.Anything, mock.Anything).Return(nil)

	s.Registrator = Registrator{Repository: s.repoMock}
	s.ctx = context.Background()
}

func TestRegister(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) TestValidLogin() {
	pw := authdomain.NewPassword("s3cre7")
	s.Register(s.ctx, RegistratorInput{
		Email:       "jd@example.com",
		Password:    pw,
		Name:        "John Smith",
		DisplayName: "John",
	})

	res := s.repoMock.Calls[0].Arguments.Get(1).(AccountUseCaseResult)
	entity := res.Entity
	events := res.Events

	s.Assert().NotZero(entity.Id)
	s.Assert().Equal("jd@example.com", entity.Email)
	s.Assert().Equal("John Smith", entity.Name)
	s.Assert().Equal("John", entity.DisplayName)

	s.Assert().Equal([]DomainEvent{authdomain.AccountRegistered{
		AccountID: entity.ID(),
	}}, events, "A AccountRegistered domain event was generated")
	s.Assert().True(res.PasswordAuthentication.Password.Validate(pw))
}
