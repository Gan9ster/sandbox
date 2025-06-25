package docker

import (
	"context"
	"errors"

	"github.com/your-org/hsr/sandbox"
)

// Runtime is a placeholder Docker implementation.
type Runtime struct{}

func New() *Runtime { return &Runtime{} }

func (r *Runtime) Run(ctx context.Context, t sandbox.Task) (string, error) {
	return "docker-task-id", errors.New("docker runtime not implemented")
}

func (r *Runtime) Wait(ctx context.Context, id string) (sandbox.Result, error) {
	return sandbox.Result{}, errors.New("docker runtime not implemented")
}

func (r *Runtime) Kill(ctx context.Context, id string) error {
	return errors.New("docker runtime not implemented")
}
