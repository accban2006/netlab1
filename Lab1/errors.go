package main

import "errors"

// Các lỗi validate dùng chung cho tính toán khoản vay và đọc input.
var (
	errInvalidDownPayment  = errors.New("down payment percent phải nằm trong [0, 100]")
	errInvalidInterestRate = errors.New("interest rate không được âm")
	errInvalidYears        = errors.New("số năm phải lớn hơn 0")
)
