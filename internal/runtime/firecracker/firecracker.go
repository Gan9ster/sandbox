package firecracker

import (
	"context"
	"errors"

	"github.com/your-org/hsr/sandbox"
)

// Runtime is a placeholder Firecracker implementation.
type Runtime struct{}

func New() *Runtime { return &Runtime{} }

func (r *Runtime) Run(ctx context.Context, t sandbox.Task) (string, error) {
	return "firecracker-task-id", errors.New("firecracker runtime not implemented")
}

func (r *Runtime) Wait(ctx context.Context, id string) (sandbox.Result, error) {
	return sandbox.Result{}, errors.New("firecracker runtime not implemented")
}

func (r *Runtime) Kill(ctx context.Context, id string) error {
	return errors.New("firecracker runtime not implemented")
}
