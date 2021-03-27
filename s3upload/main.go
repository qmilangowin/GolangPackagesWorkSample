package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"s3upload/helpers"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/namsral/flag"
)

//Application ...
type Application struct {
	awsProfile string
	bucket     string
	region     string
	dir        string
	upload     bool
	list       bool
	json       bool
}

//S3Uploader will allow us to mock uploads in the test by implementing an Upload function or method
type S3Uploader interface {
	Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

func main() {

	var (
		awsProfile string
		bucket     string
		region     string
		dir        string
		upload     bool
		list       bool
		json       bool
		output     io.Writer
	)

	flag.StringVar(&awsProfile, "aws_default_profile", "", "Set your AWS Default profile via flag or export to env variable: AWS_DEFAULT_PROFILE")
	flag.StringVar(&bucket, "bucket", "", "Set bucket via flag or env variable export: export BUCKET")
	flag.StringVar(&region, "aws_region", "us-east-1", "Set region via flag or env variable export: export AWS_REGION, default set to: us-east-1")
	flag.StringVar(&dir, "dir", "", "Directory to upload. export DIR or pass via flag")
	flag.BoolVar(&upload, "upload", false, "pass --upload flag to trigger upload")
	flag.BoolVar(&list, "list", false, "pass --list flag to list contents of bucket")
	flag.BoolVar(&json, "json", false, "include --json flag to print JSON output")

	flag.Parse()

	//values from flags
	app := Application{
		awsProfile,
		bucket,
		region,
		dir,
		upload,
		list,
		json,
	}

	//set output, in this case to os.Stdout, but can take any Writer interface
	output = os.Stdout

	if err := app.run(output); err != nil {

		fmt.Println("Ensure your variables/flags are set: ")
		flag.PrintDefaults()
	}

}

func (app *Application) run(w io.Writer) error {

	var err error
	switch {

	case app.upload:

		if app.dir == "" {
			err = errors.New("set the directory flag")
			fmt.Fprintln(w, "\033[31mError occurred:\033[0m", err)
		}

		if _, err := os.Stat(app.dir); os.IsNotExist(err) {
			err = errors.New("Directory does not exist")
			fmt.Fprintln(w, "\033[31mError occurred:\033[0m", err)
		}

		if err = app.s3UploadDir(w); err != nil {
			fmt.Fprintln(w, "\033[31mError occurred:\033[0m", err)
		}
	case app.list:
		if err = app.s3ListBucketObjects(w); err != nil {
			fmt.Fprintln(w, "\033[31mError occurred:\033[0m", err)
		}
	default:
		fmt.Println("Ensure your variables/flags are set. \nDon't forget to include the '--upload' flag to trigger upload")
		flag.PrintDefaults()

	}

	return err
}

//s3Session creates an AWS s3 Session
func (app *Application) s3Session() (*session.Session, error) {

	session, err := session.NewSessionWithOptions(session.Options{
		Profile: app.awsProfile,
		Config: aws.Config{
			Region: aws.String(app.region),
		},
	})

	return session, err
}

func (app *Application) s3UploadDir(w io.Writer) error {
	session, err := app.s3Session()
	if err != nil {
		fmt.Fprintln(w, "Could not S3 create session: ", err)
		return err
	}

	uploader := s3manager.NewUploader(session, func(u *s3manager.Uploader) {
		u.PartSize = 100 * 1024 * 1024
		u.LeavePartsOnError = false
		u.Concurrency = 100
	})

	fileList := []string{}

	filepath.Walk(app.dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	for _, pathOfFile := range fileList[1:] {
		if err := app.s3UploadFile(pathOfFile, session, w, uploader); err != nil {
			return err
		}
	}
	fmt.Fprintln(w, "\033[34mDone!")
	return nil
}

//s3UploadFile uploads the actual file.
func (app *Application) s3UploadFile(pathOfFile string, session *session.Session, w io.Writer, s3svc S3Uploader) error {

	stat, err := os.Stat(pathOfFile)
	if err != nil {
		fmt.Fprintln(w, "Error getting file size: ", err)
		return err
	}

	//open file and add to buffer
	file, err := os.Open(pathOfFile)
	if err != nil {
		fmt.Fprintln(w, "Error opening file: ", err)
		return err
	}
	defer file.Close()
	path := file.Name()
	fileSize := stat.Size()

	progReader := helpers.NewProgressReader(file, fileSize)
	fmt.Fprintln(w, "\033[32mUploading: ", pathOfFile)
	progReader.ProgBar.Format("\033[32m\x00=\x00>\x00-\x00]")
	progReader.ProgBar.Start()
	defer progReader.ProgBar.Finish()

	_, err = s3svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(app.bucket),
		Key:    aws.String(path),
		Body:   progReader, //fileBytes
	})

	if err != nil {

		//directories in S3 are initially created as zero-sized objects which causes
		//an error in the progress bar due to a zero-byte read prior to upload. We can ignore this error below.
		//rest of the contents are properly calculated once the directory is created
		//and upload starts. Other errors we will catch. So create a message to user during upload
		//that directory is being created instead.
		if strings.Contains(err.Error(), "BodyHash") {
			fmt.Fprintln(w, "\033[34m Creating Directory...\033[32m")
			return nil
		}

		fmt.Fprintln(w, "\033[31mError uploading file \033[0m", err)
		return err

	}

	return nil

}

//s3ListBucketObjects lists all the objects in the bucket (including directory tree structures)
//in the bucket
func (app *Application) s3ListBucketObjects(w io.Writer) error {
	session, err := app.s3Session()
	if err != nil {
		fmt.Fprintln(w, "Error getting file size: ", err)
		return err
	}
	svc := s3.New(session)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(app.bucket),
		MaxKeys: aws.Int64(5000), //hard-coded to max 5000 objects to return. Should be enough
	}

	result, err := svc.ListObjectsV2(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Fprintln(w, fmt.Errorf("No such bucket - please check name or credentials"))
				fmt.Fprintln(w, aerr.Error())
			}
		} else {

			fmt.Fprintln(w, err.Error())
		}
		return err
	}

	if app.json {
		fmt.Fprintln(w, *result)

	} else {
		for _, v := range result.Contents {
			fmt.Fprintln(w, *v.Key)
		}
	}

	return nil
}
