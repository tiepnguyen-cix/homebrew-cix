package executor

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/tiepnguyen-cix/cix/pkg/output"
	"github.com/tiepnguyen-cix/cix/pkg/parser"
)

// Options cấu hình cho executor
type Options struct {
	KeepOnFail bool
	Verbose    bool
	DryRun     bool
}

// Executor chạy CI jobs bằng Docker
type Executor struct {
	opts Options
}

func New(opts Options) *Executor {
	return &Executor{opts: opts}
}

// RunJob chạy một job cụ thể
func (e *Executor) RunJob(ctx context.Context, job *parser.Job, globalVars map[string]string) (*parser.JobResult, error) {
	image := job.Image
	if image == "" {
		image = "alpine:latest"
	}

	output.JobStart(job.Name, image)

	if e.opts.DryRun {
		e.printDryRun(job)
		return &parser.JobResult{JobName: job.Name}, nil
	}

	// Merge variables: global → job-level
	vars := mergeVars(globalVars, job.Variables)
	vars = injectCIVars(vars, job.Name)

	// Pull image nếu chưa có
	if err := e.pullImage(ctx, image); err != nil {
		return nil, fmt.Errorf("pull image %s: %w", image, err)
	}

	// Build tất cả commands: before_script + script + after_script
	allSteps := buildSteps(job)

	result := &parser.JobResult{JobName: job.Name}

	for _, step := range allSteps {
		stepResult, err := e.runStep(ctx, image, step, vars)
		result.Steps = append(result.Steps, stepResult)

		if err != nil && !job.AllowFailure {
			result.Failed = true
			output.StepFail(step, 0)
			output.JobFailed(job.Name, stepResult.ExitCode, stepResult.Output)
			return result, nil
		}
		output.StepOK(step, 0)
	}

	return result, nil
}

// runStep chạy một command trong container
func (e *Executor) runStep(ctx context.Context, image, command string, vars map[string]string) (parser.StepResult, error) {
	args := []string{"run", "--rm"}

	// Inject env vars
	for k, v := range vars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, image, "sh", "-c", command)

	cmd := exec.CommandContext(ctx, "docker", args...)

	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return parser.StepResult{Command: command, ExitCode: 1}, err
	}

	if err := cmd.Start(); err != nil {
		return parser.StepResult{Command: command, ExitCode: 1}, err
	}

	// Stream stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if e.opts.Verbose {
			output.LogLine(scanner.Text())
		}
	}

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	result := parser.StepResult{
		Command:  command,
		ExitCode: exitCode,
		Output:   stderrBuf.String(),
	}

	if exitCode != 0 {
		return result, fmt.Errorf("command exited with code %d", exitCode)
	}

	return result, nil
}

// pullImage pulls Docker image jika belum ada
func (e *Executor) pullImage(ctx context.Context, image string) error {
	// Check nếu image đã có local
	checkCmd := exec.CommandContext(ctx, "docker", "image", "inspect", image)
	if err := checkCmd.Run(); err == nil {
		return nil // already exists
	}

	fmt.Printf("  pulling %s...\n", image)
	pullCmd := exec.CommandContext(ctx, "docker", "pull", image)
	pullCmd.Stdout = nil
	return pullCmd.Run()
}

// printDryRun hiển thị sẽ chạy gì mà không thực sự chạy
func (e *Executor) printDryRun(job *parser.Job) {
	fmt.Printf("  [dry-run] image: %s\n", job.Image)
	for _, cmd := range job.BeforeScript {
		fmt.Printf("  [dry-run] before: %s\n", cmd)
	}
	for _, cmd := range job.Script {
		fmt.Printf("  [dry-run] script: %s\n", cmd)
	}
}

func buildSteps(job *parser.Job) []string {
	var steps []string
	steps = append(steps, job.BeforeScript...)
	steps = append(steps, job.Script...)
	// after_script luôn chạy kể cả khi fail — handle riêng sau
	return steps
}

func mergeVars(global, local map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range global {
		result[k] = v
	}
	for k, v := range local {
		result[k] = v
	}
	return result
}

// injectCIVars inject các biến GitLab CI chuẩn
func injectCIVars(vars map[string]string, jobName string) map[string]string {
	// Lấy git info từ local repo
	branch := gitOutput("git", "branch", "--show-current")
	sha := gitOutput("git", "rev-parse", "--short", "HEAD")

	vars["CI"] = "true"
	vars["CI_JOB_NAME"] = jobName
	vars["CI_COMMIT_BRANCH"] = branch
	vars["CI_COMMIT_SHORT_SHA"] = sha
	vars["CI_PIPELINE_SOURCE"] = "local"
	vars["CI_PROJECT_DIR"] = "/builds/project"

	return vars
}

func gitOutput(args ...string) string {
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}