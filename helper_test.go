package konfig

import (
	"errors"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/moorara/konfig/ptr"
	"github.com/stretchr/testify/assert"
)

func TestFlagValue(t *testing.T) {
	fv := new(flagValue)

	assert.Empty(t, fv.String())
	assert.NoError(t, fv.Set(""))
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		fieldName      string
		expectedTokens []string
	}{
		{"c", []string{"c"}},
		{"C", []string{"C"}},
		{"camel", []string{"camel"}},
		{"Camel", []string{"Camel"}},
		{"camelCase", []string{"camel", "Case"}},
		{"CamelCase", []string{"Camel", "Case"}},
		{"OneTwoThree", []string{"One", "Two", "Three"}},
		{"DatabaseURL", []string{"Database", "URL"}},
		{"DBEndpoints", []string{"DB", "Endpoints"}},
	}

	for _, tc := range tests {
		tokens := tokenize(tc.fieldName)
		assert.Equal(t, tc.expectedTokens, tokens)
	}
}

func TestGetFlagName(t *testing.T) {
	tests := []struct {
		fieldName        string
		expectedFlagName string
	}{
		{"c", "c"},
		{"C", "c"},
		{"camel", "camel"},
		{"Camel", "camel"},
		{"camelCase", "camel.case"},
		{"CamelCase", "camel.case"},
		{"OneTwoThree", "one.two.three"},
		{"DatabaseURL", "database.url"},
		{"DBEndpoints", "db.endpoints"},
	}

	for _, tc := range tests {
		flagName := getFlagName(tc.fieldName)
		assert.Equal(t, tc.expectedFlagName, flagName)
	}
}

func TestGetEnvVarName(t *testing.T) {
	tests := []struct {
		fieldName          string
		expectedEnvVarName string
	}{
		{"c", "C"},
		{"C", "C"},
		{"camel", "CAMEL"},
		{"Camel", "CAMEL"},
		{"camelCase", "CAMEL_CASE"},
		{"CamelCase", "CAMEL_CASE"},
		{"OneTwoThree", "ONE_TWO_THREE"},
		{"DatabaseURL", "DATABASE_URL"},
		{"DBEndpoints", "DB_ENDPOINTS"},
	}

	for _, tc := range tests {
		envVarName := getEnvVarName(tc.fieldName)
		assert.Equal(t, tc.expectedEnvVarName, envVarName)
	}
}

func TestGetFileEnvVarName(t *testing.T) {
	tests := []struct {
		fieldName              string
		expectedFileEnvVarName string
	}{
		{"c", "C_FILE"},
		{"C", "C_FILE"},
		{"camel", "CAMEL_FILE"},
		{"Camel", "CAMEL_FILE"},
		{"camelCase", "CAMEL_CASE_FILE"},
		{"CamelCase", "CAMEL_CASE_FILE"},
		{"OneTwoThree", "ONE_TWO_THREE_FILE"},
		{"DatabaseURL", "DATABASE_URL_FILE"},
		{"DBEndpoints", "DB_ENDPOINTS_FILE"},
	}

	for _, tc := range tests {
		fileEnvVarName := getFileEnvVarName(tc.fieldName)
		assert.Equal(t, tc.expectedFileEnvVarName, fileEnvVarName)
	}
}

func TestGetFlagValue(t *testing.T) {
	tests := []struct {
		args              []string
		flagName          string
		expectedFlagValue string
	}{
		{[]string{"exe", "invalid"}, "invalid", ""},

		{[]string{"app", "-enabled"}, "enabled", "true"},
		{[]string{"app", "--enabled"}, "enabled", "true"},
		{[]string{"app", "-enabled=false"}, "enabled", "false"},
		{[]string{"app", "--enabled=false"}, "enabled", "false"},
		{[]string{"app", "-enabled", "false"}, "enabled", "false"},
		{[]string{"app", "--enabled", "false"}, "enabled", "false"},

		{[]string{"app", "-number=-10"}, "number", "-10"},
		{[]string{"app", "--number=-10"}, "number", "-10"},
		{[]string{"app", "-number", "-10"}, "number", "-10"},
		{[]string{"app", "--number", "-10"}, "number", "-10"},

		{[]string{"app", "-text=content"}, "text", "content"},
		{[]string{"app", "--text=content"}, "text", "content"},
		{[]string{"app", "-text", "content"}, "text", "content"},
		{[]string{"app", "--text", "content"}, "text", "content"},

		{[]string{"app", "-enabled", "-text=content"}, "enabled", "true"},
		{[]string{"app", "--enabled", "--text=content"}, "enabled", "true"},
		{[]string{"app", "-enabled", "-text", "content"}, "enabled", "true"},
		{[]string{"app", "--enabled", "--text", "content"}, "enabled", "true"},

		{[]string{"app", "-name-list=alice,bob"}, "name-list", "alice,bob"},
		{[]string{"app", "--name-list=alice,bob"}, "name-list", "alice,bob"},
		{[]string{"app", "-name-list", "alice,bob"}, "name-list", "alice,bob"},
		{[]string{"app", "--name-list", "alice,bob"}, "name-list", "alice,bob"},
	}

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	for _, tc := range tests {
		os.Args = tc.args
		flagValue := getFlagValue(tc.flagName)

		assert.Equal(t, tc.expectedFlagValue, flagValue)
	}
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name          string
		s             interface{}
		expectedError error
	}{
		{
			"NonStruct",
			new(string),
			errors.New("a non-struct type is passed"),
		},
		{
			"NonPointer",
			struct{}{},
			errors.New("a non-pointer type is passed"),
		},
		{
			"OK",
			new(struct{}),
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := validateStruct(tc.s)

			if tc.expectedError == nil {
				assert.NotNil(t, v)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, v)
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}

func TestIsTypeSupported(t *testing.T) {
	u, _ := url.Parse("service-1")
	r := regexp.MustCompilePOSIX("[:digit:]")

	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{"String", "content", true},
		{"Bool", true, true},
		{"Float32", float32(3.1415), true},
		{"Float64", float64(3.14159265359), true},
		{"Int", int(-9223372036854775808), true},
		{"Int8", int8(-128), true},
		{"Int16", int16(-32768), true},
		{"Int32", int32(-2147483648), true},
		{"Int64", int64(-9223372036854775808), true},
		{"Uint", uint(18446744073709551615), true},
		{"Uint8", uint8(255), true},
		{"Uint16", uint16(65535), true},
		{"Uint32", uint32(4294967295), true},
		{"Uint64", uint64(18446744073709551615), true},
		{"URL", *u, true},
		{"Regexp", *r, true},
		{"Duration", time.Second, true},
		{"StringPointer", ptr.String("content"), true},
		{"BoolPointer", ptr.Bool(true), true},
		{"Float32Pointer", ptr.Float32(3.1415), true},
		{"Float64Pointer", ptr.Float64(3.14159265359), true},
		{"IntPointer", ptr.Int(-9223372036854775808), true},
		{"Int8Pointer", ptr.Int8(-128), true},
		{"Int16Pointer", ptr.Int16(-32768), true},
		{"Int32Pointer", ptr.Int32(-2147483648), true},
		{"Int64Pointer", ptr.Int64(-9223372036854775808), true},
		{"UintPointer", ptr.Uint(18446744073709551615), true},
		{"Uint8Pointer", ptr.Uint8(255), true},
		{"Uint16Pointer", ptr.Uint16(65535), true},
		{"Uint32Pointer", ptr.Uint32(4294967295), true},
		{"Uint64Pointer", ptr.Uint64(18446744073709551615), true},
		{"URLPointer", u, true},
		{"RegexpPointer", r, true},
		{"DurationPointer", ptr.Duration(time.Second), true},
		{"StringSlice", []string{"content"}, true},
		{"BoolSlice", []bool{true}, true},
		{"Float32Slice", []float32{3.1415}, true},
		{"Float64Slice", []float64{3.14159265359}, true},
		{"IntSlice", []int{-9223372036854775808}, true},
		{"Int8Slice", []int8{-128}, true},
		{"Int16Slice", []int16{-32768}, true},
		{"Int32Slice", []int32{-2147483648}, true},
		{"Int64Slice", []int64{-9223372036854775808}, true},
		{"UintSlice", []uint{18446744073709551615}, true},
		{"Uint8Slice", []uint8{255}, true},
		{"Uint16Slice", []uint16{65535}, true},
		{"Uint32Slice", []uint32{4294967295}, true},
		{"Uint64Slice", []uint64{18446744073709551615}, true},
		{"URLSlice", []url.URL{*u}, true},
		{"RegexpSlice", []regexp.Regexp{*r}, true},
		{"DurationSlice", []time.Duration{time.Second}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			typ := reflect.TypeOf(tc.field)

			assert.Equal(t, tc.expected, isTypeSupported(typ))
		})
	}
}
