package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var mockStdOut io.Writer = &bytes.Buffer{}

func setup() *Application {

	awsProfile := os.Getenv("AWS_DEFAULT_PROFILE")
	awsBucket := os.Getenv("BUCKET")

	app := Application{
		awsProfile: awsProfile,
		bucket:     awsBucket,
		region:     "us-east-1",
		dir:        "",
		upload:     false,
		list:       true,
		json:       true,
	}

	return &app
}

type fileUpload struct{}

func (f *fileUpload) Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	log.Println("Mock uploaded to s3")
	return &s3manager.UploadOutput{}, nil
}

func TestS3Connectivity(t *testing.T) {

	app := setup()
	_, err := app.s3Session()
	if err != nil {
		t.Error("Could not create s3 session ", err)
	}

	if err := app.run(mockStdOut); err != nil {
		fmt.Println(err)
	}

	got := fmt.Sprint(mockStdOut)
	expected := "MaxKeys"
	if !strings.Contains(got, expected) {
		t.Errorf("Expected: %s; got %s", expected, got)
	}

}

func TestMockUpload(t *testing.T) {

	mockFilePath := "/"
	f := fileUpload{}
	app := setup()
	sess, _ := app.s3Session()
	if err := app.s3UploadFile(mockFilePath, sess, mockStdOut, &f); err != nil {
		t.Errorf("s3Upload error: %s", err)
	}
}
