package konfig

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// fieldInfo has all the information for setting a struct field later.
type fieldInfo struct {
	value   reflect.Value
	name    string
	listSep string
}

// reader controls how configuration values are read.
type reader struct {
	debug         uint
	listSep       string
	skipFlag      bool
	skipEnv       bool
	skipFileEnv   bool
	prefixFlag    string
	prefixEnv     string
	prefixFileEnv string
	telepresence  bool

	subscribers   []chan Update
	filesToFields map[string]fieldInfo
}

// readerFromEnv creates a new reader with defaults and with options read from environment variables.
func readerFromEnv() *reader {
	var debug uint
	if str := os.Getenv(envDebug); str != "" {
		// debug verbosity level should not be higher than 255 (8-bits)
		if u, err := strconv.ParseUint(str, 10, 8); err == nil {
			debug = uint(u)
		}
	}

	listSep := os.Getenv(envListSep)

	// Set the default list separator
	if listSep == "" {
		listSep = ","
	}

	var skipFlag bool
	if str := os.Getenv(envSkipFlag); str != "" {
		skipFlag, _ = strconv.ParseBool(str)
	}

	var skipEnv bool
	if str := os.Getenv(envSkipEnv); str != "" {
		skipEnv, _ = strconv.ParseBool(str)
	}

	var skipFileEnv bool
	if str := os.Getenv(envSkipFileEnv); str != "" {
		skipFileEnv, _ = strconv.ParseBool(str)
	}

	prefixFlag := os.Getenv(envPrefixFlag)
	prefixEnv := os.Getenv(envPrefixEnv)
	prefixFileEnv := os.Getenv(envPrefixFileEnv)

	var telepresence bool
	if str := os.Getenv(envTelepresence); str != "" {
		telepresence, _ = strconv.ParseBool(str)
	}

	return &reader{
		debug:         debug,
		listSep:       listSep,
		skipFlag:      skipFlag,
		skipEnv:       skipEnv,
		skipFileEnv:   skipFileEnv,
		prefixFlag:    prefixFlag,
		prefixEnv:     prefixEnv,
		prefixFileEnv: prefixFileEnv,
		telepresence:  telepresence,

		subscribers:   nil,
		filesToFields: map[string]fieldInfo{},
	}
}

// String is used for printing debugging information.
// The output should fit in one line.
func (r *reader) String() string {
	strs := []string{}

	if r.debug > 0 {
		strs = append(strs, fmt.Sprintf("Debug<%d>", r.debug))
	}

	if r.listSep != "" {
		strs = append(strs, fmt.Sprintf("ListSep<%s>", r.listSep))
	}

	if r.skipFlag {
		strs = append(strs, "SkipFlag")
	}

	if r.skipEnv {
		strs = append(strs, "SkipEnv")
	}

	if r.skipFileEnv {
		strs = append(strs, "SkipFileEnv")
	}

	if r.prefixFlag != "" {
		strs = append(strs, fmt.Sprintf("PrefixFlag<%s>", r.prefixFlag))
	}

	if r.prefixEnv != "" {
		strs = append(strs, fmt.Sprintf("PrefixEnv<%s>", r.prefixEnv))
	}

	if r.prefixFileEnv != "" {
		strs = append(strs, fmt.Sprintf("PrefixFileEnv<%s>", r.prefixFileEnv))
	}

	if r.telepresence {
		strs = append(strs, "Telepresence")
	}

	if len(r.subscribers) > 0 {
		strs = append(strs, fmt.Sprintf("Subscribers<%d>", len(r.subscribers)))
	}

	return strings.Join(strs, " + ")
}

func (r *reader) log(verbosity uint, msg string, args ...interface{}) {
	if verbosity <= r.debug {
		log.Printf(msg+"\n", args...)
	}
}

// getFieldValue reads and returns the string value for a field from either
//   - command-line flags,
//   - environment variables,
//   - or configuration files
// If the value is read from a file, the second returned value will be the file path.
func (r *reader) getFieldValue(fieldName, flagName, envName, fileEnvName string) (string, string) {
	var value, filePath string

	// First, try reading from flag
	if value == "" && flagName != skip && !r.skipFlag {
		value = getFlagValue(flagName)
		r.log(5, "[%s] value read from flag %s: %s", fieldName, flagName, value)
	}

	// Second, try reading from environment variable
	if value == "" && envName != skip && !r.skipEnv {
		value = os.Getenv(envName)
		r.log(5, "[%s] value read from environment variable %s: %s", fieldName, envName, value)
	}

	// Third, try reading from file
	if value == "" && fileEnvName != skip && !r.skipFileEnv {
		// Read file environment variable
		filePath = os.Getenv(fileEnvName)
		r.log(5, "[%s] value read from file environment variable %s: %s", fieldName, fileEnvName, filePath)

		if filePath != "" {
			// Check for Telepresence
			// See https://telepresence.io/howto/volumes.html for details
			if r.telepresence {
				if mountPath := os.Getenv(envTelepresenceRoot); mountPath != "" {
					filePath = filepath.Join(mountPath, filePath)
					r.log(5, "[%s] telepresence mount path: %s", fieldName, mountPath)
				}
			}

			// Read config file
			filePath = filepath.Clean(filePath)
			if b, err := ioutil.ReadFile(filePath); err == nil {
				value = string(b)
				r.log(5, "[%s] value read from %s: %s", fieldName, filePath, value)
			}
		}
	}

	return value, filePath
}

// notifySubscribers sends an update to every subscriber channel in a new go routine.
func (r *reader) notifySubscribers(name string, value interface{}) {
	if len(r.subscribers) == 0 {
		return
	}

	r.log(4, "[%s] notifying %d subscribers ...", name, len(r.subscribers))

	update := Update{
		Name:  name,
		Value: value,
	}

	for i, sub := range r.subscribers {
		go func(id int, ch chan Update) {
			r.log(4, "[%s] notifying subscriber %d ...", name, id)
			ch <- update
			r.log(4, "[%s] subscriber %d notified", name, id)
		}(i, sub)
	}
}

func (r *reader) iterateOnFields(vStruct reflect.Value, handle func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string)) {
	// Iterate over struct fields
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)        // reflect.Value       --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		t := v.Type()                // reflect.Type        --> t.Kind(), t.PkgPath(), t.Name(), t.NumField()
		f := vStruct.Type().Field(i) // reflect.StructField --> f.Name, f.Type.Name(), f.Type.Kind(), f.Tag.Get(tag)

		// Skip unexported and unsupported fields
		if !v.CanSet() || !isTypeSupported(t) {
			continue
		}

		// `flag:"..."`
		flagName := f.Tag.Get(tagFlag)
		if flagName == "" {
			flagName = r.prefixFlag + getFlagName(f.Name)
		}

		// `env:"..."`
		envName := f.Tag.Get(tagEnv)
		if envName == "" {
			envName = r.prefixEnv + getEnvVarName(f.Name)
		}

		// `fileenv:"..."`
		fileEnvName := f.Tag.Get(tagFileEnv)
		if fileEnvName == "" {
			fileEnvName = r.prefixFileEnv + getFileEnvVarName(f.Name)
		}

		// `sep:"..."`
		listSep := f.Tag.Get(tagSep)
		if listSep == "" {
			listSep = r.listSep
		}

		handle(v, f.Name, flagName, envName, fileEnvName, listSep)
	}
}

func (r *reader) registerFlags(vStruct reflect.Value) {
	r.log(2, "Registering configuration flags ...")
	r.log(2, line)

	r.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
		if flagName == skip {
			return
		}

		var dataType string
		if v.Kind() == reflect.Slice {
			dataType = fmt.Sprintf("[]%s", reflect.TypeOf(v.Interface()).Elem())
		} else {
			dataType = v.Type().String()
		}

		defaultValue := fmt.Sprintf("%v", v.Interface())

		usage := fmt.Sprintf(
			"%s:\t\t\t\t%s\n%s:\t\t\t\t%s\n%s:\t\t\t%s\n%s:\t%s",
			"data type", dataType,
			"default value", defaultValue,
			"environment variable", envName,
			"environment variable for file path", fileEnvName,
		)

		// Define a flag for the field, so flag.Parse() can be called
		if flag.Lookup(flagName) == nil {
			switch v.Kind() {
			case reflect.Bool:
				flag.Bool(flagName, v.Bool(), usage)
			default:
				flag.Var(&flagValue{}, flagName, usage)
			}
		}

		r.log(5, "[%s] flag registered: %s", fieldName, flagName)
	})

	r.log(5, line)
}

func (r *reader) readFields(vStruct reflect.Value) {
	r.log(2, "Reading configuration values ...")
	r.log(2, line)

	r.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
		r.log(5, "[%s] expecting flag name: %s", fieldName, flagName)
		r.log(5, "[%s] expecting environment variable name: %s", fieldName, envName)
		r.log(5, "[%s] expecting file environment variable name: %s", fieldName, fileEnvName)
		r.log(5, "[%s] expecting list separator: %s", fieldName, listSep)
		defer r.log(5, line)

		// Try reading the configuration value for current field
		val, path := r.getFieldValue(fieldName, flagName, envName, fileEnvName)

		// If no value, skip this field
		if val == "" {
			r.log(5, "[%s] falling back to default value: %v", fieldName, v.Interface())
			return
		}

		f := fieldInfo{
			value:   v,
			name:    fieldName,
			listSep: listSep,
		}

		// Keep the track of which fields are read from which files
		if path != "" {
			r.filesToFields[path] = f
		}

		r.setFieldValue(f, val)
	})
}
