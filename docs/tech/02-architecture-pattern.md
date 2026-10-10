# Architecture pattern — Clean Architecture (quy ước evrone/go-clean-template)

> Trạng thái: **đã áp dụng**. Repo trải qua 2 lần tái cấu trúc:
> 1. Clean Architecture tự thiết kế (`domain` gộp entity + interface, `usecase` phẳng,
>    `repository/postgres`, `delivery/http`).
> 2. Tổ chức lại theo đúng quy ước đặt tên/thư mục của
>    [evrone/go-clean-template](https://github.com/evrone/go-clean-template) — template Clean
>    Architecture cho Go được dùng rộng rãi nhất trong thực tế (~7.7k sao trên GitHub tại thời điểm
>    viết, cập nhật thường xuyên, nhiều dự án production tham khảo).
>
> Tài liệu này mô tả trạng thái hiện tại (sau lần 2).

## Vì sao chọn evrone/go-clean-template làm tham chiếu

So sánh nhanh các lựa chọn phổ biến nhất khi tìm "chuẩn Go project layout" trên GitHub:

| Repo | Vai trò | Có phù hợp áp dụng trực tiếp không |
|---|---|---|
| `golang-standards/project-layout` (~56k sao) | Danh mục thư mục tham khảo, **không phải** template chạy được. Maintainer tự ghi rõ "không phải chuẩn chính thức của Go team", khuyến cáo chỉ lấy phần cần dùng. | Chỉ dùng để tra cứu tên thư mục, không có kiến trúc/convention cụ thể để áp |
| `ardanlabs/service` | Starter kit production-grade, phong cách riêng của Ardan Labs (khoá học Go nổi tiếng) | Ít sao hơn hẳn (~3.9k), ít ví dụ cộng đồng hơn |
| **`evrone/go-clean-template`** | Starter Clean Architecture **chạy được thật**, có ví dụ REST/gRPC/message queue, test, mock | **Chọn cái này** — nhiều sao nhất trong nhóm "template chạy được" (không tính layout thuần), update thường xuyên, đúng phong cách Clean Architecture (Robert C. Martin) mà bài trước đã phân tích |

## 4 lớp, phụ thuộc một chiều

```
controller/restapi  ──depends on──▶  usecase (contracts.go interface)
                                           ▲
usecase/<feature>   ──implements──────────┘
       │
       └──depends on──▶  repo (contracts.go interface)
                                ▲
repo/persistent/<feature> ─implements──┘
```

Khác với lần refactor đầu tiên (gộp entity + interface vào 1 package `domain`), bản này **tách
riêng interface khỏi entity**, đúng quy ước evrone:

- **`entity`** — chỉ chứa struct thuần. Không interface, không tag.
- **`usecase/contracts.go`** — interface mà **controller** (bên tiêu thụ usecase) cần, định nghĩa
  cạnh nơi nó được dùng.
- **`repo/contracts.go`** — interface mà **usecase** (bên tiêu thụ repo) cần, định nghĩa cạnh nơi
  nó được dùng.
- Implementation thật của mỗi interface nằm trong package con cùng tên feature
  (`usecase/category`, `repo/persistent/category`...), implement interface tương ứng và trả về qua
  hàm `New(...)` — bên gọi chỉ thấy kiểu interface, không thấy struct cụ thể.

## Vì sao tách interface khỏi entity thay vì gộp như bản cũ

Bản refactor đầu tiên định nghĩa interface ngay trong `domain` cùng với entity — hoạt động đúng,
nhưng không khớp quy ước phổ biến nhất: hầu hết code Go idiomatic (và chính evrone) định nghĩa
interface **ở nơi tiêu thụ nó** (consumer-side interface), không phải ở nơi entity sống. Tách ra:

- `usecase/contracts.go` nằm trong package `usecase` — vì `usecase` (gián tiếp, qua `controller`)
  là bên cần biết "usecase trông như thế nào", không phải bản thân entity.
- `repo/contracts.go` nằm trong package `repo` — vì `usecase` là bên tiêu thụ `repo`, và interface
  được evrone đặt ngay tại package cung cấp nó (`repo`) để dễ tìm tất cả port liên quan tới lưu trữ
  ở một chỗ, thay vì rải theo từng entity.

Đây là lựa chọn phong cách (có thể tranh luận), không phải đúng/sai tuyệt đối — nhưng khớp với số
đông dự án Go Clean Architecture thực tế nên dễ onboard người mới hơn.

## Những gì **không** mang theo từ evrone/go-clean-template — và vì sao

Template gốc là cho một service production đầy đủ (multi-protocol: REST + gRPC + AMQP + NATS, có
tracing OpenTelemetry, Prometheus metrics, Swagger, dùng Fiber thay vì `net/http`/chi, dùng
pgx+squirrel thay vì GORM, dùng `caarlos0/env` cho config YAML+env). Repo này **chỉ lấy phần quy
ước tổ chức thư mục/interface**, không lấy phần hạ tầng, vì:

| Phần của evrone | Bỏ qua vì |
|---|---|
| gRPC, AMQP (RabbitMQ), NATS | App này chỉ có REST, không cần multi-protocol |
| OpenTelemetry tracing wrapper quanh mỗi usecase/repo | Thêm phức tạp không cần thiết cho app quy mô nhỏ, chưa có nhu cầu observability nâng cao |
| Prometheus metrics, Swagger | Chưa có yêu cầu, có thể thêm sau nếu cần mà không ảnh hưởng kiến trúc |
| Fiber (thay cho chi) | Đổi router là quyết định công nghệ riêng, không liên quan tới Clean Architecture; chi đã chạy ổn |
| pgx + squirrel (thay cho GORM) | Đổi ORM/query builder là quyết định riêng, rủi ro cao hơn lợi ích ở quy mô 3 bảng; GORM đã hoạt động tốt |
| `caarlos0/env` + config theo struct lồng nhau (App/HTTP/PG/GRPC/RMQ/NATS/JWT...) | App chỉ có vài biến môi trường, `os.Getenv`/`time.ParseDuration` thủ công trong `internal/config/config.go` đã đủ, thêm dependency mới là phí |
| `go-playground/validator` cho request DTO | Request hiện tại chỉ vài field, validate thủ công trong usecase (`isValidType`, độ dài password...) đủ dùng |
| Refresh token / rotation | evrone's `usecase.User.Login` chỉ trả 1 access token (24h), **không** có khái niệm refresh token. App này **có** tự thêm refresh token + rotation (xem [../features/01-authentication.md](../features/01-authentication.md)) — phần tự thiết kế, không lấy từ template |

## Lợi ích giữ lại được từ lần refactor trước (không đổi)

1. Test usecase không cần Postgres thật — `usecase/category/category_test.go` dùng mock implement
   `repo.CategoryRepo`.
2. Lỗi domain (`entity.ErrNotFound`, `entity.ErrInvalidCategoryType`...) được map đúng HTTP status
   qua `errors.Is` ở `controller/restapi/v1`, không còn bug "mọi lỗi đều trả 404" như bản gốc.
3. DTO response tách khỏi entity — đổi JSON shape không đụng domain.

## Đánh đổi

- Nhiều file hơn nữa so với cả 2 bản trước — một thay đổi nhỏ giờ có thể chạm 5-6 file
  (`entity`, `usecase/contracts.go`, `usecase/<feature>`, `repo/contracts.go`,
  `repo/persistent/<feature>/model.go`, `repo/persistent/<feature>/<feature>.go`,
  `controller/restapi/v1/<feature>.go`, `request/`, `response/`).
- Việc tách `usecase/contracts.go` và `repo/contracts.go` thành 2 file riêng (thay vì gộp trong
  `domain` như bản cũ) là đánh đổi lấy sự quen thuộc với cộng đồng, không mang lại lợi ích kỹ thuật
  mới so với bản trước — cả 2 bản đều đạt dependency inversion đúng.
