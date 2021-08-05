package request

import (
	"net/url"
	"strconv"
	"strings"
)

func ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func ReadInt(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	l, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return l
}
