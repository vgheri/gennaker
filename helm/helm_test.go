package helm

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func Test_GetRepositoryName(t *testing.T) {
	name, err := GetRepositoryName("http://127.0.0.1:8879/charts")
	if err != nil {
		t.Fatalf("Expected OK, got error. Error details: %v", err)
	}
	if name != "local" {
		t.Fatalf("Expected local, got %s", name)
	}
	name, err = GetRepositoryName("http://blabla.com/charts")
	if err != nil {
		t.Fatalf("Expected OK, got error. Error details: %v", err)
	}
	if name != "" {
		t.Fatalf("Expected name to be empty with invalid repo URL")
	}
}

func Test_Fetch(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	destination := path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "charts")
	expectedDestination := path.Join(destination, "consul")

	savePath, err := Fetch("stable", "consul", "", destination)
	if err != nil {
		t.Fatalf("Expected success with stable/consul. Error details: %v", err)
	}
	if savePath != expectedDestination {
		t.Fatalf("Expected destination %s, got %s", savePath, expectedDestination)
	}

	savePath, err = Fetch("uistiti", "test", "", destination)
	if err == nil {
		t.Fatalf("Expected to get error with invalid repository, got nothing")
	}
	if savePath != "" {
		t.Fatalf("SavePath should be empty")
	}

	savePath, err = Fetch("", "test", "", destination)
	if err == nil {
		t.Fatalf("Expected error with empty repository name")
	}
	if savePath != "" {
		t.Fatalf("SavePath should be empty with empty repo name")
	}

	savePath, err = Fetch("stable", "", "", destination)
	if err == nil {
		t.Fatalf("Expected error with empty chart name")
	}
	if savePath != "" {
		t.Fatalf("SavePath should be empty with empty chart name")
	}

	err = os.RemoveAll(destination)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InstallOrUpgrade(t *testing.T) {
	var tt = []struct {
		testName       string
		releaseName    string
		namespace      string
		repositoryName string
		chartName      string
		valuesFilePath string
		releaseValues  string
		shouldErr      bool
	}{
		{testName: "Successfull install", releaseName: /*utils.GenerateRandomString(10)*/ "happy-panda",
			repositoryName: "stable", chartName: "consul", valuesFilePath: "", releaseValues: "", shouldErr: false},
		{testName: "Successfull upgrade", releaseName: /*utils.GenerateRandomString(10)*/ "happy-panda",
			repositoryName: "stable", chartName: "consul", valuesFilePath: "", releaseValues: "", shouldErr: false},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			output, err := InstallOrUpgrade(tc.releaseName, tc.namespace, tc.repositoryName, tc.chartName, tc.valuesFilePath, tc.releaseValues)
			if tc.shouldErr {
				if err == nil {
					t.Fatalf("Expected test to fail. Install output %s", output)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected InstallOrUpgrade to succeed, got %v", err)
			}
			err = deleteRelease(tc.releaseName)
			if err != nil {
				t.Fatalf("Could not delete release %s, err %v", tc.releaseName, err)
			}
		})
	}
}

func Test_Status(t *testing.T) {
	var tt = []struct {
		testName      string
		releaseName   string
		shouldInstall bool
		shouldErr     bool
	}{
		{testName: "Can read status", releaseName: "happy-panda", shouldInstall: true},
		{testName: "Error on non existing release", releaseName: "i-dont-exist", shouldErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			if tc.shouldInstall {
				_, err := InstallOrUpgrade(tc.releaseName, "default", "stable", "consul", "", "")
				if err != nil {
					t.Fatalf("Could not setup the test by installing a chart. Error details: %v", err)
				}
			}
			_, output, err := Status(tc.releaseName)
			if tc.shouldErr {
				if err == nil {
					t.Fatalf("Expected test to fail. Got output %s", output)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected test to succeed, got output %s and err %v", output, err)
			}
			if len(output) == 0 {
				t.Fatalf("Expected output not to be empty")
			}
			err = deleteRelease(tc.releaseName)
			if err != nil {
				t.Fatalf("Could not delete release %s, err %v", tc.releaseName, err)
			}
		})
	}
}

func deleteRelease(name string) error {
	cmdName := "helm"
	var cmdArgs = []string{"del", "--purge"}
	cmdArgs = append(cmdArgs, name)

	cmd := exec.Command(cmdName, cmdArgs...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
