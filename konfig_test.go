package konfig

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type config struct {
	sync.Mutex
	unexported         string
	SkipFlag           string          `flag:"-"`
	SkipFlagEnv        string          `flag:"-" env:"-"`
	SkipFlagEnvFile    string          `flag:"-" env:"-" fileenv:"-"`
	FieldString        string          // `flag:"field.string" env:"FIELD_STRING" fileenv:"FIELD_STRING_FILE"`
	FieldBool          bool            // `flag:"field.bool" env:"FIELD_BOOL" fileenv:"FIELD_BOOL_FILE"`
	FieldFloat32       float32         // `flag:"field.float32" env:"FIELD_FLOAT32" fileenv:"FIELD_FLOAT32_FILE"`
	FieldFloat64       float64         // `flag:"field.float64" env:"FIELD_FLOAT64" fileenv:"FIELD_FLOAT64_FILE"`
	FieldInt           int             // `flag:"field.int" env:"FIELD_INT" fileenv:"FIELD_INT_FILE"`
	FieldInt8          int8            // `flag:"field.int8" env:"FIELD_INT8" fileenv:"FIELD_INT8_FILE"`
	FieldInt16         int16           // `flag:"field.int16" env:"FIELD_INT16" fileenv:"FIELD_INT16_FILE"`
	FieldInt32         int32           // `flag:"field.int32" env:"FIELD_INT32" fileenv:"FIELD_INT32_FILE"`
	FieldInt64         int64           // `flag:"field.int64" env:"FIELD_INT64" fileenv:"FIELD_INT64_FILE"`
	FieldUint          uint            // `flag:"field.uint" env:"FIELD_UINT" fileenv:"FIELD_UINT_FILE"`
	FieldUint8         uint8           // `flag:"field.uint8" env:"FIELD_UINT8" fileenv:"FIELD_UINT8_FILE"`
	FieldUint16        uint16          // `flag:"field.uint16" env:"FIELD_UINT16" fileenv:"FIELD_UINT16_FILE"`
	FieldUint32        uint32          // `flag:"field.uint32" env:"FIELD_UINT32" fileenv:"FIELD_UINT32_FILE"`
	FieldUint64        uint64          // `flag:"field.uint64" env:"FIELD_UINT64" fileenv:"FIELD_UINT64_FILE"`
	FieldDuration      time.Duration   // `flag:"field.duration" env:"FIELD_DURATION" fileenv:"FIELD_DURATION_FILE"`
	FieldURL           url.URL         // `flag:"field.url" env:"FIELD_URL" fileenv:"FIELD_URL_FILE"`
	FieldStringArray   []string        // `flag:"field.string.array" env:"FIELD_STRING_ARRAY" fileenv:"FIELD_STRING_ARRAY_FILE" sep:","`
	FieldBoolArray     []bool          // `flag:"field.bool.array" env:"FIELD_BOOL_ARRAY" fileenv:"FIELD_BOOL_ARRAY_FILE" sep:","`
	FieldFloat32Array  []float32       // `flag:"field.float32.array" env:"FIELD_FLOAT32_ARRAY" fileenv:"FIELD_FLOAT32_ARRAY_FILE" sep:","`
	FieldFloat64Array  []float64       // `flag:"field.float64.array" env:"FIELD_FLOAT64_ARRAY" fileenv:"FIELD_FLOAT64_ARRAY_FILE" sep:","`
	FieldIntArray      []int           // `flag:"field.int.array" env:"FIELD_INT_ARRAY" fileenv:"FIELD_INT_ARRAY_FILE" sep:","`
	FieldInt8Array     []int8          // `flag:"field.int8.array" env:"FIELD_INT8_ARRAY" fileenv:"FIELD_INT8_ARRAY_FILE" sep:","`
	FieldInt16Array    []int16         // `flag:"field.int16.array" env:"FIELD_INT16_ARRAY" fileenv:"FIELD_INT16_ARRAY_FILE" sep:","`
	FieldInt32Array    []int32         // `flag:"field.int32.array" env:"FIELD_INT32_ARRAY" fileenv:"FIELD_INT32_ARRAY_FILE" sep:","`
	FieldInt64Array    []int64         // `flag:"field.int64.array" env:"FIELD_INT64_ARRAY" fileenv:"FIELD_INT64_ARRAY_FILE" sep:","`
	FieldUintArray     []uint          // `flag:"field.uint.array" env:"FIELD_UINT_ARRAY" fileenv:"FIELD_UINT_ARRAY_FILE" sep:","`
	FieldUint8Array    []uint8         // `flag:"field.uint8.array" env:"FIELD_UINT8_ARRAY" fileenv:"FIELD_UINT8_ARRAY_FILE" sep:","`
	FieldUint16Array   []uint16        // `flag:"field.uint16.array" env:"FIELD_UINT16_ARRAY" fileenv:"FIELD_UINT16_ARRAY_FILE" sep:","`
	FieldUint32Array   []uint32        // `flag:"field.uint32.array" env:"FIELD_UINT32_ARRAY" fileenv:"FIELD_UINT32_ARRAY_FILE" sep:","`
	FieldUint64Array   []uint64        // `flag:"field.uint64.array" env:"FIELD_UINT64_ARRAY" fileenv:"FIELD_UINT64_ARRAY_FILE" sep:","`
	FieldDurationArray []time.Duration // `flag:"field.duration.array" env:"FIELD_DURATION_ARRAY" fileenv:"FIELD_DURATION_ARRAY_FILE" sep:","`
	FieldURLArray      []url.URL       // `flag:"field.url.array" env:"FIELD_URL_ARRAY" fileenv:"FIELD_URL_ARRAY_FILE" sep:","`
}

func configEqual(c1, c2 *config) bool {
	return c1.unexported == c2.unexported &&
		c1.SkipFlag == c2.SkipFlag &&
		c1.SkipFlagEnv == c2.SkipFlagEnv &&
		c1.SkipFlagEnvFile == c2.SkipFlagEnvFile &&
		c1.FieldString == c2.FieldString &&
		c1.FieldBool == c2.FieldBool &&
		c1.FieldFloat32 == c2.FieldFloat32 &&
		c1.FieldFloat64 == c2.FieldFloat64 &&
		c1.FieldInt == c2.FieldInt &&
		c1.FieldInt8 == c2.FieldInt8 &&
		c1.FieldInt16 == c2.FieldInt16 &&
		c1.FieldInt32 == c2.FieldInt32 &&
		c1.FieldInt64 == c2.FieldInt64 &&
		c1.FieldUint == c2.FieldUint &&
		c1.FieldUint8 == c2.FieldUint8 &&
		c1.FieldUint16 == c2.FieldUint16 &&
		c1.FieldUint32 == c2.FieldUint32 &&
		c1.FieldUint64 == c2.FieldUint64 &&
		c1.FieldDuration == c2.FieldDuration &&
		c1.FieldURL == c2.FieldURL &&
		reflect.DeepEqual(c1.FieldStringArray, c2.FieldStringArray) &&
		reflect.DeepEqual(c1.FieldBoolArray, c2.FieldBoolArray) &&
		reflect.DeepEqual(c1.FieldFloat32Array, c2.FieldFloat32Array) &&
		reflect.DeepEqual(c1.FieldFloat64Array, c2.FieldFloat64Array) &&
		reflect.DeepEqual(c1.FieldIntArray, c2.FieldIntArray) &&
		reflect.DeepEqual(c1.FieldInt8Array, c2.FieldInt8Array) &&
		reflect.DeepEqual(c1.FieldInt16Array, c2.FieldInt16Array) &&
		reflect.DeepEqual(c1.FieldInt32Array, c2.FieldInt32Array) &&
		reflect.DeepEqual(c1.FieldInt64Array, c2.FieldInt64Array) &&
		reflect.DeepEqual(c1.FieldUintArray, c2.FieldUintArray) &&
		reflect.DeepEqual(c1.FieldUint8Array, c2.FieldUint8Array) &&
		reflect.DeepEqual(c1.FieldUint16Array, c2.FieldUint16Array) &&
		reflect.DeepEqual(c1.FieldUint32Array, c2.FieldUint32Array) &&
		reflect.DeepEqual(c1.FieldUint64Array, c2.FieldUint64Array) &&
		reflect.DeepEqual(c1.FieldDurationArray, c2.FieldDurationArray) &&
		reflect.DeepEqual(c1.FieldURLArray, c2.FieldURLArray)
}

func TestControllerFromEnv(t *testing.T) {
	tests := []struct {
		name               string
		env                map[string]string
		expectedController *controller
	}{
		{
			name: "NoOption",
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "InvalidDebug",
			env: map[string]string{
				envDebug: "NaN",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "OutOfRangeDebug",
			env: map[string]string{
				envDebug: "999",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "DebugLevel1",
			env: map[string]string{
				envDebug: "1",
			},
			expectedController: &controller{
				debug:         1,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "DebugLevel2",
			env: map[string]string{
				envDebug: "2",
			},
			expectedController: &controller{
				debug:         2,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "DebugLevel3",
			env: map[string]string{
				envDebug: "3",
			},
			expectedController: &controller{
				debug:         3,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "ListSep",
			env: map[string]string{
				envListSep: "|",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       "|",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "SkipFlag",
			env: map[string]string{
				envSkipFlag: "true",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      true,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "SkipEnv",
			env: map[string]string{
				envSkipEnv: "true",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       true,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "SkipEnvFile",
			env: map[string]string{
				envSkipFileEnv: "true",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   true,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "PrefixFlag",
			env: map[string]string{
				envPrefixFlag: "config.",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "config.",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "PrefixEnv",
			env: map[string]string{
				envPrefixEnv: "CONFIG_",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "CONFIG_",
				prefixFileEnv: "",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "PrefixEnv",
			env: map[string]string{
				envPrefixFileEnv: "CONFIG_",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "CONFIG_",
				telepresence:  false,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "Telepresence",
			env: map[string]string{
				envTelepresence: "true",
			},
			expectedController: &controller{
				debug:         0,
				listSep:       ",",
				skipFlag:      false,
				skipEnv:       false,
				skipFileEnv:   false,
				prefixFlag:    "",
				prefixEnv:     "",
				prefixFileEnv: "",
				telepresence:  true,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
		{
			name: "AllOptions",
			env: map[string]string{
				envDebug:         "3",
				envListSep:       "|",
				envSkipFlag:      "true",
				envSkipEnv:       "true",
				envSkipFileEnv:   "true",
				envPrefixFlag:    "config.",
				envPrefixEnv:     "CONFIG_",
				envPrefixFileEnv: "CONFIG_",
				envTelepresence:  "true",
			},
			expectedController: &controller{
				debug:         3,
				listSep:       "|",
				skipFlag:      true,
				skipEnv:       true,
				skipFileEnv:   true,
				prefixFlag:    "config.",
				prefixEnv:     "CONFIG_",
				prefixFileEnv: "CONFIG_",
				telepresence:  true,
				subscribers:   nil,
				filesToFields: map[string]fieldInfo{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for name, value := range tc.env {
				err := os.Setenv(name, value)
				assert.NoError(t, err)
				defer os.Unsetenv(name)
			}

			c := controllerFromEnv()
			assert.Equal(t, tc.expectedController, c)
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		c         *controller
		verbosity uint
		expected  *controller
	}{
		{
			&controller{},
			2,
			&controller{
				debug: 2,
			},
		},
	}

	for _, tc := range tests {
		opt := Debug(tc.verbosity)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestListSep(t *testing.T) {
	tests := []struct {
		c        *controller
		sep      string
		expected *controller
	}{
		{
			&controller{},
			"|",
			&controller{
				listSep: "|",
			},
		},
	}

	for _, tc := range tests {
		opt := ListSep(tc.sep)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestSkipFlag(t *testing.T) {
	tests := []struct {
		c        *controller
		expected *controller
	}{
		{
			&controller{},
			&controller{
				skipFlag: true,
			},
		},
	}

	for _, tc := range tests {
		opt := SkipFlag()
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestSkipEnv(t *testing.T) {
	tests := []struct {
		c        *controller
		expected *controller
	}{
		{
			&controller{},
			&controller{
				skipEnv: true,
			},
		},
	}

	for _, tc := range tests {
		opt := SkipEnv()
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestSkipFileEnv(t *testing.T) {
	tests := []struct {
		c        *controller
		expected *controller
	}{
		{
			&controller{},
			&controller{
				skipFileEnv: true,
			},
		},
	}

	for _, tc := range tests {
		opt := SkipFileEnv()
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestPrefixFlag(t *testing.T) {
	tests := []struct {
		c        *controller
		prefix   string
		expected *controller
	}{
		{
			&controller{},
			"config.",
			&controller{
				prefixFlag: "config.",
			},
		},
	}

	for _, tc := range tests {
		opt := PrefixFlag(tc.prefix)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestPrefixEnv(t *testing.T) {
	tests := []struct {
		c        *controller
		prefix   string
		expected *controller
	}{
		{
			&controller{},
			"CONFIG_",
			&controller{
				prefixEnv: "CONFIG_",
			},
		},
	}

	for _, tc := range tests {
		opt := PrefixEnv(tc.prefix)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestPrefixFileEnv(t *testing.T) {
	tests := []struct {
		c        *controller
		prefix   string
		expected *controller
	}{
		{
			&controller{},
			"CONFIG_",
			&controller{
				prefixFileEnv: "CONFIG_",
			},
		},
	}

	for _, tc := range tests {
		opt := PrefixFileEnv(tc.prefix)
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestTelepresence(t *testing.T) {
	tests := []struct {
		c        *controller
		expected *controller
	}{
		{
			&controller{},
			&controller{
				telepresence: true,
			},
		},
	}

	for _, tc := range tests {
		opt := Telepresence()
		opt(tc.c)

		assert.Equal(t, tc.expected, tc.c)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		expectedString string
	}{
		{
			"NoOption",
			&controller{},
			"",
		},
		{
			"WithDebug",
			&controller{
				debug: 2,
			},
			"Debug<2>",
		},
		{
			"WithListSep",
			&controller{
				listSep: "|",
			},
			"ListSep<|>",
		},
		{
			"WithPrefixFlag",
			&controller{
				prefixFlag: "config.",
			},
			"PrefixFlag<config.>",
		},
		{
			"WithPrefixEnv",
			&controller{
				prefixEnv: "CONFIG_",
			},
			"PrefixEnv<CONFIG_>",
		},
		{
			"WithprefixFileEnv",
			&controller{
				prefixFileEnv: "CONFIG_",
			},
			"PrefixFileEnv<CONFIG_>",
		},
		{
			"WithSkipFlag",
			&controller{
				skipFlag: true,
			},
			"SkipFlag",
		},
		{
			"WithSkipEnv",
			&controller{
				skipEnv: true,
			},
			"SkipEnv",
		},
		{
			"WithSkipFileEnv",
			&controller{
				skipFileEnv: true,
			},
			"SkipFileEnv",
		},
		{
			"WithTelepresence",
			&controller{
				telepresence: true,
			},
			"Telepresence",
		},
		{
			"WithSubscribers",
			&controller{
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"Subscribers<2>",
		},
		{
			"WithAll",
			&controller{
				debug:         2,
				listSep:       "|",
				prefixFlag:    "config.",
				prefixEnv:     "CONFIG_",
				prefixFileEnv: "CONFIG_",
				skipFlag:      true,
				skipEnv:       true,
				skipFileEnv:   true,
				telepresence:  true,
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"Debug<2> + ListSep<|> + SkipFlag + SkipEnv + SkipFileEnv + PrefixFlag<config.> + PrefixEnv<CONFIG_> + PrefixFileEnv<CONFIG_> + Telepresence + Subscribers<2>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.c.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name string
		c    *controller
		v    uint
		msg  string
		args []interface{}
	}{
		{
			"WithoutDebug",
			&controller{},
			1,
			"testing ...",
			nil,
		},
		{
			"WithDebug",
			&controller{
				debug: 2,
			},
			2,
			"testing ...",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.c.log(tc.v, tc.msg, tc.args...)
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName string
		value   string
	}

	tests := []struct {
		name                                      string
		args                                      []string
		envConfig                                 env
		fileConfig                                file
		fieldName, flagName, envName, fileEnvName string
		c                                         *controller
		expectedValue                             string
		expectFilePath                            bool
	}{
		{
			"SkipFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"info",
			false,
		},
		{
			"SkipFlagAndEnv",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "LOG_LEVEL_FILE",
			&controller{},
			"error",
			true,
		},
		{
			"SkipFlagAndEnvAndFile",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "-",
			&controller{},
			"",
			false,
		},
		{
			"SkipAllFlags",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{
				skipFlag: true,
			},
			"info",
			false,
		},
		{
			"SkipAllFlagsAndEnvs",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{
				skipFlag: true,
				skipEnv:  true,
			},
			"error",
			true,
		},
		{
			"SkipAllFlagsAndEnvsAndFileEnvs",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{
				skipFlag:    true,
				skipEnv:     true,
				skipFileEnv: true,
			},
			"",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"debug",
			false,
		},
		{
			"FromEnvVar",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"info",
			false,
		},
		{
			"FromFile",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{},
			"error",
			true,
		},
		{
			"FromFileWithTelepresenceOption",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "info"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&controller{telepresence: true},
			"info",
			true,
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set value using a flag
			os.Args = tc.args

			// Set value in an environment variable
			err := os.Setenv(tc.envConfig.varName, tc.envConfig.value)
			assert.NoError(t, err)
			defer os.Unsetenv(tc.envConfig.varName)

			// Testing Telepresence option
			if tc.c.telepresence {
				err := os.Setenv(envTelepresenceRoot, "/")
				assert.NoError(t, err)
				defer os.Unsetenv(envTelepresenceRoot)
			}

			// Write value in a temporary config file

			tmpfile, err := ioutil.TempFile("", "gotest_")
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())

			_, err = tmpfile.WriteString(tc.fileConfig.value)
			assert.NoError(t, err)

			err = tmpfile.Close()
			assert.NoError(t, err)

			err = os.Setenv(tc.fileConfig.varName, tmpfile.Name())
			assert.NoError(t, err)
			defer os.Unsetenv(tc.fileConfig.varName)

			// Verify
			value, filePath := tc.c.getFieldValue(tc.fieldName, tc.flagName, tc.envName, tc.fileEnvName)
			assert.Equal(t, tc.expectedValue, value)
			if tc.expectFilePath {
				assert.Equal(t, tmpfile.Name(), filePath)
			}
		})
	}
}

func TestNotifySubscribers(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		fieldName      string
		fieldValue     interface{}
		expectedUpdate Update
	}{
		{
			"Nil",
			&controller{},
			"FieldBool", true,
			Update{},
		},
		{
			"NoChannel",
			&controller{
				subscribers: []chan Update{},
			},
			"FieldString", "value",
			Update{},
		},
		{
			"WithBlockingChannels",
			&controller{
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"FieldInt", 27,
			Update{"FieldInt", 27},
		},
		{
			"WithBufferedChannels",
			&controller{
				subscribers: []chan Update{
					make(chan Update, 1),
					make(chan Update, 1),
				},
			},
			"FieldFloat", 3.1415,
			Update{"FieldFloat", 3.1415},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.c.notifySubscribers(tc.fieldName, tc.fieldValue)

			if tc.expectedUpdate != (Update{}) {
				for _, ch := range tc.c.subscribers {
					update := <-ch
					assert.Equal(t, tc.expectedUpdate, update)
				}
			}
		})
	}
}

func TestSetString(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          string
		fieldName      string
		fieldValue     string
		expectedValue  string
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          "",
			fieldName:      "Field",
			fieldValue:     "test",
			expectedValue:  "test",
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          "test",
			fieldName:      "Field",
			fieldValue:     "test",
			expectedValue:  "test",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setString(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetBool(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          bool
		fieldName      string
		fieldValue     string
		expectedValue  bool
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          false,
			fieldName:      "Field",
			fieldValue:     "true",
			expectedValue:  true,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          true,
			fieldName:      "Field",
			fieldValue:     "true",
			expectedValue:  true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setBool(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat32(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          float32
		fieldName      string
		fieldValue     string
		expectedValue  float32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "3.1415",
			expectedValue:  3.1415,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          3.1415,
			fieldName:      "Field",
			fieldValue:     "3.1415",
			expectedValue:  3.1415,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat32(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat64(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          float64
		fieldName      string
		fieldValue     string
		expectedValue  float64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "3.14159265359",
			expectedValue:  3.14159265359,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          3.14159265359,
			fieldName:      "Field",
			fieldValue:     "3.14159265359",
			expectedValue:  3.14159265359,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat64(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int
		fieldName      string
		fieldValue     string
		expectedValue  int
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-2147483648",
			expectedValue:  -2147483648,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -2147483648,
			fieldName:      "Field",
			fieldValue:     "-2147483648",
			expectedValue:  -2147483648,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt8(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int8
		fieldName      string
		fieldValue     string
		expectedValue  int8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-128",
			expectedValue:  -128,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -128,
			fieldName:      "Field",
			fieldValue:     "-128",
			expectedValue:  -128,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt8(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt16(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int16
		fieldName      string
		fieldValue     string
		expectedValue  int16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-32768",
			expectedValue:  -32768,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -32768,
			fieldName:      "Field",
			fieldValue:     "-32768",
			expectedValue:  -32768,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt16(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt32(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int32
		fieldName      string
		fieldValue     string
		expectedValue  int32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-2147483648",
			expectedValue:  -2147483648,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -2147483648,
			fieldName:      "Field",
			fieldValue:     "-2147483648",
			expectedValue:  -2147483648,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt32(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt64(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          int64
		fieldName      string
		fieldValue     string
		expectedValue  int64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "-9223372036854775808",
			expectedValue:  -9223372036854775808,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          -9223372036854775808,
			fieldName:      "Field",
			fieldValue:     "-9223372036854775808",
			expectedValue:  -9223372036854775808,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt64(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetDuration(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          time.Duration
		fieldName      string
		fieldValue     string
		expectedValue  time.Duration
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "1h0m0s",
			expectedValue:  time.Hour,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          time.Hour,
			fieldName:      "Field",
			fieldValue:     "1h0m0s",
			expectedValue:  time.Hour,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt64(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint
		fieldName      string
		fieldValue     string
		expectedValue  uint
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "4294967295",
			expectedValue:  4294967295,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          4294967295,
			fieldName:      "Field",
			fieldValue:     "4294967295",
			expectedValue:  4294967295,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint8(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint8
		fieldName      string
		fieldValue     string
		expectedValue  uint8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "255",
			expectedValue:  255,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          255,
			fieldName:      "Field",
			fieldValue:     "255",
			expectedValue:  255,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint8(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint16(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint16
		fieldName      string
		fieldValue     string
		expectedValue  uint16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "65535",
			expectedValue:  65535,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          65535,
			fieldName:      "Field",
			fieldValue:     "65535",
			expectedValue:  65535,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint16(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint32(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint32
		fieldName      string
		fieldValue     string
		expectedValue  uint32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "4294967295",
			expectedValue:  4294967295,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          4294967295,
			fieldName:      "Field",
			fieldValue:     "4294967295",
			expectedValue:  4294967295,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint32(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint64(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          uint64
		fieldName      string
		fieldValue     string
		expectedValue  uint64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          0,
			fieldName:      "Field",
			fieldValue:     "18446744073709551615",
			expectedValue:  18446744073709551615,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          18446744073709551615,
			fieldName:      "Field",
			fieldValue:     "18446744073709551615",
			expectedValue:  18446744073709551615,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint64(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetURL(t *testing.T) {
	u, _ := url.Parse("example.com")

	tests := []struct {
		name           string
		c              *controller
		field          url.URL
		fieldName      string
		fieldValue     string
		expectedValue  url.URL
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          url.URL{},
			fieldName:      "Field",
			fieldValue:     "example.com",
			expectedValue:  *u,
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          *u,
			fieldName:      "Field",
			fieldValue:     "example.com",
			expectedValue:  *u,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setStruct(v, tc.fieldName, tc.fieldValue)

			assert.Equal(t, tc.expectedValue, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetStringSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []string
		fieldName      string
		fieldValues    []string
		expectedValues []string
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []string{},
			fieldName:      "Field",
			fieldValues:    []string{"milad", "mona"},
			expectedValues: []string{"milad", "mona"},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []string{"milad", "mona"},
			fieldName:      "Field",
			fieldValues:    []string{"milad", "mona"},
			expectedValues: []string{"milad", "mona"},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setStringSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetBoolSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []bool
		fieldName      string
		fieldValues    []string
		expectedValues []bool
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []bool{},
			fieldName:      "Field",
			fieldValues:    []string{"false", "true"},
			expectedValues: []bool{false, true},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []bool{false, true},
			fieldName:      "Field",
			fieldValues:    []string{"false", "true"},
			expectedValues: []bool{false, true},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setBoolSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat32Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []float32
		fieldName      string
		fieldValues    []string
		expectedValues []float32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []float32{},
			fieldName:      "Field",
			fieldValues:    []string{"3.1415", "2.7182"},
			expectedValues: []float32{3.1415, 2.7182},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []float32{3.1415, 2.7182},
			fieldName:      "Field",
			fieldValues:    []string{"3.1415", "2.7182"},
			expectedValues: []float32{3.1415, 2.7182},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFloat64Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []float64
		fieldName      string
		fieldValues    []string
		expectedValues []float64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []float64{},
			fieldName:      "Field",
			fieldValues:    []string{"3.14159265", "2.71828182"},
			expectedValues: []float64{3.14159265, 2.71828182},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []float64{3.14159265, 2.71828182},
			fieldName:      "Field",
			fieldValues:    []string{"3.14159265", "2.71828182"},
			expectedValues: []float64{3.14159265, 2.71828182},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setFloat64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetIntSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int
		fieldName      string
		fieldValues    []string
		expectedValues []int
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int{},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []int{27, 69},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int{27, 69},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []int{27, 69},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setIntSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt8Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int8
		fieldName      string
		fieldValues    []string
		expectedValues []int8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int8{},
			fieldName:      "Field",
			fieldValues:    []string{"-128", "127"},
			expectedValues: []int8{-128, 127},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int8{-128, 127},
			fieldName:      "Field",
			fieldValues:    []string{"-128", "127"},
			expectedValues: []int8{-128, 127},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt8Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt16Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int16
		fieldName      string
		fieldValues    []string
		expectedValues []int16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int16{},
			fieldName:      "Field",
			fieldValues:    []string{"-32768", "32767"},
			expectedValues: []int16{-32768, 32767},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int16{-32768, 32767},
			fieldName:      "Field",
			fieldValues:    []string{"-32768", "32767"},
			expectedValues: []int16{-32768, 32767},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt16Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt32Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int32
		fieldName      string
		fieldValues    []string
		expectedValues []int32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int32{},
			fieldName:      "Field",
			fieldValues:    []string{"-2147483648", "2147483647"},
			expectedValues: []int32{-2147483648, 2147483647},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int32{-2147483648, 2147483647},
			fieldName:      "Field",
			fieldValues:    []string{"-2147483648", "2147483647"},
			expectedValues: []int32{-2147483648, 2147483647},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetInt64Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []int64
		fieldName      string
		fieldValues    []string
		expectedValues []int64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []int64{},
			fieldName:      "Field",
			fieldValues:    []string{"-9223372036854775808", "9223372036854775807"},
			expectedValues: []int64{-9223372036854775808, 9223372036854775807},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []int64{-9223372036854775808, 9223372036854775807},
			fieldName:      "Field",
			fieldValues:    []string{"-9223372036854775808", "9223372036854775807"},
			expectedValues: []int64{-9223372036854775808, 9223372036854775807},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetDurationSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []time.Duration
		fieldName      string
		fieldValues    []string
		expectedValues []time.Duration
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []time.Duration{},
			fieldName:      "Field",
			fieldValues:    []string{"1h0m0s", "1m0s"},
			expectedValues: []time.Duration{time.Hour, time.Minute},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []time.Duration{time.Hour, time.Minute},
			fieldName:      "Field",
			fieldValues:    []string{"1h0m0s", "1m0s"},
			expectedValues: []time.Duration{time.Hour, time.Minute},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setInt64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUintSlice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint
		fieldName      string
		fieldValues    []string
		expectedValues []uint
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint{},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []uint{27, 69},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint{27, 69},
			fieldName:      "Field",
			fieldValues:    []string{"27", "69"},
			expectedValues: []uint{27, 69},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUintSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint8Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint8
		fieldName      string
		fieldValues    []string
		expectedValues []uint8
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint8{},
			fieldName:      "Field",
			fieldValues:    []string{"128", "255"},
			expectedValues: []uint8{128, 255},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint8{128, 255},
			fieldName:      "Field",
			fieldValues:    []string{"128", "255"},
			expectedValues: []uint8{128, 255},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint8Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint16Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint16
		fieldName      string
		fieldValues    []string
		expectedValues []uint16
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint16{},
			fieldName:      "Field",
			fieldValues:    []string{"32768", "65535"},
			expectedValues: []uint16{32768, 65535},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint16{32768, 65535},
			fieldName:      "Field",
			fieldValues:    []string{"32768", "65535"},
			expectedValues: []uint16{32768, 65535},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint16Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint32Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint32
		fieldName      string
		fieldValues    []string
		expectedValues []uint32
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint32{},
			fieldName:      "Field",
			fieldValues:    []string{"2147483648", "4294967295"},
			expectedValues: []uint32{2147483648, 4294967295},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint32{2147483648, 4294967295},
			fieldName:      "Field",
			fieldValues:    []string{"2147483648", "4294967295"},
			expectedValues: []uint32{2147483648, 4294967295},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint32Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetUint64Slice(t *testing.T) {
	tests := []struct {
		name           string
		c              *controller
		field          []uint64
		fieldName      string
		fieldValues    []string
		expectedValues []uint64
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []uint64{},
			fieldName:      "Field",
			fieldValues:    []string{"9223372036854775808", "18446744073709551615"},
			expectedValues: []uint64{9223372036854775808, 18446744073709551615},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []uint64{9223372036854775808, 18446744073709551615},
			fieldName:      "Field",
			fieldValues:    []string{"9223372036854775808", "18446744073709551615"},
			expectedValues: []uint64{9223372036854775808, 18446744073709551615},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setUint64Slice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetURLSlice(t *testing.T) {
	u1, _ := url.Parse("localhost")
	u2, _ := url.Parse("example.com")

	tests := []struct {
		name           string
		c              *controller
		field          []url.URL
		fieldName      string
		fieldValues    []string
		expectedValues []url.URL
		expectedResult bool
	}{
		{
			name:           "NewValue",
			c:              &controller{},
			field:          []url.URL{},
			fieldName:      "Field",
			fieldValues:    []string{"localhost", "example.com"},
			expectedValues: []url.URL{*u1, *u2},
			expectedResult: true,
		},
		{
			name:           "NoNewValue",
			c:              &controller{},
			field:          []url.URL{*u1, *u2},
			fieldName:      "Field",
			fieldValues:    []string{"localhost", "example.com"},
			expectedValues: []url.URL{*u1, *u2},
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := reflect.ValueOf(&tc.field).Elem()
			res := tc.c.setURLSlice(v, tc.fieldName, tc.fieldValues)

			assert.Equal(t, tc.expectedValues, tc.field)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	d4h := 4 * time.Hour
	d8h := 8 * time.Hour

	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")
	service3URL, _ := url.Parse("service-3:8080")
	service4URL, _ := url.Parse("service-4:8080")

	tests := []struct {
		name           string
		c              *controller
		config         *config
		values         map[string]string
		expectedResult bool
		expectedConfig *config
	}{
		{
			"NewValues",
			&controller{},
			&config{
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			map[string]string{
				"FieldString":        "content",
				"FieldBool":          "true",
				"FieldFloat32":       "2.7182",
				"FieldFloat64":       "2.7182818284",
				"FieldInt":           "2147483647",
				"FieldInt8":          "127",
				"FieldInt16":         "32767",
				"FieldInt32":         "2147483647",
				"FieldInt64":         "9223372036854775807",
				"FieldUint":          "2147483648",
				"FieldUint8":         "128",
				"FieldUint16":        "32768",
				"FieldUint32":        "2147483648",
				"FieldUint64":        "9223372036854775808",
				"FieldDuration":      "4h",
				"FieldURL":           "service-3:8080",
				"FieldStringArray":   "mona,milad",
				"FieldBoolArray":     "true,false",
				"FieldFloat32Array":  "2.7182,3.1415",
				"FieldFloat64Array":  "2.71828182845,3.14159265359",
				"FieldIntArray":      "2147483647,-2147483648",
				"FieldInt8Array":     "127,-128",
				"FieldInt16Array":    "32767,-32768",
				"FieldInt32Array":    "2147483647,-2147483648",
				"FieldInt64Array":    "9223372036854775807,-9223372036854775808",
				"FieldUintArray":     "4294967295,0",
				"FieldUint8Array":    "255,0",
				"FieldUint16Array":   "65535,0",
				"FieldUint32Array":   "4294967295,0",
				"FieldUint64Array":   "18446744073709551615,0",
				"FieldDurationArray": "4h,8h",
				"FieldURLArray":      "service-3:8080,service-4:8080",
			},
			true,
			&config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       2.7182,
				FieldFloat64:       2.7182818284,
				FieldInt:           2147483647,
				FieldInt8:          127,
				FieldInt16:         32767,
				FieldInt32:         2147483647,
				FieldInt64:         9223372036854775807,
				FieldUint:          2147483648,
				FieldUint8:         128,
				FieldUint16:        32768,
				FieldUint32:        2147483648,
				FieldUint64:        9223372036854775808,
				FieldDuration:      d4h,
				FieldURL:           *service3URL,
				FieldStringArray:   []string{"mona", "milad"},
				FieldBoolArray:     []bool{true, false},
				FieldFloat32Array:  []float32{2.7182, 3.1415},
				FieldFloat64Array:  []float64{2.71828182845, 3.14159265359},
				FieldIntArray:      []int{2147483647, -2147483648},
				FieldInt8Array:     []int8{127, -128},
				FieldInt16Array:    []int16{32767, -32768},
				FieldInt32Array:    []int32{2147483647, -2147483648},
				FieldInt64Array:    []int64{9223372036854775807, -9223372036854775808},
				FieldUintArray:     []uint{4294967295, 0},
				FieldUint8Array:    []uint8{255, 0},
				FieldUint16Array:   []uint16{65535, 0},
				FieldUint32Array:   []uint32{4294967295, 0},
				FieldUint64Array:   []uint64{18446744073709551615, 0},
				FieldDurationArray: []time.Duration{d4h, d8h},
				FieldURLArray:      []url.URL{*service3URL, *service4URL},
			},
		},
		{
			"NoNewValues",
			&controller{},
			&config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       2.7182,
				FieldFloat64:       2.7182818284,
				FieldInt:           2147483647,
				FieldInt8:          127,
				FieldInt16:         32767,
				FieldInt32:         2147483647,
				FieldInt64:         9223372036854775807,
				FieldUint:          2147483648,
				FieldUint8:         128,
				FieldUint16:        32768,
				FieldUint32:        2147483648,
				FieldUint64:        9223372036854775808,
				FieldDuration:      d4h,
				FieldURL:           *service3URL,
				FieldStringArray:   []string{"mona", "milad"},
				FieldBoolArray:     []bool{true, false},
				FieldFloat32Array:  []float32{2.7182, 3.1415},
				FieldFloat64Array:  []float64{2.71828182845, 3.14159265359},
				FieldIntArray:      []int{2147483647, -2147483648},
				FieldInt8Array:     []int8{127, -128},
				FieldInt16Array:    []int16{32767, -32768},
				FieldInt32Array:    []int32{2147483647, -2147483648},
				FieldInt64Array:    []int64{9223372036854775807, -9223372036854775808},
				FieldUintArray:     []uint{4294967295, 0},
				FieldUint8Array:    []uint8{255, 0},
				FieldUint16Array:   []uint16{65535, 0},
				FieldUint32Array:   []uint32{4294967295, 0},
				FieldUint64Array:   []uint64{18446744073709551615, 0},
				FieldDurationArray: []time.Duration{d4h, d8h},
				FieldURLArray:      []url.URL{*service3URL, *service4URL},
			},
			map[string]string{
				"FieldString":        "content",
				"FieldBool":          "true",
				"FieldFloat32":       "2.7182",
				"FieldFloat64":       "2.7182818284",
				"FieldInt":           "2147483647",
				"FieldInt8":          "127",
				"FieldInt16":         "32767",
				"FieldInt32":         "2147483647",
				"FieldInt64":         "9223372036854775807",
				"FieldUint":          "2147483648",
				"FieldUint8":         "128",
				"FieldUint16":        "32768",
				"FieldUint32":        "2147483648",
				"FieldUint64":        "9223372036854775808",
				"FieldDuration":      "4h",
				"FieldURL":           "service-3:8080",
				"FieldStringArray":   "mona,milad",
				"FieldBoolArray":     "true,false",
				"FieldFloat32Array":  "2.7182,3.1415",
				"FieldFloat64Array":  "2.71828182845,3.14159265359",
				"FieldIntArray":      "2147483647,-2147483648",
				"FieldInt8Array":     "127,-128",
				"FieldInt16Array":    "32767,-32768",
				"FieldInt32Array":    "2147483647,-2147483648",
				"FieldInt64Array":    "9223372036854775807,-9223372036854775808",
				"FieldUintArray":     "4294967295,0",
				"FieldUint8Array":    "255,0",
				"FieldUint16Array":   "65535,0",
				"FieldUint32Array":   "4294967295,0",
				"FieldUint64Array":   "18446744073709551615,0",
				"FieldDurationArray": "4h,8h",
				"FieldURLArray":      "service-3:8080,service-4:8080",
			},
			false,
			&config{
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       2.7182,
				FieldFloat64:       2.7182818284,
				FieldInt:           2147483647,
				FieldInt8:          127,
				FieldInt16:         32767,
				FieldInt32:         2147483647,
				FieldInt64:         9223372036854775807,
				FieldUint:          2147483648,
				FieldUint8:         128,
				FieldUint16:        32768,
				FieldUint32:        2147483648,
				FieldUint64:        9223372036854775808,
				FieldDuration:      d4h,
				FieldURL:           *service3URL,
				FieldStringArray:   []string{"mona", "milad"},
				FieldBoolArray:     []bool{true, false},
				FieldFloat32Array:  []float32{2.7182, 3.1415},
				FieldFloat64Array:  []float64{2.71828182845, 3.14159265359},
				FieldIntArray:      []int{2147483647, -2147483648},
				FieldInt8Array:     []int8{127, -128},
				FieldInt16Array:    []int16{32767, -32768},
				FieldInt32Array:    []int32{2147483647, -2147483648},
				FieldInt64Array:    []int64{9223372036854775807, -9223372036854775808},
				FieldUintArray:     []uint{4294967295, 0},
				FieldUint8Array:    []uint8{255, 0},
				FieldUint16Array:   []uint16{65535, 0},
				FieldUint32Array:   []uint32{4294967295, 0},
				FieldUint64Array:   []uint64{18446744073709551615, 0},
				FieldDurationArray: []time.Duration{d4h, d8h},
				FieldURLArray:      []url.URL{*service3URL, *service4URL},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct := reflect.ValueOf(tc.config).Elem()
			for i := 0; i < vStruct.NumField(); i++ {
				v := vStruct.Field(i)
				f := vStruct.Type().Field(i)

				// Only consider exported and supported fields that their names start with "Field"
				if v.CanSet() && isTypeSupported(v.Type()) && strings.HasPrefix(f.Name, "Field") {
					f := fieldInfo{
						value:   v,
						name:    f.Name,
						listSep: ",",
					}

					res := tc.c.setFieldValue(f, tc.values[f.name])
					assert.Equal(t, tc.expectedResult, res)
				}
			}

			assert.True(t, configEqual(tc.expectedConfig, tc.config))
		})
	}
}

func TestIterateOnFields(t *testing.T) {
	tests := []struct {
		name                 string
		c                    *controller
		config               interface{}
		expectedValues       []reflect.Value
		expectedFieldNames   []string
		expectedFlagNames    []string
		expectedEnvNames     []string
		expectedFileEnvNames []string
		expectedListSeps     []string
	}{
		{
			name: "Default",
			c: &controller{
				listSep: ",",
			},
			config:         &config{},
			expectedValues: []reflect.Value{},
			expectedFieldNames: []string{
				"SkipFlag", "SkipFlagEnv", "SkipFlagEnvFile",
				"FieldString",
				"FieldBool",
				"FieldFloat32", "FieldFloat64",
				"FieldInt", "FieldInt8", "FieldInt16", "FieldInt32", "FieldInt64",
				"FieldUint", "FieldUint8", "FieldUint16", "FieldUint32", "FieldUint64",
				"FieldDuration", "FieldURL",
				"FieldStringArray",
				"FieldBoolArray",
				"FieldFloat32Array", "FieldFloat64Array",
				"FieldIntArray", "FieldInt8Array", "FieldInt16Array", "FieldInt32Array", "FieldInt64Array",
				"FieldUintArray", "FieldUint8Array", "FieldUint16Array", "FieldUint32Array", "FieldUint64Array",
				"FieldDurationArray", "FieldURLArray",
			},
			expectedFlagNames: []string{
				"-", "-", "-",
				"field.string",
				"field.bool",
				"field.float32", "field.float64",
				"field.int", "field.int8", "field.int16", "field.int32", "field.int64",
				"field.uint", "field.uint8", "field.uint16", "field.uint32", "field.uint64",
				"field.duration", "field.url",
				"field.string.array",
				"field.bool.array",
				"field.float32.array", "field.float64.array",
				"field.int.array", "field.int8.array", "field.int16.array", "field.int32.array", "field.int64.array",
				"field.uint.array", "field.uint8.array", "field.uint16.array", "field.uint32.array", "field.uint64.array",
				"field.duration.array", "field.url.array",
			},
			expectedEnvNames: []string{
				"SKIP_FLAG", "-", "-",
				"FIELD_STRING",
				"FIELD_BOOL",
				"FIELD_FLOAT32", "FIELD_FLOAT64",
				"FIELD_INT", "FIELD_INT8", "FIELD_INT16", "FIELD_INT32", "FIELD_INT64",
				"FIELD_UINT", "FIELD_UINT8", "FIELD_UINT16", "FIELD_UINT32", "FIELD_UINT64",
				"FIELD_DURATION", "FIELD_URL",
				"FIELD_STRING_ARRAY",
				"FIELD_BOOL_ARRAY",
				"FIELD_FLOAT32_ARRAY", "FIELD_FLOAT64_ARRAY",
				"FIELD_INT_ARRAY", "FIELD_INT8_ARRAY", "FIELD_INT16_ARRAY", "FIELD_INT32_ARRAY", "FIELD_INT64_ARRAY",
				"FIELD_UINT_ARRAY", "FIELD_UINT8_ARRAY", "FIELD_UINT16_ARRAY", "FIELD_UINT32_ARRAY", "FIELD_UINT64_ARRAY",
				"FIELD_DURATION_ARRAY", "FIELD_URL_ARRAY",
			},
			expectedFileEnvNames: []string{
				"SKIP_FLAG_FILE", "SKIP_FLAG_ENV_FILE", "-",
				"FIELD_STRING_FILE",
				"FIELD_BOOL_FILE",
				"FIELD_FLOAT32_FILE", "FIELD_FLOAT64_FILE",
				"FIELD_INT_FILE", "FIELD_INT8_FILE", "FIELD_INT16_FILE", "FIELD_INT32_FILE", "FIELD_INT64_FILE",
				"FIELD_UINT_FILE", "FIELD_UINT8_FILE", "FIELD_UINT16_FILE", "FIELD_UINT32_FILE", "FIELD_UINT64_FILE",
				"FIELD_DURATION_FILE", "FIELD_URL_FILE",
				"FIELD_STRING_ARRAY_FILE",
				"FIELD_BOOL_ARRAY_FILE",
				"FIELD_FLOAT32_ARRAY_FILE", "FIELD_FLOAT64_ARRAY_FILE",
				"FIELD_INT_ARRAY_FILE", "FIELD_INT8_ARRAY_FILE", "FIELD_INT16_ARRAY_FILE", "FIELD_INT32_ARRAY_FILE", "FIELD_INT64_ARRAY_FILE",
				"FIELD_UINT_ARRAY_FILE", "FIELD_UINT8_ARRAY_FILE", "FIELD_UINT16_ARRAY_FILE", "FIELD_UINT32_ARRAY_FILE", "FIELD_UINT64_ARRAY_FILE",
				"FIELD_DURATION_ARRAY_FILE", "FIELD_URL_ARRAY_FILE",
			},
			expectedListSeps: []string{
				",", ",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
			},
		},
		{
			name: "WithPrefixOptions",
			c: &controller{
				listSep:       ",",
				prefixFlag:    "config.",
				prefixEnv:     "CONFIG_",
				prefixFileEnv: "CONFIG_",
			},
			config:         &config{},
			expectedValues: []reflect.Value{},
			expectedFieldNames: []string{
				"SkipFlag", "SkipFlagEnv", "SkipFlagEnvFile",
				"FieldString",
				"FieldBool",
				"FieldFloat32", "FieldFloat64",
				"FieldInt", "FieldInt8", "FieldInt16", "FieldInt32", "FieldInt64",
				"FieldUint", "FieldUint8", "FieldUint16", "FieldUint32", "FieldUint64",
				"FieldDuration", "FieldURL",
				"FieldStringArray",
				"FieldBoolArray",
				"FieldFloat32Array", "FieldFloat64Array",
				"FieldIntArray", "FieldInt8Array", "FieldInt16Array", "FieldInt32Array", "FieldInt64Array",
				"FieldUintArray", "FieldUint8Array", "FieldUint16Array", "FieldUint32Array", "FieldUint64Array",
				"FieldDurationArray", "FieldURLArray",
			},
			expectedFlagNames: []string{
				"-", "-", "-",
				"config.field.string",
				"config.field.bool",
				"config.field.float32", "config.field.float64",
				"config.field.int", "config.field.int8", "config.field.int16", "config.field.int32", "config.field.int64",
				"config.field.uint", "config.field.uint8", "config.field.uint16", "config.field.uint32", "config.field.uint64",
				"config.field.duration", "config.field.url",
				"config.field.string.array",
				"config.field.bool.array",
				"config.field.float32.array", "config.field.float64.array",
				"config.field.int.array", "config.field.int8.array", "config.field.int16.array", "config.field.int32.array", "config.field.int64.array",
				"config.field.uint.array", "config.field.uint8.array", "config.field.uint16.array", "config.field.uint32.array", "config.field.uint64.array",
				"config.field.duration.array", "config.field.url.array",
			},
			expectedEnvNames: []string{
				"CONFIG_SKIP_FLAG", "-", "-",
				"CONFIG_FIELD_STRING",
				"CONFIG_FIELD_BOOL",
				"CONFIG_FIELD_FLOAT32", "CONFIG_FIELD_FLOAT64",
				"CONFIG_FIELD_INT", "CONFIG_FIELD_INT8", "CONFIG_FIELD_INT16", "CONFIG_FIELD_INT32", "CONFIG_FIELD_INT64",
				"CONFIG_FIELD_UINT", "CONFIG_FIELD_UINT8", "CONFIG_FIELD_UINT16", "CONFIG_FIELD_UINT32", "CONFIG_FIELD_UINT64",
				"CONFIG_FIELD_DURATION", "CONFIG_FIELD_URL",
				"CONFIG_FIELD_STRING_ARRAY",
				"CONFIG_FIELD_BOOL_ARRAY",
				"CONFIG_FIELD_FLOAT32_ARRAY", "CONFIG_FIELD_FLOAT64_ARRAY",
				"CONFIG_FIELD_INT_ARRAY", "CONFIG_FIELD_INT8_ARRAY", "CONFIG_FIELD_INT16_ARRAY", "CONFIG_FIELD_INT32_ARRAY", "CONFIG_FIELD_INT64_ARRAY",
				"CONFIG_FIELD_UINT_ARRAY", "CONFIG_FIELD_UINT8_ARRAY", "CONFIG_FIELD_UINT16_ARRAY", "CONFIG_FIELD_UINT32_ARRAY", "CONFIG_FIELD_UINT64_ARRAY",
				"CONFIG_FIELD_DURATION_ARRAY", "CONFIG_FIELD_URL_ARRAY",
			},
			expectedFileEnvNames: []string{
				"CONFIG_SKIP_FLAG_FILE", "CONFIG_SKIP_FLAG_ENV_FILE", "-",
				"CONFIG_FIELD_STRING_FILE",
				"CONFIG_FIELD_BOOL_FILE",
				"CONFIG_FIELD_FLOAT32_FILE", "CONFIG_FIELD_FLOAT64_FILE",
				"CONFIG_FIELD_INT_FILE", "CONFIG_FIELD_INT8_FILE", "CONFIG_FIELD_INT16_FILE", "CONFIG_FIELD_INT32_FILE", "CONFIG_FIELD_INT64_FILE",
				"CONFIG_FIELD_UINT_FILE", "CONFIG_FIELD_UINT8_FILE", "CONFIG_FIELD_UINT16_FILE", "CONFIG_FIELD_UINT32_FILE", "CONFIG_FIELD_UINT64_FILE",
				"CONFIG_FIELD_DURATION_FILE", "CONFIG_FIELD_URL_FILE",
				"CONFIG_FIELD_STRING_ARRAY_FILE",
				"CONFIG_FIELD_BOOL_ARRAY_FILE",
				"CONFIG_FIELD_FLOAT32_ARRAY_FILE", "CONFIG_FIELD_FLOAT64_ARRAY_FILE",
				"CONFIG_FIELD_INT_ARRAY_FILE", "CONFIG_FIELD_INT8_ARRAY_FILE", "CONFIG_FIELD_INT16_ARRAY_FILE", "CONFIG_FIELD_INT32_ARRAY_FILE", "CONFIG_FIELD_INT64_ARRAY_FILE",
				"CONFIG_FIELD_UINT_ARRAY_FILE", "CONFIG_FIELD_UINT8_ARRAY_FILE", "CONFIG_FIELD_UINT16_ARRAY_FILE", "CONFIG_FIELD_UINT32_ARRAY_FILE", "CONFIG_FIELD_UINT64_ARRAY_FILE",
				"CONFIG_FIELD_DURATION_ARRAY_FILE", "CONFIG_FIELD_URL_ARRAY_FILE",
			},
			expectedListSeps: []string{
				",", ",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
				",",
				",",
				",", ",",
				",", ",", ",", ",", ",",
				",", ",", ",", ",", ",",
				",", ",",
			},
		},
		{
			name: "WithListSepOption",
			c: &controller{
				listSep: "|",
			},
			config:         &config{},
			expectedValues: []reflect.Value{},
			expectedFieldNames: []string{
				"SkipFlag", "SkipFlagEnv", "SkipFlagEnvFile",
				"FieldString",
				"FieldBool",
				"FieldFloat32", "FieldFloat64",
				"FieldInt", "FieldInt8", "FieldInt16", "FieldInt32", "FieldInt64",
				"FieldUint", "FieldUint8", "FieldUint16", "FieldUint32", "FieldUint64",
				"FieldDuration", "FieldURL",
				"FieldStringArray",
				"FieldBoolArray",
				"FieldFloat32Array", "FieldFloat64Array",
				"FieldIntArray", "FieldInt8Array", "FieldInt16Array", "FieldInt32Array", "FieldInt64Array",
				"FieldUintArray", "FieldUint8Array", "FieldUint16Array", "FieldUint32Array", "FieldUint64Array",
				"FieldDurationArray", "FieldURLArray",
			},
			expectedFlagNames: []string{
				"-", "-", "-",
				"field.string",
				"field.bool",
				"field.float32", "field.float64",
				"field.int", "field.int8", "field.int16", "field.int32", "field.int64",
				"field.uint", "field.uint8", "field.uint16", "field.uint32", "field.uint64",
				"field.duration", "field.url",
				"field.string.array",
				"field.bool.array",
				"field.float32.array", "field.float64.array",
				"field.int.array", "field.int8.array", "field.int16.array", "field.int32.array", "field.int64.array",
				"field.uint.array", "field.uint8.array", "field.uint16.array", "field.uint32.array", "field.uint64.array",
				"field.duration.array", "field.url.array",
			},
			expectedEnvNames: []string{
				"SKIP_FLAG", "-", "-",
				"FIELD_STRING",
				"FIELD_BOOL",
				"FIELD_FLOAT32", "FIELD_FLOAT64",
				"FIELD_INT", "FIELD_INT8", "FIELD_INT16", "FIELD_INT32", "FIELD_INT64",
				"FIELD_UINT", "FIELD_UINT8", "FIELD_UINT16", "FIELD_UINT32", "FIELD_UINT64",
				"FIELD_DURATION", "FIELD_URL",
				"FIELD_STRING_ARRAY",
				"FIELD_BOOL_ARRAY",
				"FIELD_FLOAT32_ARRAY", "FIELD_FLOAT64_ARRAY",
				"FIELD_INT_ARRAY", "FIELD_INT8_ARRAY", "FIELD_INT16_ARRAY", "FIELD_INT32_ARRAY", "FIELD_INT64_ARRAY",
				"FIELD_UINT_ARRAY", "FIELD_UINT8_ARRAY", "FIELD_UINT16_ARRAY", "FIELD_UINT32_ARRAY", "FIELD_UINT64_ARRAY",
				"FIELD_DURATION_ARRAY", "FIELD_URL_ARRAY",
			},
			expectedFileEnvNames: []string{
				"SKIP_FLAG_FILE", "SKIP_FLAG_ENV_FILE", "-",
				"FIELD_STRING_FILE",
				"FIELD_BOOL_FILE",
				"FIELD_FLOAT32_FILE", "FIELD_FLOAT64_FILE",
				"FIELD_INT_FILE", "FIELD_INT8_FILE", "FIELD_INT16_FILE", "FIELD_INT32_FILE", "FIELD_INT64_FILE",
				"FIELD_UINT_FILE", "FIELD_UINT8_FILE", "FIELD_UINT16_FILE", "FIELD_UINT32_FILE", "FIELD_UINT64_FILE",
				"FIELD_DURATION_FILE", "FIELD_URL_FILE",
				"FIELD_STRING_ARRAY_FILE",
				"FIELD_BOOL_ARRAY_FILE",
				"FIELD_FLOAT32_ARRAY_FILE", "FIELD_FLOAT64_ARRAY_FILE",
				"FIELD_INT_ARRAY_FILE", "FIELD_INT8_ARRAY_FILE", "FIELD_INT16_ARRAY_FILE", "FIELD_INT32_ARRAY_FILE", "FIELD_INT64_ARRAY_FILE",
				"FIELD_UINT_ARRAY_FILE", "FIELD_UINT8_ARRAY_FILE", "FIELD_UINT16_ARRAY_FILE", "FIELD_UINT32_ARRAY_FILE", "FIELD_UINT64_ARRAY_FILE",
				"FIELD_DURATION_ARRAY_FILE", "FIELD_URL_ARRAY_FILE",
			},
			expectedListSeps: []string{
				"|", "|", "|",
				"|",
				"|",
				"|", "|",
				"|", "|", "|", "|", "|",
				"|", "|", "|", "|", "|",
				"|", "|",
				"|",
				"|",
				"|", "|",
				"|", "|", "|", "|", "|",
				"|", "|", "|", "|", "|",
				"|", "|",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// values := []reflect.Value{}
			fieldNames := []string{}
			flagNames := []string{}
			envNames := []string{}
			fileEnvNames := []string{}
			listSeps := []string{}

			vStruct, err := validateStruct(tc.config)
			assert.NoError(t, err)

			tc.c.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
				// values = append(values, v)
				fieldNames = append(fieldNames, fieldName)
				flagNames = append(flagNames, flagName)
				envNames = append(envNames, envName)
				fileEnvNames = append(fileEnvNames, fileEnvName)
				listSeps = append(listSeps, listSep)
			})

			// assert.Equal(t, tc.expectedValues, values)
			assert.Equal(t, tc.expectedFieldNames, fieldNames)
			assert.Equal(t, tc.expectedFlagNames, flagNames)
			assert.Equal(t, tc.expectedEnvNames, envNames)
			assert.Equal(t, tc.expectedFileEnvNames, fileEnvNames)
			assert.Equal(t, tc.expectedListSeps, listSeps)
		})
	}
}

func TestRegisterFlags(t *testing.T) {
	tests := []struct {
		name          string
		c             *controller
		config        interface{}
		expectedError error
		expectedFlags []string
	}{
		{
			name:          "Default",
			c:             &controller{},
			config:        &config{},
			expectedError: nil,
			expectedFlags: []string{
				"field.string",
				"field.bool",
				"field.float32", "field.float64",
				"field.int", "field.int8", "field.int16", "field.int32", "field.int64",
				"field.uint", "field.uint8", "field.uint16", "field.uint32", "field.uint64",
				"field.duration", "field.url",
				"field.string.array",
				"field.bool.array",
				"field.float32.array", "field.float64.array",
				"field.int.array", "field.int8.array", "field.int16.array", "field.int32.array", "field.int64.array",
				"field.uint.array", "field.uint8.array", "field.uint16.array", "field.uint32.array", "field.uint64.array",
				"field.duration.array", "field.url.array",
			},
		},
		{
			name: "WithPrefixFlagOption",
			c: &controller{
				prefixFlag: "config.",
			},
			config:        &config{},
			expectedError: nil,
			expectedFlags: []string{
				"config.field.string",
				"config.field.bool",
				"config.field.float32", "config.field.float64",
				"config.field.int", "config.field.int8", "config.field.int16", "config.field.int32", "config.field.int64",
				"config.field.uint", "config.field.uint8", "config.field.uint16", "config.field.uint32", "config.field.uint64",
				"config.field.duration", "config.field.url",
				"config.field.string.array",
				"config.field.bool.array",
				"config.field.float32.array", "config.field.float64.array",
				"config.field.int.array", "config.field.int8.array", "config.field.int16.array", "config.field.int32.array", "config.field.int64.array",
				"config.field.uint.array", "config.field.uint8.array", "config.field.uint16.array", "config.field.uint32.array", "config.field.uint64.array",
				"config.field.duration.array", "config.field.url.array",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct, err := validateStruct(tc.config)
			assert.NoError(t, err)

			tc.c.registerFlags(vStruct)

			for _, expectedFlag := range tc.expectedFlags {
				f := flag.Lookup(expectedFlag)
				assert.NotEmpty(t, f)
			}
		})
	}
}

func TestReadFields(t *testing.T) {
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName string
		value   string
	}

	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")

	tests := []struct {
		name             string
		args             []string
		envs             []env
		files            []file
		c                *controller
		config           interface{}
		expectedConfig   interface{}
		expectedFilesLen int
	}{
		{
			"Empty",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{},
			0,
		},
		{
			"AllFromDefaults",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{
				unexported:         "internal",
				SkipFlag:           "default",
				SkipFlagEnv:        "default",
				SkipFlagEnvFile:    "default",
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			&config{
				unexported:         "internal",
				SkipFlag:           "default",
				SkipFlagEnv:        "default",
				SkipFlagEnvFile:    "default",
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFlags#1",
			[]string{
				"path/to/binary",
				"-field.string", "content",
				"-field.bool",
				"-field.float32", "3.1415",
				"-field.float64", "3.14159265359",
				"-field.int", "-2147483648",
				"-field.int8", "-128",
				"-field.int16", "-32768",
				"-field.int32", "-2147483648",
				"-field.int64", "-9223372036854775808",
				"-field.uint", "4294967295",
				"-field.uint8", "255",
				"-field.uint16", "65535",
				"-field.uint32", "4294967295",
				"-field.uint64", "18446744073709551615",
				"-field.duration", "90m",
				"-field.url", "service-1:8080",
				"-field.string.array", "milad,mona",
				"-field.bool.array", "false,true",
				"-field.float32.array", "3.1415,2.7182",
				"-field.float64.array", "3.14159265359,2.71828182845",
				"-field.int.array", "-2147483648,2147483647",
				"-field.int8.array", "-128,127",
				"-field.int16.array", "-32768,32767",
				"-field.int32.array", "-2147483648,2147483647",
				"-field.int64.array", "-9223372036854775808,9223372036854775807",
				"-field.uint.array", "0,4294967295",
				"-field.uint8.array", "0,255",
				"-field.uint16.array", "0,65535",
				"-field.uint32.array", "0,4294967295",
				"-field.uint64.array", "0,18446744073709551615",
				"-field.duration.array", "90m,120m",
				"-field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFlags#2",
			[]string{
				"path/to/binary",
				"--field.string", "content",
				"--field.bool",
				"--field.float32", "3.1415",
				"--field.float64", "3.14159265359",
				"--field.int", "-2147483648",
				"--field.int8", "-128",
				"--field.int16", "-32768",
				"--field.int32", "-2147483648",
				"--field.int64", "-9223372036854775808",
				"--field.uint", "4294967295",
				"--field.uint8", "255",
				"--field.uint16", "65535",
				"--field.uint32", "4294967295",
				"--field.uint64", "18446744073709551615",
				"--field.duration", "90m",
				"--field.url", "service-1:8080",
				"--field.string.array", "milad,mona",
				"--field.bool.array", "false,true",
				"--field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
				"--field.int.array", "-2147483648,2147483647",
				"--field.int8.array", "-128,127",
				"--field.int16.array", "-32768,32767",
				"--field.int32.array", "-2147483648,2147483647",
				"--field.int64.array", "-9223372036854775808,9223372036854775807",
				"--field.uint.array", "0,4294967295",
				"--field.uint8.array", "0,255",
				"--field.uint16.array", "0,65535",
				"--field.uint32.array", "0,4294967295",
				"--field.uint64.array", "0,18446744073709551615",
				"--field.duration.array", "90m,120m",
				"--field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFlags#3",
			[]string{
				"path/to/binary",
				"-field.string=content",
				"-field.bool",
				"-field.float32=3.1415",
				"-field.float64=3.14159265359",
				"-field.int=-2147483648",
				"-field.int8=-128",
				"-field.int16=-32768",
				"-field.int32=-2147483648",
				"-field.int64=-9223372036854775808",
				"-field.uint=4294967295",
				"-field.uint8=255",
				"-field.uint16=65535",
				"-field.uint32=4294967295",
				"-field.uint64=18446744073709551615",
				"-field.duration=90m",
				"-field.url=service-1:8080",
				"-field.string.array=milad,mona",
				"-field.bool.array=false,true",
				"-field.float32.array=3.1415,2.7182",
				"-field.float64.array=3.14159265359,2.71828182845",
				"-field.int.array=-2147483648,2147483647",
				"-field.int8.array=-128,127",
				"-field.int16.array=-32768,32767",
				"-field.int32.array=-2147483648,2147483647",
				"-field.int64.array=-9223372036854775808,9223372036854775807",
				"-field.uint.array=0,4294967295",
				"-field.uint8.array=0,255",
				"-field.uint16.array=0,65535",
				"-field.uint32.array=0,4294967295",
				"-field.uint64.array=0,18446744073709551615",
				"-field.duration.array=90m,120m",
				"-field.url.array=service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFlags#4",
			[]string{
				"path/to/binary",
				"--field.string=content",
				"--field.bool",
				"--field.float32=3.1415",
				"--field.float64=3.14159265359",
				"--field.int=-2147483648",
				"--field.int8=-128",
				"--field.int16=-32768",
				"--field.int32=-2147483648",
				"--field.int64=-9223372036854775808",
				"--field.uint=4294967295",
				"--field.uint8=255",
				"--field.uint16=65535",
				"--field.uint32=4294967295",
				"--field.uint64=18446744073709551615",
				"--field.duration=90m",
				"--field.url=service-1:8080",
				"--field.string.array=milad,mona",
				"--field.bool.array=false,true",
				"--field.float32.array=3.1415,2.7182",
				"--field.float64.array=3.14159265359,2.71828182845",
				"--field.int.array=-2147483648,2147483647",
				"--field.int8.array=-128,127",
				"--field.int16.array=-32768,32767",
				"--field.int32.array=-2147483648,2147483647",
				"--field.int64.array=-9223372036854775808,9223372036854775807",
				"--field.uint.array=0,4294967295",
				"--field.uint8.array=0,255",
				"--field.uint16.array=0,65535",
				"--field.uint32.array=0,4294967295",
				"--field.uint64.array=0,18446744073709551615",
				"--field.duration.array=90m,120m",
				"--field.url.array=service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFlagsWithPrefixFlagOption",
			[]string{
				"path/to/binary",
				"-config.field.string", "content",
				"-config.field.bool",
				"-config.field.float32", "3.1415",
				"-config.field.float64", "3.14159265359",
				"-config.field.int", "-2147483648",
				"-config.field.int8", "-128",
				"-config.field.int16", "-32768",
				"-config.field.int32", "-2147483648",
				"-config.field.int64", "-9223372036854775808",
				"-config.field.uint", "4294967295",
				"-config.field.uint8", "255",
				"-config.field.uint16", "65535",
				"-config.field.uint32", "4294967295",
				"-config.field.uint64", "18446744073709551615",
				"-config.field.duration", "90m",
				"-config.field.url", "service-1:8080",
				"-config.field.string.array", "milad,mona",
				"-config.field.bool.array", "false,true",
				"-config.field.float32.array", "3.1415,2.7182",
				"-config.field.float64.array", "3.14159265359,2.71828182845",
				"-config.field.int.array", "-2147483648,2147483647",
				"-config.field.int8.array", "-128,127",
				"-config.field.int16.array", "-32768,32767",
				"-config.field.int32.array", "-2147483648,2147483647",
				"-config.field.int64.array", "-9223372036854775808,9223372036854775807",
				"-config.field.uint.array", "0,4294967295",
				"-config.field.uint8.array", "0,255",
				"-config.field.uint16.array", "0,65535",
				"-config.field.uint32.array", "0,4294967295",
				"-config.field.uint64.array", "0,18446744073709551615",
				"-config.field.duration.array", "90m,120m",
				"-config.field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&controller{
				listSep:       ",",
				prefixFlag:    "config.",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFlagsWithListSepOption",
			[]string{
				"path/to/binary",
				"-field.string", "content",
				"-field.bool",
				"-field.float32", "3.1415",
				"-field.float64", "3.14159265359",
				"-field.int", "-2147483648",
				"-field.int8", "-128",
				"-field.int16", "-32768",
				"-field.int32", "-2147483648",
				"-field.int64", "-9223372036854775808",
				"-field.uint", "4294967295",
				"-field.uint8", "255",
				"-field.uint16", "65535",
				"-field.uint32", "4294967295",
				"-field.uint64", "18446744073709551615",
				"-field.duration", "90m",
				"-field.url", "service-1:8080",
				"-field.string.array", "milad|mona",
				"-field.bool.array", "false|true",
				"-field.float32.array", "3.1415|2.7182",
				"-field.float64.array", "3.14159265359|2.71828182845",
				"-field.int.array", "-2147483648|2147483647",
				"-field.int8.array", "-128|127",
				"-field.int16.array", "-32768|32767",
				"-field.int32.array", "-2147483648|2147483647",
				"-field.int64.array", "-9223372036854775808|9223372036854775807",
				"-field.uint.array", "0|4294967295",
				"-field.uint8.array", "0|255",
				"-field.uint16.array", "0|65535",
				"-field.uint32.array", "0|4294967295",
				"-field.uint64.array", "0|18446744073709551615",
				"-field.duration.array", "90m|120m",
				"-field.url.array", "service-1:8080|service-2:8080",
			},
			[]env{},
			[]file{},
			&controller{
				listSep:       "|",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromEnvVars",
			[]string{"path/to/binary"},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_STRING", "content"},
				{"FIELD_BOOL", "true"},
				{"FIELD_FLOAT32", "3.1415"},
				{"FIELD_FLOAT64", "3.14159265359"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_UINT", "4294967295"},
				{"FIELD_UINT8", "255"},
				{"FIELD_UINT16", "65535"},
				{"FIELD_UINT32", "4294967295"},
				{"FIELD_UINT64", "18446744073709551615"},
				{"FIELD_DURATION", "90m"},
				{"FIELD_URL", "service-1:8080"},
				{"FIELD_STRING_ARRAY", "milad,mona"},
				{"FIELD_BOOL_ARRAY", "false,true"},
				{"FIELD_FLOAT32_ARRAY", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY", "0,4294967295"},
				{"FIELD_UINT8_ARRAY", "0,255"},
				{"FIELD_UINT16_ARRAY", "0,65535"},
				{"FIELD_UINT32_ARRAY", "0,4294967295"},
				{"FIELD_UINT64_ARRAY", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY", "90m,120m"},
				{"FIELD_URL_ARRAY", "service-1:8080,service-2:8080"},
			},
			[]file{},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromEnvVarsWithPrefixEnvOption",
			[]string{"path/to/binary"},
			[]env{
				{"CONFIG_SKIP_FLAG", "from_env"},
				{"CONFIG_SKIP_FLAG_ENV", "from_env"},
				{"CONFIG_SKIP_FLAG_ENV_FILE", "from_env"},
				{"CONFIG_FIELD_STRING", "content"},
				{"CONFIG_FIELD_BOOL", "true"},
				{"CONFIG_FIELD_FLOAT32", "3.1415"},
				{"CONFIG_FIELD_FLOAT64", "3.14159265359"},
				{"CONFIG_FIELD_INT", "-2147483648"},
				{"CONFIG_FIELD_INT8", "-128"},
				{"CONFIG_FIELD_INT16", "-32768"},
				{"CONFIG_FIELD_INT32", "-2147483648"},
				{"CONFIG_FIELD_INT64", "-9223372036854775808"},
				{"CONFIG_FIELD_UINT", "4294967295"},
				{"CONFIG_FIELD_UINT8", "255"},
				{"CONFIG_FIELD_UINT16", "65535"},
				{"CONFIG_FIELD_UINT32", "4294967295"},
				{"CONFIG_FIELD_UINT64", "18446744073709551615"},
				{"CONFIG_FIELD_DURATION", "90m"},
				{"CONFIG_FIELD_URL", "service-1:8080"},
				{"CONFIG_FIELD_STRING_ARRAY", "milad,mona"},
				{"CONFIG_FIELD_BOOL_ARRAY", "false,true"},
				{"CONFIG_FIELD_FLOAT32_ARRAY", "3.1415,2.7182"},
				{"CONFIG_FIELD_FLOAT64_ARRAY", "3.14159265359,2.71828182845"},
				{"CONFIG_FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT8_ARRAY", "-128,127"},
				{"CONFIG_FIELD_INT16_ARRAY", "-32768,32767"},
				{"CONFIG_FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				{"CONFIG_FIELD_UINT_ARRAY", "0,4294967295"},
				{"CONFIG_FIELD_UINT8_ARRAY", "0,255"},
				{"CONFIG_FIELD_UINT16_ARRAY", "0,65535"},
				{"CONFIG_FIELD_UINT32_ARRAY", "0,4294967295"},
				{"CONFIG_FIELD_UINT64_ARRAY", "0,18446744073709551615"},
				{"CONFIG_FIELD_DURATION_ARRAY", "90m,120m"},
				{"CONFIG_FIELD_URL_ARRAY", "service-1:8080,service-2:8080"},
			},
			[]file{},
			&controller{
				listSep:       ",",
				prefixEnv:     "CONFIG_",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromEnvVarsWithListSepOption",
			[]string{"path/to/binary"},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_STRING", "content"},
				{"FIELD_BOOL", "true"},
				{"FIELD_FLOAT32", "3.1415"},
				{"FIELD_FLOAT64", "3.14159265359"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_UINT", "4294967295"},
				{"FIELD_UINT8", "255"},
				{"FIELD_UINT16", "65535"},
				{"FIELD_UINT32", "4294967295"},
				{"FIELD_UINT64", "18446744073709551615"},
				{"FIELD_DURATION", "90m"},
				{"FIELD_URL", "service-1:8080"},
				{"FIELD_STRING_ARRAY", "milad|mona"},
				{"FIELD_BOOL_ARRAY", "false|true"},
				{"FIELD_FLOAT32_ARRAY", "3.1415|2.7182"},
				{"FIELD_FLOAT64_ARRAY", "3.14159265359|2.71828182845"},
				{"FIELD_INT_ARRAY", "-2147483648|2147483647"},
				{"FIELD_INT8_ARRAY", "-128|127"},
				{"FIELD_INT16_ARRAY", "-32768|32767"},
				{"FIELD_INT32_ARRAY", "-2147483648|2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808|9223372036854775807"},
				{"FIELD_UINT_ARRAY", "0|4294967295"},
				{"FIELD_UINT8_ARRAY", "0|255"},
				{"FIELD_UINT16_ARRAY", "0|65535"},
				{"FIELD_UINT32_ARRAY", "0|4294967295"},
				{"FIELD_UINT64_ARRAY", "0|18446744073709551615"},
				{"FIELD_DURATION_ARRAY", "90m|120m"},
				{"FIELD_URL_ARRAY", "service-1:8080|service-2:8080"},
			},
			[]file{},
			&controller{
				listSep:       "|",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"AllFromFromFiles",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			34,
		},
		{
			"AllFromFromFilesWithPrefixFileEnvOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"CONFIG_SKIP_FLAG_FILE", "from_file"},
				{"CONFIG_SKIP_FLAG_ENV_FILE", "from_file"},
				{"CONFIG_SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"CONFIG_FIELD_STRING_FILE", "content"},
				{"CONFIG_FIELD_BOOL_FILE", "true"},
				{"CONFIG_FIELD_FLOAT32_FILE", "3.1415"},
				{"CONFIG_FIELD_FLOAT64_FILE", "3.14159265359"},
				{"CONFIG_FIELD_INT_FILE", "-2147483648"},
				{"CONFIG_FIELD_INT8_FILE", "-128"},
				{"CONFIG_FIELD_INT16_FILE", "-32768"},
				{"CONFIG_FIELD_INT32_FILE", "-2147483648"},
				{"CONFIG_FIELD_INT64_FILE", "-9223372036854775808"},
				{"CONFIG_FIELD_UINT_FILE", "4294967295"},
				{"CONFIG_FIELD_UINT8_FILE", "255"},
				{"CONFIG_FIELD_UINT16_FILE", "65535"},
				{"CONFIG_FIELD_UINT32_FILE", "4294967295"},
				{"CONFIG_FIELD_UINT64_FILE", "18446744073709551615"},
				{"CONFIG_FIELD_DURATION_FILE", "90m"},
				{"CONFIG_FIELD_URL_FILE", "service-1:8080"},
				{"CONFIG_FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"CONFIG_FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"CONFIG_FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"CONFIG_FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"CONFIG_FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"CONFIG_FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"CONFIG_FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"CONFIG_FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"CONFIG_FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"CONFIG_FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"CONFIG_FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"CONFIG_FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"CONFIG_FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"CONFIG_FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			&controller{
				listSep:       ",",
				prefixFileEnv: "CONFIG_",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			34,
		},
		{
			"AllFromFromFilesWithListSepOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad|mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false|true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415|2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359|2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648|2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128|127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768|32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648|2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808|9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0|4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0|255"},
				{"FIELD_UINT16_ARRAY_FILE", "0|65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0|4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0|18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m|120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080|service-2:8080"},
			},
			&controller{
				listSep:       "|",
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			34,
		},
		{
			"FromMixedSources",
			[]string{
				"path/to/binary",
				"-field.float32=3.1415",
				"--field.float64=3.14159265359",
				"-field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
			},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
			},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
			},
			&controller{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "default",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			11,
		},
		{
			"FromMixedSourcesWithSkipOptions",
			[]string{
				"path/to/binary",
				"-field.float32=3.1415",
				"--field.float64=3.14159265359",
				"-field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
			},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
			},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
			},
			&controller{
				listSep:       ",",
				skipFlag:      true,
				skipEnv:       true,
				skipFileEnv:   true,
				filesToFields: map[string]fieldInfo{},
			},
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			0,
		},
		{
			"WithTelepresenceOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			&controller{
				listSep:       ",",
				telepresence:  true,
				filesToFields: map[string]fieldInfo{},
			},
			&config{},
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			34,
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.args

			// Set environment variables
			for _, e := range tc.envs {
				err := os.Setenv(e.varName, e.value)
				assert.NoError(t, err)
				defer os.Unsetenv(e.varName)
			}

			// Testing Telepresence option
			if tc.c.telepresence {
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

			vStruct, err := validateStruct(tc.config)
			assert.NoError(t, err)

			tc.c.readFields(vStruct)
			assert.Equal(t, tc.expectedConfig, tc.config)
			assert.Equal(t, tc.expectedFilesLen, len(tc.c.filesToFields))
		})
	}
}

func TestPick(t *testing.T) {
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName string
		value   string
	}

	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")

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
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			new(string),
			nil,
			errors.New("a non-struct type is passed"),
			&config{},
		},
		{
			"NonPointer",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			config{},
			nil,
			errors.New("a non-pointer type is passed"),
			&config{},
		},
		{
			"Empty",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{},
		},
		{
			"AllFromDefaults",
			[]string{"path/to/binary"},
			[]env{},
			[]file{},
			&config{
				unexported:         "internal",
				SkipFlag:           "default",
				SkipFlagEnv:        "default",
				SkipFlagEnvFile:    "default",
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			nil,
			nil,
			&config{
				unexported:         "internal",
				SkipFlag:           "default",
				SkipFlagEnv:        "default",
				SkipFlagEnvFile:    "default",
				FieldString:        "default",
				FieldBool:          false,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags#1",
			[]string{
				"path/to/binary",
				"-field.string", "content",
				"-field.bool",
				"-field.float32", "3.1415",
				"-field.float64", "3.14159265359",
				"-field.int", "-2147483648",
				"-field.int8", "-128",
				"-field.int16", "-32768",
				"-field.int32", "-2147483648",
				"-field.int64", "-9223372036854775808",
				"-field.uint", "4294967295",
				"-field.uint8", "255",
				"-field.uint16", "65535",
				"-field.uint32", "4294967295",
				"-field.uint64", "18446744073709551615",
				"-field.duration", "90m",
				"-field.url", "service-1:8080",
				"-field.string.array", "milad,mona",
				"-field.bool.array", "false,true",
				"-field.float32.array", "3.1415,2.7182",
				"-field.float64.array", "3.14159265359,2.71828182845",
				"-field.int.array", "-2147483648,2147483647",
				"-field.int8.array", "-128,127",
				"-field.int16.array", "-32768,32767",
				"-field.int32.array", "-2147483648,2147483647",
				"-field.int64.array", "-9223372036854775808,9223372036854775807",
				"-field.uint.array", "0,4294967295",
				"-field.uint8.array", "0,255",
				"-field.uint16.array", "0,65535",
				"-field.uint32.array", "0,4294967295",
				"-field.uint64.array", "0,18446744073709551615",
				"-field.duration.array", "90m,120m",
				"-field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags#2",
			[]string{
				"path/to/binary",
				"--field.string", "content",
				"--field.bool",
				"--field.float32", "3.1415",
				"--field.float64", "3.14159265359",
				"--field.int", "-2147483648",
				"--field.int8", "-128",
				"--field.int16", "-32768",
				"--field.int32", "-2147483648",
				"--field.int64", "-9223372036854775808",
				"--field.uint", "4294967295",
				"--field.uint8", "255",
				"--field.uint16", "65535",
				"--field.uint32", "4294967295",
				"--field.uint64", "18446744073709551615",
				"--field.duration", "90m",
				"--field.url", "service-1:8080",
				"--field.string.array", "milad,mona",
				"--field.bool.array", "false,true",
				"--field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
				"--field.int.array", "-2147483648,2147483647",
				"--field.int8.array", "-128,127",
				"--field.int16.array", "-32768,32767",
				"--field.int32.array", "-2147483648,2147483647",
				"--field.int64.array", "-9223372036854775808,9223372036854775807",
				"--field.uint.array", "0,4294967295",
				"--field.uint8.array", "0,255",
				"--field.uint16.array", "0,65535",
				"--field.uint32.array", "0,4294967295",
				"--field.uint64.array", "0,18446744073709551615",
				"--field.duration.array", "90m,120m",
				"--field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags#3",
			[]string{
				"path/to/binary",
				"-field.string=content",
				"-field.bool",
				"-field.float32=3.1415",
				"-field.float64=3.14159265359",
				"-field.int=-2147483648",
				"-field.int8=-128",
				"-field.int16=-32768",
				"-field.int32=-2147483648",
				"-field.int64=-9223372036854775808",
				"-field.uint=4294967295",
				"-field.uint8=255",
				"-field.uint16=65535",
				"-field.uint32=4294967295",
				"-field.uint64=18446744073709551615",
				"-field.duration=90m",
				"-field.url=service-1:8080",
				"-field.string.array=milad,mona",
				"-field.bool.array=false,true",
				"-field.float32.array=3.1415,2.7182",
				"-field.float64.array=3.14159265359,2.71828182845",
				"-field.int.array=-2147483648,2147483647",
				"-field.int8.array=-128,127",
				"-field.int16.array=-32768,32767",
				"-field.int32.array=-2147483648,2147483647",
				"-field.int64.array=-9223372036854775808,9223372036854775807",
				"-field.uint.array=0,4294967295",
				"-field.uint8.array=0,255",
				"-field.uint16.array=0,65535",
				"-field.uint32.array=0,4294967295",
				"-field.uint64.array=0,18446744073709551615",
				"-field.duration.array=90m,120m",
				"-field.url.array=service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlags#4",
			[]string{
				"path/to/binary",
				"--field.string=content",
				"--field.bool",
				"--field.float32=3.1415",
				"--field.float64=3.14159265359",
				"--field.int=-2147483648",
				"--field.int8=-128",
				"--field.int16=-32768",
				"--field.int32=-2147483648",
				"--field.int64=-9223372036854775808",
				"--field.uint=4294967295",
				"--field.uint8=255",
				"--field.uint16=65535",
				"--field.uint32=4294967295",
				"--field.uint64=18446744073709551615",
				"--field.duration=90m",
				"--field.url=service-1:8080",
				"--field.string.array=milad,mona",
				"--field.bool.array=false,true",
				"--field.float32.array=3.1415,2.7182",
				"--field.float64.array=3.14159265359,2.71828182845",
				"--field.int.array=-2147483648,2147483647",
				"--field.int8.array=-128,127",
				"--field.int16.array=-32768,32767",
				"--field.int32.array=-2147483648,2147483647",
				"--field.int64.array=-9223372036854775808,9223372036854775807",
				"--field.uint.array=0,4294967295",
				"--field.uint8.array=0,255",
				"--field.uint16.array=0,65535",
				"--field.uint32.array=0,4294967295",
				"--field.uint64.array=0,18446744073709551615",
				"--field.duration.array=90m,120m",
				"--field.url.array=service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&config{},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlagsWithPrefixFlagOption",
			[]string{
				"path/to/binary",
				"-config.field.string", "content",
				"-config.field.bool",
				"-config.field.float32", "3.1415",
				"-config.field.float64", "3.14159265359",
				"-config.field.int", "-2147483648",
				"-config.field.int8", "-128",
				"-config.field.int16", "-32768",
				"-config.field.int32", "-2147483648",
				"-config.field.int64", "-9223372036854775808",
				"-config.field.uint", "4294967295",
				"-config.field.uint8", "255",
				"-config.field.uint16", "65535",
				"-config.field.uint32", "4294967295",
				"-config.field.uint64", "18446744073709551615",
				"-config.field.duration", "90m",
				"-config.field.url", "service-1:8080",
				"-config.field.string.array", "milad,mona",
				"-config.field.bool.array", "false,true",
				"-config.field.float32.array", "3.1415,2.7182",
				"-config.field.float64.array", "3.14159265359,2.71828182845",
				"-config.field.int.array", "-2147483648,2147483647",
				"-config.field.int8.array", "-128,127",
				"-config.field.int16.array", "-32768,32767",
				"-config.field.int32.array", "-2147483648,2147483647",
				"-config.field.int64.array", "-9223372036854775808,9223372036854775807",
				"-config.field.uint.array", "0,4294967295",
				"-config.field.uint8.array", "0,255",
				"-config.field.uint16.array", "0,65535",
				"-config.field.uint32.array", "0,4294967295",
				"-config.field.uint64.array", "0,18446744073709551615",
				"-config.field.duration.array", "90m,120m",
				"-config.field.url.array", "service-1:8080,service-2:8080",
			},
			[]env{},
			[]file{},
			&config{},
			[]Option{
				PrefixFlag("config."),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFlagsWithListSepOption",
			[]string{
				"path/to/binary",
				"-field.string", "content",
				"-field.bool",
				"-field.float32", "3.1415",
				"-field.float64", "3.14159265359",
				"-field.int", "-2147483648",
				"-field.int8", "-128",
				"-field.int16", "-32768",
				"-field.int32", "-2147483648",
				"-field.int64", "-9223372036854775808",
				"-field.uint", "4294967295",
				"-field.uint8", "255",
				"-field.uint16", "65535",
				"-field.uint32", "4294967295",
				"-field.uint64", "18446744073709551615",
				"-field.duration", "90m",
				"-field.url", "service-1:8080",
				"-field.string.array", "milad|mona",
				"-field.bool.array", "false|true",
				"-field.float32.array", "3.1415|2.7182",
				"-field.float64.array", "3.14159265359|2.71828182845",
				"-field.int.array", "-2147483648|2147483647",
				"-field.int8.array", "-128|127",
				"-field.int16.array", "-32768|32767",
				"-field.int32.array", "-2147483648|2147483647",
				"-field.int64.array", "-9223372036854775808|9223372036854775807",
				"-field.uint.array", "0|4294967295",
				"-field.uint8.array", "0|255",
				"-field.uint16.array", "0|65535",
				"-field.uint32.array", "0|4294967295",
				"-field.uint64.array", "0|18446744073709551615",
				"-field.duration.array", "90m|120m",
				"-field.url.array", "service-1:8080|service-2:8080",
			},
			[]env{},
			[]file{},
			&config{},
			[]Option{
				ListSep("|"),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromEnvVars",
			[]string{"path/to/binary"},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_STRING", "content"},
				{"FIELD_BOOL", "true"},
				{"FIELD_FLOAT32", "3.1415"},
				{"FIELD_FLOAT64", "3.14159265359"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_UINT", "4294967295"},
				{"FIELD_UINT8", "255"},
				{"FIELD_UINT16", "65535"},
				{"FIELD_UINT32", "4294967295"},
				{"FIELD_UINT64", "18446744073709551615"},
				{"FIELD_DURATION", "90m"},
				{"FIELD_URL", "service-1:8080"},
				{"FIELD_STRING_ARRAY", "milad,mona"},
				{"FIELD_BOOL_ARRAY", "false,true"},
				{"FIELD_FLOAT32_ARRAY", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY", "0,4294967295"},
				{"FIELD_UINT8_ARRAY", "0,255"},
				{"FIELD_UINT16_ARRAY", "0,65535"},
				{"FIELD_UINT32_ARRAY", "0,4294967295"},
				{"FIELD_UINT64_ARRAY", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY", "90m,120m"},
				{"FIELD_URL_ARRAY", "service-1:8080,service-2:8080"},
			},
			[]file{},
			&config{},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromEnvVarsWithPrefixEnvOption",
			[]string{"path/to/binary"},
			[]env{
				{"CONFIG_SKIP_FLAG", "from_env"},
				{"CONFIG_SKIP_FLAG_ENV", "from_env"},
				{"CONFIG_SKIP_FLAG_ENV_FILE", "from_env"},
				{"CONFIG_FIELD_STRING", "content"},
				{"CONFIG_FIELD_BOOL", "true"},
				{"CONFIG_FIELD_FLOAT32", "3.1415"},
				{"CONFIG_FIELD_FLOAT64", "3.14159265359"},
				{"CONFIG_FIELD_INT", "-2147483648"},
				{"CONFIG_FIELD_INT8", "-128"},
				{"CONFIG_FIELD_INT16", "-32768"},
				{"CONFIG_FIELD_INT32", "-2147483648"},
				{"CONFIG_FIELD_INT64", "-9223372036854775808"},
				{"CONFIG_FIELD_UINT", "4294967295"},
				{"CONFIG_FIELD_UINT8", "255"},
				{"CONFIG_FIELD_UINT16", "65535"},
				{"CONFIG_FIELD_UINT32", "4294967295"},
				{"CONFIG_FIELD_UINT64", "18446744073709551615"},
				{"CONFIG_FIELD_DURATION", "90m"},
				{"CONFIG_FIELD_URL", "service-1:8080"},
				{"CONFIG_FIELD_STRING_ARRAY", "milad,mona"},
				{"CONFIG_FIELD_BOOL_ARRAY", "false,true"},
				{"CONFIG_FIELD_FLOAT32_ARRAY", "3.1415,2.7182"},
				{"CONFIG_FIELD_FLOAT64_ARRAY", "3.14159265359,2.71828182845"},
				{"CONFIG_FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT8_ARRAY", "-128,127"},
				{"CONFIG_FIELD_INT16_ARRAY", "-32768,32767"},
				{"CONFIG_FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
				{"CONFIG_FIELD_UINT_ARRAY", "0,4294967295"},
				{"CONFIG_FIELD_UINT8_ARRAY", "0,255"},
				{"CONFIG_FIELD_UINT16_ARRAY", "0,65535"},
				{"CONFIG_FIELD_UINT32_ARRAY", "0,4294967295"},
				{"CONFIG_FIELD_UINT64_ARRAY", "0,18446744073709551615"},
				{"CONFIG_FIELD_DURATION_ARRAY", "90m,120m"},
				{"CONFIG_FIELD_URL_ARRAY", "service-1:8080,service-2:8080"},
			},
			[]file{},
			&config{},
			[]Option{
				PrefixEnv("CONFIG_"),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromEnvVarsWithListSepOption",
			[]string{"path/to/binary"},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_STRING", "content"},
				{"FIELD_BOOL", "true"},
				{"FIELD_FLOAT32", "3.1415"},
				{"FIELD_FLOAT64", "3.14159265359"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_UINT", "4294967295"},
				{"FIELD_UINT8", "255"},
				{"FIELD_UINT16", "65535"},
				{"FIELD_UINT32", "4294967295"},
				{"FIELD_UINT64", "18446744073709551615"},
				{"FIELD_DURATION", "90m"},
				{"FIELD_URL", "service-1:8080"},
				{"FIELD_STRING_ARRAY", "milad|mona"},
				{"FIELD_BOOL_ARRAY", "false|true"},
				{"FIELD_FLOAT32_ARRAY", "3.1415|2.7182"},
				{"FIELD_FLOAT64_ARRAY", "3.14159265359|2.71828182845"},
				{"FIELD_INT_ARRAY", "-2147483648|2147483647"},
				{"FIELD_INT8_ARRAY", "-128|127"},
				{"FIELD_INT16_ARRAY", "-32768|32767"},
				{"FIELD_INT32_ARRAY", "-2147483648|2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808|9223372036854775807"},
				{"FIELD_UINT_ARRAY", "0|4294967295"},
				{"FIELD_UINT8_ARRAY", "0|255"},
				{"FIELD_UINT16_ARRAY", "0|65535"},
				{"FIELD_UINT32_ARRAY", "0|4294967295"},
				{"FIELD_UINT64_ARRAY", "0|18446744073709551615"},
				{"FIELD_DURATION_ARRAY", "90m|120m"},
				{"FIELD_URL_ARRAY", "service-1:8080|service-2:8080"},
			},
			[]file{},
			&config{},
			[]Option{
				ListSep("|"),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFromFiles",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			&config{},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFromFilesWithPrefixFileEnvOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"CONFIG_SKIP_FLAG_FILE", "from_file"},
				{"CONFIG_SKIP_FLAG_ENV_FILE", "from_file"},
				{"CONFIG_SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"CONFIG_FIELD_STRING_FILE", "content"},
				{"CONFIG_FIELD_BOOL_FILE", "true"},
				{"CONFIG_FIELD_FLOAT32_FILE", "3.1415"},
				{"CONFIG_FIELD_FLOAT64_FILE", "3.14159265359"},
				{"CONFIG_FIELD_INT_FILE", "-2147483648"},
				{"CONFIG_FIELD_INT8_FILE", "-128"},
				{"CONFIG_FIELD_INT16_FILE", "-32768"},
				{"CONFIG_FIELD_INT32_FILE", "-2147483648"},
				{"CONFIG_FIELD_INT64_FILE", "-9223372036854775808"},
				{"CONFIG_FIELD_UINT_FILE", "4294967295"},
				{"CONFIG_FIELD_UINT8_FILE", "255"},
				{"CONFIG_FIELD_UINT16_FILE", "65535"},
				{"CONFIG_FIELD_UINT32_FILE", "4294967295"},
				{"CONFIG_FIELD_UINT64_FILE", "18446744073709551615"},
				{"CONFIG_FIELD_DURATION_FILE", "90m"},
				{"CONFIG_FIELD_URL_FILE", "service-1:8080"},
				{"CONFIG_FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"CONFIG_FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"CONFIG_FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"CONFIG_FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"CONFIG_FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"CONFIG_FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"CONFIG_FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"CONFIG_FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"CONFIG_FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"CONFIG_FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"CONFIG_FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"CONFIG_FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"CONFIG_FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"CONFIG_FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"CONFIG_FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			&config{},
			[]Option{
				PrefixFileEnv("CONFIG_"),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"AllFromFromFilesWithListSepOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad|mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false|true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415|2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359|2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648|2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128|127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768|32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648|2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808|9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0|4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0|255"},
				{"FIELD_UINT16_ARRAY_FILE", "0|65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0|4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0|18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m|120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080|service-2:8080"},
			},
			&config{},
			[]Option{
				ListSep("|"),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"FromMixedSources",
			[]string{
				"path/to/binary",
				"-field.float32=3.1415",
				"--field.float64=3.14159265359",
				"-field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
			},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
			},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
			},
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			nil,
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_env",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "default",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"FromMixedSourcesWithSkipOptions",
			[]string{
				"path/to/binary",
				"-field.float32=3.1415",
				"--field.float64=3.14159265359",
				"-field.float32.array", "3.1415,2.7182",
				"--field.float64.array", "3.14159265359,2.71828182845",
			},
			[]env{
				{"SKIP_FLAG", "from_env"},
				{"SKIP_FLAG_ENV", "from_env"},
				{"SKIP_FLAG_ENV_FILE", "from_env"},
				{"FIELD_INT", "-2147483648"},
				{"FIELD_INT8", "-128"},
				{"FIELD_INT16", "-32768"},
				{"FIELD_INT32", "-2147483648"},
				{"FIELD_INT64", "-9223372036854775808"},
				{"FIELD_INT_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY", "-128,127"},
				{"FIELD_INT16_ARRAY", "-32768,32767"},
				{"FIELD_INT32_ARRAY", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY", "-9223372036854775808,9223372036854775807"},
			},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
			},
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			[]Option{
				SkipFlag(),
				SkipEnv(),
				SkipFileEnv(),
			},
			nil,
			&config{
				FieldString:        "default",
				FieldStringArray:   []string{"milad", "mona"},
				FieldBool:          true,
				FieldBoolArray:     []bool{false, true},
				FieldDuration:      d90m,
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURL:           *service1URL,
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
		{
			"WithTelepresenceOption",
			[]string{"path/to/binary"},
			[]env{},
			[]file{
				{"SKIP_FLAG_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE", "from_file"},
				{"SKIP_FLAG_ENV_FILE_FILE", "from_file"},
				{"FIELD_STRING_FILE", "content"},
				{"FIELD_BOOL_FILE", "true"},
				{"FIELD_FLOAT32_FILE", "3.1415"},
				{"FIELD_FLOAT64_FILE", "3.14159265359"},
				{"FIELD_INT_FILE", "-2147483648"},
				{"FIELD_INT8_FILE", "-128"},
				{"FIELD_INT16_FILE", "-32768"},
				{"FIELD_INT32_FILE", "-2147483648"},
				{"FIELD_INT64_FILE", "-9223372036854775808"},
				{"FIELD_UINT_FILE", "4294967295"},
				{"FIELD_UINT8_FILE", "255"},
				{"FIELD_UINT16_FILE", "65535"},
				{"FIELD_UINT32_FILE", "4294967295"},
				{"FIELD_UINT64_FILE", "18446744073709551615"},
				{"FIELD_DURATION_FILE", "90m"},
				{"FIELD_URL_FILE", "service-1:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080"},
			},
			&config{},
			[]Option{
				Telepresence(),
			},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "from_file",
				SkipFlagEnv:        "from_file",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
		},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &controller{}
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
	type env struct {
		varName string
		value   string
	}

	type file struct {
		varName   string
		initValue string
		newValue  string
	}

	updateDelay := 50 * time.Millisecond

	d90m := 90 * time.Minute
	d120m := 120 * time.Minute
	d4h := 4 * time.Hour
	d8h := 8 * time.Hour

	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")
	service3URL, _ := url.Parse("service-3:8080")
	service4URL, _ := url.Parse("service-4:8080")

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
			[]string{
				"path/to/binary",
				"-field.bool",
			},
			[]env{
				{"FIELD_STRING", "content"},
			},
			[]file{
				{"FIELD_FLOAT32_FILE", "3.1415", "2.7182"},
				{"FIELD_FLOAT64_FILE", "3.14159265359", "2.7182818284"},
				{"FIELD_INT_FILE", "-2147483648", "2147483647"},
				{"FIELD_INT8_FILE", "-128", "127"},
				{"FIELD_INT16_FILE", "-32768", "32767"},
				{"FIELD_INT32_FILE", "-2147483648", "2147483647"},
				{"FIELD_INT64_FILE", "-9223372036854775808", "9223372036854775807"},
				{"FIELD_UINT_FILE", "4294967295", "2147483648"},
				{"FIELD_UINT8_FILE", "255", "128"},
				{"FIELD_UINT16_FILE", "65535", "32768"},
				{"FIELD_UINT32_FILE", "4294967295", "2147483648"},
				{"FIELD_UINT64_FILE", "18446744073709551615", "9223372036854775808"},
				{"FIELD_DURATION_FILE", "90m", "4h"},
				{"FIELD_URL_FILE", "service-1:8080", "service-3:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona", "mona,milad"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true", "true,false"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182", "2.7182,3.1415"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845", "2.71828182845,3.14159265359"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647", "2147483647,-2147483648"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127", "127,-128"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767", "32767,-32768"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647", "2147483647,-2147483648"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807", "9223372036854775807,-9223372036854775808"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295", "4294967295,0"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255", "255,0"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535", "65535,0"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295", "4294967295,0"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615", "18446744073709551615,0"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m", "4h,8h"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080", "service-3:8080,service-4:8080"},
			},
			&config{},
			[]chan Update{
				make(chan Update),
				make(chan Update),
			},
			[]Option{},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       2.7182,
				FieldFloat64:       2.7182818284,
				FieldInt:           2147483647,
				FieldInt8:          127,
				FieldInt16:         32767,
				FieldInt32:         2147483647,
				FieldInt64:         9223372036854775807,
				FieldUint:          2147483648,
				FieldUint8:         128,
				FieldUint16:        32768,
				FieldUint32:        2147483648,
				FieldUint64:        9223372036854775808,
				FieldDuration:      d4h,
				FieldURL:           *service3URL,
				FieldStringArray:   []string{"mona", "milad"},
				FieldBoolArray:     []bool{true, false},
				FieldFloat32Array:  []float32{2.7182, 3.1415},
				FieldFloat64Array:  []float64{2.71828182845, 3.14159265359},
				FieldIntArray:      []int{2147483647, -2147483648},
				FieldInt8Array:     []int8{127, -128},
				FieldInt16Array:    []int16{32767, -32768},
				FieldInt32Array:    []int32{2147483647, -2147483648},
				FieldInt64Array:    []int64{9223372036854775807, -9223372036854775808},
				FieldUintArray:     []uint{4294967295, 0},
				FieldUint8Array:    []uint8{255, 0},
				FieldUint16Array:   []uint16{65535, 0},
				FieldUint32Array:   []uint32{4294967295, 0},
				FieldUint64Array:   []uint64{18446744073709551615, 0},
				FieldDurationArray: []time.Duration{d4h, d8h},
				FieldURLArray:      []url.URL{*service3URL, *service4URL},
			},
			[]Update{
				{"FieldString", "content"},
				{"FieldBool", true},
				{"FieldFloat32", float32(3.1415)},
				{"FieldFloat64", float64(3.14159265359)},
				{"FieldInt", int(-2147483648)},
				{"FieldInt8", int8(-128)},
				{"FieldInt16", int16(-32768)},
				{"FieldInt32", int32(-2147483648)},
				{"FieldInt64", int64(-9223372036854775808)},
				{"FieldUint", uint(4294967295)},
				{"FieldUint8", uint8(255)},
				{"FieldUint16", uint16(65535)},
				{"FieldUint32", uint32(4294967295)},
				{"FieldUint64", uint64(18446744073709551615)},
				{"FieldDuration", d90m},
				{"FieldURL", *service1URL},
				{"FieldStringArray", []string{"milad", "mona"}},
				{"FieldBoolArray", []bool{false, true}},
				{"FieldFloat32Array", []float32{3.1415, 2.7182}},
				{"FieldFloat64Array", []float64{3.14159265359, 2.71828182845}},
				{"FieldIntArray", []int{-2147483648, 2147483647}},
				{"FieldInt8Array", []int8{-128, 127}},
				{"FieldInt16Array", []int16{-32768, 32767}},
				{"FieldInt32Array", []int32{-2147483648, 2147483647}},
				{"FieldInt64Array", []int64{-9223372036854775808, 9223372036854775807}},
				{"FieldUintArray", []uint{0, 4294967295}},
				{"FieldUint8Array", []uint8{0, 255}},
				{"FieldUint16Array", []uint16{0, 65535}},
				{"FieldUint32Array", []uint32{0, 4294967295}},
				{"FieldUint64Array", []uint64{0, 18446744073709551615}},
				{"FieldDurationArray", []time.Duration{d90m, d120m}},
				{"FieldURLArray", []url.URL{*service1URL, *service2URL}},

				{"FieldFloat32", float32(2.7182)},
				{"FieldFloat64", float64(2.7182818284)},
				{"FieldInt", int(2147483647)},
				{"FieldInt8", int8(127)},
				{"FieldInt16", int16(32767)},
				{"FieldInt32", int32(2147483647)},
				{"FieldInt64", int64(9223372036854775807)},
				{"FieldUint", uint(2147483648)},
				{"FieldUint8", uint8(128)},
				{"FieldUint16", uint16(32768)},
				{"FieldUint32", uint32(2147483648)},
				{"FieldUint64", uint64(9223372036854775808)},
				{"FieldDuration", d4h},
				{"FieldURL", *service3URL},
				{"FieldStringArray", []string{"mona", "milad"}},
				{"FieldBoolArray", []bool{true, false}},
				{"FieldFloat32Array", []float32{2.7182, 3.1415}},
				{"FieldFloat64Array", []float64{2.71828182845, 3.14159265359}},
				{"FieldIntArray", []int{2147483647, -2147483648}},
				{"FieldInt8Array", []int8{127, -128}},
				{"FieldInt16Array", []int16{32767, -32768}},
				{"FieldInt32Array", []int32{2147483647, -2147483648}},
				{"FieldInt64Array", []int64{9223372036854775807, -9223372036854775808}},
				{"FieldUintArray", []uint{4294967295, 0}},
				{"FieldUint8Array", []uint8{255, 0}},
				{"FieldUint16Array", []uint16{65535, 0}},
				{"FieldUint32Array", []uint32{4294967295, 0}},
				{"FieldUint64Array", []uint64{18446744073709551615, 0}},
				{"FieldDurationArray", []time.Duration{d4h, d8h}},
				{"FieldURLArray", []url.URL{*service3URL, *service4URL}},
			},
		},
		{
			"BufferedChannels",
			[]string{
				"path/to/binary",
				"-field.bool",
			},
			[]env{
				{"FIELD_STRING", "content"},
			},
			[]file{
				{"FIELD_FLOAT32_FILE", "3.1415", "2.7182"},
				{"FIELD_FLOAT64_FILE", "3.14159265359", "2.7182818284"},
				{"FIELD_INT_FILE", "-2147483648", "2147483647"},
				{"FIELD_INT8_FILE", "-128", "127"},
				{"FIELD_INT16_FILE", "-32768", "32767"},
				{"FIELD_INT32_FILE", "-2147483648", "2147483647"},
				{"FIELD_INT64_FILE", "-9223372036854775808", "9223372036854775807"},
				{"FIELD_UINT_FILE", "4294967295", "2147483648"},
				{"FIELD_UINT8_FILE", "255", "128"},
				{"FIELD_UINT16_FILE", "65535", "32768"},
				{"FIELD_UINT32_FILE", "4294967295", "2147483648"},
				{"FIELD_UINT64_FILE", "18446744073709551615", "9223372036854775808"},
				{"FIELD_DURATION_FILE", "90m", "4h"},
				{"FIELD_URL_FILE", "service-1:8080", "service-3:8080"},
				{"FIELD_STRING_ARRAY_FILE", "milad,mona", "mona,milad"},
				{"FIELD_BOOL_ARRAY_FILE", "false,true", "true,false"},
				{"FIELD_FLOAT32_ARRAY_FILE", "3.1415,2.7182", "2.7182,3.1415"},
				{"FIELD_FLOAT64_ARRAY_FILE", "3.14159265359,2.71828182845", "2.71828182845,3.14159265359"},
				{"FIELD_INT_ARRAY_FILE", "-2147483648,2147483647", "2147483647,-2147483648"},
				{"FIELD_INT8_ARRAY_FILE", "-128,127", "127,-128"},
				{"FIELD_INT16_ARRAY_FILE", "-32768,32767", "32767,-32768"},
				{"FIELD_INT32_ARRAY_FILE", "-2147483648,2147483647", "2147483647,-2147483648"},
				{"FIELD_INT64_ARRAY_FILE", "-9223372036854775808,9223372036854775807", "9223372036854775807,-9223372036854775808"},
				{"FIELD_UINT_ARRAY_FILE", "0,4294967295", "4294967295,0"},
				{"FIELD_UINT8_ARRAY_FILE", "0,255", "255,0"},
				{"FIELD_UINT16_ARRAY_FILE", "0,65535", "65535,0"},
				{"FIELD_UINT32_ARRAY_FILE", "0,4294967295", "4294967295,0"},
				{"FIELD_UINT64_ARRAY_FILE", "0,18446744073709551615", "18446744073709551615,0"},
				{"FIELD_DURATION_ARRAY_FILE", "90m,120m", "4h,8h"},
				{"FIELD_URL_ARRAY_FILE", "service-1:8080,service-2:8080", "service-3:8080,service-4:8080"},
			},
			&config{},
			[]chan Update{
				make(chan Update, 100),
				make(chan Update, 100),
			},
			[]Option{},
			nil,
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       3.1415,
				FieldFloat64:       3.14159265359,
				FieldInt:           -2147483648,
				FieldInt8:          -128,
				FieldInt16:         -32768,
				FieldInt32:         -2147483648,
				FieldInt64:         -9223372036854775808,
				FieldUint:          4294967295,
				FieldUint8:         255,
				FieldUint16:        65535,
				FieldUint32:        4294967295,
				FieldUint64:        18446744073709551615,
				FieldDuration:      d90m,
				FieldURL:           *service1URL,
				FieldStringArray:   []string{"milad", "mona"},
				FieldBoolArray:     []bool{false, true},
				FieldFloat32Array:  []float32{3.1415, 2.7182},
				FieldFloat64Array:  []float64{3.14159265359, 2.71828182845},
				FieldIntArray:      []int{-2147483648, 2147483647},
				FieldInt8Array:     []int8{-128, 127},
				FieldInt16Array:    []int16{-32768, 32767},
				FieldInt32Array:    []int32{-2147483648, 2147483647},
				FieldInt64Array:    []int64{-9223372036854775808, 9223372036854775807},
				FieldUintArray:     []uint{0, 4294967295},
				FieldUint8Array:    []uint8{0, 255},
				FieldUint16Array:   []uint16{0, 65535},
				FieldUint32Array:   []uint32{0, 4294967295},
				FieldUint64Array:   []uint64{0, 18446744073709551615},
				FieldDurationArray: []time.Duration{d90m, d120m},
				FieldURLArray:      []url.URL{*service1URL, *service2URL},
			},
			&config{
				unexported:         "",
				SkipFlag:           "",
				SkipFlagEnv:        "",
				SkipFlagEnvFile:    "",
				FieldString:        "content",
				FieldBool:          true,
				FieldFloat32:       2.7182,
				FieldFloat64:       2.7182818284,
				FieldInt:           2147483647,
				FieldInt8:          127,
				FieldInt16:         32767,
				FieldInt32:         2147483647,
				FieldInt64:         9223372036854775807,
				FieldUint:          2147483648,
				FieldUint8:         128,
				FieldUint16:        32768,
				FieldUint32:        2147483648,
				FieldUint64:        9223372036854775808,
				FieldDuration:      d4h,
				FieldURL:           *service3URL,
				FieldStringArray:   []string{"mona", "milad"},
				FieldBoolArray:     []bool{true, false},
				FieldFloat32Array:  []float32{2.7182, 3.1415},
				FieldFloat64Array:  []float64{2.71828182845, 3.14159265359},
				FieldIntArray:      []int{2147483647, -2147483648},
				FieldInt8Array:     []int8{127, -128},
				FieldInt16Array:    []int16{32767, -32768},
				FieldInt32Array:    []int32{2147483647, -2147483648},
				FieldInt64Array:    []int64{9223372036854775807, -9223372036854775808},
				FieldUintArray:     []uint{4294967295, 0},
				FieldUint8Array:    []uint8{255, 0},
				FieldUint16Array:   []uint16{65535, 0},
				FieldUint32Array:   []uint32{4294967295, 0},
				FieldUint64Array:   []uint64{18446744073709551615, 0},
				FieldDurationArray: []time.Duration{d4h, d8h},
				FieldURLArray:      []url.URL{*service3URL, *service4URL},
			},
			[]Update{
				{"FieldString", "content"},
				{"FieldBool", true},
				{"FieldFloat32", float32(3.1415)},
				{"FieldFloat64", float64(3.14159265359)},
				{"FieldInt", int(-2147483648)},
				{"FieldInt8", int8(-128)},
				{"FieldInt16", int16(-32768)},
				{"FieldInt32", int32(-2147483648)},
				{"FieldInt64", int64(-9223372036854775808)},
				{"FieldUint", uint(4294967295)},
				{"FieldUint8", uint8(255)},
				{"FieldUint16", uint16(65535)},
				{"FieldUint32", uint32(4294967295)},
				{"FieldUint64", uint64(18446744073709551615)},
				{"FieldDuration", d90m},
				{"FieldURL", *service1URL},
				{"FieldStringArray", []string{"milad", "mona"}},
				{"FieldBoolArray", []bool{false, true}},
				{"FieldFloat32Array", []float32{3.1415, 2.7182}},
				{"FieldFloat64Array", []float64{3.14159265359, 2.71828182845}},
				{"FieldIntArray", []int{-2147483648, 2147483647}},
				{"FieldInt8Array", []int8{-128, 127}},
				{"FieldInt16Array", []int16{-32768, 32767}},
				{"FieldInt32Array", []int32{-2147483648, 2147483647}},
				{"FieldInt64Array", []int64{-9223372036854775808, 9223372036854775807}},
				{"FieldUintArray", []uint{0, 4294967295}},
				{"FieldUint8Array", []uint8{0, 255}},
				{"FieldUint16Array", []uint16{0, 65535}},
				{"FieldUint32Array", []uint32{0, 4294967295}},
				{"FieldUint64Array", []uint64{0, 18446744073709551615}},
				{"FieldDurationArray", []time.Duration{d90m, d120m}},
				{"FieldURLArray", []url.URL{*service1URL, *service2URL}},

				{"FieldFloat32", float32(2.7182)},
				{"FieldFloat64", float64(2.7182818284)},
				{"FieldInt", int(2147483647)},
				{"FieldInt8", int8(127)},
				{"FieldInt16", int16(32767)},
				{"FieldInt32", int32(2147483647)},
				{"FieldInt64", int64(9223372036854775807)},
				{"FieldUint", uint(2147483648)},
				{"FieldUint8", uint8(128)},
				{"FieldUint16", uint16(32768)},
				{"FieldUint32", uint32(2147483648)},
				{"FieldUint64", uint64(9223372036854775808)},
				{"FieldDuration", d4h},
				{"FieldURL", *service3URL},
				{"FieldStringArray", []string{"mona", "milad"}},
				{"FieldBoolArray", []bool{true, false}},
				{"FieldFloat32Array", []float32{2.7182, 3.1415}},
				{"FieldFloat64Array", []float64{2.71828182845, 3.14159265359}},
				{"FieldIntArray", []int{2147483647, -2147483648}},
				{"FieldInt8Array", []int8{127, -128}},
				{"FieldInt16Array", []int16{32767, -32768}},
				{"FieldInt32Array", []int32{2147483647, -2147483648}},
				{"FieldInt64Array", []int64{9223372036854775807, -9223372036854775808}},
				{"FieldUintArray", []uint{4294967295, 0}},
				{"FieldUint8Array", []uint8{255, 0}},
				{"FieldUint16Array", []uint16{65535, 0}},
				{"FieldUint32Array", []uint32{4294967295, 0}},
				{"FieldUint64Array", []uint64{18446744073709551615, 0}},
				{"FieldDurationArray", []time.Duration{d4h, d8h}},
				{"FieldURLArray", []url.URL{*service3URL, *service4URL}},
			},
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
			c := &controller{}
			for _, opt := range tc.opts {
				opt(c)
			}
			if c.telepresence {
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
				assert.True(t, configEqual(tc.expectedInitConfig, tc.config))
				tc.config.Unlock()

				// Wait for all files to be updated and the new values are picked up
				wg.Wait()
				time.Sleep(100 * time.Millisecond)

				tc.config.Lock()
				assert.True(t, configEqual(tc.expectedNewConfig, tc.config))
				tc.config.Unlock()
			}
		})
	}

	// flag.Parse() can be called only once
	flag.Parse()
}
