/*
 * @Author: dongtao
 * @Date: 2017-04-21 17:55:04
 * @Desc: mysql连接池
 * @Last Modified by: dongtao
 * @Last Modified time: 2017-04-21 18:19:09
 */
package utils

import (
	"database/sql"

	"go.app.dt.com/conf"
)

// mysql连接池
type DbStore struct {
	pool chan *sql.DB
}

// 创建连接池
func (d *DbStore) New() {
	d.pool = make(chan *sql.DB, 50)
	for i := 0; i < 50; i++ {
		db, err := Connect(conf.XMmall)
		checkErr(err)
		d.pool <- db
	}
}

// 获取连接池
func (d *DbStore) GetConn() (db *sql.DB) {
	// fmt.Println("d.pool:")
	// fmt.Println(d.pool)
	if d.pool == nil {
		d.New()
	}
	db = <-d.pool
	return
}

// 返回连接池
func (d *DbStore) RetConn(db *sql.DB) {
	d.pool <- db
	return
}
