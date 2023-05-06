package utils

import (
	"math/rand"
	"reflect"
	"time"
)

func Random(min, max int64) int {
	return int(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(max-min+1) + min)
}

func InSlice(val any, slice any) bool {
	values := reflect.ValueOf(slice)

	if reflect.TypeOf(slice).Kind() == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

func IsPalindrome(s string) bool {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}

	return true
}
