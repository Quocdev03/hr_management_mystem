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
		if ctx.Where("name = ?", role.Name).First(&existing).Error == nil {
			continue
		}
		if err := ctx.Create(&role).Error; err != nil {
			return err
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
		ctx.Where("code = ?", d.Code).FirstOrCreate(&d)
	}
	return nil
}

func seedUsersAndEmployees(ctx *gorm.DB) error {
	var adminRole, hrRole, empRole model.Role
	ctx.Where("name = ?", "admin").First(&adminRole)
	ctx.Where("name = ?", "hr").First(&hrRole)
	ctx.Where("name = ?", "employee").First(&empRole)

	hashedPwdAdmin, _ := bcrypt.GenerateFromPassword([]byte("12345678Quoc@"), bcrypt.DefaultCost)
	hashedPwdEmployee, _ := bcrypt.GenerateFromPassword([]byte("Password123@"), bcrypt.DefaultCost)

	users := []model.User{
		{UserName: "admin", Email: "Quocdt2003@ccquocn8n.cloud", Password: string(hashedPwdAdmin), RoleID: adminRole.ID},
		// IT
		{UserName: "it_lead", Email: "it_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "it_dev", Email: "it_dev@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// HR
		{UserName: "hr_lead", Email: "hr_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: hrRole.ID},
		{UserName: "hr_recruiter", Email: "hr_recruiter@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// FIN
		{UserName: "fin_lead", Email: "fin_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "fin_analyst", Email: "fin_analyst@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// SALES
		{UserName: "sales_lead", Email: "sales_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "sales_rep1", Email: "sales_rep1@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// MKT
		{UserName: "mkt_lead", Email: "mkt_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "mkt_exec", Email: "mkt_exec@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// CS
		{UserName: "cs_lead", Email: "cs_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "cs_agent", Email: "cs_agent@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// OPS
		{UserName: "ops_lead", Email: "ops_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "ops_exec", Email: "ops_exec@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// LEGAL
		{UserName: "legal_lead", Email: "legal_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "legal_exec", Email: "legal_exec@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// PROC
		{UserName: "proc_lead", Email: "proc_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "proc_exec", Email: "proc_exec@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		// RND
		{UserName: "rnd_lead", Email: "rnd_lead@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
		{UserName: "rnd_researcher", Email: "rnd_researcher@ccquocn8n.cloud", Password: string(hashedPwdEmployee), RoleID: empRole.ID},
	}

	for i := range users {
		var existing model.User
		if err := ctx.Where("email = ?", users[i].Email).First(&existing).Error; err != nil {
			ctx.Create(&users[i])
		} else {
			users[i] = existing
		}
	}

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

	b1 := time.Date(1988, 5, 10, 0, 0, 0, 0, time.UTC)
	b2 := time.Date(1990, 8, 20, 0, 0, 0, 0, time.UTC)
	b3 := time.Date(1993, 11, 7, 0, 0, 0, 0, time.UTC)
	b4 := time.Date(1996, 3, 3, 0, 0, 0, 0, time.UTC)

	employees := []model.Employee{
		// ===== IT (2 người) =====
		{
			DepartmentID: itDept.ID,
			FirstName:    "Chí",
			LastName:     "Cao",
			Email:        "Quocdt2003@ccquocn8n.cloud",
			Phone:        "0909123456",
			Position:     "Head of Engineering",
			Salary:       47000000,
			JoinDate:     time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b1,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: itDept.ID,
			FirstName:    "Hùng",
			LastName:     "Nguyễn",
			Email:        "it_lead@ccquocn8n.cloud",
			Phone:        "0918234567",
			Position:     "Tech Lead",
			Salary:       35000000,
			JoinDate:     time.Date(2020, 9, 25, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b2,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: itDept.ID,
			FirstName:    "Lan",
			LastName:     "Trần",
			Email:        "it_dev@ccquocn8n.cloud",
			Phone:        "0934567890",
			Position:     "Backend Developer",
			Salary:       28000000,
			JoinDate:     time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b3,
			Gender:       "female",
			Status:       "inactive",
		},

		// ===== HR (2 người) =====
		{
			DepartmentID: hrDept.ID,
			FirstName:    "Hà",
			LastName:     "Nguyễn",
			Email:        "hr_lead@ccquocn8n.cloud",
			Phone:        "0901111111",
			Position:     "HR Manager",
			Salary:       33000000,
			JoinDate:     time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b2,
			Gender:       "female",
			Status:       "active",
		},
		{
			DepartmentID: hrDept.ID,
			FirstName:    "Mai",
			LastName:     "Trần",
			Email:        "hr_recruiter@ccquocn8n.cloud",
			Phone:        "0902222222",
			Position:     "Recruiter",
			Salary:       21000000,
			JoinDate:     time.Date(2022, 7, 10, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b3,
			Gender:       "female",
			Status:       "active",
		},

		// ===== FIN (2 người) =====
		{
			DepartmentID: finDept.ID,
			FirstName:    "Quyên",
			LastName:     "Phạm",
			Email:        "fin_lead@ccquocn8n.cloud",
			Phone:        "0904444444",
			Position:     "Finance Manager",
			Salary:       36000000,
			JoinDate:     time.Date(2019, 5, 20, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b3,
			Gender:       "female",
			Status:       "active",
		},
		{
			DepartmentID: finDept.ID,
			FirstName:    "Long",
			LastName:     "Nguyễn",
			Email:        "fin_analyst@ccquocn8n.cloud",
			Phone:        "0905555555",
			Position:     "Financial Analyst",
			Salary:       25000000,
			JoinDate:     time.Date(2022, 3, 12, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b4,
			Gender:       "male",
			Status:       "active",
		},

		// ===== SALES (2 người) =====
		{
			DepartmentID: salesDept.ID,
			FirstName:    "Hoàng",
			LastName:     "Vũ",
			Email:        "sales_lead@ccquocn8n.cloud",
			Phone:        "0907777777",
			Position:     "Sales Manager",
			Salary:       38000000,
			JoinDate:     time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b4,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: salesDept.ID,
			FirstName:    "Thảo",
			LastName:     "Nguyễn",
			Email:        "sales_rep1@ccquocn8n.cloud",
			Phone:        "0908888888",
			Position:     "Sales Executive",
			Salary:       26000000,
			JoinDate:     time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b1,
			Gender:       "female",
			Status:       "inactive",
		},

		// ===== MKT (2 người) =====
		{
			DepartmentID: mktDept.ID,
			FirstName:    "Linh",
			LastName:     "Đặng",
			Email:        "mkt_lead@ccquocn8n.cloud",
			Phone:        "0911223344",
			Position:     "Marketing Manager",
			Salary:       34000000,
			JoinDate:     time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b2,
			Gender:       "female",
			Status:       "active",
		},
		{
			DepartmentID: mktDept.ID,
			FirstName:    "Tùng",
			LastName:     "Bùi",
			Email:        "mkt_exec@ccquocn8n.cloud",
			Phone:        "0922334455",
			Position:     "Marketing Executive",
			Salary:       22000000,
			JoinDate:     time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b4,
			Gender:       "male",
			Status:       "active",
		},

		// ===== CS (2 người) =====
		{
			DepartmentID: csDept.ID,
			FirstName:    "Ngọc",
			LastName:     "Lý",
			Email:        "cs_lead@ccquocn8n.cloud",
			Phone:        "0933445566",
			Position:     "CS Manager",
			Salary:       30000000,
			JoinDate:     time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b1,
			Gender:       "female",
			Status:       "active",
		},
		{
			DepartmentID: csDept.ID,
			FirstName:    "Khánh",
			LastName:     "Phan",
			Email:        "cs_agent@ccquocn8n.cloud",
			Phone:        "0944556677",
			Position:     "CS Agent",
			Salary:       18000000,
			JoinDate:     time.Date(2022, 11, 20, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b3,
			Gender:       "male",
			Status:       "inactive",
		},

		// ===== OPS (2 người) =====
		{
			DepartmentID: opsDept.ID,
			FirstName:    "Tuấn",
			LastName:     "Trịnh",
			Email:        "ops_lead@ccquocn8n.cloud",
			Phone:        "0955667788",
			Position:     "Operations Manager",
			Salary:       35000000,
			JoinDate:     time.Date(2019, 10, 5, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b2,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: opsDept.ID,
			FirstName:    "Vân",
			LastName:     "Hoàng",
			Email:        "ops_exec@ccquocn8n.cloud",
			Phone:        "0966778899",
			Position:     "Operations Executive",
			Salary:       20000000,
			JoinDate:     time.Date(2023, 3, 14, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b4,
			Gender:       "female",
			Status:       "active",
		},

		// ===== LEGAL (2 người) =====
		{
			DepartmentID: legalDept.ID,
			FirstName:    "Hải",
			LastName:     "Đinh",
			Email:        "legal_lead@ccquocn8n.cloud",
			Phone:        "0977889900",
			Position:     "Legal Manager",
			Salary:       40000000,
			JoinDate:     time.Date(2018, 8, 1, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b1,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: legalDept.ID,
			FirstName:    "Thư",
			LastName:     "Võ",
			Email:        "legal_exec@ccquocn8n.cloud",
			Phone:        "0988990011",
			Position:     "Legal Executive",
			Salary:       24000000,
			JoinDate:     time.Date(2022, 4, 18, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b3,
			Gender:       "female",
			Status:       "active",
		},

		// ===== PROC (2 người) =====
		{
			DepartmentID: procDept.ID,
			FirstName:    "Phúc",
			LastName:     "Lê",
			Email:        "proc_lead@ccquocn8n.cloud",
			Phone:        "0999001122",
			Position:     "Procurement Manager",
			Salary:       32000000,
			JoinDate:     time.Date(2020, 12, 10, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b2,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: procDept.ID,
			FirstName:    "Trang",
			LastName:     "Ngô",
			Email:        "proc_exec@ccquocn8n.cloud",
			Phone:        "0900112233",
			Position:     "Procurement Executive",
			Salary:       19000000,
			JoinDate:     time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b4,
			Gender:       "female",
			Status:       "inactive",
		},

		// ===== RND (2 người) =====
		{
			DepartmentID: rndDept.ID,
			FirstName:    "Khoa",
			LastName:     "Dương",
			Email:        "rnd_lead@ccquocn8n.cloud",
			Phone:        "0911001122",
			Position:     "R&D Manager",
			Salary:       42000000,
			JoinDate:     time.Date(2019, 6, 15, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b1,
			Gender:       "male",
			Status:       "active",
		},
		{
			DepartmentID: rndDept.ID,
			FirstName:    "Vy",
			LastName:     "Trương",
			Email:        "rnd_researcher@ccquocn8n.cloud",
			Phone:        "0922113344",
			Position:     "Researcher",
			Salary:       27000000,
			JoinDate:     time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC),
			BirthDate:    &b3,
			Gender:       "female",
			Status:       "active",
		},
	}

	for _, e := range employees {
		var u model.User
		if err := ctx.Where("email = ?", e.Email).First(&u).Error; err == nil {
			e.UserID = &u.ID
		}

		var existing model.Employee
		if err := ctx.Where("email = ?", e.Email).First(&existing).Error; err != nil {
			ctx.Create(&e)
		} else {
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

	// Gán manager cho từng phòng ban
	managersMap := map[string]string{
		"IT":    "Quocdt2003@ccquocn8n.cloud",
		"HR":    "hr_lead@ccquocn8n.cloud",
		"FIN":   "fin_lead@ccquocn8n.cloud",
		"SALES": "sales_lead@ccquocn8n.cloud",
		"MKT":   "mkt_lead@ccquocn8n.cloud",
		"CS":    "cs_lead@ccquocn8n.cloud",
		"OPS":   "ops_lead@ccquocn8n.cloud",
		"LEGAL": "legal_lead@ccquocn8n.cloud",
		"PROC":  "proc_lead@ccquocn8n.cloud",
		"RND":   "rnd_lead@ccquocn8n.cloud",
	}

	for code, email := range managersMap {
		var manager model.Employee
		if err := ctx.Where("email = ?", email).First(&manager).Error; err == nil && manager.ID != 0 {
			ctx.Model(&model.Department{}).Where("code = ?", code).Update("manager_id", manager.ID)
		}
	}

	return nil
}