package utils

/**
* This is for common functions im too lazy to keep writing
**/

import (
	"reflect"
	"strings"
)

func CurrentIndex(list interface{}, item interface{}) int {
	listValue := reflect.ValueOf(list)
	itemValue := reflect.ValueOf(item)

	for i := 0; i < listValue.Len(); i++ {
		if reflect.DeepEqual(listValue.Index(i).Interface(), itemValue.Interface()) {
			return i
		}
	}

	return -1
}

func FilterArrayByValue(arr []string, value string) []string {
	var filteredArr []string

	for _, item := range arr {
		if strings.Contains(item, value) {
			filteredArr = append(filteredArr, item)
		}
	}

	return filteredArr
}
