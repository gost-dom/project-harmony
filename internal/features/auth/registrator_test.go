package auth_test

import (
	"context"
	. "harmony/internal/features/auth"
	"harmony/internal/testing/mocks/features/auth_mock"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegistratorInput struct {
	Email    string
	Password string
}

type Registrator struct {
	Repository AccountRepository
}

func (r Registrator) Register(ctx context.Context, input RegistratorInput) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	account := Account{
		Id:    AccountID(id),
		Email: input.Email,
	}
	res := NewResult(account)
	res.AddEvent(AccountRegistered{AccountID: account.Id})
	return r.Repository.Insert(ctx, *res)
}

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
	s.Register(s.ctx, RegistratorInput{Email: "jd@example.com"})

	res := s.repoMock.Calls[0].Arguments.Get(1).(AccountUseCaseResult)
	entity := res.Entity
	events := res.Events

	s.Assert().NotZero(entity.Id)
	s.Assert().Equal("jd@example.com", entity.Email)

	s.Assert().Equal([]DomainEvent{AccountRegistered{
		AccountID: entity.ID(),
	}}, events, "A AccountRegistered domain event was generated")
}
