# 6. GROUP CHECK — 15 minutes

Giảng viên kiểm tra cả nhóm (pair) cùng lúc. Mục đích không phải "bắt" việc dùng AI, mà xác nhận cả hai thành viên đều hiểu và own được solution — bất kỳ ai trong nhóm cũng phải trả lời được bất kỳ câu hỏi nào dưới đây. File này chuẩn bị câu trả lời cho từng possible challenge trong đề, dựa trên code thật trong `main.go`.

## 6.1 Explain your struct

```go
type FoodItem struct {
    Name      string
    Price     float64
    Category  string
    Available bool
}

type OrderLine struct {
    Item     FoodItem
    Quantity int
}

type Order struct {
    Lines []OrderLine
}
```

`FoodItem` chỉ chứa đúng 4 field đề bài yêu cầu: Name, Price, Category, Available — không thêm field nào khác.

`Order` không lưu trực tiếp `[]FoodItem` và `[]int` (quantities) song song như ví dụ `CalculateSubtotal` trong đề, vì hai slice song song dễ bị lệch độ dài (length mismatch) nếu một slice được thêm/xóa mà quên đồng bộ slice kia. Thay vào đó, `Order` chứa một slice `[]OrderLine`, và mỗi `OrderLine` ghép chung `FoodItem` với `Quantity` của chính nó trong cùng một struct. Nhờ vậy, mỗi dòng order luôn mang theo đúng quantity của nó, không thể lệch.

## 6.2 Explain why you used a slice/map

- **Slice (`menu []FoodItem`)**: menu cần giữ thứ tự hiển thị (hiển thị món theo đúng thứ tự khai báo) và cần duyệt toàn bộ khi hiển thị danh sách món còn hàng (`DisplayAvailableFood`) — slice là cấu trúc tự nhiên cho duyệt tuần tự toàn bộ tập dữ liệu.
- **Map (`menuByName map[string]FoodItem`)**: tra cứu theo tên (search by name, và kiểm tra tồn tại/unavailable trong `AddFoodToOrder`) là thao tác xảy ra lặp lại nhiều lần mỗi khi khách thêm món. Map cho lookup theo key với độ phức tạp trung bình O(1), so với việc phải duyệt tuyến tính (O(n)) toàn bộ slice mỗi lần tìm theo tên.
- Map được build một lần từ slice (`buildMenuByName(menu)`) khi khởi động chương trình, không phải lưu trùng dữ liệu một cách tùy tiện — slice là nguồn dữ liệu chính (source of truth) cho việc hiển thị, map là index phục vụ tra cứu nhanh.

## 6.3 Explain one important function

Hàm `AddFoodToOrder` được chọn vì nó gói toàn bộ business rule của hệ thống vào một chỗ:

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

Hàm nhận `*Order` (con trỏ) để sửa trực tiếp `order.Lines` của order đang có, không cần trả `Order` mới rồi gán lại ở nơi gọi. Thứ tự 3 điều kiện kiểm tra là cố ý: tồn tại → còn hàng → quantity hợp lệ, theo nguyên tắc fail-fast — không cần kiểm tra quantity nếu món đó đã không tồn tại hoặc đã hết hàng. Phần phân tích đầy đủ (bao gồm AI review) nằm trong `VERIFY.md` mục 4.3.

## 6.4 Find an edge case

Một edge case không được yêu cầu trực tiếp trong bảng test bắt buộc nhưng đã được code xử lý: **thêm cùng một món hai lần vào cùng một order** (ví dụ gọi `AddFoodToOrder` với "Milk Tea" quantity 2, sau đó lại gọi với "Milk Tea" quantity 3). Thiết kế hiện tại tạo ra **hai `OrderLine` riêng biệt** cho cùng một món, thay vì cộng dồn quantity vào dòng đã có. `CalculateSubtotal` vẫn tính đúng tổng (vì nó cộng qua từng `OrderLine` độc lập), nhưng `DisplayOrderSummary` sẽ hiển thị "Milk Tea" hai lần như hai dòng khác nhau thay vì một dòng gộp quantity = 5. Đây là hành vi được chấp nhận trong phạm vi đề bài (đề không yêu cầu merge dòng trùng), nhưng là điểm cần biết nếu được hỏi.

Một edge case khác: tên món khi search/add được so khớp **chính xác và phân biệt hoa thường** (exact, case-sensitive) vì map key là string nguyên bản. Gõ "milk tea" (chữ thường) sẽ bị coi là "food not found" dù "Milk Tea" có trong menu.

## 6.5 Fix a small bug

Quy trình fix bug mẫu đã được thực hiện và ghi lại đầy đủ trong `VERIFY.md` mục 4.2: cố ý xóa `* float64(line.Quantity)` trong `CalculateSubtotal`, chạy test thấy 2 test fail (`TestNormalOrderSubtotal`, `TestSeveralDifferentItemsCorrectTotal`), thêm lại phép nhân, chạy lại test thấy PASS toàn bộ. Nếu giảng viên đưa một bug khác, quy trình áp dụng là: chạy `go test -v ./...` để xác định chính xác test nào fail và giá trị expect/actual, đọc lại hàm liên quan, sửa, rồi `go test -v ./...` lại để xác nhận.

## 6.6 Add search by category

Chưa có trong `main.go` vì không thuộc yêu cầu bắt buộc của BUILD, nhưng đây là ví dụ hàm có thể thêm ngay tại chỗ nếu được yêu cầu trong buổi group check, vì nó dùng lại đúng pattern của `DisplayAvailableFood`:

```go
// SearchFoodByCategory returns every menu item that belongs to the given category.
func SearchFoodByCategory(items []FoodItem, category string) []FoodItem {
    var result []FoodItem
    for _, item := range items {
        if item.Category == category {
            result = append(result, item)
        }
    }
    return result
}
```

Hàm này duyệt slice `menu` (không dùng map, vì map hiện tại chỉ index theo Name, không theo Category) và trả về mọi món khớp category — tương tự cách `DisplayAvailableFood` lọc theo `Available`.

## 6.7 Explain an AI-generated part

Phần CLI loop trong `main()` (đọc input từ `bufio.Reader`, switch theo lựa chọn 1-5) là phần có hỗ trợ AI nhiều nhất trong giai đoạn BUILD, vì đây là boilerplate I/O lặp lại, không chứa business logic. Điểm quan trọng là: toàn bộ validate (not found/unavailable/invalid quantity) **không** nằm trong `main()` — `main()` chỉ gọi `AddFoodToOrder` và in ra lỗi nếu có. Business rule nằm hoàn toàn trong các hàm độc lập (`AddFoodToOrder`, `IsValidQuantity`, `CalculateSubtotal`, `CalculateDiscount`) để có thể test bằng `go test` mà không cần giả lập input dòng lệnh — đây là lý do `main_test.go` test trực tiếp các hàm đó, không test qua CLI.

## 6.8 Change the discount rule

Giả sử giảng viên yêu cầu đổi rule, ví dụ discount 15% khi subtotal ≥ 500,000 VND. Thay đổi chỉ cần sửa hai hằng số trong `main.go`:

```go
const DiscountThreshold = 500000.0
const DiscountRate = 0.15
```

`CalculateDiscount` không cần sửa logic vì nó đã được viết tổng quát theo `DiscountThreshold`/`DiscountRate`, không hard-code số 300000/0.10 trực tiếp trong hàm. Đây cũng chính là nội dung đã phân tích trong `ADAPT.md` mục 5.4: dễ đổi vì chỉ có một rule tuyến tính; sẽ khó hơn nếu phải thêm **nhiều** rule cùng lúc, vì khi đó cần quyết định cách kết hợp (ưu tiên, cộng dồn, hay loại trừ nhau) — vượt ra ngoài việc chỉ đổi hằng số.
