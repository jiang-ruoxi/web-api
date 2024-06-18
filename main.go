package main

import (
	"os"

	rxMysql "github.com/jiang-ruoxi/gopkg/db"
	rxLog "github.com/jiang-ruoxi/gopkg/log"
	rxRedis "github.com/jiang-ruoxi/gopkg/redis"
	"github.com/jiang-ruoxi/gopkg/utils"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"api/commands"
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
	rxRedis.InitFromViper()
	return nil
}
