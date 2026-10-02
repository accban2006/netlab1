package main

// properties là dữ liệu mẫu. Ba property đầu lấy nguyên từ PDF; hai cái sau
// thêm vào để Part 1/2 có đủ số liệu so sánh (nhiều property trong cùng quận,
// và một property hạng LUXURY để category summary không rỗng).
var properties = []Property{
	{"Saigon Apartment", 2500000000, 75.5, 2, "District 1"},
	{"HCMC House", 4200000000, 120.0, 3, "District 7"},
	{"Budget Studio", 800000000, 35.0, 1, "Binh Thanh"},
	{"Thao Dien Villa", 9500000000, 145.0, 4, "District 2"},
	{"Ben Thanh Penthouse", 6800000000, 98.0, 3, "District 1"},
}

// monthlyRents cùng thứ tự với properties — Task 3.1 ghép theo index.
var monthlyRents = []float64{25000000, 35000000, 12000000, 48000000, 40000000}

// premiumDistricts là các quận được cộng điểm ở Task 4.1.
var premiumDistricts = map[string]bool{
	"District 1": true,
	"District 2": true,
	"District 7": true,
}

// pdfProperties và pdfRents là đúng bộ dữ liệu 3 property của PDF, dùng cho
// các task cần đối chiếu output mẫu (Task 1.1, 1.2, 4.2). Copy ra slice riêng
// thay vì dùng properties[:3] — slice cắt ra vẫn dùng chung backing array,
// nên mọi thay đổi sẽ lan sang nhau.
var pdfProperties = append([]Property(nil), properties[:3]...)

var pdfRents = append([]float64(nil), monthlyRents[:3]...)
