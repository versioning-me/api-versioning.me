package db

import (
	"api-versioning-me/config"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	Db  *gorm.DB
	Con string
)

const DBMS = "mysql"

func Set() {
	DbName := "versioningme_" + config.Env

	switch {
	case config.Env == "development" || config.Env == "staging" || config.Env == "production":
		user := os.Getenv("MYSQL_USER_NAME")
		pass := os.Getenv("MYSQL_USER_PASS")

		Con = fmt.Sprintf("%s:%s@unix(%s/%s)/%s?charset=utf8&parseTime=true&loc=Local", user, pass, "/cloudsql", os.Getenv("DB_HOST"), DbName)

	default:
		user := "user"
		pass := "pass"
		Con = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", user, pass, "localhost", 3307, DbName)
	}
}

func Init() (err error) {
	Set()
	Db, err = gorm.Open(DBMS, Con)
	if err != nil {
		return err
	}
	Db.SingularTable(true)
	return nil
}
