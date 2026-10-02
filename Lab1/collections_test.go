package main

import "testing"

// testProperties là bộ dữ liệu 3 property của PDF, tạo riêng cho test để không
// phụ thuộc vào biến package-level (nếu data.go đổi thì test vẫn đúng).
func testProperties() []Property {
	return []Property{
		{"Saigon Apartment", 2500000000, 75.5, 2, "District 1"},
		{"HCMC House", 4200000000, 120.0, 3, "District 7"},
		{"Budget Studio", 800000000, 35.0, 1, "Binh Thanh"},
	}
}

func testRents() []float64 {
	return []float64{25000000, 35000000, 12000000}
}

func TestFindPropertiesInBudget(t *testing.T) {
	tests := []struct {
		name      string
		budget    float64
		wantNames []string
	}{
		{"3_ty", 3000000000, []string{"Saigon Apartment", "Budget Studio"}},
		{"tat_ca", 5000000000, []string{"Saigon Apartment", "HCMC House", "Budget Studio"}},
		{"bang_gia_re_nhat", 800000000, []string{"Budget Studio"}},
		{"khong_khop_gi", 1000000, nil},
		{"budget_zero", 0, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := findPropertiesInBudget(testProperties(), tc.budget)

			if len(got) != len(tc.wantNames) {
				t.Fatalf("findPropertiesInBudget(%v) trả %d kết quả, want %d",
					tc.budget, len(got), len(tc.wantNames))
			}
			for i, name := range tc.wantNames {
				if got[i].Name != name {
					t.Errorf("kết quả[%d].Name = %q, want %q", i, got[i].Name, name)
				}
			}
		})
	}
}

func TestFindPropertiesByBedrooms(t *testing.T) {
	tests := []struct {
		name      string
		bedrooms  int
		wantCount int
	}{
		{"1_phong", 1, 1},
		{"2_phong", 2, 1},
		{"3_phong", 3, 1},
		{"10_phong_khong_co", 10, 0},
		{"0_phong", 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := findPropertiesByBedrooms(testProperties(), tc.bedrooms); len(got) != tc.wantCount {
				t.Errorf("findPropertiesByBedrooms(%d) trả %d kết quả, want %d",
					tc.bedrooms, len(got), tc.wantCount)
			}
		})
	}
}

func TestFindPropertiesByDistrict(t *testing.T) {
	tests := []struct {
		name      string
		district  string
		wantCount int
	}{
		{"khop_chinh_xac", "District 1", 1},
		{"khong_phan_biet_hoa_thuong", "district 1", 1},
		{"bo_khoang_trang", "  District 7  ", 1},
		{"khong_ton_tai", "District 99", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := findPropertiesByDistrict(testProperties(), tc.district); len(got) != tc.wantCount {
				t.Errorf("findPropertiesByDistrict(%q) trả %d kết quả, want %d",
					tc.district, len(got), tc.wantCount)
			}
		})
	}
}

func TestAnalyzeByDistrict(t *testing.T) {
	props := append(testProperties(), Property{"Ben Thanh Penthouse", 6800000000, 98.0, 3, "District 1"})

	districtMap := analyzeByDistrict(props)

	if len(districtMap) != 3 {
		t.Errorf("số quận = %d, want 3", len(districtMap))
	}
	if got := len(districtMap["District 1"]); got != 2 {
		t.Errorf("District 1 có %d property, want 2", got)
	}
	if got := len(districtMap["Binh Thanh"]); got != 1 {
		t.Errorf("Binh Thanh có %d property, want 1", got)
	}
}

func TestCalculateDistrictStats(t *testing.T) {
	stats := calculateDistrictStats(testProperties())

	if len(stats) != 3 {
		t.Fatalf("len(stats) = %d, want 3", len(stats))
	}

	// Phải sort giảm dần theo giá trung bình.
	wantOrder := []string{"District 7", "District 1", "Binh Thanh"}
	for i, want := range wantOrder {
		if stats[i].District != want {
			t.Errorf("stats[%d].District = %q, want %q", i, stats[i].District, want)
		}
	}
	for i := 1; i < len(stats); i++ {
		if stats[i-1].AveragePrice < stats[i].AveragePrice {
			t.Errorf("stats không sort giảm dần tại index %d: %v < %v",
				i, stats[i-1].AveragePrice, stats[i].AveragePrice)
		}
	}
}

func TestCalculateDistrictStatsMultipleInDistrict(t *testing.T) {
	props := []Property{
		{"A", 2000000000, 50, 2, "District 1"},
		{"B", 4000000000, 100, 3, "District 1"},
		{"C", 1000000000, 40, 1, "Binh Thanh"},
	}

	stats := calculateDistrictStats(props)

	var d1 DistrictStats
	for _, s := range stats {
		if s.District == "District 1" {
			d1 = s
		}
	}

	if d1.Count != 2 {
		t.Errorf("Count = %d, want 2", d1.Count)
	}
	if !almostEqual(d1.AveragePrice, 3000000000, epsilon) {
		t.Errorf("AveragePrice = %v, want 3000000000", d1.AveragePrice)
	}
	if d1.MostExpensive.Name != "B" {
		t.Errorf("MostExpensive.Name = %q, want \"B\"", d1.MostExpensive.Name)
	}
}

func TestCalculateDistrictStatsEmpty(t *testing.T) {
	if got := calculateDistrictStats(nil); len(got) != 0 {
		t.Errorf("calculateDistrictStats(nil) trả %d kết quả, want 0", len(got))
	}
}

func TestFindBestInvestment(t *testing.T) {
	props := testProperties()
	rents := testRents()

	best, roi := findBestInvestmentROI(props, rents)

	if best.Name != "Budget Studio" {
		t.Errorf("best.Name = %q, want \"Budget Studio\"", best.Name)
	}
	if !almostEqual(roi, 18.0, epsilon) {
		t.Errorf("roi = %v, want 18.0", roi)
	}
	if got := findBestInvestment(props, rents); got.Name != "Budget Studio" {
		t.Errorf("findBestInvestment().Name = %q, want \"Budget Studio\"", got.Name)
	}
}

func TestFindBestInvestmentGuards(t *testing.T) {
	// Rents lệch độ dài: phải trả zero value, không được panic index out of range.
	t.Run("rents_lech_do_dai", func(t *testing.T) {
		got := findBestInvestment(testProperties(), []float64{25000000})
		if got != (Property{}) {
			t.Errorf("got %+v, want Property{} zero value", got)
		}
	})

	t.Run("slice_rong", func(t *testing.T) {
		got := findBestInvestment(nil, nil)
		if got != (Property{}) {
			t.Errorf("got %+v, want Property{} zero value", got)
		}
	})

	t.Run("mot_property", func(t *testing.T) {
		got, roi := findBestInvestmentROI(
			[]Property{{Name: "Solo", Price: 1000000000}},
			[]float64{10000000},
		)
		if got.Name != "Solo" || !almostEqual(roi, 12.0, epsilon) {
			t.Errorf("got (%q, %v), want (\"Solo\", 12.0)", got.Name, roi)
		}
	})
}
