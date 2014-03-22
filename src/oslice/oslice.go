package oslice

import (
	"strings"
	"fmt"
)

type Slice struct {
	slice	[]interface{}
	less	func(interface{}, interface{}) bool
}

func New(lecb func(interface{}, interface{}) bool) *Slice {
	return &Slice{less: lecb}
}

func NewStringSlice() *Slice {
	return &Slice{less: func(a, b interface{}) bool {
		return strings.ToLower(a.(string)) < strings.ToLower(b.(string))
	}}
} 

func NewIntSlice() *Slice {
	return &Slice{less: func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}}
} 

func (slice *Slice) Clear() {
	slice.slice = nil
}

func (slice *Slice) Add(item interface{}) {
	if slice.slice == nil {
		slice.slice = []interface{}{item}
	} else {
		ind := bisectLeft(slice.slice, slice.less, item)
		if ind == len(slice.slice) {
			slice.slice = append(slice.slice, item)
		} else {
			temp := []interface{}{}
			temp = append(temp, slice.slice[:ind]...)
			temp = append(temp, item)
			temp = append(temp, slice.slice[ind:]...)
			slice.slice = temp
		}
	}
}

func (slice *Slice) Remove(item interface{}) bool {
	if slice.slice == nil {
		return false
	}
	
	ind := slice.Index(item)
	if ind == -1 {
		return false
	} else {
		temp := []interface{}{}
		temp = append(temp, slice.slice[:ind]...)
		temp = append(temp, slice.slice[ind+1:]...)
		slice.slice = temp
		return true
	}
	return false
}

func (slice *Slice) Index(item interface{}) int {
	if slice.slice == nil {
		return -1
	}
	
	for i, x := range slice.slice {
		if x == item {
			return i
		}
	}
	
	return -1
}

func (slice *Slice) At(ind int) interface{} {
	if ind < 0 || ind >= len(slice.slice) {
		panic(fmt.Sprintf("%d index out of range", ind))
	}
	return slice.slice[ind] 
}

func (slice *Slice) Len() int {
	return len(slice.slice)
}

func bisectLeft(slice []interface{}, 
	less func(interface{}, interface{}) bool, elem interface{}) int {
	left, right := 0, len(slice)
	for left < right {
		middle := int((left + right) / 2)
		if less(slice[middle], elem) {
			left = middle + 1
		} else {
			right = middle
		}
	}
	return left
}