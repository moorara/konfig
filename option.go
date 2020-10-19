package konfig

// Option sets optional parameters for reader.
type Option func(*reader)

// Debug is the option for enabling logs for debugging purposes.
// verbosity is the verbosity level of logs.
// You can also enable this option by setting KONFIG_DEBUG environment variable to a verbosity level.
// You should not use this option in production.
func Debug(verbosity uint) Option {
	return func(c *reader) {
		c.debug = verbosity
	}
}

// ListSep is the option for specifying list separator for all fields with slice type.
// You can specify a list separator for each field using `sep` struct tag.
// Using `tag` struct tag for a field will override this option for that field.
func ListSep(sep string) Option {
	return func(c *reader) {
		c.listSep = sep
	}
}

// SkipFlag is the option for skipping command-line flags as a source for all fields.
// You can skip command-line flag as a source for each field by setting `flag` struct tag to `-`.
func SkipFlag() Option {
	return func(c *reader) {
		c.skipFlag = true
	}
}

// SkipEnv is the option for skipping environment variables as a source for all fields.
// You can skip environment variables as a source for each field by setting `env` struct tag to `-`.
func SkipEnv() Option {
	return func(c *reader) {
		c.skipEnv = true
	}
}

// SkipFileEnv is the option for skipping file environment variables as a source for all fields.
// You can skip file environment variable as a source for each field by setting `fileenv` struct tag to `-`.
func SkipFileEnv() Option {
	return func(c *reader) {
		c.skipFileEnv = true
	}
}

// PrefixFlag is the option for prefixing all flag names with a given string.
// You can specify a custom name for command-line flag for each field using `flag` struct tag.
// Using `flag` struct tag for a field will override this option for that field.
func PrefixFlag(prefix string) Option {
	return func(c *reader) {
		c.prefixFlag = prefix
	}
}

// PrefixEnv is the option for prefixing all environment variable names with a given string.
// You can specify a custom name for environment variable for each field using `env` struct tag.
// Using `env` struct tag for a field will override this option for that field.
func PrefixEnv(prefix string) Option {
	return func(c *reader) {
		c.prefixEnv = prefix
	}
}

// PrefixFileEnv is the option for prefixing all file environment variable names with a given string.
// You can specify a custom name for file environment variable for each field using `fileenv` struct tag.
// Using `fileenv` struct tag for a field will override this option for that field.
func PrefixFileEnv(prefix string) Option {
	return func(c *reader) {
		c.prefixFileEnv = prefix
	}
}

// Telepresence is the option for reading files when running in a Telepresence shell.
// If the TELEPRESENCE_ROOT environment variable exist, files will be read from mounted volume.
// See https://telepresence.io/howto/volumes.html for details.
func Telepresence() Option {
	return func(c *reader) {
		c.telepresence = true
	}
}
