package engine

import (
	"os"
	"path"
	"testing"
)

func Test_buildPipeline(t *testing.T) {

	expected := []*PipelineStep{
		&PipelineStep{
			TargetNamespace: "int",
			AutomaticDeploy: true,
			StepNumber:      1,
			NextSteps: []*PipelineStep{
				&PipelineStep{
					TargetNamespace:  "ppd",
					AutomaticDeploy:  false,
					StepNumber:       2,
					ParentStepNumber: 1,
					NextSteps: []*PipelineStep{
						&PipelineStep{
							TargetNamespace:  "prod",
							AutomaticDeploy:  false,
							StepNumber:       3,
							ParentStepNumber: 2,
							NextSteps:        []*PipelineStep{},
						},
					},
				},
			},
		},
	}
	gopath := os.Getenv("GOPATH")
	destination := path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "examples", "simple_pipeline")
	pipeline, err := buildPipeline(destination)
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
	if len(pipeline) != len(expected) {
		t.Fatalf("Malformed pipeline, expected %+v, got %+v", expected, pipeline)
	}

	gopath = os.Getenv("GOPATH")
	destination = path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "ababa", "simple_pipeline")
	_, err = buildPipeline(destination)
	if err == nil {
		t.Fatalf("Expected error with invalid path")
	}

}
