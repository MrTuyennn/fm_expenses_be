# fm_expenses_be

Backend quản lý chi tiêu gia đình, viết bằng Go (net/http chuẩn + [chi](https://github.com/go-chi/chi) router),
dữ liệu lưu trên Postgres (host trên Supabase — chỉ dùng làm DB, không dùng Supabase Auth).
Xác thực (đăng ký/đăng nhập/refresh token) tự triển khai hoàn toàn trên backend.

## Kiến trúc

Clean Architecture, tổ chức theo quy ước của
[evrone/go-clean-template](https://github.com/evrone/go-clean-template) (~7.7k sao, template Clean
Architecture cho Go được dùng rộng rãi nhất trong thực tế) — chỉ lấy phần quy ước đặt tên/thư mục,
không mang theo phần hạ tầng không liên quan (gRPC/AMQP/NATS/tracing/Fiber/pgx+squirrel) vì project
này chỉ cần REST thuần. Chi tiết: [docs/tech/02-architecture-pattern.md](docs/tech/02-architecture-pattern.md).

```
cmd/api/main.go              # entrypoint mỏng: load config rồi gọi app.Run
internal/
  app/app.go                 # composition root thật sự — wiring toàn bộ dependency, start server
  config/                    # load biến môi trường
  middleware/                # middleware xác thực access token (tự ký, không phải Supabase)
  entity/                    # entity thuần — lớp trong cùng, không phụ thuộc gì
  usecase/
    contracts.go             # interface usecase (Auth, User, Category, Expense)
    auth/ user/ category/ expense/  # implementation theo từng feature
  repo/
    contracts.go             # interface repository (UserRepo, RefreshTokenRepo, CategoryRepo, ExpenseRepo)
    persistent/
      user/ refreshtoken/ category/ expense/  # adapter GORM/Postgres theo từng feature
      migrate.go              # AutoMigrate gộp tất cả model
  controller/restapi/        # adapter HTTP (chi router)
    v1/                      # handler, request/response DTO
pkg/
  jwt/                       # ký/parse access token (HS256)
  postgres/                  # mở connection Postgres (dùng lại được ngoài repo này)
  response/                  # helper trả JSON response chuẩn
migrations/0001_init.sql     # schema tham khảo (users, refresh_tokens, categories, expenses)
```

Xem thêm [docs/](docs/) cho tài liệu chi tiết theo từng tính năng (`docs/features/`) và theo lớp kỹ
thuật (`docs/tech/`).

## Yêu cầu

- Go 1.22+
- Một Postgres database (Supabase hoặc bất kỳ Postgres nào khác — không cần tính năng Auth của Supabase)

## Cài đặt

```bash
cp .env.example .env
# điền DATABASE_URL (connection string Postgres)
# điền JWT_SECRET (tự tạo: openssl rand -base64 32)
go mod tidy
make run
```

Server chạy ở `http://localhost:8080`.

## Biến môi trường

| Key | Mô tả |
|---|---|
| `PORT` | Cổng HTTP server (default `8080`) |
| `DATABASE_URL` | Connection string Postgres |
| `JWT_SECRET` | Secret backend tự dùng để ký access token (tự chọn, không liên quan Supabase) |
| `ACCESS_TOKEN_TTL` | Thời hạn access token (default `15m`) |
| `REFRESH_TOKEN_TTL` | Thời hạn refresh token (default `720h` = 30 ngày) |

## Xác thực

Backend tự quản lý toàn bộ vòng đời tài khoản — không còn phụ thuộc Supabase Auth hay bất kỳ
identity provider nào khác. Chi tiết đầy đủ: [docs/features/01-authentication.md](docs/features/01-authentication.md).

1. `POST /api/v1/auth/register` hoặc `POST /api/v1/auth/login` → nhận về `access_token` +
   `refresh_token`.
2. Gắn `Authorization: Bearer <access_token>` vào mọi request tới route cần auth.
3. Khi access token hết hạn (`401`), gọi `POST /api/v1/auth/refresh` với `refresh_token` để lấy cặp
   token mới (refresh token cũ bị thu hồi ngay — single-use/rotation).
4. `POST /api/v1/auth/logout` để thu hồi refresh token (đăng xuất).

## API

| Method | Path | Auth? | Mô tả |
|---|---|---|---|
| GET | `/health` | Không | Health check |
| POST | `/api/v1/auth/register` | Không | Đăng ký, trả về user + token pair |
| POST | `/api/v1/auth/login` | Không | Đăng nhập, trả về token pair |
| POST | `/api/v1/auth/refresh` | Không | Đổi refresh token lấy token pair mới (rotate) |
| POST | `/api/v1/auth/logout` | Không | Thu hồi refresh token |
| GET | `/api/v1/users/me` | Có | Lấy profile người dùng hiện tại |
| GET | `/api/v1/categories` | Có | Danh sách danh mục của người dùng |
| POST | `/api/v1/categories` | Có | Tạo danh mục mới |
| GET | `/api/v1/categories/{id}` | Có | Chi tiết danh mục |
| PUT | `/api/v1/categories/{id}` | Có | Cập nhật danh mục |
| DELETE | `/api/v1/categories/{id}` | Có | Xoá danh mục |
| GET | `/api/v1/expenses` | Có | Danh sách chi tiêu (filter `category_id`, `from`, `to` theo RFC3339) |
| POST | `/api/v1/expenses` | Có | Tạo chi tiêu mới |
| GET | `/api/v1/expenses/{id}` | Có | Chi tiết chi tiêu |
| PUT | `/api/v1/expenses/{id}` | Có | Cập nhật chi tiêu |
| DELETE | `/api/v1/expenses/{id}` | Có | Xoá chi tiêu |

## Database

Mặc định app tự `AutoMigrate` schema khi khởi động (tạo bảng `users`, `refresh_tokens`,
`categories`, `expenses` nếu chưa có). File `migrations/0001_init.sql` giữ schema tương đương, dùng
khi muốn chạy migration tay. Không còn Row Level Security / `auth.uid()` — xem comment đầu file đó
để biết lý do (backend tự làm auth, Postgres không còn hiểu JWT của backend).
