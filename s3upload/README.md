## S3Upload 

Used internally at work as part of CI/CD pipeline

Before running ensure to export the following variables which can be found
in `~/.aws/credentials`

`export AWS_DEFAULT_PROFILE=<profile>`  
`export BUCKET=<AWS_BUCKET>`  

Alternatively these can be set via command-line flags. If they have been exported, no flags have to be passed.

See `--help` for how to set these


### Building 

`make build`


### Running Tests


`make test`


### Usage

Run `s3upload --help` for list of available commands

To upload directory:

```console
s3upload --upload --dir <directory-to-upload>
```

Alternatively, export the directory to DIR env variable:

```console
 export DIR=<directory-to-upload>
 s3upload --upload
```

List directory/files in bucket:

```s3upload --list```

Note, this will list all objects in the bucket and not latest uploaded directory.
To see those, use `aws s3 ls` command for the specific directory.
