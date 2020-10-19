package konfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	r := new(reader)
	Debug(2)(r)

	expected := &reader{
		debug: 2,
	}

	assert.Equal(t, expected, r)
}

func TestListSep(t *testing.T) {
	r := new(reader)
	ListSep("|")(r)

	expected := &reader{
		listSep: "|",
	}

	assert.Equal(t, expected, r)
}

func TestSkipFlag(t *testing.T) {
	r := new(reader)
	SkipFlag()(r)

	expected := &reader{
		skipFlag: true,
	}

	assert.Equal(t, expected, r)
}

func TestSkipEnv(t *testing.T) {
	r := new(reader)
	SkipEnv()(r)

	expected := &reader{
		skipEnv: true,
	}

	assert.Equal(t, expected, r)
}

func TestSkipFileEnv(t *testing.T) {
	r := new(reader)
	SkipFileEnv()(r)

	expected := &reader{
		skipFileEnv: true,
	}

	assert.Equal(t, expected, r)
}

func TestPrefixFlag(t *testing.T) {
	r := new(reader)
	PrefixFlag("config.")(r)

	expected := &reader{
		prefixFlag: "config.",
	}

	assert.Equal(t, expected, r)
}

func TestPrefixEnv(t *testing.T) {
	r := new(reader)
	PrefixEnv("CONFIG_")(r)

	expected := &reader{
		prefixEnv: "CONFIG_",
	}

	assert.Equal(t, expected, r)
}

func TestPrefixFileEnv(t *testing.T) {
	r := new(reader)
	PrefixFileEnv("CONFIG_")(r)

	expected := &reader{
		prefixFileEnv: "CONFIG_",
	}

	assert.Equal(t, expected, r)
}

func TestTelepresence(t *testing.T) {
	r := new(reader)
	Telepresence()(r)

	expected := &reader{
		telepresence: true,
	}

	assert.Equal(t, expected, r)
}
