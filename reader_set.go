package konfig

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (r *reader) setString(v reflect.Value, name, val string) (bool, error) {
	if v.String() == val {
		return false, nil
	}

	r.log(5, "[%s] setting string value: %s", name, val)
	v.SetString(val)
	r.notifySubscribers(name, val)

	return true, nil
}

func (r *reader) setBool(v reflect.Value, name, val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	if v.Bool() == b {
		return false, nil
	}

	r.log(5, "[%s] setting bool value: %t", name, b)
	v.SetBool(b)
	r.notifySubscribers(name, b)

	return true, nil
}

func (r *reader) setFloat32(v reflect.Value, name, val string) (bool, error) {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return false, err
	}

	if v.Float() == f {
		return false, nil
	}

	r.log(5, "[%s] setting float32 value: %f", name, f)
	v.SetFloat(f)
	r.notifySubscribers(name, float32(f))

	return true, nil
}

func (r *reader) setFloat64(v reflect.Value, name, val string) (bool, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false, err
	}

	if v.Float() == f {
		return false, nil
	}

	r.log(5, "[%s] setting float64 value: %f", name, f)
	v.SetFloat(f)
	r.notifySubscribers(name, f)

	return true, nil
}

func (r *reader) setInt(v reflect.Value, name, val string) (bool, error) {
	// int size and range are platform-dependent
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	r.log(5, "[%s] setting int value: %d", name, i)
	v.SetInt(i)
	r.notifySubscribers(name, int(i))

	return true, nil
}

func (r *reader) setInt8(v reflect.Value, name, val string) (bool, error) {
	i, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	r.log(5, "[%s] setting int8 value: %d", name, i)
	v.SetInt(i)
	r.notifySubscribers(name, int8(i))

	return true, nil
}

func (r *reader) setInt16(v reflect.Value, name, val string) (bool, error) {
	i, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	r.log(5, "[%s] setting int16 value: %d", name, i)
	v.SetInt(i)
	r.notifySubscribers(name, int16(i))

	return true, nil
}

func (r *reader) setInt32(v reflect.Value, name, val string) (bool, error) {
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	r.log(5, "[%s] setting int32 value: %d", name, i)
	v.SetInt(i)
	r.notifySubscribers(name, int32(i))

	return true, nil
}

func (r *reader) setInt64(v reflect.Value, name, val string) (bool, error) {
	if t := v.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
		d, err := time.ParseDuration(val)
		if err != nil {
			return false, err
		}

		if v.Interface() == d {
			return false, nil
		}

		r.log(5, "[%s] setting duration value: %s", name, d)
		v.Set(reflect.ValueOf(d))
		r.notifySubscribers(name, d)

		return true, nil
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Int() == i {
		return false, nil
	}

	r.log(5, "[%s] setting int64 value: %d", name, i)
	v.SetInt(i)
	r.notifySubscribers(name, i)

	return true, nil
}

func (r *reader) setUint(v reflect.Value, name, val string) (bool, error) {
	// uint size and range are platform-dependent
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	r.log(5, "[%s] setting uint value: %d", name, u)
	v.SetUint(u)
	r.notifySubscribers(name, uint(u))

	return true, nil
}

func (r *reader) setUint8(v reflect.Value, name, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	r.log(5, "[%s] setting uint8 value: %d", name, u)
	v.SetUint(u)
	r.notifySubscribers(name, uint8(u))

	return true, nil
}

func (r *reader) setUint16(v reflect.Value, name, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	r.log(5, "[%s] setting uint16 value: %d", name, u)
	v.SetUint(u)
	r.notifySubscribers(name, uint16(u))

	return true, nil
}

func (r *reader) setUint32(v reflect.Value, name, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	r.log(5, "[%s] setting uint32 value: %d", name, u)
	v.SetUint(u)
	r.notifySubscribers(name, uint32(u))

	return true, nil
}

func (r *reader) setUint64(v reflect.Value, name, val string) (bool, error) {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if v.Uint() == u {
		return false, nil
	}

	r.log(5, "[%s] setting unsigned integer value: %d", name, u)
	v.SetUint(u)
	r.notifySubscribers(name, u)

	return true, nil
}

func (r *reader) setStruct(v reflect.Value, name, val string) (bool, error) {
	t := v.Type()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		u, err := url.Parse(val)
		if err != nil {
			return false, err
		}

		// u is a pointer
		if reflect.DeepEqual(v.Interface(), *u) {
			return false, nil
		}

		// u is a pointer
		r.log(5, "[%s] setting url value: %s", name, val)
		v.Set(reflect.ValueOf(u).Elem())
		r.notifySubscribers(name, *u)

		return true, nil
	} else if t.PkgPath() == "regexp" && t.Name() == "Regexp" {
		re, err := regexp.CompilePOSIX(val)
		if err != nil {
			return false, err
		}

		// r is a pointer
		if reflect.DeepEqual(v.Interface(), *re) {
			return false, nil
		}

		// r is a pointer
		r.log(5, "[%s] setting regexp value: %s", name, val)
		v.Set(reflect.ValueOf(re).Elem())
		r.notifySubscribers(name, *re)

		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

func (r *reader) setStringPtr(v reflect.Value, name, val string) (bool, error) {
	if !v.IsZero() && v.Elem().String() == val {
		return false, nil
	}

	r.log(5, "[%s] setting string pointer: %s", name, val)
	v.Set(reflect.ValueOf(&val))
	r.notifySubscribers(name, &val)

	return true, nil
}

func (r *reader) setBoolPtr(v reflect.Value, name, val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Bool() == b {
		return false, nil
	}

	r.log(5, "[%s] setting bool pointer: %t", name, b)
	v.Set(reflect.ValueOf(&b))
	r.notifySubscribers(name, &b)

	return true, nil
}

func (r *reader) setFloat32Ptr(v reflect.Value, name, val string) (bool, error) {
	f64, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Float() == f64 {
		return false, nil
	}

	f32 := float32(f64)
	r.log(5, "[%s] setting float32 pointer: %f", name, f32)
	v.Set(reflect.ValueOf(&f32))
	r.notifySubscribers(name, &f32)

	return true, nil
}

func (r *reader) setFloat64Ptr(v reflect.Value, name, val string) (bool, error) {
	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Float() == f64 {
		return false, nil
	}

	r.log(5, "[%s] setting float64 pointer: %f", name, f64)
	v.Set(reflect.ValueOf(&f64))
	r.notifySubscribers(name, &f64)

	return true, nil
}

func (r *reader) setIntPtr(v reflect.Value, name, val string) (bool, error) {
	// int size and range are platform-dependent
	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i := int(i64)
	r.log(5, "[%s] setting int pointer: %d", name, i)
	v.Set(reflect.ValueOf(&i))
	r.notifySubscribers(name, &i)

	return true, nil
}

func (r *reader) setInt8Ptr(v reflect.Value, name, val string) (bool, error) {
	i64, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i8 := int8(i64)
	r.log(5, "[%s] setting int8 pointer: %d", name, i8)
	v.Set(reflect.ValueOf(&i8))
	r.notifySubscribers(name, &i8)

	return true, nil
}

func (r *reader) setInt16Ptr(v reflect.Value, name, val string) (bool, error) {
	i64, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i16 := int16(i64)
	r.log(5, "[%s] setting int16 pointer: %d", name, i16)
	v.Set(reflect.ValueOf(&i16))
	r.notifySubscribers(name, &i16)

	return true, nil
}

func (r *reader) setInt32Ptr(v reflect.Value, name, val string) (bool, error) {
	i64, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	i32 := int32(i64)
	r.log(5, "[%s] setting int32 pointer: %d", name, i32)
	v.Set(reflect.ValueOf(&i32))
	r.notifySubscribers(name, &i32)

	return true, nil
}

func (r *reader) setInt64Ptr(v reflect.Value, name, val string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "time" && t.Name() == "Duration" {
		d, err := time.ParseDuration(val)
		if err != nil {
			return false, err
		}

		if !v.IsZero() && v.Elem().Interface() == d {
			return false, nil
		}

		r.log(5, "[%s] setting duration pointer: %s", name, d)
		v.Set(reflect.ValueOf(&d))
		r.notifySubscribers(name, &d)

		return true, nil
	}

	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Int() == i64 {
		return false, nil
	}

	r.log(5, "[%s] setting int64 pointer: %d", name, i64)
	v.Set(reflect.ValueOf(&i64))
	r.notifySubscribers(name, &i64)

	return true, nil
}

func (r *reader) setUintPtr(v reflect.Value, name, val string) (bool, error) {
	// uint size and range are platform-dependent
	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u := uint(u64)
	r.log(5, "[%s] setting uint pointer: %d", name, u)
	v.Set(reflect.ValueOf(&u))
	r.notifySubscribers(name, &u)

	return true, nil
}

func (r *reader) setUint8Ptr(v reflect.Value, name, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u8 := uint8(u64)
	r.log(5, "[%s] setting uint8 pointer: %d", name, u8)
	v.Set(reflect.ValueOf(&u8))
	r.notifySubscribers(name, &u8)

	return true, nil
}

func (r *reader) setUint16Ptr(v reflect.Value, name, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u16 := uint16(u64)
	r.log(5, "[%s] setting uint16 pointer: %d", name, u16)
	v.Set(reflect.ValueOf(&u16))
	r.notifySubscribers(name, &u16)

	return true, nil
}

func (r *reader) setUint32Ptr(v reflect.Value, name, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	u32 := uint32(u64)
	r.log(5, "[%s] setting uint32 pointer: %d", name, u32)
	v.Set(reflect.ValueOf(&u32))
	r.notifySubscribers(name, &u32)

	return true, nil
}

func (r *reader) setUint64Ptr(v reflect.Value, name, val string) (bool, error) {
	u64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return false, err
	}

	if !v.IsZero() && v.Elem().Uint() == u64 {
		return false, nil
	}

	r.log(5, "[%s] setting uint pointer: %d", name, u64)
	v.Set(reflect.ValueOf(&u64))
	r.notifySubscribers(name, &u64)

	return true, nil
}

func (r *reader) setStructPtr(v reflect.Value, name, val string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		u, err := url.Parse(val)
		if err != nil {
			return false, err
		}

		if !v.IsZero() && reflect.DeepEqual(v.Elem().Interface(), *u) {
			return false, nil
		}

		// u is a pointer
		r.log(5, "[%s] setting url pointer: %s", name, val)
		v.Set(reflect.ValueOf(u))
		r.notifySubscribers(name, u)

		return true, nil
	} else if t.PkgPath() == "regexp" && t.Name() == "Regexp" {
		re, err := regexp.CompilePOSIX(val)
		if err != nil {
			return false, err
		}

		if !v.IsZero() && reflect.DeepEqual(v.Elem().Interface(), *re) {
			return false, nil
		}

		// r is a pointer
		r.log(5, "[%s] setting regexp pointer: %s", name, val)
		v.Set(reflect.ValueOf(re))
		r.notifySubscribers(name, re)

		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

func (r *reader) setStringSlice(v reflect.Value, name string, vals []string) (bool, error) {
	if reflect.DeepEqual(v.Interface(), vals) {
		return false, nil
	}

	r.log(5, "[%s] setting string slice: %v", name, vals)
	v.Set(reflect.ValueOf(vals))
	r.notifySubscribers(name, vals)

	return true, nil
}

func (r *reader) setBoolSlice(v reflect.Value, name string, vals []string) (bool, error) {
	bools := []bool{}
	for _, val := range vals {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, err
		}

		bools = append(bools, b)
	}

	if reflect.DeepEqual(v.Interface(), bools) {
		return false, nil
	}

	r.log(5, "[%s] setting bool slice: %v", name, bools)
	v.Set(reflect.ValueOf(bools))
	r.notifySubscribers(name, bools)

	return true, nil
}

func (r *reader) setFloat32Slice(v reflect.Value, name string, vals []string) (bool, error) {
	floats := []float32{}
	for _, val := range vals {
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return false, err
		}

		floats = append(floats, float32(f))
	}

	if reflect.DeepEqual(v.Interface(), floats) {
		return false, nil
	}

	r.log(5, "[%s] setting float32 slice: %v", name, floats)
	v.Set(reflect.ValueOf(floats))
	r.notifySubscribers(name, floats)

	return true, nil
}

func (r *reader) setFloat64Slice(v reflect.Value, name string, vals []string) (bool, error) {
	floats := []float64{}
	for _, val := range vals {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return false, err
		}

		floats = append(floats, f)
	}

	if reflect.DeepEqual(v.Interface(), floats) {
		return false, nil
	}

	r.log(5, "[%s] setting float64 slice: %v", name, floats)
	v.Set(reflect.ValueOf(floats))
	r.notifySubscribers(name, floats)

	return true, nil
}

func (r *reader) setIntSlice(v reflect.Value, name string, vals []string) (bool, error) {
	// int size and range are platform-dependent
	ints := []int{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, err
		}

		ints = append(ints, int(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	r.log(5, "[%s] setting int slice: %v", name, ints)
	v.Set(reflect.ValueOf(ints))
	r.notifySubscribers(name, ints)

	return true, nil
}

func (r *reader) setInt8Slice(v reflect.Value, name string, vals []string) (bool, error) {
	ints := []int8{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return false, err
		}

		ints = append(ints, int8(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	r.log(5, "[%s] setting int8 slice: %v", name, ints)
	v.Set(reflect.ValueOf(ints))
	r.notifySubscribers(name, ints)

	return true, nil
}

func (r *reader) setInt16Slice(v reflect.Value, name string, vals []string) (bool, error) {
	ints := []int16{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return false, err
		}

		ints = append(ints, int16(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	r.log(5, "[%s] setting int16 slice: %v", name, ints)
	v.Set(reflect.ValueOf(ints))
	r.notifySubscribers(name, ints)

	return true, nil
}

func (r *reader) setInt32Slice(v reflect.Value, name string, vals []string) (bool, error) {
	ints := []int32{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return false, err
		}

		ints = append(ints, int32(i))
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	r.log(5, "[%s] setting int32 slice: %v", name, ints)
	v.Set(reflect.ValueOf(ints))
	r.notifySubscribers(name, ints)

	return true, nil
}

func (r *reader) setInt64Slice(v reflect.Value, name string, vals []string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "time" && t.Name() == "Duration" {
		durations := []time.Duration{}
		for _, val := range vals {
			d, err := time.ParseDuration(val)
			if err != nil {
				return false, err
			}

			durations = append(durations, d)
		}

		if reflect.DeepEqual(v.Interface(), durations) {
			return false, nil
		}

		r.log(5, "[%s] setting duration slice: %v", name, durations)
		v.Set(reflect.ValueOf(durations))
		r.notifySubscribers(name, durations)

		return true, nil
	}

	ints := []int64{}
	for _, val := range vals {
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, err
		}

		ints = append(ints, i)
	}

	if reflect.DeepEqual(v.Interface(), ints) {
		return false, nil
	}

	r.log(5, "[%s] setting int64 slice: %v", name, ints)
	v.Set(reflect.ValueOf(ints))
	r.notifySubscribers(name, ints)

	return true, nil
}

func (r *reader) setUintSlice(v reflect.Value, name string, vals []string) (bool, error) {
	// uint size and range are platform-dependent
	uints := []uint{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	r.log(5, "[%s] setting uint slice: %v", name, uints)
	v.Set(reflect.ValueOf(uints))
	r.notifySubscribers(name, uints)

	return true, nil
}

func (r *reader) setUint8Slice(v reflect.Value, name string, vals []string) (bool, error) {
	uints := []uint8{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint8(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	r.log(5, "[%s] setting uint8 slice: %v", name, uints)
	v.Set(reflect.ValueOf(uints))
	r.notifySubscribers(name, uints)

	return true, nil
}

func (r *reader) setUint16Slice(v reflect.Value, name string, vals []string) (bool, error) {
	uints := []uint16{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint16(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	r.log(5, "[%s] setting uint16 slice: %v", name, uints)
	v.Set(reflect.ValueOf(uints))
	r.notifySubscribers(name, uints)

	return true, nil
}

func (r *reader) setUint32Slice(v reflect.Value, name string, vals []string) (bool, error) {
	uints := []uint32{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return false, err
		}

		uints = append(uints, uint32(u))
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	r.log(5, "[%s] setting uint32 slice: %v", name, uints)
	v.Set(reflect.ValueOf(uints))
	r.notifySubscribers(name, uints)

	return true, nil
}

func (r *reader) setUint64Slice(v reflect.Value, name string, vals []string) (bool, error) {
	uints := []uint64{}
	for _, val := range vals {
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return false, err
		}

		uints = append(uints, u)
	}

	if reflect.DeepEqual(v.Interface(), uints) {
		return false, nil
	}

	r.log(5, "[%s] setting uint64 slice: %v", name, uints)
	v.Set(reflect.ValueOf(uints))
	r.notifySubscribers(name, uints)

	return true, nil
}

func (r *reader) setStructSlice(v reflect.Value, name string, vals []string) (bool, error) {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		urls := []url.URL{}
		for _, val := range vals {
			u, err := url.Parse(val)
			if err != nil {
				return false, err
			}

			urls = append(urls, *u)
		}

		// []url.URL
		if reflect.DeepEqual(v.Interface(), urls) {
			return false, nil
		}

		r.log(5, "[%s] setting url slice: %v", name, urls)
		v.Set(reflect.ValueOf(urls))
		r.notifySubscribers(name, urls)

		return true, nil
	} else if t.PkgPath() == "regexp" && t.Name() == "Regexp" {
		regexps := []regexp.Regexp{}
		for _, val := range vals {
			r, err := regexp.CompilePOSIX(val)
			if err != nil {
				return false, err
			}

			regexps = append(regexps, *r)
		}

		// []regexp.Regexp
		if reflect.DeepEqual(v.Interface(), regexps) {
			return false, nil
		}

		r.log(5, "[%s] setting regexp slice: %v", name, regexps)
		v.Set(reflect.ValueOf(regexps))
		r.notifySubscribers(name, regexps)

		return true, nil
	}

	return false, fmt.Errorf("unsupported type: %s.%s", t.PkgPath(), t.Name())
}

func (r *reader) setFieldValue(f fieldInfo, val string) (bool, error) {
	switch f.value.Kind() {
	case reflect.String:
		return r.setString(f.value, f.name, val)
	case reflect.Bool:
		return r.setBool(f.value, f.name, val)
	case reflect.Float32:
		return r.setFloat32(f.value, f.name, val)
	case reflect.Float64:
		return r.setFloat64(f.value, f.name, val)
	case reflect.Int:
		return r.setInt(f.value, f.name, val)
	case reflect.Int8:
		return r.setInt8(f.value, f.name, val)
	case reflect.Int16:
		return r.setInt16(f.value, f.name, val)
	case reflect.Int32:
		return r.setInt32(f.value, f.name, val)
	case reflect.Int64:
		return r.setInt64(f.value, f.name, val)
	case reflect.Uint:
		return r.setUint(f.value, f.name, val)
	case reflect.Uint8:
		return r.setUint8(f.value, f.name, val)
	case reflect.Uint16:
		return r.setUint16(f.value, f.name, val)
	case reflect.Uint32:
		return r.setUint32(f.value, f.name, val)
	case reflect.Uint64:
		return r.setUint64(f.value, f.name, val)
	case reflect.Struct:
		return r.setStruct(f.value, f.name, val)

	case reflect.Ptr:
		tPtr := reflect.TypeOf(f.value.Interface()).Elem()

		switch tPtr.Kind() {
		case reflect.String:
			return r.setStringPtr(f.value, f.name, val)
		case reflect.Bool:
			return r.setBoolPtr(f.value, f.name, val)
		case reflect.Float32:
			return r.setFloat32Ptr(f.value, f.name, val)
		case reflect.Float64:
			return r.setFloat64Ptr(f.value, f.name, val)
		case reflect.Int:
			return r.setIntPtr(f.value, f.name, val)
		case reflect.Int8:
			return r.setInt8Ptr(f.value, f.name, val)
		case reflect.Int16:
			return r.setInt16Ptr(f.value, f.name, val)
		case reflect.Int32:
			return r.setInt32Ptr(f.value, f.name, val)
		case reflect.Int64:
			return r.setInt64Ptr(f.value, f.name, val)
		case reflect.Uint:
			return r.setUintPtr(f.value, f.name, val)
		case reflect.Uint8:
			return r.setUint8Ptr(f.value, f.name, val)
		case reflect.Uint16:
			return r.setUint16Ptr(f.value, f.name, val)
		case reflect.Uint32:
			return r.setUint32Ptr(f.value, f.name, val)
		case reflect.Uint64:
			return r.setUint64Ptr(f.value, f.name, val)
		case reflect.Struct:
			return r.setStructPtr(f.value, f.name, val)
		}

	case reflect.Slice:
		tSlice := reflect.TypeOf(f.value.Interface()).Elem()
		vals := strings.Split(val, f.listSep)

		switch tSlice.Kind() {
		case reflect.String:
			return r.setStringSlice(f.value, f.name, vals)
		case reflect.Bool:
			return r.setBoolSlice(f.value, f.name, vals)
		case reflect.Float32:
			return r.setFloat32Slice(f.value, f.name, vals)
		case reflect.Float64:
			return r.setFloat64Slice(f.value, f.name, vals)
		case reflect.Int:
			return r.setIntSlice(f.value, f.name, vals)
		case reflect.Int8:
			return r.setInt8Slice(f.value, f.name, vals)
		case reflect.Int16:
			return r.setInt16Slice(f.value, f.name, vals)
		case reflect.Int32:
			return r.setInt32Slice(f.value, f.name, vals)
		case reflect.Int64:
			return r.setInt64Slice(f.value, f.name, vals)
		case reflect.Uint:
			return r.setUintSlice(f.value, f.name, vals)
		case reflect.Uint8:
			return r.setUint8Slice(f.value, f.name, vals)
		case reflect.Uint16:
			return r.setUint16Slice(f.value, f.name, vals)
		case reflect.Uint32:
			return r.setUint32Slice(f.value, f.name, vals)
		case reflect.Uint64:
			return r.setUint64Slice(f.value, f.name, vals)
		case reflect.Struct:
			return r.setStructSlice(f.value, f.name, vals)
		}
	}

	return false, fmt.Errorf("unsupported kind: %s", f.value.Kind())
}
