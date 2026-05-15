package model

type Dashboards struct {
	TotalUsers       int64 `json:"total_users"`
	TotalDepartments int64 `json:"total_departments"`
	TotalEmployees   int64 `json:"total_employees"`
}
