# Users (Profile)

File liên quan: `internal/entity/user.go`, `internal/usecase/contracts.go`,
`internal/usecase/user/{user.go,user_test.go}`, `internal/repo/contracts.go`,
`internal/repo/persistent/user/{model.go,user.go}`, `internal/controller/restapi/v1/user.go`,
`internal/controller/restapi/v1/response/user.go`

> User được tạo ra duy nhất qua `POST /api/v1/auth/register` — xem
> [01-authentication.md](01-authentication.md). File này chỉ nói về việc **đọc lại** profile sau
> khi đã đăng nhập.

## Tính năng

`entity.User` là bảng `users` — vừa là tài khoản đăng nhập (có `password_hash`, xem
authentication) vừa là nơi giữ thông tin hiển thị (`full_name`). Không có khái niệm "profile tách
riêng khỏi tài khoản" — một bảng duy nhất cho cả hai.

Endpoint duy nhất: `GET /api/v1/users/me` (cần `Authorization: Bearer <access_token>`).

## Vận hành

1. Request đi qua `authMW` → có `user_id` trong context (xem
   [01-authentication.md](01-authentication.md)).
2. `V1.getMe` (`controller/restapi/v1/user.go`) gọi `userUsecase.GetProfile(ctx, userID)` —
   `userUsecase` có kiểu `usecase.User` (interface ở `usecase/contracts.go`), implementation ở
   `usecase/user/user.go`.
3. `GetProfile` chỉ đơn giản `repo.FindByID` — không tồn tại thì map thành `entity.ErrNotFound` →
   handler trả `404` (trường hợp này gần như không xảy ra trong thực tế, vì `user_id` lấy từ chính
   access token đã verify hợp lệ — chỉ xảy ra nếu user bị xoá thẳng khỏi DB sau khi token đã phát
   hành).
4. Handler map `entity.User` sang `response.User` (DTO JSON riêng) trước khi trả về —
   **`PasswordHash` không nằm trong DTO này nên không bao giờ bị serialize ra ngoài**, dù entity
   có field đó.

## Giới hạn quan trọng

- Không có `PUT /users/me` — không cách nào đổi `full_name` hay email sau khi đăng ký qua API hiện
  tại (phải sửa thẳng DB).
- Không có đổi password qua API (không có `PUT /users/me/password`).
- Không có endpoint xoá tài khoản hay list toàn bộ users.
