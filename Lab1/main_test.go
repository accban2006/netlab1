package main

import "testing"

// Case: Normal order -> Correct subtotal
func TestNormalOrderSubtotal(t *testing.T) {
	order := Order{}
	if err := AddFoodToOrder(&order, "Chicken Rice", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := CalculateSubtotal(order)
	want := 90000.0
	if got != want {
		t.Errorf("subtotal = %.0f, want %.0f", got, want)
	}
}

// Case: Quantity = 1 -> Accepted
func TestQuantityOneAccepted(t *testing.T) {
	order := Order{}
	if err := AddFoodToOrder(&order, "Milk Tea", 1); err != nil {
		t.Errorf("expected quantity 1 to be accepted, got error: %v", err)
	}
}

// Case: Quantity = 0 -> Rejected
func TestQuantityZeroRejected(t *testing.T) {
	order := Order{}
	if err := AddFoodToOrder(&order, "Milk Tea", 0); err == nil {
		t.Errorf("expected quantity 0 to be rejected")
	}
}

// Case: Negative quantity -> Rejected
func TestNegativeQuantityRejected(t *testing.T) {
	order := Order{}
	if err := AddFoodToOrder(&order, "Milk Tea", -3); err == nil {
		t.Errorf("expected negative quantity to be rejected")
	}
}

// Case: Unavailable food -> Rejected
func TestUnavailableFoodRejected(t *testing.T) {
	order := Order{}
	if err := AddFoodToOrder(&order, "Fried Rice", 1); err == nil {
		t.Errorf("expected unavailable food to be rejected")
	}
}

// Case: Food not found -> Handled safely
func TestFoodNotFoundHandledSafely(t *testing.T) {
	order := Order{}
	err := AddFoodToOrder(&order, "Pho", 1)
	if err == nil {
		t.Errorf("expected error for food not found")
	}
	if len(order.Lines) != 0 {
		t.Errorf("order should remain unchanged when food is not found")
	}
}

// Case: Several different items -> Correct total
func TestSeveralDifferentItemsCorrectTotal(t *testing.T) {
	order := Order{}
	AddFoodToOrder(&order, "Chicken Rice", 1) // 45000
	AddFoodToOrder(&order, "Milk Tea", 2)     // 50000
	AddFoodToOrder(&order, "Ice Cream", 3)    // 45000
	got := CalculateSubtotal(order)
	want := 140000.0
	if got != want {
		t.Errorf("subtotal = %.0f, want %.0f", got, want)
	}
}

// ADAPT: discount tests

// Below 300,000 VND
func TestDiscountBelowThreshold(t *testing.T) {
	got := CalculateDiscount(250000)
	if got != 0 {
		t.Errorf("discount = %.0f, want 0", got)
	}
}

// Exactly 300,000 VND
func TestDiscountAtThreshold(t *testing.T) {
	got := CalculateDiscount(300000)
	want := 30000.0
	if got != want {
		t.Errorf("discount = %.0f, want %.0f", got, want)
	}
}

// Above 300,000 VND
func TestDiscountAboveThreshold(t *testing.T) {
	got := CalculateDiscount(350000)
	want := 35000.0
	if got != want {
		t.Errorf("discount = %.0f, want %.0f", got, want)
	}
}
