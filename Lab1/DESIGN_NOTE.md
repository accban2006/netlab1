# DESIGN_NOTE — Lab 01 Go Foundation

## Tổ chức file

PDF viết toàn bộ vào `main.go`, nhưng khi code lên tới Part 4 thì file đó dài
hơn 600 dòng và mỗi lần sửa một task phải cuộn qua ba task khác. Tôi tách theo
part: `part1_basics.go` … `part4_logic.go`, cộng ba file nền là `property.go`
(struct + method), `format.go` (định dạng tiền) và `data.go` (dữ liệu mẫu). Tất
cả cùng `package main` nên `go run .` vẫn chạy như một chương trình duy nhất —
không đổi gì về cách build, chỉ đổi cách đọc. Mỗi task có một hàm `runTaskXY()`
để menu gọi được và cũng demo độc lập được khi giảng viên muốn xem riêng một
phần.

## Vì sao `[]DistrictStats` chứ không `map`

Task 2.2 yêu cầu "display results sorted by average price (highest first)". Map
trong Go không giữ thứ tự — tệ hơn, `range` trên map được cố tình randomize nên
chạy hai lần cho ra hai thứ tự khác nhau. Vì vậy `analyzeByDistrict` vẫn trả map
(đúng như PDF, vì nó chỉ group dữ liệu), còn `calculateDistrictStats` trả
`[]DistrictStats` đã sort. Chỗ nào buộc phải in từ map thì lấy keys ra,
`sort.Strings(keys)` rồi range theo keys. Tôi cũng thêm tie-break theo tên quận
trong `sort.Slice`, vì `sort.Slice` không stable: hai quận cùng giá trung bình có
thể đổi chỗ nhau giữa các lần chạy, làm output không lặp lại được.

## Value receiver xuyên suốt

Mọi method đều dùng `(p Property)` chứ không `(p *Property)`. `Property` chỉ có
năm field (hai float, một int, hai string) nên copy rẻ, và không method nào cần
mutate — tất cả đều là hàm tính toán đọc-thôi. Trộn hai loại receiver trong cùng
một type là thứ dễ gây nhầm nhất cho người đọc sau, nên tôi chọn nhất quán.

## Xử lý lỗi: chia đôi theo loại hàm

Hàm tính toán trả zero value, hàm đọc input trả `error`. Lý do: `PricePerM2()`
và `CalculateROI()` được gọi trong vòng lặp in bảng, nếu trả `(float64, error)`
thì mỗi lần in phải kiểm lỗi và code in bảng sẽ ngập `if err != nil`. Chia cho 0
ở đây không phải lỗi người dùng mà là dữ liệu bất thường, trả 0 là đủ. Ngược
lại, `readInt`/`readFloat` nhận input người dùng gõ — ở đó lỗi là chuyện bình
thường và cần nói rõ sai cái gì, nên bắt buộc trả `error`.

`CalculateLoan` nằm giữa hai nhóm, nên tôi làm cả hai: `CalculateLoan` trả
`LoanInfo{}` để giữ đúng chữ ký PDF và gọi gọn trong `recommendProperty`, còn
`CalculateLoanChecked` trả `(LoanInfo, error)` cho menu dùng — người dùng nhập
150% trả trước thì phải thấy lý do, không phải thấy một bảng số 0.

## Ba chỗ cố ý lệch khỏi PDF

**1. `fmt.Scanln` → `bufio.Reader`.** Starter code Part 5 dùng
`fmt.Scanln(&choice)` với `choice` là `int`. Nhập `abc` thì `Scanln` fail nhưng
`abc` vẫn nằm trong buffer stdin, nên lần lặp sau đọc lại đúng ký tự đó và fail
tiếp — vòng `for` in menu vô hạn không dừng được. Tôi đọc cả dòng bằng
`bufio.Reader.ReadString('\n')` rồi `strconv.Atoi`, nên buffer luôn sạch sau mỗi
lượt. Thêm nữa, phân biệt EOF với lỗi parse qua `errors.Is(err, io.EOF)`: EOF thì
thoát chương trình (quan trọng khi chạy qua pipe), lỗi parse thì báo rồi hỏi lại.

**2. Saigon Apartment là PREMIUM, không phải STANDARD.** Bảng output mẫu của PDF
ghi `Saigon Apartment: STANDARD`, nhưng giá/m² của nó là 2.5e9 / 75.5 =
33,112,583, mà `categorizeProperty()` — cũng do PDF cung cấp — trả `PREMIUM` cho
mọi giá trị `> 30000000`. Hai thứ trong cùng tài liệu không khớp nhau. Tôi làm
theo hàm vì hàm là đặc tả cụ thể hơn, và ghi comment tại `part1_basics.go` cộng
một test case (`TestCategorizeProperty/saigon_apartment`) để lựa chọn này hiển
hiện chứ không âm thầm.

**3. Copy slice trước khi sort trong `optimizePortfolio`.** Nếu sort trực tiếp
tham số `properties []Property`, slice của caller bị đổi thứ tự ngay — slice là
view lên cùng backing array nên `sort` ghi đè tại chỗ. Hậu quả thực tế: chạy menu
mục 6 rồi quay lại mục 1 sẽ thấy danh sách property đã bị xáo. Tôi zip
(property, ROI) vào `[]propertyROI` mới trong `rankByROI` rồi sort bản copy đó,
và có test `TestOptimizePortfolioDoesNotMutateCaller` canh đúng điều này. Cùng
lý do, `pdfProperties` dùng `append([]Property(nil), properties[:3]...)` chứ
không `properties[:3]` — slice cắt ra vẫn dùng chung backing array.

## Giới hạn đã biết

`optimizePortfolio` dùng greedy theo ROI, đúng như PDF yêu cầu ("greedily add
properties"), nhưng greedy không giải tối ưu bài 0/1 knapsack: nó có thể tiêu
budget vào một property ROI cao rồi không còn chỗ cho tổ hợp tổng lợi nhuận lớn
hơn. Với dữ liệu mẫu kết quả vẫn đúng, nên tôi giữ greedy và ghi comment thay vì
đổi sang DP.

ROI trung bình của danh mục dùng trung bình cộng đơn giản để khớp số `13.3%` của
PDF ((18+12+10)/3). Về tài chính, bình quân gia quyền theo giá mới đúng — danh
mục này thực ra sinh lợi 11.5%, vì property ROI cao nhất lại là property nhỏ
nhất. Tôi in cả hai số và ghi nhãn rõ, để không phải chọn giữa "khớp PDF" và
"đúng về tài chính".

Giả định thuê 1.2%/tháng ở Task 4.1 có một hệ quả đáng nói: vì rent tỷ lệ thuận
với giá, ROI luôn ra 14.4%/năm cho **mọi** property, nên nhánh phân loại theo ROI
không bao giờ phân biệt được gì. Chính bonus (premium location, optimal size) và
warning (giá/m² > 60M) mới tạo khác biệt giữa các khuyến nghị. Tôi hiện thực thang
bậc bằng `[]string` + index có clamp, nâng một bậc khi có ≥ 2 bonus và hạ một bậc
khi có warning — gọn hơn `if` lồng nhau và dễ thêm bậc sau này.
