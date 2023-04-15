package utils

import "fmt"

func ParsePrice(price string) string {
	priceQuery := "&price="
	priceQuery += fmt.Sprintf("%s", price)

	return priceQuery
}
