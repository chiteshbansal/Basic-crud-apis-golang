package Controllers

import (
	"first-api/Models"
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"first-api/Models/"
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

func (m *MockUserStore) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) Validate() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockUserStore := new(MockUserStore)
	user := &models.User{
		Name: "Test User",
	}

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/user", controllers.NewUserController(mockUserStore))


	mockUserStore.On("Validate").Return(nil)
	mockUserStore.On("CreateUser", user).Return(nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	// Test
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	mockUserStore.AssertExpectations(t)
}
