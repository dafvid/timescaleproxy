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

func Prepend(i interface{}, arr []interface{}) []interface{} {
	return append([]interface{}{i}, arr...)
}
