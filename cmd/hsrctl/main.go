package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gan9ster/sandbox/internal/runtime/docker"
	"github.com/gan9ster/sandbox/internal/runtime/firecracker"
	"github.com/gan9ster/sandbox/internal/runtime/gvisor"
	"github.com/gan9ster/sandbox/sandbox"
	"github.com/gan9ster/sandbox/scheduler"
)

func main() {
	policyPath := flag.String("policy", "config/policy.yaml", "policy file")
	image := flag.String("image", "", "container image")
	flag.Parse()
	if *image == "" {
		log.Fatal("image must be specified")
	}

	p, err := scheduler.LoadPolicy(*policyPath)
	if err != nil {
		log.Fatalf("load policy: %v", err)
	}
	if err := p.Validate(); err != nil {
		log.Fatalf("invalid policy: %v", err)
	}

	task := sandbox.Task{
		Image: *image,
		Cmd:   flag.Args(),
	}

	runtimeName, _ := p.SelectRuntime(task.TaskMeta)
	var rt sandbox.Runtime
	switch runtimeName {
	case "docker":
		rt = docker.New()
	case "gvisor":
		rt = gvisor.New()
	case "firecracker":
		rt = firecracker.New()
	default:
		log.Fatalf("unknown runtime %s", runtimeName)
	}

	id, err := rt.Run(context.Background(), task)
	if err != nil {
		log.Fatalf("run task: %v", err)
	}

	res, err := rt.Wait(context.Background(), id)
	if err != nil {
		log.Fatalf("wait: %v", err)
	}

	fmt.Print(res.Stdout)
	if res.ExitCode != 0 {
		log.Fatalf("task failed: %d\n%s", res.ExitCode, res.Stderr)
	}
}
