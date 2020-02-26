package main

import (
	"api-versioning-me/config"
	"api-versioning-me/db"
	"api-versioning-me/middlewares"
	"api-versioning-me/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var logger *zap.Logger


func init() {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "localhost")
	}
	if err := config.Init(os.Getenv("GO_ENV")); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("設定ファイルが存在しません。", err)
		} else {
			log.Fatal("設定ファイルが読み込めません。", err)
		}
	}
	if err := middlewares.LoggerInit(); err != nil {
		log.Fatal("ログの初期化に失敗しました。", err)
	}
	logger = middlewares.GetLogger()

	if err := db.Init(); err != nil {
		logger.Fatal("データベースの接続に失敗しました。" + err.Error())
	}
	logger.Info("データベースの接続に成功しました。")
}

func main() {
	server.Init()
}


