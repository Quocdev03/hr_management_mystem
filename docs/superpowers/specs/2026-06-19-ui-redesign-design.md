# Tài liệu Thiết kế (Spec) - Tái cấu trúc Giao diện Đồng bộ Fresh Mint Breeze

Tài liệu này đặc tả chi tiết kế hoạch thay đổi **đồng bộ toàn bộ các file Vue (.vue) và CSS** trong hệ thống HRM sang phong cách **Bento Grid & Glassmorphism** dịu mắt, đảm bảo tính nhất quán về thẩm mỹ, cỡ chữ, khoảng cách và màu sắc.

---

## 1. Mục tiêu Thiết kế (Design Goals)
* **Nhất quán 100%:** Toàn bộ ứng dụng (tất cả các View và Component) sẽ được chuyển đổi đồng đều, không bỏ sót bất kỳ thành phần nào.
* **Thẩm mỹ Bento Grid & Glassmorphism:** Sắp xếp thông tin dạng các khối grid bất đối xứng bo tròn mượt mà trên nền kính mờ dịu mắt.
* **Typography chuẩn tiếng Việt:** Sử dụng bộ font đa dạng (`Outfit`, `Plus Jakarta Sans`, `Inter`) có hỗ trợ tiếng Việt đầy đủ và hệ thống cỡ chữ có phân cấp rõ rệt.
* **Trải nghiệm êm dịu (Low Strain):** Màu nền và điểm nhấn được thiết kế mềm mại (Forest & Mint), không bị chói mắt.

---

## 2. Hệ thống Thiết kế Toàn cục (Global CSS & Design Tokens)

### A. Tích hợp Font chữ & Màu sắc trong [main.css](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/assets/css/main.css)
* **Google Fonts Import:**
  ```css
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Outfit:wght@600;700&family=Plus+Jakarta+Sans:wght@500;600;700&display=swap');
  ```
* **CSS Variables (Tokens):**
  * Nền: `--bg-gradient: linear-gradient(135deg, #f4faf7 0%, #eaf4f0 100%)` (Tránh màu trắng chói).
  * Kính mờ: `--glass-bg: rgba(255, 255, 255, 0.75)`, `--glass-border: 1px solid rgba(255, 255, 255, 0.9)`, `--glass-shadow: 0 8px 32px 0 rgba(148, 174, 166, 0.12)`.
  * Chữ: `--text-primary: #1e293b` (Slate 800), `--text-secondary: #64748b` (Slate 500).
  * Điểm nhấn: `--primary-color: #0d9488` (Teal 600), `--primary-hover: #0f766e` (Teal 700), `--accent-mint: #10b981` (Emerald 500).
  * Bo góc: `--radius-lg: 16px`, `--radius-md: 12px`, `--radius-sm: 8px`.

---

## 3. Chi tiết Sửa đổi Đồng bộ tất cả các File Vue

### A. Layout Components

#### 1. [MainLayout.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/layout/MainLayout.vue)
* Thiết lập nền ứng dụng là `--bg-gradient`.
* Chuyển vùng chứa nội dung chính thành một ô Bento lớn bo góc `16px`.

#### 2. [Sidebar.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/Sidebar.vue)
* Chuyển nền Sidebar sang dạng kính mờ `rgba(255, 255, 255, 0.8)` có `backdrop-filter: blur(12px)`.
* Các liên kết điều hướng (`.nav-item`) đổi sang font `Inter`. Khi `active`, áp dụng gradient mềm mại từ xanh Teal sang Mint.
* Phần thông tin cá nhân (Profile) ở đáy Sidebar sử dụng vòng tròn Avatar gradient Teal/Mint.

#### 3. [Header.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/Header.vue)
* Chuyển đổi nền thanh tiêu đề rút gọn trên Mobile sang kính mờ đồng bộ.

---

### B. Module Views (Các Trang Chính)

#### 4. [DashboardView.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/views/DashboardView.vue)
* Thiết kế lại các thẻ thống kê tổng quan (Stats Card) thành các khối Bento Grid có kích thước khác nhau.
* Các con số số liệu thống kê lớn sử dụng font `Plus Jakarta Sans` với trọng lượng `700` và màu Teal sang trọng.
* Biểu đồ hoặc thanh tiến trình sử dụng dải màu lục bảo nhạt dịu mắt.

#### 5. [EmployeeView.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/views/EmployeeView.vue)
* Đưa danh sách nhân viên vào một Bento Card kính mờ lớn.
* Bảng dữ liệu `.data-table` sử dụng font `Inter`, giảm padding để tạo cấu trúc chặt chẽ, tiêu đề cột màu Slate nhạt.
* Thống nhất thiết kế cho nút thêm mới nhân viên sử dụng nút gradient Bạc hà dịu mắt.

#### 6. [DepartmentView.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/views/DepartmentView.vue)
* Chuyển đổi danh sách phòng ban sang dạng Bento Grid, mỗi phòng ban là một Glassmorphism Card riêng biệt.
* Hiển thị avatar của Trưởng phòng theo vòng tròn gradient tương ứng.

#### 7. [UserView.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/views/UserView.vue)
* Tái cấu trúc bảng quản lý tài khoản với giao diện kính mờ.
* Huy hiệu trạng thái hoạt động (Active/Inactive) chuyển thành màu xanh Teal nhạt / xám Slate nhạt dịu mát.

#### 8. [ProfileView.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/views/ProfileView.vue)
* Bố cục 2 cột Bento Grid:
  * Cột trái: Card kính mờ đứng chứa thông tin Avatar lớn chữ cái đầu, họ tên, vai trò và trạng thái hoạt động.
  * Cột phải: 2 Card kính mờ ngang chứa chi tiết Thông tin cá nhân và Thông tin công việc riêng biệt.
* Sử dụng font `Outfit` cho các tiêu đề phần và `Plus Jakarta Sans` cho phần mức lương.

#### 9. [LoginView.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/views/LoginView.vue)
* Chuyển đổi khung đăng nhập (Login Card) thành một khối kính mờ trung tâm dày dặn (`rgba(255, 255, 255, 0.85)`).
* Bổ sung hiệu ứng đốm sáng mờ nhẹ màu lục ấm ở nền sau để tạo chiều sâu thẩm mỹ.
* Các ô nhập liệu có viền bo tròn nhẹ, khi focus sẽ đổi màu viền sang xanh Teal dịu kèm hiệu ứng phát sáng mờ.

---

### C. Modals & Dialog Components (Cửa sổ nổi)

#### 10. [EmployeeModal.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/EmployeeModal.vue) & [UserModal.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/UserModal.vue) & [DepartmentModal.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/DepartmentModal.vue)
* Toàn bộ Form Modal chuyển sang sử dụng thiết kế kính mờ đè bóng sâu (`box-shadow: 0 20px 50px rgba(0, 0, 0, 0.08)`).
* Tiêu đề modal sử dụng font `Outfit`.
* Form input sử dụng font `Inter`, chiều cao chuẩn hóa `42px` dễ thao tác chạm vuốt.

#### 11. [EmployeeDetailDrawer.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/EmployeeDetailDrawer.vue)
* Khung trượt (Slide Drawer) chuyển thành tấm kính mờ trượt từ bên phải.
* Phân vùng thông tin cá nhân, công việc và tài khoản bằng các đường chia mỏng màu trắng sữa.

#### 12. [ConfirmationDialog.vue](file:///c:/Users/ccquo/Downloads/Compressed/hr_managerment_system/frontend/src/components/ConfirmationDialog.vue)
* Hộp thoại xác nhận (Xóa nhân viên/User) đổi sang thiết kế kính mờ thu gọn, nút Xác nhận màu Rose mềm (không đỏ chói).

---

## 5. Kế hoạch Kiểm tra (Verification Plan)
* Chạy biên dịch production: `npm run build` để kiểm tra độ tương thích.
* Kiểm tra tính phản hồi (Responsive) trên tất cả các file Vue đã sửa đổi.
* Đảm bảo tính nhất quán về font chữ (`Outfit`, `Plus Jakarta Sans`, `Inter`) và màu sắc (`Teal/Mint`) trên mọi màn hình.
