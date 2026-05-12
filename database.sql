-- ============================================================
-- HRM Database Schema - SQL Server
-- Mini HRM System cho Fresher Backend Golang
-- Chỉ giữ phần schema + constraint cần thiết
-- Seed data sẽ được xử lý bằng Go (seed.go)
-- ============================================================


-- ============================================================
-- 1. TẠO DATABASE
-- ============================================================

USE master;
GO

-- Xóa database cũ nếu tồn tại (dùng khi dev/test)
IF EXISTS (SELECT name FROM sys.databases WHERE name = N'hrm_db')
BEGIN
    ALTER DATABASE hrm_db SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
    DROP DATABASE hrm_db;
END
GO

-- Tạo database mới
CREATE DATABASE hrm_db
COLLATE Vietnamese_CI_AS;
GO

USE hrm_db;
GO


-- ============================================================
-- 2. BẢNG permissions
-- Lưu danh sách quyền trong hệ thống
-- Ví dụ: employee:read, employee:create
-- ============================================================

CREATE TABLE permissions (
    id          INT IDENTITY(1,1) PRIMARY KEY,
    name        NVARCHAR(100) NOT NULL UNIQUE,
    description NVARCHAR(255),
    created_at  DATETIME2 DEFAULT GETDATE(),
    updated_at  DATETIME2 DEFAULT GETDATE(),
    deleted_at  DATETIME2 NULL
);
GO


-- ============================================================
-- 3. BẢNG roles
-- Vai trò: admin | hr | employee
-- ============================================================

CREATE TABLE roles (
    id          INT IDENTITY(1,1) PRIMARY KEY,
    name        NVARCHAR(50) NOT NULL UNIQUE,
    description NVARCHAR(255),
    created_at  DATETIME2 DEFAULT GETDATE(),
    updated_at  DATETIME2 DEFAULT GETDATE(),
    deleted_at  DATETIME2 NULL
);
GO


-- ============================================================
-- 4. BẢNG role_permissions
-- Junction table - quan hệ nhiều-nhiều giữa roles và permissions
-- Không cần model Go riêng, GORM tự quản lý qua tag many2many
-- ============================================================

CREATE TABLE role_permissions (
    role_id       INT NOT NULL,
    permission_id INT NOT NULL,

    PRIMARY KEY (role_id, permission_id),

    CONSTRAINT FK_role_permissions_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT FK_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
);
GO


-- ============================================================
-- 5. BẢNG departments
-- Thông tin phòng ban
-- manager_id FK tới employees, thêm sau khi tạo xong bảng employees
-- ============================================================

CREATE TABLE departments (
    id          INT IDENTITY(1,1) PRIMARY KEY,
    name        NVARCHAR(100) NOT NULL UNIQUE,
    code        NVARCHAR(20)  NOT NULL UNIQUE,
    description NVARCHAR(255),
    manager_id  INT NULL,
    created_at  DATETIME2 DEFAULT GETDATE(),
    updated_at  DATETIME2 DEFAULT GETDATE(),
    deleted_at  DATETIME2 NULL
);
GO


-- ============================================================
-- 6. BẢNG users
-- Tài khoản đăng nhập hệ thống
-- Password lưu dạng bcrypt hash, KHÔNG lưu plaintext
-- ============================================================

CREATE TABLE users (
    id          INT IDENTITY(1,1) PRIMARY KEY,
    username    NVARCHAR(100) NOT NULL UNIQUE,
    email       NVARCHAR(150) NOT NULL UNIQUE,
    password    NVARCHAR(255) NOT NULL,
    role_id     INT NOT NULL,
    is_active   BIT DEFAULT 1,
    created_at  DATETIME2 DEFAULT GETDATE(),
    updated_at  DATETIME2 DEFAULT GETDATE(),
    deleted_at  DATETIME2 NULL,

    CONSTRAINT FK_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
);
GO


-- ============================================================
-- 7. BẢNG employees
-- Thông tin chi tiết nhân viên
-- ============================================================

CREATE TABLE employees (
    id            INT IDENTITY(1,1) PRIMARY KEY,
    user_id       INT NULL,       -- NULL nếu chưa có tài khoản hệ thống
    department_id INT NOT NULL,
    first_name    NVARCHAR(100) NOT NULL,
    last_name     NVARCHAR(100) NOT NULL,
    email         NVARCHAR(150) NOT NULL UNIQUE,
    phone         NVARCHAR(20),
    position      NVARCHAR(100),
    salary        DECIMAL(15,2) DEFAULT 0,
    join_date     DATE          DEFAULT CAST(GETDATE() AS DATE),
    status        NVARCHAR(20)  DEFAULT 'active',  -- active | inactive | resigned
    created_at    DATETIME2     DEFAULT GETDATE(),
    updated_at    DATETIME2     DEFAULT GETDATE(),
    deleted_at    DATETIME2     NULL,

    -- Chỉ cho phép status hợp lệ
    CONSTRAINT CK_employees_status
        CHECK (status IN ('active', 'inactive', 'resigned')),

    -- Lương không được âm
    CONSTRAINT CK_employees_salary
        CHECK (salary >= 0),

    -- Xoá user → set user_id = NULL, không xoá employee
    CONSTRAINT FK_employees_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    -- Không cho xoá department khi còn nhân viên
    CONSTRAINT FK_employees_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
);
GO


-- ============================================================
-- 8. FK: departments.manager_id → employees.id
-- Thêm sau khi bảng employees đã tồn tại
-- (tránh circular dependency khi CREATE TABLE)
-- ============================================================

ALTER TABLE departments
    ADD CONSTRAINT FK_departments_manager
        FOREIGN KEY (manager_id)
        REFERENCES employees(id)
        ON DELETE SET NULL;  -- xoá employee → manager_id = NULL, không xoá dept
GO


-- ============================================================
-- 9. INDEX CƠ BẢN
-- ============================================================

-- Hay JOIN employees theo department
CREATE INDEX IX_employees_department_id ON employees (department_id);
GO

-- Hay JOIN users theo role
CREATE INDEX IX_users_role_id ON users (role_id);
GO

-- Soft delete: hầu hết query đều lọc WHERE deleted_at IS NULL
CREATE INDEX IX_employees_deleted_at   ON employees   (deleted_at);
CREATE INDEX IX_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IX_users_deleted_at       ON users       (deleted_at);
GO


-- ============================================================
-- 10. KIỂM TRA KẾT QUẢ
-- ============================================================

SELECT TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_TYPE = 'BASE TABLE'
ORDER BY TABLE_NAME;
GO

PRINT N'';
PRINT N'✅ HRM Database schema created successfully!';
PRINT N'💡 Seed data sẽ được tạo bằng Go (config/seed.go)';
GO