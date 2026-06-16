package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
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
		log.Fatalf("không thể kết nối MySQL server: %v", err)
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	if err := tempDB.Exec(createDBQuery).Error; err != nil {
		log.Fatalf("Không thể tạo database: %v", err)
	}

	if sqlTempDB, err := tempDB.DB(); err == nil {
		if err := sqlTempDB.Close(); err != nil {
			log.Printf("failed to close temp db: %v", err)
		}
	}
}

func RunMigrations(db *gorm.DB) {
	if err := db.AutoMigrate(&model.Role{}, &model.Permission{}, &model.RolePermission{}, &model.UserPermission{}, &model.User{}); err != nil {
		log.Fatalf("tạo bảng thất bại: %v", err)
	}

	origDisableFK := false
	if db.Config != nil {
		origDisableFK = db.DisableForeignKeyConstraintWhenMigrating
		db.DisableForeignKeyConstraintWhenMigrating = true
	}

	if err := db.AutoMigrate(&model.Department{}, &model.Employee{}); err != nil {
		log.Fatalf("Tạo bảng thất bại: %v", err)
	}

	if db.Config != nil {
		db.DisableForeignKeyConstraintWhenMigrating = origDisableFK
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
		if err := seedPermissions(ctx); err != nil {
			return err
		}
		if err := seedRolePermissions(ctx); err != nil {
			return err
		}
		if err := seedDepartments(ctx); err != nil {
			return err
		}
		if err := seedUsers(ctx); err != nil {
			return err
		}
		if err := seedEmployees(ctx); err != nil {
			return err
		}
		if err := seedDepartmentManagers(ctx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("Chèn dữ liệu thất bại!: %v", err)
		return
	}

	log.Println("Chèn dữ liệu thành công")
	log.Println("=================================================================")
	log.Println("🔑 DANH SÁCH TÀI KHOẢN MẪU ĐÃ ĐƯỢC KHỞI TẠO:")
	log.Println("   1. [Admin]      Email: chiquoc23AD@company.vn   | Password: chiquoc23AD")
	log.Println("   2. [HR]         Email: chiquoc23HR@company.vn   | Password: chiquoc23HR")
	log.Println("   3. [Employee]   Email: chiquoc23EMP@company.vn  | Password: chiquoc23EMP")
	log.Println("=================================================================")
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

	var allPerms []model.Permission
	if err := ctx.Find(&allPerms).Error; err != nil {
		return err
	}

	adminCodes := []string{"employee.read", "employee.create", "employee.update", "employee.delete", "user.read", "user.create", "user.update", "user.delete", "department.read", "department.create", "department.update", "department.delete"}
	hrCodes := []string{"employee.read", "employee.create", "employee.update", "department.read"}
	employeeCodes := []string{"employee.read", "department.read"}

	assign := func(roleID uint, codes []string) error {
		if err := ctx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		for _, code := range codes {
			var perm model.Permission
			if err := ctx.Where("code = ?", code).First(&perm).Error; err != nil {
				return err
			}
			if err := ctx.Create(&model.RolePermission{RoleID: roleID, PermissionID: perm.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	}

	if err := assign(adminRole.ID, adminCodes); err != nil {
		return err
	}
	if err := assign(hrRole.ID, hrCodes); err != nil {
		return err
	}
	if err := assign(employeeRole.ID, employeeCodes); err != nil {
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
		{Employee: model.Employee{FirstName: "Chí", LastName: "Quốc", Phone: "0912345678", Position: "Head of Engineering", Salary: 47000000, JoinDate: time.Date(2017, 4, 1, 0, 0, 0, 0, time.UTC), BirthDate: &b1, Gender: "male", Status: "active"}, RoleName: "admin", DepartmentCode: "IT", IsManager: true, DemoUsername: "chiquoc23AD"},
		{Employee: model.Employee{FirstName: "Chí", LastName: "Quốc HR", Phone: "0912345681", Position: "HR Manager", Salary: 33000000, JoinDate: time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC), BirthDate: &b2, Gender: "male", Status: "active"}, RoleName: "hr", DepartmentCode: "HR", IsManager: true, DemoUsername: "chiquoc23HR"},
		{Employee: model.Employee{FirstName: "Lan", LastName: "Trần", Phone: "0912345680", Position: "Backend Developer", Salary: 28000000, JoinDate: time.Date(2022, 9, 5, 0, 0, 0, 0, time.UTC), BirthDate: &b3, Gender: "female", Status: "active"}, RoleName: "employee", DepartmentCode: "FIN", DemoUsername: "chiquoc23EMP"},
	}
}

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
	)
	return replacer.Replace(value)
}
