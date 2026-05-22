package model

type DepartmentEmployeeCount struct {
	DepartmentName    string `json:"department_name"`
	EmployeeCount     int64  `json:"employee_count"`
	DepartmentManager string `json:"department_manager"`
}

type Dashboards struct {
	TotalUsers           int64                     `json:"total_users"`
	TotalDepartments     int64                     `json:"total_departments"`
	TotalEmployees       int64                     `json:"total_employees"`
	TotalEmployeesActive int64                     `json:"total_employees_active"`
	DepartmentStats      []DepartmentEmployeeCount `json:"department_stats"`
}
