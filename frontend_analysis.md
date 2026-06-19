# 🔍 Frontend Analysis — HR Management System

> **Stack**: Vue 3 (Composition API) + Vite + Pinia + Vue Router + Axios  
> **Kiến trúc**: SPA (Single Page Application), Modular Store, Composable-based logic  

---

## Kiến trúc tổng quan

Dự án tuân thủ tốt các pattern của Vue 3 hiện đại:
- **State Management**: Sử dụng `Pinia` chia thành các module (auth, user, employee, department, dashboard).
- **API Layer**: Sử dụng `Axios` với Interceptors để tự động gắn token và xử lý refresh token.
- **Reusability**: Logic dùng chung được tách thành các Composables (`usePaginatedSearch`, `usePermissions`, `useModalState`).
- **Styling**: Vanilla CSS thuần với CSS Variables cho theming.

---

## 🔴 HIGH — Vấn đề nghiêm trọng cần khắc phục

### 1. Race Condition trong tính năng Search (`usePaginatedSearch.js`)
**File**: `src/helpers/usePaginatedSearch.js`

Mặc dù đã có cơ chế `debounce` để giảm số lượng request gửi đi khi người dùng gõ phím, nhưng **không có cơ chế hủy request cũ (AbortController)**.

**Kịch bản lỗi (Race Condition)**:
1. Người dùng gõ "A" → Request 1 gửi đi (mất 500ms để phản hồi).
2. Người dùng gõ tiếp "B" thành "AB" → Request 2 gửi đi (mất 100ms để phản hồi).
3. Request 2 trả về kết quả cho "AB" và hiển thị lên UI.
4. Request 1 trả về kết quả cho "A" (do mạng chậm) và đè lên UI.
**Hệ quả**: Ô tìm kiếm hiển thị "AB" nhưng kết quả bảng lại hiển thị dữ liệu của "A" (Cache/State Staleness).

> **Gợi ý**: Tích hợp `AbortController` vào `fetchFn` hoặc Axios, lưu trữ `abortController` trong `usePaginatedSearch` và gọi `abort()` trước khi gửi request mới.

---

### 2. Thiết kế Axios Interceptor gây dồn ứ Request (Thundering Herd)
**File**: `src/api/index.js`

Cơ chế refresh token bằng `failedQueue` hiện tại xử lý khá tốt việc nhiều request cùng lúc bị 401. Tuy nhiên, nó dựa hoàn toàn vào việc **đợi server trả về 401**.

**Kịch bản**: 
1. Access token hết hạn nhưng frontend không biết.
2. User vào trang Dashboard, frontend bắn đồng thời 5 request (stats, employees, depts...).
3. Cả 5 request đều đập vào backend, backend xử lý 5 lần (tốn CPU/DB) và trả về 5 lỗi 401.
4. Frontend nhận 401 đầu tiên, khóa queue, gọi refresh token, nhận token mới, gửi lại 5 request.

> **Gợi ý**: Nên kiểm tra `exp` của JWT token ngay trên frontend trước khi request rời đi (trong `request.interceptor`). Nếu token đã hết hạn, chặn request lại ngay lập tức, đưa vào queue và tự động gọi refresh token mà không cần làm phiền backend trả 401.

---

## 🟡 MEDIUM — Cải thiện chất lượng & Bảo mật

### 3. Rủi ro bảo mật XSS khi lưu Token vào LocalStorage
**File**: `src/store/auth.js` & `src/api/index.js`

```javascript
localStorage.setItem('access_token', access_token);
localStorage.setItem('refresh_token', newRefreshToken);
localStorage.setItem('user', JSON.stringify(user));
```

Việc lưu trữ `access_token` và `refresh_token` trong `localStorage` khiến ứng dụng dễ bị tấn công XSS (Cross-Site Scripting). Bất kỳ mã JS độc hại nào (từ thư viện bên thứ 3 bị compromise, hoặc XSS injection) đều có thể đọc được token và đánh cắp phiên đăng nhập.

> **Gợi ý**: 
> - Lý tưởng nhất: Backend nên set token qua `HttpOnly Cookie`.
> - Nếu phải dùng Token trong body: Nên lưu `access_token` trong RAM (Vue state) và chỉ lưu `refresh_token` ở HttpOnly Cookie (cần backend hỗ trợ). Hiện tại kiến trúc phụ thuộc hoàn toàn vào LocalStorage.

---

### 4. Hardcode Logic Role ID (Fragile Code)
**File**: `src/components/UserModal.vue` và `src/views/UserView.vue`

```javascript
// UserModal.vue
isRoleDisabled.value = user.role_id === 1 || currentUser?.id === user.id;

// UserView.vue
{{ user.role_id === 1 ? 'Quản trị' : user.role_id === 2 ? 'Quản lý' : 'Nhân viên' }}
```

Logic phân quyền đang hardcode cứng `role_id === 1` là Admin, `2` là HR. Nếu Database thay đổi ID của Role (ví dụ khi migrate sang server khác, ID bị auto-increment lệch), toàn bộ logic Frontend sẽ sai lệch hoàn toàn.

> **Gợi ý**: Thay vì check `role_id`, nên check theo trường `role.name` hoặc code (ví dụ `user.role?.name === 'admin'`).

---

### 5. Double Fetching do thiếu Cache Cục Bộ
Khi thực hiện Create/Update/Delete (CUD) nhân viên, ứng dụng gọi:
```javascript
await Promise.all([
    loadEmployees(pagination.value.page),
    departmentStore.fetchDepartments(...)
]);
```
Mặc dù điều này đảm bảo tính nhất quán của dữ liệu (tránh Cache Staleness), nhưng nó làm tăng tải cho server không cần thiết nếu dữ liệu phòng ban không thực sự thay đổi. Các state management trong ứng dụng SPA thường chỉ cập nhật cục bộ (Optimistic UI Update) hoặc invalidate đúng resource bị ảnh hưởng.

---

## 🟢 LOW — Nợ kỹ thuật / Nice-to-have

### 6. Không xử lý lỗi Global (Global Error Boundary)
Trong `main.js`, không có cơ chế `app.config.errorHandler`. Nếu một Component xảy ra lỗi runtime (ví dụ biến `undefined` render trên template), toàn bộ giao diện có thể biến mất (White Screen of Death) mà không có thông báo cho người dùng.

### 7. Dùng `JSON.stringify` để so sánh Object
**File**: `src/helpers/buildPatchPayload.js`
```javascript
if (typeof value === "object") return JSON.stringify(value);
```
So sánh hai object bằng `JSON.stringify` không an toàn do thứ tự các key trong object có thể khác nhau (dẫn đến chuỗi JSON khác nhau dù dữ liệu giống nhau). Tốt nhất nên dùng một hàm deep equal hoặc dùng thư viện (như `lodash/isEqual`).

### 8. `index.html` và file `.env` chứa URL hardcode
**File**: `.env.example`
Mặc định Vite Proxy thường được dùng trong dev. File cấu hình hiện tại trỏ trực tiếp đến `http://localhost:8080/api/v1` thay vì đi qua Vite proxy, dẫn đến có thể gặp lỗi CORS trong lúc phát triển nếu backend chưa mở CORS đúng cách.

---

## Tổng kết ưu / nhược

### ✅ Điểm mạnh
- Kiến trúc Modular rõ ràng (chia Store, API, Views, Components, Helpers).
- Token Refresh Flow với Queueing mechanism xử lý được tình huống Concurrent 401s.
- Composables (`usePermissions`, `usePaginatedSearch`) giúp code tái sử dụng rất tốt, sạch sẽ.
- UI/UX được chăm chút kỹ lưỡng, sử dụng hiệu ứng Glassmorphism và layout responsive tốt.

### ❌ Cần tập trung sửa chữa
- Bắt buộc xử lý **Race Condition** cho tính năng Search.
- Cần dọn dẹp các **Hardcoded Role IDs** rải rác trong code.
