package main

import "math"

// Property mô tả một bất động sản. Struct nhỏ nên mọi method dùng value
// receiver (p Property) cho nhất quán — không có method nào cần mutate.
type Property struct {
	Name     string
	Price    float64 // VND
	Area     float64 // m2
	Bedrooms int
	District string
}

// LoanInfo là kết quả tính khoản vay cho một property.
type LoanInfo struct {
	LoanAmount     float64
	MonthlyPayment float64
	TotalInterest  float64
}

// DistrictStats tổng hợp số liệu của một quận.
type DistrictStats struct {
	District      string
	Count         int
	AveragePrice  float64
	MostExpensive Property
}

// PricePerM2 trả về giá trên mỗi m2. Guard Area == 0 để không ra +Inf.
func (p Property) PricePerM2() float64 {
	if p.Area == 0 {
		return 0
	}
	return p.Price / p.Area
}

// IsAffordable cho biết property có nằm trong ngân sách hay không.
func (p Property) IsAffordable(budget float64) bool {
	return p.Price <= budget
}

// CalculateROI trả về ROI thường niên theo %: (rent*12 / Price) * 100.
func (p Property) CalculateROI(monthlyRent float64) float64 {
	if p.Price == 0 {
		return 0
	}
	return (monthlyRent * 12 / p.Price) * 100
}

// InvestmentGrade xếp hạng ROI. PDF viết dạng khoảng ("5-8%") nên quy ước biên
// ở đây là: > 8 EXCELLENT, [5,8] GOOD, [3,5) FAIR, < 3 POOR.
func (p Property) InvestmentGrade(roi float64) string {
	switch {
	case roi > 8:
		return "EXCELLENT"
	case roi >= 5:
		return "GOOD"
	case roi >= 3:
		return "FAIR"
	default:
		return "POOR"
	}
}

// calculateMonthlyPayment tính khoản trả hàng tháng theo công thức annuity:
//
//	M = L * r * (1+r)^n / ((1+r)^n - 1)   với r = annualRate/100/12, n = years*12
func calculateMonthlyPayment(loanAmount, annualRate float64, years int) float64 {
	numPayments := float64(years * 12) // ép float64: years*12 là int
	if numPayments == 0 {
		return 0
	}

	monthlyRate := annualRate / 100 / 12
	if annualRate == 0 {
		// Lãi 0% thì công thức annuity chia cho 0, trả về chia đều.
		return loanAmount / numPayments
	}

	growth := math.Pow(1+monthlyRate, numPayments)
	return loanAmount * monthlyRate * growth / (growth - 1)
}

// CalculateLoan tính khoản vay cho property. Input không hợp lệ trả LoanInfo{}
// zero value — xem CalculateLoanChecked nếu cần biết lý do cụ thể.
func (p Property) CalculateLoan(downPaymentPercent, interestRate float64, years int) LoanInfo {
	info, err := p.CalculateLoanChecked(downPaymentPercent, interestRate, years)
	if err != nil {
		return LoanInfo{}
	}
	return info
}

// CalculateLoanChecked là bản trả error của CalculateLoan, dùng khi cần báo lỗi
// cho người dùng (menu) thay vì âm thầm trả zero value.
func (p Property) CalculateLoanChecked(downPaymentPercent, interestRate float64, years int) (LoanInfo, error) {
	if downPaymentPercent < 0 || downPaymentPercent > 100 {
		return LoanInfo{}, errInvalidDownPayment
	}
	if interestRate < 0 {
		return LoanInfo{}, errInvalidInterestRate
	}
	if years <= 0 {
		return LoanInfo{}, errInvalidYears
	}

	loanAmount := p.Price * (1 - downPaymentPercent/100)
	monthly := calculateMonthlyPayment(loanAmount, interestRate, years)

	return LoanInfo{
		LoanAmount:     loanAmount,
		MonthlyPayment: monthly,
		TotalInterest:  monthly*float64(years*12) - loanAmount,
	}, nil
}
