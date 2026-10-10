# Project structure

Tài liệu kỹ thuật — giải thích từng thư mục làm gì và vì sao chia như vậy. Khác với các file trong
`docs/features/` (mô tả theo *tính năng*), file ở `docs/tech/` mô tả theo *lớp kỹ thuật*.

Repo tổ chức theo Clean Architecture, dùng quy ước đặt tên/thư mục của
[evrone/go-clean-template](https://github.com/evrone/go-clean-template) — template Clean
Architecture cho Go phổ biến nhất trong thực tế (~7.7k sao, nhiều dự án tham khảo). Chi tiết lý do
chọn và những gì **không** mang theo từ template đó xem
[02-architecture-pattern.md](02-architecture-pattern.md).

```
cmd/api/main.go                 # entrypoint mỏng: load config rồi gọi app.Run(cfg)
internal/
  app/app.go                    # composition root — wiring toàn bộ dependency, start HTTP server
  config/                       # đọc .env / biến môi trường
  middleware/                   # cross-cutting concern (verify access token)
  entity/                       # entity thuần — KHÔNG phụ thuộc gì cả
    user.go  auth.go  category.go  expense.go  errors.go
  usecase/
    contracts.go                # interface usecase: Auth, User, Category, Expense
    auth/auth.go  auth/auth_test.go       # register/login/refresh/logout
    user/user.go  user/user_test.go       # chỉ còn GetProfile
    category/category.go  category/category_test.go
    expense/expense.go  expense/expense_test.go
  repo/
    contracts.go                # interface repository: UserRepo, RefreshTokenRepo, CategoryRepo, ExpenseRepo
    persistent/
      migrate.go                 # AutoMigrate gộp tất cả model
      user/model.go  user/user.go              # bảng `users` (có password_hash)
      refreshtoken/model.go  refreshtoken/refreshtoken.go  # bảng `refresh_tokens`
      category/model.go  category/category.go
      expense/model.go  expense/expense.go
  controller/restapi/
    router.go                   # mount /health, group /api/v1/auth (public) + /api/v1/* (auth middleware)
    v1/
      controller.go              # struct V1 giữ 4 usecase
      router.go                  # đăng ký route (NewAuthRoutes + NewRoutes)
      auth.go  user.go  category.go  expense.go   # handler
      request/   auth.go  category.go  expense.go # DTO request body
      response/  auth.go  user.go  category.go  expense.go  # DTO response
pkg/
  jwt/jwt.go                    # ký/verify access token (HS256, stateless)
  postgres/postgres.go          # mở connection Postgres, dùng lại được ngoài repo này
  response/response.go          # helper JSON envelope {"data":...} / {"error":...}
migrations/0001_init.sql        # schema tham khảo (users, refresh_tokens, categories, expenses)
```

## `cmd/api`

Entrypoint của binary — chỉ còn 2 dòng: `config.Load()` rồi `app.Run(cfg)`. Mọi việc wiring thật sự
nằm ở `internal/app`, không phải ở đây (khác với các bản Clean Architecture Go đơn giản hơn thường
nhét hết vào `main.go`).

## `internal/app`

**Composition root thật sự** — hàm `Run(cfg)` mở DB connection, chạy `persistent.Migrate`, dựng
`jwt.Manager` + `middleware.Auth` từ `cfg.JWTSecret`, rồi dựng chuỗi dependency cho từng feature
theo đúng thứ tự `repo → usecase → controller`, cuối cùng start HTTP server. Đây là nơi duy nhất
trong repo "biết" cả 4 lớp cùng lúc.

## `internal/entity`

Lớp trong cùng — struct thuần Go (`User`, `RefreshToken`, `TokenPair`, `Category`, `Expense`),
không tag `gorm`/`json`, không import gì ngoài thư viện chuẩn. `User.PasswordHash` tồn tại ở đây
nhưng **không** có trong `response.User` (DTO riêng ở `controller`) — tách domain khỏi DTO là lý do
field nhạy cảm này không bao giờ bị serialize nhầm ra JSON. Đổi ORM, đổi DB, đổi web framework đều
không chạm vào package này.

## `internal/usecase`

`contracts.go` khai báo interface (`Auth`, `User`, `Category`, `Expense`) mà `controller` phụ thuộc
vào. Mỗi feature có package con riêng chứa struct `UseCase` implement interface tương ứng, chỉ phụ
thuộc `internal/repo` qua interface, không biết `repo/persistent` hay `controller` tồn tại.
`usecase/auth` còn phụ thuộc thêm `pkg/jwt.Manager` (ký access token) — là usecase duy nhất phụ
thuộc một package ở `pkg/` thay vì chỉ `repo`, vì việc ký JWT là một "thuật toán" thuần tuý, không
phải I/O, nên coi như một utility được tiêm vào chứ không qua interface riêng. Nhờ tách theo
interface, `usecase/auth/auth_test.go`, `usecase/category/category_test.go`... test được logic
bằng repository giả lập trong RAM, không cần Postgres thật.

## `internal/repo`

`contracts.go` khai báo interface (`UserRepo`, `RefreshTokenRepo`, `CategoryRepo`, `ExpenseRepo`)
mà `usecase` phụ thuộc vào. `repo/persistent/<feature>` implement bằng GORM: mỗi feature có
`model.go` (struct GORM riêng, tag `gorm`, có hàm `toEntity()`/`modelFromEntity()` map 2 chiều) và
file cùng tên feature chứa struct `Repo`. `repo/persistent/user/user.go` còn map lỗi Postgres unique
violation (`pgconn.PgError` code `23505`) thành `entity.ErrUserAlreadyExists` — ví dụ cụ thể cho
việc adapter chịu trách nhiệm dịch lỗi hạ tầng sang lỗi domain. `repo/persistent/migrate.go` gom
`AutoMigrate` cho mọi model — nằm ở đây, không phải `pkg/postgres`, vì chỉ package này mới cần biết
cấu trúc bảng cụ thể.

## `internal/controller/restapi`

Adapter phía transport. `router.go` ở cấp `restapi` mount `/health` (không cần auth), group
`/api/v1/auth/*` (public — `v1.NewAuthRoutes`), và group `/api/v1/*` còn lại (có
`authMW.Middleware` — `v1.NewRoutes`), dùng `chi.Router.Group` để tách 2 middleware stack trong
cùng 1 prefix `/api/v1`. `v1/` chứa handler thật: decode JSON qua DTO ở `request/`, gọi usecase
interface, map entity sang DTO ở `response/` trước khi trả JSON. `writeCategoryError`/
`writeExpenseError` map lỗi domain sang đúng HTTP status bằng `errors.Is`.

## `pkg/jwt`, `pkg/postgres`, `pkg/response`

`pkg/` là mã dùng lại được, không bị giới hạn bởi `internal/`. `pkg/jwt.Manager` ký/verify access
token HS256, không biết gì về user/DB — chỉ là một bộ ký/verify JWT thuần. `pkg/postgres` chỉ mở
connection GORM thuần (không biết schema cụ thể — đó là việc của `repo/persistent`). `pkg/response`
là helper JSON response mỏng, không chứa logic nghiệp vụ.

## `migrations/`

Không đổi cách dùng — chỉ tham khảo, app vẫn dùng `persistent.Migrate` để tự tạo bảng. Nội dung đã
đổi: không còn FK tới `auth.users` của Supabase, không còn Row Level Security (xem comment đầu file
SQL — RLS dựa trên `auth.uid()` không còn ý nghĩa gì khi backend tự làm auth, Postgres không hiểu
JWT của backend).
