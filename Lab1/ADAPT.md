# 5. ADAPT — 15 minutes

Giai đoạn ADAPT mô phỏng một tình huống thường gặp trong phát triển phần mềm: requirement thay đổi sau khi hệ thống đã được build và verify. Mục tiêu không chỉ là implement đúng yêu cầu mới, mà còn đánh giá mức độ dễ/khó khi thay đổi dựa trên thiết kế đã chọn ở THINK.

## 5.1 Yêu cầu mới

Order có subtotal từ 300,000 VND trở lên nhận discount 10% trên subtotal đó.

Format hiển thị bắt buộc:

```
Subtotal: 350000
Discount: 35000
Final Total: 315000
```

## 5.2 Thay đổi trong code

Thay đổi được thực hiện tại hai điểm trong `main.go`, không ảnh hưởng đến struct `FoodItem`, `Order`, hay các hàm đã verify ở giai đoạn trước:

```go
const DiscountThreshold = 300000.0
const DiscountRate = 0.10

func CalculateDiscount(subtotal float64) float64 {
    if subtotal >= DiscountThreshold {
        return subtotal * DiscountRate
    }
    return 0.0
}
```

Và `DisplayOrderSummary` được cập nhật để gọi `CalculateDiscount` và in thêm hai dòng `Discount` và `Final Total`:

```go
subtotal := CalculateSubtotal(order)
discount := CalculateDiscount(subtotal)
finalTotal := subtotal - discount
fmt.Printf("Subtotal: %.0f\n", subtotal)
fmt.Printf("Discount: %.0f\n", discount)
fmt.Printf("Final Total: %.0f\n", finalTotal)
```

`CalculateSubtotal` — hàm đã được verify kỹ ở giai đoạn VERIFY — không bị sửa. Discount được implement như một hàm thuần (pure function) độc lập, nhận subtotal và trả discount, không biết gì về `Order` hay `FoodItem`.

## 5.3 Test case

Ba case bắt buộc được kiểm tra bằng cả automated test (`TestDiscountBelowThreshold`, `TestDiscountAtThreshold`, `TestDiscountAboveThreshold` trong `main_test.go`) và bằng cách chạy CLI thực tế để xác nhận format hiển thị đúng yêu cầu.

### Case 1 — Below 300,000 VND

Order: Iced Coffee x5 (20,000 × 5 = 100,000).

```
Order summary:
- Iced Coffee x5 = 100000 VND
Subtotal: 100000
Discount: 0
Final Total: 100000
```

Discount = 0 vì 100,000 < 300,000 → đúng.

### Case 2 — Exactly 300,000 VND

Order: Beef Noodle x6 (50,000 × 6 = 300,000).

```
Order summary:
- Beef Noodle x6 = 300000 VND
Subtotal: 300000
Discount: 30000
Final Total: 270000
```

Discount = 300,000 × 10% = 30,000 vì điều kiện là `>=` (bao gồm đúng ngưỡng) → đúng. Đây là boundary case quan trọng nhất vì đề dùng từ "at least" (≥), không phải "more than" (>).

### Case 3 — Above 300,000 VND

Order: Chicken Rice x7 (45,000 × 7 = 315,000) + Milk Tea x2 (25,000 × 2 = 50,000) = 365,000.

```
Order summary:
- Chicken Rice x7 = 315000 VND
- Milk Tea x2 = 50000 VND
Subtotal: 365000
Discount: 36500
Final Total: 328500
```

Discount = 365,000 × 10% = 36,500; Final Total = 365,000 − 36,500 = 328,500 → đúng.

### Kết quả automated test

```
=== RUN   TestDiscountBelowThreshold
--- PASS: TestDiscountBelowThreshold (0.00s)
=== RUN   TestDiscountAtThreshold
--- PASS: TestDiscountAtThreshold (0.00s)
=== RUN   TestDiscountAboveThreshold
--- PASS: TestDiscountAboveThreshold (0.00s)
PASS
```

## 5.4 Reflection

Thiết kế ban đầu dễ thay đổi vì `CalculateSubtotal` đã được tách riêng khỏi phần hiển thị (`DisplayOrderSummary`) ngay từ đầu, nên discount có thể được chèn vào giữa hai bước đó như một hàm độc lập mà không cần sửa logic subtotal đã verify. Phần sẽ trở nên khó khăn nếu có nhiều discount rule (ví dụ: giảm theo category, giảm theo loại khách hàng, giảm theo mã khuyến mãi) là việc `CalculateDiscount` hiện tại chỉ nhận một `float64 subtotal` và áp một rule cố định — khi có nhiều rule cùng áp dụng, cần quyết định rule nào ưu tiên, có cộng dồn (stack) được không, và có thể phải đổi từ một hàm sang một danh sách rule (slice of discount rules) được áp lần lượt, đòi hỏi thiết kế lại interface của hàm discount thay vì chỉ thêm hằng số.
