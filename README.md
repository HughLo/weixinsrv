weixinsrv
=========

Codes for my weixin public account

about submodules
----------------
Currently this project uses github.com/go-sql-driver/mysql to interact with mysql database. To manage the dependency, the git submodule functions are used. After the project is cloned into the working copy, two more commands need to be executed:
```
git submodule init
git submodule update
```
how to build and run
--------------------
```go run build.go```
is used to build and install main package. The building result will be installed into the ./bin path.
```go run start_server.go```
is used to launch the server.