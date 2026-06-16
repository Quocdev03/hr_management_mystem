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
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(20)  NOT NULL,
    description VARCHAR(255),
    manager_id  INT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL
);

-- ============================================================
-- 4. BẢNG users
-- Tài khoản đăng nhập hệ thống
-- ============================================================

CREATE TABLE users (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    user_name   VARCHAR(100) NOT NULL,
    email       VARCHAR(150) NOT NULL,
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
-- 5. BẢNG permissions
-- Quản lý quyền theo hành động (RBAC)
-- ============================================================

CREATE TABLE permissions (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    code        VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  DATETIME NULL
);

CREATE TABLE role_permissions (
    role_id       INT NOT NULL,
    permission_id INT NOT NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME NULL,
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT FK_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT FK_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE user_permissions (
    user_id       INT NOT NULL,
    permission_id INT NOT NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME NULL,
    PRIMARY KEY (user_id, permission_id),
    CONSTRAINT FK_user_permissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT FK_user_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- ============================================================
-- 6. BẢNG employees
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
    gender        ENUM('male', 'female', 'other') DEFAULT 'male',
    status        ENUM('active', 'inactive') DEFAULT 'active',
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME NULL,

    CONSTRAINT CK_employees_salary
        CHECK (salary >= 0),

    CONSTRAINT FK_employees_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,

    CONSTRAINT FK_employees_department
        FOREIGN KEY (department_id) REFERENCES departments(id)
);

-- ============================================================
-- 7. FK: departments.manager_id → employees.id
-- Thêm sau khi bảng employees đã tồn tại
-- ============================================================

ALTER TABLE departments
    ADD CONSTRAINT FK_departments_manager
        FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL;

-- ============================================================
-- 8. INDEX CƠ BẢN VÀ COMPOSITE
-- ============================================================

CREATE INDEX IX_employees_department_id ON employees (department_id);
CREATE INDEX IX_users_role_id ON users (role_id);
CREATE INDEX IX_employees_deleted_at ON employees (deleted_at);
CREATE INDEX IX_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IX_users_deleted_at ON users (deleted_at);
CREATE INDEX IX_roles_deleted_at ON roles (deleted_at);

-- Indexes thay thế cho UNIQUE bị loại bỏ
CREATE INDEX IX_departments_name ON departments (name);
CREATE INDEX IX_departments_code ON departments (code);
CREATE INDEX IX_users_username ON users (user_name);
CREATE INDEX IX_users_email ON users (email);
CREATE INDEX IX_employees_user_id ON employees (user_id);

-- Composite indexes cho query thường dùng
CREATE INDEX IX_employees_dept_deleted ON employees (department_id, deleted_at);
CREATE INDEX IX_employees_status_deleted ON employees (status, deleted_at);