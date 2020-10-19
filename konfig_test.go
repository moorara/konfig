package konfig

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/konfig/ptr"
)

type config struct {
	sync.Mutex
	unexported      string
	SkipFlag        string          `flag:"-"`
	SkipFlagEnv     string          `flag:"-" env:"-"`
	SkipFlagEnvFile string          `flag:"-" env:"-" fileenv:"-"`
	String          string          // `flag:"string" env:"STRING" fileenv:"STRING_FILE"`
	Bool            bool            // `flag:"bool" env:"BOOL" fileenv:"BOOL_FILE"`
	Float32         float32         // `flag:"float32" env:"FLOAT32" fileenv:"FLOAT32_FILE"`
	Float64         float64         // `flag:"float64" env:"FLOAT64" fileenv:"FLOAT64_FILE"`
	Int             int             // `flag:"int" env:"INT" fileenv:"INT_FILE"`
	Int8            int8            // `flag:"int8" env:"INT8" fileenv:"INT8_FILE"`
	Int16           int16           // `flag:"int16" env:"INT16" fileenv:"INT16_FILE"`
	Int32           int32           // `flag:"int32" env:"INT32" fileenv:"INT32_FILE"`
	Int64           int64           // `flag:"int64" env:"INT64" fileenv:"INT64_FILE"`
	Uint            uint            // `flag:"uint" env:"UINT" fileenv:"UINT_FILE"`
	Uint8           uint8           // `flag:"uint8" env:"UINT8" fileenv:"UINT8_FILE"`
	Uint16          uint16          // `flag:"uint16" env:"UINT16" fileenv:"UINT16_FILE"`
	Uint32          uint32          // `flag:"uint32" env:"UINT32" fileenv:"UINT32_FILE"`
	Uint64          uint64          // `flag:"uint64" env:"UINT64" fileenv:"UINT64_FILE"`
	URL             url.URL         // `flag:"url" env:"URL" fileenv:"URL_FILE"`
	Regexp          regexp.Regexp   // `flag:"regexp" env:"REGEXP" fileenv:"REGEXP_FILE"`
	Duration        time.Duration   // `flag:"duration" env:"DURATION" fileenv:"DURATION_FILE"`
	StringPointer   *string         // `flag:"string.pointer" env:"STRING_POINTER" fileenv:"STRING_POINTER_FILE"`
	BoolPointer     *bool           // `flag:"bool.pointer" env:"BOOL_POINTER" fileenv:"BOOL_POINTER_FILE"`
	Float32Pointer  *float32        // `flag:"float32.pointer" env:"FLOAT32_POINTER" fileenv:"FLOAT32_POINTER_FILE"`
	Float64Pointer  *float64        // `flag:"float64.pointer" env:"FLOAT64_POINTER" fileenv:"FLOAT64_POINTER_FILE"`
	IntPointer      *int            // `flag:"int.pointer" env:"INT_POINTER" fileenv:"INT_POINTER_FILE"`
	Int8Pointer     *int8           // `flag:"int8.pointer" env:"INT8_POINTER" fileenv:"INT8_POINTER_FILE"`
	Int16Pointer    *int16          // `flag:"int16.pointer" env:"INT16_POINTER" fileenv:"INT16_POINTER_FILE"`
	Int32Pointer    *int32          // `flag:"int32.pointer" env:"INT32_POINTER" fileenv:"INT32_POINTER_FILE"`
	Int64Pointer    *int64          // `flag:"int64.pointer" env:"INT64_POINTER" fileenv:"INT64_POINTER_FILE"`
	UintPointer     *uint           // `flag:"uint.pointer" env:"UINT_POINTER" fileenv:"UINT_POINTER_FILE"`
	Uint8Pointer    *uint8          // `flag:"uint8.pointer" env:"UINT8_POINTER" fileenv:"UINT8_POINTER_FILE"`
	Uint16Pointer   *uint16         // `flag:"uint16.pointer" env:"UINT16_POINTER" fileenv:"UINT16_POINTER_FILE"`
	Uint32Pointer   *uint32         // `flag:"uint32.pointer" env:"UINT32_POINTER" fileenv:"UINT32_POINTER_FILE"`
	Uint64Pointer   *uint64         // `flag:"uint64.pointer" env:"UINT64_POINTER" fileenv:"UINT64_POINTER_FILE"`
	URLPointer      *url.URL        // `flag:"url.pointer" env:"URL_POINTER" fileenv:"URL_POINTER_FILE"`
	RegexpPointer   *regexp.Regexp  // `flag:"regexp.pointer" env:"REGEXP_POINTER" fileenv:"REGEXP_POINTER_FILE"`
	DurationPointer *time.Duration  // `flag:"duration.pointer" env:"DURATION_POINTER" fileenv:"DURATION_POINTER_FILE"`
	StringSlice     []string        // `flag:"string.slice" env:"STRING_SLICE" fileenv:"STRING_SLICE_FILE" sep:","`
	BoolSlice       []bool          // `flag:"bool.slice" env:"BOOL_SLICE" fileenv:"BOOL_SLICE_FILE" sep:","`
	Float32Slice    []float32       // `flag:"float32.slice" env:"FLOAT32_SLICE" fileenv:"FLOAT32_SLICE_FILE" sep:","`
	Float64Slice    []float64       // `flag:"float64.slice" env:"FLOAT64_SLICE" fileenv:"FLOAT64_SLICE_FILE" sep:","`
	IntSlice        []int           // `flag:"int.slice" env:"INT_SLICE" fileenv:"INT_SLICE_FILE" sep:","`
	Int8Slice       []int8          // `flag:"int8.slice" env:"INT8_SLICE" fileenv:"INT8_SLICE_FILE" sep:","`
	Int16Slice      []int16         // `flag:"int16.slice" env:"INT16_SLICE" fileenv:"INT16_SLICE_FILE" sep:","`
	Int32Slice      []int32         // `flag:"int32.slice" env:"INT32_SLICE" fileenv:"INT32_SLICE_FILE" sep:","`
	Int64Slice      []int64         // `flag:"int64.slice" env:"INT64_SLICE" fileenv:"INT64_SLICE_FILE" sep:","`
	UintSlice       []uint          // `flag:"uint.slice" env:"UINT_SLICE" fileenv:"UINT_SLICE_FILE" sep:","`
	Uint8Slice      []uint8         // `flag:"uint8.slice" env:"UINT8_SLICE" fileenv:"UINT8_SLICE_FILE" sep:","`
	Uint16Slice     []uint16        // `flag:"uint16.slice" env:"UINT16_SLICE" fileenv:"UINT16_SLICE_FILE" sep:","`
	Uint32Slice     []uint32        // `flag:"uint32.slice" env:"UINT32_SLICE" fileenv:"UINT32_SLICE_FILE" sep:","`
	Uint64Slice     []uint64        // `flag:"uint64.slice" env:"UINT64_SLICE" fileenv:"UINT64_SLICE_FILE" sep:","`
	URLSlice        []url.URL       // `flag:"url.slice" env:"URL_SLICE" fileenv:"URL_SLICE_FILE" sep:","`
	RegexpSlice     []regexp.Regexp // `flag:"regexp.slice" env:"REGEXP_SLICE" fileenv:"REGEXP_SLICE_FILE" sep:","`
	DurationSlice   []time.Duration // `flag:"duration.slice" env:"DURATION_SLICE" fileenv:"DURATION_SLICE_FILE" sep:","`
}

func (c1 *config) Equal(c2 *config) bool {
	return c1.unexported == c2.unexported &&
		c1.SkipFlag == c2.SkipFlag &&
		c1.SkipFlagEnv == c2.SkipFlagEnv &&
		c1.SkipFlagEnvFile == c2.SkipFlagEnvFile &&
		c1.String == c2.String &&
		c1.Bool == c2.Bool &&
		c1.Float32 == c2.Float32 &&
		c1.Float64 == c2.Float64 &&
		c1.Int == c2.Int &&
		c1.Int8 == c2.Int8 &&
		c1.Int16 == c2.Int16 &&
		c1.Int32 == c2.Int32 &&
		c1.Int64 == c2.Int64 &&
		c1.Uint == c2.Uint &&
		c1.Uint8 == c2.Uint8 &&
		c1.Uint16 == c2.Uint16 &&
		c1.Uint32 == c2.Uint32 &&
		c1.Uint64 == c2.Uint64 &&
		c1.URL == c2.URL &&
		reflect.DeepEqual(c1.Regexp, c2.Regexp) &&
		c1.Duration == c2.Duration &&
		reflect.DeepEqual(c1.StringPointer, c2.StringPointer) &&
		reflect.DeepEqual(c1.BoolPointer, c2.BoolPointer) &&
		reflect.DeepEqual(c1.Float32Pointer, c2.Float32Pointer) &&
		reflect.DeepEqual(c1.Float64Pointer, c2.Float64Pointer) &&
		reflect.DeepEqual(c1.IntPointer, c2.IntPointer) &&
		reflect.DeepEqual(c1.Int8Pointer, c2.Int8Pointer) &&
		reflect.DeepEqual(c1.Int16Pointer, c2.Int16Pointer) &&
		reflect.DeepEqual(c1.Int32Pointer, c2.Int32Pointer) &&
		reflect.DeepEqual(c1.Int64Pointer, c2.Int64Pointer) &&
		reflect.DeepEqual(c1.UintPointer, c2.UintPointer) &&
		reflect.DeepEqual(c1.Uint8Pointer, c2.Uint8Pointer) &&
		reflect.DeepEqual(c1.Uint16Pointer, c2.Uint16Pointer) &&
		reflect.DeepEqual(c1.Uint32Pointer, c2.Uint32Pointer) &&
		reflect.DeepEqual(c1.Uint64Pointer, c2.Uint64Pointer) &&
		reflect.DeepEqual(c1.URLPointer, c2.URLPointer) &&
		reflect.DeepEqual(c1.RegexpPointer, c2.RegexpPointer) &&
		reflect.DeepEqual(c1.DurationPointer, c2.DurationPointer) &&
		reflect.DeepEqual(c1.StringSlice, c2.StringSlice) &&
		reflect.DeepEqual(c1.BoolSlice, c2.BoolSlice) &&
		reflect.DeepEqual(c1.Float32Slice, c2.Float32Slice) &&
		reflect.DeepEqual(c1.Float64Slice, c2.Float64Slice) &&
		reflect.DeepEqual(c1.IntSlice, c2.IntSlice) &&
		reflect.DeepEqual(c1.Int8Slice, c2.Int8Slice) &&
		reflect.DeepEqual(c1.Int16Slice, c2.Int16Slice) &&
		reflect.DeepEqual(c1.Int32Slice, c2.Int32Slice) &&
		reflect.DeepEqual(c1.Int64Slice, c2.Int64Slice) &&
		reflect.DeepEqual(c1.UintSlice, c2.UintSlice) &&
		reflect.DeepEqual(c1.Uint8Slice, c2.Uint8Slice) &&
		reflect.DeepEqual(c1.Uint16Slice, c2.Uint16Slice) &&
		reflect.DeepEqual(c1.Uint32Slice, c2.Uint32Slice) &&
		reflect.DeepEqual(c1.Uint64Slice, c2.Uint64Slice) &&
		reflect.DeepEqual(c1.URLSlice, c2.URLSlice) &&
		reflect.DeepEqual(c1.RegexpSlice, c2.RegexpSlice) &&
		reflect.DeepEqual(c1.DurationSlice, c2.DurationSlice)
}

var (
	url1, _ = url.Parse("service-1")
	url2, _ = url.Parse("service-2")

	re1 = regexp.MustCompilePOSIX("[:digit:]")
	re2 = regexp.MustCompilePOSIX("[:alpha:]")

	cfg = config{
		String:          "foo",
		Bool:            true,
		Float32:         3.1415,
		Float64:         3.14159265359,
		Int:             9223372036854775807,
		Int8:            127,
		Int16:           32767,
		Int32:           2147483647,
		Int64:           9223372036854775807,
		Uint:            18446744073709551615,
		Uint8:           255,
		Uint16:          65535,
		Uint32:          4294967295,
		Uint64:          18446744073709551615,
		URL:             *url1,
		Regexp:          *re1,
		Duration:        time.Second,
		StringPointer:   ptr.String("foo"),
		BoolPointer:     ptr.Bool(true),
		Float32Pointer:  ptr.Float32(3.1415),
		Float64Pointer:  ptr.Float64(3.14159265359),
		IntPointer:      ptr.Int(9223372036854775807),
		Int8Pointer:     ptr.Int8(127),
		Int16Pointer:    ptr.Int16(32767),
		Int32Pointer:    ptr.Int32(2147483647),
		Int64Pointer:    ptr.Int64(9223372036854775807),
		UintPointer:     ptr.Uint(18446744073709551615),
		Uint8Pointer:    ptr.Uint8(255),
		Uint16Pointer:   ptr.Uint16(65535),
		Uint32Pointer:   ptr.Uint32(4294967295),
		Uint64Pointer:   ptr.Uint64(18446744073709551615),
		URLPointer:      url1,
		RegexpPointer:   re1,
		DurationPointer: ptr.Duration(time.Second),
		StringSlice:     []string{"foo", "bar"},
		BoolSlice:       []bool{false, true},
		Float32Slice:    []float32{3.1415, 2.7182},
		Float64Slice:    []float64{3.14159265359, 2.71828182845},
		IntSlice:        []int{-9223372036854775808, 9223372036854775807},
		Int8Slice:       []int8{-128, 127},
		Int16Slice:      []int16{-32768, 32767},
		Int32Slice:      []int32{-2147483648, 2147483647},
		Int64Slice:      []int64{-9223372036854775808, 9223372036854775807},
		UintSlice:       []uint{0, 18446744073709551615},
		Uint8Slice:      []uint8{0, 255},
		Uint16Slice:     []uint16{0, 65535},
		Uint32Slice:     []uint32{0, 4294967295},
		Uint64Slice:     []uint64{0, 18446744073709551615},
		URLSlice:        []url.URL{*url1, *url2},
		RegexpSlice:     []regexp.Regexp{*re1, *re2},
		DurationSlice:   []time.Duration{time.Second, time.Minute},
	}
)

func TestPick(t *testing.T) {
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName string
		value   string
	}

	tests := []struct {
		name           string
		args           []string
		envs           []env
		files          []file
		config         interface{}
		opts           []Option
		expectedError  error
		expectedConfig *config
	}{
		{
			"NonStruct",
			[]string{"app"},
			[]env{},
			[]file{},
			new(string),
			nil,
			errors.New("a non-struct type is passed"),
			&config{},
		},
		{
			"NonPointer",
			[]string{"app"},
			[]env{},
			[]file{},
			config{},
			nil,
			errors.New("a non-pointer type is passed"),
			&config{},
		},
		{
			"Empty",
			[]string{"app"},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{},
		},
		{
			"AllFromDefaults",
			[]string{"app"},
			[]env{},
			[]file{},
			&cfg,
			nil,
			nil,
			&cfg,
		},
		{
			"AllFromFlags",
			[]string{
				"app",
				"-string=foo",
				"-bool",
				"-float32=3.1415",
				"-float64=3.14159265359",
				"-int=9223372036854775807",
				"-int8=127",
				"-int16=32767",
				"-int32=2147483647",
				"-int64=9223372036854775807",
				"-uint=18446744073709551615",
				"-uint8=255",
				"-uint16=65535",
				"-uint32=4294967295",
				"-uint64=18446744073709551615",
				"-url=service-1",
				"-regexp=[:digit:]",
				"-duration=1s",
				"-string.pointer=foo",
				"-bool.pointer=true",
				"-float32.pointer=3.1415",
				"-float64.pointer=3.14159265359",
				"-int.pointer=9223372036854775807",
				"-int8.pointer=127",
				"-int16.pointer=32767",
				"-int32.pointer=2147483647",
				"-int64.pointer=9223372036854775807",
				"-uint.pointer=18446744073709551615",
				"-uint8.pointer=255",
				"-uint16.pointer=65535",
				"-uint32.pointer=4294967295",
				"-uint64.pointer=18446744073709551615",
				"-url.pointer=service-1",
				"-regexp.pointer=[:digit:]",
				"-duration.pointer=1s",
				"-string.slice=foo,bar",
				"-bool.slice=false,true",
				"-float32.slice=3.1415,2.7182",
				"-float64.slice=3.14159265359,2.71828182845",
				"-int.slice=-9223372036854775808,9223372036854775807",
				"-int8.slice=-128,127",
				"-int16.slice=-32768,32767",
				"-int32.slice=-2147483648,2147483647",
				"-int64.slice=-9223372036854775808,9223372036854775807",
				"-uint.slice=0,18446744073709551615",
				"-uint8.slice=0,255",
				"-uint16.slice=0,65535",
				"-uint32.slice=0,4294967295",
				"-uint64.slice=0,18446744073709551615",
				"-url.slice=service-1,service-2",
				"-regexp.slice=[:digit:],[:alpha:]",
				"-duration.slice=1s,1m",
			},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&cfg,
		},
		{
			"AllFromFlagsWithOptions",
			[]string{
				"app",
				"-config.string=foo",
				"-config.bool",
				"-config.float32=3.1415",
				"-config.float64=3.14159265359",
				"-config.int=9223372036854775807",
				"-config.int8=127",
				"-config.int16=32767",
				"-config.int32=2147483647",
				"-config.int64=9223372036854775807",
				"-config.uint=18446744073709551615",
				"-config.uint8=255",
				"-config.uint16=65535",
				"-config.uint32=4294967295",
				"-config.uint64=18446744073709551615",
				"-config.url=service-1",
				"-config.regexp=[:digit:]",
				"-config.duration=1s",
				"-config.string.pointer=foo",
				"-config.bool.pointer=true",
				"-config.float32.pointer=3.1415",
				"-config.float64.pointer=3.14159265359",
				"-config.int.pointer=9223372036854775807",
				"-config.int8.pointer=127",
				"-config.int16.pointer=32767",
				"-config.int32.pointer=2147483647",
				"-config.int64.pointer=9223372036854775807",
				"-config.uint.pointer=18446744073709551615",
				"-config.uint8.pointer=255",
				"-config.uint16.pointer=65535",
				"-config.uint32.pointer=4294967295",
				"-config.uint64.pointer=18446744073709551615",
				"-config.url.pointer=service-1",
				"-config.regexp.pointer=[:digit:]",
				"-config.duration.pointer=1s",
				"-config.string.slice=foo|bar",
				"-config.bool.slice=false|true",
				"-config.float32.slice=3.1415|2.7182",
				"-config.float64.slice=3.14159265359|2.71828182845",
				"-config.int.slice=-9223372036854775808|9223372036854775807",
				"-config.int8.slice=-128|127",
				"-config.int16.slice=-32768|32767",
				"-config.int32.slice=-2147483648|2147483647",
				"-config.int64.slice=-9223372036854775808|9223372036854775807",
				"-config.uint.slice=0|18446744073709551615",
				"-config.uint8.slice=0|255",
				"-config.uint16.slice=0|65535",
				"-config.uint32.slice=0|4294967295",
				"-config.uint64.slice=0|18446744073709551615",
				"-config.url.slice=service-1|service-2",
				"-config.regexp.slice=[:digit:]|[:alpha:]",
				"-config.duration.slice=1s|1m",
			},
			[]env{},
			[]file{},
			&config{},
			[]Option{
				ListSep("|"),
				PrefixFlag("config."),
			},
			nil,
			&cfg,
		},
		{
			"AllFromEnvVars",
			[]string{"app"},
			[]env{
				{"STRING", "foo"},
				{"BOOL", "true"},
				{"FLOAT32", "3.1415"},
				{"FLOAT64", "3.14159265359"},
				{"INT", "9223372036854775807"},
				{"INT8", "127"},
				{"INT16", "32767"},
				{"INT32", "2147483647"},
				{"INT64", "9223372036854775807"},
				{"UINT", "18446744073709551615"},
				{"UINT8", "255"},
				{"UINT16", "65535"},
				{"UINT32", "4294967295"},
				{"UINT64", "18446744073709551615"},
				{"URL", "service-1"},
				{"REGEXP", "[:digit:]"},
				{"DURATION", "1s"},
				{"STRING_POINTER", "foo"},
				{"BOOL_POINTER", "true"},
				{"FLOAT32_POINTER", "3.1415"},
				{"FLOAT64_POINTER", "3.14159265359"},
				{"INT_POINTER", "9223372036854775807"},
				{"INT8_POINTER", "127"},
				{"INT16_POINTER", "32767"},
				{"INT32_POINTER", "2147483647"},
				{"INT64_POINTER", "9223372036854775807"},
				{"UINT_POINTER", "18446744073709551615"},
				{"UINT8_POINTER", "255"},
				{"UINT16_POINTER", "65535"},
				{"UINT32_POINTER", "4294967295"},
				{"UINT64_POINTER", "18446744073709551615"},
				{"URL_POINTER", "service-1"},
				{"REGEXP_POINTER", "[:digit:]"},
				{"DURATION_POINTER", "1s"},
				{"STRING_SLICE", "foo,bar"},
				{"BOOL_SLICE", "false,true"},
				{"FLOAT32_SLICE", "3.1415,2.7182"},
				{"FLOAT64_SLICE", "3.14159265359,2.71828182845"},
				{"INT_SLICE", "-9223372036854775808,9223372036854775807"},
				{"INT8_SLICE", "-128,127"},
				{"INT16_SLICE", "-32768,32767"},
				{"INT32_SLICE", "-2147483648,2147483647"},
				{"INT64_SLICE", "-9223372036854775808,9223372036854775807"},
				{"UINT_SLICE", "0,18446744073709551615"},
				{"UINT8_SLICE", "0,255"},
				{"UINT16_SLICE", "0,65535"},
				{"UINT32_SLICE", "0,4294967295"},
				{"UINT64_SLICE", "0,18446744073709551615"},
				{"URL_SLICE", "service-1,service-2"},
				{"REGEXP_SLICE", "[:digit:],[:alpha:]"},
				{"DURATION_SLICE", "1s,1m"},
			},
			[]file{},
			&config{},
			nil,
			nil,
			&cfg,
		},
		{
			"AllFromEnvVarsWithOptions",
			[]string{"app"},
			[]env{
				{"CONFIG_STRING", "foo"},
				{"CONFIG_BOOL", "true"},
				{"CONFIG_FLOAT32", "3.1415"},
				{"CONFIG_FLOAT64", "3.14159265359"},
				{"CONFIG_INT", "9223372036854775807"},
				{"CONFIG_INT8", "127"},
				{"CONFIG_INT16", "32767"},
				{"CONFIG_INT32", "2147483647"},
				{"CONFIG_INT64", "9223372036854775807"},
				{"CONFIG_UINT", "18446744073709551615"},
				{"CONFIG_UINT8", "255"},
				{"CONFIG_UINT16", "65535"},
				{"CONFIG_UINT32", "4294967295"},
				{"CONFIG_UINT64", "18446744073709551615"},
				{"CONFIG_URL", "service-1"},
				{"CONFIG_REGEXP", "[:digit:]"},
				{"CONFIG_DURATION", "1s"},
				{"CONFIG_STRING_POINTER", "foo"},
				{"CONFIG_BOOL_POINTER", "true"},
				{"CONFIG_FLOAT32_POINTER", "3.1415"},
				{"CONFIG_FLOAT64_POINTER", "3.14159265359"},
				{"CONFIG_INT_POINTER", "9223372036854775807"},
				{"CONFIG_INT8_POINTER", "127"},
				{"CONFIG_INT16_POINTER", "32767"},
				{"CONFIG_INT32_POINTER", "2147483647"},
				{"CONFIG_INT64_POINTER", "9223372036854775807"},
				{"CONFIG_UINT_POINTER", "18446744073709551615"},
				{"CONFIG_UINT8_POINTER", "255"},
				{"CONFIG_UINT16_POINTER", "65535"},
				{"CONFIG_UINT32_POINTER", "4294967295"},
				{"CONFIG_UINT64_POINTER", "18446744073709551615"},
				{"CONFIG_URL_POINTER", "service-1"},
				{"CONFIG_REGEXP_POINTER", "[:digit:]"},
				{"CONFIG_DURATION_POINTER", "1s"},
				{"CONFIG_STRING_SLICE", "foo|bar"},
				{"CONFIG_BOOL_SLICE", "false|true"},
				{"CONFIG_FLOAT32_SLICE", "3.1415|2.7182"},
				{"CONFIG_FLOAT64_SLICE", "3.14159265359|2.71828182845"},
				{"CONFIG_INT_SLICE", "-9223372036854775808|9223372036854775807"},
				{"CONFIG_INT8_SLICE", "-128|127"},
				{"CONFIG_INT16_SLICE", "-32768|32767"},
				{"CONFIG_INT32_SLICE", "-2147483648|2147483647"},
				{"CONFIG_INT64_SLICE", "-9223372036854775808|9223372036854775807"},
				{"CONFIG_UINT_SLICE", "0|18446744073709551615"},
				{"CONFIG_UINT8_SLICE", "0|255"},
				{"CONFIG_UINT16_SLICE", "0|65535"},
				{"CONFIG_UINT32_SLICE", "0|4294967295"},
				{"CONFIG_UINT64_SLICE", "0|18446744073709551615"},
				{"CONFIG_URL_SLICE", "service-1|service-2"},
				{"CONFIG_REGEXP_SLICE", "[:digit:]|[:alpha:]"},
				{"CONFIG_DURATION_SLICE", "1s|1m"},
			},
			[]file{},
			&config{},
			[]Option{
				ListSep("|"),
				PrefixEnv("CONFIG_"),
			},
			nil,
			&cfg,
		},
		{
			"AllFromFromFiles",
			[]string{"app"},
			[]env{},
			[]file{
				{"STRING_FILE", "foo"},
				{"BOOL_FILE", "true"},
				{"FLOAT32_FILE", "3.1415"},
				{"FLOAT64_FILE", "3.14159265359"},
				{"INT_FILE", "9223372036854775807"},
				{"INT8_FILE", "127"},
				{"INT16_FILE", "32767"},
				{"INT32_FILE", "2147483647"},
				{"INT64_FILE", "9223372036854775807"},
				{"UINT_FILE", "18446744073709551615"},
				{"UINT8_FILE", "255"},
				{"UINT16_FILE", "65535"},
				{"UINT32_FILE", "4294967295"},
				{"UINT64_FILE", "18446744073709551615"},
				{"URL_FILE", "service-1"},
				{"REGEXP_FILE", "[:digit:]"},
				{"DURATION_FILE", "1s"},
				{"STRING_POINTER_FILE", "foo"},
				{"BOOL_POINTER_FILE", "true"},
				{"FLOAT32_POINTER_FILE", "3.1415"},
				{"FLOAT64_POINTER_FILE", "3.14159265359"},
				{"INT_POINTER_FILE", "9223372036854775807"},
				{"INT8_POINTER_FILE", "127"},
				{"INT16_POINTER_FILE", "32767"},
				{"INT32_POINTER_FILE", "2147483647"},
				{"INT64_POINTER_FILE", "9223372036854775807"},
				{"UINT_POINTER_FILE", "18446744073709551615"},
				{"UINT8_POINTER_FILE", "255"},
				{"UINT16_POINTER_FILE", "65535"},
				{"UINT32_POINTER_FILE", "4294967295"},
				{"UINT64_POINTER_FILE", "18446744073709551615"},
				{"URL_POINTER_FILE", "service-1"},
				{"REGEXP_POINTER_FILE", "[:digit:]"},
				{"DURATION_POINTER_FILE", "1s"},
				{"STRING_SLICE_FILE", "foo,bar"},
				{"BOOL_SLICE_FILE", "false,true"},
				{"FLOAT32_SLICE_FILE", "3.1415,2.7182"},
				{"FLOAT64_SLICE_FILE", "3.14159265359,2.71828182845"},
				{"INT_SLICE_FILE", "-9223372036854775808,9223372036854775807"},
				{"INT8_SLICE_FILE", "-128,127"},
				{"INT16_SLICE_FILE", "-32768,32767"},
				{"INT32_SLICE_FILE", "-2147483648,2147483647"},
				{"INT64_SLICE_FILE", "-9223372036854775808,9223372036854775807"},
				{"UINT_SLICE_FILE", "0,18446744073709551615"},
				{"UINT8_SLICE_FILE", "0,255"},
				{"UINT16_SLICE_FILE", "0,65535"},
				{"UINT32_SLICE_FILE", "0,4294967295"},
				{"UINT64_SLICE_FILE", "0,18446744073709551615"},
				{"URL_SLICE_FILE", "service-1,service-2"},
				{"REGEXP_SLICE_FILE", "[:digit:],[:alpha:]"},
				{"DURATION_SLICE_FILE", "1s,1m"},
			},
			&config{},
			nil,
			nil,
			&cfg,
		},
		{
			"AllFromFromFilesWithOptions",
			[]string{"app"},
			[]env{},
			[]file{
				{"CONFIG_STRING_FILE", "foo"},
				{"CONFIG_BOOL_FILE", "true"},
				{"CONFIG_FLOAT32_FILE", "3.1415"},
				{"CONFIG_FLOAT64_FILE", "3.14159265359"},
				{"CONFIG_INT_FILE", "9223372036854775807"},
				{"CONFIG_INT8_FILE", "127"},
				{"CONFIG_INT16_FILE", "32767"},
				{"CONFIG_INT32_FILE", "2147483647"},
				{"CONFIG_INT64_FILE", "9223372036854775807"},
				{"CONFIG_UINT_FILE", "18446744073709551615"},
				{"CONFIG_UINT8_FILE", "255"},
				{"CONFIG_UINT16_FILE", "65535"},
				{"CONFIG_UINT32_FILE", "4294967295"},
				{"CONFIG_UINT64_FILE", "18446744073709551615"},
				{"CONFIG_URL_FILE", "service-1"},
				{"CONFIG_REGEXP_FILE", "[:digit:]"},
				{"CONFIG_DURATION_FILE", "1s"},
				{"CONFIG_STRING_POINTER_FILE", "foo"},
				{"CONFIG_BOOL_POINTER_FILE", "true"},
				{"CONFIG_FLOAT32_POINTER_FILE", "3.1415"},
				{"CONFIG_FLOAT64_POINTER_FILE", "3.14159265359"},
				{"CONFIG_INT_POINTER_FILE", "9223372036854775807"},
				{"CONFIG_INT8_POINTER_FILE", "127"},
				{"CONFIG_INT16_POINTER_FILE", "32767"},
				{"CONFIG_INT32_POINTER_FILE", "2147483647"},
				{"CONFIG_INT64_POINTER_FILE", "9223372036854775807"},
				{"CONFIG_UINT_POINTER_FILE", "18446744073709551615"},
				{"CONFIG_UINT8_POINTER_FILE", "255"},
				{"CONFIG_UINT16_POINTER_FILE", "65535"},
				{"CONFIG_UINT32_POINTER_FILE", "4294967295"},
				{"CONFIG_UINT64_POINTER_FILE", "18446744073709551615"},
				{"CONFIG_URL_POINTER_FILE", "service-1"},
				{"CONFIG_REGEXP_POINTER_FILE", "[:digit:]"},
				{"CONFIG_DURATION_POINTER_FILE", "1s"},
				{"CONFIG_STRING_SLICE_FILE", "foo|bar"},
				{"CONFIG_BOOL_SLICE_FILE", "false|true"},
				{"CONFIG_FLOAT32_SLICE_FILE", "3.1415|2.7182"},
				{"CONFIG_FLOAT64_SLICE_FILE", "3.14159265359|2.71828182845"},
				{"CONFIG_INT_SLICE_FILE", "-9223372036854775808|9223372036854775807"},
				{"CONFIG_INT8_SLICE_FILE", "-128|127"},
				{"CONFIG_INT16_SLICE_FILE", "-32768|32767"},
				{"CONFIG_INT32_SLICE_FILE", "-2147483648|2147483647"},
				{"CONFIG_INT64_SLICE_FILE", "-9223372036854775808|9223372036854775807"},
				{"CONFIG_UINT_SLICE_FILE", "0|18446744073709551615"},
				{"CONFIG_UINT8_SLICE_FILE", "0|255"},
				{"CONFIG_UINT16_SLICE_FILE", "0|65535"},
				{"CONFIG_UINT32_SLICE_FILE", "0|4294967295"},
				{"CONFIG_UINT64_SLICE_FILE", "0|18446744073709551615"},
				{"CONFIG_URL_SLICE_FILE", "service-1|service-2"},
				{"CONFIG_REGEXP_SLICE_FILE", "[:digit:]|[:alpha:]"},
				{"CONFIG_DURATION_SLICE_FILE", "1s|1m"},
			},
			&config{},
			[]Option{
				ListSep("|"),
				PrefixFileEnv("CONFIG_"),
			},
			nil,
			&cfg,
		},
		{
			"WithTelepresenceOption",
			[]string{"app"},
			[]env{},
			[]file{
				{"STRING_FILE", "foo"},
				{"BOOL_FILE", "true"},
				{"FLOAT32_FILE", "3.1415"},
				{"FLOAT64_FILE", "3.14159265359"},
				{"INT_FILE", "9223372036854775807"},
				{"INT8_FILE", "127"},
				{"INT16_FILE", "32767"},
				{"INT32_FILE", "2147483647"},
				{"INT64_FILE", "9223372036854775807"},
				{"UINT_FILE", "18446744073709551615"},
				{"UINT8_FILE", "255"},
				{"UINT16_FILE", "65535"},
				{"UINT32_FILE", "4294967295"},
				{"UINT64_FILE", "18446744073709551615"},
				{"URL_FILE", "service-1"},
				{"REGEXP_FILE", "[:digit:]"},
				{"DURATION_FILE", "1s"},
				{"STRING_POINTER_FILE", "foo"},
				{"BOOL_POINTER_FILE", "true"},
				{"FLOAT32_POINTER_FILE", "3.1415"},
				{"FLOAT64_POINTER_FILE", "3.14159265359"},
				{"INT_POINTER_FILE", "9223372036854775807"},
				{"INT8_POINTER_FILE", "127"},
				{"INT16_POINTER_FILE", "32767"},
				{"INT32_POINTER_FILE", "2147483647"},
				{"INT64_POINTER_FILE", "9223372036854775807"},
				{"UINT_POINTER_FILE", "18446744073709551615"},
				{"UINT8_POINTER_FILE", "255"},
				{"UINT16_POINTER_FILE", "65535"},
				{"UINT32_POINTER_FILE", "4294967295"},
				{"UINT64_POINTER_FILE", "18446744073709551615"},
				{"URL_POINTER_FILE", "service-1"},
				{"REGEXP_POINTER_FILE", "[:digit:]"},
				{"DURATION_POINTER_FILE", "1s"},
				{"STRING_SLICE_FILE", "foo,bar"},
				{"BOOL_SLICE_FILE", "false,true"},
				{"FLOAT32_SLICE_FILE", "3.1415,2.7182"},
				{"FLOAT64_SLICE_FILE", "3.14159265359,2.71828182845"},
				{"INT_SLICE_FILE", "-9223372036854775808,9223372036854775807"},
				{"INT8_SLICE_FILE", "-128,127"},
				{"INT16_SLICE_FILE", "-32768,32767"},
				{"INT32_SLICE_FILE", "-2147483648,2147483647"},
				{"INT64_SLICE_FILE", "-9223372036854775808,9223372036854775807"},
				{"UINT_SLICE_FILE", "0,18446744073709551615"},
				{"UINT8_SLICE_FILE", "0,255"},
				{"UINT16_SLICE_FILE", "0,65535"},
				{"UINT32_SLICE_FILE", "0,4294967295"},
				{"UINT64_SLICE_FILE", "0,18446744073709551615"},
				{"URL_SLICE_FILE", "service-1,service-2"},
				{"REGEXP_SLICE_FILE", "[:digit:],[:alpha:]"},
				{"DURATION_SLICE_FILE", "1s,1m"},
			},
			&config{},
			[]Option{
				Telepresence(),
			},
			nil,
			&cfg,
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &reader{}
			for _, opt := range tc.opts {
				opt(c)
			}

			// Set arguments for flags
			os.Args = tc.args

			// Set environment variables
			for _, e := range tc.envs {
				err := os.Setenv(e.varName, e.value)
				assert.NoError(t, err)
				defer os.Unsetenv(e.varName)
			}

			// Testing Telepresence option
			if c.telepresence {
				err := os.Setenv(envTelepresenceRoot, "/")
				assert.NoError(t, err)
				defer os.Unsetenv(envTelepresenceRoot)
			}

			// Write configuration files
			for _, f := range tc.files {
				tmpfile, err := ioutil.TempFile("", "gotest_")
				assert.NoError(t, err)
				defer os.Remove(tmpfile.Name())

				_, err = tmpfile.WriteString(f.value)
				assert.NoError(t, err)

				err = tmpfile.Close()
				assert.NoError(t, err)

				err = os.Setenv(f.varName, tmpfile.Name())
				assert.NoError(t, err)
				defer os.Unsetenv(f.varName)
			}

			err := Pick(tc.config, tc.opts...)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, tc.config)
			}
		})
	}

	// flag.Parse() can be called only once
	flag.Parse()
}

func TestWatch(t *testing.T) {
	updateDelay := 50 * time.Millisecond

	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName   string
		initValue string
		newValue  string
	}

	files := []file{
		{"STRING_FILE", "foo", "bar"},
		{"BOOL_FILE", "false", "true"},
		{"FLOAT32_FILE", "2.7182", "3.1415"},
		{"FLOAT64_FILE", "2.7182818284", "3.14159265359"},
		{"INT_FILE", "-9223372036854775808", "9223372036854775807"},
		{"INT8_FILE", "-128", "127"},
		{"INT16_FILE", "-32768", "32767"},
		{"INT32_FILE", "-2147483648", "2147483647"},
		{"INT64_FILE", "-9223372036854775808", "9223372036854775807"},
		{"UINT_FILE", "0", "18446744073709551615"},
		{"UINT8_FILE", "0", "255"},
		{"UINT16_FILE", "0", "65535"},
		{"UINT32_FILE", "0", "4294967295"},
		{"UINT64_FILE", "0", "18446744073709551615"},
		{"URL_FILE", "service-1", "service-2"},
		{"REGEXP_FILE", "[:digit:]", "[:alpha:]"},
		{"DURATION_FILE", "1s", "1m"},
		{"STRING_SLICE_FILE", "foo,bar", "bar,foo"},
		{"BOOL_SLICE_FILE", "false,true", "true,false"},
		{"FLOAT32_SLICE_FILE", "2.7182,3.1415", "3.1415,2.7182"},
		{"FLOAT64_SLICE_FILE", "2.71828182845,3.14159265359", "3.14159265359,2.71828182845"},
		{"INT_SLICE_FILE", "-9223372036854775808,9223372036854775807", "9223372036854775807,-9223372036854775808"},
		{"INT8_SLICE_FILE", "-128,127", "127,-128"},
		{"INT16_SLICE_FILE", "-32768,32767", "32767,-32768"},
		{"INT32_SLICE_FILE", "-2147483648,2147483647", "2147483647,-2147483648"},
		{"INT64_SLICE_FILE", "-9223372036854775808,9223372036854775807", "9223372036854775807,-9223372036854775808"},
		{"UINT_SLICE_FILE", "0,18446744073709551615", "18446744073709551615,0"},
		{"UINT8_SLICE_FILE", "0,255", "255,0"},
		{"UINT16_SLICE_FILE", "0,65535", "65535,0"},
		{"UINT32_SLICE_FILE", "0,4294967295", "4294967295,0"},
		{"UINT64_SLICE_FILE", "0,18446744073709551615", "18446744073709551615,0"},
		{"URL_SLICE_FILE", "service-1,service-2", "service-2,service-1"},
		{"REGEXP_SLICE_FILE", "[:digit:],[:alpha:]", "[:alpha:],[:digit:]"},
		{"DURATION_SLICE_FILE", "1s,1m", "1m,1s"},
	}

	old := config{
		String:        "foo",
		Bool:          false,
		Float32:       2.7182,
		Float64:       2.7182818284,
		Int:           -9223372036854775808,
		Int8:          -128,
		Int16:         -32768,
		Int32:         -2147483648,
		Int64:         -9223372036854775808,
		Uint:          0,
		Uint8:         0,
		Uint16:        0,
		Uint32:        0,
		Uint64:        0,
		URL:           *url1,
		Regexp:        *re1,
		Duration:      time.Second,
		StringSlice:   []string{"foo", "bar"},
		BoolSlice:     []bool{false, true},
		Float32Slice:  []float32{2.7182, 3.1415},
		Float64Slice:  []float64{2.71828182845, 3.14159265359},
		IntSlice:      []int{-9223372036854775808, 9223372036854775807},
		Int8Slice:     []int8{-128, 127},
		Int16Slice:    []int16{-32768, 32767},
		Int32Slice:    []int32{-2147483648, 2147483647},
		Int64Slice:    []int64{-9223372036854775808, 9223372036854775807},
		UintSlice:     []uint{0, 18446744073709551615},
		Uint8Slice:    []uint8{0, 255},
		Uint16Slice:   []uint16{0, 65535},
		Uint32Slice:   []uint32{0, 4294967295},
		Uint64Slice:   []uint64{0, 18446744073709551615},
		URLSlice:      []url.URL{*url1, *url2},
		RegexpSlice:   []regexp.Regexp{*re1, *re2},
		DurationSlice: []time.Duration{time.Second, time.Minute},
	}

	new := config{
		String:        "bar",
		Bool:          true,
		Float32:       3.1415,
		Float64:       3.14159265359,
		Int:           9223372036854775807,
		Int8:          127,
		Int16:         32767,
		Int32:         2147483647,
		Int64:         9223372036854775807,
		Uint:          18446744073709551615,
		Uint8:         255,
		Uint16:        65535,
		Uint32:        4294967295,
		Uint64:        18446744073709551615,
		URL:           *url2,
		Regexp:        *re2,
		Duration:      time.Minute,
		StringSlice:   []string{"bar", "foo"},
		BoolSlice:     []bool{true, false},
		Float32Slice:  []float32{3.1415, 2.7182},
		Float64Slice:  []float64{3.14159265359, 2.71828182845},
		IntSlice:      []int{9223372036854775807, -9223372036854775808},
		Int8Slice:     []int8{127, -128},
		Int16Slice:    []int16{32767, -32768},
		Int32Slice:    []int32{2147483647, -2147483648},
		Int64Slice:    []int64{9223372036854775807, -9223372036854775808},
		UintSlice:     []uint{18446744073709551615, 0},
		Uint8Slice:    []uint8{255, 0},
		Uint16Slice:   []uint16{65535, 0},
		Uint32Slice:   []uint32{4294967295, 0},
		Uint64Slice:   []uint64{18446744073709551615, 0},
		URLSlice:      []url.URL{*url2, *url1},
		RegexpSlice:   []regexp.Regexp{*re2, *re1},
		DurationSlice: []time.Duration{time.Minute, time.Second},
	}

	updates := []Update{
		{"String", "foo"},
		{"Bool", false},
		{"Float32", float32(2.7182)},
		{"Float64", float64(2.7182818284)},
		{"Int", int(-9223372036854775808)},
		{"Int8", int8(-128)},
		{"Int16", int16(-32768)},
		{"Int32", int32(-2147483648)},
		{"Int64", int64(-9223372036854775808)},
		{"Uint", uint(0)},
		{"Uint8", uint8(0)},
		{"Uint16", uint16(0)},
		{"Uint32", uint32(0)},
		{"Uint64", uint64(0)},
		{"URL", *url1},
		{"Regexp", *re1},
		{"Duration", time.Second},
		{"StringSlice", []string{"foo", "bar"}},
		{"BoolSlice", []bool{false, true}},
		{"Float32Slice", []float32{2.7182, 3.1415}},
		{"Float64Slice", []float64{2.71828182845, 3.14159265359}},
		{"IntSlice", []int{-9223372036854775808, 9223372036854775807}},
		{"Int8Slice", []int8{-128, 127}},
		{"Int16Slice", []int16{-32768, 32767}},
		{"Int32Slice", []int32{-2147483648, 2147483647}},
		{"Int64Slice", []int64{-9223372036854775808, 9223372036854775807}},
		{"UintSlice", []uint{0, 18446744073709551615}},
		{"Uint8Slice", []uint8{0, 255}},
		{"Uint16Slice", []uint16{0, 65535}},
		{"Uint32Slice", []uint32{0, 4294967295}},
		{"Uint64Slice", []uint64{0, 18446744073709551615}},
		{"URLSlice", []url.URL{*url1, *url2}},
		{"RegexpSlice", []regexp.Regexp{*re1, *re2}},
		{"DurationSlice", []time.Duration{time.Second, time.Minute}},

		{"String", "bar"},
		{"Bool", true},
		{"Float32", float32(3.1415)},
		{"Float64", float64(3.14159265359)},
		{"Int", int(9223372036854775807)},
		{"Int8", int8(127)},
		{"Int16", int16(32767)},
		{"Int32", int32(2147483647)},
		{"Int64", int64(9223372036854775807)},
		{"Uint", uint(18446744073709551615)},
		{"Uint8", uint8(255)},
		{"Uint16", uint16(65535)},
		{"Uint32", uint32(4294967295)},
		{"Uint64", uint64(18446744073709551615)},
		{"URL", *url2},
		{"Regexp", *re2},
		{"Duration", time.Minute},
		{"StringSlice", []string{"bar", "foo"}},
		{"BoolSlice", []bool{true, false}},
		{"Float32Slice", []float32{3.1415, 2.7182}},
		{"Float64Slice", []float64{3.14159265359, 2.71828182845}},
		{"IntSlice", []int{9223372036854775807, -9223372036854775808}},
		{"Int8Slice", []int8{127, -128}},
		{"Int16Slice", []int16{32767, -32768}},
		{"Int32Slice", []int32{2147483647, -2147483648}},
		{"Int64Slice", []int64{9223372036854775807, -9223372036854775808}},
		{"UintSlice", []uint{18446744073709551615, 0}},
		{"Uint8Slice", []uint8{255, 0}},
		{"Uint16Slice", []uint16{65535, 0}},
		{"Uint32Slice", []uint32{4294967295, 0}},
		{"Uint64Slice", []uint64{18446744073709551615, 0}},
		{"URLSlice", []url.URL{*url2, *url1}},
		{"RegexpSlice", []regexp.Regexp{*re2, *re1}},
		{"DurationSlice", []time.Duration{time.Minute, time.Second}},
	}

	tests := []struct {
		name               string
		args               []string
		envs               []env
		files              []file
		config             *config
		subscribers        []chan Update
		opts               []Option
		expectedError      error
		expectedInitConfig *config
		expectedNewConfig  *config
		expectedUpdates    []Update
	}{
		{
			"BlockingChannels",
			[]string{"app"},
			[]env{},
			files,
			&config{},
			[]chan Update{
				make(chan Update),
				make(chan Update),
			},
			[]Option{},
			nil,
			&old,
			&new,
			updates,
		},
		{
			"BufferedChannels",
			[]string{"app"},
			[]env{},
			files,
			&config{},
			[]chan Update{
				make(chan Update, 100),
				make(chan Update, 100),
			},
			[]Option{},
			nil,
			&old,
			&new,
			updates,
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup

			// Set arguments for flags
			os.Args = tc.args

			// Set environment variables
			for _, e := range tc.envs {
				err := os.Setenv(e.varName, e.value)
				defer os.Unsetenv(e.varName)
				assert.NoError(t, err)
			}

			// Testing Telepresence option
			r := &reader{}
			for _, opt := range tc.opts {
				opt(r)
			}
			if r.telepresence {
				err := os.Setenv(envTelepresenceRoot, "/")
				defer os.Unsetenv(envTelepresenceRoot)
				assert.NoError(t, err)
			}

			// Write configuration files
			for _, f := range tc.files {
				tmpfile, err := ioutil.TempFile("", "gotest_")
				assert.NoError(t, err)
				defer os.Remove(tmpfile.Name())

				_, err = tmpfile.WriteString(f.initValue)
				assert.NoError(t, err)

				err = tmpfile.Close()
				assert.NoError(t, err)

				err = os.Setenv(f.varName, tmpfile.Name())
				assert.NoError(t, err)
				defer os.Unsetenv(f.varName)

				// Will write the new value to the file
				wg.Add(1)
				newValue := f.newValue
				time.AfterFunc(updateDelay, func() {
					err := ioutil.WriteFile(tmpfile.Name(), []byte(newValue), 0644)
					assert.NoError(t, err)
					wg.Done()
				})
			}

			// Listening for updates
			for i, sub := range tc.subscribers {
				go func(id int, ch chan Update) {
					for update := range ch {
						assert.Contains(t, tc.expectedUpdates, update)
					}
				}(i, sub)
			}

			close, err := Watch(tc.config, tc.subscribers, tc.opts...)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, close)
			} else {
				assert.NoError(t, err)
				defer close()

				tc.config.Lock()
				// assert.Equal(t, tc.expectedInitConfig, tc.config)
				assert.True(t, tc.config.Equal(tc.expectedInitConfig))
				tc.config.Unlock()

				// Wait for all files to be updated and the new values are picked up
				wg.Wait()
				time.Sleep(100 * time.Millisecond)

				tc.config.Lock()
				assert.True(t, tc.config.Equal(tc.expectedNewConfig))
				tc.config.Unlock()
			}
		})
	}

	// flag.Parse() can be called only once
	flag.Parse()
}
