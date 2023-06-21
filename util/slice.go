package util

import "sort"

func DeleteVal(values []string, val string) []string {
	newValues := []string{}
	for _, v := range values {
		if v != val {
			newValues = append(newValues, v)
		}
	}
	return newValues
}

func ContainsString(values []string, val string) bool {
	sort.Strings(values)
	return sort.SearchStrings(values, val) != len(values)
}

func ReturnAnyNotEmpty(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}
