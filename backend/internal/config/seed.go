package config

import (
	"chiquoc_hocgolang/internal/model"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	log.Println("Đang chèn dữ liệu!")

	// Nếu có lỗi giữa chừng, rollback toàn bộ
	err := db.Transaction(func(ctx *gorm.DB) error {
		// 1. Tạo data bảng Permissions
		if err := seedPermissions(ctx); err != nil {
			return err
		}

		// 2. Tạo data bảng Roles
		if err := seedRoles(ctx); err != nil {
			return err
		}

		// 3. Tạo data bảng Departments
		if err := seedDepartments(ctx); err != nil {
			return err
		}

		// 4. Tạo data bảng Users
		if err := seedUsers(ctx); err != nil {
			return err
		}

		// 5. Tạo data bảng Employees
		if err := seedEmployees(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Chèn dữ liệu thất bại!: %v", err)
		return
	}

	log.Println("Chèn dữ liệu thành công")
}

func seedPermissions(ctx *gorm.DB) error {
	permissions := []model.Permission{
		{Name: "employee:read", Description: "Xem danh sách nhân viên"},
		{Name: "employee:create", Description: "Tạo nhân viên mới"},
		{Name: "employee:update", Description: "Cập nhật thông tin nhân viên"},
		{Name: "employee:delete", Description: "Xóa nhân viên"},
		{Name: "department:read", Description: "Xem danh sách phòng ban"},
		{Name: "department:create", Description: "Tạo phòng ban mới"},
		{Name: "department:update", Description: "Cập nhật phòng ban"},
		{Name: "department:delete", Description: "Xóa phòng ban"},
		{Name: "user:manage", Description: "Quản lý tài khoản người dùng"},
	}
	for _, p := range permissions {
		// Tạo nếu chưa có
		if err := ctx.Where("name = ?", p.Name).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedRoles(ctx *gorm.DB) error {
	// Lấy tất cả permission gán cho admin
	var allPermissions []model.Permission
	ctx.Find(&allPermissions)

	// Lấy permission cho HR
	var hrPermissions []model.Permission
	ctx.Where("name IN ?", []string{
		"employee:read", "employee:create", "employee:update",
		"department:read",
	}).Find(&hrPermissions)

	// Lấy permission cho nhân viên thường
	var empPermissions []model.Permission
	ctx.Where("name In ?", []string{
		"employee:read", "department:read",
	}).Find(&empPermissions)

	roles := []struct {
		role        model.Role
		permissions []model.Permission
	}{
		{model.Role{Name: "admin", Description: "Quản trị viên - toàn quyền"}, allPermissions},
		{model.Role{Name: "hr", Description: "Nhân sự - quản lý nhân viên"}, hrPermissions},
		{model.Role{Name: "employee", Description: "Nhân viên - xem thông tin"}, empPermissions},
	}

	for _, r := range roles {
		var existing model.Role
		if ctx.Where("name = ?", r.role.Name).First(&existing).Error == nil {
			continue
		}
		r.role.Permission = r.permissions
		if err := ctx.Create(&r.role).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDepartments(ctx *gorm.DB) error {
	departments := []model.Department{
		{Name: "Công nghệ thông tin", Code: "IT", Description: "Phòng phát triển phần mềm"},
		{Name: "Nhân sự", Code: "HR", Description: "Phòng quản lý nhân sự"},
		{Name: "Tài chính - Kế toán", Code: "FIN", Description: "Phòng tài chính kế toán"},
		{Name: "Kinh doanh", Code: "SALES", Description: "Phòng kinh doanh và bán hàng"},
	}
	for _, d := range departments {
		ctx.Where("code = ?", d.Code).FirstOrCreate(&d)
	}
	return nil
}

func seedUsers(ctx *gorm.DB) error {
	// Tìm Role ID
	var adminRole, hrRole, empRole model.Role
	ctx.Where("name = ?", "admin").First(&adminRole)
	ctx.Where("name = ?", "hr").First(&hrRole)
	ctx.Where("name = ?", "employee").First(&empRole)

	// Hash pass "password123" cho user mẫu
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []model.User{
		{UserName: "admin", Email: "admin3212@hrm.com", Password: string(hashedPwd), RoleID: adminRole.ID},
		{UserName: "hr_manager", Email: "hr3212@hrm.com", Password: string(hashedPwd), RoleID: hrRole.ID},
		{UserName: "john_doe", Email: "john3212@hrm.com", Password: string(hashedPwd), RoleID: empRole.ID},
		{UserName: "jane_doe", Email: "jane3212@hrm.com", Password: string(hashedPwd), RoleID: empRole.ID},
	}

	for _, u := range users {
		ctx.Where("username = ?", u.UserName).FirstOrCreate(&u)
	}

	return nil
}

func seedEmployees(ctx *gorm.DB) error {
	// tìm id phòng ban
	var itDept, hrDept model.Department
	ctx.Where("code = ?", "IT").First(&itDept)
	ctx.Where("code = ?", "HR").First(&hrDept)

	joinDate := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	employees := []model.Employee{
		{
			DepartmentID: itDept.ID,
			FirstName:    "Nguyễn",
			LastName:     "Văn A",
			Email:        "nguyenvana@hrm.com",
			Phone:        "0901234567",
			Position:     "Backend Developer",
			Salary:       15000000,
			JoinDate:     joinDate,
			Status:       "active",
		},
		{
			DepartmentID: itDept.ID,
			FirstName:    "Trần",
			LastName:     "Thị B",
			Email:        "tranthib@hrm.com",
			Phone:        "0901234568",
			Position:     "Frontend Developer",
			Salary:       14000000,
			JoinDate:     joinDate,
			Status:       "active",
		},
		{
			DepartmentID: hrDept.ID,
			FirstName:    "Lê",
			LastName:     "Văn C",
			Email:        "levanc@hrm.com",
			Phone:        "0901234569",
			Position:     "HR Specialist",
			Salary:       12000000,
			JoinDate:     joinDate,
			Status:       "active",
		},
	}

	for _, e := range employees {
		ctx.Where("email = ?", e.Email).FirstOrCreate(&e)
	}
	return nil
}
