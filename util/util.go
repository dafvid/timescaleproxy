package util

//import "fmt"

// find string in array of strings
func InArr(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}

// find string in keys of map
func InMap(str string, m map[string]interface{}) bool {
	for k := range m {
		if str == k {
			return true
		}
	}
	return false
}
