-- ============================================================
-- HRM Database Schema - MySQL
-- Mini HRM System cho Fresher Backend Golang
-- Seed data sẽ được xử lý bằng Go (seed.go)
-- ============================================================

-- ============================================================
-- 1. TẠO DATABASE
-- ============================================================

DROP DATABASE IF EXISTS hrm_db;
CREATE DATABASE hrm_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE hrm_db;

-- ============================================================
-- 2. BẢNG roles
-- Vai trò: admin | hr | employee
-- ============================================================

CREATE TABLE roles (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL
);

-- ============================================================
-- 3. BẢNG departments
-- Thông tin phòng ban
-- ============================================================

CREATE TABLE departments (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    code        VARCHAR(20)  NOT NULL UNIQUE,
    description VARCHAR(255),
    manager_id  INT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL
);

-- ============================================================
-- 6. BẢNG users
-- Tài khoản đăng nhập hệ thống
-- ============================================================

CREATE TABLE users (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    user_name   VARCHAR(100) NOT NULL UNIQUE,
    email       VARCHAR(150) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role_id     INT NOT NULL,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL,

    CONSTRAINT FK_users_role
        FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- ============================================================
-- 7. BẢNG employees
-- Thông tin chi tiết nhân viên
-- ============================================================

CREATE TABLE employees (
    id            INT AUTO_INCREMENT PRIMARY KEY,
    user_id       INT NULL,       -- NULL nếu chưa có tài khoản hệ thống
    department_id INT NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    phone         VARCHAR(20),
    position      VARCHAR(100),
    salary        DECIMAL(15,2) DEFAULT 0,
    join_date     DATE DEFAULT (CURRENT_DATE),
    birth_date    DATE NULL,
    gender        VARCHAR(10)  DEFAULT 'male',
    status        VARCHAR(20)  DEFAULT 'active',  -- active | inactive
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME NULL,

    CONSTRAINT CK_employees_status
        CHECK (status IN ('active', 'inactive')),

    CONSTRAINT CK_employees_salary
        CHECK (salary >= 0),

    CONSTRAINT FK_employees_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,

    CONSTRAINT FK_employees_department
        FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- ============================================================
-- 8. FK: departments.manager_id → employees.id
-- Thêm sau khi bảng employees đã tồn tại
-- ============================================================

ALTER TABLE departments
    ADD CONSTRAINT FK_departments_manager
        FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL;

-- ============================================================
-- 9. INDEX CƠ BẢN
-- ============================================================

CREATE INDEX IX_employees_department_id ON employees (department_id);
CREATE INDEX IX_users_role_id ON users (role_id);
CREATE INDEX IX_employees_deleted_at ON employees (deleted_at);
CREATE INDEX IX_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IX_users_deleted_at ON users (deleted_at);

-- MySQL syntax for unique index on user_id where user_id IS NOT NULL:
-- MySQL 8.0 allows unique index containing NULLs by default, 
-- multiple NULLs are distinct in MySQL.
CREATE UNIQUE INDEX UQ_employees_user_id ON employees (user_id);