package sandbox

import "context"

// Runtime represents a concrete sandbox backend implementation.
type Runtime interface {
	// Run starts a task and returns immediately with an ID.
	Run(ctx context.Context, t Task) (string, error)

	// Wait blocks until the task exits or ctx is cancelled.
	Wait(ctx context.Context, id string) (Result, error)

	// Kill attempts to stop the task. No-op if already finished.
	Kill(ctx context.Context, id string) error
}
