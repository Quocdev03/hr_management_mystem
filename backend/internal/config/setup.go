package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"chiquoc_hocgolang/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func CreateDatabase(cfg *DatabaseConfig) {
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent),
	})
	if err != nil {
		log.Fatalf("Không thể kết nối MySQL server: %v", err)
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	if err := tempDB.Exec(createDBQuery).Error; err != nil {
		log.Fatalf("Không thể tạo database: %v", err)
	}

	if sqlTempDB, err := tempDB.DB(); err == nil {
		sqlTempDB.Close()
	}
}

func RunMigrations(db *gorm.DB) {
	if err := db.AutoMigrate(&model.Role{}, &model.User{}); err != nil {
		log.Fatalf("Tạo bảng thất bại: %v", err)
	}

	origDisableFK := false
	if db.Config != nil {
		origDisableFK = db.Config.DisableForeignKeyConstraintWhenMigrating
		db.Config.DisableForeignKeyConstraintWhenMigrating = true
	}

	if err := db.AutoMigrate(&model.Department{}, &model.Employee{}); err != nil {
		log.Fatalf("Tạo bảng thất bại: %v", err)
	}

	if db.Config != nil {
		db.Config.DisableForeignKeyConstraintWhenMigrating = origDisableFK
	}

	if err := db.Migrator().CreateConstraint(&model.Employee{}, "User"); err != nil {
		log.Fatalf("Tạo constraint User cho Employee thất bại: %v", err)
	}
	if err := db.Migrator().CreateConstraint(&model.Employee{}, "Department"); err != nil {
		log.Fatalf("Tạo constraint Department cho Employee thất bại: %v", err)
	}
	if err := db.Migrator().CreateConstraint(&model.Department{}, "Manager"); err != nil {
		log.Fatalf("Tạo constraint Manager cho Department thất bại: %v", err)
	}
}

func SeedData(db *gorm.DB) {
	log.Println("Đang chèn dữ liệu!")

	err := db.Transaction(func(ctx *gorm.DB) error {
		if err := seedRoles(ctx); err != nil {
			return err
		}
		if err := seedDepartments(ctx); err != nil {
			return err
		}
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

func seedRoles(ctx *gorm.DB) error {
	roles := []model.Role{
		{Name: "admin", Description: "Quản trị viên - toàn quyền"},
		{Name: "hr", Description: "Nhân sự - quản lý nhân viên"},
		{Name: "employee", Description: "Nhân viên - xem thông tin"},
	}

	for _, role := range roles {
		var existing model.Role
		err := ctx.Where("name = ?", role.Name).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := ctx.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func seedDepartments(ctx *gorm.DB) error {
	departments := []model.Department{
		{Name: "Công nghệ thông tin", Code: "IT", Description: "Phát triển phần mềm"},
		{Name: "Nhân sự", Code: "HR", Description: "Quản lý nhân sự"},
		{Name: "Tài chính - Kế toán", Code: "FIN", Description: "Tài chính kế toán"},
		{Name: "Kinh doanh", Code: "SALES", Description: "Kinh doanh bán hàng"},
		{Name: "Marketing", Code: "MKT", Description: "Marketing & Branding"},
		{Name: "Chăm sóc khách hàng", Code: "CS", Description: "Hỗ trợ khách hàng"},
		{Name: "Vận hành", Code: "OPS", Description: "Quản lý vận hành"},
		{Name: "Pháp lý", Code: "LEGAL", Description: "Pháp chế doanh nghiệp"},
		{Name: "Mua hàng", Code: "PROC", Description: "Thu mua & cung ứng"},
		{Name: "R&D", Code: "RND", Description: "Nghiên cứu & phát triển"},
	}

	for _, d := range departments {
		if err := ctx.Where("code = ?", d.Code).FirstOrCreate(&d).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedUsersAndEmployees(ctx *gorm.DB) error {
	var adminRole, hrRole, empRole model.Role
	if err := ctx.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := ctx.Where("name = ?", "hr").First(&hrRole).Error; err != nil {
		return err
	}
	if err := ctx.Where("name = ?", "employee").First(&empRole).Error; err != nil {
		return err
	}

	hashedPwdAdmin, _ := bcrypt.GenerateFromPassword([]byte("12345678Quoc@"), bcrypt.DefaultCost)
	hashedPwdEmployee, _ := bcrypt.GenerateFromPassword([]byte("Password123@"), bcrypt.DefaultCost)

	// ================= USERS =================
	users := []model.User{
		{UserName: "admin", Email: "quocdt2003@ccquocn8n.cloud", Password: string(hashedPwdAdmin), RoleID: adminRole.ID},

		{UserName: "it_lead", Email: "it_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "it_dev", Email: "it_dev@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "hr_lead", Email: "hr_lead@example.com", Password: string(hashedPwdEmployee), RoleID: hrRole.ID},
		{UserName: "hr_recruiter", Email: "hr_recruiter@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "fin_lead", Email: "fin_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "fin_analyst", Email: "fin_analyst@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "sales_lead", Email: "sales_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "sales_rep1", Email: "sales_rep1@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "mkt_lead", Email: "mkt_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "mkt_exec", Email: "mkt_exec@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "cs_lead", Email: "cs_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "cs_agent", Email: "cs_agent@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "ops_lead", Email: "ops_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "ops_exec", Email: "ops_exec@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "legal_lead", Email: "legal_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "legal_exec", Email: "legal_exec@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "proc_lead", Email: "proc_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "proc_exec", Email: "proc_exec@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},

		{UserName: "rnd_lead", Email: "rnd_lead@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "rnd_researcher", Email: "rnd_researcher@example.com", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
	}

	// upsert users
	for i := range users {
		var existing model.User
		err := ctx.Unscoped().Where("user_name = ?", users[i].UserName).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := ctx.Create(&users[i]).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			if existing.DeletedAt.Valid {
				if err := ctx.Unscoped().Model(&existing).Update("deleted_at", nil).Error; err != nil {
					return err
				}
			}
			users[i] = existing
		}
	}

	// ================= DEPARTMENTS =================
	var itDept, hrDept, finDept, salesDept, mktDept, csDept, opsDept, legalDept, procDept, rndDept model.Department

	ctx.Where("code = ?", "IT").First(&itDept)
	ctx.Where("code = ?", "HR").First(&hrDept)
	ctx.Where("code = ?", "FIN").First(&finDept)
	ctx.Where("code = ?", "SALES").First(&salesDept)
	ctx.Where("code = ?", "MKT").First(&mktDept)
	ctx.Where("code = ?", "CS").First(&csDept)
	ctx.Where("code = ?", "OPS").First(&opsDept)
	ctx.Where("code = ?", "LEGAL").First(&legalDept)
	ctx.Where("code = ?", "PROC").First(&procDept)
	ctx.Where("code = ?", "RND").First(&rndDept)

	// ================= EMPLOYEES =================
	b1 := time.Date(1988, 5, 10, 0, 0, 0, 0, time.UTC)
	b2 := time.Date(1990, 8, 20, 0, 0, 0, 0, time.UTC)
	b3 := time.Date(1993, 11, 7, 0, 0, 0, 0, time.UTC)
	b4 := time.Date(1991, 3, 15, 0, 0, 0, 0, time.UTC)
	b5 := time.Date(1992, 7, 22, 0, 0, 0, 0, time.UTC)

	employees := []struct {
		model.Employee
		UserName string
	}{
		// IT
		{model.Employee{DepartmentID: itDept.ID, FirstName: "Chí", LastName: "Quốc", Phone: "0912345678", Position: "Head of Engineering", Salary: 47000000, JoinDate: time.Now(), BirthDate: &b1, Gender: "male", Status: "active"}, "admin"},
		{model.Employee{DepartmentID: itDept.ID, FirstName: "Hùng", LastName: "Nguyễn", Phone: "0912345679", Position: "Tech Lead", Salary: 35000000, JoinDate: time.Now(), BirthDate: &b2, Gender: "male", Status: "active"}, "it_lead"},
		{model.Employee{DepartmentID: itDept.ID, FirstName: "Lan", LastName: "Trần", Phone: "0912345680", Position: "Backend Developer", Salary: 28000000, JoinDate: time.Now(), BirthDate: &b3, Gender: "female", Status: "inactive"}, "it_dev"},

		// HR
		{model.Employee{DepartmentID: hrDept.ID, FirstName: "Hà", LastName: "Nguyễn", Phone: "0912345681", Position: "HR Manager", Salary: 33000000, JoinDate: time.Now(), BirthDate: &b2, Gender: "female", Status: "active"}, "hr_lead"},
		{model.Employee{DepartmentID: hrDept.ID, FirstName: "Mai", LastName: "Trần", Phone: "0912345682", Position: "Recruiter", Salary: 21000000, JoinDate: time.Now(), BirthDate: &b3, Gender: "female", Status: "active"}, "hr_recruiter"},

		// FIN
		{model.Employee{DepartmentID: finDept.ID, FirstName: "Anh", LastName: "Phạm", Phone: "0912345683", Position: "Finance Manager", Salary: 32000000, JoinDate: time.Now(), BirthDate: &b4, Gender: "male", Status: "active"}, "fin_lead"},

		// SALES
		{model.Employee{DepartmentID: salesDept.ID, FirstName: "Dũng", LastName: "Lê", Phone: "0912345684", Position: "Sales Manager", Salary: 31000000, JoinDate: time.Now(), BirthDate: &b5, Gender: "male", Status: "active"}, "sales_lead"},

		// MKT
		{model.Employee{DepartmentID: mktDept.ID, FirstName: "Hương", LastName: "Trần", Phone: "0912345685", Position: "Marketing Manager", Salary: 30000000, JoinDate: time.Now(), BirthDate: &b3, Gender: "female", Status: "active"}, "mkt_lead"},

		// CS
		{model.Employee{DepartmentID: csDept.ID, FirstName: "Vân", LastName: "Phạm", Phone: "0912345686", Position: "Customer Success Manager", Salary: 29000000, JoinDate: time.Now(), BirthDate: &b2, Gender: "female", Status: "active"}, "cs_lead"},

		// OPS
		{model.Employee{DepartmentID: opsDept.ID, FirstName: "Nam", LastName: "Đặng", Phone: "0912345687", Position: "Operations Manager", Salary: 30000000, JoinDate: time.Now(), BirthDate: &b1, Gender: "male", Status: "active"}, "ops_lead"},

		// LEGAL
		{model.Employee{DepartmentID: legalDept.ID, FirstName: "Minh", LastName: "Hoàng", Phone: "0912345688", Position: "Legal Manager", Salary: 30000000, JoinDate: time.Now(), BirthDate: &b5, Gender: "male", Status: "active"}, "legal_lead"},

		// PROC
		{model.Employee{DepartmentID: procDept.ID, FirstName: "Trâm", LastName: "Nguyễn", Phone: "0912345689", Position: "Procurement Manager", Salary: 30000000, JoinDate: time.Now(), BirthDate: &b4, Gender: "female", Status: "active"}, "proc_lead"},

		// RND
		{model.Employee{DepartmentID: rndDept.ID, FirstName: "Hải", LastName: "Bùi", Phone: "0912345690", Position: "R&D Manager", Salary: 30500000, JoinDate: time.Now(), BirthDate: &b1, Gender: "male", Status: "active"}, "rnd_lead"},
	}

	for i := range employees {
		var u model.User
		if err := ctx.Where("user_name = ?", employees[i].UserName).First(&u).Error; err == nil {
			employees[i].UserID = &u.ID
		}

		var existing model.Employee
		if err := ctx.Where("first_name = ? AND last_name = ? AND department_id = ?", employees[i].FirstName, employees[i].LastName, employees[i].DepartmentID).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := ctx.Create(&employees[i].Employee).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			update := map[string]interface{}{}

			if existing.UserID == nil && employees[i].UserID != nil {
				update["user_id"] = employees[i].UserID
			}
			if len(update) > 0 {
				if err := ctx.Model(&existing).Updates(update).Error; err != nil {
					return err
				}
			}
		}
	}

	// Gán trưởng phòng cho từng phòng ban nếu lead employee đã tồn tại
	managerAssignments := map[string]string{
		"IT":    "it_lead",
		"HR":    "hr_lead",
		"FIN":   "fin_lead",
		"SALES": "sales_lead",
		"MKT":   "mkt_lead",
		"CS":    "cs_lead",
		"OPS":   "ops_lead",
		"LEGAL": "legal_lead",
		"PROC":  "proc_lead",
		"RND":   "rnd_lead",
	}

	for code, userName := range managerAssignments {
		var u model.User
		if err := ctx.Where("user_name = ?", userName).First(&u).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return err
		}

		var manager model.Employee
		if err := ctx.Where("user_id = ?", u.ID).First(&manager).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return err
		}

		if err := ctx.Model(&model.Department{}).Where("code = ?", code).Update("manager_id", manager.ID).Error; err != nil {
			return err
		}
	}

	return nil
}
