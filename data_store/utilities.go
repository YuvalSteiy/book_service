package data_store

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func GetPriceRange(priceRangeStr string) ([]float64, error) {
	priceRange := make([]float64, 2)
	var err error
	priceRangeStrArr := strings.Split(priceRangeStr, "-")
	if len(priceRangeStrArr) != 2 {
		return nil, errors.New("Invalid Price Range")
	}

	priceRange[0], err = strconv.ParseFloat(priceRangeStrArr[0], 64)
	if err != nil {
		return nil, err
	}

	priceRange[1], err = strconv.ParseFloat(priceRangeStrArr[1], 64)
	if err != nil {
		return nil, err
	}

	return priceRange, nil
}


