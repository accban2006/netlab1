package main

import (
	"math"
	"testing"
)

// epsilon dùng cho mọi so sánh float64 — so sánh bằng `==` với float64 là sai
// vì sai số làm tròn nhị phân.
const epsilon = 1e-6

// almostEqual so sánh hai float64 với sai số cho phép.
func almostEqual(got, want, tol float64) bool {
	return math.Abs(got-want) < tol
}

func TestPricePerM2(t *testing.T) {
	tests := []struct {
		name string
		p    Property
		want float64
	}{
		{"saigon_apartment", Property{Price: 2500000000, Area: 75.5}, 2500000000.0 / 75.5},
		{"hcmc_house", Property{Price: 4200000000, Area: 120}, 35000000},
		{"area_zero_guard", Property{Price: 2500000000, Area: 0}, 0},
		{"zero_value", Property{}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.p.PricePerM2(); !almostEqual(got, tc.want, epsilon) {
				t.Errorf("PricePerM2() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsAffordable(t *testing.T) {
	p := Property{Price: 2500000000}

	tests := []struct {
		name   string
		budget float64
		want   bool
	}{
		{"du_budget", 3000000000, true},
		{"bang_gia", 2500000000, true}, // <= nên bằng giá vẫn affordable
		{"thieu_budget", 2499999999, false},
		{"budget_zero", 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := p.IsAffordable(tc.budget); got != tc.want {
				t.Errorf("IsAffordable(%v) = %t, want %t", tc.budget, got, tc.want)
			}
		})
	}
}

func TestCalculateROI(t *testing.T) {
	tests := []struct {
		name string
		p    Property
		rent float64
		want float64
	}{
		{"saigon_apartment", Property{Price: 2500000000}, 25000000, 12.0},
		{"hcmc_house", Property{Price: 4200000000}, 35000000, 10.0},
		{"budget_studio", Property{Price: 800000000}, 12000000, 18.0},
		{"price_zero_guard", Property{Price: 0}, 25000000, 0},
		{"rent_zero", Property{Price: 2500000000}, 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.p.CalculateROI(tc.rent); !almostEqual(got, tc.want, epsilon) {
				t.Errorf("CalculateROI(%v) = %v, want %v", tc.rent, got, tc.want)
			}
		})
	}
}

func TestInvestmentGrade(t *testing.T) {
	var p Property

	tests := []struct {
		name string
		roi  float64
		want string
	}{
		{"excellent", 12.0, "EXCELLENT"},
		{"bien_excellent_good", 8.0, "GOOD"}, // quy ước: >8 mới EXCELLENT
		{"vuot_bien_excellent", 8.01, "EXCELLENT"},
		{"good", 6.5, "GOOD"},
		{"bien_good_fair", 5.0, "GOOD"}, // quy ước: >=5 là GOOD
		{"fair", 4.0, "FAIR"},
		{"bien_fair_poor", 3.0, "FAIR"}, // quy ước: >=3 là FAIR
		{"poor", 2.9, "POOR"},
		{"zero", 0, "POOR"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := p.InvestmentGrade(tc.roi); got != tc.want {
				t.Errorf("InvestmentGrade(%v) = %q, want %q", tc.roi, got, tc.want)
			}
		})
	}
}

func TestCalculateMonthlyPayment(t *testing.T) {
	tests := []struct {
		name  string
		loan  float64
		rate  float64
		years int
		want  float64
		tol   float64
	}{
		// Test case bắt buộc của PDF: loan 2 tỷ, 8.5%, 20 năm.
		{"pdf_saigon_apartment", 2000000000, 8.5, 20, 17356465, 1},
		{"pdf_hcmc_house", 3360000000, 8.5, 20, 29158861, 1},
		// Lãi 0%: chia đều, không dùng công thức annuity (tránh chia 0).
		{"lai_zero", 2400000000, 0, 10, 20000000, epsilon},
		{"loan_zero", 0, 8.5, 20, 0, epsilon},
		// years = 0 -> numPayments = 0, guard trả 0 thay vì NaN.
		{"years_zero_guard", 2000000000, 8.5, 0, 0, epsilon},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := calculateMonthlyPayment(tc.loan, tc.rate, tc.years)
			if !almostEqual(got, tc.want, tc.tol) {
				t.Errorf("calculateMonthlyPayment(%v, %v, %d) = %v, want %v (±%v)",
					tc.loan, tc.rate, tc.years, got, tc.want, tc.tol)
			}
		})
	}
}

func TestCalculateLoan(t *testing.T) {
	saigon := Property{Name: "Saigon Apartment", Price: 2500000000, Area: 75.5}

	// Test case bắt buộc của PDF: 20% down, 8.5%, 20 năm.
	info := saigon.CalculateLoan(20, 8.5, 20)

	if !almostEqual(info.LoanAmount, 2000000000, epsilon) {
		t.Errorf("LoanAmount = %v, want 2000000000", info.LoanAmount)
	}
	if !almostEqual(info.MonthlyPayment, 17356465, 1) {
		t.Errorf("MonthlyPayment = %v, want ≈17356465", info.MonthlyPayment)
	}
	if !almostEqual(info.TotalInterest, 2165551600, 1000) {
		t.Errorf("TotalInterest = %v, want ≈2.166e9", info.TotalInterest)
	}

	// Kiểm chứng quan hệ: TotalInterest = monthly*n - loan.
	wantInterest := info.MonthlyPayment*float64(20*12) - info.LoanAmount
	if !almostEqual(info.TotalInterest, wantInterest, epsilon) {
		t.Errorf("TotalInterest = %v, không khớp monthly*n - loan = %v",
			info.TotalInterest, wantInterest)
	}
}

func TestCalculateLoanInvalidInput(t *testing.T) {
	saigon := Property{Price: 2500000000, Area: 75.5}

	tests := []struct {
		name        string
		downPayment float64
		rate        float64
		years       int
	}{
		{"down_payment_qua_100", 150, 8.5, 20},
		{"down_payment_am", -10, 8.5, 20},
		{"lai_am", 20, -1, 20},
		{"years_zero", 20, 8.5, 0},
		{"years_am", 20, 8.5, -5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := saigon.CalculateLoan(tc.downPayment, tc.rate, tc.years)
			if got != (LoanInfo{}) {
				t.Errorf("CalculateLoan(%v, %v, %d) = %+v, want LoanInfo{} zero value",
					tc.downPayment, tc.rate, tc.years, got)
			}

			if _, err := saigon.CalculateLoanChecked(tc.downPayment, tc.rate, tc.years); err == nil {
				t.Errorf("CalculateLoanChecked(%v, %v, %d) err = nil, want lỗi validate",
					tc.downPayment, tc.rate, tc.years)
			}
		})
	}
}

func TestCalculateLoanFullDownPayment(t *testing.T) {
	// 100% trả trước là hợp lệ: loan = 0, monthly = 0, lãi = 0.
	saigon := Property{Price: 2500000000, Area: 75.5}

	info, err := saigon.CalculateLoanChecked(100, 8.5, 20)
	if err != nil {
		t.Fatalf("CalculateLoanChecked(100, 8.5, 20) err = %v, want nil", err)
	}
	if !almostEqual(info.LoanAmount, 0, epsilon) || !almostEqual(info.MonthlyPayment, 0, epsilon) {
		t.Errorf("với 100%% trả trước, got %+v, want tất cả bằng 0", info)
	}
}
