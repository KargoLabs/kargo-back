package converter

import "strconv"

// ConvertToGreaterThanZeroFloat converts string to float64
// Returns false if convertion fails or is not greater than 0
func ConvertToGreaterThanZeroFloat(floatString string) (float64, bool) {
	float, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		return 0, false
	}

	return float, float >= 0
}
