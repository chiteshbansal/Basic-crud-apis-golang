package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	middleware "first-api/internal/middlewares"
	model "first-api/internal/models"
	route "first-api/internal/route"
	"first-api/internal/service"
	"first-api/internal/utils"
	"first-api/pkg/cache"
)

// MockRepo struct holds a mock.Mock field to mock the repository.SongRepo interface. It helps in testing controller functions by mocking the associated helper functions of the repo layer.
type MockRepo struct {
	mock.Mock
}

// GetAllUsers mocks the GetAllUsers method of the repository.SongRepo interface.
func (m *MockRepo) GetAllUsers(b *[]model.User) error {
	args := m.Called(b)
	return args.Error(0)
}

// CreateUser mocks the CreateUser method of the repository.SongRepo interface.
func (m *MockRepo) CreateUser(b *model.User) error {
	args := m.Called(b)
	return args.Error(0)
}

// GetUser mocks the GetUser method of the repository.SongRepo interface.
func (m *MockRepo) GetUser(user *model.User, id string) (err error) {
	args := m.Called(user, id)

	if args.Error(0) != nil {
		return args.Error(0)
	}

	user.Id = 2
	user.Name = "test user"
	user.Email = "test@gmail.com"
	user.Phone = "9999999999"
	user.Address = "abcd efgh ijkl"

	return nil
}

// UpdateUser mocks the UpdateUser method of the repository.SongRepo interface.
func (m *MockRepo) UpdateUser(b *model.User, id string) (err error) {
	args := m.Called(b, id)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method of the repository.SongRepo interface.
func (m *MockRepo) DeleteUser(b *model.User, id string) (err error) {
	args := m.Called(b, id)
	return args.Error(0)
}

// initializeTest instantiates a MockRepo and creates a new Controller with this MockRepo as its Repo field. It also creates a new default gin.Engine and returns all three.
func initializeTest() (*MockRepo, service.User, *gin.Engine) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	gin.SetMode(gin.TestMode)

	mockRepo := new(MockRepo)
	userService := service.User{
		Store:     mockRepo,
		UserCache: cache.GetRateLimiterCache(),
	}

	return mockRepo, userService, gin.Default()
}

// // TestGetSongById function tests the GetSongById function of Controller
type Response struct {
	Status int        `json:"status"`
	User   model.User `json:"user"`
}

// TestGetAllUsers function tests the GetAllUsers function of Controller.
func TestGetAllUsers(t *testing.T) {
	users := []model.User{
		{Name: "test User 1", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"},
		{Name: "test user 2", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"},
	}

	mockRepo, userService, router := initializeTest()
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "GET",
		Handler:     userService.GetUsers,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) { middleware.VerifyJWT(ctx, userService.UserCache, "user") }},
	})

	route.InitializeRoutes(router)
	mockRepo.On("GetAllUsers", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]model.User)
		*arg = users
	})

	req, _ := http.NewRequest("GET", "/v1/user", nil)
	token, _ := utils.GenerateJWT(&model.User{Email: "test@test.com", Role: "admin"})
	token = "Bearer " + token

	req.Header.Set("Authorization", token)
	req.Header.Set("X-User-Email", "test@test.com")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	mockRepo.AssertExpectations(t)
}

// TestCreateUser function tests the CreateUser function of Controller.
func TestCreateUser(t *testing.T) {
	user := &model.User{
		Name:     "test user",
		Email:    "test@gmail.com",
		Phone:    "9999999999",
		Address:  "abcd efgh ijkl",
		Password: "pass123",
	}

	mockRepo, userService, router := initializeTest()
	mockRepo.On("CreateUser", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.User)
		*arg = *user
	})

	route.RegisterRoutes(route.RouteDef{
		Path:        "/user",
		Version:     "v1",
		Method:      "POST",
		Handler:     userService.CreateUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) { middleware.VerifyJWT(ctx, userService.UserCache, "user") }, middleware.ValidateUserData},
	})
	route.InitializeRoutes(router)
	AppReq, _ := route.StructToMapStringInterface(user)
	AppReq["confirmPassword"] = "pass123"

	body, _ := json.Marshal(AppReq)
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(body))
	token, _ := utils.GenerateJWT(&model.User{Email: "test@test.com", Role: "admin"})
	token = "Bearer " + token

	req.Header.Set("Authorization", token)
	req.Header.Set("X-User-Email", "test@test.com")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

// TestGetUser function tests the GetUser function of Controller.
func TestGetUser(t *testing.T) {
	mockRepo, userService, router := initializeTest()
	user := &model.User{Id: 2, Name: "test user", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"}

	route.RegisterRoutes(route.RouteDef{
		Path:        "/user/filter",
		Version:     "v1",
		Method:      "GET",
		Handler:     userService.GetUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) { middleware.VerifyJWT(ctx, userService.UserCache, "user") }},
	})
	route.InitializeRoutes(router)
	mockRepo.On("GetUser", mock.AnythingOfType("*model.User"), mock.AnythingOfType("string")).Return(nil)

	req, _ := http.NewRequest("GET", "/v1/user/filter?filter=id&value=2", nil)
	token, _ := utils.GenerateJWT(&model.User{Email: "test@test.com", Role: "admin"})
	token = "Bearer " + token

	req.Header.Set("Authorization", token)
	req.Header.Set("X-User-Email", "test@test.com")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var responseUser Response

	err := json.Unmarshal(bodyBytes, &responseUser)
	// If unmarshaling didn't return an error, check that the user fields match.
	if assert.NoError(t, err) {
		assert.Equal(t, user, &responseUser.User)
	}

	assert.Equal(t, http.StatusOK, responseUser.Status)
	mockRepo.AssertExpectations(t)
}

// TestUpdateUser function tests the UpdateUser function of Controller.
func TestUpdateUser(t *testing.T) {
	mockRepo, userService, router := initializeTest()
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user/:id",
		Version:     "v1",
		Method:      "PUT",
		Handler:     userService.UpdateUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) { middleware.VerifyJWT(ctx, userService.UserCache, "user") }, middleware.ValidateUserData},
	})
	route.InitializeRoutes(router)
	mockRepo.On("GetUser", mock.AnythingOfType("*model.User"), mock.AnythingOfType("string")).Return(nil)
	mockRepo.On("UpdateUser", mock.AnythingOfType("*model.User"), mock.AnythingOfType("string")).Return(nil)

	user := &model.User{Id: 1, Name: "test user", Email: "test@gmail.com", Phone: "9999999999", Address: "abcd efgh ijkl"}

	AppReq, _ := route.StructToMapStringInterface(user)

	body, _ := json.Marshal(AppReq)
	req, _ := http.NewRequest("PUT", "/v1/user/1", bytes.NewBuffer(body))
	token, _ := utils.GenerateJWT(&model.User{Email: "test@test.com", Role: "admin"})
	token = "Bearer " + token

	req.Header.Set("Authorization", token)
	req.Header.Set("X-User-Email", "test@test.com")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

// TestDeleteUser function tests the DeleteUser function of Controller.
func TestDeleteUser(t *testing.T) {
	mockRepo, userService, router := initializeTest()
	route.RegisterRoutes(route.RouteDef{
		Path:        "/user/:id",
		Version:     "v1",
		Method:      "DELETE",
		Handler:     userService.DeleteUser,
		Middlewares: []gin.HandlerFunc{func(ctx *gin.Context) { middleware.VerifyJWT(ctx, userService.UserCache, "admin") }},
	})
	route.InitializeRoutes(router)

	mockRepo.On("GetUser", mock.AnythingOfType("*model.User"), mock.AnythingOfType("string")).Return(nil)

	mockRepo.On("DeleteUser", mock.AnythingOfType("*model.User"), "1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/v1/user/1", nil)
	token, _ := utils.GenerateJWT(&model.User{Email: "test@test.com", Role: "admin"})
	token = "Bearer " + token

	req.Header.Set("Authorization", token)
	req.Header.Set("X-User-Email", "test@test.com")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var responseUser Response

	_ = json.Unmarshal(bodyBytes, &responseUser)

	assert.Equal(t, http.StatusOK, responseUser.Status)
	mockRepo.AssertExpectations(t)
}
