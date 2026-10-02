package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Part 5 — Menu tổng hợp.
//
// PDF dùng fmt.Scanln(&choice) với choice là int. Nếu người dùng nhập chữ thì
// Scanln fail nhưng để nguyên ký tự trong buffer, nên vòng for sẽ quay vô hạn
// in menu liên tục. Ở đây đọc cả dòng bằng bufio.Reader rồi tự parse, và return
// khi gặp EOF để không treo khi chạy qua pipe.

// readLine đọc một dòng từ reader. err != nil nghĩa là EOF hoặc lỗi I/O.
func readLine(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)

	line, err := reader.ReadString('\n')
	if err != nil {
		// Vẫn có thể có dữ liệu trước EOF (dòng cuối không có '\n').
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed, nil
		}
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// readInt đọc một số nguyên. Enter rỗng trả về defaultValue.
func readInt(reader *bufio.Reader, prompt string, defaultValue int) (int, error) {
	line, err := readLine(reader, prompt)
	if err != nil {
		return 0, err
	}
	if line == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("%q không phải số nguyên hợp lệ", line)
	}
	return value, nil
}

// readFloat đọc một số thực. Enter rỗng trả về defaultValue.
func readFloat(reader *bufio.Reader, prompt string, defaultValue float64) (float64, error) {
	line, err := readLine(reader, prompt)
	if err != nil {
		return 0, err
	}
	if line == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return 0, fmt.Errorf("%q không phải số hợp lệ", line)
	}
	return value, nil
}

// isEOF phân biệt lỗi EOF (phải thoát chương trình) với lỗi parse (chỉ cần
// báo rồi hỏi lại). Dùng errors.Is thay vì so sánh chuỗi err.Error().
func isEOF(err error) bool {
	return errors.Is(err, io.EOF)
}

func printMenu() {
	fmt.Println("\n=== Property Analyzer Menu ===")
	fmt.Println("1. View all properties")
	fmt.Println("2. Search by budget")
	fmt.Println("3. Investment analysis")
	fmt.Println("4. Loan calculator")
	fmt.Println("5. Get recommendations")
	fmt.Println("6. Optimize portfolio")
	fmt.Println("7. District analysis")
	fmt.Println("8. Run all parts (demo Task 1.1 -> 4.2)")
	fmt.Println("0. Exit")
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()

		choice, err := readInt(reader, "Choose option: ", -1)
		if isEOF(err) {
			fmt.Println("\nEOF nhận được. Goodbye!")
			return
		}
		if err != nil {
			fmt.Printf("Input không hợp lệ: %v\n", err)
			continue
		}

		switch choice {
		case 1:
			printHeader("All Properties")
			printPropertyTable(properties)

		case 2:
			if handleSearch(reader) {
				return
			}

		case 3:
			runTask31()

		case 4:
			if handleLoan(reader) {
				return
			}

		case 5:
			if handleRecommendations(reader) {
				return
			}

		case 6:
			if handlePortfolio(reader) {
				return
			}

		case 7:
			runTask22()

		case 8:
			runPart1()
			runPart2()
			runPart3()
			runPart4()

		case 0:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option!")
		}
	}
}

// Các handler dưới đây trả true nếu gặp EOF — caller phải thoát chương trình.

func handleSearch(reader *bufio.Reader) bool {
	budget, err := readFloat(reader, "Nhập budget tối đa (VND, Enter = 3000000000): ", 3000000000)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}
	if budget < 0 {
		fmt.Println("Budget không được âm.")
		return false
	}

	result := findPropertiesInBudget(properties, budget)
	fmt.Printf("\nProperties under %s (%s):\n", formatVND(budget), formatPrice(budget))
	printPropertyTable(result)

	bedrooms, err := readInt(reader, "\nLọc thêm theo số phòng ngủ (Enter = bỏ qua, 0 = bỏ qua): ", 0)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}
	if bedrooms > 0 {
		filtered := findPropertiesByBedrooms(result, bedrooms)
		fmt.Printf("\nTrong budget và có %d phòng ngủ:\n", bedrooms)
		printPropertyTable(filtered)
	}

	return false
}

func handleLoan(reader *bufio.Reader) bool {
	downPayment, err := readFloat(reader, "% trả trước (Enter = 20): ", 20)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}

	interestRate, err := readFloat(reader, "Lãi suất %/năm (Enter = 8.5): ", 8.5)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}

	years, err := readInt(reader, "Số năm (Enter = 20): ", 20)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}

	// Validate sớm ở đây để báo lỗi rõ ràng, thay vì để CalculateLoan trả zero.
	switch {
	case downPayment < 0 || downPayment > 100:
		fmt.Printf("Lỗi: %v (nhận %.1f)\n", errInvalidDownPayment, downPayment)
		return false
	case interestRate < 0:
		fmt.Printf("Lỗi: %v (nhận %.1f)\n", errInvalidInterestRate, interestRate)
		return false
	case years <= 0:
		fmt.Printf("Lỗi: %v (nhận %d)\n", errInvalidYears, years)
		return false
	}

	runLoanAnalysis(downPayment, interestRate, years)
	return false
}

func handleRecommendations(reader *bufio.Reader) bool {
	budget, err := readFloat(reader, "Budget (VND, Enter = 5000000000): ", 5000000000)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}

	maxMonthly, err := readFloat(reader, "Khoản trả tối đa/tháng (VND, Enter = 40000000): ", 40000000)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}

	if budget < 0 || maxMonthly < 0 {
		fmt.Println("Budget và khoản trả tối đa không được âm.")
		return false
	}

	runRecommendations(budget, maxMonthly)
	return false
}

func handlePortfolio(reader *bufio.Reader) bool {
	budget, err := readFloat(reader, "Tổng budget (VND, Enter = 8000000000): ", 8000000000)
	if isEOF(err) {
		return true
	}
	if err != nil {
		fmt.Printf("Input không hợp lệ: %v\n", err)
		return false
	}
	if budget < 0 {
		fmt.Println("Budget không được âm.")
		return false
	}

	runPortfolioOptimization(properties, monthlyRents, budget)
	return false
}
