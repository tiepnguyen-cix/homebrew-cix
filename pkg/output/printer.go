package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	cyan   = color.New(color.FgCyan)
	bold   = color.New(color.Bold)
	dim    = color.New(color.Faint)
)

// Header prints the cix banner
func Header(version string) {
	bold.Printf("cix")
	dim.Printf(" %s · GitLab CI local runner\n\n", version)
}

// JobStart prints job start line
func JobStart(jobName, image string) {
	cyan.Printf("▶ job: ")
	bold.Printf("%s", jobName)
	dim.Printf("  (%s)\n", image)
}

// StepOK prints a successful step
func StepOK(command string, dur time.Duration) {
	green.Printf("  ✓ ")
	fmt.Printf("%-45s", truncate(command, 45))
	dim.Printf("  %.1fs\n", dur.Seconds())
}

// StepFail prints a failed step
func StepFail(command string, dur time.Duration) {
	red.Printf("  ✗ ")
	fmt.Printf("%s\n", command)
}

// StepRunning prints a step that's currently running
func StepRunning(command string) {
	dim.Printf("  ○ %s\n", command)
}

// JobFailed prints job failure summary
func JobFailed(jobName string, exitCode int, stderr string) {
	fmt.Println()
	red.Printf("FAILED")
	fmt.Printf(" · %s exited with code %d\n", jobName, exitCode)

	if stderr != "" {
		fmt.Println()
		yellow.Printf("stderr output:\n")
		for _, line := range strings.Split(strings.TrimSpace(stderr), "\n") {
			fmt.Printf("  %s\n", line)
		}
	}

	fmt.Println()
	dim.Printf("tip: run with --keep-on-fail to inspect the container\n")
}

// JobSuccess prints job success
func JobSuccess(jobName string, dur time.Duration) {
	fmt.Println()
	green.Printf("PASSED")
	fmt.Printf(" · %s completed in %.1fs\n", jobName, dur.Seconds())
}

// PipelineSuccess prints overall pipeline success
func PipelineSuccess(dur time.Duration) {
	fmt.Println()
	green.Printf("Pipeline passed")
	dim.Printf(" in %.1fs\n", dur.Seconds())
}

// ValidationOK prints validation results
func ValidationOK(stages []string, jobCount int, warnings []string) {
	green.Printf("✓ syntax valid\n")
	green.Printf("✓ %d stage(s): %s\n", len(stages), strings.Join(stages, ", "))
	green.Printf("✓ %d job(s) found\n", jobCount)

	for _, w := range warnings {
		yellow.Printf("⚠ %s\n", w)
	}
}

// JobList prints list of jobs
func JobList(jobs []JobInfo) {
	maxName := 0
	for _, j := range jobs {
		if len(j.Name) > maxName {
			maxName = len(j.Name)
		}
	}

	for _, j := range jobs {
		cyan.Printf("  %-*s", maxName+2, j.Name)
		dim.Printf("stage: %-12s  image: %s\n", j.Stage, j.Image)
	}
}

// JobInfo holds display info for a job
type JobInfo struct {
	Name  string
	Stage string
	Image string
}

// LogLine prints a raw log line from container
func LogLine(line string) {
	dim.Printf("  │ ")
	fmt.Println(line)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}