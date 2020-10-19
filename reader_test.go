package konfig

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/moorara/konfig/ptr"
	"github.com/stretchr/testify/assert"
)

func TestReaderFromEnv(t *testing.T) {
	tests := []struct {
		name           string
		env            map[string]string
		expectedReader *reader
	}{
		{
			name: "NoOption",
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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
			expectedReader: &reader{
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

			assert.Equal(t, tc.expectedReader, readerFromEnv())
		})
	}
}

func TestReaderString(t *testing.T) {
	tests := []struct {
		name           string
		r              *reader
		expectedString string
	}{
		{
			"NoOption",
			&reader{},
			"",
		},
		{
			"WithDebug",
			&reader{
				debug: 2,
			},
			"Debug<2>",
		},
		{
			"WithListSep",
			&reader{
				listSep: "|",
			},
			"ListSep<|>",
		},
		{
			"WithPrefixFlag",
			&reader{
				prefixFlag: "config.",
			},
			"PrefixFlag<config.>",
		},
		{
			"WithPrefixEnv",
			&reader{
				prefixEnv: "CONFIG_",
			},
			"PrefixEnv<CONFIG_>",
		},
		{
			"WithprefixFileEnv",
			&reader{
				prefixFileEnv: "CONFIG_",
			},
			"PrefixFileEnv<CONFIG_>",
		},
		{
			"WithSkipFlag",
			&reader{
				skipFlag: true,
			},
			"SkipFlag",
		},
		{
			"WithSkipEnv",
			&reader{
				skipEnv: true,
			},
			"SkipEnv",
		},
		{
			"WithSkipFileEnv",
			&reader{
				skipFileEnv: true,
			},
			"SkipFileEnv",
		},
		{
			"WithTelepresence",
			&reader{
				telepresence: true,
			},
			"Telepresence",
		},
		{
			"WithSubscribers",
			&reader{
				subscribers: []chan Update{
					make(chan Update),
					make(chan Update),
				},
			},
			"Subscribers<2>",
		},
		{
			"WithAll",
			&reader{
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
			assert.Equal(t, tc.expectedString, tc.r.String())
		})
	}
}

func TestReaderLog(t *testing.T) {
	tests := []struct {
		name      string
		r         *reader
		verbosity uint
		msg       string
		args      []interface{}
	}{
		{
			"WithoutDebug",
			&reader{},
			1,
			"testing ...",
			nil,
		},
		{
			"WithDebug",
			&reader{
				debug: 2,
			},
			2,
			"testing ...",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.r.log(tc.verbosity, tc.msg, tc.args...)
		})
	}
}

func TestReaderGetFieldValue(t *testing.T) {
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
		r                                         *reader
		expectedValue                             string
		expectFilePath                            bool
	}{
		{
			"SkipFlag",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{},
			"info",
			false,
		},
		{
			"SkipFlagAndEnv",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "LOG_LEVEL_FILE",
			&reader{},
			"error",
			true,
		},
		{
			"SkipFlagAndEnvAndFile",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "-", "-", "-",
			&reader{},
			"",
			false,
		},
		{
			"SkipAllFlags",
			[]string{"/path/to/executable", "-log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{
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
			&reader{
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
			&reader{
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
			&reader{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level=debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "-log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{},
			"debug",
			false,
		},
		{
			"FromFlag",
			[]string{"/path/to/executable", "--log.level", "debug"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{},
			"debug",
			false,
		},
		{
			"FromEnvVar",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", "info"},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{},
			"info",
			false,
		},
		{
			"FromFile",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "error"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{},
			"error",
			true,
		},
		{
			"FromFileWithTelepresenceOption",
			[]string{"/path/to/executable"},
			env{"LOG_LEVEL", ""},
			file{"LOG_LEVEL_FILE", "info"},
			"Field", "log.level", "LOG_LEVEL", "LOG_LEVEL_FILE",
			&reader{telepresence: true},
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
			if tc.r.telepresence {
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
			value, filePath := tc.r.getFieldValue(tc.fieldName, tc.flagName, tc.envName, tc.fileEnvName)
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
		r              *reader
		fieldName      string
		fieldValue     interface{}
		expectedUpdate Update
	}{
		{
			"Nil",
			&reader{},
			"FieldBool", true,
			Update{},
		},
		{
			"NoChannel",
			&reader{
				subscribers: []chan Update{},
			},
			"FieldString", "value",
			Update{},
		},
		{
			"WithBlockingChannels",
			&reader{
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
			&reader{
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
			tc.r.notifySubscribers(tc.fieldName, tc.fieldValue)

			if tc.expectedUpdate != (Update{}) {
				for _, ch := range tc.r.subscribers {
					update := <-ch
					assert.Equal(t, tc.expectedUpdate, update)
				}
			}
		})
	}
}

func TestIterateOnFields(t *testing.T) {
	type fields struct {
		String        string
		Int           int
		StringPointer *string
		IntPointer    *int
		StringSlice   []string
		IntSlice      []int
	}

	tests := []struct {
		name                 string
		r                    *reader
		s                    interface{}
		expectedFieldNames   []string
		expectedFlagNames    []string
		expectedEnvNames     []string
		expectedFileEnvNames []string
		expectedListSeps     []string
	}{
		{
			name: "Default",
			r: &reader{
				listSep: ",",
			},
			s:                    &fields{},
			expectedFieldNames:   []string{"String", "Int", "StringPointer", "IntPointer", "StringSlice", "IntSlice"},
			expectedFlagNames:    []string{"string", "int", "string.pointer", "int.pointer", "string.slice", "int.slice"},
			expectedEnvNames:     []string{"STRING", "INT", "STRING_POINTER", "INT_POINTER", "STRING_SLICE", "INT_SLICE"},
			expectedFileEnvNames: []string{"STRING_FILE", "INT_FILE", "STRING_POINTER_FILE", "INT_POINTER_FILE", "STRING_SLICE_FILE", "INT_SLICE_FILE"},
			expectedListSeps:     []string{",", ",", ",", ",", ",", ","},
		},
		{
			name: "WithOptions",
			r: &reader{
				listSep:       "|",
				prefixFlag:    "config.",
				prefixEnv:     "CONFIG_",
				prefixFileEnv: "CONFIG_",
			},
			s:                    &fields{},
			expectedFieldNames:   []string{"String", "Int", "StringPointer", "IntPointer", "StringSlice", "IntSlice"},
			expectedFlagNames:    []string{"config.string", "config.int", "config.string.pointer", "config.int.pointer", "config.string.slice", "config.int.slice"},
			expectedEnvNames:     []string{"CONFIG_STRING", "CONFIG_INT", "CONFIG_STRING_POINTER", "CONFIG_INT_POINTER", "CONFIG_STRING_SLICE", "CONFIG_INT_SLICE"},
			expectedFileEnvNames: []string{"CONFIG_STRING_FILE", "CONFIG_INT_FILE", "CONFIG_STRING_POINTER_FILE", "CONFIG_INT_POINTER_FILE", "CONFIG_STRING_SLICE_FILE", "CONFIG_INT_SLICE_FILE"},
			expectedListSeps:     []string{"|", "|", "|", "|", "|", "|"},
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

			vStruct, err := validateStruct(tc.s)
			assert.NoError(t, err)

			tc.r.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
				fieldNames = append(fieldNames, fieldName)
				flagNames = append(flagNames, flagName)
				envNames = append(envNames, envName)
				fileEnvNames = append(fileEnvNames, fileEnvName)
				listSeps = append(listSeps, listSep)
			})

			assert.Equal(t, tc.expectedFieldNames, fieldNames)
			assert.Equal(t, tc.expectedFlagNames, flagNames)
			assert.Equal(t, tc.expectedEnvNames, envNames)
			assert.Equal(t, tc.expectedFileEnvNames, fileEnvNames)
			assert.Equal(t, tc.expectedListSeps, listSeps)
		})
	}
}

func TestRegisterFlags(t *testing.T) {
	type fields struct {
		String        string
		Int           int
		StringPointer *string
		IntPointer    *int
		StringSlice   []string
		IntSlice      []int
	}

	tests := []struct {
		name          string
		r             *reader
		s             interface{}
		expectedError error
		expectedFlags []string
	}{
		{
			name:          "Default",
			r:             &reader{},
			s:             &fields{},
			expectedError: nil,
			expectedFlags: []string{"string", "int", "string.pointer", "int.pointer", "string.slice", "int.slice"},
		},
		{
			name: "WithPrefixFlagOption",
			r: &reader{
				prefixFlag: "config.",
			},
			s:             &fields{},
			expectedError: nil,
			expectedFlags: []string{"config.string", "config.int", "config.string.pointer", "config.int.pointer", "config.string.slice", "config.int.slice"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vStruct, err := validateStruct(tc.s)
			assert.NoError(t, err)

			tc.r.registerFlags(vStruct)

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

	type fields struct {
		String        string
		Int           int
		StringPointer *string
		IntPointer    *int
		StringSlice   []string
		IntSlice      []int
	}

	tests := []struct {
		name     string
		args     []string
		envs     []env
		files    []file
		r        *reader
		s        interface{}
		expected interface{}
	}{
		{
			"Empty",
			[]string{"app"},
			[]env{},
			[]file{},
			&reader{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&fields{},
			&fields{},
		},
		{
			"AllFromDefaults",
			[]string{"app"},
			[]env{},
			[]file{},
			&reader{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&fields{
				String:        "content",
				Int:           -9223372036854775808,
				StringPointer: ptr.String("content"),
				IntPointer:    ptr.Int(-9223372036854775808),
				StringSlice:   []string{"content"},
				IntSlice:      []int{-9223372036854775808},
			},
			&fields{
				String:        "content",
				Int:           -9223372036854775808,
				StringPointer: ptr.String("content"),
				IntPointer:    ptr.Int(-9223372036854775808),
				StringSlice:   []string{"content"},
				IntSlice:      []int{-9223372036854775808},
			},
		},
		{
			"AllFromFlags",
			[]string{
				"-string=content",
				"-int=-9223372036854775808",
				"-string.pointer=content",
				"-int.pointer=-9223372036854775808",
				"-string.slice=content",
				"-int.slice=-9223372036854775808",
			},
			[]env{},
			[]file{},
			&reader{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&fields{},
			&fields{
				String:        "content",
				Int:           -9223372036854775808,
				StringPointer: ptr.String("content"),
				IntPointer:    ptr.Int(-9223372036854775808),
				StringSlice:   []string{"content"},
				IntSlice:      []int{-9223372036854775808},
			},
		},
		{
			"AllFromEnvVars",
			[]string{"app"},
			[]env{
				{"STRING", "content"},
				{"INT", "-9223372036854775808"},
				{"STRING_POINTER", "content"},
				{"INT_POINTER", "-9223372036854775808"},
				{"STRING_SLICE", "content"},
				{"INT_SLICE", "-9223372036854775808"},
			},
			[]file{},
			&reader{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&fields{},
			&fields{
				String:        "content",
				Int:           -9223372036854775808,
				StringPointer: ptr.String("content"),
				IntPointer:    ptr.Int(-9223372036854775808),
				StringSlice:   []string{"content"},
				IntSlice:      []int{-9223372036854775808},
			},
		},
		{
			"AllFromFromFiles",
			[]string{"app"},
			[]env{},
			[]file{
				{"STRING_FILE", "content"},
				{"INT_FILE", "-9223372036854775808"},
				{"STRING_POINTER_FILE", "content"},
				{"INT_POINTER_FILE", "-9223372036854775808"},
				{"STRING_SLICE_FILE", "content"},
				{"INT_SLICE_FILE", "-9223372036854775808"},
			},
			&reader{
				listSep:       ",",
				filesToFields: map[string]fieldInfo{},
			},
			&fields{},
			&fields{
				String:        "content",
				Int:           -9223372036854775808,
				StringPointer: ptr.String("content"),
				IntPointer:    ptr.Int(-9223372036854775808),
				StringSlice:   []string{"content"},
				IntSlice:      []int{-9223372036854775808},
			},
		},
		{
			"WithTelepresenceOption",
			[]string{"app"},
			[]env{},
			[]file{
				{"STRING_FILE", "content"},
				{"INT_FILE", "-9223372036854775808"},
				{"STRING_POINTER_FILE", "content"},
				{"INT_POINTER_FILE", "-9223372036854775808"},
				{"STRING_SLICE_FILE", "content"},
				{"INT_SLICE_FILE", "-9223372036854775808"},
			},
			&reader{
				listSep:       ",",
				telepresence:  true,
				filesToFields: map[string]fieldInfo{},
			},
			&fields{},
			&fields{
				String:        "content",
				Int:           -9223372036854775808,
				StringPointer: ptr.String("content"),
				IntPointer:    ptr.Int(-9223372036854775808),
				StringSlice:   []string{"content"},
				IntSlice:      []int{-9223372036854775808},
			},
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
			if tc.r.telepresence {
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

			vStruct, err := validateStruct(tc.s)
			assert.NoError(t, err)

			tc.r.readFields(vStruct)
			assert.Equal(t, tc.expected, tc.s)
		})
	}
}
