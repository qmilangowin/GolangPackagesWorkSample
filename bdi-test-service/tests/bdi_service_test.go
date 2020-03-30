package app_test

import (
	"bdi-test-service/app"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestShowAllConfigurations(t *testing.T) {

	server := app.Server{}
	server.Initialize()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/bdi_test_service/configurations", nil)

	handler := http.HandlerFunc(server.ShowAllConfigurationsRoute)
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	got := strings.TrimSpace(string(body))
	want := `{"default":{"sourcefolder":"/home/data","datasetname":"hacker"}}`
	if got != want {
		t.Errorf("GET to ShowAllConfigurationsRoute: GOT: %s WANT: %s", got, want)
	}

}

func TestCreateNewConfiguration(t *testing.T) {

	app := app.Server{}
	app.Initialize()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/v1/bdi_test_service/configurations", nil)
	app.CreateNewConfigurationRoute(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	got := strings.TrimSpace(string(body))
	want := `{"1":{"sourcefolder":"","datasetname":""},"default":{"sourcefolder":"/home/data","datasetname":"hacker"}}`
	if got != want {
		t.Errorf("PATCH to CreateNewConfigurationRoute: GOT: %s WANT: %s", got, want)
	}
}

func TestGetFilesFromDefaultConfiguration(t *testing.T) {

	apiUrl := "/v1/bdi_test_service/configurations/default/files"
	data := url.Values{}
	data.Set("sourcefolder", "/Users/igo/home/data")
	data.Set("dataset", "hacker")
	app := app.Server{}
	app.Initialize()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, apiUrl, strings.NewReader(data.Encode()))
	app.ShowFilesRoute(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
