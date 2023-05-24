package Controllers

import (
	"first-api/Models"
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)



func TestValidateUserData(t *testing.T){
	assert  := assert.New(t)

	testService := Models.User{Name:"",Email:"test@gmail.com",Phone:"9999999999",Address:"abcd efgh ijkl"}

	err := testService.Validate();
	// fmt.Println(err)

	assert.NotNil(t,err);

	testService.Name = "testing"
	testService.Email = "test"
	err = testService.Validate();
	// fmt.Println(err)

	assert.NotNil(t,err);


	testService.Phone = "999999999" 
	err = testService.Validate();
	// fmt.Println(err)

	assert.NotNil(t,err);

	testService.Address = ""
	err = testService.Validate()
	// fmt.Println(err)

	assert.NotNil(t,err);

}



type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUser(user *Models.User) ( error) {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) Validate(user Models.User) error {
	fmt.Println("VAlidate called")
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) GetAllUsers(users *[]Models.User) error{
	fmt.Println("GEt all users mock called")
	args:= m.Called(users)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockUserStore := new(MockUserStore)
	user := &Models.User{Name:"test user",Email:"test@gmail.com",Phone:"9999999999",Address:"abcd efgh ijkl"}
	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/user", NewUserController(mockUserStore))


	mockUserStore.On("Validate",*user).Return(nil)
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

	users :=[]Models.User{
		{Name:"test User 1",Email:"test@gmail.com",Phone:"9999999999",Address:"abcd efgh ijkl"},
		{Name:"test user 2",Email:"test@gmail.com",Phone:"9999999999",Address:"abcd efgh ijkl"},
	}
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/user", GetUsers(mockUserStore))


	mockUserStore.On("GetAllUsers",mock.Anything).Return(nil).Run(func(args mock.Arguments){
		arg:=args.Get(0).(*[]Models.User)
		*arg = users
	})

	req, _ := http.NewRequest(http.MethodGet, "/user", nil)
	resp := httptest.NewRecorder()

	// Test
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	mockUserStore.AssertExpectations(t)
}
