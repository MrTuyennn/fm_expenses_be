# Expenses

File liên quan: `internal/entity/expense.go`, `internal/usecase/contracts.go`,
`internal/usecase/expense/expense.go`, `internal/repo/contracts.go`,
`internal/repo/persistent/expense/{model.go,expense.go}`,
`internal/controller/restapi/v1/expense.go`,
`internal/controller/restapi/v1/{request,response}/expense.go`

## Tính năng

Giao dịch chi tiêu (hoặc thu nhập, tuỳ `category.Type`) của user, mỗi expense gắn với đúng một
category. `usecase/expense` là nơi duy nhất trong repo phụ thuộc chéo 2 port của 2 feature khác
nhau cùng lúc (`repo.ExpenseRepo` **và** `repo.CategoryRepo`) để validate quyền sở hữu category
trước khi cho tạo/sửa expense.

Endpoints: `GET/POST /api/v1/expenses`, `GET/PUT/DELETE /api/v1/expenses/{id}`.

## Model

| Cột | Ghi chú |
|---|---|
| `id` | UUID, tự sinh |
| `user_id` | chủ sở hữu |
| `category_id` | FK tới `categories.id`, `on delete restrict` |
| `amount` | `numeric(14,2)` — không có validate `> 0` ở Go lẫn DB, có thể lưu số âm hoặc 0 |
| `currency` | mặc định `'VND'`, không validate theo ISO 4217 |
| `note` | tự do, optional |
| `spent_at` | timestamp, **bắt buộc**, parse từ string RFC3339 ở tầng handler |

## Vận hành

### Create / Update — validate category thuộc về user

`usecase/expense.Create` và `Update` đều gọi `categoryRepo.FindByID(ctx, categoryID, userID)`
trước khi ghi expense (`usecase/expense/expense.go`) — nếu category không tồn tại **hoặc thuộc user
khác**, trả `entity.ErrCategoryInvalid` → `writeExpenseError` trong
`controller/restapi/v1/expense.go` map thành `400 category does not belong to user`. Đây là nơi
duy nhất dữ liệu của 2 feature (category + expense) được đối chiếu chéo, qua interface
`repo.CategoryRepo` chứ không phải struct cụ thể — không có foreign-key check nào ở tầng ứng dụng
cho `user_id` giữa 2 bảng ngoài việc này.

### Parse `spent_at`

Handler parse `spent_at` bằng `time.Parse(time.RFC3339, ...)` (`controller/restapi/v1/expense.go`)
— sai format (ví dụ chỉ gửi `2026-10-08` không có giờ/timezone) sẽ bị từ chối ngay ở tầng handler
với `400 invalid spent_at, expected RFC3339`, **trước khi** chạm tới usecase/DB.

### List — filter

`GET /expenses` hỗ trợ 3 query param, cộng dồn bằng AND
(`repo/persistent/expense/expense.go` — `ListByUser`):

- `category_id` — lọc đúng category
- `from` / `to` — khoảng `spent_at`, inclusive cả 2 đầu (`>=` và `<=`)

**Lưu ý (chưa fix qua 2 lần refactor, vẫn giữ nguyên hành vi từ bản gốc):** nếu `from`/`to` gửi sai
format RFC3339, handler **âm thầm bỏ qua** filter đó (`controller/restapi/v1/expense.go` — `List`,
lỗi parse bị nuốt bằng `if t, err := ...; err == nil`) thay vì trả lỗi 400 — nghĩa là client gửi sai
format `from`/`to` sẽ nhận về list **không lọc theo thời gian** thay vì một thông báo lỗi, dễ gây
hiểu nhầm là API đang bug.

Kết quả sort theo `spent_at desc`.

### Update — full replace

Giống category, `Update` là full-replace: load expense hiện tại (đảm bảo thuộc user), ghi đè toàn
bộ `CategoryID/Amount/Currency/Note/SpentAt` rồi `Save` (`usecase/expense/expense.go` — `Update`).
Không có partial update (PATCH). Lỗi "không tìm thấy expense" map đúng `404` qua
`entity.ErrNotFound` + `errors.Is` trong `writeExpenseError` (bản gốc trước khi áp dụng Clean
Architecture trả `400` cho trường hợp này).

### Delete

Xoá trực tiếp theo `id + user_id`, không có check gì thêm, không soft-delete (GORM `Save`/`Delete`
ở đây là hard delete vì model không có field `gorm.DeletedAt`).

## Giới hạn quan trọng

- Không validate `amount > 0` — có thể tạo expense âm hoặc bằng 0.
- Không validate `currency` theo danh sách mã tiền tệ hợp lệ.
- Filter `from`/`to` sai format bị nuốt lỗi âm thầm (xem trên) — nên sửa nếu cần UX chặt hơn.
- Đổi `category_id` khi Update sẽ re-validate quyền sở hữu category mới, nhưng **không** kiểm tra
  category mới có cùng `type` (`income`/`expense`) với trước đó hay không — có thể đổi một expense
  đang gắn category `expense` sang category `type = income` mà không có cảnh báo gì.
