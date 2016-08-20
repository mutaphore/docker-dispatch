package main

import (
	"crypto/rand"
	"encoding/base64"
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

// Find if an item is in a list
func itemInList(item string, list []string) bool {
	for _, y := range list {
		if item == y {
			return true
		}
	}
	return false
}
