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

// HOMEWORK: customer type and order type tests

func TestCustomerDiscount(t *testing.T) {
	if got := CalculateCustomerDiscount(200000, Regular); got != 0 {
		t.Errorf("regular discount = %.0f, want 0", got)
	}
	if got := CalculateCustomerDiscount(200000, Member); got != 10000 {
		t.Errorf("member discount = %.0f, want 10000", got)
	}
}

func TestDeliveryFee(t *testing.T) {
	if got := CalculateDeliveryFee(Pickup); got != 0 {
		t.Errorf("pickup fee = %.0f, want 0", got)
	}
	if got := CalculateDeliveryFee(Delivery); got != 30000 {
		t.Errorf("delivery fee = %.0f, want 30000", got)
	}
}

// orderWorth builds an order whose subtotal is exactly the given amount, using
// Iced Coffee at 20,000 VND per unit.
func orderWorth(t *testing.T, subtotal float64, customer CustomerType, orderType OrderType) Order {
	t.Helper()
	order := Order{Customer: customer, Type: orderType}
	quantity := int(subtotal / 20000)
	if err := AddFoodToOrder(&order, "Iced Coffee", quantity); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := CalculateSubtotal(order); got != subtotal {
		t.Fatalf("setup subtotal = %.0f, want %.0f", got, subtotal)
	}
	return order
}

// Every combination of customer type and order type, crossed with the three
// subtotal milestones (below / exactly / above the 300,000 VND threshold).
// Expected values follow: Subtotal - VolumeDiscount - MemberDiscount + DeliveryFee,
// where both discounts apply to the original subtotal.
func TestFinalTotalCombinations(t *testing.T) {
	cases := []struct {
		name     string
		subtotal float64
		customer CustomerType
		order    OrderType
		want     float64
	}{
		{"Regular+Pickup below", 200000, Regular, Pickup, 200000},
		{"Regular+Delivery below", 200000, Regular, Delivery, 230000},
		{"Member+Pickup below", 200000, Member, Pickup, 190000},
		{"Member+Delivery below", 200000, Member, Delivery, 220000},

		{"Regular+Pickup at", 300000, Regular, Pickup, 270000},
		{"Regular+Delivery at", 300000, Regular, Delivery, 300000},
		{"Member+Pickup at", 300000, Member, Pickup, 255000},
		{"Member+Delivery at", 300000, Member, Delivery, 285000},

		{"Regular+Pickup above", 400000, Regular, Pickup, 360000},
		{"Regular+Delivery above", 400000, Regular, Delivery, 390000},
		{"Member+Pickup above", 400000, Member, Pickup, 340000},
		{"Member+Delivery above", 400000, Member, Delivery, 370000},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			order := orderWorth(t, c.subtotal, c.customer, c.order)
			if got := CalculateFinalTotal(order); got != c.want {
				t.Errorf("final total = %.0f, want %.0f", got, c.want)
			}
		})
	}
}

// The member discount must not reduce the delivery fee: the gap between a
// pickup and a delivery order stays exactly the flat fee for both customer types.
func TestMemberDiscountDoesNotReduceDeliveryFee(t *testing.T) {
	for _, customer := range []CustomerType{Regular, Member} {
		pickup := CalculateFinalTotal(orderWorth(t, 400000, customer, Pickup))
		delivery := CalculateFinalTotal(orderWorth(t, 400000, customer, Delivery))
		if gap := delivery - pickup; gap != DeliveryFeeAmount {
			t.Errorf("%s: delivery gap = %.0f, want %.0f", customer, gap, DeliveryFeeAmount)
		}
	}
}
