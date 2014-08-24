package wxsrv

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type ExecResult sql.Result
type QueryResult struct {
	*sql.Rows
}

type SqlQuery struct {
	tableName string
	colNames  []string
	whereCond string
	vals      []string
}

type DBMgr struct {
	SqlDB     *sql.DB
	QueryInfo SqlQuery
}

func CreateDBMgr(cs string) *DBMgr {
	db, err := sql.Open("mysql", cs)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &DBMgr{SqlDB: db}
}

func (db *DBMgr) Close() {
	db.SqlDB.Close()
}

func (db *DBMgr) CreateDatabase(name string) error {
	cmdString := fmt.Sprintf("create database if not exist %s", name)
	_, err := db.SqlDB.Exec(cmdString)
	return err
}

func (db *DBMgr) UseDB(name string) error {
	cmdString := fmt.Sprintf("use %s", name)
	_, err := db.SqlDB.Exec(cmdString)
	return err
}

func (db *DBMgr) Cols(colNames []string) *DBMgr {
	db.QueryInfo.colNames = colNames
	return db
}

func (db *DBMgr) Table(tableName string) *DBMgr {
	db.QueryInfo.tableName = tableName
	return db
}

func (db *DBMgr) Where(cond string) *DBMgr {
	db.QueryInfo.whereCond = cond
	return db
}

func (db *DBMgr) Values(vals []string) *DBMgr {
	db.QueryInfo.vals = vals
	return db
}

func (db *DBMgr) Query() (*QueryResult, error) {
	cs := strings.Join(db.QueryInfo.colNames, ",")
	qs := fmt.Sprintf("select %s from %s where %s", db.QueryInfo.tableName, cs, db.QueryInfo.whereCond)
	r, err := db.SqlDB.Query(qs)
	return &QueryResult{r}, err
}

func (db *DBMgr) RawQuery(qs string) (*QueryResult, error) {
	r, err := db.SqlDB.Query(qs)
	return &QueryResult{r}, err
}

func (db *DBMgr) Insert() (ExecResult, error) {
	vs := strings.Join(db.QueryInfo.vals, ",")
	cs := strings.Join(db.QueryInfo.colNames, ",")
	qs := fmt.Sprintf(`insert into %s (%s) values (%s)`, db.QueryInfo.tableName, cs, vs)

	log.Printf("insert sql string: %s\n", qs)

	return db.SqlDB.Exec(qs)
}

func (db *DBMgr) Call(name string, params ...string) (*QueryResult, error) {
	var qs string
	if len(params) > 0 {
		qs = fmt.Sprintf("call %s(%s)", name, params)
	} else {
		qs = fmt.Sprintf("call %s()", name)
	}

	r, err := db.SqlDB.Query(qs)
	return &QueryResult{r}, err
}
