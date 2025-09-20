package main

import (
	"Aiplus_project/internal/config"
	"Aiplus_project/internal/handler"
	"Aiplus_project/internal/repository"
	"Aiplus_project/internal/service"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Запускаем миграции
	if err := runMigrations(db, cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Инициализируем зависимости
	employeeRepo := repository.NewEmployeeRepository(db)
	employeeService := service.NewEmployeeService(employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Настраиваем роутер
	router := setupRouter(employeeHandler)

	// Запускаем сервер
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func runMigrations(db *sql.DB, databaseURL string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func setupRouter(employeeHandler *handler.EmployeeHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Employee routes
	api := router.Group("/api/v1")
	{
		employees := api.Group("/employees")
		{
			employees.POST("", employeeHandler.CreateEmployee)
			employees.GET("", employeeHandler.GetAllEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)
		}
	}

	return router
}
