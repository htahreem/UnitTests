// database/mocks/mock_database.go
package mocks

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

// Query is a mock method for the Query in the Database interface
func (m *MockDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Mock behavior: Return predefined data for testing GetAllStudents function
	// argsMock := m.Called(query, args...)
	argsMock := m.Called(append([]interface{}{query}, args...)...)
	return argsMock.Get(0).(*sql.Rows), argsMock.Error(1)
}

// NewMockDatabase creates and returns a new instance of MockDatabase
func NewMockDatabase(t *testing.T) *MockDatabase {
	mockDB := new(MockDatabase)

	// Set up expectations or behaviors specific to your tests here
	// Example: mockDB.On("Query", "SELECT * FROM students").Return(mockRows, nil)

	return mockDB
}
