package main

import (
	"fmt"
	"godis/cluster"
	"godis/config"
	"godis/database"
	databaseface "godis/interface/database"
	"godis/lib/logger"
	"godis/resp/handler"
	"godis/tcp"
	"os"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// *2\r\n$3\r\nGET\r\n$3\r\nKEY\r\n
// *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
	var db databaseface.Database
	if config.Properties.Self != "" &&
		len(config.Properties.Peers) > 0 {
		db = cluster.NewClusterDatabase()
	} else {
		db = database.NewStandaloneDatabase()
	}
	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d",
				config.Properties.Bind,
				config.Properties.Port),
		},
		handler.NewRespHandler(db))
	if err != nil {
		logger.Error(err)
	}
}
