# Food Ordering System — IT096IIU Lab 01 (Go Foundation)

Báo cáo kỹ thuật và hướng dẫn sử dụng cho bài thực hành Lab 01, môn IT096IIU — Netcentric Programming, International University. Hệ thống mô phỏng một quy trình đặt món ăn (food ordering) tối giản cho một quán ăn trong khuôn viên trường, được dùng làm phương tiện để thực hành thiết kế dữ liệu (data design), slice/map trong Go, kiểm thử tự động (automated testing), và khả năng thích ứng khi yêu cầu thay đổi (requirement change adaptability).

## 1. Giới thiệu (Introduction)

Lab này tuân theo quy trình 5 giai đoạn do đề bài quy định: **THINK → BUILD → VERIFY → ADAPT → GROUP CHECK**, dựa trên nguyên tắc cốt lõi *"Think together. Build with AI. Verify together. Own the result."* Mỗi giai đoạn tương ứng với một tài liệu riêng trong repository này, để tách biệt rõ ràng giữa (a) đặc tả yêu cầu, (b) mã nguồn hiện thực, và (c) minh chứng kiểm thử/thích ứng — phục vụ cho việc đánh giá (assessment) theo từng mục tiêu độc lập.

Phạm vi chức năng của hệ thống, theo đúng Scenario trong đề bài:

1. Hiển thị các món ăn còn hàng (available).
2. Tìm món ăn theo tên.
3. Thêm món ăn vào một order.
4. Từ chối số lượng (quantity) không hợp lệ.
5. Từ chối món ăn không còn hàng.
6. Tính subtotal của order.
7. Hiển thị order summary.
8. Áp dụng discount khi yêu cầu thay đổi.

## 2. Cấu trúc thư mục (Project structure)

```
Lab1/
├── CLAUDE.md                                        # Kế hoạch chi tiết, dịch toàn bộ đề bài (mục 1-8)
├── README.md                                        # Tài liệu này
├── VERIFY.md                                        # Giai đoạn 4 — VERIFY: test case, defect injection, AI review
├── ADAPT.md                                         # Giai đoạn 5 — ADAPT: discount requirement, test, reflection
├── GROUP_CHECK.md                                   # Giai đoạn 6 — GROUP CHECK: chuẩn bị trả lời câu hỏi giảng viên
├── go.mod                                           # Khai báo Go module (module foodordering)
├── main.go                                          # Toàn bộ mã nguồn hiện thực (struct, slice, map, business logic, CLI)
├── main_test.go                                     # Bộ test tự động (go test), tương ứng các case VERIFY/ADAPT
└── IT096IIU_Lab01_Food_Ordering_Student (1).pdf      # Đề bài gốc do giảng viên cung cấp
```

Nguyên tắc tổ chức: **mã nguồn** (`main.go`, `main_test.go`) tách biệt hoàn toàn khỏi **minh chứng và lập luận** (`VERIFY.md`, `ADAPT.md`, `GROUP_CHECK.md`). Lựa chọn này phản ánh đúng cách Assessment (mục 8 trong `CLAUDE.md`) chấm điểm theo evidence của từng giai đoạn, không chỉ chấm sản phẩm cuối.

## 3. Yêu cầu hệ thống (Requirements)

- Go ≥ 1.21 (dự án được phát triển và kiểm thử với `go1.27.1`).
- Không có dependency ngoài; toàn bộ chương trình chỉ dùng Go standard library (`bufio`, `fmt`, `os`, `strconv`, `strings`).

Kiểm tra Go đã cài đặt:

```bash
go version
```

## 4. Hướng dẫn sử dụng (Usage)

### 4.1 Build

```bash
cd Lab1
go build -o foodordering.exe .
```

Lệnh này biên dịch `main.go` thành file thực thi `foodordering.exe`. Vì `go.mod` đã khai báo module `foodordering`, không cần thiết lập `GOPATH` thủ công.

### 4.2 Chạy chương trình

```bash
./foodordering.exe
```

Chương trình hiển thị một menu CLI dạng vòng lặp:

```
1. Display available food
2. Search food by name
3. Add food to order
4. Show order summary
5. Exit
Choose an option:
```

**Lựa chọn 1 — Display available food.** In toàn bộ món trong menu có `Available == true`, kèm tên, category và giá. Món `Available == false` (ví dụ "Fried Rice", "Cheesecake" trong dữ liệu mẫu) sẽ không xuất hiện.

**Lựa chọn 2 — Search food by name.** Nhập tên món (chính xác, phân biệt hoa thường). Nếu tìm thấy, in đầy đủ thông tin món kể cả trạng thái `Available`. Nếu không tìm thấy, in `"Food not found"` — không có lỗi runtime nào xảy ra.

**Lựa chọn 3 — Add food to order.** Nhập tên món, sau đó nhập quantity. Hệ thống kiểm tra theo thứ tự: món có tồn tại không → món có còn hàng không → quantity có hợp lệ (> 0) không. Nếu một trong ba điều kiện sai, chương trình in dòng `Error: ...` mô tả rõ nguyên nhân và **không** thêm gì vào order.

**Lựa chọn 4 — Show order summary.** In từng dòng order (tên, quantity, tổng tiền dòng đó), sau đó in `Subtotal`, `Discount` (theo rule ADAPT: 10% nếu subtotal ≥ 300,000 VND), và `Final Total`.

**Lựa chọn 5 — Exit.** Kết thúc chương trình.

### 4.3 Ví dụ phiên làm việc

```
Choose an option: 3
Enter food name: Chicken Rice
Enter quantity: 7
Added to order.

Choose an option: 3
Enter food name: Milk Tea
Enter quantity: 2
Added to order.

Choose an option: 4
Order summary:
- Chicken Rice x7 = 315000 VND
- Milk Tea x2 = 50000 VND
Subtotal: 365000
Discount: 36500
Final Total: 328500
```

### 4.4 Chạy test tự động

```bash
go test -v ./...
```

In chi tiết từng test case cùng kết quả PASS/FAIL. Dùng `-run` để chạy một test cụ thể, ví dụ:

```bash
go test -v -run TestDiscountAtThreshold ./...
```

Toàn bộ kết quả test và quy trình defect injection/fix được ghi chi tiết trong `VERIFY.md`.

## 5. Thiết kế dữ liệu (Data design)

### 5.1 `FoodItem`

```go
type FoodItem struct {
    Name      string
    Price     float64
    Category  string
    Available bool
}
```

Bốn field đúng theo yêu cầu Scenario của đề bài, không mở rộng thêm field nào. `Price` dùng `float64` vì VND trong đề bài được biểu diễn dưới dạng số không có phần thập phân nhưng cần cộng/nhân với các hệ số không nguyên (ví dụ `DiscountRate = 0.10`).

### 5.2 `OrderLine` và `Order`

```go
type OrderLine struct {
    Item     FoodItem
    Quantity int
}

type Order struct {
    Lines []OrderLine
}
```

Quyết định thiết kế quan trọng nhất ở giai đoạn THINK: đề bài gợi ý ví dụ dùng hai slice song song, `[]FoodItem` và `[]int` (quantities), được truyền riêng vào `CalculateSubtotal`. Thiết kế này được đánh giá và **không chọn**, vì hai slice độc lập có thể lệch độ dài nếu không được đồng bộ cẩn thận ở mọi điểm thêm/xóa — một lớp lỗi tiềm ẩn (độ dài không khớp) hoàn toàn có thể tránh được bằng cấu trúc dữ liệu khác.

Thiết kế được chọn gộp `FoodItem` và `Quantity` của nó vào cùng một struct (`OrderLine`), và `Order` chỉ giữ một slice duy nhất `[]OrderLine`. Với cách này, mỗi phần tử của slice luôn tự mang đủ thông tin (món nào, bao nhiêu), nên không tồn tại trạng thái "lệch độ dài" giữa hai danh sách — bất biến (invariant) này được đảm bảo bởi cấu trúc dữ liệu, không phải bởi quy ước lập trình phải tuân thủ thủ công.

### 5.3 Slice và map

Hệ thống dùng cả hai cấu trúc dữ liệu, mỗi cấu trúc phục vụ một mục đích riêng, không trùng lặp vai trò:

| Cấu trúc | Biến | Vai trò |
|---|---|---|
| Slice | `menu []FoodItem` | Nguồn dữ liệu chính (source of truth), giữ thứ tự, dùng khi cần duyệt toàn bộ (hiển thị available food). |
| Map | `menuByName map[string]FoodItem` | Index tra cứu theo tên, dùng khi cần tìm một món cụ thể (search, validate trong add-to-order). Được build một lần từ `menu` khi khởi động. |

Việc tách slice (duyệt tuần tự) và map (tra cứu theo khóa) phản ánh đúng hai pattern truy cập dữ liệu khác nhau trong ứng dụng: hiển thị toàn bộ (O(n), không tránh được) và tìm một phần tử theo tên (O(1) trung bình nhờ map, so với O(n) nếu phải duyệt slice mỗi lần).

## 6. Business logic (Các hàm chính)

| Hàm | Trách nhiệm |
|---|---|
| `buildMenuByName` | Khởi tạo map tra cứu từ slice menu. |
| `DisplayAvailableFood` | Duyệt slice, in các món có `Available == true`. |
| `SearchFoodByName` | Tra map, trả `(FoodItem, bool)` theo idiom Go (ok-pattern), không panic khi không tìm thấy. |
| `IsValidQuantity` | Business rule duy nhất về quantity: chỉ `quantity > 0` là hợp lệ. |
| `AddFoodToOrder` | Hàm tổng hợp mọi validation (not-found, unavailable, invalid quantity) trước khi mutate order; trả `error` mô tả rõ nguyên nhân nếu thất bại. |
| `CalculateSubtotal` | Cộng `Price × Quantity` qua toàn bộ `OrderLine` của order. |
| `CalculateDiscount` | Áp rule ADAPT: 10% nếu subtotal ≥ 300,000 VND, ngược lại 0. |
| `DisplayOrderSummary` | Kết hợp `CalculateSubtotal` và `CalculateDiscount` để in summary đầy đủ (từng dòng, Subtotal, Discount, Final Total). |

Nguyên tắc phân chia trách nhiệm (separation of concerns) được áp dụng nhất quán: validation nằm trong `AddFoodToOrder`/`IsValidQuantity`, tính toán nằm trong `CalculateSubtotal`/`CalculateDiscount`, hiển thị nằm trong `DisplayAvailableFood`/`DisplayOrderSummary`, và điều phối input/output nằm riêng trong `main()`. Nhờ đó, toàn bộ logic nghiệp vụ có thể được kiểm thử bằng `go test` mà không cần giả lập input từ dòng lệnh (xem `main_test.go`).

## 7. Tài liệu liên quan (Related documents)

| Tài liệu | Nội dung |
|---|---|
| [`CLAUDE.md`](./CLAUDE.md) | Kế hoạch chi tiết toàn bộ 8 mục của đề bài: Scenario, THINK, BUILD, VERIFY, ADAPT, GROUP CHECK, Submission, Assessment. |
| [`VERIFY.md`](./VERIFY.md) | Giai đoạn VERIFY: kết quả 7 test case bắt buộc, quy trình defect injection → identify → explain → fix → retest, và AI review của hàm `AddFoodToOrder`. |
| [`ADAPT.md`](./ADAPT.md) | Giai đoạn ADAPT: thay đổi code cho discount rule, 3 test case (dưới/đúng/trên ngưỡng), và reflection về khả năng mở rộng thiết kế. |
| [`GROUP_CHECK.md`](./GROUP_CHECK.md) | Chuẩn bị trả lời 8 possible challenge của giảng viên (giải thích struct, slice/map, hàm, edge case, fix bug, thêm search-by-category, giải thích phần AI-generated, đổi discount rule). |

## 8. Giới hạn phạm vi (Scope notes)

Theo đúng yêu cầu "tuyệt đối không được thêm bớt bất kỳ thứ gì" của đề bài, hệ thống **không** bao gồm các mục thuộc phần Homework (mục 9 trong đề: customer type Regular/Member, order type Pickup/Delivery) vì đó là phần gia hạn một tuần riêng, không thuộc phạm vi buổi lab 150 phút. Hàm `SearchFoodByCategory` được nêu trong `GROUP_CHECK.md` mục 6.6 chỉ là ví dụ chuẩn bị sẵn cho câu hỏi "Add search by category" của giảng viên, và **chưa** được thêm vào `main.go`, vì đó là một possible challenge trong GROUP CHECK, không phải yêu cầu bắt buộc của BUILD.
