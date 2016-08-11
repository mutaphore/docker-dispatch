package main

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"
)

// Generate a random base64 string of fixed size
func genRandStr(size int) (error, string) {
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		return err, ""
	}
	rs := base64.URLEncoding.EncodeToString(rb)
	return nil, rs
}

func itemInList(item interface{}, list []interface{}) bool {
	for _, y := range list {
		if reflect.DeepEqual(item, y) {
			return true
		}
	}
	return false
}
