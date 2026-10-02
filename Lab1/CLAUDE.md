# Lab 01 — Go Foundation: Property Calculator

Kế hoạch triển khai chi tiết cho `lab1-go-foundation.pdf` (2.5 giờ, 100 điểm + 15 bonus).

Mục tiêu: xây một công cụ phân tích bất động sản bằng Go — tính giá, lọc/tìm kiếm,
thống kê theo quận, tính ROI, tính khoản vay, đưa ra khuyến nghị đầu tư và tối ưu
danh mục, cuối cùng gói tất cả vào một menu tương tác.

## Môi trường

- Go đã cài: `go1.27.1 windows/amd64` (PDF yêu cầu >= 1.21, OK).
- Thư mục làm việc: `D:/appdata/netcentricstuff/Lab1`.
- Khởi tạo module: `go mod init property-calculator`.
- Không dùng thư viện ngoài. Chỉ `fmt`, `math`, `sort`, `strings`, `os`, `bufio`, `strconv`.

## Cấu trúc file dự định

PDF viết toàn bộ vào `main.go`, nhưng tách file cho dễ đọc và dễ chấm. Tất cả đều
`package main` nên `go run .` vẫn chạy như một chương trình duy nhất.

```
Lab1/
├── go.mod
├── CLAUDE.md              # file này
├── property.go            # struct Property + các method (PricePerM2, IsAffordable, ROI, Loan...)
├── errors.go              # các lỗi validate dùng chung (thêm khi code, không có trong plan ban đầu)
├── format.go              # formatPrice, formatPriceDetailed, formatVND, helper in bảng
├── data.go                # dữ liệu mẫu: properties + monthlyRents
├── part1_basics.go        # Task 1.1, 1.2
├── part2_collections.go   # Task 2.1, 2.2
├── part3_investment.go    # Task 3.1, 3.2
├── part4_logic.go         # Task 4.1, 4.2
├── main.go                # Part 5: menu + đọc input an toàn
├── *_test.go              # test cho các hàm tính toán (bonus)
│
├── TEST_RESULTS.md        # Deliverable 1: bảng test case + output thật
├── DESIGN_NOTE.md         # Deliverable 2: ghi chú thiết kế ngắn
├── AI_LOG.md              # Deliverable 3: log làm việc với AI
└── REFLECTION.md          # Deliverable 4: ADAPT reflection
```

4 file markdown cuối là **deliverable bắt buộc** — xem mục
"4 deliverable bắt buộc" ở cuối tài liệu.

Quy ước: mỗi task có một hàm `runTaskXY()` để menu gọi được và cũng chạy độc lập
được khi demo từng phần.

## Mô hình dữ liệu

```go
type Property struct {
    Name     string
    Price    float64 // VND
    Area     float64 // m2
    Bedrooms int
    District string
}

type LoanInfo struct {
    LoanAmount     float64
    MonthlyPayment float64
    TotalInterest  float64
}

type DistrictStats struct {
    District      string
    Count         int
    AveragePrice  float64
    MostExpensive Property
}
```

Dữ liệu mẫu (`data.go`) — giữ đúng 3 property của PDF rồi thêm 2-3 cái nữa để
Task 1.1/1.2 có đủ số liệu so sánh:

| Name | Price (VND) | Area | Bedrooms | District |
|---|---|---|---|---|
| Saigon Apartment | 2,500,000,000 | 75.5 | 2 | District 1 |
| HCMC House | 4,200,000,000 | 120.0 | 3 | District 7 |
| Budget Studio | 800,000,000 | 35.0 | 1 | Binh Thanh |

`monthlyRents := []float64{25000000, 35000000, 12000000}` — cùng thứ tự với
`properties`, index-based như PDF yêu cầu ở Task 3.1.

## Part 1 — Variables & Basic Types (25 phút, 20 điểm)

### Task 1.1 — Property comparison (10 điểm)

Yêu cầu: 3 property, tính giá/m² cho cả 3, tìm cái thấp nhất theo giá/m².

Cách làm:
- Phần này PDF cố tình bắt dùng biến rời (`property2Name`, `property2Price`...)
  để luyện khai báo biến. Giữ đúng tinh thần đó trong `part1_basics.go`: dùng
  `var` tường minh với kiểu, và `:=` cho giá trị suy ra được như `pricePerM2`.
- So sánh min bằng chuỗi `if`, không dùng slice/sort (slice là nội dung Part 2).
- Khởi tạo `cheapestPrice`/`cheapestName` từ property 1 rồi so dần.

Lưu ý: chia cho `area` — nếu `area == 0` thì chia cho 0 ra `+Inf`, không panic.
Vẫn nên guard để output không ra `+Inf`.

Output mẫu:
```
=== Property Comparison ===
Property 1: Saigon Apartment - 33,112,582 VND/m²
Property 2: HCMC House - 35,000,000 VND/m²
Property 3: Budget Studio - 22,857,142 VND/m²

Cheapest per m²: Budget Studio at 22,857,142 VND/m²
```

(Số `33,113` trong PDF là bản in rút gọn; giá trị đúng của 2.5e9 / 75.5 là
≈ 33,112,582. Dùng số thật và in có dấu phẩy phân cách nghìn.)

### Task 1.2 — Price categories + formatPrice (10 điểm)

Yêu cầu: dùng `categorizeProperty()`, đếm số property mỗi hạng, viết `formatPrice()`.

`categorizeProperty(pricePerM2 float64) string` — copy nguyên từ PDF:
- `> 50,000,000` → `LUXURY`
- `> 30,000,000` → `PREMIUM`
- `> 20,000,000` → `STANDARD`
- còn lại → `BUDGET`

`formatPrice(vnd float64) string` — quy đổi sang đơn vị tiếng Việt:
- `>= 1e9` → `"%.1f tỷ VND"` (2,500,000,000 → `2.5 tỷ VND`)
- `>= 1e6` → `"%.0f triệu VND"` (850,000,000 → `850 triệu VND`)
- `>= 1e3` → `"%.0f nghìn VND"`
- nhỏ hơn → `"%.0f VND"`

Đếm theo hạng: dùng `map[string]int`, nhưng in theo thứ tự cố định
`[]string{"LUXURY","PREMIUM","STANDARD","BUDGET"}` chứ không range map (range map
trong Go có thứ tự ngẫu nhiên — đây là bẫy hay mất điểm).

Output mẫu:
```
=== Property Categories ===
Saigon Apartment: PREMIUM (2.5 tỷ VND)
HCMC House: PREMIUM (4.2 tỷ VND)
Budget Studio: STANDARD (800 triệu VND)

Category Summary:
LUXURY: 0 properties
PREMIUM: 2 properties
STANDARD: 1 properties
BUDGET: 0 properties
```

PDF ghi Saigon Apartment là `STANDARD`, nhưng 33.1M/m² > 30M nên theo đúng hàm
của PDF phải là `PREMIUM`. Làm theo hàm, không theo bảng output — và ghi một
comment ngắn giải thích chênh lệch này.

## Part 2 — Arrays, Slices, Maps (30 phút, 25 điểm)

Chuyển từ biến rời sang `[]Property`. Đây là lúc `property.go` và `data.go` ra đời.

### Task 2.1 — Property search (10 điểm)

```go
func findPropertiesInBudget(properties []Property, maxBudget float64) []Property
func findPropertiesByBedrooms(properties []Property, bedrooms int) []Property
```

- Pattern: khai báo `var result []Property` (nil slice), `append` trong loop, return.
- Trả nil khi không khớp là hợp lệ — `len(nil slice) == 0` và `range` trên nil an toàn.
  Khi in thì kiểm tra `len(result) == 0` để hiện "No properties found" thay vì bảng rỗng.
- Dùng `for _, prop := range properties` — bỏ index vì không cần.
- Test với `budget := 3000000000.0` → kỳ vọng Saigon Apartment + Budget Studio.
- Thêm `findPropertiesByDistrict()` cho menu search phong phú hơn (tính vào bonus).

### Task 2.2 — District analysis with maps (15 điểm)

```go
func analyzeByDistrict(properties []Property) map[string][]Property
func calculateDistrictStats(properties []Property) []DistrictStats
```

`analyzeByDistrict` copy từ PDF — group theo `prop.District`.

`calculateDistrictStats` trả về **slice** chứ không phải map, vì yêu cầu là
"display results sorted by average price (highest first)" — map không giữ thứ tự.
Với mỗi quận tính: `Count`, `AveragePrice` (tổng giá / count), `MostExpensive`.

Sắp xếp: `sort.Slice(stats, func(i, j int) bool { return stats[i].AveragePrice > stats[j].AveragePrice })`.
Thêm tie-break bằng tên quận để output ổn định (`sort.Slice` không stable; dùng
`sort.SliceStable` hoặc so sánh thêm `District` khi giá bằng nhau).

Output mẫu:
```
=== District Analysis ===
District 1: 1 properties, Avg: 2.5 tỷ VND, Most expensive: Saigon Apartment
District 7: 1 properties, Avg: 4.2 tỷ VND, Most expensive: HCMC House
Binh Thanh: 1 properties, Avg: 800 triệu VND, Most expensive: Budget Studio

Ranking by Average Price:
1. District 7: 4.2 tỷ VND
2. District 1: 2.5 tỷ VND
3. Binh Thanh: 800 triệu VND
```

Phần "=== District Analysis ===" muốn in theo thứ tự ổn định thì lấy keys của map,
`sort.Strings(keys)`, rồi range theo keys.

## Part 3 — Functions & Methods (35 phút, 25 điểm)

### Method cơ bản (từ PDF, vào `property.go`)

```go
func (p Property) PricePerM2() float64      // guard Area == 0 → return 0
func (p Property) IsAffordable(budget float64) bool
```

Dùng value receiver `(p Property)` xuyên suốt — struct nhỏ, không cần mutate.
Nhất quán receiver là một tiêu chí của mục Code Quality.

### Task 3.1 — Investment calculator (10 điểm)

```go
func (p Property) CalculateROI(monthlyRent float64) float64  // (rent*12/Price)*100
func (p Property) InvestmentGrade(roi float64) string
func findBestInvestment(properties []Property, rents []float64) Property
```

- `CalculateROI`: guard `p.Price == 0` → return 0.
- `InvestmentGrade`: `> 8` EXCELLENT, `>= 5` GOOD, `>= 3` FAIR, còn lại POOR.
  PDF viết dạng khoảng "5-8%" nên chọn biên: `>8` / `5..8` / `3..5` / `<3`.
  Chốt `>=` ở biên dưới và ghi comment về quy ước biên.
- `findBestInvestment`: ghép `properties[i]` với `rents[i]` theo index. **Phải**
  guard `len(rents) != len(properties)` — truy cập lệch index là panic
  `index out of range`. Return `Property{}` + xử lý khi slice rỗng.
  Chữ ký PDF trả một `Property`; thêm bản `findBestInvestmentROI` trả `(Property, float64)`
  để in được ROI kèm theo mà không phải tính lại.

Kiểm chứng số: 25e6×12/2.5e9×100 = 12.0%; 35e6×12/4.2e9×100 = 10.0%;
12e6×12/0.8e9×100 = 18.0%. Khớp output PDF.

```
=== Investment Analysis ===
Saigon Apartment: ROI 12.0% per year - EXCELLENT
HCMC House: ROI 10.0% per year - EXCELLENT
Budget Studio: ROI 18.0% per year - EXCELLENT

Best Investment: Budget Studio (18.0% ROI)
```

### Task 3.2 — Loan calculator (15 điểm)

```go
func calculateMonthlyPayment(loanAmount, annualRate float64, years int) float64
func (p Property) CalculateLoan(downPaymentPercent, interestRate float64, years int) LoanInfo
```

`calculateMonthlyPayment` copy từ PDF (công thức annuity):

```
M = L × r × (1+r)^n / ((1+r)^n − 1),  r = annualRate/100/12,  n = years×12
```

Giữ nhánh `annualRate == 0 → loanAmount / numPayments` để tránh chia 0.
Cần `import "math"` cho `math.Pow`.

`CalculateLoan`:
```go
loanAmount := p.Price * (1 - downPaymentPercent/100)
monthly := calculateMonthlyPayment(loanAmount, interestRate, years)
totalInterest := monthly*float64(years*12) - loanAmount
```
Validate: `downPaymentPercent` trong [0,100], `years > 0`, `interestRate >= 0`.
Nếu sai thì trả `LoanInfo{}` zero — hoặc thêm biến thể trả `error` cho mục
"proper error handling" (3 điểm Code Quality).

Test case bắt buộc: 20% down, 8.5%, 20 năm.
- Saigon Apartment: loan 2.0 tỷ → monthly ≈ 17.36 triệu, tổng lãi ≈ 2.166 tỷ
- HCMC House: loan 3.36 tỷ → monthly ≈ 29.16 triệu, tổng lãi ≈ 3.64 tỷ

Khớp xấp xỉ output PDF (17.3 / 29.1 triệu). Số in ra lấy từ công thức, không
hard-code theo bảng.

```
=== Loan Analysis ===
Saigon Apartment:
   Loan Amount: 2.0 tỷ VND (80% of price)
   Monthly Payment: 17.4 triệu VND
   Total Interest: 2.17 tỷ VND over 20 years
```

## Part 4 — Control Flow & Logic (25 phút, 20 điểm)

### Task 4.1 — Smart recommendation (10 điểm)

```go
func recommendProperty(p Property, budget, maxMonthlyPayment float64) string
func smartRecommendProperty(p Property, budget, maxMonthlyPayment float64) (string, string)
```

`recommendProperty` lấy nguyên từ PDF (base case). `smartRecommendProperty` mở rộng,
trả `(recommendation, details)`:

1. Affordability: `!p.IsAffordable(budget)` → `"SKIP - Over budget"`, return sớm.
2. Loan check: `CalculateLoan(20, 8.5, 20).MonthlyPayment > maxMonthlyPayment`
   → `"CONSIDER - High monthly payment"`.
3. Bonus `[]string`:
   - Premium location: `District` thuộc {District 1, District 2, District 7}
     → `"Premium location"`. Dùng một `map[string]bool premiumDistricts` ở package
     level, so sánh sau `strings.TrimSpace`.
   - Optimal size: `Area >= 50 && Area <= 100` → `"Optimal size"`.
4. Warnings `[]string`:
   - `PricePerM2() > 60_000_000` → `"High price per m²"`.
5. Logic chốt khuyến nghị: tính ROI với giả định thuê 1.2%/tháng
   (`p.CalculateROI(p.Price * 0.012)` — theo PDF), rồi:
   - `roi > 10` → `BUY NOW`; `roi > 6` → `GOOD BUY`; còn lại → `MAYBE`
   - nâng một bậc nếu `len(bonus) >= 2`, hạ một bậc nếu `len(warnings) > 0`
   - Viết thang bậc thành `[]string` + chỉ số để nâng/hạ bằng clamp index, gọn hơn
     if lồng nhau.
6. `details := fmt.Sprintf("Bonus: %v, Warnings: %v", bonus, warnings)`.

Chú ý: ROI 1.2%/tháng = 14.4%/năm cho **mọi** property (vì rent tỷ lệ theo giá),
nên nhánh ROI luôn ra `BUY NOW` nếu chỉ dựa vào giả định này. Đó là lý do bonus/
warnings mới là phần tạo khác biệt. Ghi comment nêu rõ quan sát này.

### Task 4.2 — Portfolio optimizer (10 điểm)

```go
func optimizePortfolio(properties []Property, totalBudget float64) []Property
```

- Copy slice trước khi sort: `sorted := make([]Property, len(properties)); copy(sorted, properties)`.
  Sort trực tiếp tham số sẽ mutate slice gốc của caller — bug kinh điển với slice
  trong Go, và sẽ làm sai output các menu chạy sau.
- Sort giảm dần theo ROI. ROI cần rent → dùng cùng giả định 1.2%/tháng để ROI chỉ
  phụ thuộc property, hoặc nhận thêm `rents []float64` và sort cặp (property, rent).
  Chọn cách nhận thêm rents và zip thành struct phụ `type propertyROI struct{...}`
  để khỏi lệch index sau sort — chữ ký PDF giữ lại dưới dạng wrapper.
- Greedy: `if prop.Price <= remainingBudget { portfolio = append(...); remainingBudget -= prop.Price }`.
- Tính và in: tổng đã đầu tư, budget còn lại, ROI trung bình của danh mục.
  ROI trung bình nên là bình quân gia quyền theo giá; PDF in `13.3%` = trung bình
  cộng đơn giản của (18+12+10)/3 = 13.33. Dùng trung bình cộng để khớp PDF, và
  ghi comment rằng gia quyền chính xác hơn về tài chính.

Test: budget 8 tỷ → chọn cả 3 (tổng 7.5 tỷ, còn 500 triệu).

```
=== Portfolio Optimization ===
Budget: 8.0 tỷ VND
Selected Properties:
1. Budget Studio: 800 triệu VND (ROI: 18.0%)
2. Saigon Apartment: 2.5 tỷ VND (ROI: 12.0%)
3. HCMC House: 4.2 tỷ VND (ROI: 10.0%)

Total Invested: 7.5 tỷ VND
Remaining Budget: 500 triệu VND
Portfolio Average ROI: 13.3%
```

Lưu ý thuật toán: greedy theo ROI không tối ưu tuyệt đối (đây là bài knapsack).
Với dữ liệu mẫu nó cho kết quả đúng. Ghi comment về giới hạn này — đúng với yêu
cầu "Greedily add properties" của PDF.

## Part 5 — Menu tổng hợp (10 phút, +5 bonus)

`main.go`: vòng `for` vô hạn, in menu, đọc lựa chọn, `switch`.

7 mục: 1 view all, 2 search by budget, 3 investment analysis, 4 loan calculator,
5 recommendations, 6 optimize portfolio, 0 exit.

**Không dùng `fmt.Scanln(&choice)`** như starter code của PDF. `Scanln` vào `*int`
khi người dùng nhập chữ sẽ fail, để nguyên ký tự trong buffer và vòng `for` quay
vô hạn in menu liên tục. Thay bằng:

```go
reader := bufio.NewReader(os.Stdin)

func readLine(reader *bufio.Reader, prompt string) (string, error) {
    fmt.Print(prompt)
    line, err := reader.ReadString('\n')
    if err != nil {
        return "", err          // EOF: thoát sạch, không loop vô hạn
    }
    return strings.TrimSpace(line), nil
}

func readInt(...)   (int, error)      // strconv.Atoi
func readFloat(...) (float64, error)  // strconv.ParseFloat
```

Input sai → in lỗi, `continue`. Gặp EOF (`err != nil`) → return, tránh treo khi
chạy qua pipe. Đây là phần "+5 user input validation" của bonus.

Mục 2, 4, 5, 6 nên hỏi tham số (budget / % trả trước / lãi suất / số năm) với giá
trị mặc định khi người dùng Enter rỗng, thay vì hard-code.

## Lộ trình thực thi

1. `go mod init property-calculator` → tạo `go.mod`.
2. `property.go` + `format.go` + `data.go` (nền tảng cho mọi phần sau).
3. `part1_basics.go` → Task 1.1, 1.2 → `go run .` kiểm tra output.
4. `part2_collections.go` → Task 2.1, 2.2.
5. `part3_investment.go` → Task 3.1, 3.2 (import `math`).
6. `part4_logic.go` → Task 4.1, 4.2 (import `sort`, `strings`).
7. `main.go` → menu + input validation.
8. `*_test.go` → test cho `formatPrice`, `CalculateROI`, `calculateMonthlyPayment`,
   `CalculateLoan`, `optimizePortfolio`.

Sau mỗi bước: `go build ./...` rồi `go vet ./...`. Trước khi nộp: `gofmt -l .`
(không ra file nào) và `go test ./...`.

## Bẫy Go cần tránh

- Range trên map có thứ tự ngẫu nhiên → luôn sort keys trước khi in.
- Chia số nguyên: `years * 12` là `int`, phải `float64(years * 12)` khi dùng trong
  công thức float.
- `sort.Slice` không stable → dùng `sort.SliceStable` hoặc thêm tie-break.
- Sort tham số slice làm mutate slice của caller → copy trước khi sort.
- `%d` với `float64` in ra `%!d(float64=...)` → `float64` luôn dùng `%.Nf`.
- `append` vào nil slice là hợp lệ, không cần khởi tạo trước.
- `p.Area == 0` / `p.Price == 0` → guard trước mọi phép chia.
- Zip hai slice theo index → kiểm tra độ dài bằng nhau trước khi loop.
- Ký tự `²` và tiếng Việt có dấu: file lưu UTF-8. Terminal Windows cũ (cp437) có
  thể in sai — nếu gặp, chạy `chcp 65001` hoặc đổi sang `m2` trong chuỗi output.

## Đối chiếu thang điểm

| Hạng mục | Điểm | File | Trạng thái |
|---|---|---|---|
| Task 1.1 Property comparison | 10 | `part1_basics.go` | ☑ |
| Task 1.2 Categorization + formatPrice | 10 | `part1_basics.go`, `format.go` | ☑ |
| Task 2.1 Search functions | 10 | `part2_collections.go` | ☑ |
| Task 2.2 District analysis (maps) | 15 | `part2_collections.go` | ☑ |
| Task 3.1 Investment methods | 10 | `property.go`, `part3_investment.go` | ☑ |
| Task 3.2 Loan system | 15 | `property.go`, `part3_investment.go` | ☑ |
| Task 4.1 Smart recommendation | 10 | `part4_logic.go` | ☑ |
| Task 4.2 Portfolio optimizer | 10 | `part4_logic.go` | ☑ |
| Clean readable code | 5 | toàn bộ | ☑ |
| Error handling | 3 | guards + `readInt`/`readFloat` + `errors.go` | ☑ |
| Naming + comments | 2 | toàn bộ | ☑ |
| Bonus: menu | +5 | `main.go` (9 mục, thêm 7 district + 8 run-all) | ☑ |
| Bonus: input validation | +5 | `main.go` | ☑ |
| Bonus: creative (tests, thêm filter, bảng đẹp) | +5 | `*_test.go` (99 subtest) | ☑ |

Tổng: 100 + 15 bonus.

Deliverable bắt buộc (ngoài thang điểm code):

| Deliverable | File | Trạng thái |
|---|---|---|
| 1. Test cases / results | `TEST_RESULTS.md` | ☑ |
| 2. Short design note | `DESIGN_NOTE.md` | ☑ |
| 3. AI collaboration log | `AI_LOG.md` | ☑ (7 entry) |
| 4. Short ADAPT reflection | `REFLECTION.md` | ☑ |

## 4 deliverable bắt buộc

Ngoài code, phải nộp 4 thứ dưới đây. Viết **sau khi code chạy được**, nội dung lấy
từ output thật — không viết trước rồi đoán số.

### 1. Test cases / results → `TEST_RESULTS.md`

In ra bảng test case kèm kết quả thật. Hai tầng:

**A. Unit test tự động** (`go test -v ./...`) — paste nguyên output vào file.
Các case tối thiểu:

| # | Hàm | Input | Expected | Lý do chọn case |
|---|---|---|---|---|
| 1 | `formatPrice` | `2500000000` | `2.5 tỷ VND` | nhánh tỷ |
| 2 | `formatPrice` | `850000000` | `850 triệu VND` | nhánh triệu |
| 3 | `formatPrice` | `0` | `0 VND` | biên zero |
| 4 | `categorizeProperty` | `50000000` | `PREMIUM` | biên `>` không phải `>=` |
| 5 | `categorizeProperty` | `50000001` | `LUXURY` | vượt biên 1 đơn vị |
| 6 | `PricePerM2` | `Area: 0` | `0` | guard chia 0 |
| 7 | `CalculateROI` | rent `25e6`, price `2.5e9` | `12.0` | khớp PDF |
| 8 | `CalculateROI` | `Price: 0` | `0` | guard chia 0 |
| 9 | `InvestmentGrade` | `8.0` | `GOOD` | biên EXCELLENT/GOOD |
| 10 | `calculateMonthlyPayment` | `2e9, 8.5, 20` | `≈17356465` | so sánh có epsilon |
| 11 | `calculateMonthlyPayment` | `annualRate: 0` | `L/n` | nhánh lãi 0% |
| 12 | `CalculateLoan` | `20, 8.5, 20` | loan `2e9`, lãi `≈2.166e9` | case bắt buộc của PDF |
| 13 | `CalculateLoan` | `downPayment: 150` | `LoanInfo{}` | input không hợp lệ |
| 14 | `findPropertiesInBudget` | `3e9` | 2 kết quả | lọc đúng |
| 15 | `findPropertiesInBudget` | `1e6` | `len == 0` | không khớp gì |
| 16 | `findBestInvestment` | rents lệch độ dài | không panic | guard index |
| 17 | `findBestInvestment` | slice rỗng | `Property{}` | biên rỗng |
| 18 | `calculateDistrictStats` | 3 property | sort giảm dần | thứ tự đúng |
| 19 | `optimizePortfolio` | `8e9` | chọn 3, còn `500e6` | case bắt buộc của PDF |
| 20 | `optimizePortfolio` | `0` | `len == 0` | budget zero |
| 21 | `optimizePortfolio` | — | slice gốc không đổi | không mutate caller |

So sánh `float64` **phải** dùng epsilon (`math.Abs(got-want) < 1e-6`), không dùng
`==`. Dùng table-driven test — idiom chuẩn của Go:

```go
func TestFormatPrice(t *testing.T) {
    tests := []struct{ name string; in float64; want string }{
        {"ty", 2500000000, "2.5 tỷ VND"},
        {"trieu", 850000000, "850 triệu VND"},
        {"zero", 0, "0 VND"},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            if got := formatPrice(tc.in); got != tc.want {
                t.Errorf("formatPrice(%v) = %q, want %q", tc.in, got, tc.want)
            }
        })
    }
}
```

**B. Manual test cho menu** — bảng các lần chạy tay, vì input tương tác không test
tự động được:

| # | Thao tác | Kỳ vọng | Thực tế |
|---|---|---|---|
| M1 | chọn `1` | in 3 property | |
| M2 | chọn `2`, budget `3000000000` | 2 kết quả | |
| M3 | nhập `abc` | báo lỗi, hiện lại menu | |
| M4 | nhập `99` | `Invalid option!` | |
| M5 | Enter rỗng ở prompt budget | dùng giá trị mặc định | |
| M6 | nhập `-5` cho số năm | báo lỗi, không crash | |
| M7 | `Ctrl+D` / EOF | thoát sạch, không loop vô hạn | |
| M8 | chọn `0` | `Goodbye!` rồi exit | |

File phải ghi rõ: tổng số case, số pass/fail, và **case nào fail thì nói thẳng**
kèm output lỗi. Không che fail.

Cuối file paste output của: `go build ./...`, `go vet ./...`, `gofmt -l .`,
`go test -v ./...`.

### 2. Short design note → `DESIGN_NOTE.md`

Ngắn — khoảng 1 trang, 300-500 từ. Giải thích **vì sao** thiết kế như vậy, không
mô tả lại code. Các mục:

- **Tổ chức file**: vì sao tách theo part thay vì dồn hết vào `main.go` như PDF.
- **Vì sao `[]DistrictStats` chứ không `map`** cho Task 2.2: map không giữ thứ tự,
  mà yêu cầu là sort theo giá trung bình.
- **Value receiver `(p Property)`** xuyên suốt: struct nhỏ, không mutate, nhất quán.
- **Xử lý lỗi**: chọn guard + trả zero value cho hàm tính toán, trả `error` cho
  hàm đọc input. Nêu lý do phân đôi như vậy.
- **3 chỗ cố ý lệch khỏi PDF** (quan trọng nhất, giảng viên sẽ hỏi):
  1. `Scanln` → `bufio.Reader` (starter code của PDF loop vô hạn khi nhập chữ).
  2. Saigon Apartment ra `PREMIUM` không phải `STANDARD` (33.1M/m² > 30M — bảng
     output của PDF không khớp hàm của chính PDF).
  3. Copy slice trước khi sort trong `optimizePortfolio` (tránh mutate slice caller).
- **Giới hạn đã biết**: greedy theo ROI không giải tối ưu knapsack; ROI trung bình
  dùng trung bình cộng (khớp PDF) chứ không gia quyền theo giá.

Viết thành đoạn văn ngắn có tiêu đề, không phải bullet rời rạc.

### 3. AI collaboration log → `AI_LOG.md`

PDF mục "Getting Help" khuyến khích dùng AI. Log lại trung thực việc đó.

Mỗi entry một bảng:

| Trường | Nội dung |
|---|---|
| Prompt | câu đã hỏi, nguyên văn |
| Công cụ | tên AI + model |
| AI trả gì | tóm tắt 1-2 câu |
| Mình giữ / sửa / bỏ | và **vì sao** |
| Đã verify thế nào | chạy test? đối chiếu PDF? tính tay? |

Tối thiểu 5-6 entry, phủ các chủ đề PDF gợi ý:
- "How do I sort a slice of structs by a field in Go?"
- "Help me debug this Go function that calculates monthly payments"
- "What's the best way to handle user input in Go?"
- công thức annuity cho khoản vay
- vì sao `fmt.Scanln` gây loop vô hạn
- table-driven test trong Go

**Phần giá trị nhất của file này là cột "sửa / bỏ".** Ghi lại ít nhất 1-2 lần AI
đưa thứ sai hoặc không phù hợp và mình đã bắt được — ví dụ gợi ý thư viện ngoài
khi lab yêu cầu stdlib, hay so sánh `float64` bằng `==` trong test. Log chỉ có
toàn "AI đúng, mình nhận hết" thì không cho thấy mình đánh giá được gì.

Kết file bằng 3-4 câu: AI giúp nhanh ở đâu, chỗ nào vẫn phải tự hiểu mới làm được.

### 4. Short ADAPT reflection → `REFLECTION.md`

Rất ngắn — 5 đoạn, mỗi đoạn 3-5 câu, theo 5 chữ ADAPT:

- **A — Ask**: mình đã hỏi gì (AI, bạn, giảng viên)? Câu hỏi nào mở ra hướng đi,
  câu nào hỏi sai trọng tâm?
- **D — Discover**: phát hiện gì về Go mà trước đó chưa biết? (nil slice `append`
  được; range map thứ tự ngẫu nhiên; `float64(years*12)` vì chia số nguyên; sort
  mutate slice của caller)
- **A — Apply**: áp dụng vào task nào cụ thể? Dẫn tên hàm hoặc số task.
- **P — Practice**: lặp lại / refactor gì? Lần đầu viết thế nào, sau đổi sao?
  (ví dụ: ban đầu `calculateDistrictStats` trả map, sau đổi sang slice để sort được)
- **T — Teach**: giải thích lại cho người khác thế nào? Chọn một khái niệm —
  value vs pointer receiver, hoặc vì sao `Scanln` loop vô hạn — và viết 2-3 câu
  đủ rõ cho người chưa biết Go.

Viết thật, đúng những gì đã gặp trong lab. Phần này không có đáp án đúng/sai,
nhưng viết chung chung kiểu "em học được nhiều" thì không có giá trị.

### Thứ tự làm 4 deliverable

Làm **sau** bước 8 của lộ trình (sau khi test pass):

9. Chạy `go test -v ./...` → paste vào `TEST_RESULTS.md`; chạy tay 8 case menu,
   điền cột "Thực tế".
10. `DESIGN_NOTE.md` — lúc này lý do thiết kế còn mới trong đầu.
11. `AI_LOG.md` — ghi dần **trong lúc** code thì chính xác hơn là nhớ lại cuối giờ.
    Mở file này từ đầu lab, mỗi lần hỏi AI thì thêm một dòng.
12. `REFLECTION.md` — cuối cùng, vì cần nhìn lại toàn bộ.

## Checklist trước khi nộp

1. `go build ./...` không lỗi.
2. `go vet ./...` sạch.
3. `gofmt -l .` không in file nào.
4. `go test ./...` pass.
5. So sánh "cheapest" ra đúng Budget Studio.
6. Search với budget 3 tỷ ra 2 property.
7. ROI ra 12.0 / 10.0 / 18.0%.
8. Monthly payment ≈ 17.4 triệu cho Saigon Apartment (20% down, 8.5%, 20y).
9. Menu: nhập chữ, số âm, Enter rỗng — không crash, không loop vô hạn.
10. Portfolio với 8 tỷ chọn cả 3, còn 500 triệu.
11. `TEST_RESULTS.md` có output thật của `go test -v`, 8 case menu đã điền cột
    "Thực tế", và ghi rõ số pass/fail.
12. `DESIGN_NOTE.md` nêu đủ 3 chỗ cố ý lệch khỏi PDF.
13. `AI_LOG.md` có >= 5 entry, trong đó >= 1 entry mình sửa/bỏ đề xuất của AI.
14. `REFLECTION.md` đủ 5 đoạn A-D-A-P-T.

