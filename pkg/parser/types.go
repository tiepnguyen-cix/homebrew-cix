package parser

type Pipeline struct {
	Stages    []string
	Jobs      map[string]*Job
	Variables map[string]string
}

type Job struct {
	Name         string
	Stage        string
	Image        string
	Variables    map[string]string
	BeforeScript []string
	Script       []string
	AfterScript  []string
	Artifacts    *Artifacts
	Needs        []string
	AllowFailure bool
}

type Artifacts struct {
	Paths []string
	When  string // always, on_success, on_failure
}


type StepResult struct {
	Command  string
	ExitCode int
	Duration float64
	Output   string
}


type JobResult struct {
	JobName  string
	Steps    []StepResult
	Failed   bool
	Duration float64
}