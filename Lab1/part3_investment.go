package main

import "fmt"

// Part 3 — Functions and Methods.

// findBestInvestment trả về property có ROI cao nhất. Ghép properties[i] với
// rents[i] theo index nên phải kiểm tra độ dài bằng nhau trước khi loop —
// truy cập lệch index là panic "index out of range".
func findBestInvestment(properties []Property, rents []float64) Property {
	best, _ := findBestInvestmentROI(properties, rents)
	return best
}

// findBestInvestmentROI như findBestInvestment nhưng trả kèm ROI, để caller in
// được ROI mà không phải tính lại. Slice rỗng hoặc lệch độ dài trả Property{}, 0.
func findBestInvestmentROI(properties []Property, rents []float64) (Property, float64) {
	if len(properties) == 0 || len(properties) != len(rents) {
		return Property{}, 0
	}

	best := properties[0]
	bestROI := properties[0].CalculateROI(rents[0])

	for i := 1; i < len(properties); i++ {
		roi := properties[i].CalculateROI(rents[i])
		if roi > bestROI {
			best = properties[i]
			bestROI = roi
		}
	}

	return best, bestROI
}

// runTask31 in phân tích đầu tư cho toàn bộ danh sách.
func runTask31() {
	printHeader("Investment Analysis")

	if len(properties) != len(monthlyRents) {
		fmt.Println("Lỗi dữ liệu: số lượng rent không khớp số lượng property")
		return
	}

	for i, prop := range properties {
		roi := prop.CalculateROI(monthlyRents[i])
		fmt.Printf("%s: ROI %.1f%% per year - %s\n", prop.Name, roi, prop.InvestmentGrade(roi))
	}

	best, bestROI := findBestInvestmentROI(properties, monthlyRents)
	fmt.Printf("\nBest Investment: %s (%.1f%% ROI)\n", best.Name, bestROI)
}

// runTask32 in phân tích khoản vay với test case bắt buộc của PDF:
// 20% trả trước, lãi 8.5%/năm, 20 năm.
func runTask32() {
	runLoanAnalysis(20, 8.5, 20)
}

// runLoanAnalysis in bảng khoản vay cho mọi property theo tham số cho trước.
func runLoanAnalysis(downPaymentPercent, interestRate float64, years int) {
	printHeader("Loan Analysis")
	fmt.Printf("Giả định: %.0f%% trả trước, lãi %.1f%%/năm, %d năm\n\n",
		downPaymentPercent, interestRate, years)

	for _, prop := range properties {
		info, err := prop.CalculateLoanChecked(downPaymentPercent, interestRate, years)
		if err != nil {
			fmt.Printf("%s: không tính được khoản vay (%v)\n", prop.Name, err)
			continue
		}

		fmt.Printf("%s:\n", prop.Name)
		fmt.Printf("   Loan Amount: %s (%.0f%% of price)\n",
			formatPrice(info.LoanAmount), 100-downPaymentPercent)
		fmt.Printf("   Monthly Payment: %s\n", formatPriceDetailed(info.MonthlyPayment))
		fmt.Printf("   Total Interest: %s over %d years\n\n",
			formatPriceDetailed(info.TotalInterest), years)
	}
}

// runPart3 chạy cả hai task của Part 3.
func runPart3() {
	runTask31()
	runTask32()
}
