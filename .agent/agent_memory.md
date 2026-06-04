# 💾 Agent Memory

## Stage 1: Theory Scan and Interview Prep compilation

- **Task**: Quét toàn repo và lấy ra lý thuyết những cái đã áp dụng, dùng để làm gì và là gì, những cái bị hỏi khi phỏng vấn.
- **Files Touched**:
  - [interview_prep_and_repo_theory.md](file:///C:/Users/ccquo/.gemini/antigravity-ide/brain/8c5df52f-f927-4937-807d-062e0af10ad9/interview_prep_and_repo_theory.md) (Created)
- **Result**: Successfully scanned the repository, analyzed the backend layer, database connection pool, caching mechanisms, routing guards, and client-side store logic. Created a detailed report containing core concept definitions, their application inside this project, and relevant QA style interview prep guide.
- **Notes**: All code patterns follow clean layered architecture design principles, stateless JWT auth with Redis blacklisting, and Vue 3 Composition API/Pinia setup stores.

## Stage 2: Agent Skills Auditing & Redundancy Check (Executed)

- **Task**: Kiểm tra và dọn dẹp các skill trong thư mục `.agent/skills/` xem cái nào không cần thiết cho dự án.
- **Files Touched**:
  - [.agent/skills/error-handling/](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/.agent/skills/error-handling/) (Deleted)
  - [.agent/skills/database-optimizer/](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/.agent/skills/database-optimizer/) (Deleted)
  - [.agent/skills/golang-pro/](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/.agent/skills/golang-pro/) (Deleted)
  - [.agent/skills/mysql-patterns/SKILL.md](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/.agent/skills/mysql-patterns/SKILL.md) (Cleaned up Node/Python parts, added GORM config)
  - [.agent/skills/golang-style/SKILL.md](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/.agent/skills/golang-style/SKILL.md) (Cleaned up Einride organization parts)
  - [.agent/agent_memory.md](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/.agent/agent_memory.md) (Updated)
- **Result**: Successfully deleted 3 redundant skill directories and optimized 2 remaining skill configurations to align perfectly with Go + GORM + MySQL + Vue 3 architecture.
- **Notes**: Total skill count reduced from 13 to 10. Agent context space is now more clean and memory efficient.

## Stage 3: Dockerization & Seed Account Update (Executed)

- **Task**: Cấu hình Dockerfile cho backend và frontend, tích hợp docker-compose, thay đổi thông tin 3 tài khoản seed theo tên người dùng Chí Quốc và cập nhật README.
- **Files Touched**:
  - [Dockerfile](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/Dockerfile)
  - [docker-entrypoint.sh](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/docker-entrypoint.sh)
  - [Dockerfile](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/frontend/Dockerfile)
  - [nginx.conf](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/frontend/nginx.conf)
  - [docker-compose.yml](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/docker-compose.yml)
  - [README.MD](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/README.MD)
  - [setup.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/config/setup.go)
- **Result**: Hoàn thành đóng gói toàn bộ hệ thống lên Docker, đồng nhất phiên bản Go (1.26.2+) trong mọi tài liệu, cập nhật 3 tài khoản seed (`chiquoc23AD`, `chiquoc23HR`, `chiquoc23EMP`) hiển thị rõ thông tin đăng nhập khi kết thúc seed data.

## Stage 4: Refactoring Business Logic (Handler -> Service) (Executed)

- **Task**: Di chuyển logic kiểm tra phân quyền bảo mật (không tự đổi quyền bản thân/Admin khác, không tự vô hiệu hoá tài khoản chính mình) từ `user_handler.go` sang `user_service.go` theo đúng kiến trúc layered và tối ưu hóa truy vấn DB.
- **Files Touched**:
  - [user_service.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/service/user_service.go)
  - [user_handler.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/handler/user_handler.go)
- **Result**: Tách biệt logic phân quyền an toàn hoàn toàn ra khỏi tầng Handler. Khai báo các lỗi sentinel (`ErrSelfRoleChange`, `ErrSelfDeactivate`) ở tầng Service. Handler bắt lỗi sentinel và trả về status `403 Forbidden` tương ứng, loại bỏ truy vấn DB dư thừa (`GetUserByID`) giúp tăng hiệu năng hệ thống.
- **Notes**: Chạy kiểm tra biên dịch (`go build`) và phân tích tĩnh (`go vet`) thành công 100% không cảnh báo/lỗi.
