# Categories

File liên quan: `internal/entity/category.go`, `internal/usecase/contracts.go`,
`internal/usecase/category/{category.go,category_test.go}`, `internal/repo/contracts.go`,
`internal/repo/persistent/category/{model.go,category.go}`,
`internal/controller/restapi/v1/category.go`,
`internal/controller/restapi/v1/{request,response}/category.go`

## Tính năng

Danh mục thu/chi do từng user tự tạo (ví dụ: "Ăn uống", "Lương", "Di chuyển") để gắn vào mỗi
expense. CRUD đầy đủ, luôn scope theo `user_id` — user A không bao giờ thấy/sửa được category của
user B.

Endpoints: `GET/POST /api/v1/categories`, `GET/PUT/DELETE /api/v1/categories/{id}`.

## Model

| Cột | Ghi chú |
|---|---|
| `id` | UUID, tự sinh (`gen_random_uuid()`) |
| `user_id` | chủ sở hữu, không đổi được sau khi tạo |
| `name` | tên danh mục |
| `icon` | tuỳ chọn, string tự do (app client tự quy ước, ví dụ emoji hoặc icon key) |
| `type` | bắt buộc `'income'` hoặc `'expense'`, enforce bằng **Postgres check constraint**, không enforce ở tầng Go |

## Vận hành

- **Create**: decode JSON body (`request.Category`) → `categoryUsecase.Create` validate `Type`
  tường minh trước (`usecase/category/category.go` — `isValidType`), trả `400 category type must
  be income or expense` ngay nếu sai, không đợi DB check constraint từ chối nữa. Hợp lệ thì tạo
  `entity.Category{UserID, Name, Icon, Type}` rồi insert qua repo (`repo/persistent/category`).
- **List**: trả toàn bộ category của user, sắp xếp `created_at desc`
  (`repo/persistent/category/category.go` — `ListByUser`).
- **Get/Update/Delete**: mọi query đều có điều kiện `id = ? AND user_id = ?`
  (`repo/persistent/category/category.go`) — đây là cơ chế authorization duy nhất, không có check
  role/permission nào khác. Record không tồn tại hoặc thuộc user khác đều trả `404` giống nhau
  (không lộ thông tin category của người khác tồn tại hay không).
- **Update**: validate `Type` giống `Create`, load record hiện tại trước (để xác nhận thuộc về
  user), ghi đè cả 3 field `Name/Icon/Type` rồi `Save` toàn bộ struct (`usecase/category/category.go`
  — `Update`) — **không phải partial update**, client phải gửi đủ cả 3 field mỗi lần `PUT`, thiếu
  field nào sẽ bị ghi thành rỗng.

## Giới hạn quan trọng

- ~~Lỗi khi Update bị map sai thành 404 cho mọi loại lỗi~~ — **đã fix** khi refactor sang Clean
  Architecture: `categoryUsecase.Update` giờ trả `entity.ErrNotFound` chỉ khi record không tồn
  tại/không thuộc user, `entity.ErrInvalidCategoryType` khi `type` sai, hoặc lỗi gốc cho trường hợp
  khác; `controller/restapi/v1` dùng `errors.Is` (`writeCategoryError` trong
  `controller/restapi/v1/category.go`) để map đúng status: `404` / `400` / `500` tương ứng — không
  còn luôn trả `404` như trước.
- Không có check khi xoá category đang được expense nào đó tham chiếu ở tầng Go — việc này do DB
  đảm nhiệm: `migrations/0001_init.sql:30` đặt `on delete restrict` trên FK `expenses.category_id`,
  nên xoá category đang có expense sẽ bị Postgres từ chối (lỗi FK constraint), và lỗi đó hiện lộ ra
  thành `500 failed to delete category` (generic) chứ không có message rõ ràng kiểu "category đang
  được dùng".
- Không có validate độ dài `name`/`icon`, không có unique constraint theo `(user_id, name)` — user
  có thể tạo nhiều category trùng tên.
