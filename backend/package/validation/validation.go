package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ValidationError chứa thông tin lỗi cho từng field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors là danh sách các lỗi validation
type ValidationErrors struct {
	Errors []FieldError `json:"errors"`
}

func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, FieldError{
		Field:   field,
		Message: message,
	})
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return ""
	}
	var msgs []string
	for _, e := range ve.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(msgs, "; ")
}

// --- Các hàm validate dùng chung ---

// IsNotEmpty kiểm tra chuỗi không rỗng sau khi trim
func IsNotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// IsValidEmail kiểm tra định dạng email hợp lệ
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidPhone kiểm tra số điện thoại Việt Nam (bắt đầu bằng 0, đúng 10 số)
func IsValidPhone(phone string) bool {
	if phone == "" {
		return true // phone có thể không bắt buộc
	}
	re := regexp.MustCompile(`^0\d{9}$`)
	return re.MatchString(phone)
}

// IsValidDate kiểm tra ngày hợp lệ theo format YYYY-MM-DD
func IsValidDate(date string) bool {
	if date == "" {
		return true // ngày có thể không bắt buộc
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// IsDateNotInFuture kiểm tra ngày không ở tương lai
func IsDateNotInFuture(date string) bool {
	if date == "" {
		return true
	}
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	return !parsed.After(time.Now())
}

// IsValidLength kiểm tra độ dài chuỗi
func IsValidLength(s string, min, max int) bool {
	length := utf8.RuneCountInString(strings.TrimSpace(s))
	return length >= min && length <= max
}

// IsNonNegative kiểm tra số không âm
func IsNonNegative(val float64) bool {
	return val >= 0
}

// IsPositive kiểm tra số dương
func IsPositive(val float64) bool {
	return val > 0
}

// IsValidStatus kiểm tra trạng thái hợp lệ
func IsValidStatus(status string, allowedStatuses []string) bool {
	if status == "" {
		return true // không bắt buộc khi update
	}
	for _, s := range allowedStatuses {
		if status == s {
			return true
		}
	}
	return false
}

// IsValidName kiểm tra tên chỉ chứa chữ cái, dấu cách và dấu tiếng Việt
func IsValidName(name string) bool {
	if name == "" {
		return true
	}
	// Cho phép Unicode letters và dấu cách
	re := regexp.MustCompile(`^[\p{L}\s]+$`)
	return re.MatchString(strings.TrimSpace(name))
}

// IsValidCode kiểm tra mã chỉ chứa chữ cái, số, gạch ngang và gạch dưới (không có khoảng trắng)
func IsValidCode(code string) bool {
	if code == "" {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(code)
}

// IsValidPassword kiểm tra mật khẩu đủ mạnh
// Ít nhất 8 ký tự, có chữ hoa, chữ thường, số
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}

// IsValidUsername kiểm tra username hợp lệ (chữ cái, số, gạch dưới, 4-50 ký tự)
func IsValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{4,50}$`)
	return re.MatchString(username)
}

// --- Các hàm validate cho từng Request ---

// ValidateCreateEmployee validate CreateEmployeeRequest
func ValidateCreateEmployee(departmentID uint, firstName, lastName, email, phone, position, joinDate string, salary float64) *ValidationErrors {
	ve := &ValidationErrors{}

	// department_id: bắt buộc, phải > 0
	if departmentID == 0 {
		ve.Add("department_id", "Phòng ban là bắt buộc")
	}

	// first_name: bắt buộc, 2-100 ký tự, chỉ chữ cái
	if !IsNotEmpty(firstName) {
		ve.Add("first_name", "Họ là bắt buộc")
	} else {
		if !IsValidLength(firstName, 2, 100) {
			ve.Add("first_name", "Họ phải từ 2 đến 100 ký tự")
		}
		if !IsValidName(firstName) {
			ve.Add("first_name", "Họ chỉ được chứa chữ cái và dấu cách")
		}
	}

	// last_name: bắt buộc, 2-100 ký tự, chỉ chữ cái
	if !IsNotEmpty(lastName) {
		ve.Add("last_name", "Tên là bắt buộc")
	} else {
		if !IsValidLength(lastName, 2, 100) {
			ve.Add("last_name", "Tên phải từ 2 đến 100 ký tự")
		}
		if !IsValidName(lastName) {
			ve.Add("last_name", "Tên chỉ được chứa chữ cái và dấu cách")
		}
	}

	// email: bắt buộc, phải đúng format
	if !IsNotEmpty(email) {
		ve.Add("email", "Email là bắt buộc")
	} else if !IsValidEmail(email) {
		ve.Add("email", "Email không đúng định dạng")
	}

	// phone: bắt buộc, phải đúng format VN
	if !IsNotEmpty(phone) {
		ve.Add("phone", "Số điện thoại là bắt buộc")
	} else if !IsValidPhone(phone) {
		ve.Add("phone", "Số điện thoại phải bắt đầu bằng 0 và có đúng 10 số")
	}

	// position: tùy chọn, nhưng nếu có thì phải 2-100 ký tự
	if position != "" && !IsValidLength(position, 2, 100) {
		ve.Add("position", "Vị trí phải từ 2 đến 100 ký tự")
	}

	// salary: không được âm
	if !IsNonNegative(salary) {
		ve.Add("salary", "Mức lương không được nhỏ hơn 0")
	}

	// join_date: nếu có, phải đúng format YYYY-MM-DD, không ở tương lai
	if joinDate != "" {
		if !IsValidDate(joinDate) {
			ve.Add("join_date", "Ngày vào làm phải đúng định dạng YYYY-MM-DD")
		} else if !IsDateNotInFuture(joinDate) {
			ve.Add("join_date", "Ngày vào làm không được là ngày trong tương lai")
		}
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidateUpdateEmployee validate UpdateEmployeeRequest
func ValidateUpdateEmployee(departmentID uint, firstName, lastName, phone, position, status string, salary float64) *ValidationErrors {
	ve := &ValidationErrors{}

	// first_name: nếu có thì phải 2-100 ký tự
	if firstName != "" {
		if !IsValidLength(firstName, 2, 100) {
			ve.Add("first_name", "Họ phải từ 2 đến 100 ký tự")
		}
		if !IsValidName(firstName) {
			ve.Add("first_name", "Họ chỉ được chứa chữ cái và dấu cách")
		}
	}

	// last_name: nếu có thì phải 2-100 ký tự
	if lastName != "" {
		if !IsValidLength(lastName, 2, 100) {
			ve.Add("last_name", "Tên phải từ 2 đến 100 ký tự")
		}
		if !IsValidName(lastName) {
			ve.Add("last_name", "Tên chỉ được chứa chữ cái và dấu cách")
		}
	}

	// phone: nếu có, phải đúng format
	if phone != "" && !IsValidPhone(phone) {
		ve.Add("phone", "Số điện thoại phải bắt đầu bằng 0 và có đúng 10 số")
	}

	// position: nếu có, phải 2-100 ký tự
	if position != "" && !IsValidLength(position, 2, 100) {
		ve.Add("position", "Vị trí phải từ 2 đến 100 ký tự")
	}

	// salary: nếu != 0, không được âm
	if salary != 0 && !IsNonNegative(salary) {
		ve.Add("salary", "Mức lương không được nhỏ hơn 0")
	}

	// status: nếu có, phải là active hoặc inactive
	if status != "" && !IsValidStatus(status, []string{"active", "inactive"}) {
		ve.Add("status", "Trạng thái chỉ có thể là 'active' hoặc 'inactive'")
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidateCreateDepartment validate CreateDepartmentRequest
func ValidateCreateDepartment(name, code, description string) *ValidationErrors {
	ve := &ValidationErrors{}

	// name: bắt buộc, 1-100 ký tự
	if !IsNotEmpty(name) {
		ve.Add("name", "Tên phòng ban là bắt buộc")
	} else if !IsValidLength(name, 1, 100) {
		ve.Add("name", "Tên phòng ban phải từ 1 đến 100 ký tự")
	}

	// code: bắt buộc, 1-20 ký tự, chỉ chữ/số/gạch ngang/gạch dưới
	if !IsNotEmpty(code) {
		ve.Add("code", "Mã phòng ban là bắt buộc")
	} else {
		if !IsValidLength(code, 1, 20) {
			ve.Add("code", "Mã phòng ban phải từ 1 đến 20 ký tự")
		}
		if !IsValidCode(code) {
			ve.Add("code", "Mã phòng ban chỉ được chứa chữ cái, số, gạch ngang và gạch dưới")
		}
	}

	// description: tùy chọn, nhưng nếu có thì tối đa 500 ký tự
	if description != "" && !IsValidLength(description, 0, 500) {
		ve.Add("description", "Mô tả phòng ban không được vượt quá 500 ký tự")
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidateUpdateDepartment validate UpdateDepartmentRequest
func ValidateUpdateDepartment(name, description string, managerID *uint) *ValidationErrors {
	ve := &ValidationErrors{}

	// name: nếu có thì 1-100 ký tự
	if name != "" && !IsValidLength(name, 1, 100) {
		ve.Add("name", "Tên phòng ban phải từ 1 đến 100 ký tự")
	}

	// description: tùy chọn, tối đa 500 ký tự
	if description != "" && !IsValidLength(description, 0, 500) {
		ve.Add("description", "Mô tả phòng ban không được vượt quá 500 ký tự")
	}

	// manager_id: nếu có, phải > 0
	if managerID != nil && *managerID == 0 {
		ve.Add("manager_id", "ID quản lý phải lớn hơn 0")
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidateLoginRequest validate LoginRequest
func ValidateLogin(email, password string) *ValidationErrors {
	ve := &ValidationErrors{}

	if !IsNotEmpty(email) {
		ve.Add("email", "Email là bắt buộc")
	} else if !IsValidEmail(email) {
		ve.Add("email", "Email không đúng định dạng")
	}

	if !IsNotEmpty(password) {
		ve.Add("password", "Mật khẩu là bắt buộc")
	} else if len(password) < 8 {
		ve.Add("password", "Mật khẩu phải có ít nhất 8 ký tự")
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidateRegister validate RegisterRequest
func ValidateRegister(username, email, password string) *ValidationErrors {
	ve := &ValidationErrors{}

	if !IsNotEmpty(username) {
		ve.Add("user_name", "Tên đăng nhập là bắt buộc")
	} else {
		if !IsValidLength(username, 4, 50) {
			ve.Add("user_name", "Tên đăng nhập phải từ 4 đến 50 ký tự")
		}
		if !IsValidUsername(username) {
			ve.Add("user_name", "Tên đăng nhập chỉ được chứa chữ cái, số và dấu gạch dưới")
		}
	}

	if !IsNotEmpty(email) {
		ve.Add("email", "Email là bắt buộc")
	} else if !IsValidEmail(email) {
		ve.Add("email", "Email không đúng định dạng")
	}

	if !IsNotEmpty(password) {
		ve.Add("password", "Mật khẩu là bắt buộc")
	} else {
		if len(password) < 8 {
			ve.Add("password", "Mật khẩu phải có ít nhất 8 ký tự")
		}
		if !IsValidPassword(password) {
			ve.Add("password", "Mật khẩu phải chứa ít nhất 1 chữ hoa, 1 chữ thường và 1 số")
		}
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

// ValidatePagination validate PaginationQuery
func ValidatePagination(page, limit int) *ValidationErrors {
	ve := &ValidationErrors{}

	if page < 1 {
		ve.Add("page", "Số trang phải lớn hơn hoặc bằng 1")
	}

	if limit < 1 || limit > 100 {
		ve.Add("limit", "Số lượng mỗi trang phải từ 1 đến 100")
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}
