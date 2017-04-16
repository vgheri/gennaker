package engine

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_buildPipeline(t *testing.T) {
	Convey("Testing buildPipeline()", t, FailureContinues, func() {
		Convey("With a valid path and gennaker.yml file", func() {
			Convey("Should build the exptected pipeline", func() {
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
				So(err, ShouldBeNil)
				So(pipeline, ShouldResemble, expected)
			})
		})
		Convey("With an invalid path", func() {
			Convey("Should raise error", func() {
				gopath := os.Getenv("GOPATH")
				destination := path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "ababa", "simple_pipeline")
				_, err := buildPipeline(destination)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
