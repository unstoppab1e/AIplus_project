package tests

import (
	"Aiplus_project/internal/model"
	"Aiplus_project/internal/repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Жаңа қызметкерді құруды тексеру
func TestEmployeeRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	employee := &model.Employee{
		FullName: "Нұржан Нұртаев",
		Phone:    "+7 701 123 45 67",
		City:     "Астана",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO employees (full_name, phone, city, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
		WithArgs(employee.FullName, employee.Phone, employee.City, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.Create(employee)
	require.NoError(t, err)
	assert.Equal(t, 1, employee.ID)
	assert.False(t, employee.CreatedAt.IsZero())
	assert.False(t, employee.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ID бойынша табу (сәтті жағдай)
func TestEmployeeRepository_GetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	now := time.Now()
	expectedEmployee := &model.Employee{
		ID:        1,
		FullName:  "Нұржан Нұртаев",
		Phone:     "+7 701 123 45 67",
		City:      "Астана",
		CreatedAt: now,
		UpdatedAt: now,
	}

	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "city", "created_at", "updated_at"}).
		AddRow(expectedEmployee.ID, expectedEmployee.FullName, expectedEmployee.Phone,
			expectedEmployee.City, expectedEmployee.CreatedAt, expectedEmployee.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, full_name, phone, city, created_at, updated_at FROM employees WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	employee, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, expectedEmployee.ID, employee.ID)
	assert.Equal(t, expectedEmployee.FullName, employee.FullName)
	assert.Equal(t, expectedEmployee.Phone, employee.Phone)
	assert.Equal(t, expectedEmployee.City, employee.City)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ID бойынша табылмаса
func TestEmployeeRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, full_name, phone, city, created_at, updated_at FROM employees WHERE id = $1`)).
		WithArgs(999).
		WillReturnError(sqlmock.ErrCancelled) // sql.ErrNoRows-ты модельдеу

	employee, err := repo.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, employee)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Барлық қызметкерлерді шығару
func TestEmployeeRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "city", "created_at", "updated_at"}).
		AddRow(1, "Нұржан Нұртаев", "+7 701 123 45 67", "Астана", now, now).
		AddRow(2, "Айбек Әлиев", "+7 702 987 65 43", "Алматы", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, full_name, phone, city, created_at, updated_at FROM employees ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	employees, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, employees, 2)
	assert.Equal(t, "Нұржан Нұртаев", employees[0].FullName)
	assert.Equal(t, "Айбек Әлиев", employees[1].FullName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Жаңарту тесті
func TestEmployeeRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	employee := &model.Employee{
		ID:       1,
		FullName: "Ержан Сейтқазы",
		Phone:    "+7 705 111 22 33",
		City:     "Шымкент",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE employees SET full_name = $2, phone = $3, city = $4, updated_at = $5 WHERE id = $1`)).
		WithArgs(employee.ID, employee.FullName, employee.Phone, employee.City, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(employee)
	require.NoError(t, err)
	assert.False(t, employee.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Жою тесті (сәтті)
func TestEmployeeRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM employees WHERE id = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Жою тесті (қызметкер табылмады)
func TestEmployeeRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewEmployeeRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM employees WHERE id = $1`)).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 жол әсер етті

	err = repo.Delete(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "employee not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}
