package authorization

import "github.com/stretchr/testify/mock"

type MockAuthorization struct {
	mock.Mock
}

func NewAuthorizationMock() *MockAuthorization {
	return &MockAuthorization{}
}

func (m *MockAuthorization) GenerateToken(claim AppAuthorizationClaim) (token string, err error) {
	args := m.Called(claim)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorization) ValidateToken(token string, paserTo interface{}) error {
	args := m.Called(token, paserTo)
	return args.Error(1)
}
