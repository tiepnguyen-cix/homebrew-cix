package parser

import (
	"testing"
)

func TestParse_BasicPipeline(t *testing.T) {
	yaml := `
stages:
  - build
  - test

variables:
  NODE_ENV: test

build:
  stage: build
  image: node:20-alpine
  script:
    - npm run build

test:
  stage: test
  image: node:20-alpine
  before_script:
    - npm ci
  script:
    - npm run test
  allow_failure: true
`

	pipeline, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check stages
	if len(pipeline.Stages) != 2 {
		t.Errorf("expected 2 stages, got %d", len(pipeline.Stages))
	}

	// Check variables
	if pipeline.Variables["NODE_ENV"] != "test" {
		t.Errorf("expected NODE_ENV=test, got %s", pipeline.Variables["NODE_ENV"])
	}

	// Check jobs
	if len(pipeline.Jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(pipeline.Jobs))
	}

	buildJob := pipeline.Jobs["build"]
	if buildJob == nil {
		t.Fatal("build job not found")
	}
	if buildJob.Image != "node:20-alpine" {
		t.Errorf("expected image node:20-alpine, got %s", buildJob.Image)
	}
	if len(buildJob.Script) != 1 {
		t.Errorf("expected 1 script command, got %d", len(buildJob.Script))
	}

	testJob := pipeline.Jobs["test"]
	if !testJob.AllowFailure {
		t.Error("expected allow_failure=true for test job")
	}
}

func TestParse_ScriptAsString(t *testing.T) {
	yaml := `
build:
  stage: build
  script: echo hello
`
	pipeline, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	job := pipeline.Jobs["build"]
	if len(job.Script) != 1 || job.Script[0] != "echo hello" {
		t.Errorf("expected script=[echo hello], got %v", job.Script)
	}
}

func TestParse_DefaultStages(t *testing.T) {
	yaml := `
build:
  script:
    - make build
`
	pipeline, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(pipeline.Stages) != 3 {
		t.Errorf("expected 3 default stages, got %d: %v", len(pipeline.Stages), pipeline.Stages)
	}
}

func TestParse_Artifacts(t *testing.T) {
	yaml := `
build:
  stage: build
  script:
    - make build
  artifacts:
    paths:
      - dist/
      - build/
    when: always
`
	pipeline, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	job := pipeline.Jobs["build"]
	if job.Artifacts == nil {
		t.Fatal("expected artifacts, got nil")
	}
	if len(job.Artifacts.Paths) != 2 {
		t.Errorf("expected 2 artifact paths, got %d", len(job.Artifacts.Paths))
	}
	if job.Artifacts.When != "always" {
		t.Errorf("expected when=always, got %s", job.Artifacts.When)
	}
}

func TestParse_Needs(t *testing.T) {
	yaml := `
test:
  stage: test
  needs:
    - build
    - job: lint
  script:
    - npm test
`
	pipeline, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	job := pipeline.Jobs["test"]
	if len(job.Needs) != 2 {
		t.Errorf("expected 2 needs, got %d: %v", len(job.Needs), job.Needs)
	}
}