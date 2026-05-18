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

		// 4. Tạo data bảng Users và Employees
		if err := seedUsersAndEmployees(ctx); err != nil {
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
		{model.Role{Name: "manager", Description: "Quản lý phòng ban"}, empPermissions},
		{model.Role{Name: "intern", Description: "Thực tập sinh"}, empPermissions},
	}

	for _, r := range roles {
		var existing model.Role
		if ctx.Where("name = ?", r.role.Name).First(&existing).Error == nil {
			continue
		}
		r.role.Permissions = r.permissions
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

func seedUsersAndEmployees(ctx *gorm.DB) error {
	// Tìm Role ID
	var adminRole, hrRole, empRole model.Role
	ctx.Where("name = ?", "admin").First(&adminRole)
	ctx.Where("name = ?", "hr").First(&hrRole)
	ctx.Where("name = ?", "employee").First(&empRole)

	// Hash pass
	hashedPwdDefault, _ := bcrypt.GenerateFromPassword([]byte("Password123@"), bcrypt.DefaultCost)
	hashedPwdQuoc, _ := bcrypt.GenerateFromPassword([]byte("12345678Quoc@"), bcrypt.DefaultCost)

	users := []model.User{
		{UserName: "admin", Email: "admin3212@hrm.com", Password: string(hashedPwdDefault), RoleID: adminRole.ID},
		{UserName: "quocdev03", Email: "quocdev03@gmail.com", Password: string(hashedPwdQuoc), RoleID: adminRole.ID}, // Admin thứ 2
		{UserName: "hr_manager", Email: "hr3212@hrm.com", Password: string(hashedPwdDefault), RoleID: hrRole.ID},
		{UserName: "john_doe", Email: "john3212@hrm.com", Password: string(hashedPwdDefault), RoleID: empRole.ID},
		{UserName: "jane_doe", Email: "jane3212@hrm.com", Password: string(hashedPwdDefault), RoleID: empRole.ID},
		{UserName: "no_employee", Email: "noemp@company.com", Password: string(hashedPwdDefault), RoleID: empRole.ID},
	}

	for i := range users {
		var existing model.User
		if err := ctx.Where("email = ?", users[i].Email).First(&existing).Error; err != nil {
			ctx.Create(&users[i])
		} else {
			users[i] = existing
		}
	}

	// tìm id phòng ban
	var itDept, hrDept model.Department
	ctx.Where("code = ?", "IT").First(&itDept)
	ctx.Where("code = ?", "HR").First(&hrDept)

	joinDate := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	quocJoinDate := time.Date(2023, 4, 6, 0, 0, 0, 0, time.UTC)

	quocBirthDate := time.Date(2003, 4, 6, 0, 0, 0, 0, time.UTC)
	defaultBirthDate := time.Date(1995, 5, 20, 0, 0, 0, 0, time.UTC)

	employees := []model.Employee{
		// ===== ADMIN / CORE =====
		{
			DepartmentID: itDept.ID,
			FirstName:    "Chí",
			LastName:     "Quốc",
			Email:        "quocdev03@gmail.com",
			Phone:        "0825218643",
			Position:     "Senior Backend Engineer",
			Salary:       30000000,
			JoinDate:     quocJoinDate,
			BirthDate:    &quocBirthDate,
			Gender:       "male",
			Status:       "active",
		},

		// ===== IT TEAM =====
		{
			DepartmentID: itDept.ID,
			FirstName:    "Trần",
			LastName:     "Frontend",
			Email:        "frontend@company.com",
			Phone:        "0900000001",
			Position:     "Frontend Developer",
			Salary:       15000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: itDept.ID,
			FirstName:    "Nguyễn",
			LastName:     "Backend",
			Email:        "backend@company.com",
			Phone:        "0900000002",
			Position:     "Backend Developer",
			Salary:       18000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "female",
			Status:       "active",
		},
		{
			DepartmentID: itDept.ID,
			FirstName:    "Tester",
			LastName:     "QA",
			Email:        "qa@company.com",
			Phone:        "0900000003",
			Position:     "QA Engineer",
			Salary:       12000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "female",
			Status:       "inactive", // ⭐ test case
		},

		// ===== HR TEAM =====
		{
			DepartmentID: hrDept.ID,
			FirstName:    "HR",
			LastName:     "Manager",
			Email:        "hr_manager@company.com",
			Phone:        "0900000004",
			Position:     "HR Manager",
			Salary:       25000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "female",
			Status:       "active",
		},
		{
			DepartmentID: hrDept.ID,
			FirstName:    "Recruiter",
			LastName:     "Junior",
			Email:        "recruiter@company.com",
			Phone:        "0900000005",
			Position:     "HR Intern",
			Salary:       7000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "female",
			Status:       "active",
		},

		// ===== SALES TEAM =====
		{
			DepartmentID: 0, // ⭐ không có phòng ban
			FirstName:    "Sale",
			LastName:     "NoDept",
			Email:        "sales_nodept@company.com",
			Phone:        "0900000006",
			Position:     "Sales Executive",
			Salary:       10000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "male",
			Status:       "active",
		},

		// ===== EDGE CASE =====
		{
			DepartmentID: itDept.ID,
			FirstName:    "Duplicate",
			LastName:     "Name",
			Email:        "duplicate1@company.com",
			Phone:        "0900000007",
			Position:     "Developer",
			Salary:       13000000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "male",
			Status:       "inactive",
		},
		{
			DepartmentID: itDept.ID,
			FirstName:    "Duplicate",
			LastName:     "Name",
			Email:        "duplicate2@company.com",
			Phone:        "0900000008",
			Position:     "Developer",
			Salary:       13500000,
			JoinDate:     joinDate,
			BirthDate:    &defaultBirthDate,
			Gender:       "female",
			Status:       "active",
		},
	}

	for _, e := range employees {
		// ===== Set manager cho department =====

		// Lấy lại employee theo email (an toàn hơn hardcode ID)
		var itManager model.Employee
		ctx.Where("email = ?", "admin3212@hrm.com").First(&itManager)

		var hrManager model.Employee
		ctx.Where("email = ?", "hr3212@hrm.com").First(&hrManager)

		// Update manager_id
		if itManager.ID != 0 {
			ctx.Model(&model.Department{}).
				Where("code = ?", "IT").
				Update("manager_id", itManager.ID)
		}

		if hrManager.ID != 0 {
			ctx.Model(&model.Department{}).
				Where("code = ?", "HR").
				Update("manager_id", hrManager.ID)
		}

		// Tìm user có email tương ứng để gắn UserID
		var u model.User
		if err := ctx.Where("email = ?", e.Email).First(&u).Error; err == nil {
			e.UserID = &u.ID
		}

		var existing model.Employee
		if err := ctx.Where("email = ?", e.Email).First(&existing).Error; err != nil {
			ctx.Create(&e)
		} else {
			// Cập nhật UserID nếu chưa có
			if (existing.UserID == nil && e.UserID != nil) || existing.BirthDate == nil || existing.Gender == "" {
				updateFields := map[string]interface{}{}
				if existing.UserID == nil && e.UserID != nil {
					updateFields["user_id"] = e.UserID
				}
				if existing.BirthDate == nil && e.BirthDate != nil {
					updateFields["birth_date"] = e.BirthDate
				}
				if existing.Gender == "" && e.Gender != "" {
					updateFields["gender"] = e.Gender
				}
				if len(updateFields) > 0 {
					ctx.Model(&existing).Updates(updateFields)
				}
			}
		}
	}
	return nil
}
