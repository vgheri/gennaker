package route

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

	"github.com/vgheri/gennaker/api/handler"
	"github.com/vgheri/gennaker/engine"
	"github.com/vgheri/gennaker/repository/pg"
	"github.com/vgheri/gennaker/utils"
)

var testhandler *handler.Handler
var server *httptest.Server

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
	testhandler = handler.New(testengine)
	router := NewRouter(testhandler)
	server = httptest.NewServer(router)
	defer server.Close()
}

func TestCreateDeploymentRouting(t *testing.T) {
	bodyValue := handler.CreateDeploymentRequest{Name: utils.GenerateRandomString(10), ChartName: "consul", RepositoryURL: "https://kubernetes-charts.storage.googleapis.com"}
	bodyMarshaled, _ := json.Marshal(bodyValue)
	body := bytes.NewReader(bodyMarshaled)
	res, err := http.Post(fmt.Sprintf("%s/vapi/v1/deployment", server.URL), "application/json", body)
	if err == nil {
		t.Fatal("sould not have been able to POST request to /vapi/v1/deployment")
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status NotFound; got %v", res.Status)
	}

	res, err = http.Post(fmt.Sprintf("%s/api/v1/deployment", server.URL), "application/json", body)
	if err != nil {
		t.Fatalf("could not POST request to /vapi/v1/deployment, err: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	_, err = strconv.Atoi(string(bytes.TrimSpace(b)))
	if err != nil {
		t.Fatalf("expected an integer; got %s", b)
	}
}

func TestNewDeploymentReleaseNotificationRouting(t *testing.T) {
	bodyValue := handler.NewDeploymentReleaseNotificationRequest{DeploymentName: "a", ImageTag: "b", ReleaseValues: ""}
	bodyMarshaled, _ := json.Marshal(bodyValue)
	body := bytes.NewReader(bodyMarshaled)
	res, err := http.Post(fmt.Sprintf("%s/vapi/v1/deployment/newrelease", server.URL), "application/json", body)
	if err == nil {
		t.Fatal("sould not have been able to POST request to /api/v1/deployment/newrelease")
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status NotFound; got %v", res.Status)
	}
}
