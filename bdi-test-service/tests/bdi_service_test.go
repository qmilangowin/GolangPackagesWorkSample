package app_test

import (
	"bdi-test-service/handlers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShowAllConfigurations(t *testing.T) {

	handlers.Sourcepath = "/home/data/stories.table/stories.parquet"
	handlers.Dataset = "hacker"
	server := handlers.Server{}
	server.Initialize()
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/sta/v1/bdi_test_service/configurations", nil)
	if err != nil {
		t.Fatal(err)
	}

	server.ShowAllConfigurationsRoute(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("WANT %d; GOT %d", http.StatusOK, resp.StatusCode)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	got := strings.TrimSpace(string(body))
	want := `{"default":{"sourcefolder":"/home/data/stories.table/stories.parquet","datasetname":"hacker"}}`
	if got != want {
		t.Errorf("GET to ShowAllConfigurationsRoute: GOT: %s WANT: %s", got, want)
	}

}

func TestCreateNewConfiguration(t *testing.T) {

	handlers := handlers.Server{}
	handlers.Initialize()
	w := httptest.NewRecorder()
	bodyReader := strings.NewReader(`{"sourcefolder": "/home/data", "datasetname": "tpch10"}`)
	r, err := http.NewRequest(http.MethodPatch, "/sta/v1/bdi_test_service/configurations", bodyReader)
	if err != nil {
		t.Fatal(err)
	}

	handlers.CreateNewConfigurationRoute(w, r)
	resp := w.Result()
	if resp.StatusCode != 202 {
		t.Errorf("WANT %d; GOT %d", 202, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	got := strings.TrimSpace(string(body))
	want := `{"1":{"sourcefolder":"/home/data","datasetname":"tpch10"},"default":{"sourcefolder":"/home/data/stories.table/stories.parquet","datasetname":"hacker"}}`
	if got != want {
		t.Errorf("PATCH to CreateNewConfigurationRoute: GOT: %s WANT: %s", got, want)
	}
}
