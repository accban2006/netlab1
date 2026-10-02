# 4. VERIFY — 25 minutes

Giai đoạn VERIFY kiểm chứng rằng implementation ở giai đoạn BUILD thỏa mãn các business rule đã đặt ra ở THINK: validate quantity, reject unavailable food, tính đúng subtotal. Việc kiểm chứng được thực hiện bằng automated test (`go test`), không phải chỉ chạy tay qua CLI, để kết quả có thể tái lập (reproducible) và chạy lại bất cứ lúc nào.

Toàn bộ test case nằm trong `main_test.go` và tương ứng 1-1 với bảng case bắt buộc của đề bài.

## 4.1 Bảng test case bắt buộc

| Case | Expected behavior | Test function | Result |
|---|---|---|---|
| Normal order | Correct subtotal | `TestNormalOrderSubtotal` | PASS |
| Quantity = 1 | Accepted | `TestQuantityOneAccepted` | PASS |
| Quantity = 0 | Rejected | `TestQuantityZeroRejected` | PASS |
| Negative quantity | Rejected | `TestNegativeQuantityRejected` | PASS |
| Unavailable food | Rejected | `TestUnavailableFoodRejected` | PASS |
| Food not found | Handled safely | `TestFoodNotFoundHandledSafely` | PASS |
| Several different items | Correct total | `TestSeveralDifferentItemsCorrectTotal` | PASS |

### Chi tiết từng case

**Normal order — Correct subtotal.** Thêm "Chicken Rice" (45,000 VND) với quantity 2 vào một order mới. `CalculateSubtotal` phải trả về 90,000. Kết quả thực tế: 90,000 → đúng.

**Quantity = 1 — Accepted.** Thêm "Milk Tea" với quantity 1. `AddFoodToOrder` phải trả `err == nil`. Kết quả thực tế: không có lỗi → đúng. Case này xác nhận ranh giới dưới của `IsValidQuantity` (quantity > 0) hoạt động đúng cho giá trị hợp lệ nhỏ nhất.

**Quantity = 0 — Rejected.** Thêm "Milk Tea" với quantity 0. `AddFoodToOrder` phải trả về lỗi. Kết quả thực tế: trả lỗi `invalid quantity: 0` → đúng.

**Negative quantity — Rejected.** Thêm "Milk Tea" với quantity -3. Kết quả thực tế: trả lỗi `invalid quantity: -3` → đúng.

**Unavailable food — Rejected.** Thêm "Fried Rice" (có trong menu nhưng `Available: false`). Kết quả thực tế: trả lỗi `food unavailable: Fried Rice` → đúng. Case này xác nhận việc kiểm tra `Available` diễn ra trước khi item được thêm vào order, không phải sau.

**Food not found — Handled safely.** Thêm một tên không tồn tại trong menu ("Pho"). Yêu cầu là hệ thống phải "handled safely", nghĩa là không panic, không crash, và order không bị thay đổi. Kết quả thực tế: trả lỗi `food not found: Pho`, và `len(order.Lines) == 0` sau lệnh gọi → đúng.

**Several different items — Correct total.** Thêm 3 món khác nhau vào cùng một order: Chicken Rice x1 (45,000), Milk Tea x2 (50,000), Ice Cream x3 (45,000). Tổng kỳ vọng: 140,000. Kết quả thực tế: 140,000 → đúng. Case này xác nhận rằng `CalculateSubtotal` cộng đúng trên nhiều `OrderLine` khác nhau, không chỉ trên một dòng.

### Log chạy test (baseline, trước khi inject defect)

```
=== RUN   TestNormalOrderSubtotal
--- PASS: TestNormalOrderSubtotal (0.00s)
=== RUN   TestQuantityOneAccepted
--- PASS: TestQuantityOneAccepted (0.00s)
=== RUN   TestQuantityZeroRejected
--- PASS: TestQuantityZeroRejected (0.00s)
=== RUN   TestNegativeQuantityRejected
--- PASS: TestNegativeQuantityRejected (0.00s)
=== RUN   TestUnavailableFoodRejected
--- PASS: TestUnavailableFoodRejected (0.00s)
=== RUN   TestFoodNotFoundHandledSafely
--- PASS: TestFoodNotFoundHandledSafely (0.00s)
=== RUN   TestSeveralDifferentItemsCorrectTotal
--- PASS: TestSeveralDifferentItemsCorrectTotal (0.00s)
PASS
ok      foodordering    0.461s
```

Tất cả 7/7 case bắt buộc PASS.

## 4.2 Defect injection

### Lỗi được chủ động đưa vào

Theo yêu cầu của đề, một lỗi nhỏ được cố ý đưa vào hàm `CalculateSubtotal` trong `main.go`:

```go
// Trước khi inject (đúng)
for _, line := range order.Lines {
    total += line.Item.Price * float64(line.Quantity)
}
```

bị đổi thành:

```go
// Sau khi inject (lỗi)
for _, line := range order.Lines {
    total += line.Item.Price
}
```

### Identify (nhận diện lỗi)

Chạy lại test suite ngay sau khi inject lỗi:

```
=== RUN   TestNormalOrderSubtotal
    main_test.go:14: subtotal = 45000, want 90000
--- FAIL: TestNormalOrderSubtotal (0.00s)
=== RUN   TestQuantityOneAccepted
--- PASS: TestQuantityOneAccepted (0.00s)
=== RUN   TestQuantityZeroRejected
--- PASS: TestQuantityZeroRejected (0.00s)
=== RUN   TestNegativeQuantityRejected
--- PASS: TestNegativeQuantityRejected (0.00s)
=== RUN   TestUnavailableFoodRejected
--- PASS: TestUnavailableFoodRejected (0.00s)
=== RUN   TestFoodNotFoundHandledSafely
--- PASS: TestFoodNotFoundHandledSafely (0.00s)
=== RUN   TestSeveralDifferentItemsCorrectTotal
    main_test.go:71: subtotal = 85000, want 140000
--- FAIL: TestSeveralDifferentItemsCorrectTotal (0.00s)
FAIL
FAIL    foodordering    0.476s
```

Hai test fail ngay: `TestNormalOrderSubtotal` (mong đợi 90,000 nhưng nhận 45,000) và `TestSeveralDifferentItemsCorrectTotal` (mong đợi 140,000 nhưng nhận 85,000). Các test còn lại vẫn PASS vì chúng không phụ thuộc vào `CalculateSubtotal` với quantity > 1.

### Explain (giải thích nguyên nhân)

Lỗi là bỏ sót phép nhân với `quantity`. Hàm cộng `line.Item.Price` cho mỗi dòng order, bất kể dòng đó có quantity bao nhiêu — tức là coi mọi quantity như bằng 1. Với test "Normal order" (Chicken Rice x2, giá 45,000), kết quả đúng phải là 45,000 × 2 = 90,000, nhưng vì thiếu `* float64(line.Quantity)`, kết quả chỉ còn 45,000 (chỉ tính giá của 1 đơn vị). Tương tự, với "Several different items" (1 + 2 + 3 = 6 đơn vị nhưng mỗi dòng chỉ tính 1 đơn vị), subtotal bị thiếu đúng phần (quantity − 1) × price của mỗi dòng: 85,000 so với 140,000 đúng, chênh 55,000 — khớp với (2-1)×25,000 + (3-1)×15,000 = 25,000 + 30,000 = 55,000.

Đây là lớp lỗi điển hình khi tính tổng có trọng số (weighted sum): dễ quên nhân biến thứ hai nếu không có test kiểm tra quantity > 1 trên nhiều dòng khác nhau.

### Fix (sửa lỗi)

Khôi phục lại phép nhân với quantity:

```go
for _, line := range order.Lines {
    total += line.Item.Price * float64(line.Quantity)
}
```

### Retest (kiểm tra lại)

```
=== RUN   TestNormalOrderSubtotal
--- PASS: TestNormalOrderSubtotal (0.00s)
=== RUN   TestQuantityOneAccepted
--- PASS: TestQuantityOneAccepted (0.00s)
=== RUN   TestQuantityZeroRejected
--- PASS: TestQuantityZeroRejected (0.00s)
=== RUN   TestNegativeQuantityRejected
--- PASS: TestNegativeQuantityRejected (0.00s)
=== RUN   TestUnavailableFoodRejected
--- PASS: TestUnavailableFoodRejected (0.00s)
=== RUN   TestFoodNotFoundHandledSafely
--- PASS: TestFoodNotFoundHandledSafely (0.00s)
=== RUN   TestSeveralDifferentItemsCorrectTotal
--- PASS: TestSeveralDifferentItemsCorrectTotal (0.00s)
PASS
ok      foodordering    (cached)
```

Toàn bộ 7/7 test PASS trở lại. Lỗi được xác nhận đã sửa đúng.

## 4.3 AI review (không rewrite)

Yêu cầu AI review hàm `AddFoodToOrder` — một trong những hàm quan trọng nhất vì nó gói gọn toàn bộ business rule (not-found, unavailable, invalid quantity) — mà không để AI viết lại hàm.

```go
func AddFoodToOrder(order *Order, name string, quantity int) error {
    item, ok := SearchFoodByName(name)
    if !ok {
        return fmt.Errorf("food not found: %s", name)
    }
    if !item.Available {
        return fmt.Errorf("food unavailable: %s", name)
    }
    if !IsValidQuantity(quantity) {
        return fmt.Errorf("invalid quantity: %d", quantity)
    }
    order.Lines = append(order.Lines, OrderLine{Item: item, Quantity: quantity})
    return nil
}
```

### Nhận xét của AI

- Thứ tự kiểm tra (not found → unavailable → invalid quantity) hợp lý: fail nhanh (fail-fast) theo mức độ "cơ bản" của điều kiện — một món không tồn tại thì không cần kiểm tra tiếp các điều kiện khác.
- Hàm nhận `*Order` (con trỏ) để có thể mutate `order.Lines` tại chỗ, tránh phải trả về `Order` mới và gán lại ở nơi gọi — hợp lý vì `Order` được dùng xuyên suốt một session CLI.
- Dùng `fmt.Errorf` với message mô tả rõ nguyên nhân (not found/unavailable/invalid) giúp caller (CLI) hiển thị lỗi có ý nghĩa, thay vì một mã lỗi chung.
- Điểm có thể cải thiện nếu mở rộng về sau: nếu cùng một `name` được add nhiều lần, hàm hiện tại tạo nhiều `OrderLine` riêng biệt thay vì cộng dồn quantity vào dòng đã có. Đây không phải lỗi theo yêu cầu đề bài (đề không yêu cầu merge dòng trùng), nhưng là một giả định cần biết trước khi mở rộng tính năng.

### So sánh với lý luận của nhóm

Lý luận ban đầu của nhóm khi thiết kế hàm này trùng với điểm đầu tiên AI nêu: thứ tự kiểm tra được chọn chủ ý theo nguyên tắc fail-fast và theo đúng thứ tự liệt kê business rule trong đề (reject food not found trước, unavailable sau, invalid quantity sau cùng — vì quantity chỉ có ý nghĩa kiểm tra khi đã biết chắc món đó tồn tại và còn bán). Điểm về con trỏ `*Order` cũng là quyết định có chủ đích, không phải ngẫu nhiên.

Điểm AI nêu thêm (không merge dòng trùng tên) là một observation nhóm chưa viết thành test case rõ ràng, nhưng xác nhận đây là hành vi được chấp nhận trong phạm vi đề bài hiện tại — không cần sửa, chỉ cần ghi nhận làm giả định thiết kế. Không có thay đổi code nào được thực hiện dựa trên review này, đúng theo yêu cầu "review without rewriting".
