package dbhelp

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoonieslab/hoon_loghelper/helper/loghelp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DbInfo struct {
	GormDB *gorm.DB
	SqlDB  *sql.DB
	DbType string
}

/**
 * 실행 쿼리 로깅을 위한 Plugin 작성
 */
type ZapLoggingPlugin struct {
	PluginName string
	Desc       string
}

func (plugin *ZapLoggingPlugin) Name() string {
	plugin.PluginName = "DBQueryLog"
	plugin.Desc = ":: Use to loghelp queries to files and consoles after gorm behavior"
	return plugin.PluginName
}

func (plugin *ZapLoggingPlugin) Initialize(db *gorm.DB) (err error) {
	// 쿼리 이후
	_ = db.Callback().Create().After("gorm:after_create").Register("lshplugin:after_create", gormAfterActive)
	_ = db.Callback().Query().After("gorm:after_query").Register("lshplugin:after_query", gormAfterActive)
	_ = db.Callback().Delete().After("gorm:after_delete").Register("lshplugin:after_delete", gormAfterActive)
	_ = db.Callback().Update().After("gorm:after_update").Register("lshplugin:after_update", gormAfterActive)
	_ = db.Callback().Row().After("gorm:row").Register("lshplugin:after_row", gormAfterActive)
	_ = db.Callback().Raw().After("gorm:raw").Register("lshplugin:after_raw", gormAfterActive)
	return
}

func gormBeforeActive(db *gorm.DB) {
	fmt.Println(db.Statement)

	return
}

func gormAfterActive(db *gorm.DB) {
	executionSql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	//loghelp.Info(strings.ReplaceAll(strings.ReplaceAll(executionSql, "\n", " "), "\t", " "))
	loghelp.Info(executionSql)

	return
}

// InitGormCon initialize Database connection with query logging plugin.
func InitGormCon(dbConfig mysql.Config, dbType string, idleConn int, openConn int) (dbInfo *DbInfo, err error) {
	//Logger 설정
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})

	db, err := gorm.Open(mysql.New(dbConfig), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		loghelp.Panic(err.Error())
		panic(err)
	} else {
		loghelp.Info("Gorm DB Connection Completed")
	}
	sqlDB, _ := db.DB()

	//Connection Pool Set
	if idleConn == 0 {
		sqlDB.SetMaxIdleConns(10)
	} else if idleConn > 0 {
		sqlDB.SetMaxIdleConns(idleConn)
	} else {
		err := errors.New("wrong parameter :: idleConn must be positive")
		if err != nil {
			return nil, err
		}
	}

	if openConn == 0 {
		sqlDB.SetMaxOpenConns(100)
	} else if openConn > 0 {
		sqlDB.SetMaxOpenConns(openConn)
	} else {
		err := errors.New("wrong parameter :: openConn must be positive")
		if err != nil {
			return nil, err
		}
	}

	//Plugin Set
	_ = db.Use(&ZapLoggingPlugin{})
	loghelp.Info(fmt.Sprintf("%v", db.Config.Plugins["DBQueryLog"]))

	dbInfo = &DbInfo{
		GormDB: db,
		SqlDB:  sqlDB,
		DbType: dbType,
	}

	return
}

func CloseDB(sqlDB *sql.DB) {
	err := sqlDB.Close()
	if err != nil {
		panic(err)
	}
}
