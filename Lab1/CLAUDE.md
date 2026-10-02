# Lab 01 — Go Foundation: Food Ordering System

Kế hoạch chi tiết cho IT096IIU — Netcentric, Lab 01. Nhóm 2 người, 150 phút, 100 điểm.

Nguyên tắc cốt lõi: **Think together. Build with AI. Verify together. Own the result.**

Lab flow: THINK (AI OFF) → BUILD (AI ON) → VERIFY (AI ON) → ADAPT (AI ON) → GROUP CHECK (Own it).

---

## 1. Scenario (Bối cảnh)

Xây dựng một Food Ordering System nhỏ cho quán ăn trong campus. Hệ thống quản lý menu món ăn và cho phép khách tạo order.

Mỗi món ăn (FoodItem) có:
- Name
- Price
- Category
- Available

Danh mục ban đầu (Initial categories): Main Dish, Drink, Dessert.

Hệ thống phải hỗ trợ:
1. Hiển thị món ăn còn hàng (available).
2. Tìm món ăn theo tên (search by name).
3. Thêm món vào order.
4. Từ chối số lượng không hợp lệ (invalid quantities).
5. Từ chối món ăn không còn hàng (unavailable food).
6. Tính subtotal của order.
7. Hiển thị order summary.
8. Áp dụng discount khi yêu cầu thay đổi (ADAPT phase).

---

## 2. THINK — 20 phút, AI OFF

Làm việc theo cặp, **không dùng AI** trong giai đoạn này.

### Câu hỏi cần thảo luận và trả lời ngắn gọn
1. FoodItem cần những thông tin gì?
2. Order cần những thông tin gì?
3. Một order có thể chứa nhiều món ăn khác nhau không?
4. Quantity nên được biểu diễn như thế nào?
5. Food items nên được lưu ở đâu?
6. Map có thể hữu ích ở đâu?
7. Điều gì xảy ra khi món được yêu cầu không tồn tại?
8. Điều gì xảy ra khi quantity bằng 0 hoặc âm?
9. Điều gì xảy ra khi món ăn unavailable?
10. Function/method nào nên chịu trách nhiệm cho business rules?

### Thiết kế ngắn cần tạo ra
- Struct `FoodItem`
- Struct `Order`
- Ít nhất một slice và một map
- Các function/method quan trọng
- Hai thiết kế khả thi (two possible designs)
- Thiết kế được chọn và lý do chọn (one reason)
- Bốn trường hợp lỗi dự đoán (four predicted failure cases)

**Trước khi code**, phải show thiết kế cho giảng viên. Giảng viên có thể clarify requirement/concept nhưng sẽ không cho full solution.

---

## 3. BUILD — 60 phút, AI ON

Được dùng AI tools. AI có thể generate code, nhưng cả hai thành viên **phải hiểu và own code**.

### Ví dụ tham khảo trong đề (không bắt buộc copy nguyên bản)

**3.1 FoodItem struct ví dụ:**
```go
type FoodItem struct {
    Name      string
    Price     float64
    Category  string
    Available bool
}
```

**3.2 Slice ví dụ (menu):**
```go
menu := []FoodItem{
    {Name: "Chicken Rice", Price: 45000, Category: "Main Dish", Available: true},
    {Name: "Milk Tea", Price: 25000, Category: "Drink", Available: true},
    {Name: "Cheesecake", Price: 35000, Category: "Dessert", Available: false},
}
```
Cần thảo luận: tại sao slice phù hợp, cách tìm món, điều gì xảy ra khi món unavailable.

**3.3 Map lookup ví dụ:**
```go
menuByName := map[string]FoodItem{
    "Chicken Rice": {Name: "Chicken Rice", Price: 45000, Category: "Main Dish", Available: true},
    "Milk Tea":     {Name: "Milk Tea", Price: 25000, Category: "Drink", Available: true},
}

food, ok := menuByName["Milk Tea"]
if !ok {
    fmt.Println("Food not found")
} else {
    fmt.Println(food.Name, food.Price)
}
```
Không copy nguyên si — tự quyết định final design dùng slice, map, hay cả hai.

**3.4 Subtotal function ví dụ:**
```go
func CalculateSubtotal(items []FoodItem, quantities []int) float64 {
    total := 0.0
    for i := range items {
        total += items[i].Price * float64(quantities[i])
    }
    return total
}
```
Cần thảo luận các assumption: độ dài slice khác nhau, quantity âm, validation nên đặt ở đâu, có data structure nào an toàn hơn không.

**3.5 Validation function ví dụ:**
```go
func IsValidQuantity(quantity int) bool {
    return quantity > 0
}
```

### Implementation phải có
- Ít nhất 6 food items và 3 categories
- Hiển thị món còn hàng + search by name
- Thêm món, validate quantity, reject unavailable food
- Subtotal và order summary
- Ít nhất một slice và một map
- Functions và/hoặc methods
- CLI menu (optional)

### AI collaboration log
Ghi lại ít nhất 3 tương tác AI có ý nghĩa, bao gồm ít nhất một yêu cầu review/verification. Không nộp transcript dài dòng.

Bảng log mẫu:

| # | What we asked AI | What AI suggested | What we changed/verified |
|---|---|---|---|
| 1 | | | |
| 2 | | | |
| 3 | | | |

---

## 4. VERIFY — 25 phút

### Test case bắt buộc tối thiểu

| Case | Expected behavior | Result |
|---|---|---|
| Normal order | Correct subtotal | |
| Quantity = 1 | Accepted | |
| Quantity = 0 | Rejected | |
| Negative quantity | Rejected | |
| Unavailable food | Rejected | |
| Food not found | Handled safely | |
| Several different items | Correct total | |

### Defect injection
Cố ý đưa vào một lỗi nhỏ, ví dụ:
```go
total += items[i].Price
```
thay vì:
```go
total += items[i].Price * float64(quantities[i])
```
Sau đó: **identify → explain → fix → retest** lỗi này.

### AI review
Yêu cầu AI review một function quan trọng **mà không rewrite nó**. So sánh AI review với lý luận (reasoning) của chính mình.

---

## 5. ADAPT — 15 phút

### Yêu cầu mới
Order có subtotal ≥ 300,000 VND nhận discount 10%.

Hiển thị theo format:
```
Subtotal: 350000
Discount: 35000
Final Total: 315000
```

### Test case
1. Dưới 300,000 VND
2. Đúng 300,000 VND
3. Trên 300,000 VND

### Reflection (2–3 câu)
Trả lời: Thiết kế ban đầu dễ hay khó thay đổi? Phần nào sẽ trở nên khó khăn nếu có nhiều discount rule được thêm vào?

---

## 6. GROUP CHECK — 15 phút

Giảng viên kiểm tra cả nhóm (cặp). Cả hai thành viên phải giải thích được solution.

### Các câu hỏi/thử thách có thể gặp
- Giải thích struct.
- Giải thích tại sao dùng slice/map.
- Giải thích một function quan trọng.
- Tìm một edge case.
- Fix một bug nhỏ.
- Thêm search by category.
- Giải thích một phần code do AI generate.
- Thay đổi discount rule.

> Mục đích không phải để "bắt" việc dùng AI, mà để xác nhận cặp **own** solution của mình.

Giảng viên có thể hỏi riêng một thành viên một câu follow-up ngắn nếu cần clarify.

---

## 7. Submission (Nộp bài)

Cần nộp:
1. Source code
2. Test cases/results
3. Short design note
4. AI collaboration log
5. Short ADAPT reflection

Cả hai thành viên đều chịu trách nhiệm cho bài nộp cuối cùng.

---

## 8. Assessment (Thang điểm)

| Evidence | Weight |
|---|---|
| THINK | 20 |
| BUILD | 20 |
| VERIFY | 20 |
| ADAPT | 15 |
| GROUP CHECK | 25 |
| **Total** | **100** |
