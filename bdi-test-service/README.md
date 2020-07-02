# Oxpecker Tools

Repo for test tools used for the Oxpecker/BDI project

## BDI Test Service

This service will be used in conjunction with automated tests for QABDI.

### Building the docker image and running for local deployment and testing

Build the image. Optionally pass ```-t``` to create a custom tag
```docker build . ``` 

Running the docker image:

```docker run -p 8081:8081 <image_id> ```

This will run the docker container on `localhost:8081` 

### Data Appends (Local Storage)

API endpoints to help with testing of Data Appends - dataset stored on local storage (i.e. EFS)
The following endpoints are available:

*List all configurations*

```
# GET /sta/v1/bdi_test_service/configurations
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations
```

```
Response: 200OK
Response Body:

{
  "default": {
    "sourcefolder": "/home/data",
    "datasetname": "hacker"
  }
}
```

*Create new configuration*

```
# PATCH /sta/v1/bdi_test_service/configurations
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations

Request Body:

{
  "sourcefolder": "/home/data",
  "datasetname": "tpch100"
}

```

```
Response: 202 Accepted
Response Body:

{
  "default": {
    "sourcefolder": "/home/data",
    "datasetname": "hacker"
  },
  "latest": {
    "sourcefolder": "/home/data",
    "datasetname": "tpch100"
  }
}
```

*Show configuration by ID*

Pass the configuration ID to the endpoint. For example passing ```latest``` to configID below will generate the following response

```
# GET /sta/v1/bdi_test_service/configurations/{configID}
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations/{configID}
```
```
Response: 200 OK
Response Body:

{
  "sourcefolder": "/home/data",
  "datasetname": "tpch-mini"
}
```

*Delete configuration*

Pass the configuration ID to the endpoint. DEFAULT configuration cannot be deleted.

```
# DELETE /sta/v1/bdi_test_service/configurations/{configID}
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations/{configID}
```
For example if deleting the ```latest``` configID, response will look as follows (with ```default``` remaining in below example):

```
Response: 202 Accepted
Response Body:

{
  "default": {
    "sourcefolder": "/home/data",
    "datasetname": "hacker"
  }
}
```
*Show files for given configuration*

This will show all the files for a given configuration

```
# GET /sta/v1/bdi_test_service/configurations/{configID}/files
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations/{configID}/files
```

```
Response: 200 OK
Response Body:
{
  "files": [
    "filename1.parquet",
    "filename2.parquet",
    "filename3.parquet"
  ]
}
```

*Rename files for a given configuration*

This will rename files for a given configuration. Pass as many filenames as needed.

```
# PATCH /sta/v1/bdi_test_service/configurations/{configID}/files
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations/{configID}/files

Request Body:

[

{
	"oldFileName": "filename",
     "newFileName": "new_filename"
},
	{
	"oldFileName": "filename",
     "newFileName": "new_filename"
}
	]
```

```
Response: 200 OK
Response Body:
Filenames changed
```

*Delete output*

Delete the out.

```
# DELETE /sta/v1/bdi_test_service/configurations/{configID}/output
curl https://<bd-test-service-address>/sta/v1/bdi_test_service/configurations/{configID}/output
```
Will delete the output folder in a bdi-cluster:

```
Response: 202 Accepted

```