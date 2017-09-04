package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/vgheri/gennaker/engine"
	"github.com/vgheri/gennaker/repository/pg"
	"github.com/vgheri/gennaker/utils"
)

var testhandler *Handler

func TestMain(m *testing.M) {
	var host, port, username, password, dbname string
	if os.Getenv("PG_HOST") == "" {
		host = "localhost"
	}
	if os.Getenv("PG_PORT") == "" {
		port = "5432"
	}
	if os.Getenv("PG_USERNAME") == "" {
		username = "postgres"
	}
	if os.Getenv("PG_PASSWORD") == "" {
		password = "password"
	}
	if os.Getenv("PG_DBNAME") == "" {
		dbname = "gennaker"
	}
	repository, err := pg.NewClient(host, port, username, password, dbname, 10)
	if err != nil {
		panic(err)
	}
	var chartsFolder string
	if chartsFolder = os.Getenv("CHARTS_FOLDER"); chartsFolder == "" {
		gopath := os.Getenv("GOPATH")
		chartsFolder = path.Join(gopath, "src", "github.com", "vgheri", "gennaker", "charts")
	}

	testengine := engine.New(repository, chartsFolder)
	testhandler = New(testengine)
}

func TestCreateDeploymentHandler(t *testing.T) {
	tt := []struct {
		name       string
		deployName string
		chartName  string
		repository string
		shouldErr  bool
	}{
		{name: "invalid deployment name", deployName: "", chartName: "test", repository: "http://nonexistent.com:8080", shouldErr: true},
		{name: "invalid chart name", deployName: "test", chartName: "", repository: "http://nonexistent.com:8080", shouldErr: true},
		{name: "invalid chart repository url", deployName: "test", chartName: "test", repository: "http://nonexistent.com:8080", shouldErr: true},
		{name: "should create deployment", deployName: utils.GenerateRandomString(10), chartName: "consul", repository: "https://kubernetes-charts.storage.googleapis.com", shouldErr: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bodyValue := CreateDeploymentRequest{Name: tc.deployName, ChartName: tc.chartName, RepositoryURL: tc.repository}
			bodyMarshaled, _ := json.Marshal(bodyValue)
			body := bytes.NewReader(bodyMarshaled)
			req, err := http.NewRequest("POST", "localhost:8080/api/v1/deployment", body)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			testhandler.CreateDeploymentHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.shouldErr {
				// do something
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status Bad Request; got %v", res.StatusCode)
				}
				// if msg := string(bytes.TrimSpace(b)); msg != tc.err {
				// 	t.Errorf("expected message %q; got %q", tc.err, msg)
				// }
				return
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}

			_, err = strconv.Atoi(string(bytes.TrimSpace(b)))
			if err != nil {
				t.Fatalf("expected an integer; got %s", b)
			}
		})
	}
}

func TestNewDeploymentReleaseNotificationHandler(t *testing.T) {
	tt := []struct {
		name          string
		deployName    string
		imageTag      string
		releaseValues string
		shouldErr     bool
	}{
		{name: "Empty deployment name", deployName: "", imageTag: "0.0.1", shouldErr: true},
		{name: "Empty image tag", deployName: "test", imageTag: "", shouldErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bodyValue := NewDeploymentReleaseNotificationRequest{DeploymentName: tc.deployName, ImageTag: tc.imageTag, ReleaseValues: tc.releaseValues}
			bodyMarshaled, _ := json.Marshal(bodyValue)
			body := bytes.NewReader(bodyMarshaled)
			req, err := http.NewRequest("POST", "/api/v1/deployment/newrelease", body)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			testhandler.NewDeploymentReleaseNotificationHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.shouldErr {
				// do something
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status Bad Request; got %v", res.StatusCode)
				}
				// if msg := string(bytes.TrimSpace(b)); msg != tc.err {
				// 	t.Errorf("expected message %q; got %q", tc.err, msg)
				// }
				return
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}
		})
	}
}

func TestPromoteReleaseHandler(t *testing.T) {
	tt := []struct {
		name          string
		deployName    string
		namespace     string
		releaseValues string
		shouldErr     bool
	}{
		{name: "Empty deployment name", deployName: "", namespace: "int", shouldErr: true},
		{name: "Empty namespace", deployName: "test", namespace: "", shouldErr: true},
		{name: "Invalid release values", deployName: "test", namespace: "int", releaseValues: "abc", shouldErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bodyValue := PromoteReleaseRequest{FromNamespace: tc.namespace, ReleaseValues: tc.releaseValues}
			bodyMarshaled, _ := json.Marshal(bodyValue)
			body := bytes.NewReader(bodyMarshaled)
			req, err := http.NewRequest("POST", fmt.Sprintf("/api/v1/deployment/%s/release/promote", tc.deployName), body)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			testhandler.NewDeploymentReleaseNotificationHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.shouldErr {
				// do something
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status Bad Request; got %v", res.StatusCode)
				}
				// if msg := string(bytes.TrimSpace(b)); msg != tc.err {
				// 	t.Errorf("expected message %q; got %q", tc.err, msg)
				// }
				return
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}
		})
	}
}
