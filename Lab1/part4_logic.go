package main

import (
	"fmt"
	"sort"
	"strings"
)

// Part 4 — Control Flow and Logic.

// assumedMonthlyRentRate là giả định thuê 1.2% giá trị/tháng, theo PDF.
const assumedMonthlyRentRate = 0.012

// recommendProperty là bản base case lấy nguyên từ PDF.
func recommendProperty(p Property, budget, maxMonthlyPayment float64) string {
	if !p.IsAffordable(budget) {
		return "SKIP - Over budget"
	}

	// Giả định 20% trả trước, lãi 8.5%, 20 năm.
	loanInfo := p.CalculateLoan(20, 8.5, 20)
	if loanInfo.MonthlyPayment > maxMonthlyPayment {
		return "CONSIDER - High monthly payment"
	}

	roi := p.CalculateROI(p.Price * assumedMonthlyRentRate)

	if roi > 10 {
		return "BUY NOW - Excellent ROI"
	} else if roi > 6 {
		return "GOOD BUY - Solid investment"
	}

	return "MAYBE - Average investment"
}

// recommendationTiers là thang bậc khuyến nghị, từ thấp đến cao. Dùng index +
// clamp để nâng/hạ một bậc, gọn hơn if lồng nhau.
var recommendationTiers = []string{
	"MAYBE - Average investment",
	"GOOD BUY - Solid investment",
	"BUY NOW - Excellent ROI",
}

// smartRecommendProperty mở rộng recommendProperty với bonus/warning theo
// Task 4.1, trả về (recommendation, details).
func smartRecommendProperty(p Property, budget, maxMonthlyPayment float64) (string, string) {
	if !p.IsAffordable(budget) {
		return "SKIP - Over budget", "Bonus: [], Warnings: [Over budget]"
	}

	if loanInfo := p.CalculateLoan(20, 8.5, 20); loanInfo.MonthlyPayment > maxMonthlyPayment {
		return "CONSIDER - High monthly payment",
			fmt.Sprintf("Monthly payment %s > giới hạn %s",
				formatPriceDetailed(loanInfo.MonthlyPayment), formatPriceDetailed(maxMonthlyPayment))
	}

	var bonus []string
	var warnings []string

	if premiumDistricts[strings.TrimSpace(p.District)] {
		bonus = append(bonus, "Premium location")
	}
	if p.Area >= 50 && p.Area <= 100 {
		bonus = append(bonus, "Optimal size")
	}
	if p.PricePerM2() > 60000000 {
		warnings = append(warnings, "High price per m²")
	}

	// Quan sát: với giả định thuê 1.2%/tháng thì ROI = 1.2*12 = 14.4%/năm cho
	// MỌI property (rent tỷ lệ thuận với giá), nên nhánh ROI luôn ra BUY NOW.
	// Vì vậy chính bonus/warnings mới là phần tạo khác biệt giữa các property.
	roi := p.CalculateROI(p.Price * assumedMonthlyRentRate)

	tier := 0
	switch {
	case roi > 10:
		tier = 2
	case roi > 6:
		tier = 1
	}

	if len(bonus) >= 2 {
		tier++
	}
	if len(warnings) > 0 {
		tier--
	}

	// Clamp về khoảng hợp lệ của recommendationTiers.
	if tier < 0 {
		tier = 0
	}
	if tier > len(recommendationTiers)-1 {
		tier = len(recommendationTiers) - 1
	}

	details := fmt.Sprintf("Bonus: %v, Warnings: %v", bonus, warnings)
	return recommendationTiers[tier], details
}

// runTask41 in khuyến nghị cho toàn bộ danh sách.
func runTask41() {
	runRecommendations(5000000000, 40000000)
}

// runRecommendations in cả base case và bản smart để thấy khác biệt.
func runRecommendations(budget, maxMonthlyPayment float64) {
	printHeader("Property Recommendations")
	fmt.Printf("Budget: %s | Max monthly payment: %s\n\n",
		formatPrice(budget), formatPrice(maxMonthlyPayment))

	for _, prop := range properties {
		recommendation, details := smartRecommendProperty(prop, budget, maxMonthlyPayment)

		fmt.Printf("%s (%s, %.0f m²):\n", prop.Name, prop.District, prop.Area)
		fmt.Printf("   Basic: %s\n", recommendProperty(prop, budget, maxMonthlyPayment))
		fmt.Printf("   Smart: %s\n", recommendation)
		fmt.Printf("   %s\n\n", details)
	}
}

// propertyROI zip property với ROI của nó, để sort không làm lệch cặp
// (property, rent) theo index.
type propertyROI struct {
	Property Property
	ROI      float64
}

// optimizePortfolio giữ đúng chữ ký của PDF. Dùng giả định thuê 1.2%/tháng khi
// không có dữ liệu rent cụ thể.
func optimizePortfolio(properties []Property, totalBudget float64) []Property {
	rents := make([]float64, len(properties))
	for i, prop := range properties {
		rents[i] = prop.Price * assumedMonthlyRentRate
	}
	return optimizePortfolioWithRents(properties, rents, totalBudget)
}

// optimizePortfolioWithRents chọn danh mục theo greedy: sort giảm dần theo ROI
// rồi thêm dần khi còn đủ budget.
//
// Giới hạn đã biết: greedy theo ROI không giải tối ưu bài knapsack — có thể bỏ
// qua tổ hợp tốt hơn. Với dữ liệu mẫu nó cho kết quả đúng, và PDF cũng yêu cầu
// đúng cách "greedily add properties".
func optimizePortfolioWithRents(properties []Property, rents []float64, totalBudget float64) []Property {
	var portfolio []Property
	remainingBudget := totalBudget

	for _, item := range rankByROI(properties, rents) {
		if item.Property.Price <= remainingBudget {
			portfolio = append(portfolio, item.Property)
			remainingBudget -= item.Property.Price
		}
	}

	return portfolio
}

// rankByROI trả về danh sách (property, ROI) sort giảm dần theo ROI.
// Copy slice trước khi sort: sort trực tiếp tham số sẽ mutate slice của caller.
func rankByROI(properties []Property, rents []float64) []propertyROI {
	if len(properties) != len(rents) {
		return nil
	}

	ranked := make([]propertyROI, len(properties))
	for i, prop := range properties {
		ranked[i] = propertyROI{Property: prop, ROI: prop.CalculateROI(rents[i])}
	}

	// Tie-break bằng giá thấp hơn rồi theo tên, để output ổn định khi ROI bằng
	// nhau (sort.Slice không stable).
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].ROI != ranked[j].ROI {
			return ranked[i].ROI > ranked[j].ROI
		}
		if ranked[i].Property.Price != ranked[j].Property.Price {
			return ranked[i].Property.Price < ranked[j].Property.Price
		}
		return ranked[i].Property.Name < ranked[j].Property.Name
	})

	return ranked
}

// runTask42 chạy test case bắt buộc của PDF: budget 8 tỷ trên 3 property gốc.
func runTask42() {
	runPortfolioOptimization(pdfProperties, pdfRents, 8000000000)
}

// runPortfolioOptimization in kết quả tối ưu danh mục.
func runPortfolioOptimization(props []Property, rents []float64, totalBudget float64) {
	printHeader("Portfolio Optimization")
	fmt.Printf("Budget: %s\n", formatPrice(totalBudget))

	if len(props) != len(rents) {
		fmt.Println("Lỗi dữ liệu: số lượng rent không khớp số lượng property")
		return
	}

	portfolio := optimizePortfolioWithRents(props, rents, totalBudget)
	if len(portfolio) == 0 {
		fmt.Println("No properties fit within the budget")
		return
	}

	// Map tên -> rent để tra ROI sau khi sort mà không lệch index.
	rentByName := make(map[string]float64, len(props))
	for i, prop := range props {
		rentByName[prop.Name] = rents[i]
	}

	fmt.Println("Selected Properties:")
	totalInvested := 0.0
	totalROI := 0.0

	for i, prop := range portfolio {
		roi := prop.CalculateROI(rentByName[prop.Name])
		totalInvested += prop.Price
		totalROI += roi

		fmt.Printf("%d. %s: %s (ROI: %.1f%%)\n", i+1, prop.Name, formatPrice(prop.Price), roi)
	}

	// ROI trung bình dùng trung bình cộng đơn giản để khớp output PDF (13.3%).
	// Về tài chính, bình quân gia quyền theo giá chính xác hơn — in thêm để so.
	weightedROI := 0.0
	if totalInvested > 0 {
		for _, prop := range portfolio {
			weightedROI += prop.CalculateROI(rentByName[prop.Name]) * (prop.Price / totalInvested)
		}
	}

	fmt.Printf("\nTotal Invested: %s\n", formatPrice(totalInvested))
	fmt.Printf("Remaining Budget: %s\n", formatPrice(totalBudget-totalInvested))
	fmt.Printf("Portfolio Average ROI: %.1f%%\n", totalROI/float64(len(portfolio)))
	fmt.Printf("Portfolio Weighted ROI: %.1f%% (gia quyền theo giá)\n", weightedROI)
}

// runPart4 chạy cả hai task của Part 4.
func runPart4() {
	runTask41()
	runTask42()
}
