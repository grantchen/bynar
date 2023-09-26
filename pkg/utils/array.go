/**
    @author: dongjs
    @date: 2023/9/11
    @description:
**/

package utils

// IsStringArrayInclude Returns +true+ if the given +obj+ is present in +arr+
func IsStringArrayInclude(arr []string, obj string) bool {
	for _, e := range arr {
		if e == obj {
			return true
		}
	}
	return false
}
