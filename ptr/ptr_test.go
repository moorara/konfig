package ptr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var s = "string"
	assert.Equal(t, &s, String(s))
}

func TestBool(t *testing.T) {
	var b = true
	assert.Equal(t, &b, Bool(b))
}

func TestFloat32(t *testing.T) {
	var f = float32(3.1415)
	assert.Equal(t, &f, Float32(f))
}

func TestFloat64(t *testing.T) {
	var f = float64(3.14159265359)
	assert.Equal(t, &f, Float64(f))
}

func TestInt(t *testing.T) {
	var i = int(-2147483648)
	assert.Equal(t, &i, Int(i))
}

func TestInt8(t *testing.T) {
	var i = int8(-128)
	assert.Equal(t, &i, Int8(i))
}

func TestInt16(t *testing.T) {
	var i = int16(-32768)
	assert.Equal(t, &i, Int16(i))
}

func TestInt32(t *testing.T) {
	var i = int32(-2147483648)
	assert.Equal(t, &i, Int32(i))
}

func TestInt64(t *testing.T) {
	var i = int64(-9223372036854775808)
	assert.Equal(t, &i, Int64(i))
}

func TestUint(t *testing.T) {
	var u = uint(4294967295)
	assert.Equal(t, &u, Uint(u))
}

func TestUint8(t *testing.T) {
	var u = uint8(255)
	assert.Equal(t, &u, Uint8(u))
}

func TestUint16(t *testing.T) {
	var u = uint16(65535)
	assert.Equal(t, &u, Uint16(u))
}

func TestUint32(t *testing.T) {
	var u = uint32(4294967295)
	assert.Equal(t, &u, Uint32(u))
}

func TestUint64(t *testing.T) {
	var u = uint64(18446744073709551615)
	assert.Equal(t, &u, Uint64(u))
}

func TestDuration(t *testing.T) {
	var d = time.Second
	assert.Equal(t, &d, Duration(d))
}
