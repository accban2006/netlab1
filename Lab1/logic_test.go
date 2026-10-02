package main

import (
	"strings"
	"testing"
)

func TestRecommendProperty(t *testing.T) {
	tests := []struct {
		name       string
		p          Property
		budget     float64
		maxMonthly float64
		want       string
	}{
		{
			name:       "over_budget",
			p:          Property{Name: "Villa", Price: 9500000000, Area: 145, District: "District 2"},
			budget:     5000000000,
			maxMonthly: 40000000,
			want:       "SKIP - Over budget",
		},
		{
			name:       "monthly_payment_qua_cao",
			p:          Property{Name: "HCMC House", Price: 4200000000, Area: 120, District: "District 7"},
			budget:     5000000000,
			maxMonthly: 10000000,
			want:       "CONSIDER - High monthly payment",
		},
		{
			// ROI với giả định 1.2%/tháng luôn = 14.4% nên ra BUY NOW.
			name:       "roi_tot",
			p:          Property{Name: "Saigon Apartment", Price: 2500000000, Area: 75.5, District: "District 1"},
			budget:     5000000000,
			maxMonthly: 40000000,
			want:       "BUY NOW - Excellent ROI",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := recommendProperty(tc.p, tc.budget, tc.maxMonthly); got != tc.want {
				t.Errorf("recommendProperty() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSmartRecommendProperty(t *testing.T) {
	tests := []struct {
		name         string
		p            Property
		budget       float64
		maxMonthly   float64
		wantRec      string
		wantInDetail []string
	}{
		{
			name:         "premium_location_va_optimal_size",
			p:            Property{Name: "Saigon Apartment", Price: 2500000000, Area: 75.5, District: "District 1"},
			budget:       5000000000,
			maxMonthly:   40000000,
			wantRec:      "BUY NOW - Excellent ROI",
			wantInDetail: []string{"Premium location", "Optimal size"},
		},
		{
			name:         "premium_location_nhung_qua_to",
			p:            Property{Name: "HCMC House", Price: 4200000000, Area: 120, District: "District 7"},
			budget:       5000000000,
			maxMonthly:   40000000,
			wantRec:      "BUY NOW - Excellent ROI",
			wantInDetail: []string{"Premium location"},
		},
		{
			name:         "khong_bonus",
			p:            Property{Name: "Budget Studio", Price: 800000000, Area: 35, District: "Binh Thanh"},
			budget:       5000000000,
			maxMonthly:   40000000,
			wantRec:      "BUY NOW - Excellent ROI",
			wantInDetail: []string{"Bonus: []"},
		},
		{
			// Giá/m² > 60M -> warning, hạ một bậc từ BUY NOW xuống GOOD BUY.
			name:         "warning_ha_mot_bac",
			p:            Property{Name: "Expensive Box", Price: 3000000000, Area: 40, District: "Binh Thanh"},
			budget:       5000000000,
			maxMonthly:   40000000,
			wantRec:      "GOOD BUY - Solid investment",
			wantInDetail: []string{"High price per m²"},
		},
		{
			name:         "over_budget",
			p:            Property{Name: "Villa", Price: 9500000000, Area: 145, District: "District 2"},
			budget:       5000000000,
			maxMonthly:   40000000,
			wantRec:      "SKIP - Over budget",
			wantInDetail: []string{"Over budget"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec, details := smartRecommendProperty(tc.p, tc.budget, tc.maxMonthly)

			if rec != tc.wantRec {
				t.Errorf("recommendation = %q, want %q", rec, tc.wantRec)
			}
			for _, want := range tc.wantInDetail {
				if !strings.Contains(details, want) {
					t.Errorf("details = %q, thiếu %q", details, want)
				}
			}
		})
	}
}

func TestOptimizePortfolio(t *testing.T) {
	props := testProperties()
	rents := testRents()

	// Test case bắt buộc của PDF: budget 8 tỷ -> chọn cả 3, sort giảm dần ROI.
	portfolio := optimizePortfolioWithRents(props, rents, 8000000000)

	wantOrder := []string{"Budget Studio", "Saigon Apartment", "HCMC House"}
	if len(portfolio) != len(wantOrder) {
		t.Fatalf("len(portfolio) = %d, want %d", len(portfolio), len(wantOrder))
	}
	for i, want := range wantOrder {
		if portfolio[i].Name != want {
			t.Errorf("portfolio[%d].Name = %q, want %q", i, portfolio[i].Name, want)
		}
	}

	total := 0.0
	for _, prop := range portfolio {
		total += prop.Price
	}
	if !almostEqual(total, 7500000000, epsilon) {
		t.Errorf("tổng đầu tư = %v, want 7.5e9", total)
	}
	if remaining := 8000000000 - total; !almostEqual(remaining, 500000000, epsilon) {
		t.Errorf("budget còn lại = %v, want 500e6", remaining)
	}
}

func TestOptimizePortfolioBudgetLimits(t *testing.T) {
	props := testProperties()
	rents := testRents()

	tests := []struct {
		name      string
		budget    float64
		wantNames []string
	}{
		{"budget_zero", 0, nil},
		{"budget_am", -1000, nil},
		{"chi_du_studio", 1000000000, []string{"Budget Studio"}},
		// 3.5 tỷ: greedy lấy Studio (800M, ROI 18) rồi Saigon (2.5 tỷ, ROI 12),
		// còn 200M nên không đủ cho HCMC House.
		{"greedy_bo_qua_dat_nhat", 3500000000, []string{"Budget Studio", "Saigon Apartment"}},
		{"du_tat_ca", 10000000000, []string{"Budget Studio", "Saigon Apartment", "HCMC House"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := optimizePortfolioWithRents(props, rents, tc.budget)

			if len(got) != len(tc.wantNames) {
				t.Fatalf("len = %d, want %d (got %v)", len(got), len(tc.wantNames), got)
			}
			for i, want := range tc.wantNames {
				if got[i].Name != want {
					t.Errorf("portfolio[%d].Name = %q, want %q", i, got[i].Name, want)
				}
			}
		})
	}
}

// TestOptimizePortfolioDoesNotMutateCaller kiểm tra bug kinh điển với slice
// trong Go: sort trực tiếp tham số sẽ đổi thứ tự slice của caller.
func TestOptimizePortfolioDoesNotMutateCaller(t *testing.T) {
	props := testProperties()
	rents := testRents()

	before := make([]string, len(props))
	for i, prop := range props {
		before[i] = prop.Name
	}

	optimizePortfolioWithRents(props, rents, 8000000000)
	optimizePortfolio(props, 8000000000)

	for i, want := range before {
		if props[i].Name != want {
			t.Errorf("slice gốc bị mutate: props[%d].Name = %q, want %q",
				i, props[i].Name, want)
		}
	}
}

func TestOptimizePortfolioMismatchedRents(t *testing.T) {
	// Rents lệch độ dài: rankByROI trả nil nên portfolio rỗng, không panic.
	got := optimizePortfolioWithRents(testProperties(), []float64{25000000}, 8000000000)
	if len(got) != 0 {
		t.Errorf("len = %d, want 0 khi rents lệch độ dài", len(got))
	}
}

func TestRankByROI(t *testing.T) {
	ranked := rankByROI(testProperties(), testRents())

	if len(ranked) != 3 {
		t.Fatalf("len(ranked) = %d, want 3", len(ranked))
	}
	for i := 1; i < len(ranked); i++ {
		if ranked[i-1].ROI < ranked[i].ROI {
			t.Errorf("không sort giảm dần tại index %d: %v < %v",
				i, ranked[i-1].ROI, ranked[i].ROI)
		}
	}
	if !almostEqual(ranked[0].ROI, 18.0, epsilon) {
		t.Errorf("ranked[0].ROI = %v, want 18.0", ranked[0].ROI)
	}
}

// TestRankByROITieBreak kiểm tra tie-break khi ROI bằng nhau — sort.Slice không
// stable nên cần tie-break để output ổn định giữa các lần chạy.
func TestRankByROITieBreak(t *testing.T) {
	props := []Property{
		{Name: "Dat", Price: 2000000000, Area: 50},
		{Name: "Re", Price: 1000000000, Area: 25},
	}
	// Cùng ROI 12% cho cả hai.
	rents := []float64{20000000, 10000000}

	ranked := rankByROI(props, rents)

	// Tie-break: giá thấp hơn lên trước.
	if ranked[0].Property.Name != "Re" {
		t.Errorf("ranked[0].Name = %q, want \"Re\" (giá thấp lên trước khi ROI bằng nhau)",
			ranked[0].Property.Name)
	}
}
