package db

import (
	"database/sql"
	"os"

	"github.com/Moonyongjung/xpla-anchor/app"
	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/Moonyongjung/xpla-anchor/util"
	_ "github.com/go-sql-driver/mysql"
)

const (
	sqlFilePath = "./gw/db/sql/"
)

// Anchor DB initialization
func DbInit() error {
	dbConf := app.AppFile().Get().Config.Anchor.DB
	dataSource := dbConf.DBUserName + ":" + dbConf.DBPassword + "@tcp(" + dbConf.DBHost + ":" + dbConf.DBPort + ")/"

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	isDbExist, err := dbExist(db, dbConf.DBName)
	if err != nil {
		return err
	}

	if !isDbExist {
		runCreateSql(db, "db")
	}

	_, err = db.Exec("USE " + dbConf.DBName)
	if err != nil {
		return err
	}

	isErrTableExist, err := tableExist(db, "errLog")
	if err != nil {
		return err
	}

	if !isErrTableExist {
		runCreateSql(db, types.ErrLogTable+"_table")
	}

	isInfoTableExist, err := tableExist(db, "infoLog")
	if err != nil {
		return err
	}

	if !isInfoTableExist {
		runCreateSql(db, types.InfoLogTable+"_table")
	}

	logIndexInit(db)

	types.Db = db
	return nil

}

func dbExist(db *sql.DB, dbName string) (bool, error) {
	return checkExist(db, dbName, "databases")
}

func tableExist(db *sql.DB, tableName string) (bool, error) {
	return checkExist(db, tableName, "tables")
}

func checkExist(db *sql.DB, checkName string, checkType string) (bool, error) {
	isExist := true

	result, err := db.Query("show " + checkType)
	if err != nil {
		return !isExist, err
	}
	defer result.Close()

	var dbs string
	for result.Next() {
		result.Scan(&dbs)
		if dbs == checkName {
			return isExist, nil
		}
	}

	return !isExist, nil
}

func runCreateSql(db *sql.DB, name string) error {
	create, err := os.ReadFile(sqlFilePath + name + "_create.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(create))
	if err != nil {
		return err
	}

	return nil
}

func logIndexInit(db *sql.DB) {
	var errIndex string
	var infoIndex string

	util.InitLogIndex().UseDb()

	queryErrResult, _ := db.Query("select index_id from errLog order by cast(index_id as signed) desc limit 1")
	defer queryErrResult.Close()

	for queryErrResult.Next() {
		queryErrResult.Scan(&errIndex)
	}

	if errIndex == "" {
		util.InitLogIndex().NewErrLogIndex("0")
	} else {
		util.InitLogIndex().NewErrLogIndex(errIndex)
	}

	queryInfoResult, _ := db.Query("select index_id from infoLog order by cast(index_id as signed) desc limit 1")
	defer queryInfoResult.Close()

	for queryInfoResult.Next() {
		queryInfoResult.Scan(&infoIndex)
	}

	if infoIndex == "" {
		util.InitLogIndex().NewInfoLogIndex("0")
	} else {
		util.InitLogIndex().NewInfoLogIndex(infoIndex)
	}
}
