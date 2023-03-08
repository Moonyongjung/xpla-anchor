package util

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Moonyongjung/xpla-anchor/types"
	"github.com/logrusorgru/aurora"
)

var logIndexOnce sync.Once
var logIndexInstance *LogIndex

type LogIndex struct {
	LogDB     bool
	LogFile   bool
	InfoIndex string
	ErrIndex  string
}

func InitLogIndex() *LogIndex {
	logIndexOnce.Do(func() {
		logIndexInstance = &LogIndex{}
	})
	return logIndexInstance
}

func (l *LogIndex) UseDb() {
	l.LogDB = true
}

func (l *LogIndex) IsUseDb() bool {
	return l.LogDB
}

func (l *LogIndex) NewErrLogIndex(Index string) {
	temp, _ := strconv.Atoi(Index)
	l.ErrIndex = strconv.Itoa(temp + 1)
}

func (l *LogIndex) NowErrLogIndex() string {
	return l.ErrIndex
}

func (l *LogIndex) AddErrLogIndex() {
	temp, _ := strconv.Atoi(l.ErrIndex)
	l.ErrIndex = strconv.Itoa(temp + 1)
}

func (l *LogIndex) NewInfoLogIndex(Index string) {
	temp, _ := strconv.Atoi(Index)
	l.InfoIndex = strconv.Itoa(temp + 1)
}

func (l *LogIndex) NowInfoLogIndex() string {
	return l.InfoIndex
}

func (l *LogIndex) AddInfoLogIndex() {
	temp, _ := strconv.Atoi(l.InfoIndex)
	l.InfoIndex = strconv.Itoa(temp + 1)
}

func LogInfo(log ...interface{}) {
	print := logTime() + G(" Anchor   ") + ToStringTrim(log, "")
	fmt.Println(print)

	if InitLogIndex().IsUseDb() {
		saveLogsDb(ToStringTrim(log, ""), 0, types.InfoLogTable)
	}
}

func LogWait(log ...interface{}) {
	print := logTime() + Y(" Waiting  ") + ToStringTrim(log, "")
	fmt.Println(print)

	if InitLogIndex().IsUseDb() {
		saveLogsDb(ToStringTrim(log, ""), 1, types.InfoLogTable)
	}
}

func LogWarning(log ...interface{}) {
	print := logTime() + " " + BgR("WARNING") + "  " + ToStringTrim(log, "")
	fmt.Println(print)

	if InitLogIndex().IsUseDb() {
		saveLogsDb(ToStringTrim(log, ""), 2, types.ErrLogTable)
	}
}

func LogErr(errType types.XGoError, errDesc ...interface{}) error {
	print := logErr("code", errType.ErrCode(), ":", errType.Desc(), "-", errDesc)
	fmt.Println(print)

	if InitLogIndex().IsUseDb() {
		var log []interface{}
		log = append(log, errType.Desc(), "-", errDesc)
		saveLogsDb(ToStringTrim(log, ""), errType.ErrCode(), types.ErrLogTable)
	}

	return errors.New(ToStringTrim(errDesc, ""))
}

func logErr(log ...interface{}) string {
	return logTime() + R(" Error    ") + ToStringTrim(log, "")
}

func logTime() string {
	return B(time.Now().Format("2006-01-02 15:04:05"))
}

func saveLogsDb(message string, code uint64, logTable string) {
	db := types.Db

	var dbExe *sql.Stmt
	var err error

	dbExe, err = db.Prepare("insert into " + logTable + " values (?, ?, ?, ?)")

	if err != nil {
		panic(err)
	}
	defer dbExe.Close()

	var index string

	if logTable == types.ErrLogTable {
		index = InitLogIndex().NowErrLogIndex()

	} else if logTable == types.InfoLogTable {
		index = InitLogIndex().NowInfoLogIndex()

	} else {
		panic("invalid log table")
	}

	message = strings.ReplaceAll(message, "[94m", "")
	message = strings.ReplaceAll(message, "[0m", "")

	_, err = dbExe.Exec(index, code, message, time.Now())
	if err != nil {
		panic(err)
	}

	if logTable == types.ErrLogTable {
		InitLogIndex().AddErrLogIndex()

	} else if logTable == types.InfoLogTable {
		InitLogIndex().AddInfoLogIndex()

	} else {
		panic("invalid log table")
	}
}

func G(str string) string {
	return aurora.Green(str).String()
}

func R(str string) string {
	return aurora.Red(str).String()
}

func BB(str string) string {
	return aurora.BrightBlue(str).String()
}

func B(str string) string {
	return aurora.Blue(str).String()
}

func Y(str string) string {
	return aurora.Yellow(str).String()
}

func BgR(str string) string {
	return aurora.BgRed(str).String()
}
