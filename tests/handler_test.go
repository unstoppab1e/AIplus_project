package tests

import (
	"bytes"
	"Aiplus_project/internal/handler"
	"Aiplus_project/internal/model"
	"Aiplus_project/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockEmployeeService — service.EmployeeService интерфейсін жалған түрде жүзеге асыру
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) CreateEmployee(req *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) GetEmployee(id int) (*model.EmployeeResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) GetAllEmployees() ([]*model.EmployeeResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) UpdateEmployee(id int, req *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.EmployeeResponse), args.Error(1)
}

func (m *MockEmployeeService) DeleteEmployee(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Қызметкерді сәтті құру тесті
func TestEmployeeHandler_CreateEmployee_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockEmployeeService)
	employeeHandler := handler.NewEmployeeHandler(mockService)

	router := gin.New()
	router.POST("/employees", employeeHandler.CreateEmployee)

	reqBody := model.CreateEmployeeRequest{
		FullName: "Нұржан Нұртаев",
		Phone:    "+7 701 123 45 67",
		City:     "Астана",
	}

	expectedResponse := &model.EmployeeResponse{
		ID:        1,
		FullName:  "Нұржан Нұртаев",
		Phone:     "+7 701 123 45 67",
		City:      "Астана",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("CreateEmployee", mock.MatchedBy(func(req *model.CreateEmployeeRequest) bool {
		return req.FullName == reqBody.FullName && req.Phone == reqBody.Phone && req.City == reqBody.City
	})).Return(expectedResponse, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.EmployeeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.FullName, response.FullName)
	assert.Equal(t, expectedResponse.Phone, response.Phone)
	assert.Equal(t, expectedResponse.City, response.City)

	mockService.AssertExpectations(t)
}

// JSON дұрыс емес болған кездегі тест
func TestEmployeeHandler_CreateEmployee_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockEmployeeService)
	employeeHandler := handler.NewEmployeeHandler(mockService)

	router := gin.New()
	router.POST("/employees", employeeHandler.CreateEmployee)

	// Қате JSON
	invalidJSON := `{"full_name": "Нұржан", "phone": ""` // жабылмаған жақша және бос телефон

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest,
