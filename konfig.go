// Package konfig is a minimal and unopinionated library for reading configuration values in Go applications
// based on The 12-Factor App (https://12factor.net/config).
package konfig

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	skip       = "-"
	tagFlag    = "flag"
	tagEnv     = "env"
	tagFileEnv = "fileenv"
	tagSep     = "sep"

	envDebug            = "KONFIG_DEBUG"
	envListSep          = "KONFIG_LIST_SEP"
	envSkipFlag         = "KONFIG_SKIP_FLAG"
	envSkipEnv          = "KONFIG_SKIP_ENV"
	envSkipFileEnv      = "KONFIG_SKIP_FILE_ENV"
	envPrefixFlag       = "KONFIG_PREFIX_FLAG"
	envPrefixEnv        = "KONFIG_PREFIX_ENV"
	envPrefixFileEnv    = "KONFIG_PREFIX_FILE_ENV"
	envTelepresence     = "KONFIG_TELEPRESENCE"
	envTelepresenceRoot = "TELEPRESENCE_ROOT"

	line = "----------------------------------------------------------------------------------------------------"
)

// Update represents a configuration field that received a new value.
type Update struct {
	Name  string
	Value interface{}
}

// fieldInfo has all the information for setting a struct field later.
type fieldInfo struct {
	value   reflect.Value
	name    string
	listSep string
}

// controller controls how configuration values are read.
type controller struct {
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

// controllerFromEnv creates a new controller with defaults and with options read from environment variables.
func controllerFromEnv() *controller {
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

	return &controller{
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

// Option sets optional parameters for controller.
type Option func(*controller)

// Debug is the option for enabling logs for debugging purposes.
// verbosity is the verbosity level of logs.
// You can also enable this option by setting KONFIG_DEBUG environment variable to a verbosity level.
// You should not use this option in production.
func Debug(verbosity uint) Option {
	return func(c *controller) {
		c.debug = verbosity
	}
}

// ListSep is the option for specifying list separator for all fields with slice type.
// You can specify a list separator for each field using `sep` struct tag.
// Using `tag` struct tag for a field will override this option for that field.
func ListSep(sep string) Option {
	return func(c *controller) {
		c.listSep = sep
	}
}

// SkipFlag is the option for skipping command-line flags as a source for all fields.
// You can skip command-line flag as a source for each field by setting `flag` struct tag to `-`.
func SkipFlag() Option {
	return func(c *controller) {
		c.skipFlag = true
	}
}

// SkipEnv is the option for skipping environment variables as a source for all fields.
// You can skip environment variables as a source for each field by setting `env` struct tag to `-`.
func SkipEnv() Option {
	return func(c *controller) {
		c.skipEnv = true
	}
}

// SkipFileEnv is the option for skipping file environment variables as a source for all fields.
// You can skip file environment variable as a source for each field by setting `fileenv` struct tag to `-`.
func SkipFileEnv() Option {
	return func(c *controller) {
		c.skipFileEnv = true
	}
}

// PrefixFlag is the option for prefixing all flag names with a given string.
// You can specify a custom name for command-line flag for each field using `flag` struct tag.
// Using `flag` struct tag for a field will override this option for that field.
func PrefixFlag(prefix string) Option {
	return func(c *controller) {
		c.prefixFlag = prefix
	}
}

// PrefixEnv is the option for prefixing all environment variable names with a given string.
// You can specify a custom name for environment variable for each field using `env` struct tag.
// Using `env` struct tag for a field will override this option for that field.
func PrefixEnv(prefix string) Option {
	return func(c *controller) {
		c.prefixEnv = prefix
	}
}

// PrefixFileEnv is the option for prefixing all file environment variable names with a given string.
// You can specify a custom name for file environment variable for each field using `fileenv` struct tag.
// Using `fileenv` struct tag for a field will override this option for that field.
func PrefixFileEnv(prefix string) Option {
	return func(c *controller) {
		c.prefixFileEnv = prefix
	}
}

// Telepresence is the option for reading files when running in a Telepresence shell.
// If the TELEPRESENCE_ROOT environment variable exist, files will be read from mounted volume.
// See https://telepresence.io/howto/volumes.html for details.
func Telepresence() Option {
	return func(c *controller) {
		c.telepresence = true
	}
}

// String is used for printing debugging information.
// The output should fit in one line.
func (c *controller) String() string {
	strs := []string{}

	if c.debug > 0 {
		strs = append(strs, fmt.Sprintf("Debug<%d>", c.debug))
	}

	if c.listSep != "" {
		strs = append(strs, fmt.Sprintf("ListSep<%s>", c.listSep))
	}

	if c.skipFlag {
		strs = append(strs, "SkipFlag")
	}

	if c.skipEnv {
		strs = append(strs, "SkipEnv")
	}

	if c.skipFileEnv {
		strs = append(strs, "SkipFileEnv")
	}

	if c.prefixFlag != "" {
		strs = append(strs, fmt.Sprintf("PrefixFlag<%s>", c.prefixFlag))
	}

	if c.prefixEnv != "" {
		strs = append(strs, fmt.Sprintf("PrefixEnv<%s>", c.prefixEnv))
	}

	if c.prefixFileEnv != "" {
		strs = append(strs, fmt.Sprintf("PrefixFileEnv<%s>", c.prefixFileEnv))
	}

	if c.telepresence {
		strs = append(strs, "Telepresence")
	}

	if len(c.subscribers) > 0 {
		strs = append(strs, fmt.Sprintf("Subscribers<%d>", len(c.subscribers)))
	}

	return strings.Join(strs, " + ")
}

func (c *controller) log(v uint, msg string, args ...interface{}) {
	if v <= c.debug {
		log.Printf(msg+"\n", args...)
	}
}

// getFieldValue reads and returns the string value for a field from either
//   - command-line flags,
//   - environment variables,
//   - or configuration files
// If the value is read from a file, the second returned value will be the file path.
func (c *controller) getFieldValue(fieldName, flagName, envName, fileEnvName string) (string, string) {
	var value, filePath string

	// First, try reading from flag
	if value == "" && flagName != skip && !c.skipFlag {
		value = getFlagValue(flagName)
		c.log(5, "[%s] value read from flag %s: %s", fieldName, flagName, value)
	}

	// Second, try reading from environment variable
	if value == "" && envName != skip && !c.skipEnv {
		value = os.Getenv(envName)
		c.log(5, "[%s] value read from environment variable %s: %s", fieldName, envName, value)
	}

	// Third, try reading from file
	if value == "" && fileEnvName != skip && !c.skipFileEnv {
		// Read file environment variable
		filePath = os.Getenv(fileEnvName)
		c.log(5, "[%s] value read from file environment variable %s: %s", fieldName, fileEnvName, filePath)

		if filePath != "" {
			// Check for Telepresence
			// See https://telepresence.io/howto/volumes.html for details
			if c.telepresence {
				if mountPath := os.Getenv(envTelepresenceRoot); mountPath != "" {
					filePath = filepath.Join(mountPath, filePath)
					c.log(5, "[%s] telepresence mount path: %s", fieldName, mountPath)
				}
			}

			// Read config file
			filePath = filepath.Clean(filePath)
			if b, err := ioutil.ReadFile(filePath); err == nil {
				value = string(b)
				c.log(5, "[%s] value read from %s: %s", fieldName, filePath, value)
			}
		}
	}

	return value, filePath
}

// notifySubscribers sends an update to every subscriber channel in a new go routine.
func (c *controller) notifySubscribers(name string, value interface{}) {
	if len(c.subscribers) == 0 {
		return
	}

	c.log(4, "[%s] notifying %d subscribers ...", name, len(c.subscribers))

	update := Update{
		Name:  name,
		Value: value,
	}

	for i, sub := range c.subscribers {
		go func(id int, ch chan Update) {
			c.log(4, "[%s] notifying subscriber %d ...", name, id)
			ch <- update
			c.log(4, "[%s] subscriber %d notified", name, id)
		}(i, sub)
	}
}

func (c *controller) setString(v reflect.Value, name, val string) bool {
	if v.String() != val {
		c.log(5, "[%s] setting string value: %s", name, val)
		v.SetString(val)
		c.notifySubscribers(name, val)
		return true
	}

	return false
}

func (c *controller) setBool(v reflect.Value, name, val string) bool {
	if b, err := strconv.ParseBool(val); err == nil {
		if v.Bool() != b {
			c.log(5, "[%s] setting boolean value: %t", name, b)
			v.SetBool(b)
			c.notifySubscribers(name, b)
			return true
		}
	}

	return false
}

func (c *controller) setFloat32(v reflect.Value, name, val string) bool {
	if f, err := strconv.ParseFloat(val, 32); err == nil {
		if v.Float() != f {
			c.log(5, "[%s] setting float value: %f", name, f)
			v.SetFloat(f)
			c.notifySubscribers(name, float32(f))
			return true
		}
	}

	return false
}

func (c *controller) setFloat64(v reflect.Value, name, val string) bool {
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		if v.Float() != f {
			c.log(5, "[%s] setting float value: %f", name, f)
			v.SetFloat(f)
			c.notifySubscribers(name, f)
			return true
		}
	}

	return false
}

func (c *controller) setInt(v reflect.Value, name, val string) bool {
	// int size and range are platform-dependent
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			c.log(5, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt8(v reflect.Value, name, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 8); err == nil {
		if v.Int() != i {
			c.log(5, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int8(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt16(v reflect.Value, name, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 16); err == nil {
		if v.Int() != i {
			c.log(5, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int16(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt32(v reflect.Value, name, val string) bool {
	if i, err := strconv.ParseInt(val, 10, 32); err == nil {
		if v.Int() != i {
			c.log(5, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, int32(i))
			return true
		}
	}

	return false
}

func (c *controller) setInt64(v reflect.Value, name, val string) bool {
	if t := v.Type(); t.PkgPath() == "time" && t.Name() == "Duration" {
		// time.Duration
		if d, err := time.ParseDuration(val); err == nil {
			if v.Interface() != d {
				c.log(5, "[%s] setting duration value: %s", name, d)
				v.Set(reflect.ValueOf(d))
				c.notifySubscribers(name, d)
				return true
			}
		}
	} else if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		if v.Int() != i {
			c.log(5, "[%s] setting integer value: %d", name, i)
			v.SetInt(i)
			c.notifySubscribers(name, i)
			return true
		}
	}

	return false
}

func (c *controller) setUint(v reflect.Value, name, val string) bool {
	// uint size and range are platform-dependent
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			c.log(5, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint8(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 8); err == nil {
		if v.Uint() != u {
			c.log(5, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint8(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint16(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 16); err == nil {
		if v.Uint() != u {
			c.log(5, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint16(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint32(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 32); err == nil {
		if v.Uint() != u {
			c.log(5, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, uint32(u))
			return true
		}
	}

	return false
}

func (c *controller) setUint64(v reflect.Value, name, val string) bool {
	if u, err := strconv.ParseUint(val, 10, 64); err == nil {
		if v.Uint() != u {
			c.log(5, "[%s] setting unsigned integer value: %d", name, u)
			v.SetUint(u)
			c.notifySubscribers(name, u)
			return true
		}
	}

	return false
}

func (c *controller) setStruct(v reflect.Value, name, val string) bool {
	if t := v.Type(); t.PkgPath() == "net/url" && t.Name() == "URL" {
		// url.URL
		if u, err := url.Parse(val); err == nil {
			// u is a pointer
			if !reflect.DeepEqual(v.Interface(), *u) {
				c.log(5, "[%s] setting url value: %s", name, val)
				v.Set(reflect.ValueOf(u).Elem())
				c.notifySubscribers(name, *u)
				return true
			}
		}
	}

	return false
}

func (c *controller) setStringSlice(v reflect.Value, name string, vals []string) bool {
	if !reflect.DeepEqual(v.Interface(), vals) {
		c.log(5, "[%s] setting string slice: %v", name, vals)
		v.Set(reflect.ValueOf(vals))
		c.notifySubscribers(name, vals)
		return true
	}

	return false
}

func (c *controller) setBoolSlice(v reflect.Value, name string, vals []string) bool {
	bools := []bool{}
	for _, val := range vals {
		if b, err := strconv.ParseBool(val); err == nil {
			bools = append(bools, b)
		}
	}

	if !reflect.DeepEqual(v.Interface(), bools) {
		c.log(5, "[%s] setting boolean slice: %v", name, bools)
		v.Set(reflect.ValueOf(bools))
		c.notifySubscribers(name, bools)
		return true
	}

	return false
}

func (c *controller) setFloat32Slice(v reflect.Value, name string, vals []string) bool {
	floats := []float32{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 32); err == nil {
			floats = append(floats, float32(f))
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		c.log(5, "[%s] setting float32 slice: %v", name, floats)
		v.Set(reflect.ValueOf(floats))
		c.notifySubscribers(name, floats)
		return true
	}

	return false
}

func (c *controller) setFloat64Slice(v reflect.Value, name string, vals []string) bool {
	floats := []float64{}
	for _, val := range vals {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			floats = append(floats, f)
		}
	}

	if !reflect.DeepEqual(v.Interface(), floats) {
		c.log(5, "[%s] setting float64 slice: %v", name, floats)
		v.Set(reflect.ValueOf(floats))
		c.notifySubscribers(name, floats)
		return true
	}

	return false
}

func (c *controller) setIntSlice(v reflect.Value, name string, vals []string) bool {
	// int size and range are platform-dependent
	ints := []int{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			ints = append(ints, int(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(5, "[%s] setting int slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt8Slice(v reflect.Value, name string, vals []string) bool {
	ints := []int8{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 8); err == nil {
			ints = append(ints, int8(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(5, "[%s] setting int8 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt16Slice(v reflect.Value, name string, vals []string) bool {
	ints := []int16{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 16); err == nil {
			ints = append(ints, int16(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(5, "[%s] setting int16 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt32Slice(v reflect.Value, name string, vals []string) bool {
	ints := []int32{}
	for _, val := range vals {
		if i, err := strconv.ParseInt(val, 10, 32); err == nil {
			ints = append(ints, int32(i))
		}
	}

	if !reflect.DeepEqual(v.Interface(), ints) {
		c.log(5, "[%s] setting int32 slice: %v", name, ints)
		v.Set(reflect.ValueOf(ints))
		c.notifySubscribers(name, ints)
		return true
	}

	return false
}

func (c *controller) setInt64Slice(v reflect.Value, name string, vals []string) bool {
	if t := reflect.TypeOf(v.Interface()).Elem(); t.PkgPath() == "time" && t.Name() == "Duration" {
		durations := []time.Duration{}
		for _, val := range vals {
			if d, err := time.ParseDuration(val); err == nil {
				durations = append(durations, d)
			}
		}

		// []time.Duration
		if !reflect.DeepEqual(v.Interface(), durations) {
			c.log(5, "[%s] setting duration slice: %v", name, durations)
			v.Set(reflect.ValueOf(durations))
			c.notifySubscribers(name, durations)
			return true
		}
	} else {
		ints := []int64{}
		for _, val := range vals {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				ints = append(ints, i)
			}
		}

		if !reflect.DeepEqual(v.Interface(), ints) {
			c.log(5, "[%s] setting int64 slice: %v", name, ints)
			v.Set(reflect.ValueOf(ints))
			c.notifySubscribers(name, ints)
			return true
		}
	}

	return false
}

func (c *controller) setUintSlice(v reflect.Value, name string, vals []string) bool {
	// uint size and range are platform-dependent
	uints := []uint{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, uint(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(5, "[%s] setting uint slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint8Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint8{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 8); err == nil {
			uints = append(uints, uint8(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(5, "[%s] setting uint8 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint16Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint16{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 16); err == nil {
			uints = append(uints, uint16(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(5, "[%s] setting uint16 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint32Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint32{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 32); err == nil {
			uints = append(uints, uint32(u))
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(5, "[%s] setting uint32 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setUint64Slice(v reflect.Value, name string, vals []string) bool {
	uints := []uint64{}
	for _, val := range vals {
		if u, err := strconv.ParseUint(val, 10, 64); err == nil {
			uints = append(uints, u)
		}
	}

	if !reflect.DeepEqual(v.Interface(), uints) {
		c.log(5, "[%s] setting uint64 slice: %v", name, uints)
		v.Set(reflect.ValueOf(uints))
		c.notifySubscribers(name, uints)
		return true
	}

	return false
}

func (c *controller) setURLSlice(v reflect.Value, name string, vals []string) bool {
	t := reflect.TypeOf(v.Interface()).Elem()

	if t.PkgPath() == "net/url" && t.Name() == "URL" {
		urls := []url.URL{}
		for _, val := range vals {
			if u, err := url.Parse(val); err == nil {
				urls = append(urls, *u)
			}
		}

		// []url.URL
		if !reflect.DeepEqual(v.Interface(), urls) {
			c.log(5, "[%s] setting url slice: %v", name, urls)
			v.Set(reflect.ValueOf(urls))
			c.notifySubscribers(name, urls)
			return true
		}
	}

	return false
}

func (c *controller) setFieldValue(f fieldInfo, val string) bool {
	switch f.value.Kind() {
	case reflect.String:
		return c.setString(f.value, f.name, val)
	case reflect.Bool:
		return c.setBool(f.value, f.name, val)
	case reflect.Float32:
		return c.setFloat32(f.value, f.name, val)
	case reflect.Float64:
		return c.setFloat64(f.value, f.name, val)
	case reflect.Int:
		return c.setInt(f.value, f.name, val)
	case reflect.Int8:
		return c.setInt8(f.value, f.name, val)
	case reflect.Int16:
		return c.setInt16(f.value, f.name, val)
	case reflect.Int32:
		return c.setInt32(f.value, f.name, val)
	case reflect.Int64:
		return c.setInt64(f.value, f.name, val)
	case reflect.Uint:
		return c.setUint(f.value, f.name, val)
	case reflect.Uint8:
		return c.setUint8(f.value, f.name, val)
	case reflect.Uint16:
		return c.setUint16(f.value, f.name, val)
	case reflect.Uint32:
		return c.setUint32(f.value, f.name, val)
	case reflect.Uint64:
		return c.setUint64(f.value, f.name, val)
	case reflect.Struct:
		return c.setStruct(f.value, f.name, val)

	case reflect.Slice:
		tSlice := reflect.TypeOf(f.value.Interface()).Elem()
		vals := strings.Split(val, f.listSep)

		switch tSlice.Kind() {
		case reflect.String:
			return c.setStringSlice(f.value, f.name, vals)
		case reflect.Bool:
			return c.setBoolSlice(f.value, f.name, vals)
		case reflect.Float32:
			return c.setFloat32Slice(f.value, f.name, vals)
		case reflect.Float64:
			return c.setFloat64Slice(f.value, f.name, vals)
		case reflect.Int:
			return c.setIntSlice(f.value, f.name, vals)
		case reflect.Int8:
			return c.setInt8Slice(f.value, f.name, vals)
		case reflect.Int16:
			return c.setInt16Slice(f.value, f.name, vals)
		case reflect.Int32:
			return c.setInt32Slice(f.value, f.name, vals)
		case reflect.Int64:
			return c.setInt64Slice(f.value, f.name, vals)
		case reflect.Uint:
			return c.setUintSlice(f.value, f.name, vals)
		case reflect.Uint8:
			return c.setUint8Slice(f.value, f.name, vals)
		case reflect.Uint16:
			return c.setUint16Slice(f.value, f.name, vals)
		case reflect.Uint32:
			return c.setUint32Slice(f.value, f.name, vals)
		case reflect.Uint64:
			return c.setUint64Slice(f.value, f.name, vals)
		case reflect.Struct:
			return c.setURLSlice(f.value, f.name, vals)
		}
	}

	return false
}

func (c *controller) iterateOnFields(vStruct reflect.Value, handle func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string)) {
	// Iterate over struct fields
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)        // reflect.Value --> vField.Kind(), vField.Type().Name(), vField.Type().Kind(), vField.Interface()
		f := vStruct.Type().Field(i) // reflect.StructField --> tField.Name, tField.Type.Name(), tField.Type.Kind(), tField.Tag.Get(tag)

		// Skip unexported and unsupported fields
		if !v.CanSet() || !isTypeSupported(v.Type()) {
			continue
		}

		// `flag:"..."`
		flagName := f.Tag.Get(tagFlag)
		if flagName == "" {
			flagName = c.prefixFlag + getFlagName(f.Name)
		}

		// `env:"..."`
		envName := f.Tag.Get(tagEnv)
		if envName == "" {
			envName = c.prefixEnv + getEnvVarName(f.Name)
		}

		// `fileenv:"..."`
		fileEnvName := f.Tag.Get(tagFileEnv)
		if fileEnvName == "" {
			fileEnvName = c.prefixFileEnv + getFileEnvVarName(f.Name)
		}

		// `sep:"..."`
		listSep := f.Tag.Get(tagSep)
		if listSep == "" {
			listSep = c.listSep
		}

		handle(v, f.Name, flagName, envName, fileEnvName, listSep)
	}
}

func (c *controller) registerFlags(vStruct reflect.Value) {
	c.log(2, "Registering configuration flags ...")
	c.log(2, line)

	c.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
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

		c.log(5, "[%s] flag registered: %s", fieldName, flagName)
	})

	c.log(5, line)
}

func (c *controller) readFields(vStruct reflect.Value) {
	c.log(2, "Reading configuration values ...")
	c.log(2, line)

	c.iterateOnFields(vStruct, func(v reflect.Value, fieldName, flagName, envName, fileEnvName, listSep string) {
		c.log(5, "[%s] expecting flag name: %s", fieldName, flagName)
		c.log(5, "[%s] expecting environment variable name: %s", fieldName, envName)
		c.log(5, "[%s] expecting file environment variable name: %s", fieldName, fileEnvName)
		c.log(5, "[%s] expecting list separator: %s", fieldName, listSep)
		defer c.log(5, line)

		// Try reading the configuration value for current field
		val, path := c.getFieldValue(fieldName, flagName, envName, fileEnvName)

		// If no value, skip this field
		if val == "" {
			c.log(5, "[%s] falling back to default value: %v", fieldName, v.Interface())
			return
		}

		f := fieldInfo{
			value:   v,
			name:    fieldName,
			listSep: listSep,
		}

		// Keep the track of which fields are read from which files
		if path != "" {
			c.filesToFields[path] = f
		}

		c.setFieldValue(f, val)
	})
}

// Pick reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// Default values can also be specified.
// You should pass the pointer to a struct for config; otherwise you will get an error.
func Pick(config interface{}, opts ...Option) error {
	c := controllerFromEnv()
	for _, opt := range opts {
		opt(c)
	}

	c.log(2, line)
	c.log(2, "Options: %s", c)
	c.log(2, line)

	v, err := validateStruct(config)
	if err != nil {
		c.log(1, err.Error())
		return err
	}

	c.registerFlags(v)
	c.readFields(v)

	return nil
}

// Watch first reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// It then watches any change to those fields that their values are read from configuration files and notifies subscribers on a channel.
func Watch(config sync.Locker, subscribers []chan Update, opts ...Option) (func(), error) {
	c := controllerFromEnv()
	c.subscribers = subscribers
	for _, opt := range opts {
		opt(c)
	}

	c.log(2, line)
	c.log(2, "Options: %s", c)
	c.log(2, line)

	v, err := validateStruct(config)
	if err != nil {
		c.log(1, err.Error())
		return nil, err
	}

	c.registerFlags(v)
	c.readFields(v)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		c.log(1, "cannot create a watcher: %s", err)
		return nil, err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					break
				}

				path := filepath.Clean(event.Name)
				c.log(6, "event received: %s %s", event.Op, event.Name)

				// We only receive events for added files, this if check is redundant!
				if f, ok := c.filesToFields[path]; ok {
					// Write
					if event.Op&fsnotify.Write == fsnotify.Write {
						if b, err := ioutil.ReadFile(path); err == nil {
							val := string(b)
							c.log(3, "received an update from %s: %s", path, val)
							config.Lock()
							c.setFieldValue(f, val)
							config.Unlock()
						}
					}

					// Remove
					// Kubernetes injects new values from ConfigMaps and Secrets by removing the mounted files and recreating them.
					// When a watched file is removed, the fsnotify package will remove it from the watcher too.
					// This if block is a workaround for the aforementioned Kubernetes situation.
					// See https://github.com/moorara/konfig/issues/47
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						// Check if the removed file is already recreated
						if _, err := os.Stat(path); err == nil {
							if b, err := ioutil.ReadFile(path); err == nil {
								val := string(b)
								c.log(3, "received an update from %s: %s", path, val)
								config.Lock()
								c.setFieldValue(f, val)
								config.Unlock()
							}

							// Re-Add a watch for the file
							if err := watcher.Add(path); err != nil {
								c.log(1, "cannot watch file %s: %s", f, err)
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					break
				}
				c.log(1, "error watching: %s", err)
			}
		}
	}()

	for path := range c.filesToFields {
		if err := watcher.Add(path); err != nil {
			c.log(1, "cannot watch file %s: %s", path, err)
			return nil, err
		}
	}

	close := func() {
		watcher.Close()
		// TODO: closing subscriber channels causes data race if notifySubscribers is writing to any
		/* for _, sub := range c.subscribers {
			close(sub)
		} */
	}

	return close, nil
}
