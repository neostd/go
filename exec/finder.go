package exec

import (
	"errors"
	"runtime"
	"unicode"

	"github.com/neostd/go/env"
)

// Executable describes an executable and optional platform-specific lookup hints.
type Executable struct {
	Name     string
	Path     string
	Variable string
	Windows  []string
	Linux    []string
	Darwin   []string
}

// ExecutableRegistry stores named executables and lookup metadata.
type ExecutableRegistry struct {
	data map[string]Executable
}

// EnvLike describes the environment behavior used during executable resolution.
type EnvLike interface {
	Get(key string) string
	Expand(s string) (string, error)
	Set(key, value string)
	SplitPath() []string
}

// SetEnvLike overrides environment access used during executable lookup.
func SetEnvLike(e EnvLike) {
	envLike = e
}

// GetEnvLike returns the environment adapter used during executable lookup.
func GetEnvLike() EnvLike {
	return envLike
}

var envLike EnvLike

type defaultEnvLike struct{}

func (d *defaultEnvLike) Get(key string) string {
	return env.Get(key)
}

func (d *defaultEnvLike) Expand(s string) (string, error) {
	return env.Expand(s)
}

func (d *defaultEnvLike) Set(key, value string) {
	_ = env.Set(key, value)
}

func (d *defaultEnvLike) SplitPath() []string {
	return env.SplitPath()
}

func init() {
	envLike = &defaultEnvLike{}
}

// Registry is the default executable registry.
var Registry = &ExecutableRegistry{data: make(map[string]Executable)}

// Register stores an executable description in the registry.
func (r *ExecutableRegistry) Register(name string, exe *Executable) {
	r.data[name] = *exe

	if exe.Variable == "" {
		sb := underscore([]rune(name), &underscoreOptions{Screaming: true})
		exe.Variable = string(sb)
	}
}

// Set stores an executable description without modifying derived defaults.
func (r *ExecutableRegistry) Set(name string, exe *Executable) {
	r.data[name] = *exe
}

// Get returns a registered executable by name.
func (r *ExecutableRegistry) Get(name string) (*Executable, bool) {
	item, ok := r.data[name]
	return &item, ok
}

// Has reports whether an executable is registered.
func (r *ExecutableRegistry) Has(name string) bool {
	_, ok := r.data[name]
	return ok
}

// Find resolves an executable path by name using registry metadata.
func (r *ExecutableRegistry) Find(name string, options *WhichOptions) (string, error) {
	m, ok := r.data[name]
	if !ok {
		sb := underscore([]rune(name), &underscoreOptions{Screaming: true})
		m = Executable{Name: name}
		m.Variable = string(sb)
		r.data[name] = m
	}

	if options == nil {
		options = &WhichOptions{}
	}

	if options.UseCache && m.Path != "" {
		return m.Path, nil
	}

	if m.Variable != "" {
		value := envLike.Get(m.Variable)
		if value != "" {
			value, _ = envLike.Expand(value)
			if value != "" {
				next, ok := WhichFirst(value, options)
				if ok {
					m.Path = next
					return m.Path, nil
				}
			}
		}
	}

	if m.Path != "" {
		next, ok := WhichFirst(m.Path, options)
		if ok {
			m.Path = next
			return m.Path, nil
		}
	}

	if runtime.GOOS == "windows" {
		for _, path := range m.Windows {
			if emptySpace(path) {
				continue
			}

			exe2, _ := envLike.Expand(path)
			if exe2 == "" {
				continue
			}

			next, ok := WhichFirst(exe2, options)
			if ok {
				m.Path = next
				return m.Path, nil
			}
		}

		return "", errors.New("executable not found: " + name)
	}

	if runtime.GOOS == "darwin" {
		for _, path := range m.Darwin {
			if emptySpace(path) {
				continue
			}

			exe2, _ := envLike.Expand(path)
			if exe2 == "" {
				continue
			}

			next, ok := WhichFirst(exe2, options)
			if ok {
				m.Path = next
				return m.Path, nil
			}
		}

		// fallthrough to unix
	}

	for _, path := range m.Linux {
		if emptySpace(path) {
			continue
		}

		exe2, _ := envLike.Expand(path)
		if exe2 == "" {
			continue
		}

		next, ok := WhichFirst(exe2, options)
		if ok {
			m.Path = next
			return m.Path, nil
		}
	}

	return "", errors.New("executable not found: " + name)
}

// Register stores an executable description in the default registry.
func Register(name string, exe *Executable) {
	Registry.Register(name, exe)
}

// Find resolves an executable path from the default registry.
func Find(name string, options *WhichOptions) (string, error) {
	return Registry.Find(name, options)
}

type underscoreOptions struct {
	PreserveCase bool
	Screaming    bool
}

func underscore(runes []rune, options *underscoreOptions) []rune {
	if len(runes) == 0 {
		return runes
	}

	sb := make([]rune, 0)
	last := rune(0)
	if options == nil {
		options = &underscoreOptions{}
	}

	for _, r := range runes {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				if unicode.IsLetter(last) && unicode.IsLower(last) {
					sb = append(sb, '_')
					if options.PreserveCase || options.Screaming {
						sb = append(sb, r)
						last = r
						continue
					}

					sb = append(sb, unicode.ToLower(r))
					last = r
					continue
				}

				if options.PreserveCase || options.Screaming {
					sb = append(sb, r)
					last = r
					continue
				}

				sb = append(sb, unicode.ToLower(r))
				last = r
				continue
			}

			if options.Screaming {
				sb = append(sb, unicode.ToUpper(r))
			} else if options.PreserveCase {
				sb = append(sb, r)
			} else {
				sb = append(sb, unicode.ToLower(r))
			}

			last = r
			continue
		}

		if unicode.IsNumber(r) {
			sb = append(sb, r)
			last = r
			continue
		}

		if r == '_' || r == '-' || unicode.IsSpace(r) {
			if len(sb) == 0 {
				continue
			}

			if last == '_' {
				continue
			}

			last = '_'
			sb = append(sb, last)
			continue
		}

	}

	if len(sb) > 0 && sb[len(sb)-1] == '_' {
		sb = sb[:len(sb)-1]
	}

	return sb
}
