protoc --go_out=./ ./proto/center_server.proto

go build -o dbcreate ./run/dbtable_create.go
go build -o centerserver ./run/centerserver_create.go
