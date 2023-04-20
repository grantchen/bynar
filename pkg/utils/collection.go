package utils

import "strings"

func MergeMaps(maps ...[]interface{}) (result []interface{}) {
	for _, m := range maps {
		result = append(result, m...)
	}
	return result
}

// Returns the substring between two substring (start, end), within one main string(string)
func GetStringBetween(str string, start string, end string) (result string) {
	ini := strings.Index(str, start)
	if ini == 0 {
		return ""
	}
	ini += len(start)
	endStr := strings.Index(str, end)
	if endStr == 0 {
		return ""
	}
	return str[ini:endStr]
}
