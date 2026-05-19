package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// rawPipeline là raw YAML struct trước khi normalize
type rawPipeline struct {
	Stages    []string               `yaml:"stages"`
	Variables map[string]string      `yaml:"variables"`
	Jobs      map[string]rawJob      `yaml:"-"`
	Raw       map[string]interface{} `yaml:",inline"`
}

type rawJob struct {
	Stage        string            `yaml:"stage"`
	Image        string            `yaml:"image"`
	Variables    map[string]string `yaml:"variables"`
	BeforeScript []string          `yaml:"before_script"`
	Script       interface{}       `yaml:"script"` // string หรือ []string
	AfterScript  []string          `yaml:"after_script"`
	Artifacts    *rawArtifacts     `yaml:"artifacts"`
	Needs        interface{}       `yaml:"needs"`
	AllowFailure bool              `yaml:"allow_failure"`
	Extends      string            `yaml:"extends"`
}

type rawArtifacts struct {
	Paths []string `yaml:"paths"`
	When  string   `yaml:"when"`
}

// reserved keywords ở top-level — không phải job names
var reservedKeys = map[string]bool{
	"stages": true, "variables": true, "include": true,
	"workflow": true, "default": true, "image": true,
	"before_script": true, "after_script": true, "cache": true,
}

// ParseFile đọc file YAML và trả về Pipeline struct
func ParseFile(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return Parse(data)
}

// Parse nhận raw YAML bytes và trả về Pipeline struct
func Parse(data []byte) (*Pipeline, error) {
	// Unmarshal toàn bộ vào map trước để tách jobs ra
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	pipeline := &Pipeline{
		Jobs:      make(map[string]*Job),
		Variables: make(map[string]string),
	}

	// Parse top-level stages
	if stages, ok := raw["stages"]; ok {
		if stageList, ok := stages.([]interface{}); ok {
			for _, s := range stageList {
				if stage, ok := s.(string); ok {
					pipeline.Stages = append(pipeline.Stages, stage)
				}
			}
		}
	}

	// Default stages nếu không define
	if len(pipeline.Stages) == 0 {
		pipeline.Stages = []string{"build", "test", "deploy"}
	}

	// Parse top-level variables
	if vars, ok := raw["variables"]; ok {
		if varMap, ok := vars.(map[string]interface{}); ok {
			for k, v := range varMap {
				pipeline.Variables[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Parse từng job (bất kỳ key nào không phải reserved)
	for key, val := range raw {
		if reservedKeys[key] {
			continue
		}
		jobMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		job, err := parseJob(key, jobMap)
		if err != nil {
			return nil, fmt.Errorf("job %q: %w", key, err)
		}
		pipeline.Jobs[key] = job
	}

	return pipeline, nil
}

// parseJob convert raw map thành Job struct
func parseJob(name string, raw map[string]interface{}) (*Job, error) {
	// Re-marshal và unmarshal vào rawJob để leverage yaml tags
	data, _ := yaml.Marshal(raw)
	var rj rawJob
	if err := yaml.Unmarshal(data, &rj); err != nil {
		return nil, err
	}

	job := &Job{
		Name:         name,
		Stage:        rj.Stage,
		Image:        rj.Image,
		Variables:    rj.Variables,
		BeforeScript: rj.BeforeScript,
		AfterScript:  rj.AfterScript,
		AllowFailure: rj.AllowFailure,
	}

	if job.Stage == "" {
		job.Stage = "test" // GitLab default
	}
	if job.Variables == nil {
		job.Variables = make(map[string]string)
	}

	// script có thể là string đơn hoặc []string
	job.Script = normalizeScript(rj.Script)

	// artifacts
	if rj.Artifacts != nil {
		job.Artifacts = &Artifacts{
			Paths: rj.Artifacts.Paths,
			When:  rj.Artifacts.When,
		}
		if job.Artifacts.When == "" {
			job.Artifacts.When = "on_success"
		}
	}

	// needs: có thể là []string hoặc []map với "job" key
	job.Needs = normalizeNeeds(rj.Needs)

	return job, nil
}

func normalizeScript(raw interface{}) []string {
	if raw == nil {
		return nil
	}
	switch v := raw.(type) {
	case string:
		return []string{v}
	case []interface{}:
		var result []string
		for _, item := range v {
			result = append(result, fmt.Sprintf("%v", item))
		}
		return result
	}
	return nil
}

func normalizeNeeds(raw interface{}) []string {
	if raw == nil {
		return nil
	}
	var result []string
	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			switch n := item.(type) {
			case string:
				result = append(result, n)
			case map[string]interface{}:
				if jobName, ok := n["job"].(string); ok {
					result = append(result, jobName)
				}
			}
		}
	}
	return result
}