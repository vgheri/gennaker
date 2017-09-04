package engine

import (
	"strings"
	"testing"
)

func Test_PromoteRelease(t *testing.T) {
	invalidPromoteRequest := &PromoteRequest{
		DeploymentName: "",
		FromNamespace:  "int",
	}
	_, err := testEngine.PromoteRelease(invalidPromoteRequest)
	if !strings.HasPrefix(err.Error(), "Promote request is invalid") {
		t.Fatalf("Expected invalid promote request due to empty deployment name, got %v instead", err)
	}

	invalidPromoteRequest = &PromoteRequest{
		DeploymentName: "abc",
		FromNamespace:  "",
	}
	_, err = testEngine.PromoteRelease(invalidPromoteRequest)
	if !strings.HasPrefix(err.Error(), "Promote request is invalid") {
		t.Fatalf("Expected invalid promote request due to empty namespace, got %v instead", err)
	}

	invalidPromoteRequest = &PromoteRequest{
		DeploymentName: "abc",
		FromNamespace:  "int",
		ReleaseValues:  "abc",
	}
	_, err = testEngine.PromoteRelease(invalidPromoteRequest)
	if !strings.HasPrefix(err.Error(), "Promote request is invalid") {
		t.Fatalf("Expected invalid promote request due to invalid release values, got %v instead", err)
	}
}
