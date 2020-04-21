
PROJECT_PATH=$(pwd)
echo $GOPATH

#cd $GOPATH

echo install protobuf
go get github.com/golang/protobuf/

echo install mysql
go get github.com/go-sql-driver/mysql

echo install gorm
github.com/jinzhu/gorm
