package scheduler

import "github.com/gan9ster/sandbox/sandbox"

// Strategy selects a runtime for the given task metadata.
type Strategy interface {
	SelectRuntime(task sandbox.TaskMeta) (runtimeName string, reason string)
}
