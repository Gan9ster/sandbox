package sandbox

// TaskMeta holds metadata used for scheduling decisions.
type TaskMeta struct {
	Lang        string // e.g. "python", "go"
	Network     bool   // whether network access is needed
	Binary      bool   // contains compiled binaries
	UnknownLang bool   // true if language couldn't be detected
}

// Task defines an executable job.
type Task struct {
	TaskMeta

	Image string   // container image to run
	Cmd   []string // command to execute
}

// Result is returned when a task finishes executing.
type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}
