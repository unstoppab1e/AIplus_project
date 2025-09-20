package service

import (
	"Aiplus_project/internal/model"
	"Aiplus_project/internal/repository"
)

type EmployeeService interface {
	CreateEmployee(req *model.CreateEmployeeRequest) (*model.EmployeeResponse, error)
	GetEmployee(id int) (*model.EmployeeResponse, error)
	GetAllEmployees() ([]*model.EmployeeResponse, error)
	UpdateEmployee(id int, req *model.CreateEmployeeRequest) (*model.EmployeeResponse, error)
	DeleteEmployee(id int) error
}

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) CreateEmployee(req *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	employee := &model.Employee{
		FullName: req.FullName,
		Phone:    req.Phone,
		City:     req.City,
	}

	err := s.repo.Create(employee)
	if err != nil {
		return nil, err
	}

	return s.toEmployeeResponse(employee), nil
}

func (s *employeeService) GetEmployee(id int) (*model.EmployeeResponse, error) {
	employee, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toEmployeeResponse(employee), nil
}

func (s *employeeService) GetAllEmployees() ([]*model.EmployeeResponse, error) {
	employees, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*model.EmployeeResponse
	for _, employee := range employees {
		responses = append(responses, s.toEmployeeResponse(employee))
	}

	return responses, nil
}

func (s *employeeService) UpdateEmployee(id int, req *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	employee, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	employee.FullName = req.FullName
	employee.Phone = req.Phone
	employee.City = req.City

	err = s.repo.Update(employee)
	if err != nil {
		return nil, err
	}

	return s.toEmployeeResponse(employee), nil
}

func (s *employeeService) DeleteEmployee(id int) error {
	return s.repo.Delete(id)
}

func (s *employeeService) toEmployeeResponse(employee *model.Employee) *model.EmployeeResponse {
	return &model.EmployeeResponse{
		ID:        employee.ID,
		FullName:  employee.FullName,
		Phone:     employee.Phone,
		City:      employee.City,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}
}
