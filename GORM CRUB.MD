# GORM CRUD Cheat Sheet — Từ Cơ Bản Đến Thực Chiến

Trong GORM, thực chất là những lớp bao bọc (wrappers) để tự động sinh ra câu lệnh SQL tương ứng.

### Ví dụ Model cơ bản:

```go
type User struct {
    ID    uint
    Name  string
    Email string
    Age   int
}
```

---

## 1. Create — INSERT

**GORM:**

```go
user := User{
    Name:  "Khanh",
    Email: "khanh@gmail.com",
    Age:   22,
}
db.Create(&user)
```

**SQL tương ứng:**

```sql
INSERT INTO users (name, email, age) VALUES ('Khanh', 'khanh@gmail.com', 22);
```

**Kết quả:**

- Thêm dữ liệu vào Database.
- `user.ID` sẽ tự động được điền nếu trường đó là tự tăng (auto-increment).

---

## 2. First — Lấy bản ghi đầu tiên

**GORM:**

```go
var user User
db.First(&user)
```

**SQL:**

```sql
SELECT * FROM users ORDER BY id LIMIT 1;
```

---

## 3. First theo ID

**GORM:**

```go
db.First(&user, 10)
```

**SQL:**

```sql
SELECT * FROM users WHERE id = 10 LIMIT 1;
```

---

## 4. Take — Lấy ngẫu nhiên 1 bản ghi

**GORM:**

```go
db.Take(&user)
```

**SQL:**

```sql
SELECT * FROM users LIMIT 1;
```

> **Khác gì `First`?**
>
> - `First`: Có `ORDER BY id`.
> - `Take`: Không có `ORDER BY`, nhanh hơn một chút vì không cần sắp xếp.

---

## 5. Last — Lấy bản ghi cuối cùng

**GORM:**

```go
db.Last(&user)
```

**SQL:**

```sql
SELECT * FROM users ORDER BY id DESC LIMIT 1;
```

---

## 6. Find — Lấy nhiều bản ghi

**GORM:**

```go
var users []User
db.Find(&users)
```

**SQL:**

```sql
SELECT * FROM users;
```

---

## 7. Where — Điều kiện lọc

**GORM:**

```go
db.Where("age > ?", 18).Find(&users)
```

**SQL:**

```sql
SELECT * FROM users WHERE age > 18;
```

---

## 8. Where + First

**GORM:**

```go
db.Where("email = ?", "a@gmail.com").First(&user)
```

**SQL:**

```sql
SELECT * FROM users WHERE email = 'a@gmail.com' LIMIT 1;
```

---

## 9. Select — Chỉ định cột cần lấy

**GORM:**

```go
db.Select("name", "email").Find(&users)
```

**SQL:**

```sql
SELECT name, email FROM users;
```

---

## 10. Omit — Loại bỏ cột khi thao tác

**GORM:**

```go
db.Omit("email").Create(&user)
```

**SQL:**

```sql
INSERT INTO users (name, age) VALUES (...);
```

---

## 11. Update — Cập nhật một trường

**GORM:**

```go
db.Model(&user).Update("name", "new name")
```

**SQL:**

```sql
UPDATE users SET name = 'new name' WHERE id = 1;
```

---

## 12. Updates — Cập nhật nhiều trường

**GORM:**

```go
db.Model(&user).Updates(User{
    Name: "A",
    Age:  30,
})
```

**SQL:**

```sql
UPDATE users SET name = 'A', age = 30 WHERE id = 1;
```

---

## 13. Save — Insert hoặc Update

**GORM:**

```go
db.Save(&user)
```

**Logic:**

- Nếu `ID = 0` → Thực hiện `INSERT`.
- Nếu `ID != 0` → Thực hiện `UPDATE`.

---

## 14. Delete — Xóa bản ghi

**GORM:**

```go
db.Delete(&user)
```

**SQL:**

```sql
DELETE FROM users WHERE id = 1;
```

---

## 15. Delete theo điều kiện

**GORM:**

```go
db.Where("age < ?", 18).Delete(&User{})
```

**SQL:**

```sql
DELETE FROM users WHERE age < 18;
```

---

## 16. Count — Đếm số lượng

**GORM:**

```go
var total int64
db.Model(&User{}).Count(&total)
```

**SQL:**

```sql
SELECT count(*) FROM users;
```

---

## 17. Order — Sắp xếp

**GORM:**

```go
db.Order("age desc").Find(&users)
```

**SQL:**

```sql
SELECT * FROM users ORDER BY age DESC;
```

---

## 18. Limit — Giới hạn số lượng

**GORM:**

```go
db.Limit(5).Find(&users)
```

**SQL:**

```sql
SELECT * FROM users LIMIT 5;
```

---

## 19. Offset — Phân trang

**GORM:**

```go
db.Offset(10).Limit(5).Find(&users)
```

**SQL:**

```sql
SELECT * FROM users LIMIT 5 OFFSET 10;
```

---

## 20. Joins — Kết bảng

**GORM:**

```go
db.Joins("JOIN orders ON orders.user_id = users.id").Find(&users)
```

**SQL:**

```sql
SELECT * FROM users JOIN orders ON orders.user_id = users.id;
```

---

## 21. Preload — Eager Loading (Nạp trước quan hệ)

**GORM:**

```go
db.Preload("Orders").Find(&users)
```

**SQL thực tế:**
GORM sẽ chạy 2 truy vấn riêng biệt:

1. `SELECT * FROM users;`
2. `SELECT * FROM orders WHERE user_id IN (...);`

---

## 22. Raw SQL — Truy vấn thuần

**GORM:**

```go
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
```

---

## 23. Exec — Thực thi SQL không trả về dữ liệu

**GORM:**

```go
db.Exec("UPDATE users SET age = age + 1")
```

---

## 24. Transaction — Giao dịch

**GORM:**

```go
tx := db.Begin()
tx.Create(&user)
tx.Commit()
```

---

## 25. Rollback

**GORM:**

```go
tx := db.Begin()
if err != nil {
    tx.Rollback()
}
```

---

## 26. Debug SQL — Xem log SQL thực tế

**GORM:**

```go
db.Debug().Find(&users)
```

Nó sẽ in log SQL thật ra terminal để bạn kiểm tra.

---

## 27. Error Handling — Xử lý lỗi

**GORM:**

```go
result := db.First(&user)
if result.Error != nil {
    // Xử lý lỗi
}
```

---

## 28. RowsAffected — Số dòng bị ảnh hưởng

**GORM:**

```go
result := db.Delete(&user)
fmt.Println(result.RowsAffected)
```

---

## 29. AutoMigrate — Tự động đồng bộ Schema

**GORM:**

```go
db.AutoMigrate(&User{})
```

Tự động tạo hoặc cập nhật cấu trúc bảng dựa trên Model.

---

## 30. Model — Chỉ định Table

**GORM:**

```go
db.Model(&User{})
```

Chỉ định bảng mà bạn muốn thực hiện các thao tác tiếp theo.

---

## Tổng hợp Mapping GORM ↔ SQL

| GORM         | SQL         | Ghi chú                 |
| :----------- | :---------- | :---------------------- |
| `Create`     | `INSERT`    | Thêm mới                |
| `First/Find` | `SELECT`    | Truy vấn dữ liệu        |
| `Where`      | `WHERE`     | Điều kiện lọc           |
| `Update`     | `UPDATE`    | Cập nhật                |
| `Delete`     | `DELETE`    | Xóa                     |
| `Count`      | `COUNT`     | Đếm                     |
| `Joins`      | `JOIN`      | Kết bảng                |
| `Order`      | `ORDER BY`  | Sắp xếp                 |
| `Limit`      | `LIMIT`     | Giới hạn                |
| `Offset`     | `OFFSET`    | Phân trang              |
| `Raw`        | SQL thuần   | Truy vấn trực tiếp      |
| `Exec`       | Execute SQL | Chạy SQL không trả data |
| `Begin`      | `BEGIN`     | Bắt đầu Transaction     |
| `Commit`     | `COMMIT`    | Hoàn tất Transaction    |
| `Rollback`   | `ROLLBACK`  | Hoàn tác Transaction    |

---

## Lưu ý cực kỳ quan trọng

### 1. Record Not Found

Khi sử dụng `First`, `Take`, `Last`, nếu không tìm thấy dữ liệu, GORM sẽ trả về lỗi `ErrRecordNotFound`.
Bạn nên kiểm tra như sau:

```go
if errors.Is(err, gorm.ErrRecordNotFound) {
    // Xử lý khi không tìm thấy
}
```

### 2. Sử dụng Pointer (&)

Các hàm `Create`, `Update`, `Find` thường yêu cầu truyền vào một **pointer** (`&user`).

- Điều này cho phép GORM sửa đổi object gốc (ví dụ: điền giá trị `ID` sau khi `Create`).
- Ví dụ: `db.Create(&user)` để sau đó có thể lấy được `user.ID`.
