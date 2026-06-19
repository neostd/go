// Package env provides helpers for reading, writing, and expanding environment
// variables.
package env

import (
	"os"
	"strings"
)

// Get retrieves the value of the environment variable named by key.
// It returns an empty string if the variable is not present.
func Get(key string) string {
	return os.Getenv(key)
}

// Set the value of the environment variable named by key to value.
// It returns an error if the variable cannot be set.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset removes the environment variable named by key.
// It returns an error if the variable cannot be unset.
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Has reports whether the environment variable named by key exists.
func Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// All returns all environment variables as a key/value map.
func All() map[string]string {
	kv := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if len(pair) == 2 && len(pair[1]) > 0 {
			kv[pair[0]] = pair[1]
		}
	}

	return kv
}

// GetPath returns the current process path environment variable.
func GetPath() string {
	return Get(PATH)
}

// SetPath updates the current process path environment variable.
func SetPath(value string) error {
	return Set(PATH, value)
}

// SplitPath splits the current process path into path entries.
func SplitPath() []string {
	return strings.Split(GetPath(), string(os.PathListSeparator))
}

// JoinPath joins path entries using the OS-specific path list separator.
func JoinPath(paths ...string) string {
	return strings.Join(paths, string(os.PathListSeparator))
}

// PrependPath prepends a path entry if it is not already first.
func PrependPath(path string) error {
	paths := SplitPath()

	if matchPath(paths[0], path) {
		return nil
	}

	paths = append([]string{path}, paths...)
	return SetPath(JoinPath(paths...))
}

// AppendPath appends a path entry if it is not already last.
func AppendPath(path string) error {
	paths := SplitPath()

	if matchPath(paths[len(paths)-1], path) {
		return nil
	}

	paths = append(paths, path)
	return SetPath(JoinPath(paths...))
}

// HasPath reports whether a path entry exists in the current process path.
func HasPath(path string) bool {
	paths := SplitPath()
	return hasPath(path, paths)
}
