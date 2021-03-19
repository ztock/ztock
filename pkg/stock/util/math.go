package util

import "strconv"

// Calculate the percent increase/decrease from two int numbers.
func PercentageChange(old, new int) float64 {
	return (float64(new-old) / float64(old)) * 100
}

// Calculate the percent increase/decrease from two float64 numbers.
func PercentageChangeFloat(old, new float64) float64 {
	return ((new - old) / old) * 100
}

// Calculate the percent increase/decrease from two string numbers.
func PercentageChangeString(old, new string) (float64, error) {
	o, err := strconv.ParseFloat(old, 64)
	if err != nil {
		return 0, err
	}

	n, err := strconv.ParseFloat(new, 64)
	if err != nil {
		return 0, err
	}

	return PercentageChangeFloat(o, n), nil
}
