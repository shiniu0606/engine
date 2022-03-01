package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	base "github.com/shiniu0606/engine/core/base"
)

const DEFAULT_DB_DRIVER = "mysql"
const TEST_DB_DRIVER = "sqlite3"

//db实例
var gdb *gorm.DB

type Action func(db *gorm.DB) (interface{}, error)

//InitDB 初始化数据库连接
func InitDB(driver, dbConnStr string) {
	db, err := gorm.Open(driver, dbConnStr)
	if err != nil {
		base.LogError("db init err：" + err.Error())
		return
	}
	db.DB().SetMaxIdleConns(100) //SetMaxIdleConns用于设置闲置的连接数。设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetMaxOpenConns(0)   //SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制(不建议使用)，设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetConnMaxLifetime(60 * time.Second)
	gdb = db
}

//带连接池设置的DB对象
func InitDBForPoll(driver, dbConnStr string, maxConn int, idleConn int) {
	db, err := gorm.Open(driver, dbConnStr)
	if err != nil {
		base.LogError("db init error：" + err.Error())
		return
	}
	db.DB().SetMaxIdleConns(idleConn)
	db.DB().SetMaxOpenConns(maxConn)
	gdb = db
}

//GetDB 获取数据库连接对象
func GetDB() *gorm.DB {
	return gdb
}

//CloseDB 关闭数据库连接
func CloseDB() {
	gdb.Close()
}

//Execute 多表操作，使用事务处理
func Execute(action Action) (obj interface{}, err error) {
	if action == nil {
		return nil, errors.New("action can't not be nil")
	}

	obj = nil
	//开启事务
	tran := GetDB().Begin()
	if err := tran.Error; err != nil {
		return nil, err
	}

	//引入这个方法的目的是action出现panic（比如除零，数组越界，但依然捕获不到栈溢出）的时候，能够手动回滚事务，而无需等待数据库事务超时
	defer func() {
		if r := recover(); r != nil {
			//如果action发生panic错误时执行
			tran.Rollback()
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	if obj, err = action(tran); err != nil {
		tran.Rollback()
		return nil, err
	}
	tran.Commit()
	return obj, nil
}

//MapScan 从sql.Rows映射到map
func MapScan(r *sql.Rows, dest map[string]interface{}) error {
	columns, err := r.Columns()
	if err != nil {
		return err
	}
	values := make([]interface{}, len(columns))
	for i, _ := range values {
		values[i] = new(interface{})
	}
	err = r.Scan(values...)
	if err != nil {
		return err
	}
	for i, column := range columns {
		value := *(values[i].(*interface{}))
		switch v := value.(type) {
		case []byte:
			value = string(v)
		case time.Time:
			//格式化时间
			value = v.Format("2006-01-02 15:04:05")
		case nil:
			//数据库null，返回空字符串
			value = ""
		case int64:
			value = strconv.FormatInt(v, 10)
		default:
			//后续有其他类型需要特殊处理，可以直接添加case处理
		}
		dest[column] = value
	}
	return r.Err()
}

//Query 通用查询方法,返回[]map
//queryStr	查询语句
//values	参数
//返回值：数据集合，错误
// models.Query("SELECT * FROM tablename where id = ?",2)
func Query(queryStr string, values ...interface{}) ([]map[string]interface{}, error) {
	rows, err := GetDB().Raw(queryStr, values...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dataList := []map[string]interface{}{}

	for rows.Next() {
		dest := make(map[string]interface{})
		if scanErr := MapScan(rows, dest); scanErr != nil {
			base.LogError(err.Error())
			return nil, scanErr
		}

		dataList = append(dataList, dest)
	}

	return dataList, nil
}

//QueryForPage 通用分页查询方法
//pageIndex	页码
//pageSize	记录数
//queryStr	查询语句
//values	参数
//返回值：数据集合，总记录数，错误
// models.QueryForPage(1,10,"SELECT * FROM tablename where id = ?",2)
func QueryForPage(pageIndex, pageSize int64, queryStr string, values ...interface{}) ([]map[string]interface{}, int64, error) {
	//查询分页数据
	skipCount := (pageIndex - 1) * pageSize
	pageStr := fmt.Sprintf("SELECT * FROM (%s) TMP LIMIT %d,%d", queryStr, skipCount, pageSize)

	datalist, err := Query(pageStr, values...)
	if err != nil {
		return nil, 0, err
	}

	//查询数据条数
	countStr := fmt.Sprintf("SELECT COUNT(1) total FROM (%s) TMP", queryStr)
	totalMap, err := Query(countStr, values...)
	if err != nil {
		return nil, 0, err
	}

	total, err := strconv.ParseInt(totalMap[0]["total"].(string), 10, 64)
	if err != nil {
		return nil, 0, err
	}

	return datalist, total, nil
}

//返回查询行数
func QueryCount(queryStr string, values ...interface{}) (int64, error) {
	//查询数据条数
	countStr := fmt.Sprintf("SELECT COUNT(1) total FROM (%s) TMP", queryStr)
	totalMap, err := Query(countStr, values...)
	if err != nil {
		return 0, err
	}

	total, err := strconv.ParseInt(totalMap[0]["total"].(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return total, nil
}
