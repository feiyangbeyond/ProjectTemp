package util

import (
	"reflect"
	"strconv"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Itoa(a interface{}) string {
	switch at := a.(type) {
	case int, int8, int16, int64, int32:
		return strconv.FormatInt(reflect.ValueOf(a).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(int64(reflect.ValueOf(a).Uint()), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(a).Float(), 'f', 0, 64)
	case string:
		return at
	}
	return ""
}

func Atoi(a string) int {
	if a == "" {
		return 0
	}
	r, e := strconv.Atoi(a)
	if e == nil {
		return r
	}
	return 0
}
