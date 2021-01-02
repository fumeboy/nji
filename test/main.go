package main

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"unsafe"
)

type a struct {

}

func (a *a) A(){}

type AA interface {
	A()
}

type b struct {
	a
	AAA string
}

type face struct {
	typ *struct{}
	word unsafe.Pointer
}

func main() {
	var bb = b{AAA: "123"}
	var cc = b{AAA: "32"}
	var c interface{} = cc
	var b = c
	(*(*face)(unsafe.Pointer(&b))).word = unsafe.Pointer((&bb))
	jsoniter.NewDecoder(bytes.NewReader([]byte("{\"AAA\":\"222\"}"))).Decode(&b)
	fmt.Println(bb, cc)
}
