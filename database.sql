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
-- Ví dụ:
-- employee:read
-- employee:create
-- ============================================================

CREATE TABLE permissions (
    id INT IDENTITY(1,1) PRIMARY KEY,

    name NVARCHAR(100) NOT NULL UNIQUE,

    description NVARCHAR(255),

    created_at DATETIME2 DEFAULT GETDATE(),

    updated_at DATETIME2 DEFAULT GETDATE(),

    deleted_at DATETIME2 NULL
);
GO


-- ============================================================
-- 3. BẢNG roles
-- Vai trò:
-- admin
-- hr
-- employee
-- ============================================================

CREATE TABLE roles (
    id INT IDENTITY(1,1) PRIMARY KEY,

    name NVARCHAR(50) NOT NULL UNIQUE,

    description NVARCHAR(255),

    created_at DATETIME2 DEFAULT GETDATE(),

    updated_at DATETIME2 DEFAULT GETDATE(),

    deleted_at DATETIME2 NULL
);
GO


-- ============================================================
-- 4. BẢNG role_permissions
-- Quan hệ nhiều-nhiều giữa roles và permissions
-- ============================================================

CREATE TABLE role_permissions (
    role_id INT NOT NULL,

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
-- ============================================================

CREATE TABLE departments (
    id INT IDENTITY(1,1) PRIMARY KEY,

    name NVARCHAR(100) NOT NULL UNIQUE,

    code NVARCHAR(20) NOT NULL UNIQUE,

    description NVARCHAR(255),

    created_at DATETIME2 DEFAULT GETDATE(),

    updated_at DATETIME2 DEFAULT GETDATE(),

    deleted_at DATETIME2 NULL
);
GO


-- ============================================================
-- 6. BẢNG users
-- Tài khoản đăng nhập hệ thống
-- Password lưu dạng bcrypt hash
-- ============================================================

CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,

    username NVARCHAR(100) NOT NULL UNIQUE,

    email NVARCHAR(150) NOT NULL UNIQUE,

    password NVARCHAR(255) NOT NULL,

    role_id INT NOT NULL,

    is_active BIT DEFAULT 1,

    created_at DATETIME2 DEFAULT GETDATE(),

    updated_at DATETIME2 DEFAULT GETDATE(),

    deleted_at DATETIME2 NULL,

    CONSTRAINT FK_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
);
GO


-- ============================================================
-- 7. BẢNG employees
-- Thông tin nhân viên
-- ============================================================

CREATE TABLE employees (
    id INT IDENTITY(1,1) PRIMARY KEY,

    user_id INT NULL,

    department_id INT NOT NULL,

    first_name NVARCHAR(100) NOT NULL,

    last_name NVARCHAR(100) NOT NULL,

    email NVARCHAR(150) NOT NULL UNIQUE,

    phone NVARCHAR(20),

    position NVARCHAR(100),

    salary DECIMAL(15,2) DEFAULT 0,

    join_date DATE DEFAULT GETDATE(),

    -- active | inactive | resigned
    status NVARCHAR(20) DEFAULT 'active',

    created_at DATETIME2 DEFAULT GETDATE(),

    updated_at DATETIME2 DEFAULT GETDATE(),

    deleted_at DATETIME2 NULL,

    -- Chỉ cho phép status hợp lệ
    CONSTRAINT CK_employees_status
        CHECK (status IN ('active', 'inactive', 'resigned')),

    -- Lương không được âm
    CONSTRAINT CK_employees_salary
        CHECK (salary >= 0),

    CONSTRAINT FK_employees_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT FK_employees_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
);
GO


-- ============================================================
-- 8. BẢNG leave_requests
-- Đơn xin nghỉ phép
-- ============================================================

CREATE TABLE leave_requests (
    id INT IDENTITY(1,1) PRIMARY KEY,

    employee_id INT NOT NULL,

    leave_type NVARCHAR(50) NOT NULL,

    start_date DATE NOT NULL,

    end_date DATE NOT NULL,

    reason NVARCHAR(500),

    -- pending | approved | rejected
    status NVARCHAR(20) DEFAULT 'pending',

    approved_by INT NULL,

    created_at DATETIME2 DEFAULT GETDATE(),

    updated_at DATETIME2 DEFAULT GETDATE(),

    deleted_at DATETIME2 NULL,

    -- Validate loại nghỉ phép
    CONSTRAINT CK_leave_requests_type
        CHECK (
            leave_type IN (
                'annual',
                'sick',
                'unpaid'
            )
        ),

    -- Validate trạng thái
    CONSTRAINT CK_leave_requests_status
        CHECK (
            status IN (
                'pending',
                'approved',
                'rejected'
            )
        ),

    -- end_date phải >= start_date
    CONSTRAINT CK_leave_requests_dates
        CHECK (end_date >= start_date),

    CONSTRAINT FK_leave_requests_employee
        FOREIGN KEY (employee_id)
        REFERENCES employees(id),

    CONSTRAINT FK_leave_requests_approver
        FOREIGN KEY (approved_by)
        REFERENCES employees(id)
);
GO


-- ============================================================
-- 9. INDEX CƠ BẢN
-- Chỉ giữ index thường dùng
-- ============================================================

-- Query employee theo department
CREATE INDEX IX_employees_department_id
ON employees(department_id);
GO

-- Query user theo role
CREATE INDEX IX_users_role_id
ON users(role_id);
GO

-- Query leave request theo employee
CREATE INDEX IX_leave_requests_employee_id
ON leave_requests(employee_id);
GO


-- ============================================================
-- 10. KIỂM TRA KẾT QUẢ
-- ============================================================

PRINT N'';
PRINT N'✅ HRM Database schema created successfully!';
PRINT N'💡 Seed data sẽ được tạo bằng Go (config/seed.go)';
GO