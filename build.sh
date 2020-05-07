protoc --go_out=./ ./proto/center_server.proto
protoc --go_out=./ ./proto/user_front.proto

go build -o dbcreate ./run/dbtable_create.go
go build -o centerserver ./run/centerserver_create.go
go build -o userfrontserver ./run/userfront_create.go
