package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"chiquoc_hocgolang/internal/model"
	"chiquoc_hocgolang/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// CreateDatabase tạo database nếu chưa tồn tại. Idempotent — an toàn chạy nhiều lần.
func CreateDatabase(cfg *DatabaseConfig) {
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	var tempDB *gorm.DB
	var err error
	maxRetries := 15
	for i := 1; i <= maxRetries; i++ {
		tempDB, err = gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
			Logger: glogger.Default.LogMode(glogger.Silent),
		})
		if err == nil {
			break
		}
		utils.Info("[MIGRATE] Waiting for MySQL connection... (Attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		utils.Fatal("[MIGRATE] Cannot connect to MySQL server: %v", err)
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	if err := tempDB.Exec(createDBQuery).Error; err != nil {
		utils.Fatal("[MIGRATE] Cannot create database: %v", err)
	}

	if sqlTempDB, err := tempDB.DB(); err == nil {
		_ = sqlTempDB.Close()
	}
}

// RunMigrations chạy AutoMigrate. Gọi utils.Fatal nếu lỗi.
// Dùng trong các script cần panic on error.
func RunMigrations(db *gorm.DB) {
	if err := RunMigrationsWithError(db); err != nil {
		utils.Fatal("[MIGRATE] Migration failed: %v", err)
	}
}

// RunMigrationsWithError chạy AutoMigrate, trả về error (dùng trong cmd/migrate).
func RunMigrationsWithError(db *gorm.DB) error {
	utils.Info("[MIGRATE] Migrating auth tables (roles, permissions, users)...")
	if err := db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.UserPermission{}, &model.User{}); err != nil {
		return fmt.Errorf("migrate auth tables: %w", err)
	}

	utils.Info("[MIGRATE] Migrating business tables (departments, employees)...")
	origDisableFK := false
	if db.Config != nil {
		origDisableFK = db.DisableForeignKeyConstraintWhenMigrating
		db.DisableForeignKeyConstraintWhenMigrating = true
	}

	if err := db.AutoMigrate(&model.Department{}, &model.Employee{}); err != nil {
		return fmt.Errorf("migrate business tables: %w", err)
	}

	if db.Config != nil {
		db.DisableForeignKeyConstraintWhenMigrating = origDisableFK
	}

	// Tạo constraints chỉ nếu chưa tồn tại
	migratorInstance := db.Migrator()

	if !migratorInstance.HasConstraint(&model.Employee{}, "User") {
		if err := migratorInstance.CreateConstraint(&model.Employee{}, "User"); err != nil {
			return fmt.Errorf("create constraint Employee.User: %w", err)
		}
	}
	if !migratorInstance.HasConstraint(&model.Employee{}, "Department") {
		if err := migratorInstance.CreateConstraint(&model.Employee{}, "Department"); err != nil {
			return fmt.Errorf("create constraint Employee.Department: %w", err)
		}
	}
	if !migratorInstance.HasConstraint(&model.Department{}, "Manager") {
		if err := migratorInstance.CreateConstraint(&model.Department{}, "Manager"); err != nil {
			return fmt.Errorf("create constraint Department.Manager: %w", err)
		}
	}

	utils.Info("[MIGRATE] All migrations completed successfully.")
	return nil
}

// SeedData chạy seed. Gọi utils.Error nếu lỗi.
// Dùng khi không cần error propagation.
func SeedData(db *gorm.DB) {
	if err := SeedDataWithError(db); err != nil {
		utils.Error("[SEED] Seed failed: %v", err)
	}
}

// SeedDataWithError chạy seed, trả về error (dùng trong cmd/seed).
// Idempotent: chạy nhiều lần không tạo dữ liệu trùng.
func SeedDataWithError(db *gorm.DB) error {
	utils.Info("[SEED] Starting seed transaction...")

	err := db.Transaction(func(ctx *gorm.DB) error {
		utils.Info("[SEED] Seeding roles...")
		if err := seedRoles(ctx); err != nil {
			return fmt.Errorf("seed roles: %w", err)
		}

		utils.Info("[SEED] Seeding permissions...")
		if err := seedPermissions(ctx); err != nil {
			return fmt.Errorf("seed permissions: %w", err)
		}

		utils.Info("[SEED] Seeding role-permissions...")
		if err := seedRolePermissions(ctx); err != nil {
			return fmt.Errorf("seed role-permissions: %w", err)
		}

		utils.Info("[SEED] Seeding departments...")
		if err := seedDepartments(ctx); err != nil {
			return fmt.Errorf("seed departments: %w", err)
		}

		utils.Info("[SEED] Seeding users...")
		if err := seedUsers(ctx); err != nil {
			return fmt.Errorf("seed users: %w", err)
		}

		utils.Info("[SEED] Seeding employees...")
		if err := seedEmployees(ctx); err != nil {
			return fmt.Errorf("seed employees: %w", err)
		}

		utils.Info("[SEED] Seeding department managers...")
		if err := seedDepartmentManagers(ctx); err != nil {
			return fmt.Errorf("seed department managers: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	utils.Info("[SUCCESS] Seed completed.")
	utils.Info("=================================================================")
	utils.Info("🔑 DEMO ACCOUNTS:")
	utils.Info("   [Admin]    chiquoc23AD@company.vn  / chiquoc23AD")
	utils.Info("   [HR]       chiquoc23HR@company.vn  / chiquoc23HR")
	utils.Info("   [Employee] chiquoc23EMP@company.vn / chiquoc23EMP")
	utils.Info("=================================================================")
	return nil
}

// ─── Seed functions (all idempotent) ─────────────────────────────────────────

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
		// Role đã tồn tại → skip, không update
	}
	return nil
}

func seedPermissions(ctx *gorm.DB) error {
	permissions := []model.Permission{
		{Code: "employee.read", Description: "Xem danh sách nhân viên"},
		{Code: "employee.create", Description: "Tạo nhân viên"},
		{Code: "employee.update", Description: "Cập nhật nhân viên"},
		{Code: "employee.delete", Description: "Xóa nhân viên"},
		{Code: "user.read", Description: "Xem tài khoản"},
		{Code: "user.create", Description: "Tạo tài khoản"},
		{Code: "user.update", Description: "Cập nhật tài khoản"},
		{Code: "user.delete", Description: "Xóa tài khoản"},
		{Code: "department.read", Description: "Xem phòng ban"},
		{Code: "department.create", Description: "Tạo phòng ban"},
		{Code: "department.update", Description: "Cập nhật phòng ban"},
		{Code: "department.delete", Description: "Xóa phòng ban"},
	}
	for _, p := range permissions {
		if err := ctx.Where("code = ?", p.Code).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedRolePermissions(ctx *gorm.DB) error {
	var adminRole, hrRole, employeeRole model.Role
	if err := ctx.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := ctx.Where("name = ?", "hr").First(&hrRole).Error; err != nil {
		return err
	}
	if err := ctx.Where("name = ?", "employee").First(&employeeRole).Error; err != nil {
		return err
	}

	adminCodes := []string{"employee.read", "employee.create", "employee.update", "employee.delete", "user.read", "user.create", "user.update", "user.delete", "department.read", "department.create", "department.update", "department.delete"}
	hrCodes := []string{"employee.read", "employee.create", "employee.update", "department.read"}
	employeeCodes := []string{"employee.read", "department.read"}

	// Idempotent: dùng FirstOrCreate thay vì DELETE + INSERT
	// Điều này bảo toàn custom permissions đã được thêm bởi user
	assignIfNotExists := func(roleID uint, codes []string) error {
		for _, code := range codes {
			var perm model.Permission
			if err := ctx.Where("code = ?", code).First(&perm).Error; err != nil {
				return fmt.Errorf("permission not found: %s", code)
			}
			rp := model.RolePermission{RoleID: roleID, PermissionID: perm.ID}
			if err := ctx.Where(rp).FirstOrCreate(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	}

	if err := assignIfNotExists(adminRole.ID, adminCodes); err != nil {
		return err
	}
	if err := assignIfNotExists(hrRole.ID, hrCodes); err != nil {
		return err
	}
	if err := assignIfNotExists(employeeRole.ID, employeeCodes); err != nil {
		return err
	}

	return nil
}

func seedDepartments(ctx *gorm.DB) error {
	departments := []model.Department{
		{Name: "Công nghệ thông tin", Code: "IT", Description: "Phát triển phần mềm"},
		{Name: "Nhân sự", Code: "HR", Description: "Quản lý nhân sự"},
		{Name: "Tài chính - Kế toán", Code: "FIN", Description: "Tài chính kế toán"},
	}

	for _, d := range departments {
		if err := ctx.Where("code = ?", d.Code).FirstOrCreate(&d).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedUsers(ctx *gorm.DB) error {
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

	hashedPwdEmployeeDefault, _ := bcrypt.GenerateFromPassword([]byte("Emp2026@pass"), bcrypt.DefaultCost)
	hashedPwdAdminSpecial, _ := bcrypt.GenerateFromPassword([]byte("chiquoc23AD"), bcrypt.DefaultCost)
	hashedPwdHRSpecial, _ := bcrypt.GenerateFromPassword([]byte("chiquoc23HR"), bcrypt.DefaultCost)
	hashedPwdEmployeeSpecial, _ := bcrypt.GenerateFromPassword([]byte("chiquoc23EMP"), bcrypt.DefaultCost)

	roleIDs := map[string]uint{"admin": adminRole.ID, "hr": hrRole.ID, "employee": empRole.ID}

	for _, item := range getSeedEmployeeInfo() {
		if item.Employee.BirthDate == nil {
			return errors.New("birth date is required for seed user generation")
		}

		var username string
		var password []byte

		if username = resolveSeedUsername(item); username != "" {
			switch username {
			case "chiquoc23AD":
				password = hashedPwdAdminSpecial
			case "chiquoc23HR":
				password = hashedPwdHRSpecial
			case "chiquoc23EMP":
				password = hashedPwdEmployeeSpecial
			}
		} else {
			username = buildUserName(item.Employee.FirstName, item.Employee.LastName, *item.Employee.BirthDate)
			password = hashedPwdEmployeeDefault
		}

		email := buildEmail(username)

		user := model.User{
			UserName: username,
			Email:    email,
			Password: string(password),
			RoleID:   roleIDs[item.RoleName],
		}

		var existing model.User
		err := ctx.Unscoped().Where("user_name = ?", user.UserName).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := ctx.Create(&user).Error; err != nil {
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
		}
	}

	return nil
}

func seedEmployees(ctx *gorm.DB) error {
	departments := map[string]model.Department{}
	for _, code := range []string{"IT", "HR", "FIN"} {
		var d model.Department
		if err := ctx.Where("code = ?", code).First(&d).Error; err != nil {
			return err
		}
		departments[code] = d
	}

	for _, item := range getSeedEmployeeInfo() {
		dept, ok := departments[item.DepartmentCode]
		if !ok {
			return fmt.Errorf("department code not found: %s", item.DepartmentCode)
		}
		item.Employee.DepartmentID = dept.ID

		if item.Employee.BirthDate != nil {
			username := resolveSeedUsername(item)
			if username == "" {
				username = buildUserName(item.Employee.FirstName, item.Employee.LastName, *item.Employee.BirthDate)
			}

			var user model.User
			if err := ctx.Where("user_name = ?", username).First(&user).Error; err == nil {
				item.Employee.UserID = &user.ID
			}
		}

		var existing model.Employee
		if err := ctx.Where("first_name = ? AND last_name = ? AND department_id = ?", item.Employee.FirstName, item.Employee.LastName, item.Employee.DepartmentID).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := ctx.Create(&item.Employee).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			update := map[string]interface{}{}
			if existing.UserID == nil && item.Employee.UserID != nil {
				update["user_id"] = item.Employee.UserID
			}
			if len(update) > 0 {
				if err := ctx.Model(&existing).Updates(update).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func seedDepartmentManagers(ctx *gorm.DB) error {
	for _, item := range getSeedEmployeeInfo() {
		if !item.IsManager || item.Employee.BirthDate == nil {
			continue
		}

		username := resolveSeedUsername(item)
		if username == "" {
			username = buildUserName(item.Employee.FirstName, item.Employee.LastName, *item.Employee.BirthDate)
		}
		var u model.User
		if err := ctx.Where("user_name = ?", username).First(&u).Error; err != nil {
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

		if err := ctx.Model(&model.Department{}).Where("code = ?", item.DepartmentCode).Update("manager_id", manager.ID).Error; err != nil {
			return err
		}
	}

	return nil
}

// ─── Seed data definitions ────────────────────────────────────────────────────

type seedEmployeeInfo struct {
	Employee       model.Employee
	RoleName       string
	DepartmentCode string
	IsManager      bool
	DemoUsername   string
}

func resolveSeedUsername(item seedEmployeeInfo) string {
	if item.DemoUsername != "" {
		return item.DemoUsername
	}
	return ""
}

func getSeedEmployeeInfo() []seedEmployeeInfo {
	b1 := time.Date(1988, 5, 10, 0, 0, 0, 0, time.UTC)
	b2 := time.Date(1990, 8, 20, 0, 0, 0, 0, time.UTC)
	b3 := time.Date(1993, 11, 7, 0, 0, 0, 0, time.UTC)

	return []seedEmployeeInfo{
		{Employee: model.Employee{FirstName: "Chí", LastName: "Quốc AD", Phone: "0912345678", Position: "Head of Engineering", Salary: 47000000, JoinDate: time.Date(2017, 4, 1, 0, 0, 0, 0, time.UTC), BirthDate: &b1, Gender: "male", Status: "active"}, RoleName: "admin", DepartmentCode: "IT", IsManager: true, DemoUsername: "chiquoc23AD"},
		{Employee: model.Employee{FirstName: "Chí", LastName: "Quốc HR", Phone: "0912345681", Position: "HR Manager", Salary: 33000000, JoinDate: time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC), BirthDate: &b2, Gender: "male", Status: "active"}, RoleName: "hr", DepartmentCode: "HR", IsManager: true, DemoUsername: "chiquoc23HR"},
		{Employee: model.Employee{FirstName: "Chí", LastName: "Quốc EMP", Phone: "0912345680", Position: "Backend Developer", Salary: 28000000, JoinDate: time.Date(2022, 9, 5, 0, 0, 0, 0, time.UTC), BirthDate: &b3, Gender: "female", Status: "active"}, RoleName: "employee", DepartmentCode: "FIN", DemoUsername: "chiquoc23EMP"},
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func buildUserName(firstName, lastName string, birthDate time.Time) string {
	return fmt.Sprintf("%s%d", normalizeName(firstName+lastName), birthDate.Year())
}

func buildEmail(username string) string {
	return fmt.Sprintf("%s@company.vn", username)
}

func normalizeName(value string) string {
	value = strings.ToLower(value)
	value = removeVietnameseAccents(value)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, ".", "")
	return value
}

func removeVietnameseAccents(value string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ả", "a", "ã", "a", "ạ", "a",
		"â", "a", "ấ", "a", "ầ", "a", "ẩ", "a", "ẫ", "a", "ậ", "a",
		"ă", "a", "ắ", "a", "ằ", "a", "ẳ", "a", "ẵ", "a", "ặ", "a",
		"é", "e", "è", "e", "ẻ", "e", "ẽ", "e", "ẹ", "e",
		"ê", "e", "ế", "e", "ề", "e", "ể", "e", "ễ", "e", "ệ", "e",
		"í", "i", "ì", "i", "ỉ", "i", "ĩ", "i", "ị", "i",
		"ó", "o", "ò", "o", "ỏ", "o", "õ", "o", "ọ", "o",
		"ô", "o", "ố", "o", "ồ", "o", "ổ", "o", "ỗ", "o", "ộ", "o",
		"ơ", "o", "ớ", "o", "ờ", "o", "ở", "o", "ỡ", "o", "ợ", "o",
		"ú", "u", "ù", "u", "ủ", "u", "ũ", "u", "ụ", "u",
		"ư", "u", "ứ", "u", "ừ", "u", "ử", "u", "ữ", "u", "ự", "u",
		"ý", "y", "ỳ", "y", "ỷ", "y", "ỹ", "y", "ỵ", "y",
		"đ", "d",
		"Á", "A", "À", "A", "Ả", "A", "Ã", "A", "Ạ", "A",
		"Â", "A", "Ấ", "A", "Ầ", "A", "Ẩ", "A", "Ẫ", "A", "Ậ", "A",
		"Ă", "A", "Ắ", "A", "Ằ", "A", "Ẳ", "A", "Ẵ", "A", "Ặ", "A",
		"É", "E", "È", "E", "Ẻ", "E", "Ẽ", "E", "Ẹ", "E",
		"Ê", "E", "Ế", "E", "Ề", "E", "Ể", "E", "Ễ", "E", "Ệ", "E",
		"Í", "I", "Ì", "I", "Ỉ", "I", "Ĩ", "I", "Ị", "I",
		"Ó", "O", "Ò", "O", "Ỏ", "O", "Õ", "O", "Ọ", "O",
		"Ô", "O", "Ố", "O", "Ồ", "O", "Ổ", "O", "Ỗ", "O", "Ộ", "O",
		"Ơ", "O", "Ớ", "O", "Ờ", "O", "Ở", "O", "Ỡ", "O", "Ợ", "O",
		"Ú", "U", "Ù", "U", "Ủ", "U", "Ũ", "U", "Ụ", "U",
		"Ư", "U", "Ứ", "U", "Ừ", "U", "Ử", "U", "Ữ", "U", "Ự", "U",
		"Ý", "Y", "Ỳ", "Y", "Ỷ", "Y", "Ỹ", "Y", "Ỵ", "Y",
		"Đ", "D",
	)
	return replacer.Replace(value)
}
