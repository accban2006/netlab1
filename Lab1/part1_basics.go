package main

import "fmt"

// Part 1 — Variables and Basic Types.
//
// Phần này cố ý dùng biến rời (property1Name, property2Price...) thay vì slice,
// vì mục tiêu của PDF ở Part 1 là luyện khai báo biến với kiểu tường minh.
// Slice/map là nội dung của Part 2.

// categorizeProperty xếp hạng property theo giá/m². Lấy nguyên từ PDF — chú ý
// biên là `>` chứ không phải `>=`, nên đúng 50,000,000 vẫn là PREMIUM.
func categorizeProperty(pricePerM2 float64) string {
	if pricePerM2 > 50000000 {
		return "LUXURY"
	} else if pricePerM2 > 30000000 {
		return "PREMIUM"
	} else if pricePerM2 > 20000000 {
		return "STANDARD"
	}
	return "BUDGET"
}

// safePricePerM2 chia giá cho diện tích, guard area == 0 (chia 0 với float64
// không panic mà ra +Inf, output sẽ rất khó đọc).
func safePricePerM2(price, area float64) float64 {
	if area == 0 {
		return 0
	}
	return price / area
}

// runTask11 so sánh 3 property và tìm cái có giá/m² thấp nhất.
func runTask11() {
	// Property 1
	var property1Name string = "Saigon Apartment"
	var property1Price float64 = 2500000000
	var property1Area float64 = 75.5
	var property1Bedrooms int = 2
	var property1Available bool = true

	// Property 2
	var property2Name string = "HCMC House"
	var property2Price float64 = 4200000000
	var property2Area float64 = 120.0

	// Property 3
	var property3Name string = "Budget Studio"
	var property3Price float64 = 800000000
	var property3Area float64 = 35.0

	// Giá/m² — dùng := vì kiểu suy ra được từ biểu thức.
	property1PricePerM2 := safePricePerM2(property1Price, property1Area)
	property2PricePerM2 := safePricePerM2(property2Price, property2Area)
	property3PricePerM2 := safePricePerM2(property3Price, property3Area)

	printHeader("Property Information")
	fmt.Printf("Name: %s\n", property1Name)
	fmt.Printf("Price: %s VND\n", formatVND(property1Price))
	fmt.Printf("Area: %.1f m²\n", property1Area)
	fmt.Printf("Bedrooms: %d\n", property1Bedrooms)
	fmt.Printf("Available: %t\n", property1Available)
	fmt.Printf("Price per m²: %s VND\n", formatVND(property1PricePerM2))

	printHeader("Property Comparison")
	fmt.Printf("Property 1: %s - %s VND/m²\n", property1Name, formatVND(property1PricePerM2))
	fmt.Printf("Property 2: %s - %s VND/m²\n", property2Name, formatVND(property2PricePerM2))
	fmt.Printf("Property 3: %s - %s VND/m²\n", property3Name, formatVND(property3PricePerM2))

	// Tìm min bằng chuỗi if: khởi tạo từ property 1 rồi so dần.
	cheapestPrice := property1PricePerM2
	cheapestName := property1Name

	if property2PricePerM2 < cheapestPrice {
		cheapestPrice = property2PricePerM2
		cheapestName = property2Name
	}
	if property3PricePerM2 < cheapestPrice {
		cheapestPrice = property3PricePerM2
		cheapestName = property3Name
	}

	fmt.Printf("\nCheapest per m²: %s at %s VND/m²\n", cheapestName, formatVND(cheapestPrice))
}

// categoryOrder cố định thứ tự in category. Range trên map trong Go có thứ tự
// ngẫu nhiên, nên phải có slice thứ tự riêng để output ổn định.
var categoryOrder = []string{"LUXURY", "PREMIUM", "STANDARD", "BUDGET"}

// runTask12 xếp hạng từng property và đếm số lượng mỗi hạng.
func runTask12() {
	var property1Name string = "Saigon Apartment"
	var property1Price float64 = 2500000000
	var property1Area float64 = 75.5

	var property2Name string = "HCMC House"
	var property2Price float64 = 4200000000
	var property2Area float64 = 120.0

	var property3Name string = "Budget Studio"
	var property3Price float64 = 800000000
	var property3Area float64 = 35.0

	names := []string{property1Name, property2Name, property3Name}
	prices := []float64{property1Price, property2Price, property3Price}
	areas := []float64{property1Area, property2Area, property3Area}

	counts := make(map[string]int)

	printHeader("Property Categories")
	for i := range names {
		pricePerM2 := safePricePerM2(prices[i], areas[i])
		category := categorizeProperty(pricePerM2)
		counts[category]++

		fmt.Printf("%s: %s (%s)\n", names[i], category, formatPrice(prices[i]))
	}

	// Lưu ý chênh lệch với PDF: bảng output mẫu ghi Saigon Apartment là
	// STANDARD, nhưng giá/m² của nó là 33,112,582 > 30,000,000 nên theo đúng
	// categorizeProperty() của PDF phải ra PREMIUM. Làm theo hàm, không theo bảng.

	fmt.Println("\nCategory Summary:")
	for _, category := range categoryOrder {
		fmt.Printf("%s: %d properties\n", category, counts[category])
	}
}

// runPart1 chạy cả hai task của Part 1.
func runPart1() {
	runTask11()
	runTask12()
}
