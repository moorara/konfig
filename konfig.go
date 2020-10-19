// Package konfig is a minimal and unopinionated library for reading configuration values in Go applications
// based on The 12-Factor App (https://12factor.net/config).
package konfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

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

// Pick reads values for exported fields of a struct from either command-line flags, environment variables, or configuration files.
// Default values can also be specified.
// You should pass the pointer to a struct for config; otherwise you will get an error.
func Pick(config interface{}, opts ...Option) error {
	c := readerFromEnv()
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
	c := readerFromEnv()
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
