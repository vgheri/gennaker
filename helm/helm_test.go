package helm

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetRepositoryName(t *testing.T) {
	Convey("Testing GetRepositoryName()", t, FailureContinues, func() {
		Convey("With local url", func() {
			Convey("Should get `local`", func() {
				name, err := GetRepositoryName("http://127.0.0.1:8879/charts")
				So(err, ShouldBeNil)
				So(name, ShouldEqual, "local")
			})
		})
		Convey("With non existent url", func() {
			Convey("Should get empty string", func() {
				name, err := GetRepositoryName("http://blabla.com/charts")
				So(err, ShouldBeNil)
				So(name, ShouldBeEmpty)
			})
		})
	})
}

func Test_Fetch(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	destination := path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "charts")
	expectedDestination := path.Join(destination, "consul")
	Convey("Testing Fetch()", t, FailureContinues, func() {
		Convey("stable/consul to ./charts", func() {
			Convey("Should get $current_dir/charts/consul", func() {
				savePath, err := Fetch("stable", "consul", "", destination)
				So(err, ShouldBeNil)
				So(savePath, ShouldEqual, expectedDestination)
			})
		})
		Convey("With non existent combination repoName/chartName", func() {
			Convey("Should get error", func() {
				savePath, err := Fetch("uistiti", "test", "", destination)
				So(err, ShouldNotBeNil)
				So(savePath, ShouldBeEmpty)
			})
		})
		Convey("With empty repository name", func() {
			Convey("Should get error", func() {
				savePath, err := Fetch("", "test", "", destination)
				So(err, ShouldNotBeNil)
				So(savePath, ShouldBeEmpty)
			})
		})
		Convey("With empty chart name", func() {
			Convey("Should get error", func() {
				savePath, err := Fetch("stable", "", "", destination)
				So(err, ShouldNotBeNil)
				So(savePath, ShouldBeEmpty)
			})
		})
	})
	err := os.RemoveAll(destination)
	if err != nil {
		t.Fatal(err)
	}
}
