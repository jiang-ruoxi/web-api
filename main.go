package main

import (
	"api/cache"
	"github.com/jiang-ruoxi/gopkg/es"
	"github.com/urfave/cli/v2"
	"os"

	"api/commands"
	rxMysql "github.com/jiang-ruoxi/gopkg/db"
	rxLog "github.com/jiang-ruoxi/gopkg/log"
	rxQueue "github.com/jiang-ruoxi/gopkg/queue"
	rxRedis "github.com/jiang-ruoxi/gopkg/redis"
	"github.com/jiang-ruoxi/gopkg/utils"
	"github.com/spf13/viper"
)

var configFile string

func main() {

	app := cli.NewApp()
	app.Action = commands.Serve
	app.Before = initConfig
	app.Commands = commands.Commands
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       "", // 默认从config目录读取
			Usage:       "specify the location of the configuration file",
			Required:    false,
			Destination: &configFile,
		},
	}

	if err := app.Run(os.Args); err != nil {
		rxLog.Sugar().Fatal(err)
	}
	rxLog.Flush()
}

func initConfig(*cli.Context) error {
	viper.SetDefault("app", "bilingual")
	if err := utils.LoadConfigInFile(configFile); err != nil {
		return err
	}
	rxLog.InitFromViper()
	rxMysql.InitMysqlDB()
	rxRedis.InitFromViperDefault()
	rxQueue.Initialize()
	cache.HttpCache()
	es.Initialize()
	return nil
}
