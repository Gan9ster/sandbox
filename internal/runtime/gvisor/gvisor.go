package gvisor

import (
	"context"
	"errors"

	"github.com/your-org/hsr/sandbox"
)

// Runtime is a placeholder gVisor implementation.
type Runtime struct{}

func New() *Runtime { return &Runtime{} }

func (r *Runtime) Run(ctx context.Context, t sandbox.Task) (string, error) {
	return "gvisor-task-id", errors.New("gvisor runtime not implemented")
}

func (r *Runtime) Wait(ctx context.Context, id string) (sandbox.Result, error) {
	return sandbox.Result{}, errors.New("gvisor runtime not implemented")
}

func (r *Runtime) Kill(ctx context.Context, id string) error {
	return errors.New("gvisor runtime not implemented")
}
