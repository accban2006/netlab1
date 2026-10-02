# AI_LOG — Lab 01 Go Foundation

Công cụ dùng xuyên suốt: Claude Code (Claude Opus). Ghi dần trong lúc code, không
viết lại từ nhớ cuối giờ.

---

## Entry 1 — Công thức khoản vay

| Trường | Nội dung |
|---|---|
| **Prompt** | "Giải thích công thức annuity trong starter code của PDF: `loanAmount * monthlyRate * math.Pow(1+monthlyRate, numPayments) / (math.Pow(1+monthlyRate, numPayments) - 1)`. Vì sao phải có nhánh `if annualRate == 0`?" |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | Giải thích đây là công thức khoản trả cố định hàng tháng: tử số là phần lãi trên dư nợ ban đầu nhân hệ số tăng trưởng, mẫu số chuẩn hoá về chuỗi thanh toán n kỳ. Khi `annualRate == 0` thì `monthlyRate = 0`, cả tử và mẫu đều bằng 0 → `0/0 = NaN`, nên phải tách nhánh chia đều `loanAmount / numPayments`. |
| **Giữ / sửa / bỏ** | **Giữ** phần giải thích, **thêm** một guard AI không nhắc: `years = 0` làm `numPayments = 0`, và lúc đó nhánh lãi-0% cũng ra `x/0`. Tôi thêm `if numPayments == 0 { return 0 }` lên đầu hàm. |
| **Verify thế nào** | Tính tay với loan 2e9, r = 8.5/100/12 = 0.0070833, n = 240 → 17,356,464.67, khớp với `17.3 triệu` của PDF (PDF làm tròn xuống). Viết `TestCalculateMonthlyPayment/years_zero_guard` để chốt guard mới. |

---

## Entry 2 — Sort slice of structs

| Trường | Nội dung |
|---|---|
| **Prompt** | "How do I sort a slice of structs by a field in Go?" (prompt gợi ý trong mục Getting Help của PDF) |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | Gợi ý `sort.Slice(s, func(i, j int) bool { return s[i].Field > s[j].Field })`, và nói thêm `sort.Slice` dùng quicksort nên **không stable**; muốn stable thì dùng `sort.SliceStable`. |
| **Giữ / sửa / bỏ** | **Giữ** `sort.Slice` nhưng **thêm tie-break** thay vì đổi sang `sort.SliceStable`. Lý do: input của `calculateDistrictStats` đến từ `range` trên map, mà thứ tự range map vốn đã ngẫu nhiên — `SliceStable` chỉ giữ thứ tự input, mà thứ tự input ở đây không ổn định thì stable cũng vô nghĩa. Tie-break theo tên quận mới cho output lặp lại được. |
| **Verify thế nào** | Chạy `go run .` chọn mục 7 năm lần liên tiếp, thứ tự `District 1 / District 2 / District 7 / Binh Thanh` không đổi. Trước khi thêm tie-break, thứ tự đổi giữa các lần chạy. |

---

## Entry 3 — Xử lý input người dùng

| Trường | Nội dung |
|---|---|
| **Prompt** | "What's the best way to handle user input in Go?" rồi hỏi tiếp "vì sao `fmt.Scanln(&choice)` với choice là int lại gây loop vô hạn khi người dùng nhập chữ?" |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | `Scanln` parse theo kiểu của con trỏ truyền vào; gặp token không parse được thì trả error nhưng **không** consume token đó khỏi buffer. Vòng `for` lặp lại, `Scanln` đọc lại đúng `abc` cũ, fail tiếp → menu in vô hạn. Giải pháp: `bufio.Reader.ReadString('\n')` đọc hết dòng kể cả phần không hợp lệ, rồi `strconv.Atoi`. |
| **Giữ / sửa / bỏ** | **Giữ** hướng `bufio` + `strconv`. **Sửa** chỗ xử lý EOF: code AI đưa chỉ `if err != nil { continue }`, như vậy khi chạy qua pipe hết input sẽ lại loop vô hạn với `err = io.EOF` — đúng cái bug vừa định sửa, chỉ đổi nguyên nhân. Tôi tách `isEOF()` để EOF thì `return`, lỗi parse thì `continue`. |
| **Verify thế nào** | `printf 'abc\n0\n' \| go run .` → báo lỗi một lần rồi về menu (case M3). `printf '' \| go run .` → `EOF nhận được. Goodbye!` rồi exit (case M7). Hai case này trước đó đều treo. |

---

## Entry 4 — AI gợi ý so sánh float bằng `==` (đã bỏ)

| Trường | Nội dung |
|---|---|
| **Prompt** | "Viết table-driven test cho `CalculateROI` và `calculateMonthlyPayment`." |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | Khung table-driven đúng idiom (`[]struct{...}` + `t.Run(tc.name, ...)`), nhưng phần so sánh viết `if got := p.CalculateROI(tc.rent); got != tc.want`. |
| **Giữ / sửa / bỏ** | **Giữ** khung table-driven, **bỏ** phép so sánh `!=`. So sánh `float64` bằng `==`/`!=` là sai: 12.0 tính từ `25e6*12/2.5e9*100` không nhất thiết bằng đúng literal `12.0` trong biểu diễn nhị phân. Tôi viết helper `almostEqual(got, want, tol)` dùng `math.Abs(got-want) < tol`, và dùng tolerance khác nhau theo độ lớn: `1e-6` cho %, `±1 VND` cho khoản trả hàng tháng, `±1000 VND` cho tổng lãi (sai số cộng dồn qua 240 kỳ). |
| **Verify thế nào** | Thử giữ nguyên `!=` cho case `calculateMonthlyPayment(2e9, 8.5, 20)` với want `17356465`: fail vì giá trị thật là `17356464.67`. Đổi sang `almostEqual(..., 1)` thì pass. Đây là bằng chứng trực tiếp rằng `==` sẽ làm test giòn. |

---

## Entry 5 — AI gợi ý thư viện ngoài để format tiền (đã bỏ)

| Trường | Nội dung |
|---|---|
| **Prompt** | "Cách in số VND có dấu phẩy phân cách nghìn trong Go?" |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | Đề xuất `golang.org/x/text/message` với `message.NewPrinter(language.Vietnamese)` — gọn, đúng locale, `p.Printf("%d", n)` tự chèn dấu phân cách. |
| **Giữ / sửa / bỏ** | **Bỏ hoàn toàn.** Lab yêu cầu chỉ dùng stdlib (`fmt`, `math`, `sort`, `strings`, `os`, `bufio`, `strconv`), mà `golang.org/x/text` là module ngoài — thêm vào là phải `go get`, có `go.sum`, và vi phạm ràng buộc đề bài. Tôi tự viết `formatVND` bằng `strings.Builder`: in `%.0f`, rồi chèn `,` mỗi khi số chữ số còn lại chia hết cho 3. Xử lý thêm dấu âm vì `strings.Builder` không biết gì về dấu. |
| **Verify thế nào** | `TestFormatVND` với 6 case: 2.5 tỷ, giá/m² có phần thập phân (33112582.78 → `33,112,583`), số 3 chữ số (`800`, không dấu phẩy), 4 chữ số (`1,234`), zero, và số âm (`-1,234,567`). Case 3 chữ số là case tôi gần như quên — logic `(len-i)%3 == 0` sẽ chèn dấu phẩy ở đầu nếu không có điều kiện `i > 0`. |

---

## Entry 6 — Debug: slice bị mutate sau khi sort

| Trường | Nội dung |
|---|---|
| **Prompt** | "Help me debug this Go function that calculates monthly payments" — thực tế câu hỏi đổi hướng thành: "sau khi gọi `optimizePortfolio(properties, 8e9)`, biến `properties` ở ngoài bị đổi thứ tự. Vì sao?" |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | Slice trong Go là struct {ptr, len, cap} trỏ vào một backing array. Truyền slice vào hàm là copy struct đó, nhưng `ptr` vẫn trỏ cùng array — nên `sort.Slice` ghi đè tại chỗ và caller thấy thay đổi. Cách sửa: `sorted := make([]Property, len(properties)); copy(sorted, properties)` trước khi sort. |
| **Giữ / sửa / bỏ** | **Giữ** chẩn đoán, **đổi cách sửa**. Tôi cần sort theo ROI mà ROI lại nằm ở slice `rents` riêng — copy rồi sort `[]Property` sẽ làm lệch cặp (property, rent) theo index. Nên tôi zip thành `[]propertyROI{Property, ROI}` rồi sort struct đó: vừa là bản copy, vừa không thể lệch index. |
| **Verify thế nào** | Viết `TestOptimizePortfolioDoesNotMutateCaller`: lưu danh sách tên trước khi gọi, gọi cả `optimizePortfolioWithRents` và `optimizePortfolio`, rồi so lại từng tên. Test này fail với bản đầu (sort trực tiếp tham số), pass sau khi zip. Phát hiện kéo theo: `pdfProperties = properties[:3]` cũng dùng chung backing array nên tôi đổi sang `append([]Property(nil), ...)`. |

---

## Entry 7 — Tại sao ROI ở Task 4.1 luôn giống nhau

| Trường | Nội dung |
|---|---|
| **Prompt** | "Trong `recommendProperty`, PDF dùng `p.CalculateROI(p.Price * 0.012)`. Tôi chạy thử thì cả 5 property đều ra `BUY NOW`. Có phải tôi làm sai?" |
| **Công cụ** | Claude Code (Claude Opus) |
| **AI trả gì** | Không sai. Thế `monthlyRent = Price * 0.012` vào công thức ROI: `(Price*0.012*12 / Price) * 100` — `Price` triệt tiêu, còn `0.012*12*100 = 14.4` với mọi property. Giả định "thuê 1.2% giá trị mỗi tháng" làm ROI độc lập hoàn toàn với property. |
| **Giữ / sửa / bỏ** | **Giữ** đúng công thức của PDF (đề bài yêu cầu vậy), nhưng **ghi comment** nêu rõ quan sát này tại `part4_logic.go`, và dựa khác biệt giữa các khuyến nghị vào bonus/warning thay vì ROI. Không "sửa cho đẹp" bằng cách bịa rent khác — như vậy sẽ lệch đề. |
| **Verify thế nào** | Tính tay: 0.012 × 12 × 100 = 14.4 cho cả 5 property, đúng như output. Test `TestSmartRecommendProperty/warning_ha_mot_bac` dựng một property giá/m² = 75M (3e9 / 40m²) để chứng minh nhánh warning có thật sự hạ bậc BUY NOW → GOOD BUY, tức là logic phân bậc hoạt động dù ROI bằng nhau. |

---

## Tổng kết

AI nhanh nhất ở hai việc: giải thích **tại sao** một đoạn code có sẵn hành xử lạ
(Entry 1, 3, 6, 7) và dựng khung boilerplate theo idiom chuẩn (table-driven test
ở Entry 4). Cả hai đều là loại việc mà tự mò sẽ mất 10–20 phút tra tài liệu.

Chỗ AI không thay được: biết **ràng buộc của bài này**. Entry 5 là ví dụ rõ nhất —
`golang.org/x/text` là câu trả lời đúng cho câu hỏi tôi gõ, nhưng sai với lab chỉ
cho dùng stdlib; AI không biết điều đó vì tôi không nói. Entry 4 thì AI đưa code
chạy được mà ẩn một lỗi khó thấy (`==` với float), nếu nhận hết thì test vẫn
"pass" hôm nay rồi vỡ khi đổi dữ liệu.

Điều rút ra: mọi thứ AI đưa đều phải chạy thật hoặc tính tay lại mới tính là
xong. Sáu trong bảy entry trên, phần có giá trị nhất không phải câu trả lời đầu
tiên mà là cái guard hoặc test tôi thêm vào *sau* khi đã hiểu câu trả lời đó —
`years = 0` ở Entry 1, nhánh EOF ở Entry 3, tie-break ở Entry 2.
