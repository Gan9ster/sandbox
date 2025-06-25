package firecracker

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/your-org/hsr/sandbox"
)

// Runtime implements sandbox.Runtime using the ignite Firecracker CLI.
type Runtime struct{}

func New() *Runtime { return &Runtime{} }

func (r *Runtime) Run(ctx context.Context, t sandbox.Task) (string, error) {
	idBytes, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(idBytes))

	args := []string{"run", t.Image, "--name", id, "--quiet"}
	if len(t.Cmd) > 0 {
		args = append(args, "--cmd", strings.Join(t.Cmd, " "))
	}
	cmd := exec.CommandContext(ctx, "ignite", args...)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return id, nil
}

func (r *Runtime) Wait(ctx context.Context, id string) (sandbox.Result, error) {
	logsCmd := exec.CommandContext(ctx, "ignite", "logs", id)
	logs, _ := logsCmd.CombinedOutput()

	inspectCmd := exec.CommandContext(ctx, "ignite", "inspect", id, "--format", "{{.Status.ExitCode}}")
	out, err := inspectCmd.Output()
	if err != nil {
		return sandbox.Result{}, err
	}
	exitStr := strings.TrimSpace(string(out))
	exitCode := 0
	if exitStr != "" {
		if c, err := strconv.Atoi(exitStr); err == nil {
			exitCode = c
		}
	}

	exec.CommandContext(context.Background(), "ignite", "rm", "-f", id).Run()

	return sandbox.Result{ExitCode: exitCode, Stdout: string(logs)}, nil
}

func (r *Runtime) Kill(ctx context.Context, id string) error {
	cmd := exec.CommandContext(ctx, "ignite", "stop", id)
	return cmd.Run()
}
