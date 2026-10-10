# Authentication

File liên quan: `internal/entity/{user.go,auth.go,errors.go}`, `internal/usecase/contracts.go`,
`internal/usecase/auth/{auth.go,auth_test.go}`, `internal/repo/contracts.go`,
`internal/repo/persistent/user/`, `internal/repo/persistent/refreshtoken/`,
`internal/controller/restapi/v1/auth.go`, `internal/controller/restapi/v1/{request,response}/auth.go`,
`internal/middleware/auth.go`, `pkg/jwt/jwt.go`

> **Lịch sử:** phiên bản trước của backend uỷ quyền hoàn toàn cho Supabase Auth (verify JWT do
> Supabase cấp bằng `SUPABASE_JWT_SECRET`). Backend giờ **tự làm toàn bộ auth** — không còn gọi gì
> tới Supabase Auth, Postgres (dù host trên Supabase hay nơi khác) chỉ còn là nơi lưu dữ liệu.

## Tính năng

Đăng ký, đăng nhập, refresh token, đăng xuất — toàn bộ nằm trong `usecase/auth`, không phụ thuộc
identity provider bên ngoài. Password hash bằng bcrypt. Access token là JWT ngắn hạn (stateless,
tự verify không cần hit DB). Refresh token là chuỗi random lưu hash trong DB (revocable — điều một
JWT thuần không làm được).

Endpoints (không cần `Authorization` header):
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

Mọi route khác dưới `/api/v1/*` yêu cầu header:
```
Authorization: Bearer <access_token>
```

## Vận hành

### Register

1. `POST /auth/register` body `{"email","password","full_name"}`.
2. `usecase/auth.Register` validate độ dài password (`entity.ErrPasswordTooShort` nếu < 8 ký tự),
   hash bằng `bcrypt.GenerateFromPassword` (cost mặc định), insert vào bảng `users` qua
   `repo.UserRepo.Create`.
3. Email trùng → repo map lỗi Postgres unique violation (`pgconn.PgError` code `23505`) thành
   `entity.ErrUserAlreadyExists` (`internal/repo/persistent/user/user.go`) → handler trả `409`.
4. Đăng ký thành công **tự động đăng nhập luôn** — response trả cả `user` lẫn `tokens` (không cần
   gọi `/login` riêng sau khi register). Đây là lựa chọn có chủ đích để tiện cho mobile client, có
   thể đánh đổi nếu sau này cần flow xác thực email trước khi cho đăng nhập.

### Login

1. `POST /auth/login` body `{"email","password"}`.
2. Tìm user theo email (`repo.UserRepo.FindByEmail`); không có hoặc `bcrypt.CompareHashAndPassword`
   sai → **cùng một lỗi** `entity.ErrInvalidCredentials` → `401 invalid credentials` — không phân
   biệt "email không tồn tại" hay "sai password" để tránh lộ thông tin email nào đã đăng ký.

### Access token (JWT, stateless)

`pkg/jwt.Manager.GenerateToken(userID)` ký JWT HS256 bằng `JWT_SECRET`, claims chỉ có `sub`
(user ID) + `exp`. TTL mặc định `ACCESS_TOKEN_TTL=15m` (cấu hình qua env, xem `internal/config`).
`internal/middleware/auth.go` verify token này cho mọi route cần auth, inject `user_id` vào
context — handler đọc qua `middleware.UserIDFromContext(ctx)`.

**Vì JWT là stateless, access token KHÔNG thể revoke trước hạn.** Logout chỉ thu hồi refresh token
— access token đã phát hành vẫn dùng được tới khi tự hết hạn (tối đa 15 phút theo default). Đây là
đánh đổi có chủ đích: đổi lấy việc middleware verify token không cần hit DB mỗi request.

### Refresh token (opaque, revocable, lưu DB)

1. `issueTokenPair` (dùng chung cho Register/Login/Refresh) sinh refresh token: 32 byte random
   (`crypto/rand`), encode base64 URL-safe — **không phải JWT**, client không đọc được gì từ nó.
2. Chỉ **SHA-256 hash** của refresh token được lưu vào bảng `refresh_tokens`
   (`internal/repo/persistent/refreshtoken/`), kèm `expires_at` (`REFRESH_TOKEN_TTL`, default
   `720h` = 30 ngày). Token thật chỉ trả về client đúng 1 lần lúc phát hành, không bao giờ lưu lại
   dạng plaintext — rò DB không lộ được refresh token nào đang hoạt động.
3. `POST /auth/refresh` body `{"refresh_token"}`: hash token nhận được, tìm trong DB qua
   `FindValidByHash` (điều kiện `revoked_at IS NULL AND expires_at > now()`). Không tìm thấy/hết
   hạn/đã bị revoke → `entity.ErrRefreshTokenInvalid` → `401`.
4. **Rotation**: mỗi lần refresh thành công, token cũ bị `Revoke` ngay lập tức, token mới được phát
   hành thay thế — refresh token là **single-use**. Nếu một refresh token bị lộ và kẻ tấn công dùng
   nó trước người dùng thật, lần refresh tiếp theo của người dùng thật sẽ thất bại (dấu hiệu phát
   hiện token bị đánh cắp), thay vì cả 2 bên cùng dùng chung một token sống mãi.

### Logout

`POST /auth/logout` body `{"refresh_token"}` → hash rồi `Revoke` token đó. Token không tồn tại/đã
revoke thì vẫn trả `204` (im lặng) — tránh để client suy luận được token có từng tồn tại hay không.

## Giới hạn quan trọng

- Không có xác thực email, không có "forgot password" / reset password flow — tự triển khai hoàn
  toàn đồng nghĩa phải tự xây các flow này nếu cần (Supabase Auth trước đây có sẵn).
- Không có rate limiting cho `/auth/login` hay `/auth/register` — dễ bị brute-force password hoặc
  spam tạo tài khoản nếu deploy production mà chưa thêm rate limit (ví dụ qua middleware riêng hoặc
  API gateway phía trước).
- Không có cơ chế "revoke tất cả refresh token của 1 user" (ví dụ khi nghi ngờ bị chiếm tài khoản)
  — `RefreshTokenRepo` hiện chỉ revoke theo từng token, chưa có `RevokeAllForUser`.
- `full_name` chỉ set được lúc `Register`, chưa có endpoint update profile sau đó.
