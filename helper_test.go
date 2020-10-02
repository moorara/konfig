package konfig

import (
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlagValue(t *testing.T) {
	tests := []struct {
		fv               *flagValue
		expectedString   string
		setString        string
		expectedSetError error
	}{
		{
			fv:               &flagValue{},
			expectedString:   "",
			setString:        "anything",
			expectedSetError: nil,
		},
	}

	for _, tc := range tests {
		str := tc.fv.String()
		assert.Equal(t, tc.expectedString, str)

		err := tc.fv.Set(tc.setString)
		assert.Equal(t, tc.expectedSetError, err)
	}
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
			config{},
			errors.New("a non-pointer type is passed"),
		},
		{
			"OK",
			&config{},
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
	service1URL, _ := url.Parse("service-1:8080")
	service2URL, _ := url.Parse("service-2:8080")

	tests := []struct {
		name     string
		field    interface{}
		expected bool
	}{
		{"String", "dummy", true},
		{"Bool", true, true},
		{"Float32", float32(3.1415), true},
		{"Float64", float64(3.14159265359), true},
		{"Int", int(-2147483648), true},
		{"Int8", int8(-128), true},
		{"Int16", int16(-32768), true},
		{"Int32", int32(-2147483648), true},
		{"Int64", int64(-9223372036854775808), true},
		{"Duration", time.Hour, true},
		{"Uint", uint(4294967295), true},
		{"Uint8", uint8(255), true},
		{"Uint16", uint16(65535), true},
		{"Uint32", uint32(4294967295), true},
		{"Uint64", uint64(18446744073709551615), true},
		{"URL", *service1URL, true},
		{"StringSlice", []string{"foo", "bar"}, true},
		{"BoolSlice", []bool{true, false}, true},
		{"Float32Slice", []float32{3.1415, 2.7182}, true},
		{"Float64Slice", []float64{3.14159265359, 2.71828182845}, true},
		{"IntSlice", []int{}, true},
		{"Int8Slice", []int8{}, true},
		{"Int16Slice", []int16{}, true},
		{"Int32Slice", []int32{}, true},
		{"Int64Slice", []int64{}, true},
		{"DurationSlice", []time.Duration{}, true},
		{"UintSlice", []uint{}, true},
		{"Uint8Slice", []uint8{}, true},
		{"Uint16Slice", []uint16{}, true},
		{"Uint32Slice", []uint32{}, true},
		{"Uint64Slice", []uint64{}, true},
		{"URLSlice", []url.URL{*service1URL, *service2URL}, true},
		{"Unsupported", time.Now(), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tp := reflect.TypeOf(tc.field)

			assert.Equal(t, tc.expected, isTypeSupported(tp))
		})
	}
}
