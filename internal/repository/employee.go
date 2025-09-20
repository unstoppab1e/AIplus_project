package repository

import (
	"Aiplus_project/internal/model"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type EmployeeRepository interface {
	Create(employee *model.Employee) error
	GetByID(id int) (*model.Employee, error)
	GetAll() ([]*model.Employee, error)
	Update(employee *model.Employee) error
	Delete(id int) error
}

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *model.Employee) error {
	query := `
        INSERT INTO employees (full_name, phone, city, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	now := time.Now()
	employee.CreatedAt = now
	employee.UpdatedAt = now

	err := r.db.QueryRow(query, employee.FullName, employee.Phone, employee.City, now, now).Scan(&employee.ID)
	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	return nil
}

func (r *employeeRepository) GetByID(id int) (*model.Employee, error) {
	query := `
        SELECT id, full_name, phone, city, created_at, updated_at
        FROM employees
        WHERE id = $1`

	employee := &model.Employee{}
	err := r.db.QueryRow(query, id).Scan(
		&employee.ID,
		&employee.FullName,
		&employee.Phone,
		&employee.City,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return employee, nil
}

func (r *employeeRepository) GetAll() ([]*model.Employee, error) {
	query := `
        SELECT id, full_name, phone, city, created_at, updated_at
        FROM employees
        ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	defer rows.Close()

	var employees []*model.Employee
	for rows.Next() {
		employee := &model.Employee{}
		err := rows.Scan(
			&employee.ID,
			&employee.FullName,
			&employee.Phone,
			&employee.City,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *employeeRepository) Update(employee *model.Employee) error {
	query := `
        UPDATE employees
        SET full_name = $2, phone = $3, city = $4, updated_at = $5
        WHERE id = $1`

	employee.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, employee.ID, employee.FullName, employee.Phone, employee.City, employee.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

func (r *employeeRepository) Delete(id int) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}
