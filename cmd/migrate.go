package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	cloudsqlproxy "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"
	"github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	goauth "golang.org/x/oauth2/google"
)

var (
	CONNECT    string
	migrations = &migrate.FileMigrationSource{
		Dir: "migrations",
	}
)

func main() {

	credsFile := "/home/circleci/gcloud-service-key.json"
	SQLScope := "https://www.googleapis.com/auth/sqlservice.admin"
	ctx := context.Background()

	creds, err := ioutil.ReadFile(credsFile)
	if err != nil {

	}

	cfga, err := goauth.JWTConfigFromJSON(creds, SQLScope)
	if err != nil {

	}

	client := cfga.Client(ctx)
	proxy.Init(client, nil, nil)

	db, err := cloudsqlproxy.DialCfg(&mysql.Config{
		Addr:                 os.Getenv("DB_HOST"),               // インスタンス接続名
		DBName:               "sleepdays_" + os.Getenv("GO_ENV"), // データベース名
		User:                 "root",                             // ユーザ名
		Passwd:               os.Getenv("MYSQL_ROOT_PASSWORD"),   // ユーザパスワード
		Net:                  "cloudsql",                         // Cloud SQL Proxy で接続する場合は cloudsql 固定です
		ParseTime:            true,                               // DATE/DATETIME 型を time.Time へパースする
		TLSConfig:            "",                                 // TLSConfig は空文字を設定しなければなりません
		AllowNativePasswords: true,
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	appliedCount, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}
	log.Printf("Applied %v migrations", appliedCount)
}
