package hcnettoml

import "github.com/stretchr/testify/mock"

// MockClient is a mockable hcnettoml client.
type MockClient struct {
	mock.Mock
}

// GetHcnetToml is a mocking a method
func (m *MockClient) GetHcnetToml(domain string) (*Response, error) {
	a := m.Called(domain)
	return a.Get(0).(*Response), a.Error(1)
}

// GetHcnetTomlByAddress is a mocking a method
func (m *MockClient) GetHcnetTomlByAddress(address string) (*Response, error) {
	a := m.Called(address)
	return a.Get(0).(*Response), a.Error(1)
}
