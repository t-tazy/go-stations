package main

import (
	"context"
	"log"
	"time"

	"github.com/TechBowl-japan/go-stations/config"
	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"github.com/TechBowl-japan/go-stations/server"
)

func main() {
	err := realMain(context.Background())
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain(ctx context.Context) error {
	// config values
	cfg := config.New()

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(cfg.DbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB, cfg)

	// TODO: サーバーをlistenする
	s := server.NewServer(cfg.Port, mux)
	return s.Run(ctx)
}
