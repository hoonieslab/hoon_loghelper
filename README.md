# hoon_loghelper

make that easier using zap and registering log plugin using gorm.

## Installation

`go get -u github.com/hoonieslab/hoon_loghelper`

Note that package made using Go 1.18.

## Quick Start

```go
// First, set path to save *.log file and initialize loghelp
loghelp.SetLogFilePath("./logs/test-%Y-%m-%d-%H.log")
loghelp.Init()

// Registering logging query plugin for gorm
dbConfig := mysql.Config{
    DSN: "USERNAME:PW@tcp(IP:PORT)/DB?charset=utf8mb4&parseTime=True&loc=UTC",
}
dbhelp.InitGormCon(dbConfig, "MariaDB", 10, 100)
```
