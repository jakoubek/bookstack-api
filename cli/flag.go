// cli/flag.go - minimal flag parser (package main, zero deps)
package main

import (
	"fmt"
	"strings"
)

// FlagSet is a minimal flag parser.
type FlagSet struct {
	name    string
	flags   []*flag
	args    []string
	visited map[string]bool
}

type flag struct {
	name  string
	value any
	s     *string
	b     *bool
	i     *int
}

// newFlagSet creates a new flag set (lowercase to avoid shadowing stdlib flag.NewFlagSet).
func newFlagSet(name string) *FlagSet {
	return &FlagSet{name: name, visited: map[string]bool{}}
}

// String registers a string flag.
func (f *FlagSet) String(name, value, usage string) *string {
	p := new(string)
	*p = value
	f.flags = append(f.flags, &flag{name: name, value: p, s: p})
	return p
}

// Bool registers a boolean flag.
func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	*p = value
	f.flags = append(f.flags, &flag{name: name, value: p, b: p})
	return p
}

// Int registers an integer flag.
func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)
	*p = value
	f.flags = append(f.flags, &flag{name: name, value: p, i: p})
	return p
}

// Parse parses the given args.  Returns an error on unknown flags or bad values.
// Parsing stops at the first non-flag argument, so subcommands and their flags
// are passed through to Args().
func (f *FlagSet) Parse(args []string) error {
	i := 0
	for i < len(args) {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			f.args = append(f.args, args[i:]...)
			break
		}
		name := strings.TrimPrefix(a, "-")
		if strings.HasPrefix(name, "-") {
			name = strings.TrimPrefix(name, "-")
		}
		var val string
		if idx := strings.Index(name, "="); idx >= 0 {
			val = name[idx+1:]
			name = name[:idx]
		}
		fl := f.find(name)
		if fl == nil {
			return fmt.Errorf("flag provided but not defined: -%s", name)
		}
		f.visited[name] = true
		switch {
		case fl.s != nil:
			if val == "" && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				val = args[i+1]
				i++
			}
			*fl.s = val
		case fl.b != nil:
			*fl.b = true
		case fl.i != nil:
			if val == "" && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				val = args[i+1]
				i++
			}
			n, err := fmt.Sscanf(val, "%d", fl.i)
			if err != nil || n != 1 {
				return fmt.Errorf("invalid int value for -%s: %q", name, val)
			}
		}
		i++
	}
	return nil
}

// Args returns the remaining positional arguments after parsing.
func (f *FlagSet) Args() []string { return f.args }

func (f *FlagSet) find(name string) *flag {
	for _, fl := range f.flags {
		if fl.name == name {
			return fl
		}
	}
	return nil
}