# Property Calculator — Lab 01 Go Foundation

Công cụ phân tích bất động sản viết bằng Go thuần (chỉ dùng standard library),
gồm: so sánh giá, phân loại, tìm kiếm, thống kê theo quận, tính ROI, tính khoản
vay, khuyến nghị đầu tư, tối ưu danh mục, và một menu tương tác gói tất cả lại.

Chi tiết thiết kế và lý do các quyết định kỹ thuật nằm ở `DESIGN_NOTE.md`.
File này chỉ tập trung vào **cách chạy** và **từng file dùng để làm gì**.

## Yêu cầu

- Go >= 1.21 (máy đang dùng `go1.27.1`).
- Không cần cài thêm gì — module không có dependency ngoài.

## Cách chạy

Từ thư mục `Lab1/`:

```bash
go run .
```

Chương trình in menu tương tác:

```
=== Property Analyzer Menu ===
1. View all properties
2. Search by budget
3. Investment analysis
4. Loan calculator
5. Get recommendations
6. Optimize portfolio
7. District analysis
8. Run all parts (demo Task 1.1 -> 4.2)
0. Exit
```

Nhập số rồi Enter để chọn mục. Các mục 2/4/5/6 sẽ hỏi thêm tham số (budget,
% trả trước, lãi suất, số năm...) — nhấn Enter rỗng để dùng giá trị mặc định
in sẵn trong prompt. Nhập `0` để thoát, hoặc `Ctrl+D` (EOF) để thoát an toàn
bất cứ lúc nào mà không bị treo.

Muốn xem nhanh toàn bộ output của Part 1 → Part 4 mà không cần bấm từng mục,
chọn `8` trong menu.

### Build / kiểm tra tĩnh

```bash
go build ./...     # build toàn bộ, báo lỗi compile
go vet ./...        # lint tĩnh của Go
gofmt -l .           # liệt kê file chưa format đúng chuẩn (không in gì là sạch)
```

### Chạy test

```bash
go test ./...        # chạy toàn bộ test, chỉ in PASS/FAIL
go test -v ./...     # chạy chi tiết từng subtest (dùng để đối chiếu TEST_RESULTS.md)
```

## Cấu trúc file

### Code chính (`package main`, mọi file cộng lại chạy như một chương trình)

| File | Nội dung |
|---|---|
| `go.mod` | Khai báo module `property-calculator`, Go 1.27.1, không có dependency ngoài. |
| `property.go` | Struct `Property`, `LoanInfo`, `DistrictStats` và các method cốt lõi: `PricePerM2`, `IsAffordable`, `CalculateROI`, `InvestmentGrade`, `CalculateLoan`/`CalculateLoanChecked`, hàm `calculateMonthlyPayment` (công thức annuity). |
| `errors.go` | Các lỗi validate dùng chung (`errInvalidDownPayment`, `errInvalidInterestRate`, `errInvalidYears`) cho khoản vay và input. |
| `format.go` | Hàm format số/bảng: `formatPrice` (đơn vị tỷ/triệu/nghìn VND), `formatPriceDetailed` (giữ thêm số thập phân cho số cần chính xác), `formatVND` (dấu phẩy phân cách nghìn), `printHeader`, `printPropertyTable`. |
| `data.go` | Dữ liệu mẫu: 5 property (3 property gốc của đề + 2 thêm), `monthlyRents` tương ứng theo index, `premiumDistricts`, và bộ `pdfProperties`/`pdfRents` (đúng 3 property gốc, dùng cho các test case cần đối chiếu số liệu đề bài). |
| `part1_basics.go` | **Part 1** — Task 1.1 (so sánh giá/m² bằng biến rời, tìm giá thấp nhất) và Task 1.2 (`categorizeProperty`, đếm property theo hạng LUXURY/PREMIUM/STANDARD/BUDGET). |
| `part2_collections.go` | **Part 2** — Task 2.1 (`findPropertiesInBudget`, `findPropertiesByBedrooms`, `findPropertiesByDistrict`) và Task 2.2 (`analyzeByDistrict`, `calculateDistrictStats` — thống kê và sort theo giá trung bình giảm dần). |
| `part3_investment.go` | **Part 3** — Task 3.1 (`findBestInvestment`/`findBestInvestmentROI`) và Task 3.2 (`runLoanAnalysis` — in bảng khoản vay cho mọi property theo tham số). |
| `part4_logic.go` | **Part 4** — Task 4.1 (`recommendProperty` bản gốc, `smartRecommendProperty` có bonus/warning) và Task 4.2 (`optimizePortfolio`/`optimizePortfolioWithRents` — chọn danh mục greedy theo ROI trong giới hạn budget). |
| `main.go` | **Part 5** — menu tương tác dùng `bufio.Reader` (không dùng `fmt.Scanln` vì gây loop vô hạn khi nhập chữ), các hàm đọc input an toàn `readLine`/`readInt`/`readFloat`, và các handler cho từng mục menu. |

### Test (`go test ./...`)

| File | Nội dung |
|---|---|
| `property_test.go` | Test các method trên `Property` và `LoanInfo`: `PricePerM2`, `IsAffordable`, `CalculateROI`, `InvestmentGrade`, `calculateMonthlyPayment`, `CalculateLoan`/`CalculateLoanChecked`. |
| `format_test.go` | Test `formatPrice`, `formatPriceDetailed`, `formatVND` — bảng test theo từng nhánh đơn vị (tỷ/triệu/nghìn/đồng) và biên. |
| `collections_test.go` | Test `findPropertiesInBudget`, `findPropertiesByBedrooms`, `findPropertiesByDistrict`, `analyzeByDistrict`, `calculateDistrictStats` (bao gồm case rỗng và thứ tự sort). |
| `logic_test.go` | Test `recommendProperty`, `smartRecommendProperty`, `optimizePortfolio`/`optimizePortfolioWithRents` — gồm case bắt buộc của đề (budget 8 tỷ chọn cả 3 property) và test slice gốc không bị mutate sau khi sort. |
| `main_test.go` | Test các hàm parse input (`readInt`, `readFloat` qua giả lập reader) và `isEOF`. |

### Deliverable bắt buộc (markdown, không phải code)

| File | Nội dung |
|---|---|
| `TEST_RESULTS.md` | Bảng test case kèm output thật của `go test -v ./...`, cộng bảng test tay cho menu (8 case tương tác) với cột pass/fail. |
| `DESIGN_NOTE.md` | Giải thích lý do thiết kế: vì sao tách file theo part, vì sao `calculateDistrictStats` trả `[]DistrictStats` thay vì map, vì sao dùng value receiver, cách xử lý lỗi, và 3 chỗ cố ý lệch so với PDF gốc. |
| `AI_LOG.md` | Log các lần trao đổi với AI trong lúc làm lab — prompt, AI trả gì, giữ/sửa/bỏ và lý do, cách đã verify lại. |
| `REFLECTION.md` | Bài phản tư theo khung ADAPT (Ask – Discover – Apply – Practice – Teach). |

### Khác

| File | Nội dung |
|---|---|
| `CLAUDE.md` | Kế hoạch triển khai chi tiết ban đầu cho lab này (bám theo `lab1-go-foundation.pdf`) — nguồn tham chiếu khi cần đối chiếu lại yêu cầu đề bài. |
| `lab1-go-foundation.pdf` | Đề bài gốc của lab. |

## Dữ liệu mẫu

5 property trong `data.go` (3 đầu lấy từ đề, 2 sau thêm để Part 1/2 có đủ số liệu):

| Name | Price (VND) | Area (m²) | Bedrooms | District |
|---|---|---|---|---|
| Saigon Apartment | 2,500,000,000 | 75.5 | 2 | District 1 |
| HCMC House | 4,200,000,000 | 120.0 | 3 | District 7 |
| Budget Studio | 800,000,000 | 35.0 | 1 | Binh Thanh |
| Thao Dien Villa | 9,500,000,000 | 145.0 | 4 | District 2 |
| Ben Thanh Penthouse | 6,800,000,000 | 98.0 | 3 | District 1 |

## Lưu ý khi chạy trên Windows

File dùng ký tự UTF-8 (`²`, tiếng Việt có dấu). Nếu terminal in sai ký tự, chạy
`chcp 65001` trước khi `go run .`.
