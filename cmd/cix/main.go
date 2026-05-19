package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tiepnguyen-cix/cix/pkg/executor"
	"github.com/tiepnguyen-cix/cix/pkg/output"
	"github.com/tiepnguyen-cix/cix/pkg/parser"
)

var version = "dev"

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "cix",
		Short: "Run and debug GitLab CI pipelines locally",
		Long:  "cix — run any GitLab CI pipeline on your machine. No push required.",
		Version: version,
	}

	root.AddCommand(
		runCmd(),
		validateCmd(),
		listCmd(),
	)

	return root
}

// ── cix run ──────────────────────────────────────────────────────────────────

func runCmd() *cobra.Command {
	var (
		jobName    string
		stageName  string
		ciFile     string
		dryRun     bool
		keepOnFail bool
		verbose    bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a pipeline, stage, or specific job locally",
		Example: `  cix run                        # run entire pipeline
  cix run --job build            # run single job
  cix run --stage test           # run all jobs in stage
  cix run --job build --dry-run  # preview without running`,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Header(version)

			pipeline, err := parser.ParseFile(ciFile)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			exec := executor.New(executor.Options{
				KeepOnFail: keepOnFail,
				Verbose:    verbose,
				DryRun:     dryRun,
			})

			ctx := context.Background()

			// Determine which jobs to run
			jobs := selectJobs(pipeline, jobName, stageName)
			if len(jobs) == 0 {
				return fmt.Errorf("no matching jobs found")
			}

			for _, job := range jobs {
				result, err := exec.RunJob(ctx, job, pipeline.Variables)
				if err != nil {
					return err
				}
				if result.Failed {
					os.Exit(1)
				}
				output.JobSuccess(job.Name, 0)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&jobName, "job", "j", "", "Run a specific job")
	cmd.Flags().StringVarP(&stageName, "stage", "s", "", "Run all jobs in a stage")
	cmd.Flags().StringVarP(&ciFile, "file", "f", ".gitlab-ci.yml", "Path to CI config file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview what would run without executing")
	cmd.Flags().BoolVar(&keepOnFail, "keep-on-fail", false, "Keep container alive on failure for debugging")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show full command output")

	return cmd
}

// ── cix validate ─────────────────────────────────────────────────────────────

func validateCmd() *cobra.Command {
	var ciFile string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate .gitlab-ci.yml syntax and structure",
		RunE: func(cmd *cobra.Command, args []string) error {
			pipeline, err := parser.ParseFile(ciFile)
			if err != nil {
				return fmt.Errorf("invalid: %w", err)
			}

			// Collect warnings
			var warnings []string
			for _, job := range pipeline.Jobs {
				if len(job.Script) == 0 {
					warnings = append(warnings, fmt.Sprintf("job %q has no script", job.Name))
				}
				if job.Image == "" {
					warnings = append(warnings, fmt.Sprintf("job %q has no image — will use alpine:latest", job.Name))
				}
			}

			output.ValidationOK(pipeline.Stages, len(pipeline.Jobs), warnings)
			return nil
		},
	}

	cmd.Flags().StringVarP(&ciFile, "file", "f", ".gitlab-ci.yml", "Path to CI config file")
	return cmd
}

// ── cix list ─────────────────────────────────────────────────────────────────

func listCmd() *cobra.Command {
	var ciFile string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all jobs in the pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			pipeline, err := parser.ParseFile(ciFile)
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			// Sort by stage order
			var jobs []output.JobInfo
			for _, stage := range pipeline.Stages {
				for name, job := range pipeline.Jobs {
					if job.Stage == stage {
						img := job.Image
						if img == "" {
							img = "alpine:latest"
						}
						jobs = append(jobs, output.JobInfo{
							Name:  name,
							Stage: stage,
							Image: img,
						})
					}
				}
			}

			// Jobs without matching stage
			for name, job := range pipeline.Jobs {
				found := false
				for _, s := range pipeline.Stages {
					if job.Stage == s {
						found = true
						break
					}
				}
				if !found {
					jobs = append(jobs, output.JobInfo{
						Name:  name,
						Stage: job.Stage,
						Image: job.Image,
					})
				}
			}

			output.JobList(jobs)
			return nil
		},
	}

	cmd.Flags().StringVarP(&ciFile, "file", "f", ".gitlab-ci.yml", "Path to CI config file")
	return cmd
}

// ── helpers ──────────────────────────────────────────────────────────────────

func selectJobs(pipeline *parser.Pipeline, jobName, stageName string) []*parser.Job {
	if jobName != "" {
		if job, ok := pipeline.Jobs[jobName]; ok {
			return []*parser.Job{job}
		}
		return nil
	}

	if stageName != "" {
		var jobs []*parser.Job
		for _, job := range pipeline.Jobs {
			if job.Stage == stageName {
				jobs = append(jobs, job)
			}
		}
		return jobs
	}

	// Run all jobs, sorted by stage
	var jobs []*parser.Job
	for _, stage := range pipeline.Stages {
		for _, job := range pipeline.Jobs {
			if job.Stage == stage {
				jobs = append(jobs, job)
			}
		}
	}
	return jobs
}