package scheduler

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/gan9ster/sandbox/sandbox"
)

// Policy describes a set of matching rules that map tasks to runtimes.
type Policy struct {
	Rules   []Rule
	Default string
}

type Rule struct {
	Match MatchCriteria
	Use   string
}

type MatchCriteria struct {
	Lang        string
	Network     *bool
	Binary      *bool
	UnknownLang *bool
}

// LoadPolicy loads a Policy from a very small subset of YAML.
func LoadPolicy(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		p       Policy
		current *Rule
		inMatch bool
		rules   []Rule
		scanner = bufio.NewScanner(f)
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		switch {
		case line == "rules:":
			continue
		case strings.HasPrefix(line, "- match:"):
			if current != nil {
				rules = append(rules, *current)
			}
			current = &Rule{}
			inMatch = true
		case strings.HasPrefix(line, "use:"):
			if current != nil {
				current.Use = strings.Trim(strings.TrimPrefix(line, "use:"), " \"")
			}
			inMatch = false
		case strings.HasPrefix(line, "default:"):
			p.Default = strings.Trim(strings.TrimPrefix(line, "default:"), " \"")
		default:
			if inMatch && current != nil {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 {
					continue
				}
				key := strings.TrimSpace(parts[0])
				val := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				switch key {
				case "lang":
					current.Match.Lang = val
				case "network":
					b, err := strconv.ParseBool(val)
					if err == nil {
						current.Match.Network = &b
					}
				case "binary":
					b, err := strconv.ParseBool(val)
					if err == nil {
						current.Match.Binary = &b
					}
				case "unknown-lang":
					b, err := strconv.ParseBool(val)
					if err == nil {
						current.Match.UnknownLang = &b
					}
				}
			}
		}
	}

	if current != nil {
		rules = append(rules, *current)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	p.Rules = rules
	return &p, nil
}

// SelectRuntime chooses a runtime based on task metadata.
func (p *Policy) SelectRuntime(task sandbox.TaskMeta) (string, string) {
	for _, r := range p.Rules {
		if r.Match.matches(task) {
			return r.Use, "matched rule"
		}
	}
	return p.Default, "default"
}

func (m MatchCriteria) matches(task sandbox.TaskMeta) bool {
	if m.Lang != "" && m.Lang != task.Lang {
		return false
	}
	if m.Network != nil && *m.Network != task.Network {
		return false
	}
	if m.Binary != nil && *m.Binary != task.Binary {
		return false
	}
	if m.UnknownLang != nil && *m.UnknownLang != task.UnknownLang {
		return false
	}
	return true
}

// Validate ensures the policy configuration is sane.
func (p *Policy) Validate() error {
	if p.Default == "" {
		return errors.New("default runtime must be specified")
	}
	return nil
}
