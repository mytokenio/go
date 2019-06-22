package config

import (
	"testing"
)

func TestHelper(t *testing.T) {
	var i8 int8 = 108
	var i16 int16 = 1016
	var i32 int32 = 1032
	var i64 int64 = 1064
	var f32 float32 = 2032
	var f64 float64 = 2064

	var s string = "3001"

	//toInt
	assert(t, toInt(i8, 11), 108)
	assert(t, toInt(i16, 11), 1016)
	assert(t, toInt(i32, 11), 1032)
	assert(t, toInt(i64, 11), 1064)
	assert(t, toInt(f32, 11), 2032)
	assert(t, toInt(f64, 11), 2064)
	assert(t, toInt(s, 11), 3001)
	assert(t, toInt("xxx", 11), 11)

	//toInt64
	assert(t, toInt64(i8, 11), int64(108))
	assert(t, toInt64(i16, 11), int64(1016))
	assert(t, toInt64(i32, 11), int64(1032))
	assert(t, toInt64(i64, 11), int64(1064))
	assert(t, toInt64(f32, 11), int64(2032))
	assert(t, toInt64(f64, 11), int64(2064))
	assert(t, toInt64(s, 11), int64(3001))
	assert(t, toInt64("xxx", 11), int64(11))

	//toFloat
	var ff32 float32 = 1234.5678
	var ff64 float64 = 234.567
	var fs string = "333.567"
	assert(t, toFloat64(i8, 11), float64(108))
	assert(t, toFloat64(i16, 11), float64(1016))
	assert(t, toFloat64(i32, 11), float64(1032))
	assert(t, toFloat64(i64, 11), float64(1064))
	assert(t, toFloat64(f32, 11), float64(2032))
	assert(t, toFloat64(f64, 11), float64(2064))
	assert(t, toFloat64(ff32, 11), float64(1234.5677490234375)) //fuck float32
	assert(t, toFloat64(ff64, 11), float64(234.567))
	assert(t, toFloat64(fs, 11), float64(333.567))
	assert(t, toFloat64("xxxx", 11), float64(11))

	//toBool
	var bt string = "true"
	var bt2 string = "True"
	var bf string = "false"
	var bf2 string = "False"
	assert(t, toBool(bt, false), true)
	assert(t, toBool(bt2, false), true)
	assert(t, toBool(bf, true), false)
	assert(t, toBool(bf2, true), false)
	assert(t, toBool(123, true), true)
	assert(t, toBool(3222, false), false)
}

