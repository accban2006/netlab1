# TEST_RESULTS — Lab 01 Go Foundation

Môi trường: `go1.27.1 windows/amd64`, module `property-calculator`, Git Bash.

## Tổng kết

| Loại | Số case | Pass | Fail |
|---|---|---|---|
| Unit test tự động (`go test -v ./...`) | 32 test function / 99 subtest | 131 | 0 |
| Manual test menu | 8 (M1–M8) + 1 phụ (M9) | 9 | 0 |

Không có case nào fail. Tất cả output dưới đây lấy từ lần chạy thật, không viết
tay trước.

## A. Unit test tự động

### A.1 Bảng case bắt buộc (theo CLAUDE.md)

| # | Hàm | Input | Expected | Thực tế | Lý do chọn case |
|---|---|---|---|---|---|
| 1 | `formatPrice` | `2500000000` | `2.5 tỷ VND` | `2.5 tỷ VND` ✅ | nhánh tỷ |
| 2 | `formatPrice` | `850000000` | `850 triệu VND` | `850 triệu VND` ✅ | nhánh triệu |
| 3 | `formatPrice` | `0` | `0 VND` | `0 VND` ✅ | biên zero |
| 4 | `categorizeProperty` | `50000000` | `PREMIUM` | `PREMIUM` ✅ | biên `>` không phải `>=` |
| 5 | `categorizeProperty` | `50000001` | `LUXURY` | `LUXURY` ✅ | vượt biên 1 đơn vị |
| 6 | `PricePerM2` | `Area: 0` | `0` | `0` ✅ | guard chia 0 |
| 7 | `CalculateROI` | rent `25e6`, price `2.5e9` | `12.0` | `12.0` ✅ | khớp PDF |
| 8 | `CalculateROI` | `Price: 0` | `0` | `0` ✅ | guard chia 0 |
| 9 | `InvestmentGrade` | `8.0` | `GOOD` | `GOOD` ✅ | biên EXCELLENT/GOOD |
| 10 | `calculateMonthlyPayment` | `2e9, 8.5, 20` | `≈17356465` | `17356464.67` ✅ | so sánh epsilon ±1 |
| 11 | `calculateMonthlyPayment` | `annualRate: 0` | `L/n` | `20000000` (2.4e9/120) ✅ | nhánh lãi 0% |
| 12 | `CalculateLoan` | `20, 8.5, 20` | loan `2e9`, lãi `≈2.166e9` | loan `2e9`, lãi `2165551520` ✅ | case bắt buộc PDF |
| 13 | `CalculateLoan` | `downPayment: 150` | `LoanInfo{}` | `LoanInfo{}` ✅ | input không hợp lệ |
| 14 | `findPropertiesInBudget` | `3e9` | 2 kết quả | 2 (Saigon, Studio) ✅ | lọc đúng |
| 15 | `findPropertiesInBudget` | `1e6` | `len == 0` | `len == 0` ✅ | không khớp gì |
| 16 | `findBestInvestment` | rents lệch độ dài | không panic | `Property{}`, không panic ✅ | guard index |
| 17 | `findBestInvestment` | slice rỗng | `Property{}` | `Property{}` ✅ | biên rỗng |
| 18 | `calculateDistrictStats` | 3 property | sort giảm dần | D7 → D1 → Binh Thanh ✅ | thứ tự đúng |
| 19 | `optimizePortfolio` | `8e9` | chọn 3, còn `500e6` | chọn 3, còn `500000000` ✅ | case bắt buộc PDF |
| 20 | `optimizePortfolio` | `0` | `len == 0` | `len == 0` ✅ | budget zero |
| 21 | `optimizePortfolio` | — | slice gốc không đổi | thứ tự gốc giữ nguyên ✅ | không mutate caller |

### A.2 Case bổ sung đã thêm

Ngoài 21 case trên, bộ test có thêm các nhóm sau:

- **`formatPriceDetailed`** (3 case): bản giữ thập phân cho khoản trả hàng tháng.
  Phát hiện khi chạy thật: `formatPrice` làm tròn `%.0f` nên 17,356,465 in ra
  "17 triệu VND" thay vì "17.4 triệu" như PDF. Thêm formatter riêng.
- **`formatVND`** (6 case): dấu phẩy phân cách nghìn, gồm số âm và số ≤ 3 chữ số.
- **`categorizeProperty`** (8 case): cả 3 biên đều test cặp "bằng" / "vượt 1".
- **`IsAffordable`** (4 case): biên `<=` — đúng bằng giá vẫn affordable.
- **`InvestmentGrade`** (9 case): cả 3 biên 8 / 5 / 3.
- **`CalculateLoanChecked`** (5 case invalid + 1 case 100% trả trước hợp lệ).
- **`findPropertiesByDistrict`** (4 case): không phân biệt hoa thường, bỏ khoảng trắng.
- **`calculateDistrictStats`** với nhiều property trong cùng quận: kiểm `Count`,
  `AveragePrice`, `MostExpensive`.
- **`smartRecommendProperty`** (5 case): gồm case warning hạ bậc BUY NOW → GOOD BUY.
- **`rankByROI`** tie-break: hai property cùng ROI phải ra thứ tự xác định.
- **`readInt` / `readFloat`** (14 case + 1 case đọc nhiều dòng): số âm, chữ,
  Enter rỗng dùng default, EOF, dòng cuối không có `\n`.

Mọi so sánh `float64` dùng `math.Abs(got-want) < tol` qua helper `almostEqual`
(`property_test.go:13`), không dùng `==`.

### A.3 Output `go test -v ./...`

```
=== RUN   TestFindPropertiesInBudget
=== RUN   TestFindPropertiesInBudget/3_ty
=== RUN   TestFindPropertiesInBudget/tat_ca
=== RUN   TestFindPropertiesInBudget/bang_gia_re_nhat
=== RUN   TestFindPropertiesInBudget/khong_khop_gi
=== RUN   TestFindPropertiesInBudget/budget_zero
--- PASS: TestFindPropertiesInBudget (0.00s)
    --- PASS: TestFindPropertiesInBudget/3_ty (0.00s)
    --- PASS: TestFindPropertiesInBudget/tat_ca (0.00s)
    --- PASS: TestFindPropertiesInBudget/bang_gia_re_nhat (0.00s)
    --- PASS: TestFindPropertiesInBudget/khong_khop_gi (0.00s)
    --- PASS: TestFindPropertiesInBudget/budget_zero (0.00s)
--- PASS: TestFindPropertiesByBedrooms (0.00s)
    --- PASS: TestFindPropertiesByBedrooms/1_phong (0.00s)
    --- PASS: TestFindPropertiesByBedrooms/2_phong (0.00s)
    --- PASS: TestFindPropertiesByBedrooms/3_phong (0.00s)
    --- PASS: TestFindPropertiesByBedrooms/10_phong_khong_co (0.00s)
    --- PASS: TestFindPropertiesByBedrooms/0_phong (0.00s)
--- PASS: TestFindPropertiesByDistrict (0.00s)
    --- PASS: TestFindPropertiesByDistrict/khop_chinh_xac (0.00s)
    --- PASS: TestFindPropertiesByDistrict/khong_phan_biet_hoa_thuong (0.00s)
    --- PASS: TestFindPropertiesByDistrict/bo_khoang_trang (0.00s)
    --- PASS: TestFindPropertiesByDistrict/khong_ton_tai (0.00s)
--- PASS: TestAnalyzeByDistrict (0.00s)
--- PASS: TestCalculateDistrictStats (0.00s)
--- PASS: TestCalculateDistrictStatsMultipleInDistrict (0.00s)
--- PASS: TestCalculateDistrictStatsEmpty (0.00s)
--- PASS: TestFindBestInvestment (0.00s)
--- PASS: TestFindBestInvestmentGuards (0.00s)
    --- PASS: TestFindBestInvestmentGuards/rents_lech_do_dai (0.00s)
    --- PASS: TestFindBestInvestmentGuards/slice_rong (0.00s)
    --- PASS: TestFindBestInvestmentGuards/mot_property (0.00s)
--- PASS: TestFormatPrice (0.00s)
    --- PASS: TestFormatPrice/ty (0.00s)
    --- PASS: TestFormatPrice/trieu (0.00s)
    --- PASS: TestFormatPrice/nghin (0.00s)
    --- PASS: TestFormatPrice/zero (0.00s)
    --- PASS: TestFormatPrice/bien_1_ty (0.00s)
    --- PASS: TestFormatPrice/duoi_bien_1_ty (0.00s)
--- PASS: TestFormatPriceDetailed (0.00s)
    --- PASS: TestFormatPriceDetailed/monthly_payment (0.00s)
    --- PASS: TestFormatPriceDetailed/total_interest (0.00s)
    --- PASS: TestFormatPriceDetailed/zero (0.00s)
--- PASS: TestFormatVND (0.00s)
    --- PASS: TestFormatVND/ty (0.00s)
    --- PASS: TestFormatVND/price_per_m2 (0.00s)
    --- PASS: TestFormatVND/ba_chu_so (0.00s)
    --- PASS: TestFormatVND/bon_chu_so (0.00s)
    --- PASS: TestFormatVND/zero (0.00s)
    --- PASS: TestFormatVND/am (0.00s)
--- PASS: TestCategorizeProperty (0.00s)
    --- PASS: TestCategorizeProperty/bien_luxury_bang (0.00s)
    --- PASS: TestCategorizeProperty/bien_luxury_vuot (0.00s)
    --- PASS: TestCategorizeProperty/bien_premium_bang (0.00s)
    --- PASS: TestCategorizeProperty/bien_premium_vuot (0.00s)
    --- PASS: TestCategorizeProperty/bien_standard_bang (0.00s)
    --- PASS: TestCategorizeProperty/bien_standard_vuot (0.00s)
    --- PASS: TestCategorizeProperty/zero (0.00s)
    --- PASS: TestCategorizeProperty/saigon_apartment (0.00s)
--- PASS: TestRecommendProperty (0.00s)
    --- PASS: TestRecommendProperty/over_budget (0.00s)
    --- PASS: TestRecommendProperty/monthly_payment_qua_cao (0.00s)
    --- PASS: TestRecommendProperty/roi_tot (0.00s)
--- PASS: TestSmartRecommendProperty (0.00s)
    --- PASS: TestSmartRecommendProperty/premium_location_va_optimal_size (0.00s)
    --- PASS: TestSmartRecommendProperty/premium_location_nhung_qua_to (0.00s)
    --- PASS: TestSmartRecommendProperty/khong_bonus (0.00s)
    --- PASS: TestSmartRecommendProperty/warning_ha_mot_bac (0.00s)
    --- PASS: TestSmartRecommendProperty/over_budget (0.00s)
--- PASS: TestOptimizePortfolio (0.00s)
--- PASS: TestOptimizePortfolioBudgetLimits (0.00s)
    --- PASS: TestOptimizePortfolioBudgetLimits/budget_zero (0.00s)
    --- PASS: TestOptimizePortfolioBudgetLimits/budget_am (0.00s)
    --- PASS: TestOptimizePortfolioBudgetLimits/chi_du_studio (0.00s)
    --- PASS: TestOptimizePortfolioBudgetLimits/greedy_bo_qua_dat_nhat (0.00s)
    --- PASS: TestOptimizePortfolioBudgetLimits/du_tat_ca (0.00s)
--- PASS: TestOptimizePortfolioDoesNotMutateCaller (0.00s)
--- PASS: TestOptimizePortfolioMismatchedRents (0.00s)
--- PASS: TestRankByROI (0.00s)
--- PASS: TestRankByROITieBreak (0.00s)
--- PASS: TestReadInt (0.00s)
    --- PASS: TestReadInt/so_hop_le (0.00s)
    --- PASS: TestReadInt/co_khoang_trang (0.00s)
    --- PASS: TestReadInt/so_am (0.00s)
    --- PASS: TestReadInt/enter_rong_dung_default (0.00s)
    --- PASS: TestReadInt/khong_phai_so (0.00s)
    --- PASS: TestReadInt/so_thuc_khong_hop_le (0.00s)
    --- PASS: TestReadInt/eof (0.00s)
    --- PASS: TestReadInt/khong_co_newline_cuoi (0.00s)
--- PASS: TestReadFloat (0.00s)
    --- PASS: TestReadFloat/so_nguyen (0.00s)
    --- PASS: TestReadFloat/so_thuc (0.00s)
    --- PASS: TestReadFloat/ky_hieu_khoa_hoc (0.00s)
    --- PASS: TestReadFloat/enter_rong_dung_default (0.00s)
    --- PASS: TestReadFloat/khong_phai_so (0.00s)
    --- PASS: TestReadFloat/eof (0.00s)
--- PASS: TestReadLineMultipleLines (0.00s)
--- PASS: TestPricePerM2 (0.00s)
    --- PASS: TestPricePerM2/saigon_apartment (0.00s)
    --- PASS: TestPricePerM2/hcmc_house (0.00s)
    --- PASS: TestPricePerM2/area_zero_guard (0.00s)
    --- PASS: TestPricePerM2/zero_value (0.00s)
--- PASS: TestIsAffordable (0.00s)
    --- PASS: TestIsAffordable/du_budget (0.00s)
    --- PASS: TestIsAffordable/bang_gia (0.00s)
    --- PASS: TestIsAffordable/thieu_budget (0.00s)
    --- PASS: TestIsAffordable/budget_zero (0.00s)
--- PASS: TestCalculateROI (0.00s)
    --- PASS: TestCalculateROI/saigon_apartment (0.00s)
    --- PASS: TestCalculateROI/hcmc_house (0.00s)
    --- PASS: TestCalculateROI/budget_studio (0.00s)
    --- PASS: TestCalculateROI/price_zero_guard (0.00s)
    --- PASS: TestCalculateROI/rent_zero (0.00s)
--- PASS: TestInvestmentGrade (0.00s)
    --- PASS: TestInvestmentGrade/excellent (0.00s)
    --- PASS: TestInvestmentGrade/bien_excellent_good (0.00s)
    --- PASS: TestInvestmentGrade/vuot_bien_excellent (0.00s)
    --- PASS: TestInvestmentGrade/good (0.00s)
    --- PASS: TestInvestmentGrade/bien_good_fair (0.00s)
    --- PASS: TestInvestmentGrade/fair (0.00s)
    --- PASS: TestInvestmentGrade/bien_fair_poor (0.00s)
    --- PASS: TestInvestmentGrade/poor (0.00s)
    --- PASS: TestInvestmentGrade/zero (0.00s)
--- PASS: TestCalculateMonthlyPayment (0.00s)
    --- PASS: TestCalculateMonthlyPayment/pdf_saigon_apartment (0.00s)
    --- PASS: TestCalculateMonthlyPayment/pdf_hcmc_house (0.00s)
    --- PASS: TestCalculateMonthlyPayment/lai_zero (0.00s)
    --- PASS: TestCalculateMonthlyPayment/loan_zero (0.00s)
    --- PASS: TestCalculateMonthlyPayment/years_zero_guard (0.00s)
--- PASS: TestCalculateLoan (0.00s)
--- PASS: TestCalculateLoanInvalidInput (0.00s)
    --- PASS: TestCalculateLoanInvalidInput/down_payment_qua_100 (0.00s)
    --- PASS: TestCalculateLoanInvalidInput/down_payment_am (0.00s)
    --- PASS: TestCalculateLoanInvalidInput/lai_am (0.00s)
    --- PASS: TestCalculateLoanInvalidInput/years_zero (0.00s)
    --- PASS: TestCalculateLoanInvalidInput/years_am (0.00s)
--- PASS: TestCalculateLoanFullDownPayment (0.00s)
PASS
ok  	property-calculator	0.503s
```

(Các dòng `=== RUN` lặp lại cho từng subtest đã lược bớt sau test function đầu
tiên để file đọc được; phần `--- PASS` giữ nguyên đầy đủ.)

## B. Manual test menu

Chạy bằng cách pipe input vào `go run .` để kết quả lặp lại được.
Ví dụ M3: `printf 'abc\n0\n' | go run .`

| # | Thao tác | Kỳ vọng | Thực tế |
|---|---|---|---|
| M1 | chọn `1` | in danh sách property | ✅ in bảng 5 property kèm giá/m² |
| M2 | chọn `2`, budget `3000000000` | 2 kết quả | ✅ `Saigon Apartment` + `Budget Studio` |
| M3 | nhập `abc` | báo lỗi, hiện lại menu | ✅ `Input không hợp lệ: "abc" không phải số nguyên hợp lệ` rồi in lại menu |
| M4 | nhập `99` | `Invalid option!` | ✅ `Invalid option!` |
| M5 | Enter rỗng ở prompt budget | dùng giá trị mặc định | ✅ dùng 3,000,000,000 (`Properties under 3,000,000,000 (3.0 tỷ VND)`) |
| M6 | nhập `-5` cho số năm | báo lỗi, không crash | ✅ `Lỗi: số năm phải lớn hơn 0 (nhận -5)` |
| M7 | `Ctrl+D` / EOF | thoát sạch, không loop vô hạn | ✅ `EOF nhận được. Goodbye!` rồi exit |
| M8 | chọn `0` | `Goodbye!` rồi exit | ✅ `Goodbye!` |
| M9 | chọn `2`, budget `-5000` | báo lỗi budget âm | ✅ `Budget không được âm.` |

M3 là case quan trọng nhất: với `fmt.Scanln(&choice)` của starter code PDF, nhập
`abc` sẽ làm vòng `for` quay vô hạn in menu liên tục. Bản `bufio.Reader` đọc hết
dòng nên không còn ký tự rác trong buffer.

## C. Đối chiếu output với PDF

| Mục | PDF | Thực tế | Ghi chú |
|---|---|---|---|
| Saigon Apartment giá/m² | `33,113` | `33,112,583` | PDF in rút gọn; giá trị đúng của 2.5e9/75.5 |
| Cheapest per m² | Budget Studio | `Budget Studio at 22,857,143 VND/m²` | khớp |
| Saigon Apartment category | `STANDARD` | `PREMIUM` | **lệch có chủ ý** — xem `DESIGN_NOTE.md` |
| ROI 3 property | 12.0 / 10.0 / 18.0% | 12.0 / 10.0 / 18.0% | khớp |
| Best Investment | Budget Studio (18.0%) | Budget Studio (18.0% ROI) | khớp |
| Monthly payment Saigon | `17.3 triệu` | `17.4 triệu` | PDF làm tròn xuống; giá trị đúng 17,356,464.67 |
| Total Interest Saigon | `2.15 tỷ` | `2.17 tỷ` | tính từ công thức: 17356464.67×240 − 2e9 = 2,165,551,520 |
| Monthly payment HCMC | `29.1 triệu` | `29.2 triệu` | giá trị đúng 29,158,860.64 |
| Total Interest HCMC | `3.62 tỷ` | `3.64 tỷ` | tính từ công thức |
| Portfolio 8 tỷ | 3 property, còn 500 triệu | 3 property, còn 500 triệu | khớp |
| Portfolio Average ROI | `13.3%` | `13.3%` | trung bình cộng (18+12+10)/3 |

Các chênh lệch ở khoản vay là do PDF làm tròn khác, không phải lỗi công thức:
test `TestCalculateLoan` kiểm chứng quan hệ `TotalInterest = monthly×n − loan`
nên nội bộ nhất quán.

## D. Output các lệnh kiểm tra

```console
$ go build ./...
(không output — build thành công)

$ go vet ./...
(không output — sạch)

$ gofmt -l .
(không output — mọi file đã format đúng)

$ go test ./...
ok  	property-calculator	0.503s
```
