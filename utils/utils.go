package utils

import "fmt"

func ParsePrices(price []int) string {

	priceQuery := "&price="

	for i := 0; i < len(price); i++ {
		if i == len(price)-1 {
			priceQuery += fmt.Sprintf("%s", fmt.Sprint(price[i]))
		} else {
			priceQuery += fmt.Sprintf("%s,", fmt.Sprint(price[i]))
		}

	}

	return priceQuery
}
