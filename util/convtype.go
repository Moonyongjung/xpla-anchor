package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ToString(value interface{}, defaultValue string) string {
	str := strings.TrimSpace(fmt.Sprintf("%v", value))
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}

func ToStringTrim(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	s = s[1 : len(s)-1]
	str := strings.TrimSpace(s)
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}

func ToInt(value interface{}) (int, error) {
	return strconv.Atoi(value.(string))
}

func FromStringToUint32(s string) uint32 {
	number64, _ := strconv.ParseUint(s, 10, 64)
	number32 := uint32(number64)
	return number32
}

func FromStringToUint64(value string) uint64 {
	number, _ := strconv.ParseUint(value, 10, 64)
	return number
}

func FromUint64ToString(value uint64) string {
	return strconv.Itoa(int(value))
}
