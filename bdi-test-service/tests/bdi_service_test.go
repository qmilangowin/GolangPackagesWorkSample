package app_test

import (
	"bdi-test-service/app"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShowAllConfigurations(t *testing.T) {

	server := app.Server{}
	server.Initialize()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/sta/v1/bdi_test_service/configurations", nil)

	handler := http.HandlerFunc(server.ShowAllConfigurationsRoute)
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	got := strings.TrimSpace(string(body))
	want := `{"default":{"sourcefolder":"/home/data/stories.table/stories.parquet","datasetname":"hacker"}}`
	if got != want {
		t.Errorf("GET to ShowAllConfigurationsRoute: GOT: %s WANT: %s", got, want)
	}

}

func TestCreateNewConfiguration(t *testing.T) {

	app := app.Server{}
	app.Initialize()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/sta/v1/bdi_test_service/configurations", nil)
	app.CreateNewConfigurationRoute(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	got := strings.TrimSpace(string(body))
	want := `{"1":{"sourcefolder":"","datasetname":""},"default":{"sourcefolder":"/home/data/stories.table/stories.parquet","datasetname":"hacker"}}`
	if got != want {
		t.Errorf("PATCH to CreateNewConfigurationRoute: GOT: %s WANT: %s", got, want)
	}
}
