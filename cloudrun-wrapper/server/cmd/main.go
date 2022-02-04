package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	logging "cloudrun/pkg/logging"
)

type ServerApplication struct {
	Log              *logging.Logger
	DockerEntrypoint string
	TestInfo         TestInfo
}

func (app *ServerApplication) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/run", app.run)
	mux.Handle("/run", app.router("/run", mux))
	mux.Handle("/testinfo", app.router("/testinfo", mux))
	return secureHeaders(mux)
}

func initSetup() *ServerApplication {
	appLog := logging.NewLogMultiWriter(true, os.Stdout, ioutil.Discard)
	suite := "workinprogress"
	baseUrl := "<internal-url-removed>"
	parallelism := "1"
	os.Setenv("SUITE", suite)
	os.Setenv("BASEURL", baseUrl)
	os.Setenv("PARALLELEXECUTION", parallelism)
	suite = os.Getenv("SUITE")
	baseUrl = os.Getenv("BASEURL")
	parallelism = os.Getenv("PARALLELEXECUTION")
	dockerEntryPoint := "docker-entrypoint.sh"

	app := &ServerApplication{
		Log:              appLog,
		DockerEntrypoint: dockerEntryPoint,
		TestInfo: TestInfo{
			TestSuite:   suite,
			BaseUrl:     baseUrl,
			Parallelism: parallelism,
		},
	}
	return app
}

func main() {

	app := initSetup()
	srv := &http.Server{
		Handler:      app.routes(),
		Addr:         ":8080",
		WriteTimeout: 30 * time.Minute,
		ReadTimeout:  30 * time.Minute,
	}

	//start the server in own go-routine
	go func() {
		fmt.Println("Starting server on port: 8080")
		app.Log.Infoprintf("Server started on: %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			app.Log.PrintFatal("Shutting down server", err)
		}
	}()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	log.Printf("Received Terminate, graceful shutdown %s", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(tc)

}
