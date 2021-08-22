# file-uploader

File-uploader will perform scanning virus of files, uploading files in different bucket system like local, s3 etc
uploaded via http requests with `multipart/form-data`. If a virus exists, it will respond with infected file information
and also if clean it will upload to specified bucket.

## Running with docker

For running the service in docker container use ```make run-docker-compose```.

## Dependency injection

The project uses [wire](https://github.com/google/wire/blob/main/docs/guide.md) for building a dependency tree without
pain.

## Local run

The easiest and straight forward way is to have local setup of Go. There is how the GoLand build should be configured:

|Property         |Value                            |
|-----------------|---------------------------------|
|Package path     | ABSOLUTE_LOCAL_PATH...file-uploader/src |
|Output directory | ABSOLUTE_LOCAL_PATH...file-uploader/bin |
|Working directory| ABSOLUTE_LOCAL_PATH...file-uploader/    |
|Environment      | Point to .dev.env (please install ["EnvFile"](https://plugins.jetbrains.com/plugin/7861-envfile) plugin for GoLand) |

In order to generate mocks, use [mockery](https://github.com/vektra/mockery)

### Environment Variables

|Property         |Example Values                            |
|-----------------|---------------------------------|
|ENV_NAME|local|
|HTTP_PORT|8085|
|CLAMAV_SOCKET_URL|unix:/tmp/clamd.socket|
|FILE_UPLOADER_API_KEY|asdfjkasdf,sadfsad|
|DB_HOST|localhost|
|DB_PORT|3306|
|DB_NAME|file_upload|
|DB_USER|root|
|DB_PASS|password|
|DB_PASS_FILE| |
|FILE_STORAGE_DRIVER|s3|
|FILE_STORAGE_DISABLE_SSL|false|
|FILE_STORAGE_REGION|eu-central-1|
|FILE_STORAGE_PROFILE|xxx-sandbox|
|FILE_STORAGE_ACCESS_KEY| |
|FILE_STORAGE_SECRET| |
|S3_FORCE_PATH_STYLE| empty means true | 
|NEWRELIC_APM_LICENSE| XXXXXXXXX | 
|NEWRELIC_APPNAME_PREFIX| XXXXXXXXX | 
|NEWRELIC_HOST_DISPLAY_NAME| XXXXXXXXX | 

### Dependencies

Please install clamav antivirus and run clamad deamon service so that the application can communication with clamd. Also
don't forget to update virus database for more information please follow the
link [Clamav](https://www.clamav.net/documents/installing-clamav-on-unix-linux-macos-from-source)

Cookbook
============

Token Generation: To generate the token for downloading or streaming the file, you can generate JWT with validity and must be signed with your api key and file id together (e.g: <API_KEY><FILE_ID>)


## REST Endpoints

--------------------

### List Files

```
  GET /files
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### QueryParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|name         |             | <search_string>|
|created_date         |             | <YYYY-MM-DD>|
|offset         |             | Any integer value, default is 0|
|limit         |             | Any integer value, default is 500|

##### Response (Code: 200 or 500):

```json
[
  {
    "id": "957ad83c-3e7f-494e-89b7-717cad82103d",
    "name": "A-Z-Infoservices and 140 others.vcf",
    "meta_data": {
      "size": 25919,
      "mime_type": "text/x-vcard"
    },
    "owner_id": "asdfjkasdf",
    "bucket_path": "xxx-file-upload/test/957ad83c-3e7f-494e-89b7-717cad82103d.vcf",
    "provider": "s3",
    "created_at": "2021-05-10T06:37:42Z",
    "expired_at": "2021-05-12T19:21:04Z"
  }
]
```

### Create Files

```
  POST /files
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### BODY (TYPE SHOULD BE multipart/form-data):

|Name         |Type|Value     |Required| Example
|-------------|----|---------|-----|---|
|bucket_path   | string |  <BUCKET_PATH>| x|xxx-bucket or xxx-bucket/\<POSTFIX>|
|files[0].file | file |  <FILE_TO_UPLOAD>| x||
|files[0].id   | string |  <UUID_4> or will be auto generated| ||
|files[1].file | file|  <FILE_TO_UPLOAD>| x||
|files[1].id   | string|  <UUID_4> or will be auto generated| |b69724a3-8488-402c-a16e-fc0da0ea6832|
|expired_at    | date-time|  <VALID_ISO_DATE_TIME>| ||

##### Response (Code: 201 or  207):

Response 207:

```json
[
  {
    "file_name": "eicarcom2.zip",
    "success": false,
    "message": "Virus found with defination Win.Test.EICAR_HDB-1",
    "virus": true
  },
  {
    "id": "4825181a-f908-4a7b-88b4-e00eda2adb76",
    "file_name": "A-Z-Infoservices and 140 others.vcf",
    "success": true,
    "message": "",
    "virus": false
  }
]
```

Response 201, all file uploaded successfully:

```json
[
  {
    "id": "4825181a-f908-4a7b-88b4-e00eda2adb76",
    "file_name": "A-Z-Infoservices and 140 others.vcf",
    "success": true,
    "message": "",
    "virus": false
  }
]
```

### Get a File

```
  GET /files/:id
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### PathParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|id        |   x | <file_UUID>

##### Response (Code: 200 or 500):

```json
{
  "id": "62ce012e-44b7-42f3-b017-ebbbfaa6bede",
  "name": "A-Z-Infoservices and 140 others.vcf",
  "meta_data": {
    "size": 25919,
    "mime_type": "text/x-vcard"
  },
  "owner_id": "asdfjkasdf",
  "bucket_path": "xxx-file-upload/test/62ce012e-44b7-42f3-b017-ebbbfaa6bede.vcf",
  "provider": "s3",
  "created_at": "2021-05-07T18:08:02Z",
  "expired_at": "2021-05-12T19:21:04Z"
}
```

### Delete a File

```
  DELETE /files/:id
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### PathParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|id        |   x | <file_UUID>

##### Response (Code: 200 or 500)

### Generating Tokens

```
  POST /files/tokens
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### Body:

|Name      |Type   |Required     | Example Value |
|----------|---|-------------|---------|
|ids       | Array of file_UUID|   x | ["b69724a3-8488-402c-a16e-fc0da0ea6832"]
|expired_at       | date-time|  default 10 years  | 2021-05-19T12:48:32.350848+02:00

##### Response (Code: 200 or 400 or 500)

```json
{
  "b69724a3-8488-402c-a16e-fc0da0ea6832": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjEzNzYxMTksImZpbGVfaWQiOiJiNjk3MjRhMy04NDg4LTQwMmMtYTE2ZS1mYzBkYTBlYTY4MzIifQ.YSjK5MSoaUrtIK7_u0PVMd2UYWcDxSsWDUabOV69WwM"
}
```

### Stream

```
  GET /files/:id/stream
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### PathParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|id        |   x | <file_UUID>

##### QueryParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|token        |   x [If x-api-key is not present in header] | JWT token signed with \<api-key\>\<file_id\>|

##### Response (Code: 200 or 500): File stream

### Download

```
  GET /files/:id/download
```

##### Headers:

|Name         |Value     |Required|
|-------------|-------------|-----|
|x-api-key        |  <PROVIDED_API_KEY>| x|

##### PathParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|id        |   x | <file_UUID>

##### QueryParams:

|Name         |Required     | Example Value |
|-------------|-------------|---------|
|token        |   x [If x-api-key is not present in header] | JWT token signed with \<api-key\>\<file_id\>|
|disposition  |             |inline or attachment| |

##### Response (Code: 200 or 500): File content based on the disposition type, default is attachment

### Health

```
  GET /health
```

##### Response (Code: 200 or 404):

```json
"OK"
```

### Info

```
  GET /scanners/info
```

##### Response (Code: 200):

```json
{
  "file_uploader_version": "master",
  "scan_server_url": "unix:/tmp/clamd.socket",
  "ping_result": "Connected to server OK",
  "scan_server_version": "ClamAV 0.103.2/26172/Sun May 16 13:13:51 2021",
  "test_scan_virus": "Status: FOUND; Virus: true; Description: Win.Test.EICAR_HDB-1",
  "test_scan_clean": "Status: CLEAN; Virus: false"
}
```

This method will return JSON giving the current status of File Uploader and its connection to ClamAV.

### ScanFiles

```
  POST /scanners/files
```

##### RequestBody:

`Body should be multipart/form-date with files`

##### Response (Code: 200 or 418):

```json
{
  "success": false,
  "files": [
    {
      "Status": "FOUND",
      "Virus": true,
      "Description": "Win.Test.EICAR_HDB-1",
      "error": false,
      "file_name": "eicarcom2.zip",
      "message": ""
    },
    {
      "Status": "CLEAN",
      "Virus": false,
      "Description": "",
      "error": false,
      "file_name": "A-Z-Infoservices and 140 others.vcf",
      "message": ""
    }
  ],
  "message": ""
}
```

### ScanUrls

```
  POST /scanners/urls
```

##### RequestBody:

```json
{
  "urls": [
    "https://secure.eicar.org/eicar.com.txt",
    "https://secure.eicar.org/eicar_com.zip"
  ]
}
```

##### Response (Code: 200 or 418):

```json
{
  "success": false,
  "files": [
    {
      "Status": "FOUND",
      "Virus": true,
      "Description": "Win.Test.EICAR_HDB-1",
      "error": false,
      "file_name": "https://secure.eicar.org/eicar.com.txt",
      "message": ""
    },
    {
      "Status": "FOUND",
      "Virus": true,
      "Description": "Win.Test.EICAR_HDB-1",
      "error": false,
      "file_name": "https://secure.eicar.org/eicar_com.zip",
      "message": ""
    }
  ],
  "message": ""
}
```

---
