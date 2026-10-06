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
├── CLAUDE.md                                        # Yêu cầu lab và homework (mục 1-9), có clarification bổ sung
├── AGENTS.md                                        # Hướng dẫn dự án cho AI
├── README.md                                        # Tài liệu này
├── VERIFY.md                                        # Giai đoạn 4 — VERIFY: test case, defect injection, AI review
├── ADAPT.md                                         # Giai đoạn 5 — ADAPT: discount requirement, test, reflection
├── GROUP_CHECK.md                                   # Giai đoạn 6 — GROUP CHECK: chuẩn bị trả lời câu hỏi giảng viên
├── HOMEWORK.md                                      # Mục 9 — HOMEWORK: customer/order type, công thức mới, AI review
├── CODE_QA.md                                       # Câu hỏi và trả lời chi tiết về code hiện tại, logic và thiết kế
├── go.mod                                           # Khai báo Go module (module foodordering)
├── main.go                                          # Toàn bộ mã nguồn hiện thực (struct, slice, map, business logic, CLI)
├── main_test.go                                     # 14 test function, gồm 12 subtest tổ hợp homework
└── IT096IIU_Lab01_Food_Ordering_Student (1).pdf      # Đề bài gốc do giảng viên cung cấp
```

Nguyên tắc tổ chức: **mã nguồn** (`main.go`, `main_test.go`) tách biệt hoàn toàn khỏi **minh chứng và lập luận** (`VERIFY.md`, `ADAPT.md`, `GROUP_CHECK.md`). Lựa chọn này phản ánh đúng cách Assessment (mục 8 trong `CLAUDE.md`) chấm điểm theo evidence của từng giai đoạn, không chỉ chấm sản phẩm cuối.

## 3. Yêu cầu hệ thống (Requirements)

- `go.mod` khai báo `go 1.27.1`; môi trường kiểm tra ngày 06/10/2026 dùng `go1.27.1 windows/amd64`. Dùng toolchain đáp ứng phiên bản khai báo trong module; không coi Go 1.21 là phiên bản tối thiểu của cấu hình hiện tại.
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
5. Set customer type (Regular/Member)
6. Set order type (Pickup/Delivery)
7. Exit
Choose an option:
```

**Lựa chọn 1 — Display available food.** In toàn bộ món trong menu có `Available == true`, kèm tên, category và giá. Món `Available == false` (ví dụ "Fried Rice", "Cheesecake" trong dữ liệu mẫu) sẽ không xuất hiện.

**Lựa chọn 2 — Search food by name.** Nhập tên món (chính xác, phân biệt hoa thường). Nếu tìm thấy, in đầy đủ thông tin món kể cả trạng thái `Available`. Nếu không tìm thấy, in `"Food not found"` — không có lỗi runtime nào xảy ra.

**Lựa chọn 3 — Add food to order.** Nhập tên món, sau đó nhập quantity. Hệ thống kiểm tra theo thứ tự: món có tồn tại không → món có còn hàng không → quantity có hợp lệ (> 0) không. Nếu một trong ba điều kiện sai, chương trình in dòng `Error: ...` mô tả rõ nguyên nhân và **không** thêm gì vào order.

**Lựa chọn 4 — Show order summary.** In từng dòng order (tên, quantity, tổng tiền dòng đó), sau đó in loại khách/loại đơn, `Subtotal`, `Discount` (rule ADAPT: 10% nếu subtotal ≥ 300,000 VND), `Member Discount` (5% nếu là Member), `Delivery Fee`, và `Final Total`.

**Lựa chọn 5 — Set customer type.** Nhập `Regular` hoặc `Member` (chính xác, phân biệt hoa thường). Member được giảm 5% trên food subtotal. Input khác bị từ chối và order không đổi.

**Lựa chọn 6 — Set order type.** Nhập `Pickup` hoặc `Delivery`. Delivery cộng phí 30,000 VND sau khi đã trừ mọi discount. Input khác bị từ chối và order không đổi.

**Lựa chọn 7 — Exit.** Kết thúc chương trình.

Order mới mặc định là `Regular` + `Pickup`, tức không có Member discount và không phụ phí. Volume discount 10% vẫn áp dụng nếu subtotal đạt 300,000 VND. Chi tiết công thức và cách hai discount kết hợp nằm trong [`HOMEWORK.md`](./HOMEWORK.md).

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
Customer: Regular | Order: Pickup
Subtotal: 365000
Discount: 36500
Member Discount: 0
Delivery Fee: 0
Final Total: 328500
```

Cùng order đó nhưng đặt Member + Delivery (lựa chọn 5 và 6) cho kết quả `Final Total: 340250`, tức `365000 − 36500 − 18250 + 30000`.

### 4.4 Chạy test tự động

```bash
go test -count=1 -v ./...
go vet ./...
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

Bốn field đúng theo yêu cầu Scenario của đề bài. `Price` dùng `float64` theo cách biểu diễn trong ví dụ của đề, thuận tiện cho phép nhân với rate 0.10/0.05. Đây là lựa chọn đơn giản cho lab, không bắt buộc để tính phần trăm: hệ thống thực tế có thể dùng `int64` lưu VND và quy định cách làm tròn. Code hiện tại chỉ định dạng `%.0f` khi in, chưa có chính sách làm tròn tiền riêng.

### 5.2 `OrderLine` và `Order`

```go
type OrderLine struct {
    Item     FoodItem
    Quantity int
}

type Order struct {
    Lines    []OrderLine
    Customer CustomerType
    Type     OrderType
}
```

Hai field `Customer` và `Type` được thêm ở phần Homework (mục 9), đặt tại `Order` chứ không tại `OrderLine` vì loại khách và hình thức nhận hàng là thuộc tính của cả đơn hàng — xem [`HOMEWORK.md`](./HOMEWORK.md) mục 9.3.

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
| `CalculateDiscount` | Áp rule ADAPT: volume discount 10% nếu subtotal ≥ 300,000 VND, ngược lại 0. |
| `CalculateCustomerDiscount` | Homework: 5% trên food subtotal nếu khách là Member, ngược lại 0. |
| `CalculateDeliveryFee` | Homework: 30,000 VND nếu order là Delivery, ngược lại 0. |
| `CalculateFinalTotal` | Homework: nơi duy nhất viết ra công thức `Subtotal − VolumeDiscount − MemberDiscount + DeliveryFee`. |
| `DisplayOrderSummary` | Kết hợp các hàm tính toán để in summary đầy đủ (từng dòng, Subtotal, cả hai discount, Delivery Fee, Final Total). |

Nguyên tắc phân chia trách nhiệm (separation of concerns) được áp dụng nhất quán: validation nằm trong `AddFoodToOrder`/`IsValidQuantity`, tính toán nằm trong `CalculateSubtotal`/`CalculateDiscount`, hiển thị nằm trong `DisplayAvailableFood`/`DisplayOrderSummary`, và điều phối input/output nằm riêng trong `main()`. Nhờ đó, toàn bộ logic nghiệp vụ có thể được kiểm thử bằng `go test` mà không cần giả lập input từ dòng lệnh (xem `main_test.go`).

## 7. Tài liệu liên quan (Related documents)

| Tài liệu | Nội dung |
|---|---|
| [`CLAUDE.md`](./CLAUDE.md) | Yêu cầu mục 1-9; bổ sung cách làm rõ việc kết hợp hai discount và ma trận test homework. |
| [`VERIFY.md`](./VERIFY.md) | Giai đoạn VERIFY: kết quả 7 test case bắt buộc, quy trình defect injection → identify → explain → fix → retest, và AI review của hàm `AddFoodToOrder`. |
| [`ADAPT.md`](./ADAPT.md) | Giai đoạn ADAPT: thay đổi code cho discount rule, 3 test case (dưới/đúng/trên ngưỡng), và reflection về khả năng mở rộng thiết kế. |
| [`GROUP_CHECK.md`](./GROUP_CHECK.md) | Chuẩn bị trả lời 8 possible challenge của giảng viên (giải thích struct, slice/map, hàm, edge case, fix bug, thêm search-by-category, giải thích phần AI-generated, đổi discount rule). |
| [`HOMEWORK.md`](./HOMEWORK.md) | Mục 9 — Homework: customer type Regular/Member, order type Pickup/Delivery, quyết định cách hai discount kết hợp, 12 test case tổ hợp, AI review của `CalculateFinalTotal`. |
| [`CODE_QA.md`](./CODE_QA.md) | Câu hỏi bảo vệ bài và trả lời dựa trên code hiện tại, gồm trace tính tiền, con trỏ, slice/map, validation, test và giới hạn thiết kế. |

## 8. Giới hạn phạm vi (Scope notes)

Phần Homework đã được hiện thực và ghi lại trong [`HOMEWORK.md`](./HOMEWORK.md). Các hàm tính subtotal, volume discount và thêm món vẫn tách khỏi các hàm customer discount/delivery fee. `DisplayOrderSummary` và `main` có tích hợp chức năng mới. Lần kiểm tra hiện tại xác nhận 10 test VERIFY/ADAPT vẫn PASS; thư mục không có Git history nên không thể dùng kết quả này để khẳng định lịch sử rằng mọi hàm cũ chưa từng bị sửa.

Hàm `SearchFoodByCategory` được nêu trong `GROUP_CHECK.md` mục 6.6 vẫn **chưa** được thêm vào `main.go`, vì đó chỉ là ví dụ chuẩn bị sẵn cho possible challenge "Add search by category" của giảng viên trong GROUP CHECK, không phải yêu cầu bắt buộc của BUILD hay Homework.

## 9. Đối chiếu homework với dự án hiện tại

### 9.1 Kết luận và điểm khác nhau giữa hai tài liệu yêu cầu

**Dự án đáp ứng các chức năng và minh chứng kỹ thuật của homework theo `CLAUDE.md`, với chính sách cộng dồn discount đã được ghi rõ.** Tuy nhiên, không nên khẳng định công thức hiện tại trùng hoàn toàn với công thức riêng của mục 9 trong PDF.

PDF trang 5-6 yêu cầu Regular/Member, Pickup/Delivery, Member giảm 5% trên food subtotal và công thức `Final Total = Subtotal - Customer Discount + Delivery Fee`. PDF không nói rõ có giữ/cộng dồn discount 10% từ ADAPT trong homework hay không. `CLAUDE.md` mục 9 bổ sung yêu cầu làm rõ cách kết hợp hai cơ chế và kiểm thử đủ 4 tổ hợp × 3 mốc subtotal. Đây là clarification của tài liệu dự án, không phải câu được ghi nguyên văn trong PDF.

Code chọn giữ cả hai discount, cùng tính trên subtotal gốc:

```text
S = tổng Price × Quantity
V = S × 10% nếu S >= 300000; ngược lại 0
C = S × 5% nếu Customer == Member; ngược lại 0
F = 30000 nếu Type == Delivery; ngược lại 0
Final Total = S - V - C + F
```

Regular không có **customer discount**, nhưng vẫn có volume discount ở ngưỡng quy định. Với subtotal 300,000, Member + Delivery trả 285,000 theo code; nếu chỉ áp công thức homework trong PDF và bỏ volume discount thì trả 315,000. Vì vậy cần trình bày chính sách cộng dồn khi nộp/bảo vệ bài; nếu giảng viên yêu cầu homework thay thế ADAPT thì công thức và kỳ vọng test phải đổi tương ứng.

### 9.2 Checklist và bằng chứng

| Yêu cầu | Bằng chứng trong dự án | Đánh giá |
|---|---|---|
| Regular và Member, 5% food subtotal | `CustomerType`, `Order.Customer`, `CalculateCustomerDiscount` | Đạt |
| Pickup 0; Delivery 30,000 VND | `OrderType`, `Order.Type`, `CalculateDeliveryFee` | Đạt |
| Member không giảm delivery fee | Fee cộng ngoài hai discount; `TestMemberDiscountDoesNotReduceDeliveryFee` | Đạt |
| Nêu cách kết hợp ADAPT và Member | Comment trong `main.go`, `HOMEWORK.md` mục 9.2 và công thức trên | Đạt theo clarification của `CLAUDE.md` |
| Summary hiển thị discount, fee, final total | `DisplayOrderSummary`: `Discount`, `Member Discount`, `Delivery Fee`, `Final Total` | Đạt; `Discount` là volume discount |
| 4 tổ hợp × dưới/đúng/trên 300,000 | 12 subtest trong `TestFinalTotalCombinations` | Đạt |
| Source code cập nhật | `main.go`, `go.mod` | Có |
| Test cases/results | `main_test.go`, `HOMEWORK.md` mục 9.4 và lần chạy bên dưới | Có |
| Short design note | `HOMEWORK.md` mục 9.2-9.3; README mục 5-6 | Có |
| Một AI review | `HOMEWORK.md` mục 9.5; review hiện tại tại mục 9.4 bên dưới | Có nội dung review |
| Giải thích requirement change | `HOMEWORK.md` mục 9.6 | Có |
| Cả hai thành viên hiểu solution | Cần mỗi người tự giải thích và thực hiện challenge | Không xác nhận được chỉ bằng đọc file/chạy test |

### 9.3 Kết quả kiểm tra ngày 06/10/2026

Chạy trên source hiện tại bằng Go 1.27.1, không dùng kết quả cache:

| Kiểm tra | Kết quả |
|---|---|
| `go test -count=1 -v ./...` | PASS: 14 test function cấp cao nhất, gồm 12 subtest tổ hợp |
| `go vet ./...` | PASS, không có diagnostic |
| Build source vào file thực thi tạm | PASS |
| CLI: Chicken Rice ×7 + Milk Tea ×2, Member + Delivery | Subtotal 365,000; volume 36,500; Member 18,250; fee 30,000; final 340,250 |
| Thử lỗi `>=` thành `>` trên bản sao source tạm | `TestDiscountAtThreshold` và cả 4 subtest tại ngưỡng FAIL đúng như dự đoán; source chính không bị sửa |

Ma trận Final Total được test (VND):

| Subtotal | Regular + Pickup | Regular + Delivery | Member + Pickup | Member + Delivery |
|---|---:|---:|---:|---:|
| 200,000 | 200,000 | 230,000 | 190,000 | 220,000 |
| 300,000 | 270,000 | 300,000 | 255,000 | 285,000 |
| 400,000 | 360,000 | 390,000 | 340,000 | 370,000 |

Lưu ý cách đếm: `HOMEWORK.md` ghi “25/25 test” theo cách gộp 13 test khác và 12 case con. Output Go thực tế có **14 test function + 12 subtest**, trong đó một function là cha của 12 subtest; nên dùng cách đếm rõ ràng này.

### 9.4 AI review hiện tại: `CalculateFinalTotal` (không rewrite)

Hàm tính subtotal một lần, dùng cùng subtotal gốc cho hai discount và cộng delivery fee riêng. Theo chính sách cộng dồn của dự án, logic đúng và 12 case tổ hợp kiểm tra được các trường hợp tương ứng. Tách các hàm nhỏ giúp xác định lỗi nằm ở subtotal, discount hay fee mà không cần chạy CLI.

Các giới hạn cần biết khi giải thích code:

- Named string type tránh nhầm giữa hai loại đã định kiểu, nhưng không giới hạn giá trị như enum đóng: `CustomerType("VIP")` vẫn hợp lệ về kiểu. CLI từ chối input sai; hàm tính toán coi loại khách không phải Member là discount 0 và loại đơn không phải Delivery là fee 0.
- Các field được export nên người gọi có thể tạo `OrderLine` với quantity âm hoặc giá âm, bỏ qua `AddFoodToOrder`. Các hàm tính tiền dựa trên giả định order đã được tạo hợp lệ.
- `ReadString` trong CLI đang bỏ qua lỗi I/O. Khi stdin kết thúc mà không chọn 7, vòng lặp có thể in `Invalid option` liên tục. Luồng tương tác bình thường có Exit; xử lý EOF là điểm nên cải thiện.
- Map lưu bản sao `FoodItem` và chỉ được build một lần. Nếu thay menu trong runtime thì cần đồng bộ map; code hiện tại dùng menu tĩnh.
- Summary và CLI chưa có automated test; phiên CLI trên chỉ xác nhận một trường hợp. `float64` và `%.0f` cũng chưa thay thế được một chính sách làm tròn tiền rõ ràng.

Đây là các giới hạn/cải tiến, không phải chức năng homework bị thiếu trong luồng sử dụng hợp lệ. Review này chỉ cập nhật tài liệu, không thay business logic.

### 9.5 Cách chuẩn bị nộp bài

Nộp `main.go`, `go.mod`, `main_test.go`, `HOMEWORK.md` và README để có source, test, design note, AI review và giải thích thay đổi. `CODE_QA.md` hỗ trợ luyện giải thích; không thay thế việc cả hai thành viên tự hiểu bài. `VERIFY.md`, `ADAPT.md`, `GROUP_CHECK.md` chủ yếu ghi giai đoạn lab trước homework; một số đoạn mô tả `Order`/CLI/summary ở các file đó là phiên bản cũ. `README.docx` cũng chưa được đồng bộ với README Markdown trong lần cập nhật này.
