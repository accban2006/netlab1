# REFLECTION — ADAPT

## A — Ask

Câu hỏi mở ra hướng đi tốt nhất là "vì sao `fmt.Scanln` gây loop vô hạn khi nhập
chữ" — nó không phải câu hỏi về cú pháp mà về cơ chế buffer của stdin, và trả lời
được nó thì cả Part 5 sáng ra. Câu hỏi hỏi sai trọng tâm là "cách in số VND có
dấu phẩy phân cách nghìn trong Go": tôi gõ câu đó mà quên nói lab chỉ cho dùng
stdlib, nên nhận về đề xuất `golang.org/x/text` — đúng câu hỏi, sai bài. Bài học
là phải đưa ràng buộc vào trong prompt, không chờ AI đoán. Lần sau tôi hỏi
"...chỉ dùng stdlib" thì câu trả lời ra ngay đúng hướng.

## D — Discover

Bốn thứ về Go mà trước đó tôi chưa biết. Một, `append` vào nil slice là hợp lệ —
`var result []Property` rồi `append` ngay không cần `make`, và `len(nil) == 0`,
`range` trên nil cũng an toàn, nên không phải viết guard nil ở mọi hàm search.
Hai, `range` trên map được **cố tình** randomize, không phải "tình cờ không có
thứ tự" — Go làm vậy để lập trình viên không vô tình phụ thuộc vào thứ tự. Ba,
`years * 12` là phép chia/nhân số nguyên, phải `float64(years*12)` trước khi đưa
vào công thức float, nếu không kết quả bị cắt phần thập phân một cách âm thầm.
Bốn, và đáng nhớ nhất: slice là `{ptr, len, cap}` nên truyền vào hàm là copy
struct nhưng dùng chung backing array — `sort` trong hàm đổi luôn slice của
caller.

## A — Apply

Phát hiện về map áp vào Task 2.2: `analyzeByDistrict` vẫn trả map đúng như PDF,
nhưng `calculateDistrictStats` trả `[]DistrictStats` đã sort, và chỗ nào buộc in
từ map thì lấy keys ra `sort.Strings` trước. Phát hiện về chia số nguyên áp vào
`calculateMonthlyPayment` và `CalculateLoan` — chỗ tính `TotalInterest` có
`monthly*float64(years*12)`, thiếu ép kiểu là sai hẳn kết quả. Phát hiện về slice
áp vào Task 4.2: `rankByROI` zip (property, ROI) vào slice mới rồi sort bản đó,
và kéo theo việc sửa `pdfProperties` từ `properties[:3]` thành
`append([]Property(nil), properties[:3]...)` vì slice cắt ra cũng dùng chung
backing array.

## P — Practice

`calculateDistrictStats` tôi viết lại hai lần. Bản đầu trả `map[string]DistrictStats`
vì nghĩ "dữ liệu theo quận thì để map cho tra nhanh". Đến lúc in theo yêu cầu
"sorted by average price" mới thấy không sort được map, phải đổ ra slice rồi
sort — tức là map chỉ là bước trung gian vô ích. Bản hai trả luôn `[]DistrictStats`.
`optimizePortfolio` cũng qua ba vòng: sort trực tiếp tham số (mutate caller) →
copy `[]Property` rồi sort (lệch index với `rents`) → zip thành `[]propertyROI`
rồi sort. Mỗi lần sửa tôi thêm một test trước khi sửa, nên biết chắc bản mới
thật sự khác bản cũ chứ không phải trông có vẻ khác.

## T — Teach

Giải thích "vì sao `fmt.Scanln` gây loop vô hạn" cho người chưa biết Go: hãy
tưởng tượng stdin là một hàng người chờ vào cửa. `fmt.Scanln(&choice)` với
`choice` là số sẽ nhìn người đầu hàng, thấy không phải số thì nói "không hợp lệ"
rồi **bỏ đi mà không cho người đó ra khỏi hàng**. Vòng `for` quay lại, lại nhìn
đúng người đó, lại báo lỗi — mãi mãi. `bufio.Reader.ReadString('\n')` thì khác:
nó đưa cả người đó ra khỏi hàng trước, rồi mới xem có phải số hay không. Hàng
luôn ngắn đi sau mỗi lượt, nên không bao giờ lặp vô hạn. Đó là lý do đọc-rồi-parse
an toàn hơn parse-trực-tiếp khi xử lý input người dùng.
