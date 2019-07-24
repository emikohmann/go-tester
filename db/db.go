package db

import (
    "os"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var (
    Client *sql.DB
)

const (
    host          = "localhost"
    db            = "noop"
    driver        = "mysql"
    errConnecting = "error connecting to database"
)

func init() {
    connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", db, db, host, db)
    var err error
    if Client, err = sql.Open(driver, connection); err != nil {
        fmt.Println(errConnecting, err)
        os.Exit(0)
    }
}
