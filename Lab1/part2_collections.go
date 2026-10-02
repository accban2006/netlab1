package main

import (
	"fmt"
	"sort"
	"strings"
)

// Part 2 — Arrays, Slices and Maps.

// findPropertiesInBudget trả về các property có giá <= maxBudget.
// Trả nil khi không khớp gì: len(nil slice) == 0 và range trên nil là an toàn.
func findPropertiesInBudget(properties []Property, maxBudget float64) []Property {
	var result []Property
	for _, prop := range properties {
		if prop.IsAffordable(maxBudget) {
			result = append(result, prop)
		}
	}
	return result
}

// findPropertiesByBedrooms trả về các property có đúng số phòng ngủ yêu cầu.
func findPropertiesByBedrooms(properties []Property, bedrooms int) []Property {
	var result []Property
	for _, prop := range properties {
		if prop.Bedrooms == bedrooms {
			result = append(result, prop)
		}
	}
	return result
}

// findPropertiesByDistrict lọc theo quận, so sánh không phân biệt hoa thường và
// bỏ khoảng trắng đầu/cuối để chịu được input gõ tay từ menu.
func findPropertiesByDistrict(properties []Property, district string) []Property {
	needle := strings.ToLower(strings.TrimSpace(district))

	var result []Property
	for _, prop := range properties {
		if strings.ToLower(strings.TrimSpace(prop.District)) == needle {
			result = append(result, prop)
		}
	}
	return result
}

// analyzeByDistrict group property theo quận. Lấy nguyên từ PDF.
func analyzeByDistrict(properties []Property) map[string][]Property {
	districtMap := make(map[string][]Property)
	for _, prop := range properties {
		districtMap[prop.District] = append(districtMap[prop.District], prop)
	}
	return districtMap
}

// calculateDistrictStats trả về slice (không phải map) đã sort giảm dần theo
// giá trung bình — yêu cầu của PDF là "sorted by average price (highest first)",
// mà map trong Go không giữ thứ tự.
func calculateDistrictStats(properties []Property) []DistrictStats {
	districtMap := analyzeByDistrict(properties)

	var stats []DistrictStats
	for district, props := range districtMap {
		total := 0.0
		mostExpensive := props[0]

		for _, prop := range props {
			total += prop.Price
			if prop.Price > mostExpensive.Price {
				mostExpensive = prop
			}
		}

		stats = append(stats, DistrictStats{
			District:      district,
			Count:         len(props),
			AveragePrice:  total / float64(len(props)),
			MostExpensive: mostExpensive,
		})
	}

	// sort.Slice không stable, nên tie-break bằng tên quận để output ổn định
	// giữa các lần chạy (thứ tự range map là ngẫu nhiên).
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].AveragePrice != stats[j].AveragePrice {
			return stats[i].AveragePrice > stats[j].AveragePrice
		}
		return stats[i].District < stats[j].District
	})

	return stats
}

// runTask21 demo các hàm search.
func runTask21() {
	printHeader("All Properties")
	printPropertyTable(properties)

	budget := 3000000000.0
	affordable := findPropertiesInBudget(properties, budget)
	fmt.Printf("\nProperties under %s VND (%s):\n", formatVND(budget), formatPrice(budget))
	printPropertyTable(affordable)

	bedrooms := 3
	byBedrooms := findPropertiesByBedrooms(properties, bedrooms)
	fmt.Printf("\nProperties with %d bedrooms:\n", bedrooms)
	printPropertyTable(byBedrooms)

	district := "District 1"
	byDistrict := findPropertiesByDistrict(properties, district)
	fmt.Printf("\nProperties in %s:\n", district)
	printPropertyTable(byDistrict)

	// Case không khớp gì — in thông báo thay vì bảng rỗng.
	tiny := findPropertiesInBudget(properties, 1000000)
	fmt.Printf("\nProperties under %s VND:\n", formatVND(1000000))
	printPropertyTable(tiny)
}

// runTask22 in phân tích theo quận và bảng xếp hạng.
func runTask22() {
	districtMap := analyzeByDistrict(properties)

	// Sort keys trước khi in: range map có thứ tự ngẫu nhiên.
	keys := make([]string, 0, len(districtMap))
	for district := range districtMap {
		keys = append(keys, district)
	}
	sort.Strings(keys)

	printHeader("District Analysis")
	for _, district := range keys {
		props := districtMap[district]

		total := 0.0
		mostExpensive := props[0]
		for _, prop := range props {
			total += prop.Price
			if prop.Price > mostExpensive.Price {
				mostExpensive = prop
			}
		}

		fmt.Printf("%s: %d properties, Avg: %s, Most expensive: %s\n",
			district, len(props), formatPrice(total/float64(len(props))), mostExpensive.Name)
	}

	fmt.Println("\nRanking by Average Price:")
	for i, s := range calculateDistrictStats(properties) {
		fmt.Printf("%d. %s: %s\n", i+1, s.District, formatPrice(s.AveragePrice))
	}
}

// runPart2 chạy cả hai task của Part 2.
func runPart2() {
	runTask21()
	runTask22()
}
