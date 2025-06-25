package gvisor

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gan9ster/sandbox/sandbox"
)

// Runtime implements sandbox.Runtime using Docker with the runsc runtime.
type Runtime struct{}

func New() *Runtime { return &Runtime{} }

func (r *Runtime) Run(ctx context.Context, t sandbox.Task) (string, error) {
	args := append([]string{"run", "-d", "--runtime=runsc", "--rm", t.Image}, t.Cmd...)
	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(out))
	return id, nil
}

func (r *Runtime) Wait(ctx context.Context, id string) (sandbox.Result, error) {
	waitCmd := exec.CommandContext(ctx, "docker", "wait", id)
	out, err := waitCmd.Output()
	if err != nil {
		return sandbox.Result{}, err
	}
	exitStr := strings.TrimSpace(string(out))
	exitCode, err := strconv.Atoi(exitStr)
	if err != nil {
		return sandbox.Result{}, err
	}

	logsCmd := exec.CommandContext(ctx, "docker", "logs", id)
	logs, _ := logsCmd.CombinedOutput()

	exec.CommandContext(context.Background(), "docker", "rm", id).Run()

	return sandbox.Result{ExitCode: exitCode, Stdout: string(logs)}, nil
}

func (r *Runtime) Kill(ctx context.Context, id string) error {
	cmd := exec.CommandContext(ctx, "docker", "kill", id)
	return cmd.Run()
}
