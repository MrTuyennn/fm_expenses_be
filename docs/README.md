# Docs

Tài liệu chi tiết theo từng tính năng của `fm_expenses_be`. Mỗi file mô tả: tính năng làm gì,
luồng vận hành (request → service → DB), và những hành vi/giới hạn cần biết khi dùng hoặc sửa code.

- [features/01-authentication.md](features/01-authentication.md) — đăng ký/đăng nhập/refresh/logout tự triển khai (không qua Supabase Auth)
- [features/02-users.md](features/02-users.md) — đọc profile người dùng hiện tại
- [features/03-categories.md](features/03-categories.md) — danh mục thu/chi
- [features/04-expenses.md](features/04-expenses.md) — giao dịch chi tiêu

## tech/

Tài liệu theo *lớp kỹ thuật* thay vì tính năng — cấu trúc thư mục, pattern kiến trúc đang dùng.

- [tech/01-project-structure.md](tech/01-project-structure.md) — `cmd/`, `internal/`, `pkg/`,
  `migrations/` làm gì
- [tech/02-architecture-pattern.md](tech/02-architecture-pattern.md) — Clean Architecture theo quy
  ước evrone/go-clean-template, lý do chọn, và những gì chủ động không mang theo từ template đó

Xem thêm `README.md` ở root cho bảng API tổng hợp, và `CLAUDE.md` cho các gotcha khi setup/chạy dự án.
