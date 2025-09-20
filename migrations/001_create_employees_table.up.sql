CREATE TABLE employees (
                           id SERIAL PRIMARY KEY,
                           full_name VARCHAR(100) NOT NULL,
                           phone VARCHAR(20) NOT NULL,
                           city VARCHAR(50) NOT NULL,
                           created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                           updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);