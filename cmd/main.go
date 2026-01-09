package main

import "github.com/zchelalo/expense-control-back/pkg/bootstrap"

func main() {
	_, err := bootstrap.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	log := bootstrap.GetLogger()
	defer bootstrap.SyncLogger()

	log.Info("application starting")
}