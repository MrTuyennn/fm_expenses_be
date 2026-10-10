# fm_expenses_be

Go backend quản lý chi tiêu gia đình. Chi tiết kiến trúc, API, biến môi trường xem ở `README.md`.
File này chỉ ghi lại những điều **không hiển nhiên** khi làm việc với repo.

## Chạy project

Entrypoint nằm ở `cmd/api`, không phải root — `go run .` sẽ báo "no go files listed".

```bash
go run ./cmd/api
# hoặc
make run
```

## `.env`

Không có file `.env` mặc định (đã bị gitignore). Copy từ `.env.example` và điền:

- `DATABASE_URL` — connection string Postgres thật. Supabase ở đây **chỉ đóng vai trò host
  Postgres**, không phải identity provider — lấy ở Supabase Dashboard → Project Settings →
  Database → Connect → "Direct connection" hoặc "Session pooler". Dashboard hiển thị dạng
  `postgresql://postgres:[YOUR-PASSWORD]@...` — phải thay `[YOUR-PASSWORD]` bằng password thật
  (bỏ luôn dấu ngoặc vuông `[ ]`), không phải điền nguyên văn placeholder.
- `JWT_SECRET` — secret **tự chọn**, dùng để backend tự ký access token của chính nó (không liên
  quan gì tới Supabase Auth nữa — xem phần auth bên dưới). Generate bằng
  `openssl rand -base64 32`.

## Auth tự triển khai hoàn toàn trên backend (không còn dùng Supabase Auth)

`internal/usecase/auth` tự làm register/login/refresh/logout — xem
[docs/features/01-authentication.md](docs/features/01-authentication.md) để biết chi tiết luồng
(bcrypt hash password, access token JWT ngắn hạn ký bằng `JWT_SECRET`, refresh token random lưu
hash trong bảng `refresh_tokens`, rotate mỗi lần refresh). Supabase project vẫn được dùng, nhưng
**chỉ để host Postgres** — không còn liên quan gì tới việc verify JWT hay quản lý user nữa. Đừng
nhầm lẫn với phiên bản trước đó của backend (đã bị thay thế), lúc đó dùng Supabase Auth JWT trực
tiếp và cần `SUPABASE_JWT_SECRET`/JWKS từ Supabase dashboard.

## Log "SLOW SQL >= 200ms" lúc khởi động

Bình thường, không phải lỗi. Do `database.AutoMigrate` (GORM) chạy loạt query kiểm tra schema mỗi
lần start, cộng độ trễ mạng tới Supabase (region xa, ví dụ `ap-south-1`). Chỉ cần lo nếu thấy dòng
`error` hoặc server không in ra `fm_expenses_be listening on :PORT`.
