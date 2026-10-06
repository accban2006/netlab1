# Câu hỏi bảo vệ code: Food Ordering System

Tài liệu dựa trên `main.go`, `main_test.go`, `go.mod` tại thời điểm kiểm tra 06/10/2026. Câu trả lời mô tả hành vi hiện tại; đoạn code ở phần challenge là gợi ý thay đổi, chưa được thêm vào chương trình. Mỗi thành viên nên tự trace ví dụ và giải thích lý do, thay vì học thuộc câu trả lời.

## 1. Model và lựa chọn cấu trúc dữ liệu

### Câu 1. Vì sao tách `FoodItem`, `OrderLine` và `Order` thành ba struct? Nếu đặt `Quantity` ngay trong `FoodItem` thì vấn đề gì xảy ra?

`FoodItem` mô tả món trong menu: tên, giá, category và trạng thái còn bán. `OrderLine` mô tả một lần đặt món, nên ghép một `FoodItem` với quantity. `Order` gom nhiều dòng và giữ loại khách, cách nhận hàng cho cả đơn. Quantity không phải thuộc tính cố định của món: cùng Milk Tea, người A mua 2, người B mua 5. Đặt quantity trong menu sẽ trộn dữ liệu món với dữ liệu giao dịch. Tách ba struct giúp thay đổi quantity của một order mà không phải thay menu.

### Câu 2. Tại sao chọn `[]OrderLine` thay cho hai slice song song `[]FoodItem` và `[]int` như ví dụ trong đề? Thiết kế này loại bỏ lỗi nào và chưa loại bỏ lỗi nào?

Hai slice song song phải cùng độ dài và cùng thứ tự. Nếu thêm món mà quên thêm quantity, truy cập `quantities[i]` có thể panic; nếu đổi thứ tự một slice thì quantity có thể ghép nhầm món. `[]OrderLine` giữ cặp món–quantity trong cùng phần tử, loại bỏ việc lệch hai danh sách. Tuy nhiên struct vẫn cho phép `Quantity: -2`: cấu trúc dữ liệu bảo đảm quan hệ giữa món và quantity, còn tính hợp lệ của quantity phải do `AddFoodToOrder` kiểm tra.

### Câu 3. Tại sao `Customer` và `Type` nằm trong `Order`, không đặt ở từng `OrderLine`? Nếu một đơn có ba món, delivery fee được tính mấy lần?

Hai field mô tả người mua và cách nhận của cả đơn. Dự án chưa hỗ trợ giao từng món theo hình thức riêng. Nếu đặt chúng ở từng dòng, có thể tạo trạng thái một đơn vừa Pickup vừa Delivery hoặc thu phí ba lần. `CalculateFinalTotal` gọi `CalculateDeliveryFee(order.Type)` một lần sau khi tính subtotal toàn đơn, nên ba món vẫn chỉ có một phí 30,000 VND.

### Câu 4. Tại sao quantity dùng `int`, còn price dùng `float64`? Có bắt buộc dùng số thực để tính discount không?

Quantity là số phần nguyên, nên `int` phù hợp và `strconv.Atoi` chuyển input thành số nguyên. Price dùng `float64` theo ví dụ đề và để tính rate 0.10/0.05 đơn giản trong lab. Không bắt buộc dùng số thực: có thể lưu tiền bằng `int64` VND và tính `subtotal * 5 / 100`, kèm chính sách làm tròn và kiểm soát overflow. `float64` có thể sai số khi biểu diễn phân số; `%.0f` chỉ làm tròn lúc hiển thị, không thay đổi giá trị đang được dùng để tính toán.

### Câu 5. Tại sao dùng đồng thời slice `menu` và map `menuByName`? Với chỉ bảy món, lợi ích có đáng với việc giữ hai cấu trúc không?

Slice giữ thứ tự khai báo và thuận tiện duyệt tất cả món để hiển thị. Map tra một tên cụ thể trung bình O(1); tìm tuyến tính trong slice là O(n). Với bảy món, chênh lệch hiệu năng nhỏ, nên không nên nói map là tối ưu bắt buộc. Lựa chọn này phù hợp hai kiểu truy cập và yêu cầu bài có cả slice/map. Đổi lại, map cần O(n) bộ nhớ bổ sung và phải đồng bộ nếu menu thay đổi. Dự án hiện tại dùng menu tĩnh nên build map một lần đủ cho luồng đang có.

### Câu 6. `buildMenuByName` chạy lúc nào? Vì sao tạo map từ slice thay vì khai báo một danh sách món riêng cho map?

`menuByName = buildMenuByName(menu)` là initializer cấp package, được thực hiện trước `main` sau khi `menu` đã khởi tạo. Hàm tạo map bằng `make`, duyệt slice rồi gán `m[item.Name] = item`. Dùng cùng dữ liệu đầu vào tránh phải nhập lại giá/trạng thái từng món ở hai nơi. Nó tạo index ban đầu thống nhất, nhưng không tạo liên kết tự cập nhật giữa slice và map.

### Câu 7. Nếu menu có hai món cùng tên hoặc thay đổi giá sau khi map đã build, chương trình sẽ làm gì?

Map chỉ giữ một value cho mỗi key. Hai món cùng tên khiến phần tử gặp sau ghi đè phần tử trước; slice vẫn có thể hiển thị cả hai. Map lưu `FoodItem` theo value nên sửa `menu[i].Price` không tự sửa map. Nếu cần menu động, có thể dùng ID duy nhất, xây lại index sau mỗi cập nhật hoặc lưu index/vị trí vào slice thay cho bản sao. Phải xác định một nơi cập nhật dữ liệu; không nên chỉ sửa hiển thị và bỏ quên tra cứu.

### Câu 8. `OrderLine.Item` là bản sao hay con trỏ tới món trong menu? Giá trong order có thay đổi nếu menu thay giá không?

`SearchFoodByName` lấy value từ map; `OrderLine{Item: item}` lưu một bản sao `FoodItem`. Code không dùng `*FoodItem` và không tìm lại giá khi tính subtotal, nên giá đã lưu trong dòng order không tự thay theo menu. Điều này có thể xem như giữ giá tại thời điểm thêm món. Đổi sang con trỏ có thể làm đơn cũ thay tổng khi menu đổi giá, nên phải quyết định chính sách trước khi đổi kiểu dữ liệu.

### Câu 9. Vì sao khai báo `type CustomerType string` và `type OrderType string`? Chúng có bảo đảm chỉ có đúng bốn giá trị được phép không?

Đây là hai named type khác nhau: biến `OrderType` không thể truyền trực tiếp vào tham số `CustomerType` nếu không chuyển kiểu. Các constant có tên làm lời gọi dễ đọc và dùng `%s` để hiển thị được. Nhưng chúng không phải enum đóng: literal `"VIP"` có thể được gán vào `CustomerType`, hoặc viết `CustomerType("VIP")`. Code CLI kiểm tra Regular/Member và Pickup/Delivery; các hàm tính toán chỉ so sánh với Member/Delivery, không báo lỗi cho giá trị khác. Nếu thêm API, cần validation tại đầu vào model/API.

## 2. Tra cứu, validation và cập nhật order

### Câu 10. Tại sao `SearchFoodByName` trả `(FoodItem, bool)`? Nếu bỏ `ok` thì món không tồn tại được hiểu thành gì?

Tra map không có key trả zero value của `FoodItem`: tên rỗng, giá 0, category rỗng, available false. `ok` phân biệt không có key với một value có field bằng zero. `AddFoodToOrder` kiểm tra `!ok` trước khi đọc availability để báo đúng `food not found`, thay vì nhầm thành `food unavailable`. Đây là pattern tra cứu map phù hợp khi không tìm thấy là trường hợp bình thường của luồng nhập tên.

### Câu 11. Vì sao search có thể tìm thấy món unavailable, nhưng danh sách available lại không hiển thị nó?

Hai thao tác trả lời hai câu hỏi khác nhau. `DisplayAvailableFood` lọc để cho biết đang bán được món gì. `SearchFoodByName` tra thông tin món trong toàn menu và trả trạng thái `Available`, giúp phân biệt tồn tại nhưng hết hàng với không tồn tại. Quyết định có được mua hay không thuộc `AddFoodToOrder`, không thuộc search. Nếu loại món unavailable khỏi map, hệ thống sẽ mất khả năng báo đúng nguyên nhân hết hàng.

### Câu 12. Vì sao `AddFoodToOrder` kiểm tra tên → availability → quantity trước khi append? Nếu cả tên và quantity đều sai, lỗi nào được trả?

Hàm return ngay khi điều kiện đầu tiên thất bại. Tên không tồn tại trả `food not found`, kể cả quantity cũng bằng 0; tên tồn tại nhưng unavailable trả `food unavailable` trước lỗi quantity. Đây là lựa chọn thứ tự ưu tiên thông báo, không phải thứ tự duy nhất đúng. Tất cả kiểm tra diễn ra trước `append`, nên lỗi không làm thêm dòng vào order. Cách return sớm giữ đường đi thành công rõ ràng và tránh phải undo thao tác đã thực hiện.

### Câu 13. Tại sao validation nghiệp vụ nằm trong `AddFoodToOrder`, trong khi CLI vẫn gọi `strconv.Atoi`? Hai bước này có trùng nhau không?

`Atoi` kiểm tra biểu diễn input: `"abc"`, `"1.5"` hay số vượt miền `int` không chuyển được. `IsValidQuantity` kiểm tra ý nghĩa: `"0"` và `"-3"` chuyển được thành số nguyên nhưng không được phép đặt. CLI chịu trách nhiệm chuyển text; hàm thêm món chịu trách nhiệm business rule để người gọi từ test hay giao diện khác cũng được kiểm tra. Không nên chỉ chặn số âm trong CLI vì lời gọi trực tiếp sẽ bỏ qua rule đó.

### Câu 14. Vì sao `AddFoodToOrder` nhận `*Order`, còn `CalculateSubtotal` nhận `Order`? Đặc biệt, slice vốn tham chiếu mảng thì có cần con trỏ không?

Thêm món phải cập nhật field `Lines` của order ở nơi gọi. Slice gồm một header chứa con trỏ mảng, length và capacity; truyền `Order` theo value sao chép header này. `append` có thể cấp phát mảng mới và luôn trả header mới với length mới. Nếu chỉ gán `order.Lines = append(...)` trên bản sao `Order`, header trong order của caller không được cập nhật. `*Order` cho phép cập nhật field thật. Hàm subtotal chỉ đọc dữ liệu nên không cần con trỏ. Truyền `Order` theo value cũng không sao chép sâu toàn bộ các dòng; hàm đọc hiện tại không sửa mảng dùng chung.

### Câu 15. Tại sao trả `error` thay vì in lỗi ngay trong `AddFoodToOrder`? `nil` và `fmt.Errorf` được caller hiểu như thế nào?

`nil` báo thành công; `fmt.Errorf` tạo lỗi với tên món/quantity để giải thích nguyên nhân. CLI quyết định in `Error: ...`, còn test kiểm tra có lỗi hay không mà không cần bắt stdout. Điều này tách nghiệp vụ khỏi giao diện. Khi ứng dụng lớn hơn, có thể dùng sentinel error hoặc kiểu lỗi riêng để phân biệt nguyên nhân bằng `errors.Is`/`errors.As`, thay vì so sánh text. Lab hiện tại chưa có yêu cầu đó.

### Câu 16. Nếu thêm Milk Tea ×2 rồi Milk Tea ×3 thì có gộp dòng không? Tổng và summary ra sao?

Hàm luôn append một `OrderLine` mới nên tạo hai dòng, không tìm dòng trùng để gộp. Subtotal vẫn đúng: 25,000 ×2 + 25,000 ×3 = 125,000. Summary in hai dòng tương ứng. Đề không bắt buộc gộp món trùng. Nếu được yêu cầu gộp, phải xác định cùng tên nhưng khác giá có được gộp không; không nên chỉ cộng quantity khi hai dòng lưu giá khác nhau.

### Câu 17. Search có bỏ qua chữ hoa/thường và khoảng trắng không? Vì sao gọi hàm trực tiếp có thể khác nhập qua CLI?

Map tra đúng string key nên `"milk tea"` không khớp `"Milk Tea"`. CLI dùng `strings.TrimSpace` trước khi gọi, nên loại khoảng trắng đầu/cuối; `SearchFoodByName` tự nó không trim. Vì vậy gọi trực tiếp `SearchFoodByName(" Milk Tea ")` không thấy món. Nếu thêm tìm không phân biệt hoa thường, phải chuẩn hóa cả key lúc build map lẫn input lúc tra, và xử lý tên trùng sau chuẩn hóa. Chỉ chuyển input sang chữ thường sẽ không khớp map đang có key viết hoa.

### Câu 18. Các hàm có bảo vệ được `nil *Order`, quantity cực lớn hoặc một `OrderLine` được tạo trực tiếp với giá âm không?

Không bảo vệ đầy đủ. `AddFoodToOrder(nil, "Milk Tea", 1)` qua validation rồi dereference `order.Lines` sẽ panic. Quantity chỉ cần >0, không có giới hạn tối đa nghiệp vụ. Field export cho phép tạo dòng có giá/quantity âm rồi gọi subtotal trực tiếp, bỏ qua hàm thêm món. Với menu cố định và CLI khởi tạo order thật, đường đi bình thường không tạo các trạng thái này. Nếu tái sử dụng như API công khai, cần constructor/validation và giới hạn quantity, đồng thời quyết định có cho phép caller sửa field trực tiếp không.

## 3. Subtotal, discount, delivery fee và summary

### Câu 19. Trace `CalculateSubtotal` cho Chicken Rice ×2 và Milk Tea ×3. Tại sao phải có `float64(line.Quantity)`?

Hàm khởi tạo `total = 0`, duyệt từng dòng và cộng giá nhân số lượng. Chicken Rice đóng góp 45,000 ×2 = 90,000; Milk Tea đóng góp 25,000 ×3 = 75,000; subtotal là 165,000. Go không tự nhân một `float64` với biến `int`, nên chuyển quantity sang `float64`. Nếu bỏ phép nhân quantity, subtotal chỉ còn 70,000, tương đương tính một phần mỗi dòng. Hàm có O(k) thời gian với k dòng và O(1) bộ nhớ bổ sung.

### Câu 20. Vì sao tách `CalculateSubtotal`, `CalculateDiscount`, `CalculateCustomerDiscount`, `CalculateDeliveryFee` thay vì viết tất cả trong `DisplayOrderSummary`?

Subtotal mô tả tiền món; hai discount mô tả chính sách giảm giá; fee mô tả cách nhận. Tách chúng cho phép kiểm tra từng rule trực tiếp bằng test và thay đổi rate/fee mà không sửa cách duyệt món hay in. `CalculateFinalTotal` ghép kết quả theo chính sách; summary dùng các hàm để trình bày. Với chỉ vài rule, các function nhỏ đủ rõ, chưa cần interface/policy engine. Cần lưu ý summary vẫn tính subtotal để in rồi `CalculateFinalTotal` tính lại; tách trách nhiệm không có nghĩa mọi kết quả đã được cache.

### Câu 21. Tại sao volume discount dùng `>=`, không dùng `>`? Test nào chứng minh đúng tại ranh giới?

ADAPT yêu cầu subtotal ít nhất 300,000 nên đúng 300,000 cũng được giảm. `CalculateDiscount(300000)` trả 30,000. `TestDiscountAtThreshold` và bốn subtest có hậu tố `at` kiểm tra ranh giới này. Trong review ngày 06/10/2026, đổi `>=` thành `>` ở bản sao tạm làm test discount và bốn subtest tại ngưỡng FAIL. Điều này chứng minh suite bắt được lỗi cụ thể đó, không chứng minh suite bắt mọi lỗi.

### Câu 22. Regular có phải luôn không được giảm giá không? Phân biệt `Discount` với `Member Discount` trong summary.

Regular chỉ không có customer discount. `Discount` trong summary gọi `CalculateDiscount` và là volume discount từ ADAPT, áp cho cả Regular lẫn Member khi subtotal đạt ngưỡng. `Member Discount` gọi `CalculateCustomerDiscount` và chỉ áp nếu `Customer == Member`. Vì vậy Regular + Pickup ở 300,000 trả 270,000. Nói “Regular luôn trả nguyên subtotal” sẽ sai với chính sách hiện tại.

### Câu 23. Tại sao hai discount cùng dùng subtotal gốc? So sánh cộng dồn với áp tuần tự trên subtotal 350,000.

Code chủ động chọn cộng dồn vì cả hai rule được mô tả trên food subtotal. Với S=350,000, volume discount=35,000 và Member discount=17,500; Pickup trả 297,500. Nếu lấy 5% sau giảm 10%, Member discount chỉ 15,750 và trả 299,250. Hai cách khác nhau 1,750 VND. Code hiện tại truyền cùng `subtotal` cho hai hàm nên thứ tự gọi không thay đổi kết quả. Đây là quyết định dự án đã ghi trong `HOMEWORK.md`, không phải quy tắc kết hợp được PDF nói rõ.

### Câu 24. Công thức hiện tại có trùng nguyên văn công thức homework trong PDF không? Nếu giảng viên yêu cầu bỏ volume discount thì thay đổi ở đâu?

Không trùng nguyên văn: PDF ghi `Subtotal - Customer Discount + Delivery Fee`; code còn trừ `volumeDiscount`. PDF không làm rõ số phận rule ADAPT, còn `CLAUDE.md` yêu cầu chọn và giải thích cách kết hợp. Dự án chọn giữ cả hai. Nếu giảng viên xác định homework thay thế rule ADAPT, cần bỏ volume discount khỏi `CalculateFinalTotal`, điều chỉnh summary để không hiển thị khoản giảm không áp dụng, cập nhật expected value trong test tổ hợp và tài liệu. Ví dụ Member + Delivery ở 300,000 đổi từ 285,000 thành 315,000.

### Câu 25. Trace đầy đủ Member + Delivery với Chicken Rice ×7 và Milk Tea ×2. Chỗ nào bảo đảm không giảm phí giao hàng?

Subtotal = 45,000 ×7 + 25,000 ×2 = 365,000. Volume discount = 36,500; Member discount = 18,250; delivery fee = 30,000. Final = 365,000 −36,500 −18,250 +30,000 = 340,250. Hàm discount chỉ nhận subtotal chưa cộng fee; `CalculateDeliveryFee` được cộng như một khoản riêng. Việc tách cơ sở tính discount khỏi fee mới bảo đảm fee không bị giảm, không chỉ việc đặt dấu `+ fee` ở cuối biểu thức.

### Câu 26. Nếu tính Member discount trên `(subtotal + deliveryFee)` thì sai bao nhiêu tiền? Test nào phát hiện?

5% của fee 30,000 là 1,500, nên Member + Delivery sẽ bị giảm thêm 1,500 không đúng yêu cầu. Trong ví dụ trên, final sai thành 338,750. `TestMemberDiscountDoesNotReduceDeliveryFee` so sánh hai order cùng subtotal/customer, chỉ khác Pickup/Delivery: gap phải đúng 30,000. Nếu fee bị giảm 5%, gap chỉ còn 28,500; test FAIL. Ma trận 12 expected total cũng kiểm tra các giá trị cụ thể.

### Câu 27. Order rỗng trả bao nhiêu? `Order{}` có giống `Order{Customer: Regular, Type: Pickup}` hoàn toàn không?

Slice rỗng hoặc nil không có dòng để cộng nên subtotal=0; volume/member discount trên 0 đều bằng 0. Pickup final=0; Delivery final=30,000 vì code không kiểm tra phải có món trước khi thu phí. `Order{}` có `Customer` và `Type` bằng chuỗi rỗng, được các hàm tính tiền xử lý như discount 0/fee 0 nên tổng tương tự Regular + Pickup. Nhưng summary in nhãn rỗng, không phải Regular/Pickup. CLI khởi tạo tường minh hai field, không dựa vào zero value để gán mặc định.

### Câu 28. Đổi customer type sau khi thêm món có cần sửa từng dòng hay lưu lại total không? Dữ liệu có mất khi thoát không?

Các dòng chỉ lưu món và quantity. Đổi `order.Customer` hoặc `order.Type` không thay subtotal; mỗi lần xem summary, các khoản discount/fee được tính lại từ field hiện tại. Không có field tổng được cache nên không phải đồng bộ total khi đổi loại đơn. Order chỉ nằm trong biến `order` của tiến trình, không lưu file/database; thoát mất dữ liệu và lần chạy mới tạo Regular + Pickup rỗng.

### Câu 29. Summary có tính subtotal hai lần không? Có nên tối ưu ngay bằng cách lưu total trong `Order`?

Có. Summary gọi `CalculateSubtotal` để in, sau đó `CalculateFinalTotal` gọi lại. Hai lượt O(k) vẫn là O(k) về độ phức tạp và nhỏ với lab. Lưu total trong model tạo thêm nghĩa vụ cập nhật khi thêm/xóa/đổi quantity hoặc giá, dễ sinh dữ liệu cũ. Nếu cần một breakdown thống nhất và giảm lặp, có thể thêm hàm trả struct gồm Subtotal, VolumeDiscount, CustomerDiscount, DeliveryFee, FinalTotal, rồi summary và test dùng kết quả đó. Đây là hướng mở rộng, chưa implement.

## 4. Function, method và luồng CLI

### Câu 30. Đề nhắc function/method: chương trình hiện tại có method không? Vì sao dùng function thay vì `order.AddFood(...)` và `order.Subtotal()`?

Code hiện tại dùng function, không có method: không khai báo receiver như `func (o *Order) AddFood(...)`. Đề cho phép functions và/hoặc methods. Function nhận tham số rõ ràng, dễ gọi và đủ đơn giản cho bài Go foundation. Method có thể gom hành vi quanh `Order`, ví dụ `func (o Order) Subtotal() float64`, nhưng không tự bảo đảm đúng nghiệp vụ và không tự loại bỏ phụ thuộc global menu. Nếu chuyển sang method, vẫn phải giữ validation, semantics con trỏ và policy discount như cũ. Chọn kiểu nào nên theo mức dễ hiểu của model, không chỉ vì cú pháp gọi ngắn hơn.

### Câu 31. Vì sao dùng `bufio.Reader.ReadString('\n')` thay vì đọc tên món bằng token? `TrimSpace`, `switch`, `continue` làm gì?

Tên như Chicken Rice có khoảng trắng; đọc theo dòng lấy được cả tên. `TrimSpace` bỏ newline và khoảng trắng ngoài cùng, giữ khoảng trắng giữa tên. `switch choice` chọn hành động 1-7; nhánh 3 dùng `continue` khi `Atoi` thất bại để quay lại menu mà không thêm món. Nhánh 7 `return` kết thúc `main`. Đây là điều phối giao diện; tính tiền và validation món nằm trong các function được gọi.

### Câu 32. Nhập `member` hoặc `Express` có làm order đổi loại không? Nếu stdin hết dữ liệu, vòng lặp hiện tại phản ứng thế nào?

CLI chỉ nhận đúng Regular/Member, Pickup/Delivery sau trim; input khác in lỗi và không gán field nên giữ loại cũ. Nhưng mọi `ReadString` đang bỏ qua error bằng `_`. Khi EOF và không có lựa chọn 7, choice thường thành chuỗi rỗng, chạy default `Invalid option`, rồi lặp tiếp; có thể in liên tục. Cải thiện cần xử lý `io.EOF`/lỗi đọc tại các prompt và thoát hợp lý, đồng thời quyết định xử lý dòng cuối không có newline. Test tính tiền hiện tại không kiểm tra đường đi EOF.

## 5. Kiểm thử, AI review và challenge thay đổi yêu cầu

### Câu 33. Vì sao homework dùng table-driven test với `t.Run` thay vì viết 12 function gần giống nhau? Expected total có nên được tính bằng chính các hàm đang test không?

Mỗi row chứa tên case, subtotal, customer, order type và `want`. Vòng lặp dựng order rồi gọi `CalculateFinalTotal`; `t.Run` gắn tên riêng để thấy tổ hợp nào lỗi. Cách này giảm lặp phần setup/assertion mà vẫn liệt kê đủ 4×3 case. `want` là số cố định đã tính độc lập; không nên tạo expected bằng cách gọi lại `CalculateDiscount`/`CalculateCustomerDiscount`, vì bug chung có thể làm cả actual và expected cùng sai mà test PASS.

### Câu 34. Helper `orderWorth` dùng Iced Coffee 20,000 để dựng subtotal. Vì sao có `t.Helper()` và bước kiểm tra subtotal? Dùng helper này cho 299,999 được không?

Helper lấy `int(subtotal / 20000)` làm quantity, thêm Coffee, rồi kiểm tra subtotal thực tế. `t.Helper()` giúp lỗi assertion trong helper được báo tại nơi gọi hữu ích hơn. `t.Fatalf` dừng test nếu setup sai, tránh chạy assertion final trên order không đúng. Các mốc hiện có 200,000/300,000/400,000 đều chia hết 20,000. Với 299,999, quantity bị cắt còn 14 và subtotal=280,000 nên helper FAIL. Muốn kiểm tra sát biên, có thể test `CalculateDiscount` trực tiếp hoặc dùng dòng món giả lập với giá phù hợp và nêu rõ mục tiêu test.

### Câu 35. Vì sao vừa cần test discount riêng, test final theo tổ hợp, vừa cần test delivery gap? Chỉ test 300,000 Regular + Delivery có đủ không?

Test riêng xác định rule discount/fee sai ở đâu. Test tổ hợp kiểm tra chúng được ghép đúng khi có nhiều điều kiện. Test gap diễn đạt trực tiếp invariant “delivery tăng đúng 30,000 dù khách nào”. Riêng 300,000 Regular + Delivery có volume discount 30,000 và fee 30,000 triệt tiêu, nên tổng vẫn 300,000; implementation bỏ cả hai khoản có thể tình cờ PASS case đó. Những case khác và test từng rule mới loại bỏ kiểu sai này.

### Câu 36. Toàn bộ test PASS có chứng minh mọi yêu cầu và mọi edge case đều đúng không? Đếm test hiện tại như thế nào cho chính xác?

PASS xác nhận các assertion đã viết trên các input đã chạy. Hiện có 14 test function cấp cao nhất và 12 subtest nằm trong `TestFinalTotalCombinations`; không nên gọi mơ hồ “25 function”. Suite kiểm tra validation chính, subtotal, ba mốc discount và tổ hợp homework. Chưa tự động kiểm tra output summary, input customer/order sai, EOF, nil pointer, menu động hay mọi boundary sát ngưỡng. Một số test reject chỉ kiểm tra error; chỉ test not-found kiểm tra thêm order không đổi. Nếu mở rộng, ưu tiên test hành vi có nguy cơ thay đổi thay vì test sao chép y hệt code.

### Câu 37. Nếu giảng viên cố ý bỏ nhân quantity, hoặc đổi `>=` thành `>`, nên identify → explain → fix → retest thế nào?

Chạy test và xem actual/want. Bỏ nhân quantity khiến test Chicken Rice ×2 nhận 45,000 thay 90,000; nguyên nhân là chỉ cộng giá mỗi dòng. Đổi điều kiện ngưỡng làm case đúng 300,000 mất discount. Khôi phục đúng biểu thức theo yêu cầu, chạy lại test liên quan rồi suite. Không sửa expected value chỉ để test xanh nếu yêu cầu chưa đổi. Review hiện tại đã thử lỗi ngưỡng trên bản sao tạm; source chính không bị thay.

### Câu 38. AI review `CalculateFinalTotal` nên kết luận gì mà không rewrite? Làm sao phân biệt lỗi với lựa chọn thiết kế?

Logic đúng theo chính sách cộng dồn đã chọn: cùng subtotal gốc cho hai discount, fee cộng riêng. Điều chưa rõ ở đặc tả PDF là có giữ ADAPT; cần ghi lựa chọn, không gọi nó là quy tắc bắt buộc từ PDF. Các giới hạn gồm giá trị enum không hợp lệ bị coi như không giảm/không phí, order được tạo trực tiếp có thể bypass validation, và float chưa có chính sách làm tròn. Review nên nêu đầu vào gây vấn đề và phạm vi ảnh hưởng; chưa cần rewrite function chỉ vì có thể thiết kế khác. Nhóm phải tự so sánh review với lý luận của mình, không coi tài liệu AI là bằng chứng cả hai người đã hiểu.

### Câu 39. Nếu được yêu cầu thêm search by category, nên dùng slice hay map hiện có? Có lọc unavailable luôn không?

Map hiện tại index theo tên nên không tra category trực tiếp được. Với menu nhỏ, duyệt slice O(n) và trả slice các món khớp category là đủ. Phải hỏi/đọc rõ yêu cầu có chỉ lấy món available hay lấy tất cả; hai policy không giống nhau. Ví dụ trả tất cả món đúng category, giữ thứ tự menu:

```go
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

Không có match trả nil slice, `len(result)==0`; có thể dùng trong `range` bình thường. Nếu cần tìm category lặp lại trên menu rất lớn, cân nhắc index `map[string][]FoodItem` cùng chính sách đồng bộ.

### Câu 40. Nếu đổi volume rule thành 15% từ 500,000, hoặc thêm rule mới, điểm nào cần sửa và test nào cần cập nhật?

Đổi rate/ngưỡng của một rule: sửa `DiscountRate`/`DiscountThreshold`, cập nhật test boundary và expected final ở các tổ hợp, summary/documentation theo chính sách mới. `CalculateDiscount` vẫn dùng các constant nên không cần đổi cấu trúc hàm. Thêm rule mới khác loại cần quyết định cơ sở tính, điều kiện, cộng dồn/loại trừ và thứ tự; không chỉ thêm một constant. Với hai rule hiện tại, gọi tên từng hàm trong `CalculateFinalTotal` rõ ràng. Nếu số rule tăng và thay đổi thường xuyên, mới cân nhắc danh sách rule hoặc pricing policy; test vẫn phải giữ expected độc lập và kiểm tra fee không bị giảm.

### Câu 41. Homework được hấp thụ bằng thiết kế ban đầu như thế nào? Những phần nào thực sự phải thay đổi?

`OrderLine` đã giữ món và quantity nên không cần đổi cách cộng subtotal. Homework thêm hai named type, hai field order, các constant và ba hàm `CalculateCustomerDiscount`, `CalculateDeliveryFee`, `CalculateFinalTotal`; dùng lại subtotal và ADAPT discount. Summary phải hiển thị các khoản mới, CLI phải cho chọn loại khách/loại đơn, và test mở rộng ma trận. Vì vậy nói “chỉ thêm, không sửa hàm cũ nào” là quá rộng: ít nhất giao diện hiển thị và điều phối phải tích hợp thay đổi. Thư mục không có Git history để chứng minh lịch sử thay từng dòng; ta chỉ xác nhận code hiện tại và kết quả test hiện tại.

### Câu 42. Khi nộp bài, đâu là bằng chứng đủ chức năng và đâu là phần cả hai thành viên phải tự chứng minh?

Source + test chứng minh hiện thực và các case đã chạy. `HOMEWORK.md` có design note, AI review và giải thích thay đổi; README ghi checklist và policy discount. Nhưng không file nào tự chứng minh nhóm làm trong giới hạn thời gian, THINK đã được thực hiện AI OFF hay cả hai người hiểu bài. Mỗi thành viên nên tự trace Member + Delivery, giải thích slice header/con trỏ, tìm một edge case và sửa một rule kèm expected test. Đây là mục tiêu “Own the result” của đề.
