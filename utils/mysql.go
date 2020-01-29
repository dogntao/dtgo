/*
 * @Author: dongtao
 * @Date: 2017-04-21 16:55:04
 * @Desc: mysql 方法封装
 * @Last Modified by: dongtao
 * @Last Modified time: 2017-04-21 17:55:33
 */
package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"io"

	_ "github.com/go-sql-driver/mysql"
)

type Conf map[string]string

type Mysql struct {
	host    string
	user    string
	pass    string
	dbname  string
	charset string

	DbStore
	rows    *sql.Rows
	cols    []string
	RetMap  map[string]string
	RetRows []map[string]string
}

func Connect(c Conf) (db *sql.DB, err error) {
	mysql := new(Mysql)
	mysql.host = c["host"]
	mysql.user = c["user"]
	mysql.pass = c["pass"]
	mysql.dbname = c["dbname"]
	mysql.charset = c["charset"]

	// db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/dt?charset=utf8")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", mysql.user, mysql.pass, mysql.host, mysql.dbname, mysql.charset)
	db, err = sql.Open("mysql", dsn)
	checkErr(err)
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// query
func (mysql *Mysql) Query(field interface{}, table, con string, bind []string) (err error) {
	// 拼装查询字段
	// SELECT FIELD FROM TABLE WHERE CONDITION
	t := reflect.TypeOf(field)
	fieldArr := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		// 处理count
		if t.Field(i).Name == "QueryCount" {
			fieldArr[i] = "count(*) as queryCount"
		} else {
			fieldArr[i] = t.Field(i).Tag.Get("json")
		}
	}
	fieldString := strings.Join(fieldArr, ",")

	// 获取链接
	db := mysql.GetConn()
	defer mysql.RetConn(db)
	// 处理无condition
	selSql := fmt.Sprintf("SELECT %s FROM %s", fieldString, table)
	if con != "" {
		selSql = fmt.Sprintf("%s WHERE %s", selSql, con)
	}
	stmt, err := db.Prepare(selSql)
	// fmt.Println(selSql)
	// fmt.Println(bind)
	// 处理无绑定
	if len(bind) > 0 {
		// fmt.Println("bind is not nil")
		bindArr := make([]interface{}, 0)
		for _, val := range bind {
			bindArr = append(bindArr, val)
		}
		mysql.rows, err = stmt.Query(bindArr...)
	} else {
		// fmt.Println("bind is nil")
		mysql.rows, err = stmt.Query()
	}
	return
}

// fetchRow(返回结果中的第一行)
func (mysql *Mysql) FetchRow() (err error) {
	if mysql.rows.Next() {
		mysql.cols, err = mysql.rows.Columns()
		if err != nil {
			return
		}
		// fmt.Println(mysql.cols)
		scanArgs := make([]interface{}, len(mysql.cols))
		values := make([]sql.RawBytes, len(mysql.cols))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = mysql.rows.Scan(scanArgs...)
		if err != nil {
			return
		}

		mysql.RetMap = make(map[string]string)
		for i := 0; i < len(mysql.cols); i++ {
			mysql.RetMap[mysql.cols[i]] = string(values[i])
		}
		return
	}
	return io.EOF
}

// fetchAll(获取所有返回值)
func (mysql *Mysql) FetchAll() (err error) {
	// 清空数据
	mysql.RetRows = make([]map[string]string, 0)
	for {
		err = mysql.FetchRow()
		if err == io.EOF {
			err = nil
			return
		} else if err != nil {
			return
		}
		mysql.RetRows = append(mysql.RetRows, mysql.RetMap)
	}
	return nil
}

// insert(插入一行数据)
func (mysql *Mysql) Insert(tableName string, data map[string]interface{}) (lastId int64, err error) {
	// INSERT INTO TABLE(keyArr) VALUES(?,?,...)
	// 获取链接
	db := mysql.GetConn()
	defer mysql.RetConn(db)

	// 根据data拼出key和value以及bind
	keyArr := []string{}
	valArr := []string{}
	bindArr := make([]interface{}, 0)
	for key, val := range data {
		keyArr = append(keyArr, key)
		valArr = append(valArr, "?")
		bindArr = append(bindArr, val)
	}
	keyStr := strings.Join(keyArr, ",")
	valStr := strings.Join(valArr, ",")

	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, keyStr, valStr)
	result, err := db.Exec(sql, bindArr...)
	if err == nil {
		lastId, _ = result.LastInsertId()
	}
	return
}

// InsertMulti(批量插入数据)
func (mysql *Mysql) InsertMulti(tableName string, data []map[string]string) (lastId int64, err error) {
	// INSERT INTO TABLE(keyArr) VALUES (?,?,...),(?,?,...),(?,?,...)
	// 获取链接
	db := mysql.GetConn()
	defer mysql.RetConn(db)

	// 根据data拼出key和value以及bind
	keyArr := []string{}
	valArr := []string{}
	perValstr := ""
	perValArr := []string{}
	bindArr := make([]interface{}, 0)

	for key, val := range data {
		perValstr = ""
		perValArr = []string{}

		if len(val) > 0 {
			for k, v := range val {
				// 生成key
				if key == 0 {
					keyArr = append(keyArr, k)
				}
				perValArr = append(perValArr, "?")
				bindArr = append(bindArr, v)
			}
			// 生成每个value
			perValstr = "(" + strings.Join(perValArr, ",") + ")"
			valArr = append(valArr, perValstr)
		}
	}

	keyStr := strings.Join(keyArr, ",")
	valStr := strings.Join(valArr, ",")
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES %s", tableName, keyStr, valStr)
	result, err := db.Exec(sql, bindArr...)
	if err == nil {
		lastId, _ = result.LastInsertId()
	}
	return
}

// update(更新数据)
func (mysql *Mysql) Update(tableName string, upData map[string]interface{}, con string, bind []string) (affRow int64, err error) {
	// UPDATE TABLE SET keyStr WHERE conStr
	// 获取链接
	db := mysql.GetConn()
	defer mysql.RetConn(db)

	keyArr := []string{}
	bindArr := make([]interface{}, 0)
	for key, val := range upData {
		keyArr = append(keyArr, fmt.Sprintf("%s=?", key))
		bindArr = append(bindArr, val)
	}
	keyStr := strings.Join(keyArr, ",")
	for _, val := range bind {
		bindArr = append(bindArr, val)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, keyStr, con)
	result, err := db.Exec(sql, bindArr...)
	if err == nil {
		affRow, _ = result.RowsAffected()
	}
	return
}

// updateMulti(批量更新数据)
func (mysql *Mysql) UpdateMulti(tableName string, upDatas []map[string]string, id string) (affRow int64, err error) {
	/*
		UPDATE TABLE SET
			field1 = CASE
				WHEN id=1 THEN "field1_value1"
				WHEN id=2 THEN "field1_value2"
			END,
			field2 = CASE
				WHEN id=1 THEN "field2_value1"
				WHEN id=2 THEN "field2_value2"
			END
		WHERE id in(1,2)
	*/
	// 获取链接
	db := mysql.GetConn()
	defer mysql.RetConn(db)

	idArr := []string{}
	fieldStr := ""
	fieldArr := []string{}
	bindArr := make([]interface{}, 0)

	// 拼装field map
	fieldMap := make(map[string][]string)
	bindMap := make(map[string][]string)
	for _, value := range upDatas {
		for k, v := range value {
			if k == id {
				idArr = append(idArr, v)
			} else {
				// str := fmt.Sprintf("WHEN %s=%s THEN '%s'", id, value[id], v)
				str := "WHEN " + id + "=?  THEN ?"
				fieldMap[k] = append(fieldMap[k], str)
				// 生成绑定
				bindMap[k] = append(bindMap[k], value[id])
				bindMap[k] = append(bindMap[k], v)
			}
		}
	}

	// 合并field map
	for k, v := range fieldMap {
		fieldStr = k + "=" + " CASE " + strings.Join(v, " ") + " END"
		fieldArr = append(fieldArr, fieldStr)
		// 处理绑定
		for _, b_v := range bindMap[k] {
			bindArr = append(bindArr, b_v)
		}
	}

	fields := strings.Join(fieldArr, ",")
	con := id + " IN(" + strings.Join(idArr, ",") + ")"
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, fields, con)
	result, err := db.Exec(sql, bindArr...)
	if err == nil {
		affRow, _ = result.RowsAffected()
	}
	return
}
