package helper

import (
	"os"
	"reflect"
	"strconv"
)

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	strValue := Getenv(key, "")

	if strValue == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(strValue)

	if err != nil {
		return fallback
	}

	return intValue
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// DiffStringSlices Returns elements existing in slice1 but not in slice2 (strings)
func DiffStringSlices(slice1 []string, slice2 []string) []string {
	mb := make(map[string]struct{}, len(slice2))
	for _, x := range slice2 {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range slice1 {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// DiffIntegerSlices Returns elements existing in slice1 but not in slice2 (integers)
func DiffIntegerSlices(slice1 []int, slice2 []int) []int {
	mb := make(map[int]struct{}, len(slice2))
	for _, x := range slice2 {
		mb[x] = struct{}{}
	}
	var diff []int
	for _, x := range slice1 {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
