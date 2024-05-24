package utils

import "strconv"

func ConvertToInt(s string) int {
	res, _ := strconv.Atoi(s)
	return res
}
