package main

import "testing"

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want string
	}{
		{"ty", 2500000000, "2.5 tỷ VND"},
		{"trieu", 850000000, "850 triệu VND"},
		{"nghin", 5000, "5 nghìn VND"},
		{"zero", 0, "0 VND"},
		{"bien_1_ty", 1000000000, "1.0 tỷ VND"},
		{"duoi_bien_1_ty", 999999999, "1000 triệu VND"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatPrice(tc.in); got != tc.want {
				t.Errorf("formatPrice(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatPriceDetailed(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want string
	}{
		{"monthly_payment", 17356465, "17.4 triệu VND"},
		{"total_interest", 2165551600, "2.17 tỷ VND"},
		{"zero", 0, "0 VND"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatPriceDetailed(tc.in); got != tc.want {
				t.Errorf("formatPriceDetailed(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatVND(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want string
	}{
		{"ty", 2500000000, "2,500,000,000"},
		{"price_per_m2", 33112582.78, "33,112,583"},
		{"ba_chu_so", 800, "800"},
		{"bon_chu_so", 1234, "1,234"},
		{"zero", 0, "0"},
		{"am", -1234567, "-1,234,567"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := formatVND(tc.in); got != tc.want {
				t.Errorf("formatVND(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestCategorizeProperty(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want string
	}{
		// Biên của PDF dùng `>` chứ không phải `>=`.
		{"bien_luxury_bang", 50000000, "PREMIUM"},
		{"bien_luxury_vuot", 50000001, "LUXURY"},
		{"bien_premium_bang", 30000000, "STANDARD"},
		{"bien_premium_vuot", 30000001, "PREMIUM"},
		{"bien_standard_bang", 20000000, "BUDGET"},
		{"bien_standard_vuot", 20000001, "STANDARD"},
		{"zero", 0, "BUDGET"},
		// Saigon Apartment: 2.5e9/75.5 = 33.1M/m² -> PREMIUM, không phải
		// STANDARD như bảng output của PDF.
		{"saigon_apartment", 33112582.78, "PREMIUM"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := categorizeProperty(tc.in); got != tc.want {
				t.Errorf("categorizeProperty(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
