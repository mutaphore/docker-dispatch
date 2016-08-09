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

func itemsInList(items []interface{}, list []interface{}) []bool {
	var table map[interface{}]bool
	for _, x := range items {
		table[x] = true
	}
	var exist []bool
	for _, y := range list {
		if table[y] {
			exist = append(exist, true)
		} else {
			exist = append(exist, false)
		}
	}
	return exist
}
