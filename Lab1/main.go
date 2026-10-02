package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FoodItem represents a menu item.
type FoodItem struct {
	Name      string
	Price     float64
	Category  string
	Available bool
}

// OrderLine represents one food item and its quantity inside an order.
type OrderLine struct {
	Item     FoodItem
	Quantity int
}

// Order represents a customer order containing several order lines.
type Order struct {
	Lines []OrderLine
}

// menu is the slice holding all food items.
var menu = []FoodItem{
	{Name: "Chicken Rice", Price: 45000, Category: "Main Dish", Available: true},
	{Name: "Beef Noodle", Price: 50000, Category: "Main Dish", Available: true},
	{Name: "Fried Rice", Price: 40000, Category: "Main Dish", Available: false},
	{Name: "Milk Tea", Price: 25000, Category: "Drink", Available: true},
	{Name: "Iced Coffee", Price: 20000, Category: "Drink", Available: true},
	{Name: "Cheesecake", Price: 35000, Category: "Dessert", Available: false},
	{Name: "Ice Cream", Price: 15000, Category: "Dessert", Available: true},
}

// menuByName is a map used for fast lookup of a food item by its name.
var menuByName = buildMenuByName(menu)

// buildMenuByName builds the lookup map from the menu slice.
func buildMenuByName(items []FoodItem) map[string]FoodItem {
	m := make(map[string]FoodItem)
	for _, item := range items {
		m[item.Name] = item
	}
	return m
}

// DisplayAvailableFood prints every food item in the menu that is available.
func DisplayAvailableFood(items []FoodItem) {
	fmt.Println("Available food:")
	for _, item := range items {
		if item.Available {
			fmt.Printf("- %s | %s | %.0f VND\n", item.Name, item.Category, item.Price)
		}
	}
}

// SearchFoodByName looks up a food item by exact name using the map.
func SearchFoodByName(name string) (FoodItem, bool) {
	item, ok := menuByName[name]
	return item, ok
}

// IsValidQuantity rejects zero or negative quantities.
func IsValidQuantity(quantity int) bool {
	return quantity > 0
}

// AddFoodToOrder validates the food name, availability, and quantity, then
// appends the item to the order. It returns an error describing why the
// item could not be added, if any.
func AddFoodToOrder(order *Order, name string, quantity int) error {
	item, ok := SearchFoodByName(name)
	if !ok {
		return fmt.Errorf("food not found: %s", name)
	}
	if !item.Available {
		return fmt.Errorf("food unavailable: %s", name)
	}
	if !IsValidQuantity(quantity) {
		return fmt.Errorf("invalid quantity: %d", quantity)
	}
	order.Lines = append(order.Lines, OrderLine{Item: item, Quantity: quantity})
	return nil
}

// CalculateSubtotal sums price * quantity for every line in the order.
func CalculateSubtotal(order Order) float64 {
	total := 0.0
	for _, line := range order.Lines {
		total += line.Item.Price * float64(line.Quantity)
	}
	return total
}

// DiscountThreshold and DiscountRate implement the ADAPT requirement:
// orders with a subtotal of at least 300,000 VND receive a 10% discount.
const DiscountThreshold = 300000.0
const DiscountRate = 0.10

// CalculateDiscount returns the discount amount for a given subtotal.
func CalculateDiscount(subtotal float64) float64 {
	if subtotal >= DiscountThreshold {
		return subtotal * DiscountRate
	}
	return 0.0
}

// DisplayOrderSummary prints each line, the subtotal, the discount, and the
// final total of the order.
func DisplayOrderSummary(order Order) {
	fmt.Println("Order summary:")
	for _, line := range order.Lines {
		lineTotal := line.Item.Price * float64(line.Quantity)
		fmt.Printf("- %s x%d = %.0f VND\n", line.Item.Name, line.Quantity, lineTotal)
	}
	subtotal := CalculateSubtotal(order)
	discount := CalculateDiscount(subtotal)
	finalTotal := subtotal - discount
	fmt.Printf("Subtotal: %.0f\n", subtotal)
	fmt.Printf("Discount: %.0f\n", discount)
	fmt.Printf("Final Total: %.0f\n", finalTotal)
}

func main() {
	order := Order{}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("1. Display available food")
		fmt.Println("2. Search food by name")
		fmt.Println("3. Add food to order")
		fmt.Println("4. Show order summary")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			DisplayAvailableFood(menu)
		case "2":
			fmt.Print("Enter food name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			item, ok := SearchFoodByName(name)
			if !ok {
				fmt.Println("Food not found")
			} else {
				fmt.Printf("%s | %s | %.0f VND | Available: %v\n", item.Name, item.Category, item.Price, item.Available)
			}
		case "3":
			fmt.Print("Enter food name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Enter quantity: ")
			qtyStr, _ := reader.ReadString('\n')
			qtyStr = strings.TrimSpace(qtyStr)
			quantity, err := strconv.Atoi(qtyStr)
			if err != nil {
				fmt.Println("Invalid quantity input")
				continue
			}

			if err := AddFoodToOrder(&order, name, quantity); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Added to order.")
			}
		case "4":
			DisplayOrderSummary(order)
		case "5":
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}
