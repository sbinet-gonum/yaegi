// Code generated by 'goexports crypto/sha512'. DO NOT EDIT.

// +build go1.11,!go1.12

package stdlib

import (
	"crypto/sha512"
	"reflect"
)

func init() {
	Symbols["crypto/sha512"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"BlockSize":  reflect.ValueOf(sha512.BlockSize),
		"New":        reflect.ValueOf(sha512.New),
		"New384":     reflect.ValueOf(sha512.New384),
		"New512_224": reflect.ValueOf(sha512.New512_224),
		"New512_256": reflect.ValueOf(sha512.New512_256),
		"Size":       reflect.ValueOf(sha512.Size),
		"Size224":    reflect.ValueOf(sha512.Size224),
		"Size256":    reflect.ValueOf(sha512.Size256),
		"Size384":    reflect.ValueOf(sha512.Size384),
		"Sum384":     reflect.ValueOf(sha512.Sum384),
		"Sum512":     reflect.ValueOf(sha512.Sum512),
		"Sum512_224": reflect.ValueOf(sha512.Sum512_224),
		"Sum512_256": reflect.ValueOf(sha512.Sum512_256),
	}
}
