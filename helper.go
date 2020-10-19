package konfig

import (
	"errors"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// flagValue implements the flag.Value interface.
type flagValue struct{}

func (v flagValue) String() string {
	return ""
}

func (v flagValue) Set(string) error {
	return nil
}

// tokenize breaks a field name into its tokens (generally words).
//   UserID       -->  User, ID
//   DatabaseURL  -->  Database, URL
func tokenize(name string) []string {
	tokens := []string{}
	current := string(name[0])
	lastLower := unicode.IsLower(rune(name[0]))

	add := func(slice []string, str string) []string {
		if str == "" {
			return slice
		}
		return append(slice, str)
	}

	for i := 1; i < len(name); i++ {
		r := rune(name[i])

		if unicode.IsUpper(r) && lastLower {
			// The case is changing from lower to upper
			tokens = add(tokens, current)
			current = string(name[i])
		} else if unicode.IsLower(r) && !lastLower {
			// The case is changing from upper to lower
			l := len(current) - 1
			tokens = add(tokens, current[:l])
			current = current[l:] + string(name[i])
		} else {
			// Increment current token
			current += string(name[i])
		}

		lastLower = unicode.IsLower(r)
	}

	tokens = append(tokens, string(current))

	return tokens
}

// getFlagName returns a canonical flag name for a field.
//   UserID       -->  user.id
//   DatabaseURL  -->  database.url
func getFlagName(name string) string {
	parts := tokenize(name)
	result := strings.Join(parts, ".")
	result = strings.ToLower(result)

	return result
}

// getFlagName returns a canonical environment variable name for a field.
//   UserID       -->  USER_ID
//   DatabaseURL  -->  DATABASE_URL
func getEnvVarName(name string) string {
	parts := tokenize(name)
	result := strings.Join(parts, "_")
	result = strings.ToUpper(result)

	return result
}

// getFileVarName returns a canonical environment variable name for value file of a field.
//   UserID       -->  USER_ID_FILE
//   DatabaseURL  -->  DATABASE_URL_FILE
func getFileEnvVarName(name string) string {
	parts := tokenize(name)
	result := strings.Join(parts, "_")
	result = strings.ToUpper(result)
	result = result + "_FILE"

	return result
}

// getFlagValue returns the value set for a flag.
//   - The flag name can start with - or --
//   - The flag value can be separated by space or =
func getFlagValue(flagName string) string {
	flagRegex := regexp.MustCompile("-{1,2}" + flagName)
	genericRegex := regexp.MustCompile("^-{1,2}[A-Za-z].*")

	for i, arg := range os.Args {
		if flagRegex.MatchString(arg) {
			if s := strings.Index(arg, "="); s > 0 {
				return arg[s+1:]
			}

			if i+1 < len(os.Args) {
				val := os.Args[i+1]
				if !genericRegex.MatchString(val) {
					return val
				}
			}

			return "true"
		}
	}

	return ""
}

func validateStruct(s interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(s) // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(s)  // reflect.Type --> t.Name(), t.Kind(), t.NumField()

	// A pointer to a struct should be passed
	if t.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("a non-pointer type is passed")
	}

	// Navigate to the pointer value
	v = v.Elem()
	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("a non-struct type is passed")
	}

	return v, nil
}

func isTypeSupported(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Ptr, reflect.Slice:
		return isTypeSupported(t.Elem())
	case reflect.Struct:
		return (t.PkgPath() == "net/url" && t.Name() == "URL") ||
			(t.PkgPath() == "regexp" && t.Name() == "Regexp")
	}

	return false
}
