# Overview
  
# For Developers

## Setup

### Prerequisite
1\. Install [GIT](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
2\. Install [Go 1.11.4](https://golang.org/doc/install) or later.
3\. Set the GOPATH and add go binaries. Below is the sample config in ~/.bash_profile in mac osx.
```bash
export GOPATH=$HOME/workspace/go
export PATH=$PATH:$GOPATH/bin
```
4\. Install [dep](https://golang.github.io/dep/docs/installation.html). For Windows users, you can download from [releases](https://github.com/golang/dep/releases) page. Please take sometime reading the usage of [go dep](https://golang.github.io/dep/docs/daily-dep.html) specially in adding new dependencies part.

### Steps

1\. Checkout the [repository](https://github.com/ramirezag/gocleancode.git)
```bash 
cd $GOPATH/src
git clone https://github.com/ramirezag/gocleancode.git
```
2\. Pull the project dependencies
```bash
  $ cd gocleancode
  $ dep ensure
``` 

3\. Execute gocleancode/db/mysqlSchema.sql to your mysql db.

### Configure the APP config

**Option 1: Using environment variables**

Set the following environment variables
```bash
HOST= # For mac users, use localhost due avoid the annoying popup.
APP_PORT= # defaults to 8000 if not set.
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=
```  
  
**Option 2: Using config files**

This approach is useful when you want to run the app in commandline and would not like to set the environment variables. Note that environment variables will take **precedence** over config files.
  
1\. Copy `config.<env>.json.sample` and name it `config.local.json`  
2\. Fill the properties according to your needs.  
3\. Execute  
  ```bash
  $ cd $GOPATH/src/gocleancode
  $ export ENV=local && go run main.go
  ``` 
3\.1\. For windows users, manually set `ENV` environment variable to `local`. Then execute `go run main.go`.
  
## Run the server

Execute `go run main.go`

## Development and Unit Test Guidelines

- Read [How to Write Go Code](https://golang.org/doc/code.html).
- Prefix tests packages with `_test` to avoid `import-cycle-not-allowed` error. Eg, `package repository_test`
- **Always** program to interface so that it replace implementations should we need it - eg, mocking dependencies.
- use [mockery](https://github.com/vektra/mockery#installation) in generating mocks of interfaces.
- To avoid issues, capitalize the first character of the function in interfaces
- Steps to generate mocks  
  - cd $GOPATH/src/gocleancode/<to_dir_with_interface>
  - mockery -name <interface_name>
- Example steps to generate mocks
    ```bash
        cd repository
        mockery -name FileRepo
    ```  
- As much as possible, don't use log.Fatal or os.Exit(1) unless you're sure that no further clean is needed. Those commands will immediately terminate the app and won't invoke shutdown hooks (eg, in db.go and main.go).      

# APIs

Below are the list of available APIs exposed by this service.

| API        | Success Response | Description |
| ------------- | ------------- | ------------- |
| POST /files  | `{ "success": true, "message": "Created file with id 1." }` | Multipart Upload files. Note: Parameter `file` should be used. Eg, `<input type="file" name="file" />` |
| GET /files/{fileId}      | File Stream | Download file by file id. Use a browser to see the file. |
| DELETE /files/{fileId}      | `{ "success": true, "message": "Successfully deleted file with id 1" }` | Delete file by id. |

Sample usage  
```bash
curl -X GET http://localhost:8000/files/1
```
