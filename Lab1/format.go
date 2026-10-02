package main

import (
	"fmt"
	"strings"
)

// formatPrice quy đổi số VND sang đơn vị đọc được bằng tiếng Việt.
//
//	2500000000 -> "2.5 tỷ VND"
//	 850000000 -> "850 triệu VND"
func formatPrice(vnd float64) string {
	switch {
	case vnd >= 1e9:
		return fmt.Sprintf("%.1f tỷ VND", vnd/1e9)
	case vnd >= 1e6:
		return fmt.Sprintf("%.0f triệu VND", vnd/1e6)
	case vnd >= 1e3:
		return fmt.Sprintf("%.0f nghìn VND", vnd/1e3)
	default:
		return fmt.Sprintf("%.0f VND", vnd)
	}
}

// formatPriceDetailed như formatPrice nhưng giữ thêm chữ số thập phân, dùng cho
// các số cần độ chính xác như khoản trả hàng tháng hay tổng lãi.
//
//	2166351654 -> "2.17 tỷ VND"   (formatPrice ra "2.2 tỷ VND")
//	  17356465 -> "17.4 triệu VND" (formatPrice ra "17 triệu VND")
func formatPriceDetailed(vnd float64) string {
	switch {
	case vnd >= 1e9:
		return fmt.Sprintf("%.2f tỷ VND", vnd/1e9)
	case vnd >= 1e6:
		return fmt.Sprintf("%.1f triệu VND", vnd/1e6)
	case vnd >= 1e3:
		return fmt.Sprintf("%.1f nghìn VND", vnd/1e3)
	default:
		return fmt.Sprintf("%.0f VND", vnd)
	}
}

// formatVND in số nguyên VND có dấu phẩy phân cách nghìn: 33112582 ->
// "33,112,582". Làm tròn xuống về số nguyên trước khi chèn dấu phẩy.
func formatVND(vnd float64) string {
	digits := fmt.Sprintf("%.0f", vnd)

	neg := strings.HasPrefix(digits, "-")
	if neg {
		digits = digits[1:]
	}

	var b strings.Builder
	for i, ch := range digits {
		// Chèn dấu phẩy mỗi khi số chữ số còn lại chia hết cho 3.
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}

	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// printHeader in tiêu đề một section theo đúng style của PDF.
func printHeader(title string) {
	fmt.Printf("\n=== %s ===\n", title)
}

// printPropertyTable in danh sách property dạng bảng, hoặc thông báo rỗng.
func printPropertyTable(properties []Property) {
	if len(properties) == 0 {
		fmt.Println("No properties found")
		return
	}

	fmt.Printf("%-4s %-20s %-16s %8s %5s %-12s %-18s\n",
		"#", "Name", "Price", "Area", "BR", "District", "Price/m²")
	fmt.Println(strings.Repeat("-", 88))

	for i, prop := range properties {
		fmt.Printf("%-4d %-20s %-16s %8.1f %5d %-12s %-18s\n",
			i+1, prop.Name, formatPrice(prop.Price), prop.Area,
			prop.Bedrooms, prop.District, formatVND(prop.PricePerM2())+" VND")
	}
}
