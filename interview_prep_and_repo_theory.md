# 🧠 CẨM NANG LÝ THUYẾT DỰ ÁN & BỘ ĐỀ HỎI ĐÁP PHỎNG VẤN (HR MANAGEMENT SYSTEM)

Tài liệu này tổng hợp toàn bộ các kiến thức lý thuyết, mô hình kiến trúc, các công nghệ và giải pháp kỹ thuật đã áp dụng trong codebase của **Hệ thống Quản lý Nhân sự (HRM)**. Mỗi chủ đề được cấu trúc chi tiết: **Định nghĩa là gì -> Dùng để làm gì -> Minh họa thực tế trong dự án -> Các câu hỏi phỏng vấn thường gặp & cách trả lời**.

---

## MỤC LỤC
1. [Kiến Trúc Hệ Thống & Thiết Kế Phần Mềm (Architecture & System Design)](#1-kiến-trúc-hệ-thống--thiết-kế-phần-mềm-architecture--system-design)
2. [Ngôn Ngữ Go & Backend Patterns](#2-ngôn-ngữ-go--backend-patterns)
3. [Cơ Sở Dữ Liệu & ORM (MySQL & GORM)](#3-cơ-sở-dữ-liệu--orm-mysql--gorm)
4. [Caching & Redis (Tối Ưu Hiệu Năng & Giới Hạn Tần Suất)](#4-caching--redis-tối-ưu-hiệu-năng--giới-hạn-tần-suất)
5. [Bảo Mật & Xác Thực (Security, Authentication & Authorization)](#5-bảo-mật--xác-thực-security-authentication--authorization)
6. [Frontend & Vue 3 SPA (Composition API & Pinia)](#6-frontend--vue-3-spa-composition-api--pinia)

---

## 1. Kiến Trúc Hệ Thống & Thiết Kế Phần Mềm (Architecture & System Design)

### 1.1 Layered Architecture (Mô hình 3 lớp)
*   **Là gì:** Kiến trúc phân tầng chia mã nguồn thành các lớp có trách nhiệm riêng biệt: **Handler (Controller) ➔ Service (Business Logic) ➔ Repository (Data Access)**.
*   **Dùng để làm gì:** 
    *   Tách biệt các mối quan tâm (Separation of Concerns).
    *   Hạn chế sự phụ thuộc chéo, giúp code dễ đọc, dễ bảo trì.
    *   Dễ dàng viết Unit Test độc lập cho từng lớp bằng cách sử dụng Interface.
*   **Áp dụng trong dự án:**
    *   `Handler` (ví dụ: [employee_handler.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/handler/employee_handler.go)): Nhận HTTP Request từ Gin, validate payload đầu vào sơ bộ, gọi Service và định dạng HTTP Response.
    *   `Service` (ví dụ: [employee_service.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/service/employee_service.go)): Xử lý các quy tắc nghiệp vụ (ví dụ: ngày vào làm không thể ở tương lai, tính tuổi từ ngày sinh, kiểm tra điều kiện gắn tài khoản user).
    *   `Repository` (ví dụ: [employee_repository.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/repository/employee_repository.go)): Chứa các câu truy vấn GORM tương tác trực tiếp với cơ sở dữ liệu MySQL.

### 1.2 Dependency Injection (DI) thủ công
*   **Là gì:** Là một design pattern thiết kế sao cho một thành phần (struct/class) nhận các phụ thuộc (dependencies) của nó từ bên ngoài truyền vào thông qua constructor (Constructor Injection), thay vì tự khởi tạo bên trong.
*   **Dùng để làm gì:** 
    *   Loại bỏ sự phụ thuộc cứng (hard-coded dependencies).
    *   Giúp việc Mocking dependencies khi viết Unit Test trở nên dễ dàng.
*   **Áp dụng trong dự án:**
    *   Trong [server/main.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/cmd/server/main.go#L44-L61):
        ```go
        userRepo := repository.NewUserRepository(db)
        userScv := service.NewUserService(userRepo, cacheSvc)
        userHandler := handler.NewUserHandler(userScv)
        ```
        `NewUserService` nhận `UserRepository` thông qua interface, giúp tách biệt hoàn toàn database khỏi tầng logic.

### 1.3 Single Page Application (SPA)
*   **Là gì:** Ứng dụng web chạy trên trình duyệt, chỉ tải một trang HTML duy nhất từ đầu. Mọi tương tác chuyển trang tiếp theo sẽ được xử lý bằng cách render động dữ liệu (qua Javascript) mà không cần tải lại toàn bộ trang từ server.
*   **Dùng để làm gì:** Tạo trải nghiệm mượt mà giống ứng dụng desktop, giảm băng thông truyền tải và tăng tốc độ chuyển trang.
*   **Áp dụng trong dự án:** Frontend Vue 3 sử dụng **Vue Router** để định tuyến phía client và giao tiếp với Backend qua RESTful APIs bằng **Axios**.

---

### ❓ Câu Hỏi Phỏng Vấn Thường Gặp
> [!NOTE]
> **Q1: Tại sao bạn chọn Layered Architecture cho Backend Go? Nó có nhược điểm gì không?**
> *   *Trả lời:* Layered Architecture rất trực quan, dễ chia việc cho team và phù hợp với các ứng dụng CRUD truyền thống. Nhược điểm của nó là có thể tạo ra code "boilerplate" khi có những API cực kỳ đơn giản (chỉ lấy dữ liệu và trả ra) vẫn phải đi qua cả 3 lớp.
> 
> **Q2: Interface có vai trò gì trong Dependency Injection của Go?**
> *   *Trả lời:* Interface hoạt động như một "hợp đồng" (contract). Lớp Service chỉ phụ thuộc vào Interface của Repository mà không phụ thuộc vào struct cụ thể. Nhờ đó, trong Unit Test, ta có thể tạo ra một `MockRepository` thực thi Interface này để kiểm tra logic của Service mà không cần kết nối tới database thật.
> 
> **Q3: Sự khác biệt giữa SPA (Single Page Application) và MPA (Multi-Page Application)?**
> *   *Trả lời:* SPA tải một trang duy nhất, chuyển trang bằng JS (Client-side routing), mang lại trải nghiệm mượt mà hơn nhưng SEO ban đầu khó hơn (cần SSR/Prerendering). MPA tải lại toàn bộ trang mới từ Server khi chuyển hướng, SEO tốt hơn nhưng tải chậm hơn và tốn tài nguyên server hơn.

---

## 2. Ngôn Ngữ Go & Backend Patterns

### 2.1 Concurrency & Graceful Shutdown (Tắt nguồn an toàn)
*   **Là gì:** Quy trình tắt máy chủ một cách êm ái khi nhận tín hiệu kết thúc từ hệ điều hành (`SIGINT`, `SIGTERM`), cho phép các request đang xử lý được hoàn thành và đóng các kết nối ngoại vi một cách an toàn.
*   **Dùng để làm gì:** 
    *   Tránh ngắt đột ngột các request của người dùng (gây lỗi 502/chợt mất dữ liệu khi deploy code).
    *   Giải phóng tài nguyên (DB connections, Redis connections) để tránh rò rỉ tài nguyên (resource leak).
*   **Áp dụng trong dự án:**
    *   Trong [server/main.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/cmd/server/main.go#L77-L112):
        *   Máy chủ HTTP chạy trong một goroutine riêng để không chặn luồng chính.
        *   Sử dụng channel `quit := make(chan os.Signal, 1)` phối hợp với `signal.Notify` để lắng nghe tín hiệu tắt từ hệ thống.
        *   Khi có tín hiệu, dùng `srv.Shutdown(ctx)` với timeout là 30 giây để hoàn tất các request đang chạy. Sau đó đóng cơ sở dữ liệu (`sqlDB.Close()`) và Redis (`rdb.Close()`).

### 2.2 Error Wrapping (Bọc lỗi)
*   **Là gì:** Kỹ thuật bọc một lỗi bên trong một lỗi khác bằng động từ định dạng `%w` thông qua `fmt.Errorf` (được giới thiệu từ Go 1.13).
*   **Dùng để làm gì:** Giúp giữ nguyên vết ngăn xếp lỗi (trace stack) từ tầng sâu nhất (như DB), đồng thời cho phép các tầng trên thêm ngữ cảnh vào lỗi và kiểm tra lỗi gốc bằng `errors.Is` hoặc `errors.As`.
*   **Áp dụng trong dự án:**
    *   Ví dụ: `fmt.Errorf("Tạo nhân viên không thành công: %w", err)`
    *   Ở tầng handler hoặc service, chúng ta có thể kiểm tra xem lỗi có phải do GORM không tìm thấy bản ghi không bằng cách sử dụng `errors.Is(err, gorm.ErrRecordNotFound)`.

### 2.3 Middleware Recovery (Khôi phục Panic)
*   **Là gì:** Sử dụng hàm `recover()` bên trong một trì hoãn thực thi (`defer`) để bắt các lỗi nghiêm trọng làm sập chương trình (`panic`).
*   **Dùng để làm gì:** Ngăn chặn việc máy chủ bị crash hoàn toàn khi gặp lỗi logic không mong muốn (ví dụ: Null Pointer Dereference). Thay vào đó, server vẫn chạy bình thường và trả về lỗi 500 cho client bị panic.
*   **Áp dụng trong dự án:**
    *   Trong [middleware.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/middleware/middleware.go#L46-L63):
        ```go
        func Recovery() gin.HandlerFunc {
            return func(ctx *gin.Context) {
                defer func() {
                    if err := recover(); err != nil {
                        utils.Error("Panic recovered: %v", err)
                        ctx.JSON(500, response)
                        ctx.Abort()
                    }
                }()
                ctx.Next()
            }
        }
        ```

---

### ❓ Câu Hỏi Phỏng Vấn Thường Gặp
> [!NOTE]
> **Q1: Cơ chế Graceful Shutdown trong Go hoạt động như thế nào?**
> *   *Trả lời:* Graceful Shutdown sử dụng channel để lắng nghe các tín hiệu hệ thống (`SIGINT`, `SIGTERM`). Khi nhận được tín hiệu, ta gọi phương thức `Shutdown()` của đối tượng `http.Server` đi kèm với một `context.WithTimeout`. Phương thức này sẽ ngưng nhận request mới, đợi các request cũ xử lý xong hoặc hết hạn timeout, sau đó đóng toàn bộ connection pool và thoát ứng dụng.
> 
> **Q2: So sánh `errors.Is` và `errors.As` trong Go?**
> *   *Trả lời:* 
>     *   `errors.Is` dùng để so sánh giá trị lỗi gốc với một sentinel error cụ thể (ví dụ: `errors.Is(err, gorm.ErrRecordNotFound)`).
>     *   `errors.As` dùng để ép kiểu (cast) lỗi sang một struct error tự định nghĩa để lấy ra thông tin chi tiết của lỗi đó (ví dụ: lấy danh sách các field bị lỗi validation).
> 
> **Q3: Panic và Recover trong Go hoạt động ra sao? Có nên lạm dụng Panic để xử lý lỗi thông thường không?**
> *   *Trả lời:* Khi panic xảy ra, luồng chạy thông thường của hàm bị dừng, các hàm `defer` được gọi ngược lên call stack cho đến khi gặp `recover()`. Không nên lạm dụng panic để xử lý lỗi thông thường. Trong Go, thiết kế chuẩn là trả về lỗi dạng `error` rõ ràng ở cuối hàm. Chỉ sử dụng panic cho các lỗi không thể phục hồi khi khởi động hệ thống (ví dụ: không cấu hình được DB, mất kết nối mạng thiết yếu).

---

## 3. Cơ Sở Dữ Liệu & ORM (MySQL & GORM)

### 3.1 Soft Delete (Xóa mềm)
*   **Là gì:** Phương pháp đánh dấu bản ghi là đã xóa (bằng cách cập nhật trường thời gian `deleted_at`) thay vì xóa vật lý dòng dữ liệu khỏi ổ đĩa bằng lệnh `DELETE`.
*   **Dùng để làm gì:** 
    *   Bảo toàn dữ liệu lịch sử để đối chiếu.
    *   Cho phép khôi phục lại dữ liệu bị xóa nhầm một cách dễ dàng.
*   **Áp dụng trong dự án:**
    *   Tất cả các model kế thừa [base_model.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/model/base_model.go#L13): `DeletedAt gorm.DeletedAt gorm:"index"`.
    *   Khi gọi `db.Delete(&employee, id)`, GORM sẽ tự động cập nhật trường `deleted_at`. Các câu lệnh truy vấn sau đó như `Find` sẽ tự động có thêm điều kiện `WHERE deleted_at IS NULL`.

### 3.2 Eager Loading (Tải trước)
*   **Là gì:** Kỹ thuật truy vấn trước các mối quan hệ (đối tượng liên kết) đi kèm bản ghi chính bằng cách sử dụng JOIN hoặc các truy vấn phụ gom nhóm tự động.
*   **Dùng để làm gì:** Giải quyết vấn đề **N+1 Queries** kinh điển trong các ORM. Thay vì thực hiện thêm N truy vấn phụ để lấy thông tin phòng ban của N nhân viên, ta chỉ cần 1 hoặc 2 truy vấn lớn.
*   **Áp dụng trong dự án:**
    *   Trong [employee_repository.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/repository/employee_repository.go#L67-L74):
        ```go
        db.Preload("Department").Preload("User").Preload("User.Role").Find(&employees)
        ```
        GORM sẽ tự động gom các ID của phòng ban và tài khoản của nhân viên để thực hiện truy vấn `IN (...)`, hạn chế tối đa số lần gửi lệnh đến database.

### 3.3 Transactions (Giao dịch ACID)
*   **Là gì:** Nhóm các thao tác ghi đọc cơ sở dữ liệu vào một khối duy nhất, bảo đảm tính **Atomicity** (tất cả cùng thành công hoặc tất cả cùng rollback để hủy bỏ).
*   **Dùng để làm gì:** Bảo vệ tính nhất quán dữ liệu khi thực hiện các cập nhật phức tạp liên quan nhiều bảng.
*   **Áp dụng trong dự án:**
    *   Khi cập nhật hoặc tạo mới Nhân viên (Employee), nếu có yêu cầu gán làm Quản lý phòng ban (Manager), ta cần cập nhật đồng thời bảng `employees` và bảng `departments`.
    *   Trong [employee_service.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/service/employee_service.go#L107-L151):
        ```go
        es.db.Transaction(func(tx *gorm.DB) error {
            txUserRepo := es.userRepo.WithTx(tx)
            txEmpRepo := es.empRepo.WithTx(tx)
            // Nếu có bất cứ lỗi nào trả về từ closure, GORM sẽ tự động gọi ROLLBACK.
            // Nếu trả về nil, GORM sẽ gọi COMMIT.
        })
        ```

### 3.4 Database Connection Pooling & Indexing
*   **Là gì:** 
    *   *Connection Pool:* Tập hợp các kết nối cơ sở dữ liệu được duy trì sẵn để dùng lại thay vì tạo mới liên tục.
    *   *Indexing (Chỉ mục):* Cấu trúc dữ liệu (thường là B-Tree) giúp tăng tốc độ tìm kiếm bản ghi trên một cột cụ thể.
*   **Dùng để làm gì:** 
    *   Tối ưu hóa thời gian kết nối database và tránh làm nghẽn server khi tải cao.
    *   Tăng tốc độ câu truy vấn tìm kiếm từ độ phức tạp $O(N)$ (quét toàn bảng) xuống $O(\log N)$.
*   **Áp dụng trong dự án:**
    *   Thiết lập connection pool trong [database.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/config/database.go#L35-L43):
        *   `SetMaxOpenConns(10)`: Giới hạn tối đa 10 kết nối đồng thời.
        *   `SetMaxIdleConns(5)`: Giữ tối thiểu 5 kết nối nhàn rỗi trong pool.
    *   Thiết lập Index trong model [user_model.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/model/user_model.go#L16-L17): Cột `user_name` và `email` được đánh dấu `index` để tìm kiếm thông tin tài khoản siêu tốc khi Login.

---

### ❓ Câu Hỏi Phỏng Vấn Thường Gặp
> [!NOTE]
> **Q1: Vấn đề N+1 queries là gì? Làm sao để phát hiện và xử lý nó trong GORM?**
> *   *Trả lời:* Vấn đề N+1 xảy ra khi ta truy vấn danh sách gồm N bản ghi chính, sau đó với mỗi bản ghi, ta lại chạy thêm 1 câu truy vấn để lấy dữ liệu liên quan (tổng cộng N+1 câu lệnh). Cách phát hiện là bật log SQL của GORM ở chế độ phát triển để xem số lượng câu truy vấn chạy. Cách khắc phục là sử dụng **Eager Loading** thông qua phương thức `.Preload()` của GORM để lấy toàn bộ dữ liệu quan hệ trong 1 hoặc vài câu lệnh tối ưu.
> 
> **Q2: Giải thích cơ chế Soft Delete. Nếu tôi muốn truy vấn cả những dữ liệu đã bị xóa mềm thì làm thế nào?**
> *   *Trả lời:* Soft Delete không xóa bản ghi khỏi ổ đĩa mà chỉ ghi thời gian xóa vào trường `deleted_at`. Trong các truy vấn thông thường, ORM sẽ tự động bỏ qua các bản ghi này. Để lấy ra cả các bản ghi đã xóa mềm, ta sử dụng phương thức `.Unscoped()` trước các điều kiện truy vấn của GORM (ví dụ: `db.Unscoped().Find(&users)`).
> 
> **Q3: Trình bày ý nghĩa của việc cấu hình Connection Pool (`MaxOpen`, `MaxIdle`, `MaxLifetime`)?**
> *   *Trả lời:*
>     *   `MaxOpenConns` giới hạn tối đa số kết nối đồng thời mở ra để bảo vệ DB không bị cạn kiệt tài nguyên (Database Starvation).
>     *   `MaxIdleConns` giữ một lượng kết nối nhàn rỗi nhất định để phục vụ các yêu cầu tiếp theo ngay lập tức mà không mất thời gian thiết lập bắt tay (handshake).
>     *   `MaxLifetime` xác định tuổi thọ tối đa của 1 kết nối để dọn dẹp các kết nối quá hạn, tránh lỗi rò rỉ bộ nhớ ở tầng mạng.

---

## 4. Caching & Redis (Tối Ưu Hiệu Năng & Giới Hạn Tần Suất)

### 4.1 Fixed Window Rate Limiting (Giới hạn tần suất Login)
*   **Là gì:** Thuật toán giới hạn số lượng request được phép từ một IP cụ thể trong một khung thời gian cố định (Fixed Window).
*   **Dùng để làm gì:** Ngăn chặn các cuộc tấn công Brute-Force mật khẩu và DDoS vào API Đăng nhập.
*   **Áp dụng trong dự án:**
    *   Trong [rate_limiter.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/middleware/rate_limiter.go):
        *   Tạo key Redis dạng `rate_limit:login:<Client_IP>`.
        *   Dùng hàm nguyên tử `INCR` của Redis để tăng số lần đăng nhập.
        *   Đặt thời gian sống (TTL) của key bằng `EXPIRE` khi lượt đếm bằng 1.
        *   Nếu bộ đếm vượt quá hạn định (ví dụ: 5 lần/phút), trả về HTTP status `429 Too Many Requests` và thông báo chặn.

### 4.2 Fail-Open Design (Thiết kế phòng ngừa sự cố)
*   **Là gì:** Triết lý thiết kế hệ thống sao cho nếu một dịch vụ phụ trợ (như Redis cache) bị sập, hệ thống chính vẫn tiếp tục chạy (fail-open) thay vì từ chối phục vụ (fail-closed).
*   **Dùng để làm gì:** Tăng tính sẵn sàng (Availability) và độ bền bỉ của hệ thống khi dịch vụ bên thứ ba hoặc cache gặp sự cố.
*   **Áp dụng trong dự án:**
    *   Trong cả `RateLimiter` và `CacheResponse` middleware:
        *   Nếu biến client `rdb == nil` hoặc có lỗi kết nối Redis xảy ra trong quá trình `INCR`/`GET`, middleware sẽ ghi log lỗi (`utils.Error`) và gọi `ctx.Next()` cho request đi qua trực tiếp database thay vì chặn người dùng bằng mã lỗi 500.

### 4.3 Redis Caching cho API GET & Cache Invalidation
*   **Là gì:** Lưu trữ kết quả JSON của các API GET (đọc dữ liệu) vào bộ nhớ RAM siêu tốc Redis để phục vụ ngay các lượt truy cập sau mà không cần truy vấn MySQL. Khi có thay đổi dữ liệu (POST/PUT/DELETE), ta thực hiện xóa cache tương ứng (Cache Invalidation).
*   **Dùng để làm gì:** Giảm tải cho database MySQL, tăng tốc phản hồi API (từ vài chục ms xuống dưới 5ms).
*   **Áp dụng trong dự án:**
    *   **Cache Response Middleware** [cache.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/middleware/cache.go):
        *   Tạo cache key theo URI: `"cache:" + RequestURI()`.
        *   *Cache Hit:* Trả dữ liệu JSON lưu trong Redis ngay lập tức và gọi `ctx.Abort()`.
        *   *Cache Miss:* Chặn kết quả đầu ra bằng `bodyLogWriter`, cho request chạy xuống DB rồi lưu kết quả vào Redis với TTL là 15 phút.
    *   **Cache Invalidation** [cache.go:L68-L105](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/middleware/cache.go#L68-L105):
        *   Khi gọi các API chỉnh sửa dữ liệu, middleware `ClearCache` sẽ tự động kích hoạt.
        *   Sử dụng lệnh `SCAN` của Redis để tìm các key theo pattern (ví dụ: `cache:/api/v1/employees*`) và xóa chúng để đảm bảo dữ liệu hiển thị không bị cũ (stale data).

---

### ❓ Câu Hỏi Phỏng Vấn Thường Gặp
> [!NOTE]
> **Q1: Tại sao bạn lại sử dụng lệnh `SCAN` thay vì `KEYS` khi thực hiện Cache Invalidation trong Redis?**
> *   *Trả lời:* Lệnh `KEYS` trong Redis hoạt động bằng cách quét tuyến tính toàn bộ các key trong bộ nhớ trong một luồng đơn (Redis là single-threaded). Nếu cơ sở dữ liệu có hàng triệu key, `KEYS` sẽ chặn hoàn toàn mọi yêu cầu khác trong nhiều giây (Blocking). Lệnh `SCAN` sử dụng con trỏ (cursor) để chia nhỏ quá trình quét thành nhiều lượt nhỏ không gây chặn hệ thống, bảo vệ hiệu năng cho môi trường sản xuất (production).
> 
> **Q2: Thiết kế "Fail-Open" trong ứng dụng của bạn hoạt động thế nào? Hãy lấy một ví dụ thực tế.**
> *   *Trả lời:* Fail-Open nghĩa là nếu Redis (dùng làm cache hoặc rate limiter) gặp sự cố ngừng hoạt động, ứng dụng sẽ ghi lại lỗi cảnh báo nhưng vẫn cho phép luồng chạy thông thường tiếp tục xử lý thông qua database MySQL. Ví dụ, trong middleware `RateLimiter` của tôi, nếu `rdb.Incr` trả về lỗi kết nối Redis, chúng tôi ghi log cảnh báo và gọi `ctx.Next()` để khách hàng vẫn có thể đăng nhập bình thường, chấp nhận tạm thời không giới hạn tần suất thay vì khóa cứng toàn hệ thống.
> 
> **Q3: Trình bày ưu nhược điểm của thuật toán Fixed Window trong Rate Limiting?**
> *   *Trả lời:* 
>     *   *Ưu điểm:* Dễ cài đặt, tốn cực ít bộ nhớ trong Redis vì chỉ cần lưu một key đếm số lượng cùng một giá trị TTL.
>     *   *Nhược điểm:* Bị hiện tượng "lưu lượng bùng nổ ở ranh giới" (Boundary Bursting). Một client có thể spam toàn bộ hạn ngạch cho phép ở cuối cửa sổ thời gian thứ nhất và đầu cửa sổ thời gian thứ hai, dẫn đến số request thực tế cao gấp đôi giới hạn trong một khoảng thời gian ngắn.

---

## 5. Bảo Mật & Xác Thực (Security, Authentication & Authorization)

### 5.1 Bcrypt Password Hashing & Salt
*   **Là gì:** Bcrypt là thuật toán băm mật khẩu một chiều có tích hợp muối ngẫu nhiên (Salt) và cấu hình chi phí băm (Work Factor/Cost).
*   **Dùng để làm gì:** Bảo vệ mật khẩu người dùng. Nếu hacker chiếm được database MySQL, họ cũng không thể giải mã ngược để lấy mật khẩu gốc, và không thể dùng bảng băm sẵn (Rainbow Tables) để dò mật khẩu nhờ vào muối ngẫu nhiên.
*   **Áp dụng trong dự án:**
    *   Khi tạo user trong file seed hoặc qua API: Mật khẩu được băm bằng `bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)`.
    *   Khi Login: So sánh bằng `bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))`.

### 5.2 JWT Blacklist Logout (Đăng xuất vô hiệu hóa Token)
*   **Là gì:** JWT là token phi trạng thái (stateless) nên không thể thu hồi trước hạn trừ khi thay đổi khóa bí mật (làm mất hiệu lực toàn bộ user). Giải pháp là lưu các token đã đăng xuất vào danh sách đen (Blacklist) tạm thời.
*   **Dùng để làm gì:** Vô hiệu hóa ngay lập tức một token khi người dùng nhấn nút Logout, ngăn chặn kẻ xấu sử dụng lại token cũ này.
*   **Áp dụng trong dự án:**
    *   Trong [auth_service.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/service/auth_service.go#L43-L49):
        *   Khi Logout, token được ghi vào Redis với key `blacklist:<token>` và giá trị `true`.
        *   Thời gian sống (TTL) được đặt chính xác bằng thời hạn còn lại của token (`TokenRemainingTime`). Sau khi hết hạn tự nhiên, Redis tự động xóa để giải phóng RAM.
        *   Mọi API cần xác thực đều được kiểm tra qua [middleware.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/middleware/middleware.go#L102-L110) xem token có nằm trong blacklist của Redis không.

### 5.3 Defense in Depth (Bảo vệ đa lớp)
*   **Là gì:** Nguyên tắc thiết kế bảo mật nhiều lớp độc lập. Nếu một lớp bảo mật bị phá vỡ hoặc bỏ qua, các lớp sau vẫn bảo vệ được hệ thống.
*   **Dùng để làm gì:** Ngăn chặn tuyệt đối việc xâm nhập dữ liệu từ các client bên ngoài đã bị sửa đổi code.
*   **Áp dụng trong dự án:**
    *   *Lớp 1 (Client-side):* Vue Router Guard chặn truy cập trang nếu không có Token hoặc sai Role.
    *   *Lớp 2 (Server-side Auth):* Middleware `AuthJWT` kiểm tra chữ ký và tính hợp lệ của Token.
    *   *Lớp 3 (Server-side Role):* Middleware `RequireRole` kiểm tra vai trò cụ thể của user đối với endpoint đó.
    *   *Lớp 4 (Database check):* Ràng buộc khóa ngoại, validate nghiệp vụ trong Service trước khi thực thi ghi DB.

### 5.4 Phòng chống lỗ hổng Signature Bypass & User Enumeration
*   **Là gì:** 
    *   *Signature Bypass:* Lỗ hổng kẻ tấn công đổi thuật toán mã hóa thành `"none"` hoặc dùng khóa công khai để ký giả mạo token đối xứng HS256.
    *   *User Enumeration:* Kẻ xấu thu thập danh sách email tồn tại trong hệ thống dựa trên thông báo lỗi khác biệt khi nhập sai email so với sai mật khẩu.
*   **Dùng để làm gì:** Tránh bị giả danh tài khoản quản trị và bảo vệ quyền riêng tư của khách hàng.
*   **Áp dụng trong dự án:**
    *   *Ngăn Signature Bypass:* Trong [token.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/utils/token.go#L30-L35), bắt buộc kiểm tra kiểu thuật toán ký trước khi xác thực:
        ```go
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("Thuật toán ký không hợp lệ!")
        }
        ```
    *   *Ngăn User Enumeration:* Trong [auth_service.go](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/backend/internal/service/auth_service.go#L63-L76), dù sai email hay sai mật khẩu, hệ thống đều trả về một thông báo lỗi duy nhất: `"Email hoặc mật khẩu không hợp lệ!"`.

---

### ❓ Câu Hỏi Phỏng Vấn Thường Gặp
> [!NOTE]
> **Q1: Vì sao JWT có tính chất Stateless? Làm thế nào để thực hiện chức năng đăng xuất đối với JWT?**
> *   *Trả lời:* JWT mang tính chất Stateless vì mọi thông tin nhận diện người dùng (claims) đều được mã hóa và lưu trực tiếp trong token ở client. Server không cần lưu trạng thái session. Để đăng xuất, vì server không kiểm soát trạng thái token nên ta phải thiết lập một **Blacklist** trên Redis. Khi user logout, token được gửi lên và lưu vào Redis Blacklist với TTL bằng thời gian sống còn lại của token. Ở mỗi request kế tiếp, middleware xác thực sẽ tra cứu Redis xem token có nằm trong blacklist không để từ chối truy cập.
> 
> **Q2: Lỗ hổng "alg: none" trong JWT là gì và cách phòng chống?**
> *   *Trả lời:* Đây là lỗ hổng xảy ra khi thư viện JWT ở phía server tin tưởng vào tham số thuật toán `"alg"` ghi ở phần header của token. Nếu kẻ tấn công thay đổi header thành `"alg": "none"` và xóa bỏ phần chữ ký (signature), server cấu hình lỗi sẽ chấp nhận token đó mà không cần xác minh chữ ký. Cách phòng chống là cấu hình thư viện xác thực bắt buộc phải kiểm tra và chỉ chấp nhận một hoặc một số thuật toán ký cụ thể (như bắt buộc phải là `SigningMethodHMAC` đối với HS256).
> 
> **Q3: Tại sao lại chọn Bcrypt thay vì MD5 hay SHA256 để lưu trữ mật khẩu?**
> *   *Trả lời:* MD5 và SHA256 là các thuật toán băm tốc độ rất nhanh (fast hashes), hacker có thể dùng GPU mạnh chạy hàng tỷ phép thử mỗi giây để bẻ khóa (brute force). Bcrypt được thiết kế có chi phí tính toán cao (slow hash) và có tích hợp muối ngẫu nhiên (salt). Việc điều chỉnh cost factor làm kéo dài thời gian tính toán băm (ví dụ mất 100ms cho 1 lượt băm), khiến việc brute-force mật khẩu quy mô lớn trở nên bất khả thi về mặt chi phí và thời gian.

---

## 6. Frontend & Vue 3 SPA (Composition API & Pinia)

### 6.1 Composition API & Setup Stores
*   **Là gì:**
    *   *Composition API:* Cách viết Vue component gom nhóm logic theo tính năng bằng cách sử dụng các hàm phản trị như `ref`, `computed`, `watch` thay vì phân mảnh vào các option như `data`, `methods`, `mounted` (Options API).
    *   *Setup Store:* Cách định nghĩa store của Pinia tương tự như Composition API, sử dụng một hàm setup trả ra state và action.
*   **Dùng để làm gì:** Giúp dễ tái sử dụng code (qua Composables), tổ chức logic sạch sẽ hơn cho các component lớn và tối ưu hóa hệ thống kiểm tra kiểu dữ liệu (TypeScript).
*   **Áp dụng trong dự án:**
    *   Tất cả các store trong `frontend/src/store/` (ví dụ: [auth.js](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/frontend/src/store/auth.js)) đều viết theo dạng Setup Store:
        *   `ref()` đại diện cho State.
        *   `computed()` đại diện cho Getters.
        *   Các hàm thông thường hoặc `async function` đại diện cho Actions.

### 6.2 Component Lazy Loading (Tải chậm thành phần)
*   **Là gì:** Kỹ thuật tách nhỏ file bundle JavaScript của ứng dụng ra thành nhiều file nhỏ tương ứng với từng màn hình (Code Splitting). Trình duyệt chỉ tải file JS của trang khi người dùng truy cập trang đó.
*   **Dùng để làm gì:** Tăng tốc độ hiển thị trang đầu tiên (First Contentful Paint) vì không cần tải toàn bộ JS của cả website ngay từ lúc vào trang chủ.
*   **Áp dụng trong dự án:**
    *   Trong cấu hình router [router/index.js](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/frontend/src/router/index.js#L3-L9):
        ```javascript
        const EmployeeView = () => import("@/views/EmployeeView.vue");
        const DepartmentView = () => import("@/views/DepartmentView.vue");
        ```
        Sử dụng cú pháp hàm `() => import(...)` giúp Vite tự động cắt các View này thành các file chunk độc lập khi build production.

### 6.3 Axios Interceptors (Bộ chặn HTTP)
*   **Là gì:** Các hàm trung gian được tự động thực thi trước khi request được gửi đi (Request Interceptor) hoặc sau khi nhận phản hồi từ server (Response Interceptor).
*   **Dùng để làm gì:** 
    *   Tự động chèn token xác thực vào header của tất cả các API mà không cần viết lặp đi lặp lại ở từng trang.
    *   Xử lý lỗi tập trung (ví dụ: nếu server phản hồi mã 401 Unauthorized, lập tức xóa token cục bộ và chuyển hướng người dùng về trang đăng nhập `/login`).
*   **Áp dụng trong dự án:**
    *   Thiết lập trong [api/index.js](file:///c:/Users/ccquo/Downloads/Compressed/HR_Managerment_System/frontend/src/api/index.js):
        *   Request interceptor lấy `access_token` từ `localStorage` gắn vào header `Authorization: Bearer <token>`.
        *   Response interceptor kiểm tra lỗi `status === 401` để xóa dữ liệu lưu trữ cục bộ và buộc redirect về `/login`.

---

### ❓ Câu Hỏi Phỏng Vấn Thường Gặp
> [!NOTE]
> **Q1: Hãy so sánh Composition API với Options API trong Vue 3?**
> *   *Trả lời:* Options API bắt buộc viết code theo các thuộc tính cố định (`data`, `methods`, `computed`). Khi dự án phình to, các logic liên quan đến một tính năng bị chia nhỏ ở nhiều nơi, gây khó khăn khi đọc và bảo trì. Composition API cho phép gom các biến phản xạ, computed và hàm xử lý của cùng một tính năng lại một chỗ, thậm chí tách ra thành các file Composable độc lập để tái sử dụng ở nhiều component khác nhau.
> 
> **Q2: Pinia có điểm gì cải tiến so với Vuex trong Vue?**
> *   *Trả lời:* Pinia lược bỏ hoàn toàn khái niệm `mutations` rườm rà, cho phép thay đổi state trực tiếp hoặc thông qua các `actions` (giao diện đơn giản hơn nhiều). Pinia hỗ trợ hoàn hảo hệ thống kiểm tra kiểu dữ liệu TypeScript mà không cần cấu hình phức tạp, đồng thời hỗ trợ cú pháp Setup Store hiện đại giống hệt Composition API.
> 
> **Q3: Trình bày cơ chế hoạt động và tầm quan trọng của Axios Interceptors?**
> *   *Trả lời:* Axios Interceptors đóng vai trò như các middleware trung gian. Request interceptor cho phép chỉnh sửa cấu hình trước khi gửi đi (như đính kèm JWT header, thiết lập thời gian timeout cụ thể). Response interceptor chặn kết quả trả về từ server, giúp chuẩn hóa định dạng dữ liệu (bóc tách lấy trường `.data` trực tiếp) và xử lý các lỗi HTTP hệ thống (như tự động logout khi hết hạn token 401, hiển thị thông báo toast lỗi khi nhận 500) một cách tập trung, giúp code ở các view sạch sẽ và không bị trùng lặp logic xử lý lỗi.
