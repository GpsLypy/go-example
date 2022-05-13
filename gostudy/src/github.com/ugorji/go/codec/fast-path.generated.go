// +build !notfastpath

// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED from fast-path.go.tmpl
// ************************************************************

package codec

// Fast path functions try to create a fast path encode or decode implementation
// for common maps and slices.
//
// We define the functions and register then in this single file
// so as not to pollute the encode.go and decode.go, and create a dependency in there.
// This file can be omitted without causing a build failure.
//
// The advantage of fast paths is:
//    - Many calls bypass reflection altogether
//
// Currently support
//    - slice of all builtin types,
//    - map of all builtin types to string or interface value
//    - symmetrical maps of all builtin types (e.g. str-str, uint8-uint8)
// This should provide adequate "typical" implementations.
//
// Note that fast track decode functions must handle values for which an address cannot be obtained.
// For example:
//   m2 := map[string]int{}
//   p2 := []interface{}{m2}
//   // decoding into p2 will bomb if fast track functions do not treat like unaddressable.
//

import (
	"reflect"
	"sort"
)

const fastpathEnabled = true

const fastpathCheckNilFalse = false // for reflect
const fastpathCheckNilTrue = true   // for type switch

type fastpathT struct{}

var fastpathTV fastpathT

type fastpathE struct {
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encFnInfo, reflect.Value)
	decfn func(*decFnInfo, reflect.Value)
}

type fastpathA [271]fastpathE

func (x *fastpathA) index(rtid uintptr) int {
	// use binary search to grab the index (adapted from sort/search.go)
	h, i, j := 0, 0, 271 // len(x)
	for i < j {
		h = i + (j-i)/2
		if x[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
	}
	if i < 271 && x[i].rtid == rtid {
		return i
	}
	return -1
}

type fastpathAslice []fastpathE

func (x fastpathAslice) Len() int           { return len(x) }
func (x fastpathAslice) Less(i, j int) bool { return x[i].rtid < x[j].rtid }
func (x fastpathAslice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

var fastpathAV fastpathA

// due to possible initialization loop error, make fastpath in an init()
func init() {
	i := 0
	fn := func(v interface{}, fe func(*encFnInfo, reflect.Value), fd func(*decFnInfo, reflect.Value)) (f fastpathE) {
		xrt := reflect.TypeOf(v)
		xptr := reflect.ValueOf(xrt).Pointer()
		fastpathAV[i] = fastpathE{xptr, xrt, fe, fd}
		i++
		return
	}

	fn([]interface{}(nil), (*encFnInfo).fastpathEncSliceIntfR, (*decFnInfo).fastpathDecSliceIntfR)
	fn([]string(nil), (*encFnInfo).fastpathEncSliceStringR, (*decFnInfo).fastpathDecSliceStringR)
	fn([]float32(nil), (*encFnInfo).fastpathEncSliceFloat32R, (*decFnInfo).fastpathDecSliceFloat32R)
	fn([]float64(nil), (*encFnInfo).fastpathEncSliceFloat64R, (*decFnInfo).fastpathDecSliceFloat64R)
	fn([]uint(nil), (*encFnInfo).fastpathEncSliceUintR, (*decFnInfo).fastpathDecSliceUintR)
	fn([]uint16(nil), (*encFnInfo).fastpathEncSliceUint16R, (*decFnInfo).fastpathDecSliceUint16R)
	fn([]uint32(nil), (*encFnInfo).fastpathEncSliceUint32R, (*decFnInfo).fastpathDecSliceUint32R)
	fn([]uint64(nil), (*encFnInfo).fastpathEncSliceUint64R, (*decFnInfo).fastpathDecSliceUint64R)
	fn([]uintptr(nil), (*encFnInfo).fastpathEncSliceUintptrR, (*decFnInfo).fastpathDecSliceUintptrR)
	fn([]int(nil), (*encFnInfo).fastpathEncSliceIntR, (*decFnInfo).fastpathDecSliceIntR)
	fn([]int8(nil), (*encFnInfo).fastpathEncSliceInt8R, (*decFnInfo).fastpathDecSliceInt8R)
	fn([]int16(nil), (*encFnInfo).fastpathEncSliceInt16R, (*decFnInfo).fastpathDecSliceInt16R)
	fn([]int32(nil), (*encFnInfo).fastpathEncSliceInt32R, (*decFnInfo).fastpathDecSliceInt32R)
	fn([]int64(nil), (*encFnInfo).fastpathEncSliceInt64R, (*decFnInfo).fastpathDecSliceInt64R)
	fn([]bool(nil), (*encFnInfo).fastpathEncSliceBoolR, (*decFnInfo).fastpathDecSliceBoolR)

	fn(map[interface{}]interface{}(nil), (*encFnInfo).fastpathEncMapIntfIntfR, (*decFnInfo).fastpathDecMapIntfIntfR)
	fn(map[interface{}]string(nil), (*encFnInfo).fastpathEncMapIntfStringR, (*decFnInfo).fastpathDecMapIntfStringR)
	fn(map[interface{}]uint(nil), (*encFnInfo).fastpathEncMapIntfUintR, (*decFnInfo).fastpathDecMapIntfUintR)
	fn(map[interface{}]uint8(nil), (*encFnInfo).fastpathEncMapIntfUint8R, (*decFnInfo).fastpathDecMapIntfUint8R)
	fn(map[interface{}]uint16(nil), (*encFnInfo).fastpathEncMapIntfUint16R, (*decFnInfo).fastpathDecMapIntfUint16R)
	fn(map[interface{}]uint32(nil), (*encFnInfo).fastpathEncMapIntfUint32R, (*decFnInfo).fastpathDecMapIntfUint32R)
	fn(map[interface{}]uint64(nil), (*encFnInfo).fastpathEncMapIntfUint64R, (*decFnInfo).fastpathDecMapIntfUint64R)
	fn(map[interface{}]uintptr(nil), (*encFnInfo).fastpathEncMapIntfUintptrR, (*decFnInfo).fastpathDecMapIntfUintptrR)
	fn(map[interface{}]int(nil), (*encFnInfo).fastpathEncMapIntfIntR, (*decFnInfo).fastpathDecMapIntfIntR)
	fn(map[interface{}]int8(nil), (*encFnInfo).fastpathEncMapIntfInt8R, (*decFnInfo).fastpathDecMapIntfInt8R)
	fn(map[interface{}]int16(nil), (*encFnInfo).fastpathEncMapIntfInt16R, (*decFnInfo).fastpathDecMapIntfInt16R)
	fn(map[interface{}]int32(nil), (*encFnInfo).fastpathEncMapIntfInt32R, (*decFnInfo).fastpathDecMapIntfInt32R)
	fn(map[interface{}]int64(nil), (*encFnInfo).fastpathEncMapIntfInt64R, (*decFnInfo).fastpathDecMapIntfInt64R)
	fn(map[interface{}]float32(nil), (*encFnInfo).fastpathEncMapIntfFloat32R, (*decFnInfo).fastpathDecMapIntfFloat32R)
	fn(map[interface{}]float64(nil), (*encFnInfo).fastpathEncMapIntfFloat64R, (*decFnInfo).fastpathDecMapIntfFloat64R)
	fn(map[interface{}]bool(nil), (*encFnInfo).fastpathEncMapIntfBoolR, (*decFnInfo).fastpathDecMapIntfBoolR)
	fn(map[string]interface{}(nil), (*encFnInfo).fastpathEncMapStringIntfR, (*decFnInfo).fastpathDecMapStringIntfR)
	fn(map[string]string(nil), (*encFnInfo).fastpathEncMapStringStringR, (*decFnInfo).fastpathDecMapStringStringR)
	fn(map[string]uint(nil), (*encFnInfo).fastpathEncMapStringUintR, (*decFnInfo).fastpathDecMapStringUintR)
	fn(map[string]uint8(nil), (*encFnInfo).fastpathEncMapStringUint8R, (*decFnInfo).fastpathDecMapStringUint8R)
	fn(map[string]uint16(nil), (*encFnInfo).fastpathEncMapStringUint16R, (*decFnInfo).fastpathDecMapStringUint16R)
	fn(map[string]uint32(nil), (*encFnInfo).fastpathEncMapStringUint32R, (*decFnInfo).fastpathDecMapStringUint32R)
	fn(map[string]uint64(nil), (*encFnInfo).fastpathEncMapStringUint64R, (*decFnInfo).fastpathDecMapStringUint64R)
	fn(map[string]uintptr(nil), (*encFnInfo).fastpathEncMapStringUintptrR, (*decFnInfo).fastpathDecMapStringUintptrR)
	fn(map[string]int(nil), (*encFnInfo).fastpathEncMapStringIntR, (*decFnInfo).fastpathDecMapStringIntR)
	fn(map[string]int8(nil), (*encFnInfo).fastpathEncMapStringInt8R, (*decFnInfo).fastpathDecMapStringInt8R)
	fn(map[string]int16(nil), (*encFnInfo).fastpathEncMapStringInt16R, (*decFnInfo).fastpathDecMapStringInt16R)
	fn(map[string]int32(nil), (*encFnInfo).fastpathEncMapStringInt32R, (*decFnInfo).fastpathDecMapStringInt32R)
	fn(map[string]int64(nil), (*encFnInfo).fastpathEncMapStringInt64R, (*decFnInfo).fastpathDecMapStringInt64R)
	fn(map[string]float32(nil), (*encFnInfo).fastpathEncMapStringFloat32R, (*decFnInfo).fastpathDecMapStringFloat32R)
	fn(map[string]float64(nil), (*encFnInfo).fastpathEncMapStringFloat64R, (*decFnInfo).fastpathDecMapStringFloat64R)
	fn(map[string]bool(nil), (*encFnInfo).fastpathEncMapStringBoolR, (*decFnInfo).fastpathDecMapStringBoolR)
	fn(map[float32]interface{}(nil), (*encFnInfo).fastpathEncMapFloat32IntfR, (*decFnInfo).fastpathDecMapFloat32IntfR)
	fn(map[float32]string(nil), (*encFnInfo).fastpathEncMapFloat32StringR, (*decFnInfo).fastpathDecMapFloat32StringR)
	fn(map[float32]uint(nil), (*encFnInfo).fastpathEncMapFloat32UintR, (*decFnInfo).fastpathDecMapFloat32UintR)
	fn(map[float32]uint8(nil), (*encFnInfo).fastpathEncMapFloat32Uint8R, (*decFnInfo).fastpathDecMapFloat32Uint8R)
	fn(map[float32]uint16(nil), (*encFnInfo).fastpathEncMapFloat32Uint16R, (*decFnInfo).fastpathDecMapFloat32Uint16R)
	fn(map[float32]uint32(nil), (*encFnInfo).fastpathEncMapFloat32Uint32R, (*decFnInfo).fastpathDecMapFloat32Uint32R)
	fn(map[float32]uint64(nil), (*encFnInfo).fastpathEncMapFloat32Uint64R, (*decFnInfo).fastpathDecMapFloat32Uint64R)
	fn(map[float32]uintptr(nil), (*encFnInfo).fastpathEncMapFloat32UintptrR, (*decFnInfo).fastpathDecMapFloat32UintptrR)
	fn(map[float32]int(nil), (*encFnInfo).fastpathEncMapFloat32IntR, (*decFnInfo).fastpathDecMapFloat32IntR)
	fn(map[float32]int8(nil), (*encFnInfo).fastpathEncMapFloat32Int8R, (*decFnInfo).fastpathDecMapFloat32Int8R)
	fn(map[float32]int16(nil), (*encFnInfo).fastpathEncMapFloat32Int16R, (*decFnInfo).fastpathDecMapFloat32Int16R)
	fn(map[float32]int32(nil), (*encFnInfo).fastpathEncMapFloat32Int32R, (*decFnInfo).fastpathDecMapFloat32Int32R)
	fn(map[float32]int64(nil), (*encFnInfo).fastpathEncMapFloat32Int64R, (*decFnInfo).fastpathDecMapFloat32Int64R)
	fn(map[float32]float32(nil), (*encFnInfo).fastpathEncMapFloat32Float32R, (*decFnInfo).fastpathDecMapFloat32Float32R)
	fn(map[float32]float64(nil), (*encFnInfo).fastpathEncMapFloat32Float64R, (*decFnInfo).fastpathDecMapFloat32Float64R)
	fn(map[float32]bool(nil), (*encFnInfo).fastpathEncMapFloat32BoolR, (*decFnInfo).fastpathDecMapFloat32BoolR)
	fn(map[float64]interface{}(nil), (*encFnInfo).fastpathEncMapFloat64IntfR, (*decFnInfo).fastpathDecMapFloat64IntfR)
	fn(map[float64]string(nil), (*encFnInfo).fastpathEncMapFloat64StringR, (*decFnInfo).fastpathDecMapFloat64StringR)
	fn(map[float64]uint(nil), (*encFnInfo).fastpathEncMapFloat64UintR, (*decFnInfo).fastpathDecMapFloat64UintR)
	fn(map[float64]uint8(nil), (*encFnInfo).fastpathEncMapFloat64Uint8R, (*decFnInfo).fastpathDecMapFloat64Uint8R)
	fn(map[float64]uint16(nil), (*encFnInfo).fastpathEncMapFloat64Uint16R, (*decFnInfo).fastpathDecMapFloat64Uint16R)
	fn(map[float64]uint32(nil), (*encFnInfo).fastpathEncMapFloat64Uint32R, (*decFnInfo).fastpathDecMapFloat64Uint32R)
	fn(map[float64]uint64(nil), (*encFnInfo).fastpathEncMapFloat64Uint64R, (*decFnInfo).fastpathDecMapFloat64Uint64R)
	fn(map[float64]uintptr(nil), (*encFnInfo).fastpathEncMapFloat64UintptrR, (*decFnInfo).fastpathDecMapFloat64UintptrR)
	fn(map[float64]int(nil), (*encFnInfo).fastpathEncMapFloat64IntR, (*decFnInfo).fastpathDecMapFloat64IntR)
	fn(map[float64]int8(nil), (*encFnInfo).fastpathEncMapFloat64Int8R, (*decFnInfo).fastpathDecMapFloat64Int8R)
	fn(map[float64]int16(nil), (*encFnInfo).fastpathEncMapFloat64Int16R, (*decFnInfo).fastpathDecMapFloat64Int16R)
	fn(map[float64]int32(nil), (*encFnInfo).fastpathEncMapFloat64Int32R, (*decFnInfo).fastpathDecMapFloat64Int32R)
	fn(map[float64]int64(nil), (*encFnInfo).fastpathEncMapFloat64Int64R, (*decFnInfo).fastpathDecMapFloat64Int64R)
	fn(map[float64]float32(nil), (*encFnInfo).fastpathEncMapFloat64Float32R, (*decFnInfo).fastpathDecMapFloat64Float32R)
	fn(map[float64]float64(nil), (*encFnInfo).fastpathEncMapFloat64Float64R, (*decFnInfo).fastpathDecMapFloat64Float64R)
	fn(map[float64]bool(nil), (*encFnInfo).fastpathEncMapFloat64BoolR, (*decFnInfo).fastpathDecMapFloat64BoolR)
	fn(map[uint]interface{}(nil), (*encFnInfo).fastpathEncMapUintIntfR, (*decFnInfo).fastpathDecMapUintIntfR)
	fn(map[uint]string(nil), (*encFnInfo).fastpathEncMapUintStringR, (*decFnInfo).fastpathDecMapUintStringR)
	fn(map[uint]uint(nil), (*encFnInfo).fastpathEncMapUintUintR, (*decFnInfo).fastpathDecMapUintUintR)
	fn(map[uint]uint8(nil), (*encFnInfo).fastpathEncMapUintUint8R, (*decFnInfo).fastpathDecMapUintUint8R)
	fn(map[uint]uint16(nil), (*encFnInfo).fastpathEncMapUintUint16R, (*decFnInfo).fastpathDecMapUintUint16R)
	fn(map[uint]uint32(nil), (*encFnInfo).fastpathEncMapUintUint32R, (*decFnInfo).fastpathDecMapUintUint32R)
	fn(map[uint]uint64(nil), (*encFnInfo).fastpathEncMapUintUint64R, (*decFnInfo).fastpathDecMapUintUint64R)
	fn(map[uint]uintptr(nil), (*encFnInfo).fastpathEncMapUintUintptrR, (*decFnInfo).fastpathDecMapUintUintptrR)
	fn(map[uint]int(nil), (*encFnInfo).fastpathEncMapUintIntR, (*decFnInfo).fastpathDecMapUintIntR)
	fn(map[uint]int8(nil), (*encFnInfo).fastpathEncMapUintInt8R, (*decFnInfo).fastpathDecMapUintInt8R)
	fn(map[uint]int16(nil), (*encFnInfo).fastpathEncMapUintInt16R, (*decFnInfo).fastpathDecMapUintInt16R)
	fn(map[uint]int32(nil), (*encFnInfo).fastpathEncMapUintInt32R, (*decFnInfo).fastpathDecMapUintInt32R)
	fn(map[uint]int64(nil), (*encFnInfo).fastpathEncMapUintInt64R, (*decFnInfo).fastpathDecMapUintInt64R)
	fn(map[uint]float32(nil), (*encFnInfo).fastpathEncMapUintFloat32R, (*decFnInfo).fastpathDecMapUintFloat32R)
	fn(map[uint]float64(nil), (*encFnInfo).fastpathEncMapUintFloat64R, (*decFnInfo).fastpathDecMapUintFloat64R)
	fn(map[uint]bool(nil), (*encFnInfo).fastpathEncMapUintBoolR, (*decFnInfo).fastpathDecMapUintBoolR)
	fn(map[uint8]interface{}(nil), (*encFnInfo).fastpathEncMapUint8IntfR, (*decFnInfo).fastpathDecMapUint8IntfR)
	fn(map[uint8]string(nil), (*encFnInfo).fastpathEncMapUint8StringR, (*decFnInfo).fastpathDecMapUint8StringR)
	fn(map[uint8]uint(nil), (*encFnInfo).fastpathEncMapUint8UintR, (*decFnInfo).fastpathDecMapUint8UintR)
	fn(map[uint8]uint8(nil), (*encFnInfo).fastpathEncMapUint8Uint8R, (*decFnInfo).fastpathDecMapUint8Uint8R)
	fn(map[uint8]uint16(nil), (*encFnInfo).fastpathEncMapUint8Uint16R, (*decFnInfo).fastpathDecMapUint8Uint16R)
	fn(map[uint8]uint32(nil), (*encFnInfo).fastpathEncMapUint8Uint32R, (*decFnInfo).fastpathDecMapUint8Uint32R)
	fn(map[uint8]uint64(nil), (*encFnInfo).fastpathEncMapUint8Uint64R, (*decFnInfo).fastpathDecMapUint8Uint64R)
	fn(map[uint8]uintptr(nil), (*encFnInfo).fastpathEncMapUint8UintptrR, (*decFnInfo).fastpathDecMapUint8UintptrR)
	fn(map[uint8]int(nil), (*encFnInfo).fastpathEncMapUint8IntR, (*decFnInfo).fastpathDecMapUint8IntR)
	fn(map[uint8]int8(nil), (*encFnInfo).fastpathEncMapUint8Int8R, (*decFnInfo).fastpathDecMapUint8Int8R)
	fn(map[uint8]int16(nil), (*encFnInfo).fastpathEncMapUint8Int16R, (*decFnInfo).fastpathDecMapUint8Int16R)
	fn(map[uint8]int32(nil), (*encFnInfo).fastpathEncMapUint8Int32R, (*decFnInfo).fastpathDecMapUint8Int32R)
	fn(map[uint8]int64(nil), (*encFnInfo).fastpathEncMapUint8Int64R, (*decFnInfo).fastpathDecMapUint8Int64R)
	fn(map[uint8]float32(nil), (*encFnInfo).fastpathEncMapUint8Float32R, (*decFnInfo).fastpathDecMapUint8Float32R)
	fn(map[uint8]float64(nil), (*encFnInfo).fastpathEncMapUint8Float64R, (*decFnInfo).fastpathDecMapUint8Float64R)
	fn(map[uint8]bool(nil), (*encFnInfo).fastpathEncMapUint8BoolR, (*decFnInfo).fastpathDecMapUint8BoolR)
	fn(map[uint16]interface{}(nil), (*encFnInfo).fastpathEncMapUint16IntfR, (*decFnInfo).fastpathDecMapUint16IntfR)
	fn(map[uint16]string(nil), (*encFnInfo).fastpathEncMapUint16StringR, (*decFnInfo).fastpathDecMapUint16StringR)
	fn(map[uint16]uint(nil), (*encFnInfo).fastpathEncMapUint16UintR, (*decFnInfo).fastpathDecMapUint16UintR)
	fn(map[uint16]uint8(nil), (*encFnInfo).fastpathEncMapUint16Uint8R, (*decFnInfo).fastpathDecMapUint16Uint8R)
	fn(map[uint16]uint16(nil), (*encFnInfo).fastpathEncMapUint16Uint16R, (*decFnInfo).fastpathDecMapUint16Uint16R)
	fn(map[uint16]uint32(nil), (*encFnInfo).fastpathEncMapUint16Uint32R, (*decFnInfo).fastpathDecMapUint16Uint32R)
	fn(map[uint16]uint64(nil), (*encFnInfo).fastpathEncMapUint16Uint64R, (*decFnInfo).fastpathDecMapUint16Uint64R)
	fn(map[uint16]uintptr(nil), (*encFnInfo).fastpathEncMapUint16UintptrR, (*decFnInfo).fastpathDecMapUint16UintptrR)
	fn(map[uint16]int(nil), (*encFnInfo).fastpathEncMapUint16IntR, (*decFnInfo).fastpathDecMapUint16IntR)
	fn(map[uint16]int8(nil), (*encFnInfo).fastpathEncMapUint16Int8R, (*decFnInfo).fastpathDecMapUint16Int8R)
	fn(map[uint16]int16(nil), (*encFnInfo).fastpathEncMapUint16Int16R, (*decFnInfo).fastpathDecMapUint16Int16R)
	fn(map[uint16]int32(nil), (*encFnInfo).fastpathEncMapUint16Int32R, (*decFnInfo).fastpathDecMapUint16Int32R)
	fn(map[uint16]int64(nil), (*encFnInfo).fastpathEncMapUint16Int64R, (*decFnInfo).fastpathDecMapUint16Int64R)
	fn(map[uint16]float32(nil), (*encFnInfo).fastpathEncMapUint16Float32R, (*decFnInfo).fastpathDecMapUint16Float32R)
	fn(map[uint16]float64(nil), (*encFnInfo).fastpathEncMapUint16Float64R, (*decFnInfo).fastpathDecMapUint16Float64R)
	fn(map[uint16]bool(nil), (*encFnInfo).fastpathEncMapUint16BoolR, (*decFnInfo).fastpathDecMapUint16BoolR)
	fn(map[uint32]interface{}(nil), (*encFnInfo).fastpathEncMapUint32IntfR, (*decFnInfo).fastpathDecMapUint32IntfR)
	fn(map[uint32]string(nil), (*encFnInfo).fastpathEncMapUint32StringR, (*decFnInfo).fastpathDecMapUint32StringR)
	fn(map[uint32]uint(nil), (*encFnInfo).fastpathEncMapUint32UintR, (*decFnInfo).fastpathDecMapUint32UintR)
	fn(map[uint32]uint8(nil), (*encFnInfo).fastpathEncMapUint32Uint8R, (*decFnInfo).fastpathDecMapUint32Uint8R)
	fn(map[uint32]uint16(nil), (*encFnInfo).fastpathEncMapUint32Uint16R, (*decFnInfo).fastpathDecMapUint32Uint16R)
	fn(map[uint32]uint32(nil), (*encFnInfo).fastpathEncMapUint32Uint32R, (*decFnInfo).fastpathDecMapUint32Uint32R)
	fn(map[uint32]uint64(nil), (*encFnInfo).fastpathEncMapUint32Uint64R, (*decFnInfo).fastpathDecMapUint32Uint64R)
	fn(map[uint32]uintptr(nil), (*encFnInfo).fastpathEncMapUint32UintptrR, (*decFnInfo).fastpathDecMapUint32UintptrR)
	fn(map[uint32]int(nil), (*encFnInfo).fastpathEncMapUint32IntR, (*decFnInfo).fastpathDecMapUint32IntR)
	fn(map[uint32]int8(nil), (*encFnInfo).fastpathEncMapUint32Int8R, (*decFnInfo).fastpathDecMapUint32Int8R)
	fn(map[uint32]int16(nil), (*encFnInfo).fastpathEncMapUint32Int16R, (*decFnInfo).fastpathDecMapUint32Int16R)
	fn(map[uint32]int32(nil), (*encFnInfo).fastpathEncMapUint32Int32R, (*decFnInfo).fastpathDecMapUint32Int32R)
	fn(map[uint32]int64(nil), (*encFnInfo).fastpathEncMapUint32Int64R, (*decFnInfo).fastpathDecMapUint32Int64R)
	fn(map[uint32]float32(nil), (*encFnInfo).fastpathEncMapUint32Float32R, (*decFnInfo).fastpathDecMapUint32Float32R)
	fn(map[uint32]float64(nil), (*encFnInfo).fastpathEncMapUint32Float64R, (*decFnInfo).fastpathDecMapUint32Float64R)
	fn(map[uint32]bool(nil), (*encFnInfo).fastpathEncMapUint32BoolR, (*decFnInfo).fastpathDecMapUint32BoolR)
	fn(map[uint64]interface{}(nil), (*encFnInfo).fastpathEncMapUint64IntfR, (*decFnInfo).fastpathDecMapUint64IntfR)
	fn(map[uint64]string(nil), (*encFnInfo).fastpathEncMapUint64StringR, (*decFnInfo).fastpathDecMapUint64StringR)
	fn(map[uint64]uint(nil), (*encFnInfo).fastpathEncMapUint64UintR, (*decFnInfo).fastpathDecMapUint64UintR)
	fn(map[uint64]uint8(nil), (*encFnInfo).fastpathEncMapUint64Uint8R, (*decFnInfo).fastpathDecMapUint64Uint8R)
	fn(map[uint64]uint16(nil), (*encFnInfo).fastpathEncMapUint64Uint16R, (*decFnInfo).fastpathDecMapUint64Uint16R)
	fn(map[uint64]uint32(nil), (*encFnInfo).fastpathEncMapUint64Uint32R, (*decFnInfo).fastpathDecMapUint64Uint32R)
	fn(map[uint64]uint64(nil), (*encFnInfo).fastpathEncMapUint64Uint64R, (*decFnInfo).fastpathDecMapUint64Uint64R)
	fn(map[uint64]uintptr(nil), (*encFnInfo).fastpathEncMapUint64UintptrR, (*decFnInfo).fastpathDecMapUint64UintptrR)
	fn(map[uint64]int(nil), (*encFnInfo).fastpathEncMapUint64IntR, (*decFnInfo).fastpathDecMapUint64IntR)
	fn(map[uint64]int8(nil), (*encFnInfo).fastpathEncMapUint64Int8R, (*decFnInfo).fastpathDecMapUint64Int8R)
	fn(map[uint64]int16(nil), (*encFnInfo).fastpathEncMapUint64Int16R, (*decFnInfo).fastpathDecMapUint64Int16R)
	fn(map[uint64]int32(nil), (*encFnInfo).fastpathEncMapUint64Int32R, (*decFnInfo).fastpathDecMapUint64Int32R)
	fn(map[uint64]int64(nil), (*encFnInfo).fastpathEncMapUint64Int64R, (*decFnInfo).fastpathDecMapUint64Int64R)
	fn(map[uint64]float32(nil), (*encFnInfo).fastpathEncMapUint64Float32R, (*decFnInfo).fastpathDecMapUint64Float32R)
	fn(map[uint64]float64(nil), (*encFnInfo).fastpathEncMapUint64Float64R, (*decFnInfo).fastpathDecMapUint64Float64R)
	fn(map[uint64]bool(nil), (*encFnInfo).fastpathEncMapUint64BoolR, (*decFnInfo).fastpathDecMapUint64BoolR)
	fn(map[uintptr]interface{}(nil), (*encFnInfo).fastpathEncMapUintptrIntfR, (*decFnInfo).fastpathDecMapUintptrIntfR)
	fn(map[uintptr]string(nil), (*encFnInfo).fastpathEncMapUintptrStringR, (*decFnInfo).fastpathDecMapUintptrStringR)
	fn(map[uintptr]uint(nil), (*encFnInfo).fastpathEncMapUintptrUintR, (*decFnInfo).fastpathDecMapUintptrUintR)
	fn(map[uintptr]uint8(nil), (*encFnInfo).fastpathEncMapUintptrUint8R, (*decFnInfo).fastpathDecMapUintptrUint8R)
	fn(map[uintptr]uint16(nil), (*encFnInfo).fastpathEncMapUintptrUint16R, (*decFnInfo).fastpathDecMapUintptrUint16R)
	fn(map[uintptr]uint32(nil), (*encFnInfo).fastpathEncMapUintptrUint32R, (*decFnInfo).fastpathDecMapUintptrUint32R)
	fn(map[uintptr]uint64(nil), (*encFnInfo).fastpathEncMapUintptrUint64R, (*decFnInfo).fastpathDecMapUintptrUint64R)
	fn(map[uintptr]uintptr(nil), (*encFnInfo).fastpathEncMapUintptrUintptrR, (*decFnInfo).fastpathDecMapUintptrUintptrR)
	fn(map[uintptr]int(nil), (*encFnInfo).fastpathEncMapUintptrIntR, (*decFnInfo).fastpathDecMapUintptrIntR)
	fn(map[uintptr]int8(nil), (*encFnInfo).fastpathEncMapUintptrInt8R, (*decFnInfo).fastpathDecMapUintptrInt8R)
	fn(map[uintptr]int16(nil), (*encFnInfo).fastpathEncMapUintptrInt16R, (*decFnInfo).fastpathDecMapUintptrInt16R)
	fn(map[uintptr]int32(nil), (*encFnInfo).fastpathEncMapUintptrInt32R, (*decFnInfo).fastpathDecMapUintptrInt32R)
	fn(map[uintptr]int64(nil), (*encFnInfo).fastpathEncMapUintptrInt64R, (*decFnInfo).fastpathDecMapUintptrInt64R)
	fn(map[uintptr]float32(nil), (*encFnInfo).fastpathEncMapUintptrFloat32R, (*decFnInfo).fastpathDecMapUintptrFloat32R)
	fn(map[uintptr]float64(nil), (*encFnInfo).fastpathEncMapUintptrFloat64R, (*decFnInfo).fastpathDecMapUintptrFloat64R)
	fn(map[uintptr]bool(nil), (*encFnInfo).fastpathEncMapUintptrBoolR, (*decFnInfo).fastpathDecMapUintptrBoolR)
	fn(map[int]interface{}(nil), (*encFnInfo).fastpathEncMapIntIntfR, (*decFnInfo).fastpathDecMapIntIntfR)
	fn(map[int]string(nil), (*encFnInfo).fastpathEncMapIntStringR, (*decFnInfo).fastpathDecMapIntStringR)
	fn(map[int]uint(nil), (*encFnInfo).fastpathEncMapIntUintR, (*decFnInfo).fastpathDecMapIntUintR)
	fn(map[int]uint8(nil), (*encFnInfo).fastpathEncMapIntUint8R, (*decFnInfo).fastpathDecMapIntUint8R)
	fn(map[int]uint16(nil), (*encFnInfo).fastpathEncMapIntUint16R, (*decFnInfo).fastpathDecMapIntUint16R)
	fn(map[int]uint32(nil), (*encFnInfo).fastpathEncMapIntUint32R, (*decFnInfo).fastpathDecMapIntUint32R)
	fn(map[int]uint64(nil), (*encFnInfo).fastpathEncMapIntUint64R, (*decFnInfo).fastpathDecMapIntUint64R)
	fn(map[int]uintptr(nil), (*encFnInfo).fastpathEncMapIntUintptrR, (*decFnInfo).fastpathDecMapIntUintptrR)
	fn(map[int]int(nil), (*encFnInfo).fastpathEncMapIntIntR, (*decFnInfo).fastpathDecMapIntIntR)
	fn(map[int]int8(nil), (*encFnInfo).fastpathEncMapIntInt8R, (*decFnInfo).fastpathDecMapIntInt8R)
	fn(map[int]int16(nil), (*encFnInfo).fastpathEncMapIntInt16R, (*decFnInfo).fastpathDecMapIntInt16R)
	fn(map[int]int32(nil), (*encFnInfo).fastpathEncMapIntInt32R, (*decFnInfo).fastpathDecMapIntInt32R)
	fn(map[int]int64(nil), (*encFnInfo).fastpathEncMapIntInt64R, (*decFnInfo).fastpathDecMapIntInt64R)
	fn(map[int]float32(nil), (*encFnInfo).fastpathEncMapIntFloat32R, (*decFnInfo).fastpathDecMapIntFloat32R)
	fn(map[int]float64(nil), (*encFnInfo).fastpathEncMapIntFloat64R, (*decFnInfo).fastpathDecMapIntFloat64R)
	fn(map[int]bool(nil), (*encFnInfo).fastpathEncMapIntBoolR, (*decFnInfo).fastpathDecMapIntBoolR)
	fn(map[int8]interface{}(nil), (*encFnInfo).fastpathEncMapInt8IntfR, (*decFnInfo).fastpathDecMapInt8IntfR)
	fn(map[int8]string(nil), (*encFnInfo).fastpathEncMapInt8StringR, (*decFnInfo).fastpathDecMapInt8StringR)
	fn(map[int8]uint(nil), (*encFnInfo).fastpathEncMapInt8UintR, (*decFnInfo).fastpathDecMapInt8UintR)
	fn(map[int8]uint8(nil), (*encFnInfo).fastpathEncMapInt8Uint8R, (*decFnInfo).fastpathDecMapInt8Uint8R)
	fn(map[int8]uint16(nil), (*encFnInfo).fastpathEncMapInt8Uint16R, (*decFnInfo).fastpathDecMapInt8Uint16R)
	fn(map[int8]uint32(nil), (*encFnInfo).fastpathEncMapInt8Uint32R, (*decFnInfo).fastpathDecMapInt8Uint32R)
	fn(map[int8]uint64(nil), (*encFnInfo).fastpathEncMapInt8Uint64R, (*decFnInfo).fastpathDecMapInt8Uint64R)
	fn(map[int8]uintptr(nil), (*encFnInfo).fastpathEncMapInt8UintptrR, (*decFnInfo).fastpathDecMapInt8UintptrR)
	fn(map[int8]int(nil), (*encFnInfo).fastpathEncMapInt8IntR, (*decFnInfo).fastpathDecMapInt8IntR)
	fn(map[int8]int8(nil), (*encFnInfo).fastpathEncMapInt8Int8R, (*decFnInfo).fastpathDecMapInt8Int8R)
	fn(map[int8]int16(nil), (*encFnInfo).fastpathEncMapInt8Int16R, (*decFnInfo).fastpathDecMapInt8Int16R)
	fn(map[int8]int32(nil), (*encFnInfo).fastpathEncMapInt8Int32R, (*decFnInfo).fastpathDecMapInt8Int32R)
	fn(map[int8]int64(nil), (*encFnInfo).fastpathEncMapInt8Int64R, (*decFnInfo).fastpathDecMapInt8Int64R)
	fn(map[int8]float32(nil), (*encFnInfo).fastpathEncMapInt8Float32R, (*decFnInfo).fastpathDecMapInt8Float32R)
	fn(map[int8]float64(nil), (*encFnInfo).fastpathEncMapInt8Float64R, (*decFnInfo).fastpathDecMapInt8Float64R)
	fn(map[int8]bool(nil), (*encFnInfo).fastpathEncMapInt8BoolR, (*decFnInfo).fastpathDecMapInt8BoolR)
	fn(map[int16]interface{}(nil), (*encFnInfo).fastpathEncMapInt16IntfR, (*decFnInfo).fastpathDecMapInt16IntfR)
	fn(map[int16]string(nil), (*encFnInfo).fastpathEncMapInt16StringR, (*decFnInfo).fastpathDecMapInt16StringR)
	fn(map[int16]uint(nil), (*encFnInfo).fastpathEncMapInt16UintR, (*decFnInfo).fastpathDecMapInt16UintR)
	fn(map[int16]uint8(nil), (*encFnInfo).fastpathEncMapInt16Uint8R, (*decFnInfo).fastpathDecMapInt16Uint8R)
	fn(map[int16]uint16(nil), (*encFnInfo).fastpathEncMapInt16Uint16R, (*decFnInfo).fastpathDecMapInt16Uint16R)
	fn(map[int16]uint32(nil), (*encFnInfo).fastpathEncMapInt16Uint32R, (*decFnInfo).fastpathDecMapInt16Uint32R)
	fn(map[int16]uint64(nil), (*encFnInfo).fastpathEncMapInt16Uint64R, (*decFnInfo).fastpathDecMapInt16Uint64R)
	fn(map[int16]uintptr(nil), (*encFnInfo).fastpathEncMapInt16UintptrR, (*decFnInfo).fastpathDecMapInt16UintptrR)
	fn(map[int16]int(nil), (*encFnInfo).fastpathEncMapInt16IntR, (*decFnInfo).fastpathDecMapInt16IntR)
	fn(map[int16]int8(nil), (*encFnInfo).fastpathEncMapInt16Int8R, (*decFnInfo).fastpathDecMapInt16Int8R)
	fn(map[int16]int16(nil), (*encFnInfo).fastpathEncMapInt16Int16R, (*decFnInfo).fastpathDecMapInt16Int16R)
	fn(map[int16]int32(nil), (*encFnInfo).fastpathEncMapInt16Int32R, (*decFnInfo).fastpathDecMapInt16Int32R)
	fn(map[int16]int64(nil), (*encFnInfo).fastpathEncMapInt16Int64R, (*decFnInfo).fastpathDecMapInt16Int64R)
	fn(map[int16]float32(nil), (*encFnInfo).fastpathEncMapInt16Float32R, (*decFnInfo).fastpathDecMapInt16Float32R)
	fn(map[int16]float64(nil), (*encFnInfo).fastpathEncMapInt16Float64R, (*decFnInfo).fastpathDecMapInt16Float64R)
	fn(map[int16]bool(nil), (*encFnInfo).fastpathEncMapInt16BoolR, (*decFnInfo).fastpathDecMapInt16BoolR)
	fn(map[int32]interface{}(nil), (*encFnInfo).fastpathEncMapInt32IntfR, (*decFnInfo).fastpathDecMapInt32IntfR)
	fn(map[int32]string(nil), (*encFnInfo).fastpathEncMapInt32StringR, (*decFnInfo).fastpathDecMapInt32StringR)
	fn(map[int32]uint(nil), (*encFnInfo).fastpathEncMapInt32UintR, (*decFnInfo).fastpathDecMapInt32UintR)
	fn(map[int32]uint8(nil), (*encFnInfo).fastpathEncMapInt32Uint8R, (*decFnInfo).fastpathDecMapInt32Uint8R)
	fn(map[int32]uint16(nil), (*encFnInfo).fastpathEncMapInt32Uint16R, (*decFnInfo).fastpathDecMapInt32Uint16R)
	fn(map[int32]uint32(nil), (*encFnInfo).fastpathEncMapInt32Uint32R, (*decFnInfo).fastpathDecMapInt32Uint32R)
	fn(map[int32]uint64(nil), (*encFnInfo).fastpathEncMapInt32Uint64R, (*decFnInfo).fastpathDecMapInt32Uint64R)
	fn(map[int32]uintptr(nil), (*encFnInfo).fastpathEncMapInt32UintptrR, (*decFnInfo).fastpathDecMapInt32UintptrR)
	fn(map[int32]int(nil), (*encFnInfo).fastpathEncMapInt32IntR, (*decFnInfo).fastpathDecMapInt32IntR)
	fn(map[int32]int8(nil), (*encFnInfo).fastpathEncMapInt32Int8R, (*decFnInfo).fastpathDecMapInt32Int8R)
	fn(map[int32]int16(nil), (*encFnInfo).fastpathEncMapInt32Int16R, (*decFnInfo).fastpathDecMapInt32Int16R)
	fn(map[int32]int32(nil), (*encFnInfo).fastpathEncMapInt32Int32R, (*decFnInfo).fastpathDecMapInt32Int32R)
	fn(map[int32]int64(nil), (*encFnInfo).fastpathEncMapInt32Int64R, (*decFnInfo).fastpathDecMapInt32Int64R)
	fn(map[int32]float32(nil), (*encFnInfo).fastpathEncMapInt32Float32R, (*decFnInfo).fastpathDecMapInt32Float32R)
	fn(map[int32]float64(nil), (*encFnInfo).fastpathEncMapInt32Float64R, (*decFnInfo).fastpathDecMapInt32Float64R)
	fn(map[int32]bool(nil), (*encFnInfo).fastpathEncMapInt32BoolR, (*decFnInfo).fastpathDecMapInt32BoolR)
	fn(map[int64]interface{}(nil), (*encFnInfo).fastpathEncMapInt64IntfR, (*decFnInfo).fastpathDecMapInt64IntfR)
	fn(map[int64]string(nil), (*encFnInfo).fastpathEncMapInt64StringR, (*decFnInfo).fastpathDecMapInt64StringR)
	fn(map[int64]uint(nil), (*encFnInfo).fastpathEncMapInt64UintR, (*decFnInfo).fastpathDecMapInt64UintR)
	fn(map[int64]uint8(nil), (*encFnInfo).fastpathEncMapInt64Uint8R, (*decFnInfo).fastpathDecMapInt64Uint8R)
	fn(map[int64]uint16(nil), (*encFnInfo).fastpathEncMapInt64Uint16R, (*decFnInfo).fastpathDecMapInt64Uint16R)
	fn(map[int64]uint32(nil), (*encFnInfo).fastpathEncMapInt64Uint32R, (*decFnInfo).fastpathDecMapInt64Uint32R)
	fn(map[int64]uint64(nil), (*encFnInfo).fastpathEncMapInt64Uint64R, (*decFnInfo).fastpathDecMapInt64Uint64R)
	fn(map[int64]uintptr(nil), (*encFnInfo).fastpathEncMapInt64UintptrR, (*decFnInfo).fastpathDecMapInt64UintptrR)
	fn(map[int64]int(nil), (*encFnInfo).fastpathEncMapInt64IntR, (*decFnInfo).fastpathDecMapInt64IntR)
	fn(map[int64]int8(nil), (*encFnInfo).fastpathEncMapInt64Int8R, (*decFnInfo).fastpathDecMapInt64Int8R)
	fn(map[int64]int16(nil), (*encFnInfo).fastpathEncMapInt64Int16R, (*decFnInfo).fastpathDecMapInt64Int16R)
	fn(map[int64]int32(nil), (*encFnInfo).fastpathEncMapInt64Int32R, (*decFnInfo).fastpathDecMapInt64Int32R)
	fn(map[int64]int64(nil), (*encFnInfo).fastpathEncMapInt64Int64R, (*decFnInfo).fastpathDecMapInt64Int64R)
	fn(map[int64]float32(nil), (*encFnInfo).fastpathEncMapInt64Float32R, (*decFnInfo).fastpathDecMapInt64Float32R)
	fn(map[int64]float64(nil), (*encFnInfo).fastpathEncMapInt64Float64R, (*decFnInfo).fastpathDecMapInt64Float64R)
	fn(map[int64]bool(nil), (*encFnInfo).fastpathEncMapInt64BoolR, (*decFnInfo).fastpathDecMapInt64BoolR)
	fn(map[bool]interface{}(nil), (*encFnInfo).fastpathEncMapBoolIntfR, (*decFnInfo).fastpathDecMapBoolIntfR)
	fn(map[bool]string(nil), (*encFnInfo).fastpathEncMapBoolStringR, (*decFnInfo).fastpathDecMapBoolStringR)
	fn(map[bool]uint(nil), (*encFnInfo).fastpathEncMapBoolUintR, (*decFnInfo).fastpathDecMapBoolUintR)
	fn(map[bool]uint8(nil), (*encFnInfo).fastpathEncMapBoolUint8R, (*decFnInfo).fastpathDecMapBoolUint8R)
	fn(map[bool]uint16(nil), (*encFnInfo).fastpathEncMapBoolUint16R, (*decFnInfo).fastpathDecMapBoolUint16R)
	fn(map[bool]uint32(nil), (*encFnInfo).fastpathEncMapBoolUint32R, (*decFnInfo).fastpathDecMapBoolUint32R)
	fn(map[bool]uint64(nil), (*encFnInfo).fastpathEncMapBoolUint64R, (*decFnInfo).fastpathDecMapBoolUint64R)
	fn(map[bool]uintptr(nil), (*encFnInfo).fastpathEncMapBoolUintptrR, (*decFnInfo).fastpathDecMapBoolUintptrR)
	fn(map[bool]int(nil), (*encFnInfo).fastpathEncMapBoolIntR, (*decFnInfo).fastpathDecMapBoolIntR)
	fn(map[bool]int8(nil), (*encFnInfo).fastpathEncMapBoolInt8R, (*decFnInfo).fastpathDecMapBoolInt8R)
	fn(map[bool]int16(nil), (*encFnInfo).fastpathEncMapBoolInt16R, (*decFnInfo).fastpathDecMapBoolInt16R)
	fn(map[bool]int32(nil), (*encFnInfo).fastpathEncMapBoolInt32R, (*decFnInfo).fastpathDecMapBoolInt32R)
	fn(map[bool]int64(nil), (*encFnInfo).fastpathEncMapBoolInt64R, (*decFnInfo).fastpathDecMapBoolInt64R)
	fn(map[bool]float32(nil), (*encFnInfo).fastpathEncMapBoolFloat32R, (*decFnInfo).fastpathDecMapBoolFloat32R)
	fn(map[bool]float64(nil), (*encFnInfo).fastpathEncMapBoolFloat64R, (*decFnInfo).fastpathDecMapBoolFloat64R)
	fn(map[bool]bool(nil), (*encFnInfo).fastpathEncMapBoolBoolR, (*decFnInfo).fastpathDecMapBoolBoolR)

	sort.Sort(fastpathAslice(fastpathAV[:]))
}

// -- encode

// -- -- fast path type switch
func fastpathEncodeTypeSwitch(iv interface{}, e *Encoder) bool {
	switch v := iv.(type) {

	case []interface{}:
		fastpathTV.EncSliceIntfV(v, fastpathCheckNilTrue, e)
	case *[]interface{}:
		fastpathTV.EncSliceIntfV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]interface{}:
		fastpathTV.EncMapIntfIntfV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]interface{}:
		fastpathTV.EncMapIntfIntfV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]string:
		fastpathTV.EncMapIntfStringV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]string:
		fastpathTV.EncMapIntfStringV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint:
		fastpathTV.EncMapIntfUintV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint:
		fastpathTV.EncMapIntfUintV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint8:
		fastpathTV.EncMapIntfUint8V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint8:
		fastpathTV.EncMapIntfUint8V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint16:
		fastpathTV.EncMapIntfUint16V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint16:
		fastpathTV.EncMapIntfUint16V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint32:
		fastpathTV.EncMapIntfUint32V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint32:
		fastpathTV.EncMapIntfUint32V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint64:
		fastpathTV.EncMapIntfUint64V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint64:
		fastpathTV.EncMapIntfUint64V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uintptr:
		fastpathTV.EncMapIntfUintptrV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uintptr:
		fastpathTV.EncMapIntfUintptrV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int:
		fastpathTV.EncMapIntfIntV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int:
		fastpathTV.EncMapIntfIntV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int8:
		fastpathTV.EncMapIntfInt8V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int8:
		fastpathTV.EncMapIntfInt8V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int16:
		fastpathTV.EncMapIntfInt16V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int16:
		fastpathTV.EncMapIntfInt16V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int32:
		fastpathTV.EncMapIntfInt32V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int32:
		fastpathTV.EncMapIntfInt32V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int64:
		fastpathTV.EncMapIntfInt64V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int64:
		fastpathTV.EncMapIntfInt64V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]float32:
		fastpathTV.EncMapIntfFloat32V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]float32:
		fastpathTV.EncMapIntfFloat32V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]float64:
		fastpathTV.EncMapIntfFloat64V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]float64:
		fastpathTV.EncMapIntfFloat64V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]bool:
		fastpathTV.EncMapIntfBoolV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]bool:
		fastpathTV.EncMapIntfBoolV(*v, fastpathCheckNilTrue, e)

	case []string:
		fastpathTV.EncSliceStringV(v, fastpathCheckNilTrue, e)
	case *[]string:
		fastpathTV.EncSliceStringV(*v, fastpathCheckNilTrue, e)

	case map[string]interface{}:
		fastpathTV.EncMapStringIntfV(v, fastpathCheckNilTrue, e)
	case *map[string]interface{}:
		fastpathTV.EncMapStringIntfV(*v, fastpathCheckNilTrue, e)

	case map[string]string:
		fastpathTV.EncMapStringStringV(v, fastpathCheckNilTrue, e)
	case *map[string]string:
		fastpathTV.EncMapStringStringV(*v, fastpathCheckNilTrue, e)

	case map[string]uint:
		fastpathTV.EncMapStringUintV(v, fastpathCheckNilTrue, e)
	case *map[string]uint:
		fastpathTV.EncMapStringUintV(*v, fastpathCheckNilTrue, e)

	case map[string]uint8:
		fastpathTV.EncMapStringUint8V(v, fastpathCheckNilTrue, e)
	case *map[string]uint8:
		fastpathTV.EncMapStringUint8V(*v, fastpathCheckNilTrue, e)

	case map[string]uint16:
		fastpathTV.EncMapStringUint16V(v, fastpathCheckNilTrue, e)
	case *map[string]uint16:
		fastpathTV.EncMapStringUint16V(*v, fastpathCheckNilTrue, e)

	case map[string]uint32:
		fastpathTV.EncMapStringUint32V(v, fastpathCheckNilTrue, e)
	case *map[string]uint32:
		fastpathTV.EncMapStringUint32V(*v, fastpathCheckNilTrue, e)

	case map[string]uint64:
		fastpathTV.EncMapStringUint64V(v, fastpathCheckNilTrue, e)
	case *map[string]uint64:
		fastpathTV.EncMapStringUint64V(*v, fastpathCheckNilTrue, e)

	case map[string]uintptr:
		fastpathTV.EncMapStringUintptrV(v, fastpathCheckNilTrue, e)
	case *map[string]uintptr:
		fastpathTV.EncMapStringUintptrV(*v, fastpathCheckNilTrue, e)

	case map[string]int:
		fastpathTV.EncMapStringIntV(v, fastpathCheckNilTrue, e)
	case *map[string]int:
		fastpathTV.EncMapStringIntV(*v, fastpathCheckNilTrue, e)

	case map[string]int8:
		fastpathTV.EncMapStringInt8V(v, fastpathCheckNilTrue, e)
	case *map[string]int8:
		fastpathTV.EncMapStringInt8V(*v, fastpathCheckNilTrue, e)

	case map[string]int16:
		fastpathTV.EncMapStringInt16V(v, fastpathCheckNilTrue, e)
	case *map[string]int16:
		fastpathTV.EncMapStringInt16V(*v, fastpathCheckNilTrue, e)

	case map[string]int32:
		fastpathTV.EncMapStringInt32V(v, fastpathCheckNilTrue, e)
	case *map[string]int32:
		fastpathTV.EncMapStringInt32V(*v, fastpathCheckNilTrue, e)

	case map[string]int64:
		fastpathTV.EncMapStringInt64V(v, fastpathCheckNilTrue, e)
	case *map[string]int64:
		fastpathTV.EncMapStringInt64V(*v, fastpathCheckNilTrue, e)

	case map[string]float32:
		fastpathTV.EncMapStringFloat32V(v, fastpathCheckNilTrue, e)
	case *map[string]float32:
		fastpathTV.EncMapStringFloat32V(*v, fastpathCheckNilTrue, e)

	case map[string]float64:
		fastpathTV.EncMapStringFloat64V(v, fastpathCheckNilTrue, e)
	case *map[string]float64:
		fastpathTV.EncMapStringFloat64V(*v, fastpathCheckNilTrue, e)

	case map[string]bool:
		fastpathTV.EncMapStringBoolV(v, fastpathCheckNilTrue, e)
	case *map[string]bool:
		fastpathTV.EncMapStringBoolV(*v, fastpathCheckNilTrue, e)

	case []float32:
		fastpathTV.EncSliceFloat32V(v, fastpathCheckNilTrue, e)
	case *[]float32:
		fastpathTV.EncSliceFloat32V(*v, fastpathCheckNilTrue, e)

	case map[float32]interface{}:
		fastpathTV.EncMapFloat32IntfV(v, fastpathCheckNilTrue, e)
	case *map[float32]interface{}:
		fastpathTV.EncMapFloat32IntfV(*v, fastpathCheckNilTrue, e)

	case map[float32]string:
		fastpathTV.EncMapFloat32StringV(v, fastpathCheckNilTrue, e)
	case *map[float32]string:
		fastpathTV.EncMapFloat32StringV(*v, fastpathCheckNilTrue, e)

	case map[float32]uint:
		fastpathTV.EncMapFloat32UintV(v, fastpathCheckNilTrue, e)
	case *map[float32]uint:
		fastpathTV.EncMapFloat32UintV(*v, fastpathCheckNilTrue, e)

	case map[float32]uint8:
		fastpathTV.EncMapFloat32Uint8V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint8:
		fastpathTV.EncMapFloat32Uint8V(*v, fastpathCheckNilTrue, e)

	case map[float32]uint16:
		fastpathTV.EncMapFloat32Uint16V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint16:
		fastpathTV.EncMapFloat32Uint16V(*v, fastpathCheckNilTrue, e)

	case map[float32]uint32:
		fastpathTV.EncMapFloat32Uint32V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint32:
		fastpathTV.EncMapFloat32Uint32V(*v, fastpathCheckNilTrue, e)

	case map[float32]uint64:
		fastpathTV.EncMapFloat32Uint64V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint64:
		fastpathTV.EncMapFloat32Uint64V(*v, fastpathCheckNilTrue, e)

	case map[float32]uintptr:
		fastpathTV.EncMapFloat32UintptrV(v, fastpathCheckNilTrue, e)
	case *map[float32]uintptr:
		fastpathTV.EncMapFloat32UintptrV(*v, fastpathCheckNilTrue, e)

	case map[float32]int:
		fastpathTV.EncMapFloat32IntV(v, fastpathCheckNilTrue, e)
	case *map[float32]int:
		fastpathTV.EncMapFloat32IntV(*v, fastpathCheckNilTrue, e)

	case map[float32]int8:
		fastpathTV.EncMapFloat32Int8V(v, fastpathCheckNilTrue, e)
	case *map[float32]int8:
		fastpathTV.EncMapFloat32Int8V(*v, fastpathCheckNilTrue, e)

	case map[float32]int16:
		fastpathTV.EncMapFloat32Int16V(v, fastpathCheckNilTrue, e)
	case *map[float32]int16:
		fastpathTV.EncMapFloat32Int16V(*v, fastpathCheckNilTrue, e)

	case map[float32]int32:
		fastpathTV.EncMapFloat32Int32V(v, fastpathCheckNilTrue, e)
	case *map[float32]int32:
		fastpathTV.EncMapFloat32Int32V(*v, fastpathCheckNilTrue, e)

	case map[float32]int64:
		fastpathTV.EncMapFloat32Int64V(v, fastpathCheckNilTrue, e)
	case *map[float32]int64:
		fastpathTV.EncMapFloat32Int64V(*v, fastpathCheckNilTrue, e)

	case map[float32]float32:
		fastpathTV.EncMapFloat32Float32V(v, fastpathCheckNilTrue, e)
	case *map[float32]float32:
		fastpathTV.EncMapFloat32Float32V(*v, fastpathCheckNilTrue, e)

	case map[float32]float64:
		fastpathTV.EncMapFloat32Float64V(v, fastpathCheckNilTrue, e)
	case *map[float32]float64:
		fastpathTV.EncMapFloat32Float64V(*v, fastpathCheckNilTrue, e)

	case map[float32]bool:
		fastpathTV.EncMapFloat32BoolV(v, fastpathCheckNilTrue, e)
	case *map[float32]bool:
		fastpathTV.EncMapFloat32BoolV(*v, fastpathCheckNilTrue, e)

	case []float64:
		fastpathTV.EncSliceFloat64V(v, fastpathCheckNilTrue, e)
	case *[]float64:
		fastpathTV.EncSliceFloat64V(*v, fastpathCheckNilTrue, e)

	case map[float64]interface{}:
		fastpathTV.EncMapFloat64IntfV(v, fastpathCheckNilTrue, e)
	case *map[float64]interface{}:
		fastpathTV.EncMapFloat64IntfV(*v, fastpathCheckNilTrue, e)

	case map[float64]string:
		fastpathTV.EncMapFloat64StringV(v, fastpathCheckNilTrue, e)
	case *map[float64]string:
		fastpathTV.EncMapFloat64StringV(*v, fastpathCheckNilTrue, e)

	case map[float64]uint:
		fastpathTV.EncMapFloat64UintV(v, fastpathCheckNilTrue, e)
	case *map[float64]uint:
		fastpathTV.EncMapFloat64UintV(*v, fastpathCheckNilTrue, e)

	case map[float64]uint8:
		fastpathTV.EncMapFloat64Uint8V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint8:
		fastpathTV.EncMapFloat64Uint8V(*v, fastpathCheckNilTrue, e)

	case map[float64]uint16:
		fastpathTV.EncMapFloat64Uint16V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint16:
		fastpathTV.EncMapFloat64Uint16V(*v, fastpathCheckNilTrue, e)

	case map[float64]uint32:
		fastpathTV.EncMapFloat64Uint32V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint32:
		fastpathTV.EncMapFloat64Uint32V(*v, fastpathCheckNilTrue, e)

	case map[float64]uint64:
		fastpathTV.EncMapFloat64Uint64V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint64:
		fastpathTV.EncMapFloat64Uint64V(*v, fastpathCheckNilTrue, e)

	case map[float64]uintptr:
		fastpathTV.EncMapFloat64UintptrV(v, fastpathCheckNilTrue, e)
	case *map[float64]uintptr:
		fastpathTV.EncMapFloat64UintptrV(*v, fastpathCheckNilTrue, e)

	case map[float64]int:
		fastpathTV.EncMapFloat64IntV(v, fastpathCheckNilTrue, e)
	case *map[float64]int:
		fastpathTV.EncMapFloat64IntV(*v, fastpathCheckNilTrue, e)

	case map[float64]int8:
		fastpathTV.EncMapFloat64Int8V(v, fastpathCheckNilTrue, e)
	case *map[float64]int8:
		fastpathTV.EncMapFloat64Int8V(*v, fastpathCheckNilTrue, e)

	case map[float64]int16:
		fastpathTV.EncMapFloat64Int16V(v, fastpathCheckNilTrue, e)
	case *map[float64]int16:
		fastpathTV.EncMapFloat64Int16V(*v, fastpathCheckNilTrue, e)

	case map[float64]int32:
		fastpathTV.EncMapFloat64Int32V(v, fastpathCheckNilTrue, e)
	case *map[float64]int32:
		fastpathTV.EncMapFloat64Int32V(*v, fastpathCheckNilTrue, e)

	case map[float64]int64:
		fastpathTV.EncMapFloat64Int64V(v, fastpathCheckNilTrue, e)
	case *map[float64]int64:
		fastpathTV.EncMapFloat64Int64V(*v, fastpathCheckNilTrue, e)

	case map[float64]float32:
		fastpathTV.EncMapFloat64Float32V(v, fastpathCheckNilTrue, e)
	case *map[float64]float32:
		fastpathTV.EncMapFloat64Float32V(*v, fastpathCheckNilTrue, e)

	case map[float64]float64:
		fastpathTV.EncMapFloat64Float64V(v, fastpathCheckNilTrue, e)
	case *map[float64]float64:
		fastpathTV.EncMapFloat64Float64V(*v, fastpathCheckNilTrue, e)

	case map[float64]bool:
		fastpathTV.EncMapFloat64BoolV(v, fastpathCheckNilTrue, e)
	case *map[float64]bool:
		fastpathTV.EncMapFloat64BoolV(*v, fastpathCheckNilTrue, e)

	case []uint:
		fastpathTV.EncSliceUintV(v, fastpathCheckNilTrue, e)
	case *[]uint:
		fastpathTV.EncSliceUintV(*v, fastpathCheckNilTrue, e)

	case map[uint]interface{}:
		fastpathTV.EncMapUintIntfV(v, fastpathCheckNilTrue, e)
	case *map[uint]interface{}:
		fastpathTV.EncMapUintIntfV(*v, fastpathCheckNilTrue, e)

	case map[uint]string:
		fastpathTV.EncMapUintStringV(v, fastpathCheckNilTrue, e)
	case *map[uint]string:
		fastpathTV.EncMapUintStringV(*v, fastpathCheckNilTrue, e)

	case map[uint]uint:
		fastpathTV.EncMapUintUintV(v, fastpathCheckNilTrue, e)
	case *map[uint]uint:
		fastpathTV.EncMapUintUintV(*v, fastpathCheckNilTrue, e)

	case map[uint]uint8:
		fastpathTV.EncMapUintUint8V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint8:
		fastpathTV.EncMapUintUint8V(*v, fastpathCheckNilTrue, e)

	case map[uint]uint16:
		fastpathTV.EncMapUintUint16V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint16:
		fastpathTV.EncMapUintUint16V(*v, fastpathCheckNilTrue, e)

	case map[uint]uint32:
		fastpathTV.EncMapUintUint32V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint32:
		fastpathTV.EncMapUintUint32V(*v, fastpathCheckNilTrue, e)

	case map[uint]uint64:
		fastpathTV.EncMapUintUint64V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint64:
		fastpathTV.EncMapUintUint64V(*v, fastpathCheckNilTrue, e)

	case map[uint]uintptr:
		fastpathTV.EncMapUintUintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint]uintptr:
		fastpathTV.EncMapUintUintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint]int:
		fastpathTV.EncMapUintIntV(v, fastpathCheckNilTrue, e)
	case *map[uint]int:
		fastpathTV.EncMapUintIntV(*v, fastpathCheckNilTrue, e)

	case map[uint]int8:
		fastpathTV.EncMapUintInt8V(v, fastpathCheckNilTrue, e)
	case *map[uint]int8:
		fastpathTV.EncMapUintInt8V(*v, fastpathCheckNilTrue, e)

	case map[uint]int16:
		fastpathTV.EncMapUintInt16V(v, fastpathCheckNilTrue, e)
	case *map[uint]int16:
		fastpathTV.EncMapUintInt16V(*v, fastpathCheckNilTrue, e)

	case map[uint]int32:
		fastpathTV.EncMapUintInt32V(v, fastpathCheckNilTrue, e)
	case *map[uint]int32:
		fastpathTV.EncMapUintInt32V(*v, fastpathCheckNilTrue, e)

	case map[uint]int64:
		fastpathTV.EncMapUintInt64V(v, fastpathCheckNilTrue, e)
	case *map[uint]int64:
		fastpathTV.EncMapUintInt64V(*v, fastpathCheckNilTrue, e)

	case map[uint]float32:
		fastpathTV.EncMapUintFloat32V(v, fastpathCheckNilTrue, e)
	case *map[uint]float32:
		fastpathTV.EncMapUintFloat32V(*v, fastpathCheckNilTrue, e)

	case map[uint]float64:
		fastpathTV.EncMapUintFloat64V(v, fastpathCheckNilTrue, e)
	case *map[uint]float64:
		fastpathTV.EncMapUintFloat64V(*v, fastpathCheckNilTrue, e)

	case map[uint]bool:
		fastpathTV.EncMapUintBoolV(v, fastpathCheckNilTrue, e)
	case *map[uint]bool:
		fastpathTV.EncMapUintBoolV(*v, fastpathCheckNilTrue, e)

	case map[uint8]interface{}:
		fastpathTV.EncMapUint8IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint8]interface{}:
		fastpathTV.EncMapUint8IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint8]string:
		fastpathTV.EncMapUint8StringV(v, fastpathCheckNilTrue, e)
	case *map[uint8]string:
		fastpathTV.EncMapUint8StringV(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint:
		fastpathTV.EncMapUint8UintV(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint:
		fastpathTV.EncMapUint8UintV(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint8:
		fastpathTV.EncMapUint8Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint8:
		fastpathTV.EncMapUint8Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint16:
		fastpathTV.EncMapUint8Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint16:
		fastpathTV.EncMapUint8Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint32:
		fastpathTV.EncMapUint8Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint32:
		fastpathTV.EncMapUint8Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint64:
		fastpathTV.EncMapUint8Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint64:
		fastpathTV.EncMapUint8Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uintptr:
		fastpathTV.EncMapUint8UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint8]uintptr:
		fastpathTV.EncMapUint8UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint8]int:
		fastpathTV.EncMapUint8IntV(v, fastpathCheckNilTrue, e)
	case *map[uint8]int:
		fastpathTV.EncMapUint8IntV(*v, fastpathCheckNilTrue, e)

	case map[uint8]int8:
		fastpathTV.EncMapUint8Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int8:
		fastpathTV.EncMapUint8Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint8]int16:
		fastpathTV.EncMapUint8Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int16:
		fastpathTV.EncMapUint8Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint8]int32:
		fastpathTV.EncMapUint8Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int32:
		fastpathTV.EncMapUint8Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint8]int64:
		fastpathTV.EncMapUint8Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int64:
		fastpathTV.EncMapUint8Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint8]float32:
		fastpathTV.EncMapUint8Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint8]float32:
		fastpathTV.EncMapUint8Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint8]float64:
		fastpathTV.EncMapUint8Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint8]float64:
		fastpathTV.EncMapUint8Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint8]bool:
		fastpathTV.EncMapUint8BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint8]bool:
		fastpathTV.EncMapUint8BoolV(*v, fastpathCheckNilTrue, e)

	case []uint16:
		fastpathTV.EncSliceUint16V(v, fastpathCheckNilTrue, e)
	case *[]uint16:
		fastpathTV.EncSliceUint16V(*v, fastpathCheckNilTrue, e)

	case map[uint16]interface{}:
		fastpathTV.EncMapUint16IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint16]interface{}:
		fastpathTV.EncMapUint16IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint16]string:
		fastpathTV.EncMapUint16StringV(v, fastpathCheckNilTrue, e)
	case *map[uint16]string:
		fastpathTV.EncMapUint16StringV(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint:
		fastpathTV.EncMapUint16UintV(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint:
		fastpathTV.EncMapUint16UintV(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint8:
		fastpathTV.EncMapUint16Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint8:
		fastpathTV.EncMapUint16Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint16:
		fastpathTV.EncMapUint16Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint16:
		fastpathTV.EncMapUint16Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint32:
		fastpathTV.EncMapUint16Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint32:
		fastpathTV.EncMapUint16Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint64:
		fastpathTV.EncMapUint16Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint64:
		fastpathTV.EncMapUint16Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uintptr:
		fastpathTV.EncMapUint16UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint16]uintptr:
		fastpathTV.EncMapUint16UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint16]int:
		fastpathTV.EncMapUint16IntV(v, fastpathCheckNilTrue, e)
	case *map[uint16]int:
		fastpathTV.EncMapUint16IntV(*v, fastpathCheckNilTrue, e)

	case map[uint16]int8:
		fastpathTV.EncMapUint16Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int8:
		fastpathTV.EncMapUint16Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint16]int16:
		fastpathTV.EncMapUint16Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int16:
		fastpathTV.EncMapUint16Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint16]int32:
		fastpathTV.EncMapUint16Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int32:
		fastpathTV.EncMapUint16Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint16]int64:
		fastpathTV.EncMapUint16Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int64:
		fastpathTV.EncMapUint16Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint16]float32:
		fastpathTV.EncMapUint16Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint16]float32:
		fastpathTV.EncMapUint16Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint16]float64:
		fastpathTV.EncMapUint16Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint16]float64:
		fastpathTV.EncMapUint16Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint16]bool:
		fastpathTV.EncMapUint16BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint16]bool:
		fastpathTV.EncMapUint16BoolV(*v, fastpathCheckNilTrue, e)

	case []uint32:
		fastpathTV.EncSliceUint32V(v, fastpathCheckNilTrue, e)
	case *[]uint32:
		fastpathTV.EncSliceUint32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]interface{}:
		fastpathTV.EncMapUint32IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint32]interface{}:
		fastpathTV.EncMapUint32IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint32]string:
		fastpathTV.EncMapUint32StringV(v, fastpathCheckNilTrue, e)
	case *map[uint32]string:
		fastpathTV.EncMapUint32StringV(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint:
		fastpathTV.EncMapUint32UintV(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint:
		fastpathTV.EncMapUint32UintV(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint8:
		fastpathTV.EncMapUint32Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint8:
		fastpathTV.EncMapUint32Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint16:
		fastpathTV.EncMapUint32Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint16:
		fastpathTV.EncMapUint32Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint32:
		fastpathTV.EncMapUint32Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint32:
		fastpathTV.EncMapUint32Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint64:
		fastpathTV.EncMapUint32Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint64:
		fastpathTV.EncMapUint32Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uintptr:
		fastpathTV.EncMapUint32UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint32]uintptr:
		fastpathTV.EncMapUint32UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint32]int:
		fastpathTV.EncMapUint32IntV(v, fastpathCheckNilTrue, e)
	case *map[uint32]int:
		fastpathTV.EncMapUint32IntV(*v, fastpathCheckNilTrue, e)

	case map[uint32]int8:
		fastpathTV.EncMapUint32Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int8:
		fastpathTV.EncMapUint32Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint32]int16:
		fastpathTV.EncMapUint32Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int16:
		fastpathTV.EncMapUint32Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint32]int32:
		fastpathTV.EncMapUint32Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int32:
		fastpathTV.EncMapUint32Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]int64:
		fastpathTV.EncMapUint32Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int64:
		fastpathTV.EncMapUint32Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint32]float32:
		fastpathTV.EncMapUint32Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint32]float32:
		fastpathTV.EncMapUint32Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]float64:
		fastpathTV.EncMapUint32Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint32]float64:
		fastpathTV.EncMapUint32Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint32]bool:
		fastpathTV.EncMapUint32BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint32]bool:
		fastpathTV.EncMapUint32BoolV(*v, fastpathCheckNilTrue, e)

	case []uint64:
		fastpathTV.EncSliceUint64V(v, fastpathCheckNilTrue, e)
	case *[]uint64:
		fastpathTV.EncSliceUint64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]interface{}:
		fastpathTV.EncMapUint64IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint64]interface{}:
		fastpathTV.EncMapUint64IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint64]string:
		fastpathTV.EncMapUint64StringV(v, fastpathCheckNilTrue, e)
	case *map[uint64]string:
		fastpathTV.EncMapUint64StringV(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint:
		fastpathTV.EncMapUint64UintV(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint:
		fastpathTV.EncMapUint64UintV(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint8:
		fastpathTV.EncMapUint64Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint8:
		fastpathTV.EncMapUint64Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint16:
		fastpathTV.EncMapUint64Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint16:
		fastpathTV.EncMapUint64Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint32:
		fastpathTV.EncMapUint64Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint32:
		fastpathTV.EncMapUint64Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint64:
		fastpathTV.EncMapUint64Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint64:
		fastpathTV.EncMapUint64Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uintptr:
		fastpathTV.EncMapUint64UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint64]uintptr:
		fastpathTV.EncMapUint64UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint64]int:
		fastpathTV.EncMapUint64IntV(v, fastpathCheckNilTrue, e)
	case *map[uint64]int:
		fastpathTV.EncMapUint64IntV(*v, fastpathCheckNilTrue, e)

	case map[uint64]int8:
		fastpathTV.EncMapUint64Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int8:
		fastpathTV.EncMapUint64Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint64]int16:
		fastpathTV.EncMapUint64Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int16:
		fastpathTV.EncMapUint64Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint64]int32:
		fastpathTV.EncMapUint64Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int32:
		fastpathTV.EncMapUint64Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint64]int64:
		fastpathTV.EncMapUint64Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int64:
		fastpathTV.EncMapUint64Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]float32:
		fastpathTV.EncMapUint64Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint64]float32:
		fastpathTV.EncMapUint64Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint64]float64:
		fastpathTV.EncMapUint64Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint64]float64:
		fastpathTV.EncMapUint64Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]bool:
		fastpathTV.EncMapUint64BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint64]bool:
		fastpathTV.EncMapUint64BoolV(*v, fastpathCheckNilTrue, e)

	case []uintptr:
		fastpathTV.EncSliceUintptrV(v, fastpathCheckNilTrue, e)
	case *[]uintptr:
		fastpathTV.EncSliceUintptrV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]interface{}:
		fastpathTV.EncMapUintptrIntfV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]interface{}:
		fastpathTV.EncMapUintptrIntfV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]string:
		fastpathTV.EncMapUintptrStringV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]string:
		fastpathTV.EncMapUintptrStringV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint:
		fastpathTV.EncMapUintptrUintV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint:
		fastpathTV.EncMapUintptrUintV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint8:
		fastpathTV.EncMapUintptrUint8V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint8:
		fastpathTV.EncMapUintptrUint8V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint16:
		fastpathTV.EncMapUintptrUint16V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint16:
		fastpathTV.EncMapUintptrUint16V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint32:
		fastpathTV.EncMapUintptrUint32V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint32:
		fastpathTV.EncMapUintptrUint32V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint64:
		fastpathTV.EncMapUintptrUint64V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint64:
		fastpathTV.EncMapUintptrUint64V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uintptr:
		fastpathTV.EncMapUintptrUintptrV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uintptr:
		fastpathTV.EncMapUintptrUintptrV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int:
		fastpathTV.EncMapUintptrIntV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int:
		fastpathTV.EncMapUintptrIntV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int8:
		fastpathTV.EncMapUintptrInt8V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int8:
		fastpathTV.EncMapUintptrInt8V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int16:
		fastpathTV.EncMapUintptrInt16V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int16:
		fastpathTV.EncMapUintptrInt16V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int32:
		fastpathTV.EncMapUintptrInt32V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int32:
		fastpathTV.EncMapUintptrInt32V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int64:
		fastpathTV.EncMapUintptrInt64V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int64:
		fastpathTV.EncMapUintptrInt64V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]float32:
		fastpathTV.EncMapUintptrFloat32V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]float32:
		fastpathTV.EncMapUintptrFloat32V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]float64:
		fastpathTV.EncMapUintptrFloat64V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]float64:
		fastpathTV.EncMapUintptrFloat64V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]bool:
		fastpathTV.EncMapUintptrBoolV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]bool:
		fastpathTV.EncMapUintptrBoolV(*v, fastpathCheckNilTrue, e)

	case []int:
		fastpathTV.EncSliceIntV(v, fastpathCheckNilTrue, e)
	case *[]int:
		fastpathTV.EncSliceIntV(*v, fastpathCheckNilTrue, e)

	case map[int]interface{}:
		fastpathTV.EncMapIntIntfV(v, fastpathCheckNilTrue, e)
	case *map[int]interface{}:
		fastpathTV.EncMapIntIntfV(*v, fastpathCheckNilTrue, e)

	case map[int]string:
		fastpathTV.EncMapIntStringV(v, fastpathCheckNilTrue, e)
	case *map[int]string:
		fastpathTV.EncMapIntStringV(*v, fastpathCheckNilTrue, e)

	case map[int]uint:
		fastpathTV.EncMapIntUintV(v, fastpathCheckNilTrue, e)
	case *map[int]uint:
		fastpathTV.EncMapIntUintV(*v, fastpathCheckNilTrue, e)

	case map[int]uint8:
		fastpathTV.EncMapIntUint8V(v, fastpathCheckNilTrue, e)
	case *map[int]uint8:
		fastpathTV.EncMapIntUint8V(*v, fastpathCheckNilTrue, e)

	case map[int]uint16:
		fastpathTV.EncMapIntUint16V(v, fastpathCheckNilTrue, e)
	case *map[int]uint16:
		fastpathTV.EncMapIntUint16V(*v, fastpathCheckNilTrue, e)

	case map[int]uint32:
		fastpathTV.EncMapIntUint32V(v, fastpathCheckNilTrue, e)
	case *map[int]uint32:
		fastpathTV.EncMapIntUint32V(*v, fastpathCheckNilTrue, e)

	case map[int]uint64:
		fastpathTV.EncMapIntUint64V(v, fastpathCheckNilTrue, e)
	case *map[int]uint64:
		fastpathTV.EncMapIntUint64V(*v, fastpathCheckNilTrue, e)

	case map[int]uintptr:
		fastpathTV.EncMapIntUintptrV(v, fastpathCheckNilTrue, e)
	case *map[int]uintptr:
		fastpathTV.EncMapIntUintptrV(*v, fastpathCheckNilTrue, e)

	case map[int]int:
		fastpathTV.EncMapIntIntV(v, fastpathCheckNilTrue, e)
	case *map[int]int:
		fastpathTV.EncMapIntIntV(*v, fastpathCheckNilTrue, e)

	case map[int]int8:
		fastpathTV.EncMapIntInt8V(v, fastpathCheckNilTrue, e)
	case *map[int]int8:
		fastpathTV.EncMapIntInt8V(*v, fastpathCheckNilTrue, e)

	case map[int]int16:
		fastpathTV.EncMapIntInt16V(v, fastpathCheckNilTrue, e)
	case *map[int]int16:
		fastpathTV.EncMapIntInt16V(*v, fastpathCheckNilTrue, e)

	case map[int]int32:
		fastpathTV.EncMapIntInt32V(v, fastpathCheckNilTrue, e)
	case *map[int]int32:
		fastpathTV.EncMapIntInt32V(*v, fastpathCheckNilTrue, e)

	case map[int]int64:
		fastpathTV.EncMapIntInt64V(v, fastpathCheckNilTrue, e)
	case *map[int]int64:
		fastpathTV.EncMapIntInt64V(*v, fastpathCheckNilTrue, e)

	case map[int]float32:
		fastpathTV.EncMapIntFloat32V(v, fastpathCheckNilTrue, e)
	case *map[int]float32:
		fastpathTV.EncMapIntFloat32V(*v, fastpathCheckNilTrue, e)

	case map[int]float64:
		fastpathTV.EncMapIntFloat64V(v, fastpathCheckNilTrue, e)
	case *map[int]float64:
		fastpathTV.EncMapIntFloat64V(*v, fastpathCheckNilTrue, e)

	case map[int]bool:
		fastpathTV.EncMapIntBoolV(v, fastpathCheckNilTrue, e)
	case *map[int]bool:
		fastpathTV.EncMapIntBoolV(*v, fastpathCheckNilTrue, e)

	case []int8:
		fastpathTV.EncSliceInt8V(v, fastpathCheckNilTrue, e)
	case *[]int8:
		fastpathTV.EncSliceInt8V(*v, fastpathCheckNilTrue, e)

	case map[int8]interface{}:
		fastpathTV.EncMapInt8IntfV(v, fastpathCheckNilTrue, e)
	case *map[int8]interface{}:
		fastpathTV.EncMapInt8IntfV(*v, fastpathCheckNilTrue, e)

	case map[int8]string:
		fastpathTV.EncMapInt8StringV(v, fastpathCheckNilTrue, e)
	case *map[int8]string:
		fastpathTV.EncMapInt8StringV(*v, fastpathCheckNilTrue, e)

	case map[int8]uint:
		fastpathTV.EncMapInt8UintV(v, fastpathCheckNilTrue, e)
	case *map[int8]uint:
		fastpathTV.EncMapInt8UintV(*v, fastpathCheckNilTrue, e)

	case map[int8]uint8:
		fastpathTV.EncMapInt8Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint8:
		fastpathTV.EncMapInt8Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int8]uint16:
		fastpathTV.EncMapInt8Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint16:
		fastpathTV.EncMapInt8Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int8]uint32:
		fastpathTV.EncMapInt8Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint32:
		fastpathTV.EncMapInt8Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int8]uint64:
		fastpathTV.EncMapInt8Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint64:
		fastpathTV.EncMapInt8Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int8]uintptr:
		fastpathTV.EncMapInt8UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int8]uintptr:
		fastpathTV.EncMapInt8UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int8]int:
		fastpathTV.EncMapInt8IntV(v, fastpathCheckNilTrue, e)
	case *map[int8]int:
		fastpathTV.EncMapInt8IntV(*v, fastpathCheckNilTrue, e)

	case map[int8]int8:
		fastpathTV.EncMapInt8Int8V(v, fastpathCheckNilTrue, e)
	case *map[int8]int8:
		fastpathTV.EncMapInt8Int8V(*v, fastpathCheckNilTrue, e)

	case map[int8]int16:
		fastpathTV.EncMapInt8Int16V(v, fastpathCheckNilTrue, e)
	case *map[int8]int16:
		fastpathTV.EncMapInt8Int16V(*v, fastpathCheckNilTrue, e)

	case map[int8]int32:
		fastpathTV.EncMapInt8Int32V(v, fastpathCheckNilTrue, e)
	case *map[int8]int32:
		fastpathTV.EncMapInt8Int32V(*v, fastpathCheckNilTrue, e)

	case map[int8]int64:
		fastpathTV.EncMapInt8Int64V(v, fastpathCheckNilTrue, e)
	case *map[int8]int64:
		fastpathTV.EncMapInt8Int64V(*v, fastpathCheckNilTrue, e)

	case map[int8]float32:
		fastpathTV.EncMapInt8Float32V(v, fastpathCheckNilTrue, e)
	case *map[int8]float32:
		fastpathTV.EncMapInt8Float32V(*v, fastpathCheckNilTrue, e)

	case map[int8]float64:
		fastpathTV.EncMapInt8Float64V(v, fastpathCheckNilTrue, e)
	case *map[int8]float64:
		fastpathTV.EncMapInt8Float64V(*v, fastpathCheckNilTrue, e)

	case map[int8]bool:
		fastpathTV.EncMapInt8BoolV(v, fastpathCheckNilTrue, e)
	case *map[int8]bool:
		fastpathTV.EncMapInt8BoolV(*v, fastpathCheckNilTrue, e)

	case []int16:
		fastpathTV.EncSliceInt16V(v, fastpathCheckNilTrue, e)
	case *[]int16:
		fastpathTV.EncSliceInt16V(*v, fastpathCheckNilTrue, e)

	case map[int16]interface{}:
		fastpathTV.EncMapInt16IntfV(v, fastpathCheckNilTrue, e)
	case *map[int16]interface{}:
		fastpathTV.EncMapInt16IntfV(*v, fastpathCheckNilTrue, e)

	case map[int16]string:
		fastpathTV.EncMapInt16StringV(v, fastpathCheckNilTrue, e)
	case *map[int16]string:
		fastpathTV.EncMapInt16StringV(*v, fastpathCheckNilTrue, e)

	case map[int16]uint:
		fastpathTV.EncMapInt16UintV(v, fastpathCheckNilTrue, e)
	case *map[int16]uint:
		fastpathTV.EncMapInt16UintV(*v, fastpathCheckNilTrue, e)

	case map[int16]uint8:
		fastpathTV.EncMapInt16Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint8:
		fastpathTV.EncMapInt16Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int16]uint16:
		fastpathTV.EncMapInt16Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint16:
		fastpathTV.EncMapInt16Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int16]uint32:
		fastpathTV.EncMapInt16Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint32:
		fastpathTV.EncMapInt16Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int16]uint64:
		fastpathTV.EncMapInt16Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint64:
		fastpathTV.EncMapInt16Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int16]uintptr:
		fastpathTV.EncMapInt16UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int16]uintptr:
		fastpathTV.EncMapInt16UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int16]int:
		fastpathTV.EncMapInt16IntV(v, fastpathCheckNilTrue, e)
	case *map[int16]int:
		fastpathTV.EncMapInt16IntV(*v, fastpathCheckNilTrue, e)

	case map[int16]int8:
		fastpathTV.EncMapInt16Int8V(v, fastpathCheckNilTrue, e)
	case *map[int16]int8:
		fastpathTV.EncMapInt16Int8V(*v, fastpathCheckNilTrue, e)

	case map[int16]int16:
		fastpathTV.EncMapInt16Int16V(v, fastpathCheckNilTrue, e)
	case *map[int16]int16:
		fastpathTV.EncMapInt16Int16V(*v, fastpathCheckNilTrue, e)

	case map[int16]int32:
		fastpathTV.EncMapInt16Int32V(v, fastpathCheckNilTrue, e)
	case *map[int16]int32:
		fastpathTV.EncMapInt16Int32V(*v, fastpathCheckNilTrue, e)

	case map[int16]int64:
		fastpathTV.EncMapInt16Int64V(v, fastpathCheckNilTrue, e)
	case *map[int16]int64:
		fastpathTV.EncMapInt16Int64V(*v, fastpathCheckNilTrue, e)

	case map[int16]float32:
		fastpathTV.EncMapInt16Float32V(v, fastpathCheckNilTrue, e)
	case *map[int16]float32:
		fastpathTV.EncMapInt16Float32V(*v, fastpathCheckNilTrue, e)

	case map[int16]float64:
		fastpathTV.EncMapInt16Float64V(v, fastpathCheckNilTrue, e)
	case *map[int16]float64:
		fastpathTV.EncMapInt16Float64V(*v, fastpathCheckNilTrue, e)

	case map[int16]bool:
		fastpathTV.EncMapInt16BoolV(v, fastpathCheckNilTrue, e)
	case *map[int16]bool:
		fastpathTV.EncMapInt16BoolV(*v, fastpathCheckNilTrue, e)

	case []int32:
		fastpathTV.EncSliceInt32V(v, fastpathCheckNilTrue, e)
	case *[]int32:
		fastpathTV.EncSliceInt32V(*v, fastpathCheckNilTrue, e)

	case map[int32]interface{}:
		fastpathTV.EncMapInt32IntfV(v, fastpathCheckNilTrue, e)
	case *map[int32]interface{}:
		fastpathTV.EncMapInt32IntfV(*v, fastpathCheckNilTrue, e)

	case map[int32]string:
		fastpathTV.EncMapInt32StringV(v, fastpathCheckNilTrue, e)
	case *map[int32]string:
		fastpathTV.EncMapInt32StringV(*v, fastpathCheckNilTrue, e)

	case map[int32]uint:
		fastpathTV.EncMapInt32UintV(v, fastpathCheckNilTrue, e)
	case *map[int32]uint:
		fastpathTV.EncMapInt32UintV(*v, fastpathCheckNilTrue, e)

	case map[int32]uint8:
		fastpathTV.EncMapInt32Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint8:
		fastpathTV.EncMapInt32Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int32]uint16:
		fastpathTV.EncMapInt32Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint16:
		fastpathTV.EncMapInt32Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int32]uint32:
		fastpathTV.EncMapInt32Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint32:
		fastpathTV.EncMapInt32Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int32]uint64:
		fastpathTV.EncMapInt32Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint64:
		fastpathTV.EncMapInt32Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int32]uintptr:
		fastpathTV.EncMapInt32UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int32]uintptr:
		fastpathTV.EncMapInt32UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int32]int:
		fastpathTV.EncMapInt32IntV(v, fastpathCheckNilTrue, e)
	case *map[int32]int:
		fastpathTV.EncMapInt32IntV(*v, fastpathCheckNilTrue, e)

	case map[int32]int8:
		fastpathTV.EncMapInt32Int8V(v, fastpathCheckNilTrue, e)
	case *map[int32]int8:
		fastpathTV.EncMapInt32Int8V(*v, fastpathCheckNilTrue, e)

	case map[int32]int16:
		fastpathTV.EncMapInt32Int16V(v, fastpathCheckNilTrue, e)
	case *map[int32]int16:
		fastpathTV.EncMapInt32Int16V(*v, fastpathCheckNilTrue, e)

	case map[int32]int32:
		fastpathTV.EncMapInt32Int32V(v, fastpathCheckNilTrue, e)
	case *map[int32]int32:
		fastpathTV.EncMapInt32Int32V(*v, fastpathCheckNilTrue, e)

	case map[int32]int64:
		fastpathTV.EncMapInt32Int64V(v, fastpathCheckNilTrue, e)
	case *map[int32]int64:
		fastpathTV.EncMapInt32Int64V(*v, fastpathCheckNilTrue, e)

	case map[int32]float32:
		fastpathTV.EncMapInt32Float32V(v, fastpathCheckNilTrue, e)
	case *map[int32]float32:
		fastpathTV.EncMapInt32Float32V(*v, fastpathCheckNilTrue, e)

	case map[int32]float64:
		fastpathTV.EncMapInt32Float64V(v, fastpathCheckNilTrue, e)
	case *map[int32]float64:
		fastpathTV.EncMapInt32Float64V(*v, fastpathCheckNilTrue, e)

	case map[int32]bool:
		fastpathTV.EncMapInt32BoolV(v, fastpathCheckNilTrue, e)
	case *map[int32]bool:
		fastpathTV.EncMapInt32BoolV(*v, fastpathCheckNilTrue, e)

	case []int64:
		fastpathTV.EncSliceInt64V(v, fastpathCheckNilTrue, e)
	case *[]int64:
		fastpathTV.EncSliceInt64V(*v, fastpathCheckNilTrue, e)

	case map[int64]interface{}:
		fastpathTV.EncMapInt64IntfV(v, fastpathCheckNilTrue, e)
	case *map[int64]interface{}:
		fastpathTV.EncMapInt64IntfV(*v, fastpathCheckNilTrue, e)

	case map[int64]string:
		fastpathTV.EncMapInt64StringV(v, fastpathCheckNilTrue, e)
	case *map[int64]string:
		fastpathTV.EncMapInt64StringV(*v, fastpathCheckNilTrue, e)

	case map[int64]uint:
		fastpathTV.EncMapInt64UintV(v, fastpathCheckNilTrue, e)
	case *map[int64]uint:
		fastpathTV.EncMapInt64UintV(*v, fastpathCheckNilTrue, e)

	case map[int64]uint8:
		fastpathTV.EncMapInt64Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint8:
		fastpathTV.EncMapInt64Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int64]uint16:
		fastpathTV.EncMapInt64Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint16:
		fastpathTV.EncMapInt64Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int64]uint32:
		fastpathTV.EncMapInt64Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint32:
		fastpathTV.EncMapInt64Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int64]uint64:
		fastpathTV.EncMapInt64Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint64:
		fastpathTV.EncMapInt64Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int64]uintptr:
		fastpathTV.EncMapInt64UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int64]uintptr:
		fastpathTV.EncMapInt64UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int64]int:
		fastpathTV.EncMapInt64IntV(v, fastpathCheckNilTrue, e)
	case *map[int64]int:
		fastpathTV.EncMapInt64IntV(*v, fastpathCheckNilTrue, e)

	case map[int64]int8:
		fastpathTV.EncMapInt64Int8V(v, fastpathCheckNilTrue, e)
	case *map[int64]int8:
		fastpathTV.EncMapInt64Int8V(*v, fastpathCheckNilTrue, e)

	case map[int64]int16:
		fastpathTV.EncMapInt64Int16V(v, fastpathCheckNilTrue, e)
	case *map[int64]int16:
		fastpathTV.EncMapInt64Int16V(*v, fastpathCheckNilTrue, e)

	case map[int64]int32:
		fastpathTV.EncMapInt64Int32V(v, fastpathCheckNilTrue, e)
	case *map[int64]int32:
		fastpathTV.EncMapInt64Int32V(*v, fastpathCheckNilTrue, e)

	case map[int64]int64:
		fastpathTV.EncMapInt64Int64V(v, fastpathCheckNilTrue, e)
	case *map[int64]int64:
		fastpathTV.EncMapInt64Int64V(*v, fastpathCheckNilTrue, e)

	case map[int64]float32:
		fastpathTV.EncMapInt64Float32V(v, fastpathCheckNilTrue, e)
	case *map[int64]float32:
		fastpathTV.EncMapInt64Float32V(*v, fastpathCheckNilTrue, e)

	case map[int64]float64:
		fastpathTV.EncMapInt64Float64V(v, fastpathCheckNilTrue, e)
	case *map[int64]float64:
		fastpathTV.EncMapInt64Float64V(*v, fastpathCheckNilTrue, e)

	case map[int64]bool:
		fastpathTV.EncMapInt64BoolV(v, fastpathCheckNilTrue, e)
	case *map[int64]bool:
		fastpathTV.EncMapInt64BoolV(*v, fastpathCheckNilTrue, e)

	case []bool:
		fastpathTV.EncSliceBoolV(v, fastpathCheckNilTrue, e)
	case *[]bool:
		fastpathTV.EncSliceBoolV(*v, fastpathCheckNilTrue, e)

	case map[bool]interface{}:
		fastpathTV.EncMapBoolIntfV(v, fastpathCheckNilTrue, e)
	case *map[bool]interface{}:
		fastpathTV.EncMapBoolIntfV(*v, fastpathCheckNilTrue, e)

	case map[bool]string:
		fastpathTV.EncMapBoolStringV(v, fastpathCheckNilTrue, e)
	case *map[bool]string:
		fastpathTV.EncMapBoolStringV(*v, fastpathCheckNilTrue, e)

	case map[bool]uint:
		fastpathTV.EncMapBoolUintV(v, fastpathCheckNilTrue, e)
	case *map[bool]uint:
		fastpathTV.EncMapBoolUintV(*v, fastpathCheckNilTrue, e)

	case map[bool]uint8:
		fastpathTV.EncMapBoolUint8V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint8:
		fastpathTV.EncMapBoolUint8V(*v, fastpathCheckNilTrue, e)

	case map[bool]uint16:
		fastpathTV.EncMapBoolUint16V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint16:
		fastpathTV.EncMapBoolUint16V(*v, fastpathCheckNilTrue, e)

	case map[bool]uint32:
		fastpathTV.EncMapBoolUint32V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint32:
		fastpathTV.EncMapBoolUint32V(*v, fastpathCheckNilTrue, e)

	case map[bool]uint64:
		fastpathTV.EncMapBoolUint64V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint64:
		fastpathTV.EncMapBoolUint64V(*v, fastpathCheckNilTrue, e)

	case map[bool]uintptr:
		fastpathTV.EncMapBoolUintptrV(v, fastpathCheckNilTrue, e)
	case *map[bool]uintptr:
		fastpathTV.EncMapBoolUintptrV(*v, fastpathCheckNilTrue, e)

	case map[bool]int:
		fastpathTV.EncMapBoolIntV(v, fastpathCheckNilTrue, e)
	case *map[bool]int:
		fastpathTV.EncMapBoolIntV(*v, fastpathCheckNilTrue, e)

	case map[bool]int8:
		fastpathTV.EncMapBoolInt8V(v, fastpathCheckNilTrue, e)
	case *map[bool]int8:
		fastpathTV.EncMapBoolInt8V(*v, fastpathCheckNilTrue, e)

	case map[bool]int16:
		fastpathTV.EncMapBoolInt16V(v, fastpathCheckNilTrue, e)
	case *map[bool]int16:
		fastpathTV.EncMapBoolInt16V(*v, fastpathCheckNilTrue, e)

	case map[bool]int32:
		fastpathTV.EncMapBoolInt32V(v, fastpathCheckNilTrue, e)
	case *map[bool]int32:
		fastpathTV.EncMapBoolInt32V(*v, fastpathCheckNilTrue, e)

	case map[bool]int64:
		fastpathTV.EncMapBoolInt64V(v, fastpathCheckNilTrue, e)
	case *map[bool]int64:
		fastpathTV.EncMapBoolInt64V(*v, fastpathCheckNilTrue, e)

	case map[bool]float32:
		fastpathTV.EncMapBoolFloat32V(v, fastpathCheckNilTrue, e)
	case *map[bool]float32:
		fastpathTV.EncMapBoolFloat32V(*v, fastpathCheckNilTrue, e)

	case map[bool]float64:
		fastpathTV.EncMapBoolFloat64V(v, fastpathCheckNilTrue, e)
	case *map[bool]float64:
		fastpathTV.EncMapBoolFloat64V(*v, fastpathCheckNilTrue, e)

	case map[bool]bool:
		fastpathTV.EncMapBoolBoolV(v, fastpathCheckNilTrue, e)
	case *map[bool]bool:
		fastpathTV.EncMapBoolBoolV(*v, fastpathCheckNilTrue, e)

	default:
		_ = v // TODO: workaround https://github.com/golang/go/issues/12927 (remove after go 1.6 release)
		return false
	}
	return true
}

func fastpathEncodeTypeSwitchSlice(iv interface{}, e *Encoder) bool {
	switch v := iv.(type) {

	case []interface{}:
		fastpathTV.EncSliceIntfV(v, fastpathCheckNilTrue, e)
	case *[]interface{}:
		fastpathTV.EncSliceIntfV(*v, fastpathCheckNilTrue, e)

	case []string:
		fastpathTV.EncSliceStringV(v, fastpathCheckNilTrue, e)
	case *[]string:
		fastpathTV.EncSliceStringV(*v, fastpathCheckNilTrue, e)

	case []float32:
		fastpathTV.EncSliceFloat32V(v, fastpathCheckNilTrue, e)
	case *[]float32:
		fastpathTV.EncSliceFloat32V(*v, fastpathCheckNilTrue, e)

	case []float64:
		fastpathTV.EncSliceFloat64V(v, fastpathCheckNilTrue, e)
	case *[]float64:
		fastpathTV.EncSliceFloat64V(*v, fastpathCheckNilTrue, e)

	case []uint:
		fastpathTV.EncSliceUintV(v, fastpathCheckNilTrue, e)
	case *[]uint:
		fastpathTV.EncSliceUintV(*v, fastpathCheckNilTrue, e)

	case []uint16:
		fastpathTV.EncSliceUint16V(v, fastpathCheckNilTrue, e)
	case *[]uint16:
		fastpathTV.EncSliceUint16V(*v, fastpathCheckNilTrue, e)

	case []uint32:
		fastpathTV.EncSliceUint32V(v, fastpathCheckNilTrue, e)
	case *[]uint32:
		fastpathTV.EncSliceUint32V(*v, fastpathCheckNilTrue, e)

	case []uint64:
		fastpathTV.EncSliceUint64V(v, fastpathCheckNilTrue, e)
	case *[]uint64:
		fastpathTV.EncSliceUint64V(*v, fastpathCheckNilTrue, e)

	case []uintptr:
		fastpathTV.EncSliceUintptrV(v, fastpathCheckNilTrue, e)
	case *[]uintptr:
		fastpathTV.EncSliceUintptrV(*v, fastpathCheckNilTrue, e)

	case []int:
		fastpathTV.EncSliceIntV(v, fastpathCheckNilTrue, e)
	case *[]int:
		fastpathTV.EncSliceIntV(*v, fastpathCheckNilTrue, e)

	case []int8:
		fastpathTV.EncSliceInt8V(v, fastpathCheckNilTrue, e)
	case *[]int8:
		fastpathTV.EncSliceInt8V(*v, fastpathCheckNilTrue, e)

	case []int16:
		fastpathTV.EncSliceInt16V(v, fastpathCheckNilTrue, e)
	case *[]int16:
		fastpathTV.EncSliceInt16V(*v, fastpathCheckNilTrue, e)

	case []int32:
		fastpathTV.EncSliceInt32V(v, fastpathCheckNilTrue, e)
	case *[]int32:
		fastpathTV.EncSliceInt32V(*v, fastpathCheckNilTrue, e)

	case []int64:
		fastpathTV.EncSliceInt64V(v, fastpathCheckNilTrue, e)
	case *[]int64:
		fastpathTV.EncSliceInt64V(*v, fastpathCheckNilTrue, e)

	case []bool:
		fastpathTV.EncSliceBoolV(v, fastpathCheckNilTrue, e)
	case *[]bool:
		fastpathTV.EncSliceBoolV(*v, fastpathCheckNilTrue, e)

	default:
		_ = v // TODO: workaround https://github.com/golang/go/issues/12927 (remove after go 1.6 release)
		return false
	}
	return true
}

func fastpathEncodeTypeSwitchMap(iv interface{}, e *Encoder) bool {
	switch v := iv.(type) {

	case map[interface{}]interface{}:
		fastpathTV.EncMapIntfIntfV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]interface{}:
		fastpathTV.EncMapIntfIntfV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]string:
		fastpathTV.EncMapIntfStringV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]string:
		fastpathTV.EncMapIntfStringV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint:
		fastpathTV.EncMapIntfUintV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint:
		fastpathTV.EncMapIntfUintV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint8:
		fastpathTV.EncMapIntfUint8V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint8:
		fastpathTV.EncMapIntfUint8V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint16:
		fastpathTV.EncMapIntfUint16V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint16:
		fastpathTV.EncMapIntfUint16V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint32:
		fastpathTV.EncMapIntfUint32V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint32:
		fastpathTV.EncMapIntfUint32V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uint64:
		fastpathTV.EncMapIntfUint64V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uint64:
		fastpathTV.EncMapIntfUint64V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]uintptr:
		fastpathTV.EncMapIntfUintptrV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]uintptr:
		fastpathTV.EncMapIntfUintptrV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int:
		fastpathTV.EncMapIntfIntV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int:
		fastpathTV.EncMapIntfIntV(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int8:
		fastpathTV.EncMapIntfInt8V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int8:
		fastpathTV.EncMapIntfInt8V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int16:
		fastpathTV.EncMapIntfInt16V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int16:
		fastpathTV.EncMapIntfInt16V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int32:
		fastpathTV.EncMapIntfInt32V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int32:
		fastpathTV.EncMapIntfInt32V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]int64:
		fastpathTV.EncMapIntfInt64V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]int64:
		fastpathTV.EncMapIntfInt64V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]float32:
		fastpathTV.EncMapIntfFloat32V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]float32:
		fastpathTV.EncMapIntfFloat32V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]float64:
		fastpathTV.EncMapIntfFloat64V(v, fastpathCheckNilTrue, e)
	case *map[interface{}]float64:
		fastpathTV.EncMapIntfFloat64V(*v, fastpathCheckNilTrue, e)

	case map[interface{}]bool:
		fastpathTV.EncMapIntfBoolV(v, fastpathCheckNilTrue, e)
	case *map[interface{}]bool:
		fastpathTV.EncMapIntfBoolV(*v, fastpathCheckNilTrue, e)

	case map[string]interface{}:
		fastpathTV.EncMapStringIntfV(v, fastpathCheckNilTrue, e)
	case *map[string]interface{}:
		fastpathTV.EncMapStringIntfV(*v, fastpathCheckNilTrue, e)

	case map[string]string:
		fastpathTV.EncMapStringStringV(v, fastpathCheckNilTrue, e)
	case *map[string]string:
		fastpathTV.EncMapStringStringV(*v, fastpathCheckNilTrue, e)

	case map[string]uint:
		fastpathTV.EncMapStringUintV(v, fastpathCheckNilTrue, e)
	case *map[string]uint:
		fastpathTV.EncMapStringUintV(*v, fastpathCheckNilTrue, e)

	case map[string]uint8:
		fastpathTV.EncMapStringUint8V(v, fastpathCheckNilTrue, e)
	case *map[string]uint8:
		fastpathTV.EncMapStringUint8V(*v, fastpathCheckNilTrue, e)

	case map[string]uint16:
		fastpathTV.EncMapStringUint16V(v, fastpathCheckNilTrue, e)
	case *map[string]uint16:
		fastpathTV.EncMapStringUint16V(*v, fastpathCheckNilTrue, e)

	case map[string]uint32:
		fastpathTV.EncMapStringUint32V(v, fastpathCheckNilTrue, e)
	case *map[string]uint32:
		fastpathTV.EncMapStringUint32V(*v, fastpathCheckNilTrue, e)

	case map[string]uint64:
		fastpathTV.EncMapStringUint64V(v, fastpathCheckNilTrue, e)
	case *map[string]uint64:
		fastpathTV.EncMapStringUint64V(*v, fastpathCheckNilTrue, e)

	case map[string]uintptr:
		fastpathTV.EncMapStringUintptrV(v, fastpathCheckNilTrue, e)
	case *map[string]uintptr:
		fastpathTV.EncMapStringUintptrV(*v, fastpathCheckNilTrue, e)

	case map[string]int:
		fastpathTV.EncMapStringIntV(v, fastpathCheckNilTrue, e)
	case *map[string]int:
		fastpathTV.EncMapStringIntV(*v, fastpathCheckNilTrue, e)

	case map[string]int8:
		fastpathTV.EncMapStringInt8V(v, fastpathCheckNilTrue, e)
	case *map[string]int8:
		fastpathTV.EncMapStringInt8V(*v, fastpathCheckNilTrue, e)

	case map[string]int16:
		fastpathTV.EncMapStringInt16V(v, fastpathCheckNilTrue, e)
	case *map[string]int16:
		fastpathTV.EncMapStringInt16V(*v, fastpathCheckNilTrue, e)

	case map[string]int32:
		fastpathTV.EncMapStringInt32V(v, fastpathCheckNilTrue, e)
	case *map[string]int32:
		fastpathTV.EncMapStringInt32V(*v, fastpathCheckNilTrue, e)

	case map[string]int64:
		fastpathTV.EncMapStringInt64V(v, fastpathCheckNilTrue, e)
	case *map[string]int64:
		fastpathTV.EncMapStringInt64V(*v, fastpathCheckNilTrue, e)

	case map[string]float32:
		fastpathTV.EncMapStringFloat32V(v, fastpathCheckNilTrue, e)
	case *map[string]float32:
		fastpathTV.EncMapStringFloat32V(*v, fastpathCheckNilTrue, e)

	case map[string]float64:
		fastpathTV.EncMapStringFloat64V(v, fastpathCheckNilTrue, e)
	case *map[string]float64:
		fastpathTV.EncMapStringFloat64V(*v, fastpathCheckNilTrue, e)

	case map[string]bool:
		fastpathTV.EncMapStringBoolV(v, fastpathCheckNilTrue, e)
	case *map[string]bool:
		fastpathTV.EncMapStringBoolV(*v, fastpathCheckNilTrue, e)

	case map[float32]interface{}:
		fastpathTV.EncMapFloat32IntfV(v, fastpathCheckNilTrue, e)
	case *map[float32]interface{}:
		fastpathTV.EncMapFloat32IntfV(*v, fastpathCheckNilTrue, e)

	case map[float32]string:
		fastpathTV.EncMapFloat32StringV(v, fastpathCheckNilTrue, e)
	case *map[float32]string:
		fastpathTV.EncMapFloat32StringV(*v, fastpathCheckNilTrue, e)

	case map[float32]uint:
		fastpathTV.EncMapFloat32UintV(v, fastpathCheckNilTrue, e)
	case *map[float32]uint:
		fastpathTV.EncMapFloat32UintV(*v, fastpathCheckNilTrue, e)

	case map[float32]uint8:
		fastpathTV.EncMapFloat32Uint8V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint8:
		fastpathTV.EncMapFloat32Uint8V(*v, fastpathCheckNilTrue, e)

	case map[float32]uint16:
		fastpathTV.EncMapFloat32Uint16V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint16:
		fastpathTV.EncMapFloat32Uint16V(*v, fastpathCheckNilTrue, e)

	case map[float32]uint32:
		fastpathTV.EncMapFloat32Uint32V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint32:
		fastpathTV.EncMapFloat32Uint32V(*v, fastpathCheckNilTrue, e)

	case map[float32]uint64:
		fastpathTV.EncMapFloat32Uint64V(v, fastpathCheckNilTrue, e)
	case *map[float32]uint64:
		fastpathTV.EncMapFloat32Uint64V(*v, fastpathCheckNilTrue, e)

	case map[float32]uintptr:
		fastpathTV.EncMapFloat32UintptrV(v, fastpathCheckNilTrue, e)
	case *map[float32]uintptr:
		fastpathTV.EncMapFloat32UintptrV(*v, fastpathCheckNilTrue, e)

	case map[float32]int:
		fastpathTV.EncMapFloat32IntV(v, fastpathCheckNilTrue, e)
	case *map[float32]int:
		fastpathTV.EncMapFloat32IntV(*v, fastpathCheckNilTrue, e)

	case map[float32]int8:
		fastpathTV.EncMapFloat32Int8V(v, fastpathCheckNilTrue, e)
	case *map[float32]int8:
		fastpathTV.EncMapFloat32Int8V(*v, fastpathCheckNilTrue, e)

	case map[float32]int16:
		fastpathTV.EncMapFloat32Int16V(v, fastpathCheckNilTrue, e)
	case *map[float32]int16:
		fastpathTV.EncMapFloat32Int16V(*v, fastpathCheckNilTrue, e)

	case map[float32]int32:
		fastpathTV.EncMapFloat32Int32V(v, fastpathCheckNilTrue, e)
	case *map[float32]int32:
		fastpathTV.EncMapFloat32Int32V(*v, fastpathCheckNilTrue, e)

	case map[float32]int64:
		fastpathTV.EncMapFloat32Int64V(v, fastpathCheckNilTrue, e)
	case *map[float32]int64:
		fastpathTV.EncMapFloat32Int64V(*v, fastpathCheckNilTrue, e)

	case map[float32]float32:
		fastpathTV.EncMapFloat32Float32V(v, fastpathCheckNilTrue, e)
	case *map[float32]float32:
		fastpathTV.EncMapFloat32Float32V(*v, fastpathCheckNilTrue, e)

	case map[float32]float64:
		fastpathTV.EncMapFloat32Float64V(v, fastpathCheckNilTrue, e)
	case *map[float32]float64:
		fastpathTV.EncMapFloat32Float64V(*v, fastpathCheckNilTrue, e)

	case map[float32]bool:
		fastpathTV.EncMapFloat32BoolV(v, fastpathCheckNilTrue, e)
	case *map[float32]bool:
		fastpathTV.EncMapFloat32BoolV(*v, fastpathCheckNilTrue, e)

	case map[float64]interface{}:
		fastpathTV.EncMapFloat64IntfV(v, fastpathCheckNilTrue, e)
	case *map[float64]interface{}:
		fastpathTV.EncMapFloat64IntfV(*v, fastpathCheckNilTrue, e)

	case map[float64]string:
		fastpathTV.EncMapFloat64StringV(v, fastpathCheckNilTrue, e)
	case *map[float64]string:
		fastpathTV.EncMapFloat64StringV(*v, fastpathCheckNilTrue, e)

	case map[float64]uint:
		fastpathTV.EncMapFloat64UintV(v, fastpathCheckNilTrue, e)
	case *map[float64]uint:
		fastpathTV.EncMapFloat64UintV(*v, fastpathCheckNilTrue, e)

	case map[float64]uint8:
		fastpathTV.EncMapFloat64Uint8V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint8:
		fastpathTV.EncMapFloat64Uint8V(*v, fastpathCheckNilTrue, e)

	case map[float64]uint16:
		fastpathTV.EncMapFloat64Uint16V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint16:
		fastpathTV.EncMapFloat64Uint16V(*v, fastpathCheckNilTrue, e)

	case map[float64]uint32:
		fastpathTV.EncMapFloat64Uint32V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint32:
		fastpathTV.EncMapFloat64Uint32V(*v, fastpathCheckNilTrue, e)

	case map[float64]uint64:
		fastpathTV.EncMapFloat64Uint64V(v, fastpathCheckNilTrue, e)
	case *map[float64]uint64:
		fastpathTV.EncMapFloat64Uint64V(*v, fastpathCheckNilTrue, e)

	case map[float64]uintptr:
		fastpathTV.EncMapFloat64UintptrV(v, fastpathCheckNilTrue, e)
	case *map[float64]uintptr:
		fastpathTV.EncMapFloat64UintptrV(*v, fastpathCheckNilTrue, e)

	case map[float64]int:
		fastpathTV.EncMapFloat64IntV(v, fastpathCheckNilTrue, e)
	case *map[float64]int:
		fastpathTV.EncMapFloat64IntV(*v, fastpathCheckNilTrue, e)

	case map[float64]int8:
		fastpathTV.EncMapFloat64Int8V(v, fastpathCheckNilTrue, e)
	case *map[float64]int8:
		fastpathTV.EncMapFloat64Int8V(*v, fastpathCheckNilTrue, e)

	case map[float64]int16:
		fastpathTV.EncMapFloat64Int16V(v, fastpathCheckNilTrue, e)
	case *map[float64]int16:
		fastpathTV.EncMapFloat64Int16V(*v, fastpathCheckNilTrue, e)

	case map[float64]int32:
		fastpathTV.EncMapFloat64Int32V(v, fastpathCheckNilTrue, e)
	case *map[float64]int32:
		fastpathTV.EncMapFloat64Int32V(*v, fastpathCheckNilTrue, e)

	case map[float64]int64:
		fastpathTV.EncMapFloat64Int64V(v, fastpathCheckNilTrue, e)
	case *map[float64]int64:
		fastpathTV.EncMapFloat64Int64V(*v, fastpathCheckNilTrue, e)

	case map[float64]float32:
		fastpathTV.EncMapFloat64Float32V(v, fastpathCheckNilTrue, e)
	case *map[float64]float32:
		fastpathTV.EncMapFloat64Float32V(*v, fastpathCheckNilTrue, e)

	case map[float64]float64:
		fastpathTV.EncMapFloat64Float64V(v, fastpathCheckNilTrue, e)
	case *map[float64]float64:
		fastpathTV.EncMapFloat64Float64V(*v, fastpathCheckNilTrue, e)

	case map[float64]bool:
		fastpathTV.EncMapFloat64BoolV(v, fastpathCheckNilTrue, e)
	case *map[float64]bool:
		fastpathTV.EncMapFloat64BoolV(*v, fastpathCheckNilTrue, e)

	case map[uint]interface{}:
		fastpathTV.EncMapUintIntfV(v, fastpathCheckNilTrue, e)
	case *map[uint]interface{}:
		fastpathTV.EncMapUintIntfV(*v, fastpathCheckNilTrue, e)

	case map[uint]string:
		fastpathTV.EncMapUintStringV(v, fastpathCheckNilTrue, e)
	case *map[uint]string:
		fastpathTV.EncMapUintStringV(*v, fastpathCheckNilTrue, e)

	case map[uint]uint:
		fastpathTV.EncMapUintUintV(v, fastpathCheckNilTrue, e)
	case *map[uint]uint:
		fastpathTV.EncMapUintUintV(*v, fastpathCheckNilTrue, e)

	case map[uint]uint8:
		fastpathTV.EncMapUintUint8V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint8:
		fastpathTV.EncMapUintUint8V(*v, fastpathCheckNilTrue, e)

	case map[uint]uint16:
		fastpathTV.EncMapUintUint16V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint16:
		fastpathTV.EncMapUintUint16V(*v, fastpathCheckNilTrue, e)

	case map[uint]uint32:
		fastpathTV.EncMapUintUint32V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint32:
		fastpathTV.EncMapUintUint32V(*v, fastpathCheckNilTrue, e)

	case map[uint]uint64:
		fastpathTV.EncMapUintUint64V(v, fastpathCheckNilTrue, e)
	case *map[uint]uint64:
		fastpathTV.EncMapUintUint64V(*v, fastpathCheckNilTrue, e)

	case map[uint]uintptr:
		fastpathTV.EncMapUintUintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint]uintptr:
		fastpathTV.EncMapUintUintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint]int:
		fastpathTV.EncMapUintIntV(v, fastpathCheckNilTrue, e)
	case *map[uint]int:
		fastpathTV.EncMapUintIntV(*v, fastpathCheckNilTrue, e)

	case map[uint]int8:
		fastpathTV.EncMapUintInt8V(v, fastpathCheckNilTrue, e)
	case *map[uint]int8:
		fastpathTV.EncMapUintInt8V(*v, fastpathCheckNilTrue, e)

	case map[uint]int16:
		fastpathTV.EncMapUintInt16V(v, fastpathCheckNilTrue, e)
	case *map[uint]int16:
		fastpathTV.EncMapUintInt16V(*v, fastpathCheckNilTrue, e)

	case map[uint]int32:
		fastpathTV.EncMapUintInt32V(v, fastpathCheckNilTrue, e)
	case *map[uint]int32:
		fastpathTV.EncMapUintInt32V(*v, fastpathCheckNilTrue, e)

	case map[uint]int64:
		fastpathTV.EncMapUintInt64V(v, fastpathCheckNilTrue, e)
	case *map[uint]int64:
		fastpathTV.EncMapUintInt64V(*v, fastpathCheckNilTrue, e)

	case map[uint]float32:
		fastpathTV.EncMapUintFloat32V(v, fastpathCheckNilTrue, e)
	case *map[uint]float32:
		fastpathTV.EncMapUintFloat32V(*v, fastpathCheckNilTrue, e)

	case map[uint]float64:
		fastpathTV.EncMapUintFloat64V(v, fastpathCheckNilTrue, e)
	case *map[uint]float64:
		fastpathTV.EncMapUintFloat64V(*v, fastpathCheckNilTrue, e)

	case map[uint]bool:
		fastpathTV.EncMapUintBoolV(v, fastpathCheckNilTrue, e)
	case *map[uint]bool:
		fastpathTV.EncMapUintBoolV(*v, fastpathCheckNilTrue, e)

	case map[uint8]interface{}:
		fastpathTV.EncMapUint8IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint8]interface{}:
		fastpathTV.EncMapUint8IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint8]string:
		fastpathTV.EncMapUint8StringV(v, fastpathCheckNilTrue, e)
	case *map[uint8]string:
		fastpathTV.EncMapUint8StringV(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint:
		fastpathTV.EncMapUint8UintV(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint:
		fastpathTV.EncMapUint8UintV(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint8:
		fastpathTV.EncMapUint8Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint8:
		fastpathTV.EncMapUint8Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint16:
		fastpathTV.EncMapUint8Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint16:
		fastpathTV.EncMapUint8Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint32:
		fastpathTV.EncMapUint8Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint32:
		fastpathTV.EncMapUint8Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uint64:
		fastpathTV.EncMapUint8Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint8]uint64:
		fastpathTV.EncMapUint8Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint8]uintptr:
		fastpathTV.EncMapUint8UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint8]uintptr:
		fastpathTV.EncMapUint8UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint8]int:
		fastpathTV.EncMapUint8IntV(v, fastpathCheckNilTrue, e)
	case *map[uint8]int:
		fastpathTV.EncMapUint8IntV(*v, fastpathCheckNilTrue, e)

	case map[uint8]int8:
		fastpathTV.EncMapUint8Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int8:
		fastpathTV.EncMapUint8Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint8]int16:
		fastpathTV.EncMapUint8Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int16:
		fastpathTV.EncMapUint8Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint8]int32:
		fastpathTV.EncMapUint8Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int32:
		fastpathTV.EncMapUint8Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint8]int64:
		fastpathTV.EncMapUint8Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint8]int64:
		fastpathTV.EncMapUint8Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint8]float32:
		fastpathTV.EncMapUint8Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint8]float32:
		fastpathTV.EncMapUint8Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint8]float64:
		fastpathTV.EncMapUint8Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint8]float64:
		fastpathTV.EncMapUint8Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint8]bool:
		fastpathTV.EncMapUint8BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint8]bool:
		fastpathTV.EncMapUint8BoolV(*v, fastpathCheckNilTrue, e)

	case map[uint16]interface{}:
		fastpathTV.EncMapUint16IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint16]interface{}:
		fastpathTV.EncMapUint16IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint16]string:
		fastpathTV.EncMapUint16StringV(v, fastpathCheckNilTrue, e)
	case *map[uint16]string:
		fastpathTV.EncMapUint16StringV(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint:
		fastpathTV.EncMapUint16UintV(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint:
		fastpathTV.EncMapUint16UintV(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint8:
		fastpathTV.EncMapUint16Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint8:
		fastpathTV.EncMapUint16Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint16:
		fastpathTV.EncMapUint16Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint16:
		fastpathTV.EncMapUint16Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint32:
		fastpathTV.EncMapUint16Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint32:
		fastpathTV.EncMapUint16Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uint64:
		fastpathTV.EncMapUint16Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint16]uint64:
		fastpathTV.EncMapUint16Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint16]uintptr:
		fastpathTV.EncMapUint16UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint16]uintptr:
		fastpathTV.EncMapUint16UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint16]int:
		fastpathTV.EncMapUint16IntV(v, fastpathCheckNilTrue, e)
	case *map[uint16]int:
		fastpathTV.EncMapUint16IntV(*v, fastpathCheckNilTrue, e)

	case map[uint16]int8:
		fastpathTV.EncMapUint16Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int8:
		fastpathTV.EncMapUint16Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint16]int16:
		fastpathTV.EncMapUint16Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int16:
		fastpathTV.EncMapUint16Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint16]int32:
		fastpathTV.EncMapUint16Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int32:
		fastpathTV.EncMapUint16Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint16]int64:
		fastpathTV.EncMapUint16Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint16]int64:
		fastpathTV.EncMapUint16Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint16]float32:
		fastpathTV.EncMapUint16Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint16]float32:
		fastpathTV.EncMapUint16Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint16]float64:
		fastpathTV.EncMapUint16Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint16]float64:
		fastpathTV.EncMapUint16Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint16]bool:
		fastpathTV.EncMapUint16BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint16]bool:
		fastpathTV.EncMapUint16BoolV(*v, fastpathCheckNilTrue, e)

	case map[uint32]interface{}:
		fastpathTV.EncMapUint32IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint32]interface{}:
		fastpathTV.EncMapUint32IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint32]string:
		fastpathTV.EncMapUint32StringV(v, fastpathCheckNilTrue, e)
	case *map[uint32]string:
		fastpathTV.EncMapUint32StringV(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint:
		fastpathTV.EncMapUint32UintV(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint:
		fastpathTV.EncMapUint32UintV(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint8:
		fastpathTV.EncMapUint32Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint8:
		fastpathTV.EncMapUint32Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint16:
		fastpathTV.EncMapUint32Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint16:
		fastpathTV.EncMapUint32Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint32:
		fastpathTV.EncMapUint32Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint32:
		fastpathTV.EncMapUint32Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uint64:
		fastpathTV.EncMapUint32Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint32]uint64:
		fastpathTV.EncMapUint32Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint32]uintptr:
		fastpathTV.EncMapUint32UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint32]uintptr:
		fastpathTV.EncMapUint32UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint32]int:
		fastpathTV.EncMapUint32IntV(v, fastpathCheckNilTrue, e)
	case *map[uint32]int:
		fastpathTV.EncMapUint32IntV(*v, fastpathCheckNilTrue, e)

	case map[uint32]int8:
		fastpathTV.EncMapUint32Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int8:
		fastpathTV.EncMapUint32Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint32]int16:
		fastpathTV.EncMapUint32Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int16:
		fastpathTV.EncMapUint32Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint32]int32:
		fastpathTV.EncMapUint32Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int32:
		fastpathTV.EncMapUint32Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]int64:
		fastpathTV.EncMapUint32Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint32]int64:
		fastpathTV.EncMapUint32Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint32]float32:
		fastpathTV.EncMapUint32Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint32]float32:
		fastpathTV.EncMapUint32Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint32]float64:
		fastpathTV.EncMapUint32Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint32]float64:
		fastpathTV.EncMapUint32Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint32]bool:
		fastpathTV.EncMapUint32BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint32]bool:
		fastpathTV.EncMapUint32BoolV(*v, fastpathCheckNilTrue, e)

	case map[uint64]interface{}:
		fastpathTV.EncMapUint64IntfV(v, fastpathCheckNilTrue, e)
	case *map[uint64]interface{}:
		fastpathTV.EncMapUint64IntfV(*v, fastpathCheckNilTrue, e)

	case map[uint64]string:
		fastpathTV.EncMapUint64StringV(v, fastpathCheckNilTrue, e)
	case *map[uint64]string:
		fastpathTV.EncMapUint64StringV(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint:
		fastpathTV.EncMapUint64UintV(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint:
		fastpathTV.EncMapUint64UintV(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint8:
		fastpathTV.EncMapUint64Uint8V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint8:
		fastpathTV.EncMapUint64Uint8V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint16:
		fastpathTV.EncMapUint64Uint16V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint16:
		fastpathTV.EncMapUint64Uint16V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint32:
		fastpathTV.EncMapUint64Uint32V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint32:
		fastpathTV.EncMapUint64Uint32V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uint64:
		fastpathTV.EncMapUint64Uint64V(v, fastpathCheckNilTrue, e)
	case *map[uint64]uint64:
		fastpathTV.EncMapUint64Uint64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]uintptr:
		fastpathTV.EncMapUint64UintptrV(v, fastpathCheckNilTrue, e)
	case *map[uint64]uintptr:
		fastpathTV.EncMapUint64UintptrV(*v, fastpathCheckNilTrue, e)

	case map[uint64]int:
		fastpathTV.EncMapUint64IntV(v, fastpathCheckNilTrue, e)
	case *map[uint64]int:
		fastpathTV.EncMapUint64IntV(*v, fastpathCheckNilTrue, e)

	case map[uint64]int8:
		fastpathTV.EncMapUint64Int8V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int8:
		fastpathTV.EncMapUint64Int8V(*v, fastpathCheckNilTrue, e)

	case map[uint64]int16:
		fastpathTV.EncMapUint64Int16V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int16:
		fastpathTV.EncMapUint64Int16V(*v, fastpathCheckNilTrue, e)

	case map[uint64]int32:
		fastpathTV.EncMapUint64Int32V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int32:
		fastpathTV.EncMapUint64Int32V(*v, fastpathCheckNilTrue, e)

	case map[uint64]int64:
		fastpathTV.EncMapUint64Int64V(v, fastpathCheckNilTrue, e)
	case *map[uint64]int64:
		fastpathTV.EncMapUint64Int64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]float32:
		fastpathTV.EncMapUint64Float32V(v, fastpathCheckNilTrue, e)
	case *map[uint64]float32:
		fastpathTV.EncMapUint64Float32V(*v, fastpathCheckNilTrue, e)

	case map[uint64]float64:
		fastpathTV.EncMapUint64Float64V(v, fastpathCheckNilTrue, e)
	case *map[uint64]float64:
		fastpathTV.EncMapUint64Float64V(*v, fastpathCheckNilTrue, e)

	case map[uint64]bool:
		fastpathTV.EncMapUint64BoolV(v, fastpathCheckNilTrue, e)
	case *map[uint64]bool:
		fastpathTV.EncMapUint64BoolV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]interface{}:
		fastpathTV.EncMapUintptrIntfV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]interface{}:
		fastpathTV.EncMapUintptrIntfV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]string:
		fastpathTV.EncMapUintptrStringV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]string:
		fastpathTV.EncMapUintptrStringV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint:
		fastpathTV.EncMapUintptrUintV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint:
		fastpathTV.EncMapUintptrUintV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint8:
		fastpathTV.EncMapUintptrUint8V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint8:
		fastpathTV.EncMapUintptrUint8V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint16:
		fastpathTV.EncMapUintptrUint16V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint16:
		fastpathTV.EncMapUintptrUint16V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint32:
		fastpathTV.EncMapUintptrUint32V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint32:
		fastpathTV.EncMapUintptrUint32V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uint64:
		fastpathTV.EncMapUintptrUint64V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uint64:
		fastpathTV.EncMapUintptrUint64V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]uintptr:
		fastpathTV.EncMapUintptrUintptrV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]uintptr:
		fastpathTV.EncMapUintptrUintptrV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int:
		fastpathTV.EncMapUintptrIntV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int:
		fastpathTV.EncMapUintptrIntV(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int8:
		fastpathTV.EncMapUintptrInt8V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int8:
		fastpathTV.EncMapUintptrInt8V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int16:
		fastpathTV.EncMapUintptrInt16V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int16:
		fastpathTV.EncMapUintptrInt16V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int32:
		fastpathTV.EncMapUintptrInt32V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int32:
		fastpathTV.EncMapUintptrInt32V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]int64:
		fastpathTV.EncMapUintptrInt64V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]int64:
		fastpathTV.EncMapUintptrInt64V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]float32:
		fastpathTV.EncMapUintptrFloat32V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]float32:
		fastpathTV.EncMapUintptrFloat32V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]float64:
		fastpathTV.EncMapUintptrFloat64V(v, fastpathCheckNilTrue, e)
	case *map[uintptr]float64:
		fastpathTV.EncMapUintptrFloat64V(*v, fastpathCheckNilTrue, e)

	case map[uintptr]bool:
		fastpathTV.EncMapUintptrBoolV(v, fastpathCheckNilTrue, e)
	case *map[uintptr]bool:
		fastpathTV.EncMapUintptrBoolV(*v, fastpathCheckNilTrue, e)

	case map[int]interface{}:
		fastpathTV.EncMapIntIntfV(v, fastpathCheckNilTrue, e)
	case *map[int]interface{}:
		fastpathTV.EncMapIntIntfV(*v, fastpathCheckNilTrue, e)

	case map[int]string:
		fastpathTV.EncMapIntStringV(v, fastpathCheckNilTrue, e)
	case *map[int]string:
		fastpathTV.EncMapIntStringV(*v, fastpathCheckNilTrue, e)

	case map[int]uint:
		fastpathTV.EncMapIntUintV(v, fastpathCheckNilTrue, e)
	case *map[int]uint:
		fastpathTV.EncMapIntUintV(*v, fastpathCheckNilTrue, e)

	case map[int]uint8:
		fastpathTV.EncMapIntUint8V(v, fastpathCheckNilTrue, e)
	case *map[int]uint8:
		fastpathTV.EncMapIntUint8V(*v, fastpathCheckNilTrue, e)

	case map[int]uint16:
		fastpathTV.EncMapIntUint16V(v, fastpathCheckNilTrue, e)
	case *map[int]uint16:
		fastpathTV.EncMapIntUint16V(*v, fastpathCheckNilTrue, e)

	case map[int]uint32:
		fastpathTV.EncMapIntUint32V(v, fastpathCheckNilTrue, e)
	case *map[int]uint32:
		fastpathTV.EncMapIntUint32V(*v, fastpathCheckNilTrue, e)

	case map[int]uint64:
		fastpathTV.EncMapIntUint64V(v, fastpathCheckNilTrue, e)
	case *map[int]uint64:
		fastpathTV.EncMapIntUint64V(*v, fastpathCheckNilTrue, e)

	case map[int]uintptr:
		fastpathTV.EncMapIntUintptrV(v, fastpathCheckNilTrue, e)
	case *map[int]uintptr:
		fastpathTV.EncMapIntUintptrV(*v, fastpathCheckNilTrue, e)

	case map[int]int:
		fastpathTV.EncMapIntIntV(v, fastpathCheckNilTrue, e)
	case *map[int]int:
		fastpathTV.EncMapIntIntV(*v, fastpathCheckNilTrue, e)

	case map[int]int8:
		fastpathTV.EncMapIntInt8V(v, fastpathCheckNilTrue, e)
	case *map[int]int8:
		fastpathTV.EncMapIntInt8V(*v, fastpathCheckNilTrue, e)

	case map[int]int16:
		fastpathTV.EncMapIntInt16V(v, fastpathCheckNilTrue, e)
	case *map[int]int16:
		fastpathTV.EncMapIntInt16V(*v, fastpathCheckNilTrue, e)

	case map[int]int32:
		fastpathTV.EncMapIntInt32V(v, fastpathCheckNilTrue, e)
	case *map[int]int32:
		fastpathTV.EncMapIntInt32V(*v, fastpathCheckNilTrue, e)

	case map[int]int64:
		fastpathTV.EncMapIntInt64V(v, fastpathCheckNilTrue, e)
	case *map[int]int64:
		fastpathTV.EncMapIntInt64V(*v, fastpathCheckNilTrue, e)

	case map[int]float32:
		fastpathTV.EncMapIntFloat32V(v, fastpathCheckNilTrue, e)
	case *map[int]float32:
		fastpathTV.EncMapIntFloat32V(*v, fastpathCheckNilTrue, e)

	case map[int]float64:
		fastpathTV.EncMapIntFloat64V(v, fastpathCheckNilTrue, e)
	case *map[int]float64:
		fastpathTV.EncMapIntFloat64V(*v, fastpathCheckNilTrue, e)

	case map[int]bool:
		fastpathTV.EncMapIntBoolV(v, fastpathCheckNilTrue, e)
	case *map[int]bool:
		fastpathTV.EncMapIntBoolV(*v, fastpathCheckNilTrue, e)

	case map[int8]interface{}:
		fastpathTV.EncMapInt8IntfV(v, fastpathCheckNilTrue, e)
	case *map[int8]interface{}:
		fastpathTV.EncMapInt8IntfV(*v, fastpathCheckNilTrue, e)

	case map[int8]string:
		fastpathTV.EncMapInt8StringV(v, fastpathCheckNilTrue, e)
	case *map[int8]string:
		fastpathTV.EncMapInt8StringV(*v, fastpathCheckNilTrue, e)

	case map[int8]uint:
		fastpathTV.EncMapInt8UintV(v, fastpathCheckNilTrue, e)
	case *map[int8]uint:
		fastpathTV.EncMapInt8UintV(*v, fastpathCheckNilTrue, e)

	case map[int8]uint8:
		fastpathTV.EncMapInt8Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint8:
		fastpathTV.EncMapInt8Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int8]uint16:
		fastpathTV.EncMapInt8Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint16:
		fastpathTV.EncMapInt8Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int8]uint32:
		fastpathTV.EncMapInt8Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint32:
		fastpathTV.EncMapInt8Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int8]uint64:
		fastpathTV.EncMapInt8Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int8]uint64:
		fastpathTV.EncMapInt8Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int8]uintptr:
		fastpathTV.EncMapInt8UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int8]uintptr:
		fastpathTV.EncMapInt8UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int8]int:
		fastpathTV.EncMapInt8IntV(v, fastpathCheckNilTrue, e)
	case *map[int8]int:
		fastpathTV.EncMapInt8IntV(*v, fastpathCheckNilTrue, e)

	case map[int8]int8:
		fastpathTV.EncMapInt8Int8V(v, fastpathCheckNilTrue, e)
	case *map[int8]int8:
		fastpathTV.EncMapInt8Int8V(*v, fastpathCheckNilTrue, e)

	case map[int8]int16:
		fastpathTV.EncMapInt8Int16V(v, fastpathCheckNilTrue, e)
	case *map[int8]int16:
		fastpathTV.EncMapInt8Int16V(*v, fastpathCheckNilTrue, e)

	case map[int8]int32:
		fastpathTV.EncMapInt8Int32V(v, fastpathCheckNilTrue, e)
	case *map[int8]int32:
		fastpathTV.EncMapInt8Int32V(*v, fastpathCheckNilTrue, e)

	case map[int8]int64:
		fastpathTV.EncMapInt8Int64V(v, fastpathCheckNilTrue, e)
	case *map[int8]int64:
		fastpathTV.EncMapInt8Int64V(*v, fastpathCheckNilTrue, e)

	case map[int8]float32:
		fastpathTV.EncMapInt8Float32V(v, fastpathCheckNilTrue, e)
	case *map[int8]float32:
		fastpathTV.EncMapInt8Float32V(*v, fastpathCheckNilTrue, e)

	case map[int8]float64:
		fastpathTV.EncMapInt8Float64V(v, fastpathCheckNilTrue, e)
	case *map[int8]float64:
		fastpathTV.EncMapInt8Float64V(*v, fastpathCheckNilTrue, e)

	case map[int8]bool:
		fastpathTV.EncMapInt8BoolV(v, fastpathCheckNilTrue, e)
	case *map[int8]bool:
		fastpathTV.EncMapInt8BoolV(*v, fastpathCheckNilTrue, e)

	case map[int16]interface{}:
		fastpathTV.EncMapInt16IntfV(v, fastpathCheckNilTrue, e)
	case *map[int16]interface{}:
		fastpathTV.EncMapInt16IntfV(*v, fastpathCheckNilTrue, e)

	case map[int16]string:
		fastpathTV.EncMapInt16StringV(v, fastpathCheckNilTrue, e)
	case *map[int16]string:
		fastpathTV.EncMapInt16StringV(*v, fastpathCheckNilTrue, e)

	case map[int16]uint:
		fastpathTV.EncMapInt16UintV(v, fastpathCheckNilTrue, e)
	case *map[int16]uint:
		fastpathTV.EncMapInt16UintV(*v, fastpathCheckNilTrue, e)

	case map[int16]uint8:
		fastpathTV.EncMapInt16Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint8:
		fastpathTV.EncMapInt16Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int16]uint16:
		fastpathTV.EncMapInt16Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint16:
		fastpathTV.EncMapInt16Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int16]uint32:
		fastpathTV.EncMapInt16Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint32:
		fastpathTV.EncMapInt16Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int16]uint64:
		fastpathTV.EncMapInt16Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int16]uint64:
		fastpathTV.EncMapInt16Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int16]uintptr:
		fastpathTV.EncMapInt16UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int16]uintptr:
		fastpathTV.EncMapInt16UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int16]int:
		fastpathTV.EncMapInt16IntV(v, fastpathCheckNilTrue, e)
	case *map[int16]int:
		fastpathTV.EncMapInt16IntV(*v, fastpathCheckNilTrue, e)

	case map[int16]int8:
		fastpathTV.EncMapInt16Int8V(v, fastpathCheckNilTrue, e)
	case *map[int16]int8:
		fastpathTV.EncMapInt16Int8V(*v, fastpathCheckNilTrue, e)

	case map[int16]int16:
		fastpathTV.EncMapInt16Int16V(v, fastpathCheckNilTrue, e)
	case *map[int16]int16:
		fastpathTV.EncMapInt16Int16V(*v, fastpathCheckNilTrue, e)

	case map[int16]int32:
		fastpathTV.EncMapInt16Int32V(v, fastpathCheckNilTrue, e)
	case *map[int16]int32:
		fastpathTV.EncMapInt16Int32V(*v, fastpathCheckNilTrue, e)

	case map[int16]int64:
		fastpathTV.EncMapInt16Int64V(v, fastpathCheckNilTrue, e)
	case *map[int16]int64:
		fastpathTV.EncMapInt16Int64V(*v, fastpathCheckNilTrue, e)

	case map[int16]float32:
		fastpathTV.EncMapInt16Float32V(v, fastpathCheckNilTrue, e)
	case *map[int16]float32:
		fastpathTV.EncMapInt16Float32V(*v, fastpathCheckNilTrue, e)

	case map[int16]float64:
		fastpathTV.EncMapInt16Float64V(v, fastpathCheckNilTrue, e)
	case *map[int16]float64:
		fastpathTV.EncMapInt16Float64V(*v, fastpathCheckNilTrue, e)

	case map[int16]bool:
		fastpathTV.EncMapInt16BoolV(v, fastpathCheckNilTrue, e)
	case *map[int16]bool:
		fastpathTV.EncMapInt16BoolV(*v, fastpathCheckNilTrue, e)

	case map[int32]interface{}:
		fastpathTV.EncMapInt32IntfV(v, fastpathCheckNilTrue, e)
	case *map[int32]interface{}:
		fastpathTV.EncMapInt32IntfV(*v, fastpathCheckNilTrue, e)

	case map[int32]string:
		fastpathTV.EncMapInt32StringV(v, fastpathCheckNilTrue, e)
	case *map[int32]string:
		fastpathTV.EncMapInt32StringV(*v, fastpathCheckNilTrue, e)

	case map[int32]uint:
		fastpathTV.EncMapInt32UintV(v, fastpathCheckNilTrue, e)
	case *map[int32]uint:
		fastpathTV.EncMapInt32UintV(*v, fastpathCheckNilTrue, e)

	case map[int32]uint8:
		fastpathTV.EncMapInt32Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint8:
		fastpathTV.EncMapInt32Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int32]uint16:
		fastpathTV.EncMapInt32Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint16:
		fastpathTV.EncMapInt32Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int32]uint32:
		fastpathTV.EncMapInt32Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint32:
		fastpathTV.EncMapInt32Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int32]uint64:
		fastpathTV.EncMapInt32Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int32]uint64:
		fastpathTV.EncMapInt32Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int32]uintptr:
		fastpathTV.EncMapInt32UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int32]uintptr:
		fastpathTV.EncMapInt32UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int32]int:
		fastpathTV.EncMapInt32IntV(v, fastpathCheckNilTrue, e)
	case *map[int32]int:
		fastpathTV.EncMapInt32IntV(*v, fastpathCheckNilTrue, e)

	case map[int32]int8:
		fastpathTV.EncMapInt32Int8V(v, fastpathCheckNilTrue, e)
	case *map[int32]int8:
		fastpathTV.EncMapInt32Int8V(*v, fastpathCheckNilTrue, e)

	case map[int32]int16:
		fastpathTV.EncMapInt32Int16V(v, fastpathCheckNilTrue, e)
	case *map[int32]int16:
		fastpathTV.EncMapInt32Int16V(*v, fastpathCheckNilTrue, e)

	case map[int32]int32:
		fastpathTV.EncMapInt32Int32V(v, fastpathCheckNilTrue, e)
	case *map[int32]int32:
		fastpathTV.EncMapInt32Int32V(*v, fastpathCheckNilTrue, e)

	case map[int32]int64:
		fastpathTV.EncMapInt32Int64V(v, fastpathCheckNilTrue, e)
	case *map[int32]int64:
		fastpathTV.EncMapInt32Int64V(*v, fastpathCheckNilTrue, e)

	case map[int32]float32:
		fastpathTV.EncMapInt32Float32V(v, fastpathCheckNilTrue, e)
	case *map[int32]float32:
		fastpathTV.EncMapInt32Float32V(*v, fastpathCheckNilTrue, e)

	case map[int32]float64:
		fastpathTV.EncMapInt32Float64V(v, fastpathCheckNilTrue, e)
	case *map[int32]float64:
		fastpathTV.EncMapInt32Float64V(*v, fastpathCheckNilTrue, e)

	case map[int32]bool:
		fastpathTV.EncMapInt32BoolV(v, fastpathCheckNilTrue, e)
	case *map[int32]bool:
		fastpathTV.EncMapInt32BoolV(*v, fastpathCheckNilTrue, e)

	case map[int64]interface{}:
		fastpathTV.EncMapInt64IntfV(v, fastpathCheckNilTrue, e)
	case *map[int64]interface{}:
		fastpathTV.EncMapInt64IntfV(*v, fastpathCheckNilTrue, e)

	case map[int64]string:
		fastpathTV.EncMapInt64StringV(v, fastpathCheckNilTrue, e)
	case *map[int64]string:
		fastpathTV.EncMapInt64StringV(*v, fastpathCheckNilTrue, e)

	case map[int64]uint:
		fastpathTV.EncMapInt64UintV(v, fastpathCheckNilTrue, e)
	case *map[int64]uint:
		fastpathTV.EncMapInt64UintV(*v, fastpathCheckNilTrue, e)

	case map[int64]uint8:
		fastpathTV.EncMapInt64Uint8V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint8:
		fastpathTV.EncMapInt64Uint8V(*v, fastpathCheckNilTrue, e)

	case map[int64]uint16:
		fastpathTV.EncMapInt64Uint16V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint16:
		fastpathTV.EncMapInt64Uint16V(*v, fastpathCheckNilTrue, e)

	case map[int64]uint32:
		fastpathTV.EncMapInt64Uint32V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint32:
		fastpathTV.EncMapInt64Uint32V(*v, fastpathCheckNilTrue, e)

	case map[int64]uint64:
		fastpathTV.EncMapInt64Uint64V(v, fastpathCheckNilTrue, e)
	case *map[int64]uint64:
		fastpathTV.EncMapInt64Uint64V(*v, fastpathCheckNilTrue, e)

	case map[int64]uintptr:
		fastpathTV.EncMapInt64UintptrV(v, fastpathCheckNilTrue, e)
	case *map[int64]uintptr:
		fastpathTV.EncMapInt64UintptrV(*v, fastpathCheckNilTrue, e)

	case map[int64]int:
		fastpathTV.EncMapInt64IntV(v, fastpathCheckNilTrue, e)
	case *map[int64]int:
		fastpathTV.EncMapInt64IntV(*v, fastpathCheckNilTrue, e)

	case map[int64]int8:
		fastpathTV.EncMapInt64Int8V(v, fastpathCheckNilTrue, e)
	case *map[int64]int8:
		fastpathTV.EncMapInt64Int8V(*v, fastpathCheckNilTrue, e)

	case map[int64]int16:
		fastpathTV.EncMapInt64Int16V(v, fastpathCheckNilTrue, e)
	case *map[int64]int16:
		fastpathTV.EncMapInt64Int16V(*v, fastpathCheckNilTrue, e)

	case map[int64]int32:
		fastpathTV.EncMapInt64Int32V(v, fastpathCheckNilTrue, e)
	case *map[int64]int32:
		fastpathTV.EncMapInt64Int32V(*v, fastpathCheckNilTrue, e)

	case map[int64]int64:
		fastpathTV.EncMapInt64Int64V(v, fastpathCheckNilTrue, e)
	case *map[int64]int64:
		fastpathTV.EncMapInt64Int64V(*v, fastpathCheckNilTrue, e)

	case map[int64]float32:
		fastpathTV.EncMapInt64Float32V(v, fastpathCheckNilTrue, e)
	case *map[int64]float32:
		fastpathTV.EncMapInt64Float32V(*v, fastpathCheckNilTrue, e)

	case map[int64]float64:
		fastpathTV.EncMapInt64Float64V(v, fastpathCheckNilTrue, e)
	case *map[int64]float64:
		fastpathTV.EncMapInt64Float64V(*v, fastpathCheckNilTrue, e)

	case map[int64]bool:
		fastpathTV.EncMapInt64BoolV(v, fastpathCheckNilTrue, e)
	case *map[int64]bool:
		fastpathTV.EncMapInt64BoolV(*v, fastpathCheckNilTrue, e)

	case map[bool]interface{}:
		fastpathTV.EncMapBoolIntfV(v, fastpathCheckNilTrue, e)
	case *map[bool]interface{}:
		fastpathTV.EncMapBoolIntfV(*v, fastpathCheckNilTrue, e)

	case map[bool]string:
		fastpathTV.EncMapBoolStringV(v, fastpathCheckNilTrue, e)
	case *map[bool]string:
		fastpathTV.EncMapBoolStringV(*v, fastpathCheckNilTrue, e)

	case map[bool]uint:
		fastpathTV.EncMapBoolUintV(v, fastpathCheckNilTrue, e)
	case *map[bool]uint:
		fastpathTV.EncMapBoolUintV(*v, fastpathCheckNilTrue, e)

	case map[bool]uint8:
		fastpathTV.EncMapBoolUint8V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint8:
		fastpathTV.EncMapBoolUint8V(*v, fastpathCheckNilTrue, e)

	case map[bool]uint16:
		fastpathTV.EncMapBoolUint16V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint16:
		fastpathTV.EncMapBoolUint16V(*v, fastpathCheckNilTrue, e)

	case map[bool]uint32:
		fastpathTV.EncMapBoolUint32V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint32:
		fastpathTV.EncMapBoolUint32V(*v, fastpathCheckNilTrue, e)

	case map[bool]uint64:
		fastpathTV.EncMapBoolUint64V(v, fastpathCheckNilTrue, e)
	case *map[bool]uint64:
		fastpathTV.EncMapBoolUint64V(*v, fastpathCheckNilTrue, e)

	case map[bool]uintptr:
		fastpathTV.EncMapBoolUintptrV(v, fastpathCheckNilTrue, e)
	case *map[bool]uintptr:
		fastpathTV.EncMapBoolUintptrV(*v, fastpathCheckNilTrue, e)

	case map[bool]int:
		fastpathTV.EncMapBoolIntV(v, fastpathCheckNilTrue, e)
	case *map[bool]int:
		fastpathTV.EncMapBoolIntV(*v, fastpathCheckNilTrue, e)

	case map[bool]int8:
		fastpathTV.EncMapBoolInt8V(v, fastpathCheckNilTrue, e)
	case *map[bool]int8:
		fastpathTV.EncMapBoolInt8V(*v, fastpathCheckNilTrue, e)

	case map[bool]int16:
		fastpathTV.EncMapBoolInt16V(v, fastpathCheckNilTrue, e)
	case *map[bool]int16:
		fastpathTV.EncMapBoolInt16V(*v, fastpathCheckNilTrue, e)

	case map[bool]int32:
		fastpathTV.EncMapBoolInt32V(v, fastpathCheckNilTrue, e)
	case *map[bool]int32:
		fastpathTV.EncMapBoolInt32V(*v, fastpathCheckNilTrue, e)

	case map[bool]int64:
		fastpathTV.EncMapBoolInt64V(v, fastpathCheckNilTrue, e)
	case *map[bool]int64:
		fastpathTV.EncMapBoolInt64V(*v, fastpathCheckNilTrue, e)

	case map[bool]float32:
		fastpathTV.EncMapBoolFloat32V(v, fastpathCheckNilTrue, e)
	case *map[bool]float32:
		fastpathTV.EncMapBoolFloat32V(*v, fastpathCheckNilTrue, e)

	case map[bool]float64:
		fastpathTV.EncMapBoolFloat64V(v, fastpathCheckNilTrue, e)
	case *map[bool]float64:
		fastpathTV.EncMapBoolFloat64V(*v, fastpathCheckNilTrue, e)

	case map[bool]bool:
		fastpathTV.EncMapBoolBoolV(v, fastpathCheckNilTrue, e)
	case *map[bool]bool:
		fastpathTV.EncMapBoolBoolV(*v, fastpathCheckNilTrue, e)

	default:
		_ = v // TODO: workaround https://github.com/golang/go/issues/12927 (remove after go 1.6 release)
		return false
	}
	return true
}

// -- -- fast path functions

func (f *encFnInfo) fastpathEncSliceIntfR(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceIntfV(rv.Interface().([]interface{}), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceIntfV(rv.Interface().([]interface{}), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceIntfV(v []interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		e.encode(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceIntfV(v []interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		e.encode(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceStringR(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceStringV(rv.Interface().([]string), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceStringV(rv.Interface().([]string), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceStringV(v []string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeString(c_UTF8, v2)
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceStringV(v []string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeString(c_UTF8, v2)
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceFloat32R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceFloat32V(rv.Interface().([]float32), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceFloat32V(rv.Interface().([]float32), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceFloat32V(v []float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeFloat32(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceFloat32V(v []float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeFloat32(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceFloat64R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceFloat64V(rv.Interface().([]float64), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceFloat64V(rv.Interface().([]float64), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceFloat64V(v []float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeFloat64(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceFloat64V(v []float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeFloat64(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceUintR(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUintV(rv.Interface().([]uint), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceUintV(rv.Interface().([]uint), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceUintV(v []uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceUintV(v []uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceUint16R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint16V(rv.Interface().([]uint16), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceUint16V(rv.Interface().([]uint16), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceUint16V(v []uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceUint16V(v []uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceUint32R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint32V(rv.Interface().([]uint32), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceUint32V(rv.Interface().([]uint32), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceUint32V(v []uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceUint32V(v []uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceUint64R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUint64V(rv.Interface().([]uint64), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceUint64V(rv.Interface().([]uint64), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceUint64V(v []uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceUint64V(v []uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeUint(uint64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceUintptrR(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceUintptrV(rv.Interface().([]uintptr), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceUintptrV(rv.Interface().([]uintptr), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceUintptrV(v []uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		e.encode(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceUintptrV(v []uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		e.encode(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceIntR(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceIntV(rv.Interface().([]int), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceIntV(rv.Interface().([]int), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceIntV(v []int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceIntV(v []int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceInt8R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt8V(rv.Interface().([]int8), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceInt8V(rv.Interface().([]int8), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceInt8V(v []int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceInt8V(v []int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceInt16R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt16V(rv.Interface().([]int16), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceInt16V(rv.Interface().([]int16), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceInt16V(v []int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceInt16V(v []int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceInt32R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt32V(rv.Interface().([]int32), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceInt32V(rv.Interface().([]int32), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceInt32V(v []int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceInt32V(v []int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceInt64R(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceInt64V(rv.Interface().([]int64), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceInt64V(rv.Interface().([]int64), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceInt64V(v []int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceInt64V(v []int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeInt(int64(v2))
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncSliceBoolR(rv reflect.Value) {
	if f.ti.mbs {
		fastpathTV.EncAsMapSliceBoolV(rv.Interface().([]bool), fastpathCheckNilFalse, f.e)
	} else {
		fastpathTV.EncSliceBoolV(rv.Interface().([]bool), fastpathCheckNilFalse, f.e)
	}
}
func (_ fastpathT) EncSliceBoolV(v []bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeArrayStart(len(v))
	for _, v2 := range v {
		if cr != nil {
			cr.sendContainerState(containerArrayElem)
		}
		ee.EncodeBool(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerArrayEnd)
	}
}

func (_ fastpathT) EncAsMapSliceBoolV(v []bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	if len(v)%2 == 1 {
		e.errorf("mapBySlice requires even slice length, but got %v", len(v))
		return
	}
	ee.EncodeMapStart(len(v) / 2)
	for j, v2 := range v {
		if cr != nil {
			if j%2 == 0 {
				cr.sendContainerState(containerMapKey)
			} else {
				cr.sendContainerState(containerMapValue)
			}
		}
		ee.EncodeBool(v2)
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfIntfR(rv reflect.Value) {
	fastpathTV.EncMapIntfIntfV(rv.Interface().(map[interface{}]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfIntfV(v map[interface{}]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfStringR(rv reflect.Value) {
	fastpathTV.EncMapIntfStringV(rv.Interface().(map[interface{}]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfStringV(v map[interface{}]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfUintR(rv reflect.Value) {
	fastpathTV.EncMapIntfUintV(rv.Interface().(map[interface{}]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfUintV(v map[interface{}]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfUint8R(rv reflect.Value) {
	fastpathTV.EncMapIntfUint8V(rv.Interface().(map[interface{}]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfUint8V(v map[interface{}]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfUint16R(rv reflect.Value) {
	fastpathTV.EncMapIntfUint16V(rv.Interface().(map[interface{}]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfUint16V(v map[interface{}]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfUint32R(rv reflect.Value) {
	fastpathTV.EncMapIntfUint32V(rv.Interface().(map[interface{}]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfUint32V(v map[interface{}]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfUint64R(rv reflect.Value) {
	fastpathTV.EncMapIntfUint64V(rv.Interface().(map[interface{}]uint64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfUint64V(v map[interface{}]uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfUintptrR(rv reflect.Value) {
	fastpathTV.EncMapIntfUintptrV(rv.Interface().(map[interface{}]uintptr), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfUintptrV(v map[interface{}]uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfIntR(rv reflect.Value) {
	fastpathTV.EncMapIntfIntV(rv.Interface().(map[interface{}]int), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfIntV(v map[interface{}]int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfInt8R(rv reflect.Value) {
	fastpathTV.EncMapIntfInt8V(rv.Interface().(map[interface{}]int8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfInt8V(v map[interface{}]int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfInt16R(rv reflect.Value) {
	fastpathTV.EncMapIntfInt16V(rv.Interface().(map[interface{}]int16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfInt16V(v map[interface{}]int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfInt32R(rv reflect.Value) {
	fastpathTV.EncMapIntfInt32V(rv.Interface().(map[interface{}]int32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfInt32V(v map[interface{}]int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfInt64R(rv reflect.Value) {
	fastpathTV.EncMapIntfInt64V(rv.Interface().(map[interface{}]int64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfInt64V(v map[interface{}]int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfFloat32R(rv reflect.Value) {
	fastpathTV.EncMapIntfFloat32V(rv.Interface().(map[interface{}]float32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfFloat32V(v map[interface{}]float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfFloat64R(rv reflect.Value) {
	fastpathTV.EncMapIntfFloat64V(rv.Interface().(map[interface{}]float64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfFloat64V(v map[interface{}]float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapIntfBoolR(rv reflect.Value) {
	fastpathTV.EncMapIntfBoolV(rv.Interface().(map[interface{}]bool), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapIntfBoolV(v map[interface{}]bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		var mksv []byte = make([]byte, 0, len(v)*16) // temporary byte slice for the encoding
		e2 := NewEncoderBytes(&mksv, e.hh)
		v2 := make([]bytesI, len(v))
		var i, l int
		var vp *bytesI
		for k2, _ := range v {
			l = len(mksv)
			e2.MustEncode(k2)
			vp = &v2[i]
			vp.v = mksv[l:]
			vp.i = k2
			i++
		}
		sort.Sort(bytesISlice(v2))
		for j := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.asis(v2[j].v)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[v2[j].i])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			e.encode(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringIntfR(rv reflect.Value) {
	fastpathTV.EncMapStringIntfV(rv.Interface().(map[string]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringIntfV(v map[string]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[string(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringStringR(rv reflect.Value) {
	fastpathTV.EncMapStringStringV(rv.Interface().(map[string]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringStringV(v map[string]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v[string(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringUintR(rv reflect.Value) {
	fastpathTV.EncMapStringUintV(rv.Interface().(map[string]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringUintV(v map[string]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringUint8R(rv reflect.Value) {
	fastpathTV.EncMapStringUint8V(rv.Interface().(map[string]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringUint8V(v map[string]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringUint16R(rv reflect.Value) {
	fastpathTV.EncMapStringUint16V(rv.Interface().(map[string]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringUint16V(v map[string]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringUint32R(rv reflect.Value) {
	fastpathTV.EncMapStringUint32V(rv.Interface().(map[string]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringUint32V(v map[string]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringUint64R(rv reflect.Value) {
	fastpathTV.EncMapStringUint64V(rv.Interface().(map[string]uint64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringUint64V(v map[string]uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringUintptrR(rv reflect.Value) {
	fastpathTV.EncMapStringUintptrV(rv.Interface().(map[string]uintptr), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringUintptrV(v map[string]uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[string(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringIntR(rv reflect.Value) {
	fastpathTV.EncMapStringIntV(rv.Interface().(map[string]int), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringIntV(v map[string]int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringInt8R(rv reflect.Value) {
	fastpathTV.EncMapStringInt8V(rv.Interface().(map[string]int8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringInt8V(v map[string]int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringInt16R(rv reflect.Value) {
	fastpathTV.EncMapStringInt16V(rv.Interface().(map[string]int16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringInt16V(v map[string]int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringInt32R(rv reflect.Value) {
	fastpathTV.EncMapStringInt32V(rv.Interface().(map[string]int32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringInt32V(v map[string]int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringInt64R(rv reflect.Value) {
	fastpathTV.EncMapStringInt64V(rv.Interface().(map[string]int64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringInt64V(v map[string]int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[string(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringFloat32R(rv reflect.Value) {
	fastpathTV.EncMapStringFloat32V(rv.Interface().(map[string]float32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringFloat32V(v map[string]float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v[string(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringFloat64R(rv reflect.Value) {
	fastpathTV.EncMapStringFloat64V(rv.Interface().(map[string]float64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringFloat64V(v map[string]float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v[string(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapStringBoolR(rv reflect.Value) {
	fastpathTV.EncMapStringBoolV(rv.Interface().(map[string]bool), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapStringBoolV(v map[string]bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	asSymbols := e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	if e.h.Canonical {
		v2 := make([]string, len(v))
		var i int
		for k, _ := range v {
			v2[i] = string(k)
			i++
		}
		sort.Sort(stringSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v[string(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			if asSymbols {
				ee.EncodeSymbol(k2)
			} else {
				ee.EncodeString(c_UTF8, k2)
			}
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32IntfR(rv reflect.Value) {
	fastpathTV.EncMapFloat32IntfV(rv.Interface().(map[float32]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32IntfV(v map[float32]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[float32(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32StringR(rv reflect.Value) {
	fastpathTV.EncMapFloat32StringV(rv.Interface().(map[float32]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32StringV(v map[float32]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v[float32(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32UintR(rv reflect.Value) {
	fastpathTV.EncMapFloat32UintV(rv.Interface().(map[float32]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32UintV(v map[float32]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Uint8R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Uint8V(rv.Interface().(map[float32]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Uint8V(v map[float32]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Uint16R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Uint16V(rv.Interface().(map[float32]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Uint16V(v map[float32]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Uint32R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Uint32V(rv.Interface().(map[float32]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Uint32V(v map[float32]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Uint64R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Uint64V(rv.Interface().(map[float32]uint64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Uint64V(v map[float32]uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32UintptrR(rv reflect.Value) {
	fastpathTV.EncMapFloat32UintptrV(rv.Interface().(map[float32]uintptr), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32UintptrV(v map[float32]uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[float32(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32IntR(rv reflect.Value) {
	fastpathTV.EncMapFloat32IntV(rv.Interface().(map[float32]int), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32IntV(v map[float32]int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Int8R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Int8V(rv.Interface().(map[float32]int8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Int8V(v map[float32]int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Int16R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Int16V(rv.Interface().(map[float32]int16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Int16V(v map[float32]int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Int32R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Int32V(rv.Interface().(map[float32]int32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Int32V(v map[float32]int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Int64R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Int64V(rv.Interface().(map[float32]int64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Int64V(v map[float32]int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float32(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Float32R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Float32V(rv.Interface().(map[float32]float32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Float32V(v map[float32]float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v[float32(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32Float64R(rv reflect.Value) {
	fastpathTV.EncMapFloat32Float64V(rv.Interface().(map[float32]float64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32Float64V(v map[float32]float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v[float32(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat32BoolR(rv reflect.Value) {
	fastpathTV.EncMapFloat32BoolV(rv.Interface().(map[float32]bool), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat32BoolV(v map[float32]bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(float32(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v[float32(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat32(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64IntfR(rv reflect.Value) {
	fastpathTV.EncMapFloat64IntfV(rv.Interface().(map[float64]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64IntfV(v map[float64]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[float64(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64StringR(rv reflect.Value) {
	fastpathTV.EncMapFloat64StringV(rv.Interface().(map[float64]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64StringV(v map[float64]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v[float64(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64UintR(rv reflect.Value) {
	fastpathTV.EncMapFloat64UintV(rv.Interface().(map[float64]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64UintV(v map[float64]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Uint8R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Uint8V(rv.Interface().(map[float64]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Uint8V(v map[float64]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Uint16R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Uint16V(rv.Interface().(map[float64]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Uint16V(v map[float64]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Uint32R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Uint32V(rv.Interface().(map[float64]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Uint32V(v map[float64]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Uint64R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Uint64V(rv.Interface().(map[float64]uint64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Uint64V(v map[float64]uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64UintptrR(rv reflect.Value) {
	fastpathTV.EncMapFloat64UintptrV(rv.Interface().(map[float64]uintptr), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64UintptrV(v map[float64]uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[float64(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64IntR(rv reflect.Value) {
	fastpathTV.EncMapFloat64IntV(rv.Interface().(map[float64]int), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64IntV(v map[float64]int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Int8R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Int8V(rv.Interface().(map[float64]int8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Int8V(v map[float64]int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Int16R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Int16V(rv.Interface().(map[float64]int16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Int16V(v map[float64]int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Int32R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Int32V(rv.Interface().(map[float64]int32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Int32V(v map[float64]int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Int64R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Int64V(rv.Interface().(map[float64]int64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Int64V(v map[float64]int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[float64(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Float32R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Float32V(rv.Interface().(map[float64]float32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Float32V(v map[float64]float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v[float64(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64Float64R(rv reflect.Value) {
	fastpathTV.EncMapFloat64Float64V(rv.Interface().(map[float64]float64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64Float64V(v map[float64]float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v[float64(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapFloat64BoolR(rv reflect.Value) {
	fastpathTV.EncMapFloat64BoolV(rv.Interface().(map[float64]bool), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapFloat64BoolV(v map[float64]bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]float64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = float64(k)
			i++
		}
		sort.Sort(floatSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(float64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v[float64(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeFloat64(k2)
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintIntfR(rv reflect.Value) {
	fastpathTV.EncMapUintIntfV(rv.Interface().(map[uint]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintIntfV(v map[uint]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintStringR(rv reflect.Value) {
	fastpathTV.EncMapUintStringV(rv.Interface().(map[uint]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintStringV(v map[uint]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintUintR(rv reflect.Value) {
	fastpathTV.EncMapUintUintV(rv.Interface().(map[uint]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintUintV(v map[uint]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintUint8R(rv reflect.Value) {
	fastpathTV.EncMapUintUint8V(rv.Interface().(map[uint]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintUint8V(v map[uint]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintUint16R(rv reflect.Value) {
	fastpathTV.EncMapUintUint16V(rv.Interface().(map[uint]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintUint16V(v map[uint]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintUint32R(rv reflect.Value) {
	fastpathTV.EncMapUintUint32V(rv.Interface().(map[uint]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintUint32V(v map[uint]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintUint64R(rv reflect.Value) {
	fastpathTV.EncMapUintUint64V(rv.Interface().(map[uint]uint64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintUint64V(v map[uint]uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintUintptrR(rv reflect.Value) {
	fastpathTV.EncMapUintUintptrV(rv.Interface().(map[uint]uintptr), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintUintptrV(v map[uint]uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintIntR(rv reflect.Value) {
	fastpathTV.EncMapUintIntV(rv.Interface().(map[uint]int), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintIntV(v map[uint]int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintInt8R(rv reflect.Value) {
	fastpathTV.EncMapUintInt8V(rv.Interface().(map[uint]int8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintInt8V(v map[uint]int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintInt16R(rv reflect.Value) {
	fastpathTV.EncMapUintInt16V(rv.Interface().(map[uint]int16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintInt16V(v map[uint]int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintInt32R(rv reflect.Value) {
	fastpathTV.EncMapUintInt32V(rv.Interface().(map[uint]int32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintInt32V(v map[uint]int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintInt64R(rv reflect.Value) {
	fastpathTV.EncMapUintInt64V(rv.Interface().(map[uint]int64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintInt64V(v map[uint]int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintFloat32R(rv reflect.Value) {
	fastpathTV.EncMapUintFloat32V(rv.Interface().(map[uint]float32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintFloat32V(v map[uint]float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintFloat64R(rv reflect.Value) {
	fastpathTV.EncMapUintFloat64V(rv.Interface().(map[uint]float64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintFloat64V(v map[uint]float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUintBoolR(rv reflect.Value) {
	fastpathTV.EncMapUintBoolV(rv.Interface().(map[uint]bool), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUintBoolV(v map[uint]bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v[uint(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8IntfR(rv reflect.Value) {
	fastpathTV.EncMapUint8IntfV(rv.Interface().(map[uint8]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8IntfV(v map[uint8]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8StringR(rv reflect.Value) {
	fastpathTV.EncMapUint8StringV(rv.Interface().(map[uint8]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8StringV(v map[uint8]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8UintR(rv reflect.Value) {
	fastpathTV.EncMapUint8UintV(rv.Interface().(map[uint8]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8UintV(v map[uint8]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Uint8R(rv reflect.Value) {
	fastpathTV.EncMapUint8Uint8V(rv.Interface().(map[uint8]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Uint8V(v map[uint8]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Uint16R(rv reflect.Value) {
	fastpathTV.EncMapUint8Uint16V(rv.Interface().(map[uint8]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Uint16V(v map[uint8]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Uint32R(rv reflect.Value) {
	fastpathTV.EncMapUint8Uint32V(rv.Interface().(map[uint8]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Uint32V(v map[uint8]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Uint64R(rv reflect.Value) {
	fastpathTV.EncMapUint8Uint64V(rv.Interface().(map[uint8]uint64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Uint64V(v map[uint8]uint64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8UintptrR(rv reflect.Value) {
	fastpathTV.EncMapUint8UintptrV(rv.Interface().(map[uint8]uintptr), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8UintptrV(v map[uint8]uintptr, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8IntR(rv reflect.Value) {
	fastpathTV.EncMapUint8IntV(rv.Interface().(map[uint8]int), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8IntV(v map[uint8]int, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Int8R(rv reflect.Value) {
	fastpathTV.EncMapUint8Int8V(rv.Interface().(map[uint8]int8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Int8V(v map[uint8]int8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Int16R(rv reflect.Value) {
	fastpathTV.EncMapUint8Int16V(rv.Interface().(map[uint8]int16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Int16V(v map[uint8]int16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Int32R(rv reflect.Value) {
	fastpathTV.EncMapUint8Int32V(rv.Interface().(map[uint8]int32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Int32V(v map[uint8]int32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Int64R(rv reflect.Value) {
	fastpathTV.EncMapUint8Int64V(rv.Interface().(map[uint8]int64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Int64V(v map[uint8]int64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v[uint8(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeInt(int64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Float32R(rv reflect.Value) {
	fastpathTV.EncMapUint8Float32V(rv.Interface().(map[uint8]float32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Float32V(v map[uint8]float32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat32(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8Float64R(rv reflect.Value) {
	fastpathTV.EncMapUint8Float64V(rv.Interface().(map[uint8]float64), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8Float64V(v map[uint8]float64, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeFloat64(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint8BoolR(rv reflect.Value) {
	fastpathTV.EncMapUint8BoolV(rv.Interface().(map[uint8]bool), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint8BoolV(v map[uint8]bool, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint8(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v[uint8(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeBool(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16IntfR(rv reflect.Value) {
	fastpathTV.EncMapUint16IntfV(rv.Interface().(map[uint16]interface{}), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint16IntfV(v map[uint16]interface{}, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint16(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v[uint16(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			e.encode(v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16StringR(rv reflect.Value) {
	fastpathTV.EncMapUint16StringV(rv.Interface().(map[uint16]string), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint16StringV(v map[uint16]string, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint16(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v[uint16(k2)])
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeString(c_UTF8, v2)
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16UintR(rv reflect.Value) {
	fastpathTV.EncMapUint16UintV(rv.Interface().(map[uint16]uint), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint16UintV(v map[uint16]uint, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint16(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint16(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16Uint8R(rv reflect.Value) {
	fastpathTV.EncMapUint16Uint8V(rv.Interface().(map[uint16]uint8), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint16Uint8V(v map[uint16]uint8, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint16(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint16(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16Uint16R(rv reflect.Value) {
	fastpathTV.EncMapUint16Uint16V(rv.Interface().(map[uint16]uint16), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint16Uint16V(v map[uint16]uint16, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint16(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint16(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16Uint32R(rv reflect.Value) {
	fastpathTV.EncMapUint16Uint32V(rv.Interface().(map[uint16]uint32), fastpathCheckNilFalse, f.e)
}
func (_ fastpathT) EncMapUint16Uint32V(v map[uint16]uint32, checkNil bool, e *Encoder) {
	ee := e.e
	cr := e.cr
	if checkNil && v == nil {
		ee.EncodeNil()
		return
	}
	ee.EncodeMapStart(len(v))
	if e.h.Canonical {
		v2 := make([]uint64, len(v))
		var i int
		for k, _ := range v {
			v2[i] = uint64(k)
			i++
		}
		sort.Sort(uintSlice(v2))
		for _, k2 := range v2 {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(uint16(k2)))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v[uint16(k2)]))
		}
	} else {
		for k2, v2 := range v {
			if cr != nil {
				cr.sendContainerState(containerMapKey)
			}
			ee.EncodeUint(uint64(k2))
			if cr != nil {
				cr.sendContainerState(containerMapValue)
			}
			ee.EncodeUint(uint64(v2))
		}
	}
	if cr != nil {
		cr.sendContainerState(containerMapEnd)
	}
}

func (f *encFnInfo) fastpathEncMapUint16Uint64R(rv reflect.Value) {
	fastpathTV.EncMapUint16Uint64V(rv.Interface().(map[uint16]uint64), fastpathCheckNilFalse, f.e)
}
f