package main_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	main "Users/htahreem/VSCode/unittests"
	database "Users/htahreem/VSCode/unittests/database"
	databasemock "Users/htahreem/VSCode/unittests/database/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set up database connection here first, without connection to the origional database it will not work
	database.ConnectDatabase()

	// Run tests
	exitCode := m.Run()

	// Exit with the test status code
	os.Exit(exitCode)
}

func TestGetAllStudents(t *testing.T) {
	// Create an instance of the MockDatabase
	mockDB := databasemock.NewMockDatabase(t)

	// Set expectations for the Query method
	mockRows := &sql.Rows{} // create mock rows 
	mockDB.On("Query", "SELECT * FROM students").Return(mockRows, nil)

	// Create a mock Gin context for testing
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Call the GetAllStudents function with the mock context and mockDB
	main.GetAllStudents(ctx)

	// Check the response
	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

}

func TestAddStudent(t *testing.T) {
	mockDB := databasemock.NewMockDatabase(t)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// sample student data
	stu := main.Student{
		Name:      "John Doe",
		RollNo:    123,
		ContactNo: 3210,
		Email:     "john@example.com",
		ID:        "abc123",
	}

	// Convert student data to JSON
	jsonData, _ := json.Marshal(stu)

	// Set up the mock database expectations
	mockDB.On("Exec", `INSERT INTO students VALUES ($1, $2, $3, $4, $5)`, stu.Name, stu.RollNo, stu.ContactNo, stu.Email, stu.ID).Return(nil)

	// Set the request body with student data
	ctx.Request = httptest.NewRequest("POST", "/addStudent", strings.NewReader(string(jsonData)))

	// Call the AddStudent function with the mock context and mockDB
	main.AddStudent(ctx)

	// Check the response
	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

}

func TestUpdateUser(t *testing.T) {
	mockDB := databasemock.NewMockDatabase(t)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	stu := main.Student{
		Name:      "Updated Name",
		RollNo:    456,
		ContactNo: 9211,
		Email:     "updated@example.com",
		ID:        "xyz789",
	}

	jsonData, _ := json.Marshal(stu)

	ctx.Request = httptest.NewRequest("PUT", fmt.Sprintf("/updateStudent/%s", stu.ID), strings.NewReader(string(jsonData)))

	// Set up the mock database expectations for the transaction and update
	mockDB.On("Begin").Return(&sql.Tx{}, nil)
	mockDB.On("Exec", `
        UPDATE students
        SET "Name" = $1,
            "RollNo" = $2,
            "ContactNo" = $3,
            "Email" = $4
        WHERE "ID" = $5`,
		stu.Name, stu.RollNo, stu.ContactNo, stu.Email, stu.ID).Return(nil)
	mockDB.On("Commit").Return(nil)

	// Call the UpdateUser function with the mock context and mockDB
	main.UpdateUser(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	
}

func TestDeleteStudent(t *testing.T) {
	mockDB := databasemock.NewMockDatabase(t)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Set the request parameter with a sample student ID
	ctx.Params = append(ctx.Params, gin.Param{Key: "ID", Value: "xyz123"})

	// Set up the mock database expectations for the transaction and delete
	mockDB.On("Begin").Return(&sql.Tx{}, nil)
	mockDB.On("Exec", "DELETE FROM students WHERE \"ID\" = $1", "xyz123").Return(nil)
	mockDB.On("Commit").Return(nil)

	// Call the DeleteStudent function with the mock context and mockDB
	main.DeleteStudent(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

}
