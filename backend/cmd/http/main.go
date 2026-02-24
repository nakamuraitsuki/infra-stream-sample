package main

import "log"

func main() {
	app, err := InitializeHTTPServer()
	if err != nil {
		log.Fatalf("サーバーの初期化に失敗しました: %v", err)
	}
	defer app.DB.Close()

	if err := app.Echo.Start(":8080"); err != nil {
		app.Echo.Logger.Errorf("サーバーの起動に失敗しました: %v", err)
	}
}
