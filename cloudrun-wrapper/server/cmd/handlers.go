package main

import (
	"bufio"
	"cloudrun/pkg/logging"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type TestInfo struct {
	TestSuite   string `json:"testsuite"`
	BaseUrl     string `json:"baseurl"`
	Parallelism string `json:"parallelism"`
}

type encoder interface {
	setEnvsAndEncode() error
}

//helper func
func (app *ServerApplication) setEnvsAndEncode(w http.ResponseWriter, r *http.Request) error {
	app.Log = logging.NewLogMultiWriter(true, os.Stdout, w)

	if err := json.NewDecoder(r.Body).Decode(&app.TestInfo); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.Log.PrintWarn("Could not decode JSON")
		return err
	}
	err := os.Setenv("SUITE", app.TestInfo.TestSuite)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.Log.PrintFatal("Could not read env variable for SUITE")
		return err
	}
	err = os.Setenv("BASEURL", app.TestInfo.BaseUrl)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.Log.PrintFatal("Could not read env variable for BASEURL")
		return err
	}
	err = os.Setenv("PARALLELEXECUTION", app.TestInfo.Parallelism)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.Log.PrintFatal("Could not read env variable for PARALLELEXECUTION")
		return err
	}

	return nil
}

func (app *ServerApplication) home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Method Now Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//TODO: Landing page
	w.Write([]byte("Wrong Route, check your path"))
}

//TODO:look at better context use. Though CloudRun will time-out after 15 minutes
//including control-c to end tests. Or send cancellation via a route
func (app *ServerApplication) run(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text")
	if r.Method == "POST" {
		if err := app.setEnvsAndEncode(w, r); err != nil {
			return
		}
	}

	cmd := exec.Command("/bin/bash", app.DockerEntrypoint)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		app.Log.PrintFatal(err)
	}

	if err := cmd.Start(); err != nil {
		app.Log.PrintFatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println("system-tests--> ", scanner.Text())
		fmt.Fprintf(w, "\nsystem-tests--> %s", scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		app.Log.PrintWarn(err)
	}

	app.Log.Infoprint("tests completed, check output for more information")

}

func (app *ServerApplication) setTestInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	app.Log = logging.NewLogMultiWriter(true, os.Stdout, w)
	if err := app.setEnvsAndEncode(w, r); err != nil {
		return
	}

	app.Log = logging.NewLogMultiWriter(false, os.Stdout)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Test Info Updated\n")
	app.Log.Infoprintf("Test Info Updated: Test Suite: %s, BaseURL %s", app.TestInfo.TestSuite, app.TestInfo.BaseUrl)
}

func (app *ServerApplication) getTestInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	app.Log = logging.NewLogMultiWriter(true, os.Stdout, w)

	if err := json.NewEncoder(w).Encode(&app.TestInfo); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.Log.PrintWarn("Could not endcode JSON")
		return
	}
}
