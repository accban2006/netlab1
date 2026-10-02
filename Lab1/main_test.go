package main

import (
	"bufio"
	"strings"
	"testing"
)

// readerFrom tạo bufio.Reader từ chuỗi, để test input mà không cần stdin thật.
func readerFrom(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

func TestReadInt(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal int
		want       int
		wantErr    bool
		wantEOF    bool
	}{
		{name: "so_hop_le", input: "42\n", defaultVal: 0, want: 42},
		{name: "co_khoang_trang", input: "  7  \n", defaultVal: 0, want: 7},
		{name: "so_am", input: "-5\n", defaultVal: 0, want: -5},
		{name: "enter_rong_dung_default", input: "\n", defaultVal: 20, want: 20},
		{name: "khong_phai_so", input: "abc\n", defaultVal: 0, wantErr: true},
		{name: "so_thuc_khong_hop_le", input: "1.5\n", defaultVal: 0, wantErr: true},
		{name: "eof", input: "", defaultVal: 0, wantErr: true, wantEOF: true},
		// Dòng cuối không có '\n' vẫn phải đọc được.
		{name: "khong_co_newline_cuoi", input: "9", defaultVal: 0, want: 9},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readInt(readerFrom(tc.input), "", tc.defaultVal)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("readInt(%q) err = nil, want lỗi", tc.input)
				}
				if isEOF(err) != tc.wantEOF {
					t.Errorf("isEOF(err) = %t, want %t (err = %v)", isEOF(err), tc.wantEOF, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("readInt(%q) err = %v, want nil", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("readInt(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestReadFloat(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal float64
		want       float64
		wantErr    bool
	}{
		{name: "so_nguyen", input: "3000000000\n", want: 3000000000},
		{name: "so_thuc", input: "8.5\n", want: 8.5},
		{name: "ky_hieu_khoa_hoc", input: "2.5e9\n", want: 2500000000},
		{name: "enter_rong_dung_default", input: "\n", defaultVal: 8.5, want: 8.5},
		{name: "khong_phai_so", input: "abc\n", wantErr: true},
		{name: "eof", input: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readFloat(readerFrom(tc.input), "", tc.defaultVal)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("readFloat(%q) err = nil, want lỗi", tc.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("readFloat(%q) err = %v, want nil", tc.input, err)
			}
			if !almostEqual(got, tc.want, epsilon) {
				t.Errorf("readFloat(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestReadLineMultipleLines(t *testing.T) {
	// Đọc nhiều dòng liên tiếp từ cùng một reader — mô phỏng menu hỏi nhiều
	// tham số trong một lượt.
	reader := readerFrom("20\n8.5\n20\n")

	down, err := readFloat(reader, "", 0)
	if err != nil || !almostEqual(down, 20, epsilon) {
		t.Fatalf("dòng 1 = %v (err %v), want 20", down, err)
	}

	rate, err := readFloat(reader, "", 0)
	if err != nil || !almostEqual(rate, 8.5, epsilon) {
		t.Fatalf("dòng 2 = %v (err %v), want 8.5", rate, err)
	}

	years, err := readInt(reader, "", 0)
	if err != nil || years != 20 {
		t.Fatalf("dòng 3 = %v (err %v), want 20", years, err)
	}

	// Hết input -> EOF.
	if _, err := readInt(reader, "", 0); !isEOF(err) {
		t.Errorf("sau khi hết input, isEOF = false (err = %v), want true", err)
	}
}
