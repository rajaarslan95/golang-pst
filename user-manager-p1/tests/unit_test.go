package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"user-manager/handlers"
	"user-manager/schemas"
	"user-manager/services"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserStore implementing the methods for the interface
type mockUserStore struct {
	mock.Mock
}

func (m *mockUserStore) AddUser(user schemas.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserStore) GetUser(id int) (schemas.User, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(schemas.User); ok {
		return user, args.Error(1)
	}
	return schemas.User{}, args.Error(1)
}

func (m *mockUserStore) UpdateUser(user schemas.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserStore) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Setup mock service
func setupMockService() services.UserService {
	mockStore := new(mockUserStore)
	userService := services.UserService{
		UserStore: mockStore,
	}
	return userService
}

// Setup Handler with mock service
func setupHandlerWithMockService() *handlers.Handler {
	mockSvc := setupMockService()
	ch := make(chan string, 10)
	var wg sync.WaitGroup
	return handlers.NewUserHandler(mockSvc, ch, &wg)
}

// Test AddUser handler
func TestAddUser(t *testing.T) {
	handler := setupHandlerWithMockService()

	mockUser := schemas.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	// Mock the expected behavior
	handler.GetSvc().UserStore.(*mockUserStore).On("AddUser", mockUser).Return(nil)

	// Create HTTP request and recorder
	body, _ := json.Marshal(mockUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call the handler
	handler.AddUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	handler.GetSvc().UserStore.(*mockUserStore).AssertExpectations(t)
}

// Test UpdateUser handler
func TestUpdateUser(t *testing.T) {
	handler := setupHandlerWithMockService()

	mockUser := schemas.User{
		ID:    1,
		Name:  "Updated User",
		Email: "updated@example.com",
		Age:   35,
	}

	// Mock the expected behavior
	handler.GetSvc().UserStore.(*mockUserStore).On("GetUser", mockUser.ID).Return(mockUser, nil)
	handler.GetSvc().UserStore.(*mockUserStore).On("UpdateUser", mockUser).Return(nil)

	// Create HTTP request and recorder
	body, _ := json.Marshal(mockUser)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(mockUser.ID)})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call the handler
	handler.UpdateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	handler.GetSvc().UserStore.(*mockUserStore).AssertExpectations(t)
}

// Test DeleteUser handler
func TestDeleteUser(t *testing.T) {
	handler := setupHandlerWithMockService()

	mockID := 1

	// Mock the expected behavior
	handler.GetSvc().UserStore.(*mockUserStore).On("GetUser", mockID).Return(schemas.User{ID: mockID}, nil)
	handler.GetSvc().UserStore.(*mockUserStore).On("DeleteUser", mockID).Return(nil)

	// Create HTTP request and recorder
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(mockID)})
	w := httptest.NewRecorder()

	// Call the handler
	handler.DeleteUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	handler.GetSvc().UserStore.(*mockUserStore).AssertExpectations(t)
}

// Test GetUser handler
func TestGetUser(t *testing.T) {
	handler := setupHandlerWithMockService()

	mockUser := schemas.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	// Mock the expected behavior
	handler.GetSvc().UserStore.(*mockUserStore).On("GetUser", mockUser.ID).Return(mockUser, nil)

	// Create HTTP request and recorder
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(mockUser.ID)})
	w := httptest.NewRecorder()

	// Call the handler
	handler.GetUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var user schemas.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)

	handler.GetSvc().UserStore.(*mockUserStore).AssertExpectations(t)
}
