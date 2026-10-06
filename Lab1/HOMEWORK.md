# 9. HOMEWORK — One-week extension

Phần mở rộng một tuần của Lab 01: thêm `CustomerType` (Regular/Member) và `OrderType` (Pickup/Delivery) vào model đã có, cập nhật công thức tính Final Total, và viết test cho mọi tổ hợp. Tài liệu này là phần nộp bài theo mục 9 trong `CLAUDE.md`.

## 9.1 Yêu cầu và công thức

```
Final Total = Subtotal - Customer Discount + Delivery Fee
```

- **Regular** — không giảm giá. **Member** — giảm 5% trên food subtotal.
- **Pickup** — phí 0 VND. **Delivery** — phí 30,000 VND.

Clarification bắt buộc của đề: member discount chỉ áp trên food subtotal, **không** làm giảm delivery fee; delivery fee luôn được cộng vào **sau** khi đã trừ discount.

## 9.2 Quyết định thiết kế: hai discount kết hợp thế nào

Đề bài yêu cầu làm rõ cách volume discount (10% khi subtotal ≥ 300,000 VND, từ ADAPT) kết hợp với customer discount (Member 5%). Ba phương án được cân nhắc:

| Phương án | Cách tính với subtotal 350,000 | Kết quả |
|---|---|---|
| Cộng dồn trên subtotal gốc | 350,000 − 35,000 − 17,500 | 297,500 |
| Áp tuần tự (5% trên số đã giảm 10%) | 350,000 − 35,000 → 315,000 − 15,750 | 299,250 |
| Loại trừ, lấy mức cao hơn | 350,000 − 35,000 (bỏ qua 5%) | 315,000 |

**Phương án được chọn: cộng dồn trên subtotal gốc.** Lý do: đề bài nói member discount áp dụng "trên food subtotal" (*on the food subtotal*), và volume discount ở ADAPT cũng được định nghĩa trên subtotal. Hiểu nhất quán hai câu đó thì cả hai rate đều lấy cùng một cơ sở là subtotal gốc, không rate nào lấy kết quả của rate kia làm cơ sở. Phương án này cũng khiến hai discount độc lập với nhau — thứ tự gọi hai hàm không ảnh hưởng kết quả, nên không tồn tại câu hỏi "áp cái nào trước".

Công thức được implement đầy đủ:

```
Final Total = Subtotal - VolumeDiscount(Subtotal) - MemberDiscount(Subtotal) + DeliveryFee(OrderType)
```

## 9.3 Thay đổi trong code

### Model

`CustomerType` và `OrderType` được khai báo là named string type kèm hằng số, không phải `string` thuần:

```go
type CustomerType string

const (
    Regular CustomerType = "Regular"
    Member  CustomerType = "Member"
)

type OrderType string

const (
    Pickup   OrderType = "Pickup"
    Delivery OrderType = "Delivery"
)
```

Lý do chọn named type thay vì `string`: compiler chặn việc truyền lẫn hai loại (không thể đưa một `OrderType` vào tham số `CustomerType`), trong khi vẫn in ra được trực tiếp bằng `%s` trong order summary. Nếu dùng `int`/`iota` thì sẽ phải viết thêm hàm chuyển sang text để hiển thị.

Hai field mới được thêm vào `Order`, `FoodItem` và `OrderLine` **không** đổi:

```go
type Order struct {
    Lines    []OrderLine
    Customer CustomerType
    Type     OrderType
}
```

Đặt hai field này ở `Order` (không phải ở `OrderLine`) vì loại khách và hình thức nhận hàng là thuộc tính của cả đơn hàng, không phải của từng món — một order không thể vừa Pickup vừa Delivery.

### Business logic

Ba hàm mới, mỗi hàm là pure function và chỉ giữ một trách nhiệm:

```go
func CalculateCustomerDiscount(subtotal float64, customer CustomerType) float64 {
    if customer == Member {
        return subtotal * MemberDiscountRate
    }
    return 0.0
}

func CalculateDeliveryFee(orderType OrderType) float64 {
    if orderType == Delivery {
        return DeliveryFeeAmount
    }
    return 0.0
}

func CalculateFinalTotal(order Order) float64 {
    subtotal := CalculateSubtotal(order)
    volumeDiscount := CalculateDiscount(subtotal)
    customerDiscount := CalculateCustomerDiscount(subtotal, order.Customer)
    return subtotal - volumeDiscount - customerDiscount + CalculateDeliveryFee(order.Type)
}
```

`CalculateFinalTotal` là nơi duy nhất công thức được viết ra, và thứ tự các phép toán trong đó phản ánh đúng clarification của đề: hai discount bị trừ trước, delivery fee được cộng sau cùng — nên fee không bao giờ bị discount tác động.

`CalculateSubtotal`, `CalculateDiscount`, `AddFoodToOrder`, `IsValidQuantity` — toàn bộ các hàm đã verify ở giai đoạn VERIFY — **không bị sửa một dòng nào**. Thay đổi duy nhất ở phần hiển thị là `DisplayOrderSummary` in thêm ba dòng (`Member Discount`, `Delivery Fee`, và dòng `Customer | Order`), và gọi `CalculateFinalTotal` thay vì tự tính `subtotal - discount`.

### CLI

Thêm hai option để đổi customer type và order type ngay trong phiên làm việc; Exit chuyển từ option 5 sang 7:

```
1. Display available food
2. Search food by name
3. Add food to order
4. Show order summary
5. Set customer type (Regular/Member)
6. Set order type (Pickup/Delivery)
7. Exit
```

Order mới khởi tạo mặc định là `Regular` + `Pickup` — tức là trường hợp không giảm giá, không phụ phí, nên hành vi mặc định của chương trình giống hệt trước khi làm homework. Input không khớp `Regular`/`Member` (hoặc `Pickup`/`Delivery`) bị từ chối kèm thông báo lỗi và **không** làm thay đổi order, cùng nguyên tắc với cách `AddFoodToOrder` xử lý input sai.

## 9.4 Test cases và kết quả

### Ma trận tổ hợp bắt buộc

Đề yêu cầu test 4 tổ hợp (Regular/Member × Pickup/Delivery) nhân với 3 mốc subtotal (dưới / đúng / trên 300,000 VND) = **12 case**, hiện thực bằng table-driven test `TestFinalTotalCombinations` trong `main_test.go`.

| Subtotal | Customer | Order | Volume (10%) | Member (5%) | Fee | Final Total | Result |
|---|---|---|---|---|---|---|---|
| 200,000 | Regular | Pickup | 0 | 0 | 0 | 200,000 | PASS |
| 200,000 | Regular | Delivery | 0 | 0 | 30,000 | 230,000 | PASS |
| 200,000 | Member | Pickup | 0 | 10,000 | 0 | 190,000 | PASS |
| 200,000 | Member | Delivery | 0 | 10,000 | 30,000 | 220,000 | PASS |
| 300,000 | Regular | Pickup | 30,000 | 0 | 0 | 270,000 | PASS |
| 300,000 | Regular | Delivery | 30,000 | 0 | 30,000 | 300,000 | PASS |
| 300,000 | Member | Pickup | 30,000 | 15,000 | 0 | 255,000 | PASS |
| 300,000 | Member | Delivery | 30,000 | 15,000 | 30,000 | 285,000 | PASS |
| 400,000 | Regular | Pickup | 40,000 | 0 | 0 | 360,000 | PASS |
| 400,000 | Regular | Delivery | 40,000 | 0 | 30,000 | 390,000 | PASS |
| 400,000 | Member | Pickup | 40,000 | 20,000 | 0 | 340,000 | PASS |
| 400,000 | Member | Delivery | 40,000 | 20,000 | 30,000 | 370,000 | PASS |

Mốc 300,000 là boundary case quan trọng nhất: điều kiện là `>=` nên đúng 300,000 **có** nhận volume discount. Một chi tiết đáng chú ý ở dòng `300,000 | Regular | Delivery`: Final Total = 300,000, bằng đúng subtotal — vì discount 30,000 và delivery fee 30,000 triệt tiêu nhau. Con số trùng hợp này có thể che lỗi nếu chỉ kiểm tra riêng nó, nên các dòng còn lại trong bảng mới là phần xác nhận thật.

### Test bổ sung

| Test function | Mục đích | Result |
|---|---|---|
| `TestCustomerDiscount` | Regular → 0; Member → 5% subtotal | PASS |
| `TestDeliveryFee` | Pickup → 0; Delivery → 30,000 | PASS |
| `TestFinalTotalCombinations` | 12 case trong bảng trên | PASS |
| `TestMemberDiscountDoesNotReduceDeliveryFee` | Kiểm tra trực tiếp clarification của đề | PASS |

`TestMemberDiscountDoesNotReduceDeliveryFee` được viết riêng vì nó test đúng cái ràng buộc dễ làm sai nhất của homework: nếu ai đó lỡ implement thành `(subtotal + fee) * 0.95` thì member discount sẽ ăn vào cả delivery fee. Test xác nhận khoảng cách giữa Final Total của Pickup và Delivery luôn đúng bằng 30,000 cho **cả** Regular và Member — nếu fee bị discount, khoảng cách của Member sẽ chỉ còn 28,500 và test fail.

Helper `orderWorth` dựng order có subtotal chính xác bằng số cần test (dùng Iced Coffee 20,000 VND/đơn vị) và tự `t.Fatalf` nếu subtotal dựng ra không đúng kỳ vọng — tránh trường hợp test pass vì lý do sai.

### Log chạy test

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
=== RUN   TestDiscountBelowThreshold
--- PASS: TestDiscountBelowThreshold (0.00s)
=== RUN   TestDiscountAtThreshold
--- PASS: TestDiscountAtThreshold (0.00s)
=== RUN   TestDiscountAboveThreshold
--- PASS: TestDiscountAboveThreshold (0.00s)
=== RUN   TestCustomerDiscount
--- PASS: TestCustomerDiscount (0.00s)
=== RUN   TestDeliveryFee
--- PASS: TestDeliveryFee (0.00s)
=== RUN   TestFinalTotalCombinations
    --- PASS: TestFinalTotalCombinations/Regular+Pickup_below (0.00s)
    --- PASS: TestFinalTotalCombinations/Regular+Delivery_below (0.00s)
    --- PASS: TestFinalTotalCombinations/Member+Pickup_below (0.00s)
    --- PASS: TestFinalTotalCombinations/Member+Delivery_below (0.00s)
    --- PASS: TestFinalTotalCombinations/Regular+Pickup_at (0.00s)
    --- PASS: TestFinalTotalCombinations/Regular+Delivery_at (0.00s)
    --- PASS: TestFinalTotalCombinations/Member+Pickup_at (0.00s)
    --- PASS: TestFinalTotalCombinations/Member+Delivery_at (0.00s)
    --- PASS: TestFinalTotalCombinations/Regular+Pickup_above (0.00s)
    --- PASS: TestFinalTotalCombinations/Regular+Delivery_above (0.00s)
    --- PASS: TestFinalTotalCombinations/Member+Pickup_above (0.00s)
    --- PASS: TestFinalTotalCombinations/Member+Delivery_above (0.00s)
--- PASS: TestFinalTotalCombinations (0.00s)
=== RUN   TestMemberDiscountDoesNotReduceDeliveryFee
--- PASS: TestMemberDiscountDoesNotReduceDeliveryFee (0.00s)
PASS
ok      foodordering    0.488s
```

25/25 test PASS, bao gồm toàn bộ 10 test cũ của VERIFY/ADAPT — xác nhận thay đổi của homework không gây regression. `go build ./...` và `go vet ./...` đều sạch.

### Xác nhận bằng CLI thật

Order summary với Member + Delivery (Chicken Rice x7 = 315,000 + Milk Tea x2 = 50,000):

```
Order summary:
- Chicken Rice x7 = 315000 VND
- Milk Tea x2 = 50000 VND
Customer: Member | Order: Delivery
Subtotal: 365000
Discount: 36500
Member Discount: 18250
Delivery Fee: 30000
Final Total: 340250
```

Kiểm tra tay: 365,000 − 36,500 − 18,250 + 30,000 = 340,250 → khớp.

## 9.5 AI review

Hàm được đưa cho AI review (không rewrite) là `CalculateFinalTotal`, vì đây là nơi duy nhất công thức của homework được viết ra và là điểm dễ sai nhất về thứ tự phép toán.

### Nhận xét của AI

- Thứ tự phép toán trong `return` thể hiện đúng clarification của đề: `DeliveryFee` nằm ngoài mọi phép trừ discount, nên không có đường nào để discount tác động lên fee. Nếu viết thành `(subtotal + fee) - volumeDiscount - customerDiscount` thì kết quả số vẫn giống hệt, nhưng ý định thiết kế không còn đọc được từ code.
- `subtotal` được tính một lần và truyền cho cả hai hàm discount, nên cả hai chắc chắn dùng cùng một cơ sở. Đây chính là chỗ thể hiện lựa chọn "cộng dồn trên subtotal gốc" ở mục 9.2 — nếu muốn đổi sang áp tuần tự, chỉ cần truyền `subtotal - volumeDiscount` vào `CalculateCustomerDiscount`.
- Mỗi hàm con là pure function (chỉ phụ thuộc tham số, không mutate gì), nên test được độc lập mà không cần dựng cả một `Order`.
- Điểm cần biết khi mở rộng: `CalculateFinalTotal` có thể trả về số âm nếu tổng discount vượt subtotal. Với hai rate hiện tại (10% + 5% = 15%) điều này không thể xảy ra, nhưng nếu thêm rule tới mức tổng vượt 100% thì cần một `math.Max(total, 0)` hoặc chặn ở tầng cấu hình rate.
- Một `Order` có `Customer` là zero value (`""`, do khởi tạo `Order{}` mà không set field) sẽ được tính như Regular, vì điều kiện là `customer == Member`. Đây là hành vi an toàn (không giảm giá sai cho khách không rõ loại) nhưng là một giả định ngầm đáng ghi lại.

### So sánh với lý luận của nhóm

Điểm đầu tiên trùng với lý luận của nhóm khi viết hàm: thứ tự phép toán được chọn **có chủ ý** để code tự đọc được thành công thức trong đề, dù về mặt số học có nhiều cách viết tương đương. Điểm thứ hai cũng vậy — việc tính `subtotal` một lần rồi truyền vào cả hai hàm là cách nhóm hiện thực hóa quyết định ở mục 9.2.

Hai điểm cuối là observation nhóm chưa nghĩ tới. Về khả năng âm: nhóm xác nhận không cần sửa vì 15% tổng discount còn rất xa ngưỡng 100%, thêm `math.Max` lúc này là phòng vệ cho tình huống không thể xảy ra. Về zero value của `Customer`: nhóm đánh giá đây là hành vi đúng hướng (mặc định không giảm giá thì an toàn hơn mặc định giảm giá), và CLI đã luôn khởi tạo tường minh `Order{Customer: Regular, Type: Pickup}` nên trường hợp này không xảy ra trong luồng chạy thật. **Không có thay đổi code nào** được thực hiện từ review này, đúng yêu cầu "review without rewriting".

## 9.6 Requirement change được xử lý thế nào

Thay đổi được hấp thụ bằng cách **thêm**, gần như không **sửa**: 3 hàm mới, 2 field mới trên `Order`, 2 named type mới. Không một hàm nào đã verify ở VERIFY/ADAPT bị sửa logic — bằng chứng là cả 10 test cũ vẫn PASS mà không phải chỉnh một dòng test nào.

Điều này khả thi nhờ hai quyết định có từ giai đoạn trước. Thứ nhất, `CalculateSubtotal` được tách khỏi phần hiển thị ngay từ BUILD, nên các tầng tính toán mới (discount, fee, final total) chèn được vào giữa tính toán và hiển thị mà không chạm vào nó. Thứ hai, `CalculateDiscount` ở ADAPT đã là pure function nhận `float64` và trả `float64`, nên `CalculateCustomerDiscount` chỉ cần theo đúng khuôn đó và cả hai cộng dồn được trong `CalculateFinalTotal`.

Đúng như reflection ở `ADAPT.md` mục 5.4 đã dự đoán: phần khó nhất của homework **không** phải viết code, mà là quyết định hai discount kết hợp ra sao (mục 9.2) — chính là câu hỏi "ưu tiên, cộng dồn, hay loại trừ" đã nêu từ trước. Việc thêm rule thứ hai làm lộ ra giới hạn của thiết kế hiện tại: `CalculateFinalTotal` giờ gọi tên từng rule một cách tường minh, nên mỗi rule mới đều buộc phải sửa hàm này. Với 2 rule thì vẫn rõ ràng và dễ đọc hơn mọi phương án trừu tượng hóa; nếu lên tới 4–5 rule thì điểm cắt hợp lý là đổi sang một `[]DiscountRule` được áp lần lượt, lúc đó `CalculateFinalTotal` chỉ còn duyệt danh sách và không cần sửa khi thêm rule.
