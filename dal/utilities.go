package dal

import (
	"strconv"
	"strings"
)

func GetPriceRange(priceRangeStr string) []float64 {
	priceRange := make([]float64, 2)
	var err error
	priceRangeStrArr := strings.Split(priceRangeStr, "-")
	if len(priceRangeStrArr) != 2 {
		return nil
	}
	priceRange[0], err = strconv.ParseFloat(priceRangeStrArr[0], 64)
	if err != nil {
		return nil
	}
	priceRange[1], err = strconv.ParseFloat(priceRangeStrArr[1], 64)
	if err != nil {
		return nil
	}
	return priceRange
}
func FieldExist(check string) bool {
	if check == "" {
		return false
	}
	return true
}
