package engine

import (
	"io/ioutil"
	"path"
	"sort"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

type YamlPipelineStep struct {
	Step       int
	Namespace  string
	Autodeploy bool
	ParentStep int `yaml:"parent_step,omitempty"`
}

type YamlPipeline struct {
	Steps []*YamlPipelineStep
}

type YamlContent struct {
	Version  int
	Pipeline *YamlPipeline
}

type ByOrder []*YamlPipelineStep

func (a ByOrder) Len() int           { return len(a) }
func (a ByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool { return a[i].Step < a[j].Step }

func buildPipeline(pathToChartOnDisk string) ([]*PipelineStep, error) {
	yamlContent := YamlContent{}
	fullPath := path.Join(pathToChartOnDisk, "gennaker.yml")
	gennakerContent, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, "Read file gennaker.yml failed")
	}
	err = yaml.Unmarshal(gennakerContent, &yamlContent)
	if err != nil {
		return nil, errors.Wrap(err, "yaml file Umarshal failed")
	}
	sort.Sort(ByOrder(yamlContent.Pipeline.Steps))
	stepsMap := make(map[int]*PipelineStep)
	for _, s := range yamlContent.Pipeline.Steps {
		step := &PipelineStep{
			StepNumber:       s.Step,
			ParentStepNumber: s.ParentStep,
			TargetNamespace:  s.Namespace,
			AutomaticDeploy:  s.Autodeploy,
			NextSteps:        []*PipelineStep{},
		}
		stepsMap[step.StepNumber] = step
		// Add itself to list of nextsteps of parent step, if any
		if step.ParentStepNumber > 0 {
			if parent, found := stepsMap[step.ParentStepNumber]; found {
				parent.NextSteps = append(parent.NextSteps, step)
			}
		}
	}
	pipeline := []*PipelineStep{}
	for _, s := range stepsMap {
		if s.ParentStepNumber == 0 {
			pipeline = append(pipeline, s)
		}
	}

	return pipeline, nil
}
