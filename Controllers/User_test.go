package controller

import (
	"bytes"
	"encoding/json"
	model "first-api/Models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	// "fmt"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func TestValidateUserData(t *testing.T) {
	assert := assert.New(t)

	testService := model.User{Name: "", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"}

	err := testService.Validate()
	// fmt.Println(err)

	assert.NotNil(t, err)

	testService.Name = "testing"
	testService.Email = "test"
	err = testService.Validate()
	// fmt.Println(err)

	assert.NotNil(t, err)

	testService.Phone = "999999999"
	err = testService.Validate()
	// fmt.Println(err)

	assert.NotNil(t, err)

	testService.Address = ""
	err = testService.Validate()
	// fmt.Println(err)

	assert.NotNil(t, err)

}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) Validate(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) GetAllUsers(users *[]model.User) error {
	args := m.Called(users)
	return args.Error(0)
}

func (m *MockUserStore) GetUserByID(user *model.User, id string) error {
	args := m.Called(user, id)

	// if the mock is set to return an error, return it
	if args.Error(0) != nil {
		return args.Error(0)
	}

	// otherwise, set the passed user object fields
	user.Id = 1
	user.Name = "test user"
	user.Email = "test@gmail.com"
	user.Phone = "9999999999"
	user.Address = "abcd efgh ijkl"

	return nil
}
func (m *MockUserStore) UpdateUser(user *model.User, id string) error {
	args := m.Called(user, id)
	return args.Error(0)
}
func (m *MockUserStore) DeleteUser(user *model.User, id string) error {
	args := m.Called(user, id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockUserStore := new(MockUserStore)
	user := &model.User{
		Name:  "test user",
		Email: "test@gmail.com",

		Phone:   "9999999999",
		Address: "abcd efgh ijkl",
	}
	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/user", NewUserController(mockUserStore))

	mockUserStore.On("Validate", *user).Return(nil)
	mockUserStore.On("CreateUser", user).Return(nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	// Test
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	mockUserStore.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {
	mockUserStore := new(MockUserStore)

	users := []model.User{
		{Name: "test User 1", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"},
		{Name: "test user 2", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"},
	}
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/user", GetUsers(mockUserStore))

	mockUserStore.On("GetAllUsers", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]model.User)
		*arg = users
	})

	req, _ := http.NewRequest(http.MethodGet, "/user", nil)
	resp := httptest.NewRecorder()

	// Test
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	mockUserStore.AssertExpectations(t)
}

// GetUserByID test function

func TestGetUserById(t *testing.T) {
	mockUserStore := new(MockUserStore)
	user := &model.User{Id: 1, Name: "test user", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/user/:id", GetUserByIDController(mockUserStore))

	// Configure the mock to expect a call to GetUserByID with any user object and the id "1", and to return no error.
	// It should also set the user object's fields to match the fields of the `user` variable.
	mockUserStore.On("GetUserByID", mock.AnythingOfType("*model.User"), "1").Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.User)
		*arg = *user
	})

	// Create the request to get the user.
	req, _ := http.NewRequest(http.MethodGet, "/user/1", nil)
	resp := httptest.NewRecorder()

	// Test
	router.ServeHTTP(resp, req)

	// Check that the response code is 200 OK.
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal the response body.
	var responseUser model.User
	err := json.Unmarshal(resp.Body.Bytes(), &responseUser)

	// If unmarshaling didn't return an error, check that the user fields match.
	if assert.NoError(t, err) {
		assert.Equal(t, user, &responseUser)
	}

	// Check that the mock expectations were met.
	mockUserStore.AssertExpectations(t)
}
